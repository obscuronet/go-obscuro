package network

import (
	"fmt"
	"math/big"
	"time"

	"github.com/obscuronet/obscuro-playground/go/ethclient/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"
	"github.com/obscuronet/obscuro-playground/integration/erc20contract"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/config"

	"github.com/obscuronet/obscuro-playground/integration"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"

	"github.com/obscuronet/obscuro-playground/go/ethclient"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/integration/simulation/stats"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
)

const enclavePort = ":11000"

// creates Obscuro nodes with their own enclave servers that communicate with peers via sockets, wires them up, and populates the network objects
type networkWithAzureEnclaves struct {
	gethNetwork  *gethnetwork.GethNetwork
	wallets      []wallet.Wallet
	contracts    []string
	workerWallet wallet.Wallet

	obscuroClients   []*obscuroclient.Client
	enclaveAddresses []string
}

func NewNetworkWithAzureEnclaves(enclaveAddresses []string, wallets []wallet.Wallet, workerWallet wallet.Wallet, contracts []string) Network {
	if len(enclaveAddresses) == 0 {
		panic("Cannot create azure enclaves network without at least one enclave address.")
	}
	return &networkWithAzureEnclaves{
		enclaveAddresses: enclaveAddresses,
		wallets:          wallets,
		contracts:        contracts,
		workerWallet:     workerWallet,
	}
}

func (n *networkWithAzureEnclaves) Create(params *params.SimParams, stats *stats.Stats) ([]ethclient.EthClient, []*obscuroclient.Client, []string, error) {
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

	mgmtContractAddr, err := DeployContract(tmpEthClient, n.workerWallet, common.Hex2Bytes(mgmtcontractlib.MgmtContractByteCode))
	if err != nil {
		panic(fmt.Sprintf("failed to deploy management contract. Cause: %s", err))
	}
	erc20ContractAddr, err := DeployContract(tmpEthClient, n.workerWallet, common.Hex2Bytes(erc20contract.ContractByteCode))
	if err != nil {
		panic(fmt.Sprintf("failed to deploy ERC20 contract. Cause: %s", err))
	}
	fmt.Printf("Please start the docker image on the azure server with with:\n")
	for i := 0; i < len(n.enclaveAddresses); i++ {
		fmt.Printf("sudo docker run -e OE_SIMULATION=0 --privileged -v /dev/sgx:/dev/sgx -p 11000:11000/tcp obscuro_enclave --hostID %d --address :11000 --managementContractAddress %s  --erc20ContractAddresses %s\n", i, mgmtContractAddr.Hex(), erc20ContractAddr.Hex())
	}
	time.Sleep(10 * time.Second)

	params.MgmtContractAddr = mgmtContractAddr
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

	// set up nodes with azure enclave
	for i := 0; i < len(n.enclaveAddresses); i++ {
		isGenesis := i == 0
		// create the L1 client, the Obscuro host and enclave service, and the L2 client
		l1Client := createEthClientConnection(
			int64(i),
			n.gethNetwork.WebSocketPorts[i],
		)
		rpcAddress := fmt.Sprintf("%s:%d", Localhost, params.StartPort+DefaultHostRPCOffset+i)
		agg := createSocketObscuroNode(
			int64(i),
			isGenesis,
			params.AvgGossipPeriod,
			stats,
			nodeP2pAddrs[i],
			nodeP2pAddrs,
			n.enclaveAddresses[i]+enclavePort,
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

	// set up nodes with mock enclaves
	for i := len(n.enclaveAddresses); i < params.NumberOfNodes; i++ {
		l1Client := createEthClientConnection(
			int64(i),
			n.gethNetwork.WebSocketPorts[i],
		)

		// create a remote enclave server
		nonAzureEnclaveAddress := fmt.Sprintf("%s:%d", Localhost, params.StartPort+DefaultEnclaveOffset+i)
		enclaveConfig := config.EnclaveConfig{
			HostID:           common.BigToAddress(big.NewInt(int64(i))),
			Address:          nonAzureEnclaveAddress,
			L1ChainID:        integration.EthereumChainID,
			ObscuroChainID:   integration.ObscuroChainID,
			WillAttest:       false,
			ValidateL1Blocks: false,
			GenesisJSON:      nil,
			UseInMemoryDB:    true,
		}
		_, err := enclave.StartServer(
			enclaveConfig,
			params.MgmtContractLib,
			params.ERC20ContractLib,
			stats,
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create enclave server: %v", err))
		}

		// create the in memory l1 and l2 node
		rpcAddress := fmt.Sprintf("%s:%d", Localhost, params.StartPort+DefaultHostRPCOffset+i)
		agg := createSocketObscuroNode(
			int64(i),
			false,
			params.AvgGossipPeriod,
			stats,
			nodeP2pAddrs[i],
			nodeP2pAddrs,
			nonAzureEnclaveAddress,
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

	for _, m := range obscuroNodes {
		t := m
		go t.Start()
		time.Sleep(params.AvgBlockDuration / 3)
	}

	return l1Clients, n.obscuroClients, nodeP2pAddrs, nil
}

func (n *networkWithAzureEnclaves) TearDown() {
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
}
