package db

import (
	"crypto/ecdsa"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
)

// BlockResolver stores new blocks and returns information on existing blocks
type BlockResolver interface {
	// FetchBlock returns the L1 Block with the given hash and true, or (nil, false) if no such Block is stored
	FetchBlock(hash common.L1RootHash) (*types.Block, bool)
	// StoreBlock persists the L1 Block
	StoreBlock(block *types.Block) bool
	// ParentBlock returns the L1 Block's parent and true, or (nil, false)  if no parent Block is stored
	ParentBlock(block *types.Block) (*types.Block, bool)
	// IsAncestor returns true if maybeAncestor is an ancestor of the L1 Block, and false otherwise
	IsAncestor(block *types.Block, maybeAncestor *types.Block) bool
	// IsBlockAncestor returns true if maybeAncestor is an ancestor of the L1 Block, and false otherwise
	// Takes into consideration that the Block to verify might be on a branch we haven't received yet
	// Todo - this is super confusing, analyze the usage
	IsBlockAncestor(block *types.Block, maybeAncestor common.L1RootHash) bool
	// FetchHeadBlock - returns the head of the current chain
	FetchHeadBlock() *types.Block
	// ProofHeight - return the height of the L1 proof, or GenesisHeight - if the block is not known
	ProofHeight(rollup *core.Rollup) int64
	// Proof - returns the block used as proof for the rollup
	Proof(rollup *core.Rollup) *types.Block
}

type RollupResolver interface {
	// FetchRollup returns the rollup with the given hash and true, or (nil, false) if no such rollup is stored
	FetchRollup(hash common.L2RootHash) (*core.Rollup, bool)
	// FetchRollups returns all the proposed rollups with the given height
	FetchRollups(height uint64) []*core.Rollup
	// StoreRollup persists the rollup
	StoreRollup(rollup *core.Rollup)
	// ParentRollup returns the rollup's parent rollup
	ParentRollup(rollup *core.Rollup) *core.Rollup
	// StoreGenesisRollup stores the rollup genesis
	StoreGenesisRollup(rol *core.Rollup)
	// FetchGenesisRollup returns the rollup genesis
	FetchGenesisRollup() *core.Rollup
	// FetchHeadRollup returns the head rollup
	FetchHeadRollup() *core.Rollup
}

type BlockStateStorage interface {
	// FetchBlockState returns the head rollup found in the block
	FetchBlockState(blockHash common.L1RootHash) (*core.BlockState, bool)
	// FetchHeadState returns the head rollup. Returns nil if nothing recorded yet
	FetchHeadState() *core.BlockState
	// SaveNewHead save the rollup-block mapping
	SaveNewHead(state *core.BlockState, rollup *core.Rollup, receipts []*types.Receipt)
	// CreateStateDB create a database that can be used to execute transactions
	CreateStateDB(hash common.L2RootHash) *state.StateDB
	// GenesisStateDB create the original empty StateDB
	GenesisStateDB() *state.StateDB
}

type SharedSecretStorage interface {
	// FetchSecret returns the enclave's secret, returns nil if not found
	FetchSecret() *core.SharedEnclaveSecret
	// StoreSecret stores a secret in the enclave
	StoreSecret(secret core.SharedEnclaveSecret)
}

type TransactionStorage interface {
	// GetReceiptsByHash retrieves the receipts for all transactions in a given rollup.
	GetReceiptsByHash(hash gethcommon.Hash) types.Receipts

	// GetTransaction - returns the positional metadata of the tx by hash
	GetTransaction(txHash gethcommon.Hash) (*types.Transaction, gethcommon.Hash, uint64, uint64)

	// GetTransactionReceipt - returns the receipt of a tx by tx hash
	GetTransactionReceipt(txHash gethcommon.Hash) (*types.Receipt, error)

	GetSender(txHash gethcommon.Hash) (gethcommon.Address, error)
}

type AttestationStorage interface {
	// FetchAttestedKey returns the public key of an attested aggregator, returns nil if not found
	FetchAttestedKey(aggregator gethcommon.Address) *ecdsa.PublicKey

	// StoreAttestedKey - store the public key of an attested aggregator
	StoreAttestedKey(aggregator gethcommon.Address, key *ecdsa.PublicKey)
}

// Storage is the enclave's interface for interacting with the enclave's datastore
type Storage interface {
	BlockResolver
	RollupResolver
	SharedSecretStorage
	BlockStateStorage
	TransactionStorage
	AttestationStorage
}
