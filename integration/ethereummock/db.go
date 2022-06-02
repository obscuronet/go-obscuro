package ethereummock

import (
	"sync"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/db"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

// Received blocks ar stored here
type blockResolverInMem struct {
	blockCache map[obscurocommon.L1RootHash]*types.Block
	m          sync.RWMutex
}

func (n *blockResolverInMem) ProofHeight(_ *core.Rollup) int64 {
	panic("implement me")
}

func (n *blockResolverInMem) Proof(_ *core.Rollup) *types.Block {
	panic("implement me")
}

func NewResolver() db.BlockResolver {
	return &blockResolverInMem{
		blockCache: map[obscurocommon.L1RootHash]*types.Block{},
		m:          sync.RWMutex{},
	}
}

func (n *blockResolverInMem) StoreBlock(block *types.Block) bool {
	n.m.Lock()
	defer n.m.Unlock()
	n.blockCache[block.Hash()] = block
	return true
}

func (n *blockResolverInMem) FetchBlock(hash obscurocommon.L1RootHash) (*types.Block, bool) {
	n.m.RLock()
	defer n.m.RUnlock()
	if hash.Hex() == obscurocommon.GenesisHash.Hex() {
		return obscurocommon.GenesisBlock, true
	}

	block, f := n.blockCache[hash]

	return block, f
}

func (n *blockResolverInMem) FetchHeadBlock() *types.Block {
	n.m.RLock()
	defer n.m.RUnlock()
	var max *types.Block
	for k := range n.blockCache {
		bh := n.blockCache[k]
		if max == nil || max.NumberU64() < bh.NumberU64() {
			max = bh
		}
	}
	return max
}

func (n *blockResolverInMem) ParentBlock(b *types.Block) (*types.Block, bool) {
	return n.FetchBlock(b.Header().ParentHash)
}

func (n *blockResolverInMem) IsAncestor(block *types.Block, maybeAncestor *types.Block) bool {
	if maybeAncestor.Hash() == block.Hash() {
		return true
	}

	if maybeAncestor.NumberU64() >= block.NumberU64() {
		return false
	}

	p, f := n.ParentBlock(block)
	if !f {
		return false
	}

	return n.IsAncestor(p, maybeAncestor)
}

func (n *blockResolverInMem) IsBlockAncestor(block *types.Block, maybeAncestor obscurocommon.L1RootHash) bool {
	if maybeAncestor == block.Hash() {
		return true
	}

	if maybeAncestor == obscurocommon.GenesisBlock.Hash() {
		return true
	}

	if block.NumberU64() == obscurocommon.L1GenesisHeight {
		return false
	}

	resolvedBlock, found := n.FetchBlock(maybeAncestor)
	if found {
		if resolvedBlock.NumberU64() >= block.NumberU64() {
			return false
		}
	}

	p, f := n.ParentBlock(block)
	if !f {
		return false
	}

	return n.IsBlockAncestor(p, maybeAncestor)
}

// The cache of included transactions
type txDBInMem struct {
	transactionsPerBlockCache map[obscurocommon.L1RootHash]map[obscurocommon.TxHash]*types.Transaction
	rpbcM                     *sync.RWMutex
}

func NewTxDB() TxDB {
	return &txDBInMem{
		transactionsPerBlockCache: make(map[obscurocommon.L1RootHash]map[obscurocommon.TxHash]*types.Transaction),
		rpbcM:                     &sync.RWMutex{},
	}
}

func (n *txDBInMem) Txs(b *types.Block) (map[obscurocommon.TxHash]*types.Transaction, bool) {
	n.rpbcM.RLock()
	val, found := n.transactionsPerBlockCache[b.Hash()]
	n.rpbcM.RUnlock()

	return val, found
}

func (n *txDBInMem) AddTxs(b *types.Block, newMap map[obscurocommon.TxHash]*types.Transaction) {
	n.rpbcM.Lock()
	n.transactionsPerBlockCache[b.Hash()] = newMap
	n.rpbcM.Unlock()
}

// removeCommittedTransactions returns a copy of `mempool` where all transactions that are exactly `committedBlocks`
// deep have been removed.
func removeCommittedTransactions(
	cb *types.Block,
	mempool []*types.Transaction,
	resolver db.BlockResolver,
	db TxDB,
) []*types.Transaction {
	if cb.NumberU64() <= obscurocommon.HeightCommittedBlocks {
		return mempool
	}

	b := cb
	i := 0

	for {
		if i == obscurocommon.HeightCommittedBlocks {
			break
		}

		p, f := resolver.ParentBlock(b)
		if !f {
			panic("wtf")
		}

		b = p
		i++
	}

	val, _ := db.Txs(b)
	//if !found {
	//	panic("should not fail here")
	//}

	return removeExisting(mempool, val)
}
