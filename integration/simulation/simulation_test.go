package simulation

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"
)

func TestSimulation(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	uuid.EnableRandPool()

	// create a folder specific for the test
	err := os.MkdirAll("../.build/simulations/", 0o700)
	if err != nil {
		panic(err)
	}
	fileName := fmt.Sprintf("simulation-result-%d-*.txt", time.Now().Unix())
	f, err := os.CreateTemp("../.build/simulations", fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	obscurocommon.SetLog(f)

	blockDuration := uint64(20_000)
	l1netw, l2netw := RunSimulation(5, 10, 15, blockDuration, blockDuration/15, blockDuration/3)
	firstNode := l2netw.nodes[0]
	checkBlockchainValidity(t, l1netw, l2netw, firstNode.Enclave.TestDB(), firstNode.Enclave.TestPeekHead().Head, l1netw.nodes[0].Resolver)
	stats := l1netw.Stats
	fmt.Printf("%+v\n", stats)
	// pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
}

func checkBlockchainValidity(t *testing.T, l1Network L1NetworkCfg, l2Network L2NetworkCfg, db enclave.DB, r *enclave.Rollup, resolver obscurocommon.BlockResolver) {
	stats := l1Network.Stats
	p := r.Proof(resolver)
	validateL1(t, p, stats, db, resolver)
	totalWithdrawn := validateL2(t, r, stats, db)
	validateL2State(t, l2Network, stats, totalWithdrawn)
}

// For this simulation, this represents an acceptable "dead blocks" percentage.
// dead blocks - Blocks that are produced and gossiped, but don't make it into the canonical chain.
// We test the results against this threshold to catch eventual protocol errors.
const L1EfficiencyThreashold = 0.2
const L2EfficiencyThreashold = 0.3

// Sanity check
func validateL1(t *testing.T, b *types.Block, s *Stats, db enclave.DB, resolver obscurocommon.BlockResolver) {
	deposits := make([]common.Hash, 0)
	rollups := make([]obscurocommon.L2RootHash, 0)
	s.l1Height = db.HeightBlock(b)
	totalDeposited := uint64(0)

	blockchain := ethereum_mock.BlocksBetween(obscurocommon.GenesisBlock, b, resolver)
	headRollup := &enclave.GenesisRollup
	for _, block := range blockchain {
		for _, tr := range block.Transactions() {
			currentRollups := make([]*enclave.Rollup, 0)
			tx := obscurocommon.TxData(tr)
			switch tx.TxType {
			case obscurocommon.DepositTx:
				deposits = append(deposits, tr.Hash())
				totalDeposited += tx.Amount
			case obscurocommon.RollupTx:
				r := nodecommon.DecodeRollup(tx.Rollup)
				rollups = append(rollups, r.Hash())
				if obscurocommon.IsBlockAncestor(r.Header.L1Proof, b, resolver) {
					// only count the rollup if it is published in the right branch
					// todo - once logic is added to the l1 - this can be made into a check
					currentRollups = append(currentRollups, enclave.DecryptRollup(r))
					s.NewRollup(r)
				}
			case obscurocommon.RequestSecretTx:
			case obscurocommon.StoreSecretTx:
			}
			r, _ := enclave.FindWinner(headRollup, currentRollups, db, resolver)
			if r != nil {
				headRollup = r
			}
		}
	}

	if len(obscurocommon.FindHashDups(deposits)) > 0 {
		dups := obscurocommon.FindHashDups(deposits)
		t.Errorf("Found Deposit duplicates: %v", dups)
	}
	if len(obscurocommon.FindRollupDups(rollups)) > 0 {
		dups := obscurocommon.FindRollupDups(rollups)
		t.Errorf("Found Rollup duplicates: %v", dups)
	}
	if totalDeposited != s.totalDepositedAmount {
		t.Errorf("Deposit amounts don't match. Found %d , expected %d", totalDeposited, s.totalDepositedAmount)
	}

	efficiency := float64(s.totalL1Blocks-s.l1Height) / float64(s.totalL1Blocks)
	if efficiency > L1EfficiencyThreashold {
		t.Errorf("Efficiency in L1 is %f. Expected:%f", efficiency, L1EfficiencyThreashold)
	}

	// todo
	for nodeID, reorgs := range s.noL1Reorgs {
		eff := float64(reorgs) / float64(s.l1Height)
		if eff > L1EfficiencyThreashold {
			t.Errorf("Efficiency for node %d in L1 is %f. Expected:%f", nodeID, eff, L1EfficiencyThreashold)
		}
	}
}

func validateL2(t *testing.T, r *enclave.Rollup, s *Stats, db enclave.DB) uint64 {
	s.l2Height = db.HeightRollup(r)
	transfers := make([]common.Hash, 0)
	withdrawalTxs := make([]enclave.L2Tx, 0)
	withdrawalRequests := make([]nodecommon.Withdrawal, 0)
	for {
		if db.HeightRollup(r) == obscurocommon.L2GenesisHeight {
			break
		}
		for i := range r.Transactions {
			tx := r.Transactions[i]
			txData := enclave.TxData(&tx)
			switch txData.Type {
			case enclave.TransferTx:
				transfers = append(transfers, tx.Hash())
			case enclave.WithdrawalTx:
				withdrawalTxs = append(withdrawalTxs, tx)
			default:
				panic("Invalid tx type")
			}
		}
		withdrawalRequests = append(withdrawalRequests, r.Header.Withdrawals...)
		r = db.ParentRollup(r)
	}
	// todo - check that proofs are on the canonical chain

	if len(obscurocommon.FindHashDups(transfers)) > 0 {
		dups := obscurocommon.FindHashDups(transfers)
		t.Errorf("Found L2 txs duplicates: %v", dups)
	}
	if len(transfers) != s.nrTransferTransactions {
		t.Errorf("Nr of transfers don't match. Found %d , expected %d", len(transfers), s.nrTransferTransactions)
	}
	if sumWithdrawalTxs(withdrawalTxs) != s.totalWithdrawnAmount {
		t.Errorf("Withdrawal tx amounts don't match. Found %d , expected %d", sumWithdrawalTxs(withdrawalTxs), s.totalWithdrawnAmount)
	}
	if sumWithdrawals(withdrawalRequests) > s.totalWithdrawnAmount {
		t.Errorf("The amount withdrawn %d exceeds the actual amount requested %d", sumWithdrawals(withdrawalRequests), s.totalWithdrawnAmount)
	}
	efficiency := float64(s.totalL2Blocks-s.l2Height) / float64(s.totalL2Blocks)
	if efficiency > L2EfficiencyThreashold {
		t.Errorf("Efficiency in L2 is %f. Expected:%f", efficiency, L2EfficiencyThreashold)
	}

	return sumWithdrawals(withdrawalRequests)
}

func sumWithdrawals(w []nodecommon.Withdrawal) uint64 {
	sum := uint64(0)
	for _, r := range w {
		sum += r.Amount
	}
	return sum
}

func sumWithdrawalTxs(t []enclave.L2Tx) uint64 {
	sum := uint64(0)
	for i := range t {
		txData := enclave.TxData(&t[i])
		sum += txData.Amount
	}

	return sum
}

func validateL2State(t *testing.T, l2Network L2NetworkCfg, s *Stats, totalWithdrawn uint64) {
	finalAmount := s.totalDepositedAmount - totalWithdrawn
	// Check that the state on all nodes is valid
	for _, observer := range l2Network.nodes {
		// read the last state
		lastState := observer.Enclave.TestPeekHead()
		total := totalBalance(lastState)
		if total != finalAmount {
			t.Errorf("The amount of money in accounts on node %d does not match the amount deposited. Found %d , expected %d", observer.ID, total, finalAmount)
		}
	}

	// TODO Check that processing transactions in the order specified in the list results in the same balances
	// walk the blocks in reverse direction, execute deposits and transactions and compare to the state in the rollup
}

func totalBalance(s enclave.BlockState) uint64 {
	tot := uint64(0)
	for _, bal := range s.State {
		tot += bal
	}
	return tot
}
