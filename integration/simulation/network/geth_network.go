package network

import (
	"math/big"
	"time"

	"github.com/obscuronet/obscuro-playground/go/ethclient/wallet"

	"github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	"github.com/obscuronet/obscuro-playground/integration/simulation/p2p"
	"github.com/obscuronet/obscuro-playground/integration/simulation/params"
	"github.com/obscuronet/obscuro-playground/integration/simulation/stats"
)

type networkInMemGeth struct {
	obscuroNodes     []*host.Node
	obscuroAddresses []string
}

func NewNetworkInMemoryGeth() Network {
	return &networkInMemGeth{}
}

// Create inits and starts the nodes, wires them up, and populates the network objects
func (n *networkInMemGeth) Create(params params.SimParams, stats *stats.Stats) ([]ethclient.EthClient, []*host.Node, []string) {
	l1Clients := make([]ethclient.EthClient, params.NumberOfNodes)
	n.obscuroNodes = make([]*host.Node, params.NumberOfNodes)

	// all nodes use the same wallet for now
	// TODO create a wallet loading key mechanism for each node attached to a prefunded network genesis
	wallets := []wallet.Wallet{
		wallet.NewInMemoryWallet("5dbbff1b5ff19f1ad6ea656433be35f6846e890b3f3ec6ef2b2e2137a8cab4ae"),
		wallet.NewInMemoryWallet("b728cd9a9f54cede03a82fc189eab4830a612703d48b7ef43ceed2cbad1a06c7"),
		wallet.NewInMemoryWallet("1e1e76d5c0ea1382b6acf76e873977fd223c7fa2a6dc57db2b94e93eb303ba85"),
	}
	for i := 0; i < params.NumberOfNodes; i++ {
		genesis := false
		if i == 0 {
			genesis = true
		}

		// create the in memory l1 and l2 node
		miner := createRealEthNode(int64(i), wallets[i], params.MgmtContractAddr)
		agg := createInMemObscuroNode(int64(i), genesis, params.TxHandler, params.AvgGossipPeriod, params.AvgBlockDuration, params.AvgNetworkLatency, stats)

		// and connect them to each other
		agg.ConnectToEthNode(miner)

		n.obscuroNodes[i] = agg
		l1Clients[i] = miner
	}

	// populate the nodes field of each network
	for i := 0; i < params.NumberOfNodes; i++ {
		n.obscuroNodes[i].P2p.(*p2p.MockP2P).Nodes = n.obscuroNodes
	}

	n.obscuroAddresses = nil

	for _, m := range n.obscuroNodes {
		t := m
		go t.Start()
		time.Sleep(params.AvgBlockDuration / 3)
	}

	return l1Clients, n.obscuroNodes, nil
}

func (n *networkInMemGeth) TearDown() {
	// Nop
}

func createRealEthNode(id int64, wallet wallet.Wallet, contractAddr common.Address) ethclient.EthClient {
	ethnode, err := ethclient.NewEthClient(common.BigToAddress(big.NewInt(id)), "127.0.0.1", 7545, wallet, contractAddr)
	if err != nil {
		panic(err)
	}
	return ethnode
}
