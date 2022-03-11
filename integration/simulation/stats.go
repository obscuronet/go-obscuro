package simulation

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// Stats - collects information during the simulation. It can be checked programmatically.
type Stats struct {
	nrMiners int

	totalL1Blocks uint

	totalL2Blocks      uint
	l2Head             *nodecommon.Rollup
	maxRollupsPerBlock uint32
	nrEmptyBlocks      int

	totalL2Txs  int
	noL1Reorgs  map[common.Address]int
	noL2Recalcs map[common.Address]int
	// todo - actual avg block Duration

	totalDepositedAmount      uint64
	totalWithdrawnAmount      uint64
	rollupWithMoreRecentProof uint64
	nrTransferTransactions    int
	statsMu                   *sync.RWMutex
}

func NewStats(nrMiners int) *Stats {
	return &Stats{
		nrMiners:    nrMiners,
		noL1Reorgs:  map[common.Address]int{},
		noL2Recalcs: map[common.Address]int{},
		statsMu:     &sync.RWMutex{},
	}
}

func (s *Stats) L1Reorg(id common.Address) {
	s.statsMu.Lock()
	s.noL1Reorgs[id]++
	s.statsMu.Unlock()
}

func (s *Stats) L2Recalc(id common.Address) {
	s.statsMu.Lock()
	s.noL2Recalcs[id]++
	s.statsMu.Unlock()
}

func (s *Stats) NewBlock(b *types.Block) {
	s.statsMu.Lock()
	// s.l1Height = nodecommon.MaxInt(s.l1Height, b.Height)
	s.totalL1Blocks++
	s.maxRollupsPerBlock = obscurocommon.MaxInt(s.maxRollupsPerBlock, uint32(len(b.Transactions())))
	if len(b.Transactions()) == 0 {
		s.nrEmptyBlocks++
	}
	s.statsMu.Unlock()
}

func (s *Stats) NewRollup(r *nodecommon.Rollup) {
	s.statsMu.Lock()
	s.l2Head = r
	s.totalL2Blocks++
	s.totalL2Txs += len(r.Transactions)
	s.statsMu.Unlock()
}

func (s *Stats) Deposit(v uint64) {
	s.statsMu.Lock()
	s.totalDepositedAmount += v
	s.statsMu.Unlock()
}

func (s *Stats) Transfer() {
	s.statsMu.Lock()
	s.nrTransferTransactions++
	s.statsMu.Unlock()
}

func (s *Stats) Withdrawal(v uint64) {
	s.statsMu.Lock()
	s.totalWithdrawnAmount += v
	s.statsMu.Unlock()
}

func (s *Stats) RollupWithMoreRecentProof() {
	s.statsMu.Lock()
	s.rollupWithMoreRecentProof++
	s.statsMu.Unlock()
}
