package db

import (
	"sync"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

// DB lives purely in the encrypted memory space of an enclave.
// Unlike Storage, methods in this class should have minimal logic, to map them more easily to our chosen datastore.
type InMemoryDB struct {
	rollupGenesisHash common.Hash // TODO add lock protection, not needed atm

	stateMutex sync.RWMutex // Controls access to `statePerBlock`, `statePerRollup`, `headBlock`, `rollupsByHeight` and `rollups`
	blockMutex sync.RWMutex // Controls access to `blockCache`
	txMutex    sync.RWMutex // Controls access to `txsPerRollupCache`

	statePerBlock     map[obscurocommon.L1RootHash]*BlockState
	statePerRollup    map[obscurocommon.L2RootHash]*State
	headBlock         obscurocommon.L1RootHash
	rollupsByHeight   map[uint64][]*core.Rollup
	rollups           map[obscurocommon.L2RootHash]*core.Rollup
	blockCache        map[obscurocommon.L1RootHash]*types.Block
	txsPerRollupCache map[obscurocommon.L2RootHash]map[common.Hash]nodecommon.L2Tx

	sharedEnclaveSecret core.SharedEnclaveSecret
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		statePerBlock:     make(map[obscurocommon.L1RootHash]*BlockState),
		stateMutex:        sync.RWMutex{},
		rollupsByHeight:   make(map[uint64][]*core.Rollup),
		rollups:           make(map[obscurocommon.L2RootHash]*core.Rollup),
		statePerRollup:    make(map[obscurocommon.L2RootHash]*State),
		blockCache:        make(map[obscurocommon.L1RootHash]*types.Block),
		blockMutex:        sync.RWMutex{},
		txsPerRollupCache: make(map[obscurocommon.L2RootHash]map[common.Hash]nodecommon.L2Tx),
		txMutex:           sync.RWMutex{},
	}
}

func (db *InMemoryDB) StoreGenesisRollup(rol *core.Rollup) {
	db.rollupGenesisHash = rol.Hash()
	db.StoreRollup(rol)
}

func (db *InMemoryDB) FetchGenesisRollup() *core.Rollup {
	r, _ := db.FetchRollup(db.rollupGenesisHash)
	return r
}

func (db *InMemoryDB) FetchBlockState(hash obscurocommon.L1RootHash) (*BlockState, bool) {
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()

	val, found := db.statePerBlock[hash]
	return val, found
}

func (db *InMemoryDB) SetBlockState(hash obscurocommon.L1RootHash, state *BlockState) {
	db.stateMutex.Lock()
	defer db.stateMutex.Unlock()

	db.statePerBlock[hash] = state
	// todo - is there any other logic here?
	db.headBlock = hash
}

func (db *InMemoryDB) SetBlockStateNewRollup(hash obscurocommon.L1RootHash, state *BlockState) {
	db.stateMutex.Lock()
	defer db.stateMutex.Unlock()

	db.statePerBlock[hash] = state
	db.rollups[state.Head.Hash()] = state.Head
	// todo - is there any other logic here?
	db.headBlock = hash
}

func (db *InMemoryDB) SetRollupState(hash obscurocommon.L2RootHash, state *State) {
	db.stateMutex.Lock()
	defer db.stateMutex.Unlock()

	db.statePerRollup[hash] = state
}

func (db *InMemoryDB) FetchHeadBlock() obscurocommon.L1RootHash {
	return db.headBlock
}

// TODO - Pull this logic into the storage layer.
func (db *InMemoryDB) StoreRollup(rollup *core.Rollup) {
	db.stateMutex.Lock()
	defer db.stateMutex.Unlock()

	db.rollups[rollup.Hash()] = rollup
	val, found := db.rollupsByHeight[rollup.Header.Number]
	if found {
		db.rollupsByHeight[rollup.Header.Number] = append(val, rollup)
	} else {
		db.rollupsByHeight[rollup.Header.Number] = []*core.Rollup{rollup}
	}
}

func (db *InMemoryDB) FetchRollup(hash obscurocommon.L2RootHash) (*core.Rollup, bool) {
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()

	r, f := db.rollups[hash]
	return r, f
}

func (db *InMemoryDB) FetchRollups(height uint64) []*core.Rollup {
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()

	return db.rollupsByHeight[height]
}

func (db *InMemoryDB) FetchRollupState(hash obscurocommon.L2RootHash) *State {
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()

	return db.statePerRollup[hash]
}

func (db *InMemoryDB) StoreBlock(b *types.Block) {
	db.blockMutex.Lock()
	defer db.blockMutex.Unlock()

	db.blockCache[b.Hash()] = b
}

func (db *InMemoryDB) FetchBlock(hash obscurocommon.L1RootHash) (*types.Block, bool) {
	db.blockMutex.RLock()
	defer db.blockMutex.RUnlock()

	val, f := db.blockCache[hash]
	return val, f
}

func (db *InMemoryDB) FetchRollupTxs(r *core.Rollup) (map[common.Hash]nodecommon.L2Tx, bool) {
	db.txMutex.RLock()
	defer db.txMutex.RUnlock()

	val, found := db.txsPerRollupCache[r.Hash()]
	return val, found
}

func (db *InMemoryDB) StoreRollupTxs(r *core.Rollup, newMap map[common.Hash]nodecommon.L2Tx) {
	db.txMutex.Lock()
	defer db.txMutex.Unlock()

	db.txsPerRollupCache[r.Hash()] = newMap
}

func (db *InMemoryDB) StoreSecret(secret core.SharedEnclaveSecret) {
	db.sharedEnclaveSecret = secret
}

func (db *InMemoryDB) FetchSecret() core.SharedEnclaveSecret {
	return db.sharedEnclaveSecret
}
