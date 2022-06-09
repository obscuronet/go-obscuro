package ethereummock

import (
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/ethclient/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/db"

	ethclient_ethereum "github.com/ethereum/go-ethereum/ethclient"
)

type L1Network interface {
	// BroadcastBlock - send the block and the parent to make sure there are no gaps
	BroadcastBlock(b obscurocommon.EncodedBlock, p obscurocommon.EncodedBlock)
	BroadcastTx(tx *types.Transaction)
}

type MiningConfig struct {
	PowTime obscurocommon.Latency
}

type TxDB interface {
	Txs(block *types.Block) (map[obscurocommon.TxHash]*types.Transaction, bool)
	AddTxs(*types.Block, map[obscurocommon.TxHash]*types.Transaction)
}

type StatsCollector interface {
	// Register when a miner has to process a reorg (a winning block from a fork)
	L1Reorg(id common.Address)
}

type NotifyNewBlock interface {
	MockedNewHead(b obscurocommon.EncodedBlock, p obscurocommon.EncodedBlock)
	MockedNewFork(b []obscurocommon.EncodedBlock)
}

type Node struct {
	ID       common.Address
	cfg      MiningConfig
	clients  []NotifyNewBlock
	Network  L1Network
	mining   bool
	stats    StatsCollector
	Resolver db.BlockResolver
	db       TxDB

	// Channels
	exitCh       chan bool // the Node stops
	exitMiningCh chan bool // the mining loop is notified to stop
	interrupt    *int32

	p2pCh       chan *types.Block       // this is where blocks received from peers are dropped
	miningCh    chan *types.Block       // this is where blocks created by the mining setup of the current node are dropped
	canonicalCh chan *types.Block       // this is where the main processing routine drops blocks that are canonical
	mempoolCh   chan *types.Transaction // where l1 transactions to be published in the next block are added

	// internal
	headInCh         chan bool
	headOutCh        chan *types.Block
	erc20ContractLib erc20contractlib.ERC20ContractLib
	mgmtContractLib  mgmtcontractlib.MgmtContractLib
}

func (m *Node) SendTransaction(tx *types.Transaction) error {
	m.Network.BroadcastTx(tx)
	return nil
}

func (m *Node) TransactionReceipt(_ common.Hash) (*types.Receipt, error) {
	// all transactions are immediately processed
	return &types.Receipt{
		Status: types.ReceiptStatusSuccessful,
	}, nil
}

func (m *Node) Nonce(common.Address) (uint64, error) {
	return 0, nil
}

// BlockListener is not used in the mock
func (m *Node) BlockListener() chan *types.Header {
	return make(chan *types.Header)
}

func (m *Node) BlockByNumber(n *big.Int) (*types.Block, error) {
	if n.Int64() == 0 {
		return obscurocommon.GenesisBlock, nil
	}
	// TODO this should be a method in the resolver
	var f bool
	for blk := m.Resolver.FetchHeadBlock(); blk.ParentHash() != obscurocommon.GenesisHash; {
		if blk.NumberU64() == n.Uint64() {
			return blk, nil
		}

		blk, f = m.Resolver.FetchBlock(blk.ParentHash())
		if !f {
			return nil, fmt.Errorf("block in the chain without a parent")
		}
	}
	return nil, ethereum.NotFound
}

func (m *Node) BlockByHash(id common.Hash) (*types.Block, error) {
	blk, f := m.Resolver.FetchBlock(id)
	if !f {
		return nil, fmt.Errorf("blk not found")
	}
	return blk, nil
}

func (m *Node) FetchHeadBlock() *types.Block {
	return m.Resolver.FetchHeadBlock()
}

func (m *Node) Info() ethclient.Info {
	return ethclient.Info{
		ID: m.ID,
	}
}

func (m *Node) IsBlockAncestor(block *types.Block, proof obscurocommon.L1RootHash) bool {
	return m.Resolver.IsBlockAncestor(block, proof)
}

