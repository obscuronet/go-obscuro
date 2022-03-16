package host

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/rpc"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"github.com/ethereum/go-ethereum/common"
)

type AggregatorCfg struct {
	// duration of the gossip round
	GossipRoundDuration uint64
}

type L2Network interface {
	BroadcastRollup(r obscurocommon.EncodedRollup)
	BroadcastTx(tx nodecommon.EncryptedTx)
}

type StatsCollector interface {
	// Register when a node has to discard the speculative work built on top of the winner of the gossip round.
	L2Recalc(id common.Address)
	NewBlock(block *types.Block)
	NewRollup(rollup *nodecommon.Rollup)
	RollupWithMoreRecentProof()
}

// Node this will become the Obscuro "Node" type
type Node struct {
	ID common.Address

	l2Network L2Network
	L1Node    obscurocommon.L1Node

	mining  bool // true -if this is an aggregator, false if it is a validator
	genesis bool // true - if this is the first Obscuro node which has to initialize the network
	cfg     AggregatorCfg

	stats StatsCollector

	// control the lifecycle
	exitNodeCh chan bool
	interrupt  *int32

	// blockRPCCh is where the connected L1Node node drops new blocks
	blockRPCCh chan blockAndParent

	// forkRPCCh is where new forks from the L1 notify the obscuro node
	forkRPCCh chan []obscurocommon.EncodedBlock

	// rollupsP2PCh is the channel where new rollups are gossiped to
	rollupsP2PCh chan obscurocommon.EncodedRollup

	// Interface to the logic running inside the TEE
	EnclaveClient rpc.EnclaveClient

	// Node nodeDB - stores the node public available data
	nodeDB *obscuronode.DB
}

func NewAgg(
	id common.Address,
	cfg AggregatorCfg,
	l1 obscurocommon.L1Node,
	l2Network L2Network,
	collector StatsCollector,
	genesis bool,
	enclaveClient rpc.EnclaveClient,
) Node {
	return Node{
		// config
		ID:        id,
		cfg:       cfg,
		mining:    true,
		genesis:   genesis,
		L1Node:    l1,
		l2Network: l2Network,

		stats: collector,

		// lifecycle channels
		exitNodeCh: make(chan bool),
		interrupt:  new(int32),

		// incoming data
		blockRPCCh:   make(chan blockAndParent),
		forkRPCCh:    make(chan []obscurocommon.EncodedBlock),
		rollupsP2PCh: make(chan obscurocommon.EncodedRollup),

		// State processing
		EnclaveClient: enclaveClient,

		// Initialized the node nodeDB
		nodeDB: obscuronode.NewDB(),
	}
}

// Start initializes the main loop of the node
func (a *Node) Start() {
	if a.genesis {
		// Create the shared secret and submit it to the management contract for storage
		secret, err := a.EnclaveClient.GenerateSecret()
		if err != nil {
			log.Log(fmt.Sprintf(">   Agg%d: Genesis node could not generate secret: %v", obscurocommon.ShortAddress(a.ID), err))
		}
		attestation, err := a.EnclaveClient.Attestation()
		if err != nil {
			log.Log(fmt.Sprintf(">   Agg%d: Genesis node could not produce attestation: %v", obscurocommon.ShortAddress(a.ID), err))
		}
		txData := obscurocommon.L1TxData{
			TxType:      obscurocommon.StoreSecretTx,
			Secret:      secret,
			Attestation: attestation,
		}
		a.broadcastTx(*obscurocommon.NewL1Tx(txData))
	}

	isInitialised, err := a.EnclaveClient.IsInitialised()
	if err != nil {
		log.Log(fmt.Sprintf(">   Agg%d: Could not check enclave initialisation: %v", obscurocommon.ShortAddress(a.ID), err))
	}
	if !isInitialised {
		a.requestSecret()
	}

	// todo create a channel between request secret and start processing
	a.startProcessing()
}

