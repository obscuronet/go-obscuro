package db

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"

	"github.com/obscuronet/go-obscuro/go/common/log"

	obscurorawdb "github.com/obscuronet/go-obscuro/go/enclave/db/rawdb"

	"github.com/ethereum/go-ethereum/params"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
)

type storageImpl struct {
	db          ethdb.Database
	stateDB     state.Database
	nodeID      uint64
	chainConfig *params.ChainConfig
}

// ErrTxNotFound indicates that a transaction could not be found.
var ErrTxNotFound = errors.New("transaction not found")

func NewStorage(backingDB ethdb.Database, nodeID uint64, chainConfig *params.ChainConfig) Storage {
	return &storageImpl{
		db:          backingDB,
		stateDB:     state.NewDatabase(backingDB),
		nodeID:      nodeID,
		chainConfig: chainConfig,
	}
}

func (s *storageImpl) StoreGenesisRollup(rol *core.Rollup) {
	obscurorawdb.WriteGenesisHash(s.db, rol.Hash())
	s.StoreRollup(rol)
}

func (s *storageImpl) FetchGenesisRollup() *core.Rollup {
	hash := obscurorawdb.ReadGenesisHash(s.db)
	if hash == nil {
		return nil
	}
	r, _ := s.FetchRollup(*hash)
	return r
}

func (s *storageImpl) FetchHeadRollup() *core.Rollup {
	hash := obscurorawdb.ReadHeadRollupHash(s.db)
	if hash == (gethcommon.Hash{}) {
		return nil
	}
	r, _ := s.FetchRollup(hash)
	return r
}

func (s *storageImpl) StoreRollup(rollup *core.Rollup) {
	s.assertSecretAvailable()

	batch := s.db.NewBatch()
	obscurorawdb.WriteRollup(batch, rollup)
	if err := batch.Write(); err != nil {
		log.Panic("could not write rollup to storage. Cause: %s", err)
	}
}

func (s *storageImpl) FetchRollup(hash common.L2RootHash) (*core.Rollup, bool) {
	s.assertSecretAvailable()
	r := obscurorawdb.ReadRollup(s.db, hash)
	if r != nil {
		return r, true
	}
	return nil, false
}

func (s *storageImpl) FetchRollups(height uint64) []*core.Rollup {
	s.assertSecretAvailable()
	return obscurorawdb.ReadRollupsForHeight(s.db, height)
}

func (s *storageImpl) StoreBlock(b *types.Block) bool {
	s.assertSecretAvailable()
	rawdb.WriteBlock(s.db, b)
	return true
}

func (s *storageImpl) FetchBlock(hash common.L1RootHash) (*types.Block, bool) {
	s.assertSecretAvailable()
	height := rawdb.ReadHeaderNumber(s.db, hash)
	if height == nil {
		return nil, false
	}
	b := rawdb.ReadBlock(s.db, hash, *height)
	if b != nil {
		return b, true
	}
	return nil, false
}

func (s *storageImpl) FetchHeadBlock() *types.Block {
	s.assertSecretAvailable()
	b, _ := s.FetchBlock(rawdb.ReadHeadHeaderHash(s.db))
	return b
}

func (s *storageImpl) StoreSecret(secret core.SharedEnclaveSecret) {
	obscurorawdb.WriteSharedSecret(s.db, secret)
}

func (s *storageImpl) FetchSecret() *core.SharedEnclaveSecret {
	return obscurorawdb.ReadSharedSecret(s.db)
}

func (s *storageImpl) ParentRollup(r *core.Rollup) *core.Rollup {
	s.assertSecretAvailable()
	parent, found := s.FetchRollup(r.Header.ParentHash)
	if !found {
		common.LogWithID(s.nodeID, "Could not find rollup: r_%d", common.ShortHash(r.Hash()))
		return nil
	}
	return parent
}

func (s *storageImpl) ParentBlock(b *types.Block) (*types.Block, bool) {
	s.assertSecretAvailable()
	return s.FetchBlock(b.Header().ParentHash)
}

func (s *storageImpl) IsAncestor(block *types.Block, maybeAncestor *types.Block) bool {
	s.assertSecretAvailable()
	if bytes.Equal(maybeAncestor.Hash().Bytes(), block.Hash().Bytes()) {
		return true
	}

	if maybeAncestor.NumberU64() >= block.NumberU64() {
		return false
	}

	p, f := s.ParentBlock(block)
	if !f {
		return false
	}

	return s.IsAncestor(p, maybeAncestor)
}

func (s *storageImpl) IsBlockAncestor(block *types.Block, maybeAncestor common.L1RootHash) bool {
	s.assertSecretAvailable()
	if bytes.Equal(maybeAncestor.Bytes(), block.Hash().Bytes()) {
		return true
	}

	if bytes.Equal(maybeAncestor.Bytes(), common.GenesisBlock.Hash().Bytes()) {
		return true
	}

	if block.NumberU64() == common.L1GenesisHeight {
		return false
	}

	resolvedBlock, found := s.FetchBlock(maybeAncestor)
	if found {
		if resolvedBlock.NumberU64() >= block.NumberU64() {
			return false
		}
	}

	p, f := s.ParentBlock(block)
	if !f {
		return false
	}

	return s.IsBlockAncestor(p, maybeAncestor)
}