// Start runs an infinite loop that listens to the two block producing channels and processes them.
// it outputs the winning blocks to the roundWinnerCh channel
func (m *Node) Start() {
	if m.mining {
		// This starts the mining
		go m.startMining()
	}

	m.Resolver.StoreBlock(obscurocommon.GenesisBlock)
	head := m.setHead(obscurocommon.GenesisBlock)

	for {
		select {
		case p2pb := <-m.p2pCh: // Received from peers
			_, received := m.Resolver.FetchBlock(p2pb.Hash())
			// only process blocks if they haven't been processed before
			if !received {
				head = m.processBlock(p2pb, head)
			}

		case mb := <-m.miningCh: // Received from the local mining
			head = m.processBlock(mb, head)
			if head.Hash() == mb.Hash() { // Ignore the locally produced block if someone else found one already
				p, found := m.Resolver.ParentBlock(mb)
				if !found {
					panic("noo")
				}
				m.Network.BroadcastBlock(obscurocommon.EncodeBlock(mb), obscurocommon.EncodeBlock(p))
			}
		case <-m.headInCh:
			m.headOutCh <- head
		case <-m.exitCh:
			return
		}
	}
}

func (m *Node) processBlock(b *types.Block, head *types.Block) *types.Block {
	m.Resolver.StoreBlock(b)
	_, f := m.Resolver.FetchBlock(b.Header().ParentHash)

	// only proceed if the parent is available
	if !f {
		log.Info(fmt.Sprintf("> M%d: Parent block not found=b_%d", obscurocommon.ShortAddress(m.ID), obscurocommon.ShortHash(b.Header().ParentHash)))
		return head
	}

	// Ignore superseded blocks
	if b.NumberU64() <= head.NumberU64() {
		return head
	}

	// Check for Reorgs
	if !m.Resolver.IsAncestor(b, head) {
		m.stats.L1Reorg(m.ID)
		fork := LCA(head, b, m.Resolver)
		log.Info(fmt.Sprintf("> M%d: L1Reorg new=b_%d(%d), old=b_%d(%d), fork=b_%d(%d)", obscurocommon.ShortAddress(m.ID), obscurocommon.ShortHash(b.Hash()), b.NumberU64(), obscurocommon.ShortHash(head.Hash()), head.NumberU64(), obscurocommon.ShortHash(fork.Hash()), fork.NumberU64()))
		return m.setFork(m.BlocksBetween(fork, b))
	}

	if b.NumberU64() > (head.NumberU64() + 1) {
		panic(fmt.Sprintf("> M%d: Should not happen", obscurocommon.ShortAddress(m.ID)))
	}

	return m.setHead(b)
}

// Notifies the Miner to start mining on the new block and the aggregator to produce rollups
func (m *Node) setHead(b *types.Block) *types.Block {
	if atomic.LoadInt32(m.interrupt) == 1 {
		return b
	}

	// notify the clients
	for _, c := range m.clients {
		t := c
		if b.NumberU64() == obscurocommon.L1GenesisHeight {
			go t.MockedNewHead(obscurocommon.EncodeBlock(b), nil)
		} else {
			p, f := m.Resolver.ParentBlock(b)
			if !f {
				panic("This should not happen")
			}
			go t.MockedNewHead(obscurocommon.EncodeBlock(b), obscurocommon.EncodeBlock(p))
		}
	}
	m.canonicalCh <- b

	return b
}

func (m *Node) setFork(blocks []*types.Block) *types.Block {
	head := blocks[len(blocks)-1]
	if atomic.LoadInt32(m.interrupt) == 1 {
		return head
	}

	fork := make([]obscurocommon.EncodedBlock, len(blocks))
	for i, block := range blocks {
		fork[i] = obscurocommon.EncodeBlock(block)
	}

	// notify the clients
	for _, c := range m.clients {
		go c.MockedNewFork(fork)
	}
	m.canonicalCh <- head

	return head
}

// P2PReceiveBlock is called by counterparties when there is a block to broadcast
// All it does is drop the blocks in a channel for processing.
func (m *Node) P2PReceiveBlock(b obscurocommon.EncodedBlock, p obscurocommon.EncodedBlock) {
	if atomic.LoadInt32(m.interrupt) == 1 {
		return
	}
	m.p2pCh <- p.DecodeBlock()
	m.p2pCh <- b.DecodeBlock()
}

