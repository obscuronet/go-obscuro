package network

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/contracts"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/ethclient/wallet"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"
	"github.com/obscuronet/obscuro-playground/integration/simulation/p2p"
	"github.com/obscuronet/obscuro-playground/integration/simulation/params"
	"github.com/obscuronet/obscuro-playground/integration/simulation/stats"
)

type networkInMemGeth struct {
	obscuroNodes   []*host.Node
	obscuroClients []*obscuroclient.Client
	gethNetwork    *gethnetwork.GethNetwork
}

func NewNetworkInMemoryGeth() Network {
	return &networkInMemGeth{}
}

// Create inits and starts the nodes, wires them up, and populates the network objects
func (n *networkInMemGeth) Create(params *params.SimParams, stats *stats.Stats) ([]ethclient.EthClient, []*obscuroclient.Client, []string) {
	// make sure the geth network binaries exist
	path, err := gethnetwork.EnsureBinariesExist(gethnetwork.LatestVersion)
	if err != nil {
		panic(err)
	}

	// convert the wallets to strings
	walletAddresses := make([]string, params.NumberOfObscuroWallets)
	for i := 0; i < params.NumberOfObscuroWallets; i++ {
		walletAddresses[i] = params.EthWallets[i].Address().String()
	}

	// kickoff the network with the prefunded wallet addresses
	gn := gethnetwork.NewGethNetwork(
		40000,
		path,
		params.NumberOfNodes,
		int(params.AvgBlockDuration.Seconds()),
		walletAddresses,
	)
	n.gethNetwork = &gn
	// take the first random wallet and deploy the contract in the network
	contractAddr := deployContract(params.EthWallets[0], gn.WebSocketPorts[0])

	params.MgmtContractAddr = contractAddr
	params.TxHandler = mgmtcontractlib.NewEthMgmtContractTxHandler(contractAddr)

	// Create the obscuro node, each connected to a geth node
	l1Clients := make([]ethclient.EthClient, params.NumberOfNodes)
	n.obscuroNodes = make([]*host.Node, params.NumberOfNodes)
	n.obscuroClients = make([]*obscuroclient.Client, params.NumberOfNodes)

	for i := 0; i < params.NumberOfNodes; i++ {
		isGenesis := i == 0

		// create the in memory l1 and l2 node
		miner := createEthClientConnection(
			int64(i),
			n.gethNetwork.WebSocketPorts[i],
			params.EthWallets[i],
			params.MgmtContractAddr,
		)
		agg := createInMemObscuroNode(
			int64(i),
			isGenesis,
			params.TxHandler,
			params.AvgGossipPeriod,
			params.AvgBlockDuration,
			params.AvgNetworkLatency,
			stats,
			true,
			n.gethNetwork.GenesisJSON,
		)
		obscuroClient := host.NewInMemObscuroClient(int64(i), &agg.P2p, agg.DB(), &agg.EnclaveClient)

		// and connect them to each other
		agg.ConnectToEthNode(miner)

		n.obscuroNodes[i] = agg
		n.obscuroClients[i] = &obscuroClient
		l1Clients[i] = miner
	}

	// make sure the aggregators can talk to each other
	for i := 0; i < params.NumberOfNodes; i++ {
		n.obscuroNodes[i].P2p.(*p2p.MockP2P).Nodes = n.obscuroNodes
	}

	// start each obscuro node
	for _, m := range n.obscuroNodes {
		t := m
		go t.Start()
		time.Sleep(params.AvgBlockDuration / 10)
	}

	return l1Clients, n.obscuroClients, nil
}

func (n *networkInMemGeth) TearDown() {
	for _, client := range n.obscuroClients {
		temp := client
		go (*temp).Stop()
	}

	for _, node := range n.obscuroNodes {
		temp := node
		go temp.Stop()
	}

	n.gethNetwork.StopNodes()
}

func createEthClientConnection(id int64, port uint, wallet wallet.Wallet, contractAddr common.Address) ethclient.EthClient {
	ethnode, err := ethclient.NewEthClient(common.BigToAddress(big.NewInt(id)), "127.0.0.1", port, wallet, contractAddr)
	if err != nil {
		panic(err)
	}
	return ethnode
}

func deployContract(w wallet.Wallet, port uint) common.Address {
	tmpClient, err := ethclient.NewEthClient(common.Address{}, "127.0.0.1", port, w, common.Address{})
	if err != nil {
		panic(err)
	}

	deployContractTx := types.LegacyTx{
		Nonce:    0, // relies on a clean env
		GasPrice: big.NewInt(2000000000),
		Gas:      1025_000_000,
		Data:     common.Hex2Bytes(contracts.MgmtContractByteCode),
	}

	signedTx, err := tmpClient.SubmitTransaction(&deployContractTx)
	if err != nil {
		panic(err)
	}

	var receipt *types.Receipt
	for start := time.Now(); time.Since(start) < 80*time.Second; time.Sleep(2 * time.Second) {
		receipt, err = tmpClient.FetchTxReceipt(signedTx.Hash())
		if err == nil && receipt != nil {
			break
		}
		if !errors.Is(err, ethereum.NotFound) {
			panic(err)
		}
		fmt.Printf("Contract deploy tx has not been mined into a block after %s...\n", time.Since(start))
	}

	fmt.Printf("Contract deployed to %s - using port %d\n", receipt.ContractAddress, port)
	return receipt.ContractAddress
}