func (s *storageImpl) assertSecretAvailable() {
	// TODO uncomment this
	//if s.FetchSecret() == nil {
	//	panic("Enclave not initialized")
	//}
}

// ProofHeight - return the height of the L1 proof, or GenesisHeight - if the block is not known
// todo - find a better way. This is a workaround to handle rollups created with proofs that haven't propagated yet
func (s *storageImpl) ProofHeight(r *core.Rollup) int64 {
	v, f := s.FetchBlock(r.Header.L1Proof)
	if !f {
		return -1
	}
	return int64(v.NumberU64())
}

func (s *storageImpl) Proof(r *core.Rollup) *types.Block {
	v, f := s.FetchBlock(r.Header.L1Proof)
	if !f {
		log.Panic("could not find proof for this rollup")
	}
	return v
}

func (s *storageImpl) FetchBlockState(hash common.L1RootHash) (*core.BlockState, bool) {
	bs := obscurorawdb.ReadBlockState(s.db, hash)
	if bs != nil {
		return bs, true
	}
	return nil, false
}

func (s *storageImpl) SaveNewHead(state *core.BlockState, rollup *core.Rollup, receipts []*types.Receipt) {
	batch := s.db.NewBatch()

	if state.FoundNewRollup {
		obscurorawdb.WriteRollup(batch, rollup)
		obscurorawdb.WriteHeadHeaderHash(batch, rollup.Hash())
		obscurorawdb.WriteCanonicalHash(batch, rollup.Hash(), rollup.NumberU64())
		obscurorawdb.WriteTxLookupEntriesByBlock(batch, rollup)
		obscurorawdb.WriteHeadRollupHash(batch, rollup.Hash())
		obscurorawdb.WriteReceipts(batch, rollup.Hash(), rollup.NumberU64(), receipts)
	}

	obscurorawdb.WriteBlockState(batch, state)
	rawdb.WriteHeadHeaderHash(batch, state.Block)

	if err := batch.Write(); err != nil {
		log.Panic("could not save new head. Cause: %s", err)
	}
}

func (s *storageImpl) CreateStateDB(hash common.L2RootHash) *state.StateDB {
	rollup, f := s.FetchRollup(hash)
	if !f {
		log.Panic("could not retrieve rollup for hash %s", hash.String())
	}
	// todo - snapshots?
	statedb, err := state.New(rollup.Header.Root, s.stateDB, nil)
	if err != nil {
		log.Panic("could not create state DB. Cause: %s", err)
	}
	return statedb
}

func (s *storageImpl) GenesisStateDB() *state.StateDB {
	return s.CreateStateDB(s.FetchGenesisRollup().Hash())
}

func (s *storageImpl) FetchHeadState() *core.BlockState {
	h := rawdb.ReadHeadHeaderHash(s.db)
	if (bytes.Equal(h.Bytes(), gethcommon.Hash{}.Bytes())) {
		log.Error("Agg%d: could not read head header hash from storage", s.nodeID)
		return nil
	}
	return obscurorawdb.ReadBlockState(s.db, h)
}

// GetReceiptsByHash retrieves the receipts for all transactions in a given rollup.
func (s *storageImpl) GetReceiptsByHash(hash gethcommon.Hash) types.Receipts {
	number := obscurorawdb.ReadHeaderNumber(s.db, hash)
	if number == nil {
		return nil
	}
	receipts := obscurorawdb.ReadReceipts(s.db, hash, *number, s.chainConfig)
	return receipts
}

func (s *storageImpl) GetTransaction(txHash gethcommon.Hash) (*types.Transaction, gethcommon.Hash, uint64, uint64) {
	tx, blockHash, blockNumber, index := obscurorawdb.ReadTransaction(s.db, txHash)
	return tx, blockHash, blockNumber, index
}

func (s *storageImpl) GetSender(txHash gethcommon.Hash) (gethcommon.Address, error) {
	tx, _, _, _ := s.GetTransaction(txHash) //nolint:dogsled
	if tx == nil {
		return gethcommon.Address{}, ErrTxNotFound
	}
	// todo - make the signer a field of the rollup chain
	msg, err := tx.AsMessage(types.NewLondonSigner(tx.ChainId()), nil)
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("could not convert transaction to message to retrieve sender address in eth_getTransactionReceipt request. Cause: %w", err)
	}
	return msg.From(), nil
}

func (s *storageImpl) GetTransactionReceipt(txHash gethcommon.Hash) (*types.Receipt, error) {
	tx, blockHash, _, index := obscurorawdb.ReadTransaction(s.db, txHash)
	if tx == nil {
		return nil, ErrTxNotFound
	}

	receipts := s.GetReceiptsByHash(blockHash)
	if len(receipts) <= int(index) {
		return nil, fmt.Errorf("receipt index not matching the transactions in block: %s", blockHash.Hex())
	}
	receipt := receipts[index]

	return receipt, nil
}

func (s *storageImpl) FetchAttestedKey(aggregator gethcommon.Address) *ecdsa.PublicKey {
	return obscurorawdb.ReadAttestationKey(s.db, aggregator)
}

func (s *storageImpl) StoreAttestedKey(aggregator gethcommon.Address, key *ecdsa.PublicKey) {
	obscurorawdb.WriteAttestationKey(s.db, aggregator, key)
}
