package network

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/pkg/stdcopy"

	"github.com/ethereum/go-ethereum/log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/ethclient/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/config"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/enclaverunner"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"
	"github.com/obscuronet/obscuro-playground/integration/erc20contract"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"

	"github.com/obscuronet/obscuro-playground/go/ethclient"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/integration/simulation/stats"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
)

const (
	enclaveDockerImg  = "obscuro_enclave"
	enclaveAddress    = ":11000"
	enclaveDockerPort = "11000/tcp"
)

// creates Obscuro nodes with their own enclave servers that communicate with peers via sockets, wires them up, and populates the network objects
type basicNetworkOfNodesWithDockerEnclave struct {
	obscuroClients   []*obscuroclient.Client
	gethNetwork      *gethnetwork.GethNetwork
	wallets          []wallet.Wallet
	contracts        []string
	workerWallet     wallet.Wallet
	ctx              context.Context
	client           *client.Client
	containerIDs     map[string]string
	containerStreams map[string]*types.HijackedResponse
}

func NewBasicNetworkOfNodesWithDockerEnclave(wallets []wallet.Wallet, workerWallet wallet.Wallet, contracts []string) Network {
	return &basicNetworkOfNodesWithDockerEnclave{
		wallets:          wallets,
		contracts:        contracts,
		workerWallet:     workerWallet,
		containerStreams: map[string]*types.HijackedResponse{},
	}
}

// Create initializes Obscuro nodes with their own Dockerised enclave servers that communicate with peers via sockets, wires them up, and populates the network objects
// TODO - Use individual Docker containers for the Obscuro nodes and Ethereum nodes.
func (n *basicNetworkOfNodesWithDockerEnclave) Create(params *params.SimParams, stats *stats.Stats) ([]ethclient.EthClient, []*obscuroclient.Client, []string, error) {
	// We create a Docker client.
	n.ctx = context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}
	n.client = cli
	// We check the required Docker images are available.
	if !dockerImagesAvailable(n.ctx, cli) {
		// We don't cause the test to fail here, because we want users to be able to run all the tests in the repo
		// without having to build the Docker images.
		return nil, nil, nil, fmt.Errorf("this test requires the `%s` Docker image to be built using `dockerfiles/enclave.Dockerfile`. Terminating", enclaveDockerImg)
	}

	// make sure the geth network binaries exist
	path, err := gethnetwork.EnsureBinariesExist(gethnetwork.LatestVersion)
	if err != nil {
		panic(err)
	}

	// get wallet addresses to prefund them
	walletAddresses := make([]string, len(n.wallets))
	for i, w := range n.wallets {
		walletAddresses[i] = w.Address().String()
	}

	// kickoff the network with the prefunded wallet addresses
	n.gethNetwork = gethnetwork.NewGethNetwork(
		params.StartPort,
		params.StartPort+DefaultWsPortOffset,
		path,
		params.NumberOfNodes,
		int(params.AvgBlockDuration.Seconds()),
		walletAddresses,
	)

	tmpHostConfig := config.HostConfig{
		L1NodeHost:          Localhost,
		L1NodeWebsocketPort: n.gethNetwork.WebSocketPorts[0],
		L1ConnectionTimeout: DefaultL1ConnectionTimeout,
	}
	tmpEthClient, err := ethclient.NewEthClient(tmpHostConfig)
	if err != nil {
		panic(err)
	}

	mgmtContractBlkHash, mgmtContractAddr, err := DeployContract(tmpEthClient, n.workerWallet, common.Hex2Bytes(mgmtcontractlib.MgmtContractByteCode))
	if err != nil {
		panic(fmt.Sprintf("failed to deploy management contract. Cause: %s", err))
	}
	_, erc20ContractAddr, err := DeployContract(tmpEthClient, n.workerWallet, common.Hex2Bytes(erc20contract.ContractByteCode))
	if err != nil {
		panic(fmt.Sprintf("failed to deploy ERC20 contract. Cause: %s", err))
	}

	// We create the Docker containers and set up a hook to terminate them at the end of the test.
	containerIDs := createDockerContainers(n.ctx, cli, params.NumberOfNodes, params.StartPort, mgmtContractAddr.Hex(), erc20ContractAddr.Hex())
	n.containerIDs = containerIDs

	// We start the Docker containers.
	for id := range containerIDs {
		if err = cli.ContainerStart(n.ctx, id, types.ContainerStartOptions{}); err != nil {
			panic(err)
		}
		waiter, err := cli.ContainerAttach(n.ctx, id, types.ContainerAttachOptions{
			Stderr: true,
			Stdout: true,
			Stdin:  false,
			Stream: true,
		})

		go func() {
			_, err := stdcopy.StdCopy(os.Stdout, os.Stderr, waiter.Reader)
			if err != nil {
				log.Error("Could not copy output from the docker container")
			}
		}()

		if err != nil {
			panic(err)
		}
		n.containerStreams[id] = &waiter
	}

	params.MgmtContractAddr = mgmtContractAddr
	params.MgmtContractBlkHash = mgmtContractBlkHash
	params.StableTokenContractAddr = erc20ContractAddr
	params.MgmtContractLib = mgmtcontractlib.NewMgmtContractLib(mgmtContractAddr)
	params.ERC20ContractLib = erc20contractlib.NewERC20ContractLib(mgmtContractAddr, erc20ContractAddr)

	l1Clients := make([]ethclient.EthClient, params.NumberOfNodes)
	obscuroNodes := make([]*host.Node, params.NumberOfNodes)
	n.obscuroClients = make([]*obscuroclient.Client, params.NumberOfNodes)
	nodeP2pAddrs := make([]string, params.NumberOfNodes)

	for i := 0; i < params.NumberOfNodes; i++ {
		// We assign a P2P address to each node on the network.
		nodeP2pAddrs[i] = fmt.Sprintf("%s:%d", Localhost, params.StartPort+DefaultHostP2pOffset+i)
	}

	for i := 0; i < params.NumberOfNodes; i++ {
		isGenesis := i == 0

		// create a remote enclave server
		enclavePort := uint64(params.StartPort + DefaultEnclaveOffset + i)
		// create the L1 client, the Obscuro host and enclave service, and the L2 client
		l1Client := createEthClientConnection(
			int64(i),
			n.gethNetwork.WebSocketPorts[i],
		)
		rpcAddress := fmt.Sprintf("%s:%d", Localhost, params.StartPort+DefaultHostRPCOffset+i)
		agg := createSocketObscuroNode(
			int64(i+1),
			isGenesis,
			params.AvgGossipPeriod,
			stats,
			nodeP2pAddrs[i],
			nodeP2pAddrs,
			fmt.Sprintf("%s:%d", Localhost, enclavePort),
			rpcAddress,
			params.NodeEthWallets[i],
			params.MgmtContractLib,
			params.MgmtContractBlkHash,
		)
		obscuroClient := obscuroclient.NewClient(rpcAddress)

		// connect the L1 and L2 nodes
		agg.ConnectToEthNode(l1Client)

		obscuroNodes[i] = agg
		n.obscuroClients[i] = &obscuroClient
		l1Clients[i] = l1Client
	}

	// start each obscuro node
	for _, m := range obscuroNodes {
		t := m
		go t.Start()
		time.Sleep(params.AvgBlockDuration / 3)
	}

	return l1Clients, n.obscuroClients, nodeP2pAddrs, nil
}