func (a *Node) startProcessing() {
	// Todo: This is a naive implementation.
	// It feeds the entire L1 blockchain into the enclave when it starts
	allblocks := a.L1Node.RPCBlockchainFeed()
	err := a.EnclaveClient.IngestBlocks(allblocks)
	if err != nil {
		log.Log(fmt.Sprintf(">   Agg%d: Could not ingest blocks: %s", obscurocommon.ShortAddress(a.ID), err))
	}
	// todo - what happens with the blocks received while processing ?
	err = a.EnclaveClient.Start(*allblocks[len(allblocks)-1])
	if err != nil {
		log.Log(fmt.Sprintf(">   Agg%d: Could not start: %s", obscurocommon.ShortAddress(a.ID), err))
	}

	if a.genesis {
		a.initialiseProtocol()
	}

	// used as a signaling mechanism to stop processing the old block if a new L1 block arrives earlier
	i := int32(0)
	interrupt := &i

	// Main loop - Listen for notifications From the L1 node and process them
	// Note that during processing, more recent notifications can be received.
	for {
		select {
		case b := <-a.blockRPCCh:
			interrupt = sendInterrupt(interrupt)
			a.processBlocks([]obscurocommon.EncodedBlock{b.p, b.b}, interrupt)

		case f := <-a.forkRPCCh:
			interrupt = sendInterrupt(interrupt)
			a.processBlocks(f, interrupt)

		case r := <-a.rollupsP2PCh:
			rol, err := nodecommon.Decode(r)
			if err != nil {
				log.Log(fmt.Sprintf(">   Agg%d: Could not check enclave initialisation: %v", obscurocommon.ShortAddress(a.ID), err))
			}

			go func() {
				err := a.EnclaveClient.SubmitRollup(nodecommon.ExtRollup{
					Header: rol.Header,
					Txs:    rol.Transactions,
				})
				if err != nil {
					log.Log(fmt.Sprintf(">   Agg%d: Could not submit rollup: %v", obscurocommon.ShortAddress(a.ID), err))
				}
			}()

		case <-a.exitNodeCh:
			err := a.EnclaveClient.Stop()
			if err != nil {
				log.Log(fmt.Sprintf(">   Agg%d: Could not stop: %v", obscurocommon.ShortAddress(a.ID), err))
			}
			return
		}
	}
}

// RPCNewHead receives the notification of new blocks from the L1Node Node
func (a *Node) RPCNewHead(b obscurocommon.EncodedBlock, p obscurocommon.EncodedBlock) {
	if atomic.LoadInt32(a.interrupt) == 1 {
		return
	}
	a.blockRPCCh <- blockAndParent{b, p}
}

// RPCNewFork receives the notification of a new fork from the L1Node
func (a *Node) RPCNewFork(b []obscurocommon.EncodedBlock) {
	if atomic.LoadInt32(a.interrupt) == 1 {
		return
	}
	a.forkRPCCh <- b
}

// P2PGossipRollup is called by counterparties when there is a Rollup to broadcast
// All it does is forward the rollup for processing to the enclave
func (a *Node) P2PGossipRollup(r obscurocommon.EncodedRollup) {
	if atomic.LoadInt32(a.interrupt) == 1 {
		return
	}
	a.rollupsP2PCh <- r
}

// P2PReceiveTx receives a new transactions from the P2P network
func (a *Node) P2PReceiveTx(tx nodecommon.EncryptedTx) {
	if atomic.LoadInt32(a.interrupt) == 1 {
		return
	}
	// Ignore gossiped transactions while the node is still initialising
	isInitialised, err := a.EnclaveClient.IsInitialised()
	if err != nil {
		log.Log(fmt.Sprintf(">   Agg%d: Could not check enclave initialisation: %v", obscurocommon.ShortAddress(a.ID), err))
	}
	if isInitialised {
		go func() {
			err := a.EnclaveClient.SubmitTx(tx)
			if err != nil {
				log.Log(fmt.Sprintf(">   Agg%d: Could not submit transaction: %s", obscurocommon.ShortAddress(a.ID), err))
			}
		}()
	}
}

