package simulation

import (
	"testing"
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"

	"github.com/obscuronet/obscuro-playground/integration"
	"github.com/obscuronet/obscuro-playground/integration/datagenerator"
	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"
)

// This test creates a network of in memory L1 and L2 nodes, then injects transactions, and finally checks the resulting output blockchain.
// Running it long enough with various parameters will test many corner cases without having to explicitly write individual tests for them.
// The unit of time is the "AvgBlockDuration" - which is the average time between L1 blocks, which are the carriers of rollups.
// Everything else is reported to this value. This number has to be adjusted in conjunction with the number of nodes. If it's too low,
// the CPU usage will be very high during the simulation which might result in inconclusive results.
func TestInMemoryMonteCarloSimulation(t *testing.T) {
	setupTestLog("in-mem")

	// state the contract was deployed at the genesis block
	// the l2 now considers the l1 genesis block as the starting point for bootstrapping blocks
	fakeMgmtContractBlkHash := obscurocommon.GenesisHash
	simParams := params.SimParams{
		NumberOfNodes:             7,
		AvgBlockDuration:          50 * time.Millisecond,
		SimulationTime:            25 * time.Second,
		L1EfficiencyThreshold:     0.2,
		L2EfficiencyThreshold:     0.5,
		L2ToL1EfficiencyThreshold: 0.5,
		MgmtContractLib:           ethereum_mock.NewMgmtContractLibMock(),
		ERC20ContractLib:          ethereum_mock.NewERC20ContractLibMock(),
		StartPort:                 integration.StartPortSimulationInMem,
		MgmtContractBlkHash:       &fakeMgmtContractBlkHash,
	}

	simParams.AvgNetworkLatency = simParams.AvgBlockDuration / 15
	simParams.AvgGossipPeriod = simParams.AvgBlockDuration * 2 / 7

	for i := 0; i < simParams.NumberOfNodes+1; i++ {
		simParams.NodeEthWallets = append(simParams.NodeEthWallets, datagenerator.RandomWallet(integration.EthereumChainID))
		simParams.SimEthWallets = append(simParams.SimEthWallets, datagenerator.RandomWallet(integration.EthereumChainID))
	}

	testSimulation(t, network.NewBasicNetworkOfInMemoryNodes(), &simParams)
}
