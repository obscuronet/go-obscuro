package params

import (
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"

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

	// MgmtContractBlkHash defines the hash of the block where the management contract was deployed
	MgmtContractBlkHash *common.Hash

	// StableTokenContractAddr defines an erc20 contract address instance that has bee deployed
	StableTokenContractAddr *common.Address

	// NodeEthWallets are the wallets the obscuro aggregators use to create rollups and other management contract related ethereum txs
	NodeEthWallets []wallet.Wallet

	// SimEthWallets are the wallets used by the simulation aka fake users to generate traffic
	SimEthWallets []wallet.Wallet

	StartPort int // The port from which to start allocating ports. Must be unique across all simulations.
}
