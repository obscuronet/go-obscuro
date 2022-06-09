package ethclient

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/config"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

// gethRPCClient implements the EthClient interface and allows connection to a real ethereum node
type gethRPCClient struct {
	client *ethclient.Client // the underlying eth rpc client
	id     common.Address    // TODO remove the id common.Address
}

// NewEthClient instantiates a new ethclient.EthClient that connects to an ethereum node
func NewEthClient(config config.HostConfig) (EthClient, error) {
	client, err := connect(config.L1NodeHost, config.L1NodeWebsocketPort, config.L1ConnectionTimeout)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the eth node - %w", err)
	}

	log.Trace("Initialized eth node connection - port: %d - id: %s", config.L1NodeWebsocketPort, config.ID)
	return &gethRPCClient{
		client: client,
		id:     config.ID,
	}, nil
}

func (e *gethRPCClient) FetchHeadBlock() *types.Block {
	blk, err := e.client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Panic("could not fetch head block. Cause: %s", err)
	}
	return blk
}

func (e *gethRPCClient) Info() Info {
	return Info{
		ID: e.id,
	}
}

func (e *gethRPCClient) BlocksBetween(startingBlock *types.Block, lastBlock *types.Block) []*types.Block {
	// TODO this should be a stream
	var blocksBetween []*types.Block
	var err error

	for currentBlk := lastBlock; currentBlk != nil && currentBlk.Hash() != startingBlock.Hash() && currentBlk.ParentHash() != common.HexToHash(""); {
		currentBlk, err = e.BlockByHash(currentBlk.ParentHash())
		if err != nil {
			log.Panic("could not fetch parent block with hash %s. Cause: %s", currentBlk.ParentHash().String(), err)
		}
		blocksBetween = append(blocksBetween, currentBlk)
	}

	return blocksBetween
}

func (e *gethRPCClient) IsBlockAncestor(block *types.Block, maybeAncestor obscurocommon.L1RootHash) bool {
	if maybeAncestor == block.Hash() || maybeAncestor == obscurocommon.GenesisBlock.Hash() {
		return true
	}

	if block.Number().Int64() == int64(obscurocommon.L1GenesisHeight) {
		return false
	}

	resolvedBlock, err := e.BlockByHash(maybeAncestor)
	if err != nil {
		log.Panic("could not fetch parent block with hash %s. Cause: %s", maybeAncestor.String(), err)
	}
	if resolvedBlock == nil {
		if resolvedBlock.Number().Int64() >= block.Number().Int64() {
			return false
		}
	}

	p, err := e.BlockByHash(block.ParentHash())
	if err != nil {
		log.Panic("could not fetch parent block with hash %s. Cause: %s", block.ParentHash().String(), err)
	}
	if p == nil {
		return false
	}

	return e.IsBlockAncestor(p, maybeAncestor)
}

func (e *gethRPCClient) SendTransaction(signedTx *types.Transaction) error {
	return e.client.SendTransaction(context.Background(), signedTx)
}

func (e *gethRPCClient) TransactionReceipt(hash common.Hash) (*types.Receipt, error) {
	return e.client.TransactionReceipt(context.Background(), hash)
}

func (e *gethRPCClient) Nonce(account common.Address) (uint64, error) {
	return e.client.PendingNonceAt(context.Background(), account)
}

func (e *gethRPCClient) BlockListener() chan *types.Header {
	ch := make(chan *types.Header, 1)
	// TODO this should return the subscription and cleanly Unsubscribe() when the node shutsdown
	_, err := e.client.SubscribeNewHead(context.Background(), ch)
	if err != nil {
		log.Panic("could not subscribe for new head blocks. Cause: %s", err)
	}

	return ch
}

func (e *gethRPCClient) BlockByNumber(n *big.Int) (*types.Block, error) {
	return e.client.BlockByNumber(context.Background(), n)
}

func (e *gethRPCClient) BlockByHash(hash common.Hash) (*types.Block, error) {
	return e.client.BlockByHash(context.Background(), hash)
}

func (e *gethRPCClient) EthClient() *ethclient.Client {
	return e.client
}

func (e *gethRPCClient) Stop() {
	e.client.Close()
}

func connect(ipaddress string, port uint, connectionTimeout time.Duration) (*ethclient.Client, error) {
	var err error
	var c *ethclient.Client
	for start := time.Now(); time.Since(start) < connectionTimeout; time.Sleep(time.Second) {
		c, err = ethclient.Dial(fmt.Sprintf("ws://%s:%d", ipaddress, port))
		if err == nil {
			break
		}
	}

	return c, err
}
