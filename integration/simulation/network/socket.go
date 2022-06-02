package network

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/config"

	"github.com/obscuronet/obscuro-playground/go/ethclient/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"
	"github.com/obscuronet/obscuro-playground/integration/erc20contract"

	"github.com/obscuronet/obscuro-playground/integration"

	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"

	"github.com/obscuronet/obscuro-playground/go/ethclient"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/integration/simulation/stats"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
)

// creates Obscuro nodes with their own enclave servers that communicate with peers via sockets, wires them up, and populates the network objects
type networkOfSocketNodes struct {
	obscuroClients []*obscuroclient.Client
	gethNetwork    *gethnetwork.GethNetwork
	wallets        []wallet.Wallet
	contracts      []string
	workerWallet   wallet.Wallet
}

func NewNetworkOfSocketNodes(wallets []wallet.Wallet, workerWallet wallet.Wallet, contracts []string) Network {
	return &networkOfSocketNodes{
		wallets:      wallets,
		contracts:    contracts,
		workerWallet: workerWallet,
	}
}

func (n *networkOfSocketNodes) Create(params *params.SimParams, stats *stats.Stats) ([]ethclient.EthClient, []*obscuroclient.Client, []string, error) {
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
		nodeP2pAddrs[i] = fmt.Sprintf("%s:%d", Localhost, params.StartPort+200+i)
	}

	for i := 0; i < params.NumberOfNodes; i++ {
		isGenesis := i == 0

		// create a remote enclave server
		enclaveAddr := fmt.Sprintf("%s:%d", Localhost, params.StartPort+DefaultEnclaveOffset+i)
		enclaveConfig := config.EnclaveConfig{
			HostID:           common.BigToAddress(big.NewInt(int64(i))),
			Address:          enclaveAddr,
			L1ChainID:        integration.EthereumChainID,
			ObscuroChainID:   integration.ObscuroChainID,
			ValidateL1Blocks: false,
			WillAttest:       false,
			GenesisJSON:      nil,
			UseInMemoryDB:    false,
		}
		_, err := enclave.StartServer(enclaveConfig, params.MgmtContractLib, params.ERC20ContractLib, stats)
		if err != nil {
			panic(fmt.Sprintf("failed to create enclave server: %v", err))
		}

		// create the L1 client, the Obscuro host and enclave service, and the L2 client
		l1Client := createEthClientConnection(
			int64(i),
			n.gethNetwork.WebSocketPorts[i],
		)
		obscuroClientAddr := fmt.Sprintf("%s:%d", Localhost, params.StartPort+400+i)
		obscuroClient := obscuroclient.NewClient(obscuroClientAddr)
		agg := createSocketObscuroNode(
			int64(i),
			isGenesis,
			params.AvgGossipPeriod,
			stats,
			nodeP2pAddrs[i],
			nodeP2pAddrs,
			enclaveAddr,
			obscuroClientAddr,
			params.NodeEthWallets[i],
			params.MgmtContractLib,
			params.MgmtContractBlkHash,
		)

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

func (n *networkOfSocketNodes) TearDown() {
	defer n.gethNetwork.StopNodes()

	for _, client := range n.obscuroClients {
		temp := client
		go func() {
			defer (*temp).Stop()
			err := (*temp).Call(nil, obscuroclient.RPCStopHost)
			if err != nil {
				log.Error("Failed to stop client %s", err)
			}
		}()
	}
}