func (n *basicNetworkOfNodesWithDockerEnclave) TearDown() {
	n.gethNetwork.StopNodes()
	for _, c := range n.obscuroClients {
		temp := c
		go func() {
			defer (*temp).Stop()
			err := (*temp).Call(nil, obscuroclient.RPCStopHost)
			if err != nil {
				log.Error("Failed to stop node: %s", err)
			}
		}()
	}
	terminateDockerContainers(n.ctx, n.client, n.containerIDs, n.containerStreams)
}

// Checks the required Docker images exist.
func dockerImagesAvailable(ctx context.Context, cli *client.Client) bool {
	images, _ := cli.ImageList(ctx, types.ImageListOptions{})
	for _, image := range images {
		for _, tag := range image.RepoTags {
			if strings.Contains(tag, enclaveDockerImg) {
				return true
			}
		}
	}
	return false
}

// Creates the test Docker containers.
func createDockerContainers(ctx context.Context, client *client.Client, numOfNodes int, startPort int, mngmtCtrAddr string, erc20Addr string) map[string]string {
	var enclavePorts []string
	for i := 0; i < numOfNodes; i++ {
		// We assign an enclave port to each enclave service on the network.
		enclavePorts = append(enclavePorts, fmt.Sprintf("%d", startPort+DefaultEnclaveOffset+i))
	}

	containerIDs := map[string]string{}
	for idx, port := range enclavePorts {
		nodeID := common.BigToAddress(big.NewInt(int64(idx + 1))).Hex()
		containerConfig := &container.Config{
			Image: enclaveDockerImg,
			Cmd: []string{
				"--" + enclaverunner.HostIDName, nodeID,
				"--" + enclaverunner.AddressName, enclaveAddress,
				"--" + enclaverunner.ManagementContractAddressName, mngmtCtrAddr,
				"--" + enclaverunner.Erc20ContractAddrsName, erc20Addr,
			},
		}
		r := container.Resources{
			Memory:     2 * 1024 * 1024 * 1024, // 2GB
			MemorySwap: -1,
		}
		hostConfig := &container.HostConfig{
			PortBindings: nat.PortMap{nat.Port(enclaveDockerPort): []nat.PortBinding{{HostIP: Localhost, HostPort: port}}},
			Resources:    r,
		}

		resp, err := client.ContainerCreate(ctx, containerConfig, hostConfig, nil, "")
		if err != nil {
			panic(err)
		}
		containerIDs[resp.ID] = port
	}

	return containerIDs
}

// Stops and removes the test Docker containers.
func terminateDockerContainers(ctx context.Context, cli *client.Client, containerIDs map[string]string, containerStreams map[string]*types.HijackedResponse) {
	for id := range containerIDs {
		if containerStreams[id] != nil {
			containerStreams[id].Close()
		}
		err1 := cli.ContainerStop(ctx, id, nil)
		if err1 != nil {
			fmt.Printf("Could not stop the container %v : %s\n", id, err1)
			continue
		}

		err2 := cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{
			RemoveVolumes: true,
			RemoveLinks:   false,
			Force:         true,
		})
		if err2 != nil {
			fmt.Printf("Could not remove the container %v : %s\n", id, err2)
			continue
		}

		fmt.Printf("Stopped and removed container %v\n", id)
	}

	if err := cli.Close(); err != nil {
		fmt.Printf("Could not close cli: %s\n", err)
	}
}
