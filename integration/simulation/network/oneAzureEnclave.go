package network

import (
	"fmt"
	"math/big"
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"

	"github.com/obscuronet/obscuro-playground/go/ethclient"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/integration/simulation/stats"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"
)

// creates Obscuro nodes with their own enclave servers that communicate with peers via sockets, wires them up, and populates the network objects
type networkWithOneAzureEnclave struct {
	ethNodes       []*ethereum_mock.Node
	obscuroNodes   []*host.Node
	obscuroClients []*obscuroclient.Client
	enclaveAddress string
}

func NewNetworkWithOneAzureEnclave(enclaveAddress string) Network {
	return &networkWithOneAzureEnclave{enclaveAddress: enclaveAddress}
}

func (n *networkWithOneAzureEnclave) Create(params *params.SimParams, stats *stats.Stats) ([]ethclient.EthClient, []*obscuroclient.Client, []string) {
	l1Clients := make([]ethclient.EthClient, params.NumberOfNodes)
	n.ethNodes = make([]*ethereum_mock.Node, params.NumberOfNodes)
	n.obscuroNodes = make([]*host.Node, params.NumberOfNodes)
	n.obscuroClients = make([]*obscuroclient.Client, params.NumberOfNodes)
	nodeP2pAddrs := make([]string, params.NumberOfNodes)

	for i := 0; i < params.NumberOfNodes; i++ {
		// We assign a P2P address to each node on the network.
		nodeP2pAddrs[i] = fmt.Sprintf("%s:%d", Localhost, params.StartPort+i)
	}

	for i := 0; i < params.NumberOfNodes; i++ {
		isGenesis := i == 0

		if isGenesis {
			// create the in memory l1 and l2 node
			miner := createMockEthNode(int64(i), params.NumberOfNodes, params.AvgBlockDuration, params.AvgNetworkLatency, stats)
			obscuroClientAddr := fmt.Sprintf("%s:%d", Localhost, params.StartPort+200+i)
			obscuroClient := obscuroclient.NewClient(obscuroClientAddr)
			agg := createSocketObscuroNode(int64(i), isGenesis, params.AvgGossipPeriod, stats, nodeP2pAddrs[i], nodeP2pAddrs, n.enclaveAddress, obscuroClientAddr)

			// and connect them to each other
			agg.ConnectToEthNode(miner)
			miner.AddClient(agg)

			l1Clients[i] = miner
			n.ethNodes[i] = miner
			n.obscuroNodes[i] = agg
			n.obscuroClients[i] = &obscuroClient
		} else {
			// create a remote enclave server
			nodeID := common.BigToAddress(big.NewInt(int64(i)))
			enclavePort := uint64(params.StartPort + 100 + i)
			enclaveAddress := fmt.Sprintf("localhost:%d", enclavePort)
			_, err := enclave.StartServer(enclaveAddress, nodeID, params.TxHandler, false, nil, stats)
			if err != nil {
				panic(fmt.Sprintf("failed to create enclave server: %v", err))
			}

			// create the in memory l1 and l2 node
			miner := createMockEthNode(int64(i), params.NumberOfNodes, params.AvgBlockDuration, params.AvgNetworkLatency, stats)
			obscuroClientAddr := fmt.Sprintf("%s:%d", Localhost, params.StartPort+200+i)
			obscuroClient := obscuroclient.NewClient(obscuroClientAddr)
			agg := createSocketObscuroNode(int64(i), isGenesis, params.AvgGossipPeriod, stats, nodeP2pAddrs[i], nodeP2pAddrs, enclaveAddress, obscuroClientAddr)

			// and connect them to each other
			agg.ConnectToEthNode(miner)
			miner.AddClient(agg)

			n.ethNodes[i] = miner
			n.obscuroNodes[i] = agg
			n.obscuroClients[i] = &obscuroClient
			l1Clients[i] = miner
		}
	}

	// populate the nodes field of the L1 network
	for i := 0; i < params.NumberOfNodes; i++ {
		n.ethNodes[i].Network.(*ethereum_mock.MockEthNetwork).AllNodes = n.ethNodes
	}

	// The sequence of starting the nodes is important to catch various edge cases.
	// Here we first start the mock layer 1 nodes, with a pause between them of a fraction of a block duration.
	// The reason is to make sure that they catch up correctly.
	// Then we pause for a while, to give the L1 network enough time to create a number of blocks, which will have to be ingested by the Obscuro nodes
	// Then, we begin the starting sequence of the Obscuro nodes, again with a delay between them, to test that they are able to cach up correctly.
	// Note: Other simulations might test variations of this pattern.
	for _, m := range n.ethNodes {
		t := m
		go t.Start()
		time.Sleep(params.AvgBlockDuration / 8)
	}

	time.Sleep(params.AvgBlockDuration * 20)
	for _, m := range n.obscuroNodes {
		t := m
		go t.Start()
		time.Sleep(params.AvgBlockDuration / 3)
	}

	return l1Clients, n.obscuroClients, nodeP2pAddrs
}

func (n *networkWithOneAzureEnclave) TearDown() {
	for _, client := range n.obscuroClients {
		temp := client
		go (*temp).Call(nil, obscuroclient.RPCStopHost) //nolint:errcheck
		go (*temp).Stop()
	}

	for _, node := range n.ethNodes {
		temp := node
		go temp.Stop()
	}
}
