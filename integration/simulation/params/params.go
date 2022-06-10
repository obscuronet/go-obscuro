package params

import (
	"time"

	"github.com/obscuronet/obscuro-playground/go/ethclient/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"

	"github.com/ethereum/go-ethereum/common"
)

// SimParams are the parameters for setting up the simulation.
type SimParams struct {
	NumberOfNodes int

	// A critical parameter of the simulation. The value should be as low as possible, as long as the test is still meaningful
	AvgBlockDuration  time.Duration
	AvgNetworkLatency time.Duration // artificial latency injected between sending and receiving messages on the mock network
	AvgGossipPeriod   time.Duration // POBI protocol setting

	SimulationTime time.Duration // how long the simulations should run for

	// EfficiencyThresholds represents an acceptable "dead blocks" percentage for this simulation.
	// dead blocks - Blocks that are produced and gossiped, but don't make it into the canonical chain.
	// We test the results against this threshold to catch eventual protocol errors.
	L1EfficiencyThreshold     float64
	L2EfficiencyThreshold     float64 // number of dead obscuro blocks
	L2ToL1EfficiencyThreshold float64 // number of ethereum blocks that don't include an obscuro node

	// MgmtContractLib allows parsing MgmtContract txs to and from the eth txs
	MgmtContractLib mgmtcontractlib.MgmtContractLib
	// MgmtContractLib allows parsing ERC20Contract txs to and from the eth txs
	ERC20ContractLib erc20contractlib.ERC20ContractLib

	// MgmtContractAddr defines the management contract address
	MgmtContractAddr *common.Address

	// Erc20Address defines an erc20 contract address instance that has bee deployed
	Erc20Address *common.Address

	// Contains all the wallets required by the simulation
	Wallets *SimWallets

	StartPort int // The port from which to start allocating ports. Must be unique across all simulations.

	UseViewingKeys bool // Whether the enclave should encrypt responses to sensitive requests with viewing keys
}
