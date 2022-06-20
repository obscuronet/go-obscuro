package simulation

import (
	"testing"
	"time"

	"github.com/obscuronet/obscuro-playground/integration"
	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
	"github.com/obscuronet/obscuro-playground/integration/simulation/params"
)

// This test creates a network of L2 nodes, then injects transactions, and finally checks the resulting output blockchain
// The L2 nodes communicate with each other via sockets, and with their enclave servers via RPC.
// All nodes live in the same process, the enclaves run in individual Docker containers, and the Ethereum nodes are mocked out.
// $> docker rm $(docker stop $(docker ps -a -q --filter ancestor=obscuro_enclave --format="{{.ID}}") will stop and remove all images
func TestDockerNodesMonteCarloSimulation(t *testing.T) {
	return
	setupTestLog("docker")

	numberOfNodes := 5
	numberOfSimWallets := 5
	wallets := params.NewSimWallets(numberOfSimWallets, numberOfNodes, 1, integration.EthereumChainID, integration.ObscuroChainID)

	simParams := params.SimParams{
		NumberOfNodes:         numberOfNodes,
		AvgBlockDuration:      1 * time.Second,
		SimulationTime:        35 * time.Second,
		L1EfficiencyThreshold: 0.2,
		// Very hard to have precision here as blocks are continually produced and not dependent on the simulation execution thread
		L2EfficiencyThreshold:     0.6, // nodes might produce rollups because they receive a new block
		L2ToL1EfficiencyThreshold: 0.7, // nodes might stop producing rollups but the geth network is still going
		Wallets:                   wallets,
		StartPort:                 integration.StartPortSimulationDocker,
		ViewingKeysEnabled:        false,
	}
	simParams.AvgNetworkLatency = simParams.AvgBlockDuration / 20
	simParams.AvgGossipPeriod = simParams.AvgBlockDuration / 2

	testSimulation(t, network.NewBasicNetworkOfNodesWithDockerEnclave(wallets), &simParams)
}