// RPCBalance allows to fetch the balance of one address
func (a *Node) RPCBalance(address common.Address) uint64 {
	balance, err := a.EnclaveClient.Balance(address)
	if err != nil {
		log.Log(fmt.Sprintf(">   Agg%d: Could not check balance: %v", obscurocommon.ShortAddress(a.ID), err))
	}
	return balance
}

// RPCCurrentBlockHead returns the current head of the blocks (l1)
func (a *Node) RPCCurrentBlockHead() *obscuronode.BlockHeader {
	return a.nodeDB.GetCurrentBlockHead()
}

// RPCCurrentRollupHead returns the current head of the rollups (l2)
func (a *Node) RPCCurrentRollupHead() *obscuronode.RollupHeader {
	return a.nodeDB.GetCurrentRollupHead()
}

// DB returns the DB of the node
func (a *Node) DB() *obscuronode.DB {
	return a.nodeDB
}

// Stop gracefully stops the node execution
func (a *Node) Stop() {
	// block all requests
	atomic.StoreInt32(a.interrupt, 1)
	err := a.EnclaveClient.Stop()
	if err != nil {
		log.Log(fmt.Sprintf(">   Agg%d: Could not stop: %v", obscurocommon.ShortAddress(a.ID), err))
	}

	time.Sleep(time.Millisecond * 1000)
	a.exitNodeCh <- true

	err = a.EnclaveClient.StopClient()
	if err != nil {
		log.Log(fmt.Sprintf(">   Agg%d: Could not stop enclave client: %v", obscurocommon.ShortAddress(a.ID), err))
	}
}

func sendInterrupt(interrupt *int32) *int32 {
	// Notify the previous round to stop work
	atomic.StoreInt32(interrupt, 1)
	i := int32(0)
	return &i
}

type blockAndParent struct {
	b obscurocommon.EncodedBlock
	p obscurocommon.EncodedBlock
}

func (a *Node) processBlocks(blocks []obscurocommon.EncodedBlock, interrupt *int32) {
	var result enclave.BlockSubmissionResponse
	for _, block := range blocks {
		// For the genesis block the parent is nil
		if block != nil {
			a.checkForSharedSecretRequests(block)

			// submit each block to the enclave for ingestion plus validation
			var err error
			result, err = a.EnclaveClient.SubmitBlock(*block.DecodeBlock())
			if err != nil {
				log.Log(fmt.Sprintf(">   Agg%d: Could not check enclave initialisation: %v", obscurocommon.ShortAddress(a.ID), err))
			}

			// only update the node rollup headers if the enclave has ingested it
			if result.IngestedNewRollup {
				// adding a header will update the head if it has a higher height
				a.DB().AddRollupHeader(
					&obscuronode.RollupHeader{
						ID:          result.L2Hash,
						Parent:      result.L2Parent,
						Withdrawals: result.Withdrawals,
						Height:      result.L2Height,
					},
				)
			}

			// adding a header will update the head if it has a higher height
			a.DB().AddBlockHeader(
				&obscuronode.BlockHeader{
					ID:     result.L1Hash,
					Parent: result.L1Parent,
					Height: result.L1Height,
				},
			)
		}
	}

	if !result.IngestedBlock {
		b := blocks[len(blocks)-1].DecodeBlock()
		log.Log(fmt.Sprintf(">   Agg%d: Could not process block b_%d", obscurocommon.ShortAddress(a.ID), obscurocommon.ShortHash(b.Hash())))
		return
	}

	a.l2Network.BroadcastRollup(nodecommon.EncodeRollup(result.ProducedRollup.ToRollup()))

	obscurocommon.ScheduleInterrupt(a.cfg.GossipRoundDuration, interrupt, func() {
		if atomic.LoadInt32(a.interrupt) == 1 {
			return
		}
		// Request the round winner for the current head
		winnerRollup, submit, err := a.EnclaveClient.RoundWinner(result.L2Hash)
		if err != nil {
			log.Log(fmt.Sprintf(">   Agg%d: Could not determine round winner: %v", obscurocommon.ShortAddress(a.ID), err))
		}
		if submit {
			txData := obscurocommon.L1TxData{TxType: obscurocommon.RollupTx, Rollup: nodecommon.EncodeRollup(winnerRollup.ToRollup())}
			tx := obscurocommon.NewL1Tx(txData)
			t, err := obscurocommon.EncodeTx(tx)
			if err != nil {
				panic(err)
			}
			a.L1Node.BroadcastTx(t)
			// collect Stats
			// a.stats.NewRollup(DecodeRollup(winnerRollup))
		}
	})
}

