//go:build docker
// +build docker

package simulation

import (
	"testing"
	"time"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
)

const (
	vmIP = "20.90.208.251" // Todo: replace with the IP of the vm
)

// This test creates a network of L2 nodes, then injects transactions, and finally checks the resulting output blockchain
// The L2 nodes communicate with each other via sockets, and with their enclave servers via RPC.
// All nodes and enclaves live in the same process, and the Ethereum nodes are mocked out.
func TestOnAzureEnclaveNodesMonteCarloSimulation(t *testing.T) {
	logFile := setupTestLog()
	defer logFile.Close()

	params := params.SimParams{
		NumberOfNodes:             10,
		NumberOfWallets:           5,
		AvgBlockDurationUSecs:     uint64(1_000_000),
		SimulationTime:            30 * time.Second,
		L1EfficiencyThreshold:     0.2,
		L2EfficiencyThreshold:     0.3,
		L2ToL1EfficiencyThreshold: 0.4,
	}
	params.AvgNetworkLatency = params.AvgBlockDurationUSecs / 15
	params.AvgGossipPeriod = params.AvgBlockDurationUSecs / 3

	testSimulation(t, network.NewNetworkWithOneAzureEnclave(vmIp+":11000"), params)
}