// startMining - listens on the canonicalCh and schedule a go routine that produces a block after a PowTime and drop it
// on the miningCh channel
func (m *Node) startMining() {
	log.Info(fmt.Sprintf("Node-%d: starting miner...", obscurocommon.ShortAddress(m.ID)))
	// stores all transactions seen from the beginning of time.
	mempool := make([]*types.Transaction, 0)
	z := int32(0)
	interrupt := &z

	for {
		select {
		case <-m.exitMiningCh:
			return
		case tx := <-m.mempoolCh:
			mempool = append(mempool, tx)

		case canonicalBlock := <-m.canonicalCh:
			// A new canonical block was found. Start a new round based on that block.

			// remove transactions that are already considered committed
			mempool = removeCommittedTransactions(canonicalBlock, mempool, m.Resolver, m.db)

			// notify the existing mining go routine to stop mining
			atomic.StoreInt32(interrupt, 1)
			c := int32(0)
			interrupt = &c

			// Generate a random number, and wait for that number of ms. Equivalent to PoW
			// Include all rollups received during this period.
			obscurocommon.ScheduleInterrupt(m.cfg.PowTime(), interrupt, func() {
				toInclude := findNotIncludedTxs(canonicalBlock, mempool, m.Resolver, m.db)
				// todo - iterate through the rollup transactions and include only the ones with the proof on the canonical chain
				if atomic.LoadInt32(m.interrupt) == 1 {
					return
				}

				m.miningCh <- obscurocommon.NewBlock(canonicalBlock, m.ID, toInclude)
			})
		}
	}
}

// P2PGossipTx receive rollups to publish from the linked aggregators
func (m *Node) P2PGossipTx(tx *types.Transaction) {
	if atomic.LoadInt32(m.interrupt) == 1 {
		return
	}

	m.mempoolCh <- tx
}

func (m *Node) BroadcastTx(tx types.TxData) {
	m.Network.BroadcastTx(types.NewTx(tx))
}

func (m *Node) Stop() {
	// block all requests
	atomic.StoreInt32(m.interrupt, 1)
	time.Sleep(time.Millisecond * 100)

	m.exitMiningCh <- true
	m.exitCh <- true
}

func (m *Node) AddClient(client NotifyNewBlock) {
	m.clients = append(m.clients, client)
}

func (m *Node) BlocksBetween(blockA *types.Block, blockB *types.Block) []*types.Block {
	if blockA.Hash() == blockB.Hash() {
		return []*types.Block{blockA}
	}
	blocks := make([]*types.Block, 0)
	tempBlock := blockB
	var found bool
	for {
		blocks = append(blocks, tempBlock)
		if tempBlock.Hash() == blockA.Hash() {
			break
		}
		tempBlock, found = m.Resolver.ParentBlock(tempBlock)
		if !found {
			panic("should not happen")
		}
	}
	n := len(blocks)
	result := make([]*types.Block, n)
	for i, block := range blocks {
		result[n-i-1] = block
	}
	return result
}

func (m *Node) EthClient() *ethclient_ethereum.Client {
	return nil
}

func NewMiner(
	id common.Address,
	cfg MiningConfig,
	network L1Network,
	statsCollector StatsCollector,
) *Node {
	return &Node{
		ID:               id,
		mining:           true,
		cfg:              cfg,
		stats:            statsCollector,
		Resolver:         NewResolver(),
		db:               NewTxDB(),
		Network:          network,
		exitCh:           make(chan bool),
		exitMiningCh:     make(chan bool),
		interrupt:        new(int32),
		p2pCh:            make(chan *types.Block),
		miningCh:         make(chan *types.Block),
		canonicalCh:      make(chan *types.Block),
		mempoolCh:        make(chan *types.Transaction),
		headInCh:         make(chan bool),
		headOutCh:        make(chan *types.Block),
		erc20ContractLib: NewERC20ContractLibMock(),
		mgmtContractLib:  NewMgmtContractLibMock(),
	}
}