// Called only by the first enclave to bootstrap the network
func (a *Node) initialiseProtocol() obscurocommon.L2RootHash {
	// Create the genesis rollup and submit it to the MC
	genesis, err := a.EnclaveClient.ProduceGenesis()
	if err != nil {
		log.Log(fmt.Sprintf(">   Agg%d: Could not produce genesis block: %v", obscurocommon.ShortAddress(a.ID), err))
	}
	txData := obscurocommon.L1TxData{TxType: obscurocommon.RollupTx, Rollup: nodecommon.EncodeRollup(genesis.ProducedRollup.ToRollup())}
	a.broadcastTx(*obscurocommon.NewL1Tx(txData))

	return genesis.L2Hash
}

func (a *Node) broadcastTx(tx obscurocommon.L1Tx) {
	t, err := obscurocommon.EncodeTx(&tx)
	if err != nil {
		panic(err)
	}
	a.L1Node.BroadcastTx(t)
}

// This method implements the procedure by which a node obtains the secret
func (a *Node) requestSecret() {
	attestation, err := a.EnclaveClient.Attestation()
	if err != nil {
		log.Log(fmt.Sprintf(">   Agg%d: Could not produce attestation: %v", obscurocommon.ShortAddress(a.ID), err))
	}
	txData := obscurocommon.L1TxData{
		TxType:      obscurocommon.RequestSecretTx,
		Attestation: attestation,
	}
	a.broadcastTx(*obscurocommon.NewL1Tx(txData))

	// start listening for l1 blocks that contain the response to the request
	for {
		select {
		case b := <-a.blockRPCCh:
			txs := b.b.DecodeBlock().Transactions()
			for _, tx := range txs {
				t := obscurocommon.TxData(tx)
				if t.TxType == obscurocommon.StoreSecretTx && t.Attestation.Owner == a.ID {
					// someone has replied
					err := a.EnclaveClient.Init(t.Secret)
					if err != nil {
						log.Log(fmt.Sprintf(">   Agg%d: Could not initialise enclave: %v", obscurocommon.ShortAddress(a.ID), err))
					}
					return
				}
			}

		case <-a.forkRPCCh:
			// todo

		case <-a.rollupsP2PCh:
			// ignore rolllups from peers as we're not part of the network just yet

		case <-a.exitNodeCh:
			return
		}
	}
}

func (a *Node) checkForSharedSecretRequests(block obscurocommon.EncodedBlock) {
	b := block.DecodeBlock()
	for _, tx := range b.Transactions() {
		t := obscurocommon.TxData(tx)
		if t.TxType == obscurocommon.RequestSecretTx {
			secret, err := a.EnclaveClient.FetchSecret(t.Attestation)
			if err != nil {
				log.Log(fmt.Sprintf(">   Agg%d: Could not fetch secret: %v", obscurocommon.ShortAddress(a.ID), err))
			}
			txData := obscurocommon.L1TxData{
				TxType:      obscurocommon.StoreSecretTx,
				Secret:      secret,
				Attestation: t.Attestation,
			}
			a.broadcastTx(*obscurocommon.NewL1Tx(txData))
		}
	}
}
