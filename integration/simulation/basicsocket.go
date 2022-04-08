package simulation

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"
)

// creates Obscuro nodes with their own enclave servers that communicate with peers via sockets, wires them up, and populates the network objects
type basicNetworkOfSocketNodes struct {
	ethNodes         []*ethereum_mock.Node
	obscuroNodes     []*host.Node
	obscuroAddresses []string
}

func NewBasicNetworkOfSocketNodes() Network {
	return &basicNetworkOfSocketNodes{}
}

func (n *basicNetworkOfSocketNodes) Create(params SimParams, stats *Stats) ([]*ethereum_mock.Node, []*host.Node, []string) {
	// todo - add observer nodes
	l1Nodes := make([]*ethereum_mock.Node, params.NumberOfNodes)
	l2Nodes := make([]*host.Node, params.NumberOfNodes)

	var nodeP2pAddrs []string
	for i := 0; i < params.NumberOfNodes; i++ {
		// We assign a P2P address to each node on the network.
		nodeP2pAddrs = append(nodeP2pAddrs, fmt.Sprintf("%s:%d", Localhost, p2pStartPort+i))
	}

	for i := 1; i <= params.NumberOfNodes; i++ {
		genesis := false
		if i == 1 {
			genesis = true
		}

		// create a remote enclave server
		nodeID := common.BigToAddress(big.NewInt(int64(i)))
		enclavePort := uint64(EnclaveStartPort + i)
		enclaveAddress := fmt.Sprintf("localhost:%d", enclavePort)
		err := enclave.StartServer(enclaveAddress, nodeID, stats)
		if err != nil {
			panic(fmt.Sprintf("failed to create enclave server: %v", err))
		}

		// create the in memory l1 and l2 node
		miner := createMockEthNode(int64(i), params.NumberOfNodes, params.AvgBlockDurationUSecs, params.AvgNetworkLatency, stats)
		agg := createSocketObscuroNode(int64(i), genesis, params.AvgGossipPeriod, stats, nodeP2pAddrs[i-1], nodeP2pAddrs, enclavePort)

		// and connect them to each other
		agg.ConnectToEthNode(miner)
		miner.AddClient(agg)

		l1Nodes[i-1] = miner
		l2Nodes[i-1] = agg
	}

	// populate the nodes field of the L1 network
	for i := 0; i < params.NumberOfNodes; i++ {
		l1Nodes[i].Network.(*ethereum_mock.MockEthNetwork).AllNodes = l1Nodes
	}

	n.ethNodes = l1Nodes
	n.obscuroNodes = l2Nodes
	n.obscuroAddresses = nodeP2pAddrs

	// The sequence of starting the nodes is important to catch various edge cases.
	// Here we first start the mock layer 1 nodes, with a pause between them of a fraction of a block duration.
	// The reason is to make sure that they catch up correctly.
	// Then we pause for a while, to give the L1 network enough time to create a number of blocks, which will have to be ingested by the Obscuro nodes
	// Then, we begin the starting sequence of the Obscuro nodes, again with a delay between them, to test that they are able to cach up correctly.
	// Note: Other simulations might test variations of this pattern.
	for _, m := range n.ethNodes {
		t := m
		go t.Start()
		time.Sleep(time.Duration(params.AvgBlockDurationUSecs / 8))
	}

	time.Sleep(time.Duration(params.AvgBlockDurationUSecs * 20))
	for _, m := range n.obscuroNodes {
		t := m
		go t.Start()
		time.Sleep(time.Duration(params.AvgBlockDurationUSecs / 3))
	}

	return l1Nodes, l2Nodes, nodeP2pAddrs
}

func (n *basicNetworkOfSocketNodes) TearDown() {
	go func() {
		for _, n := range n.obscuroNodes {
			n.Stop()
		}
	}()
	go func() {
		for _, m := range n.ethNodes {
			t := m
			go t.Stop()
		}
	}()
}