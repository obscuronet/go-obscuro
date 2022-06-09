package ethclient

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

// EthClient defines the interface for RPC communications with the ethereum nodes
// TODO Some of these methods are composed calls that should be decoupled in the future (ie: BlocksBetween or IsBlockAncestor)
type EthClient interface {
	BlockByHash(id common.Hash) (*types.Block, error)            // retrieves a block given a hash
	BlockByNumber(n *big.Int) (*types.Block, error)              // retrieves a block given a number - returns head block if n is nil
	SendTransaction(signedTx *types.Transaction) error           // issues an ethereum transaction (expects signed tx)
	TransactionReceipt(hash common.Hash) (*types.Receipt, error) // fetches the ethereum transaction receipt
	Nonce(address common.Address) (uint64, error)                // fetches the account nonce to use in the next transaction

	Info() Info                                                              // retrieves the node Info
	FetchHeadBlock() *types.Block                                            // retrieves the block at head height
	BlocksBetween(block *types.Block, head *types.Block) []*types.Block      // returns the blocks between two blocks
	IsBlockAncestor(block *types.Block, proof obscurocommon.L1RootHash) bool // returns if the node considers a block the ancestor
	BlockListener() chan *types.Header                                       // subscribes to new blocks and returns a listener with the blocks heads

	Stop() // tries to cleanly stop the client and release any resources

	EthClient() *ethclient.Client // returns the underlying eth client
}

// Info forces the RPC EthClient to return the data in the same format (independently of it's implementation)
type Info struct {
	ID common.Address
}
