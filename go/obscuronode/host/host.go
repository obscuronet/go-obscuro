package host

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/config"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"
)

// P2PCallback -the glue between the P2p layer and the node. Notifies the node when rollups and transactions are received from peers
type P2PCallback interface {
	ReceiveRollup(r obscurocommon.EncodedRollup)
	ReceiveTx(tx nodecommon.EncryptedTx)
}

// P2P is the layer responsible for sending and receiving messages to Obscuro network peers.
type P2P interface {
	StartListening(callback P2PCallback)
	StopListening()
	BroadcastRollup(r obscurocommon.EncodedRollup)
	BroadcastTx(tx nodecommon.EncryptedTx)
}

// ClientServer is the layer responsible for handling requests from Obscuro client applications.
type ClientServer interface {
	Start()
	Stop()
}

type StatsCollector interface {
	// L2Recalc - called when a node has to discard the speculative work built on top of the winner of the gossip round.
	L2Recalc(id common.Address)
	NewBlock(block *types.Block)
	NewRollup(node common.Address, rollup *nodecommon.Rollup)
	RollupWithMoreRecentProof()
}

// Node this will become the Obscuro "Node" type
type Node struct {
	config  config.HostConfig
	ID      common.Address
	shortID uint64

	P2p           P2P                 // For communication with other Obscuro nodes
	ethClient     ethclient.EthClient // For communication with the L1 node
	EnclaveClient nodecommon.Enclave  // For communication with the enclave
	clientServer  ClientServer        // For communication with Obscuro client applications

	stats StatsCollector

	// control the lifecycle
	exitNodeCh chan bool
	interrupt  *int32

	blockRPCCh   chan blockAndParent               // The channel that new blocks from the L1 node are sent to
	forkRPCCh    chan []obscurocommon.EncodedBlock // The channel that new forks from the L1 node are sent to
	rollupsP2PCh chan obscurocommon.EncodedRollup  // The channel that new rollups from peers are sent to
	txP2PCh      chan nodecommon.EncryptedTx       // The channel that new transactions from peers are sent to

	nodeDB       *DB    // Stores the node's publicly-available data
	readyForWork *int32 // Whether the node has bootstrapped the existing blocks and has the enclave secret

	mgmtContractLib mgmtcontractlib.MgmtContractLib

	// Wallet used to issue ethereum transactions
	ethWallet wallet.Wallet
}

func NewHost(
	config config.HostConfig,
	collector StatsCollector,
	p2p P2P,
	ethClient ethclient.EthClient,
	enclaveClient nodecommon.Enclave,
	ethWallet wallet.Wallet,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
) Node {
	db := NewDB()

	host := Node{
		// config
		config:  config,
		ID:      config.ID,
		shortID: obscurocommon.ShortAddress(config.ID),

		// Communication layers.
		P2p:           p2p,
		ethClient:     ethClient,
		EnclaveClient: enclaveClient,

		stats: collector,

		// lifecycle channels
		exitNodeCh: make(chan bool),
		interrupt:  new(int32),

		// incoming data
		blockRPCCh:   make(chan blockAndParent),
		forkRPCCh:    make(chan []obscurocommon.EncodedBlock),
		rollupsP2PCh: make(chan obscurocommon.EncodedRollup),
		txP2PCh:      make(chan nodecommon.EncryptedTx),

		// Initialize the node DB
		nodeDB:       db,
		readyForWork: new(int32),

		mgmtContractLib: mgmtContractLib,
		ethWallet:       ethWallet,
	}

	if config.HasClientRPC {
		host.clientServer = NewClientServer(config.ClientRPCAddress, &host)
	}

	return host
}

// Start initializes the main loop of the node
func (a *Node) Start() {
	// TODO - Log out node config.
	a.waitForEnclave()

	if a.config.IsGenesis {
		// Create the shared secret and submit it to the management contract for storage
		attestation := a.EnclaveClient.Attestation()
		encodedAttestation := nodecommon.EncodeAttestation(attestation)
		l1tx := &obscurocommon.L1StoreSecretTx{
			Secret:      a.EnclaveClient.GenerateSecret(),
			Attestation: encodedAttestation,
		}
		a.broadcastTx(a.mgmtContractLib.CreateStoreSecret(l1tx, a.ethWallet.GetNonceAndIncrement()))
	}

	if !a.EnclaveClient.IsInitialised() {
		a.requestSecret()
	}

	if a.clientServer != nil {
		a.clientServer.Start()
	}

	// todo create a channel between request secret and start processing
	a.startProcessing()
}

// Waits for enclave to be available, printing a wait message every two seconds.
func (a *Node) waitForEnclave() {
	counter := 0
	for a.EnclaveClient.IsReady() != nil {
		if counter >= 20 {
			nodecommon.LogWithID(a.shortID, "Waiting for enclave. Error: %v", a.EnclaveClient.IsReady())
			counter = 0
		}

		time.Sleep(100 * time.Millisecond)
		counter++
	}
	nodecommon.LogWithID(a.shortID, "Connected to enclave service...")
}

// Waits for initial blocks from the L1 node, printing a wait message every two seconds.
func (a *Node) waitForL1Blocks() []*types.Block {
	// It feeds the entire L1 blockchain into the enclave when it starts
	// todo - what happens with the blocks received while processing ?
	allBlocks := a.ethClient.RPCBlockchainFeed()
	counter := 0

	for len(allBlocks) == 0 {
		if counter >= 20 {
			nodecommon.LogWithID(a.shortID, "Waiting for blocks from L1 node...")
			counter = 0
		}

		time.Sleep(100 * time.Millisecond)
		allBlocks = a.ethClient.RPCBlockchainFeed()
		counter++
	}
	nodecommon.LogWithID(a.shortID, "Received %d initial blocks from L1 node.", len(allBlocks))

	return allBlocks
}

func (a *Node) startProcessing() {
	nodecommon.LogWithID(a.shortID, "Starting processing.")
	allBlocks := a.waitForL1Blocks()

	// Todo: This is a naive implementation.
	results := a.EnclaveClient.IngestBlocks(allBlocks)
	for _, result := range results {
		if !result.IngestedBlock && result.BlockNotIngestedCause != "" {
			nodecommon.LogWithID(a.shortID, "Failed to ingest block b_%d. Cause: %s",
				obscurocommon.ShortHash(result.BlockHeader.Hash()),
				result.BlockNotIngestedCause,
			)
		}
		a.storeBlockProcessingResult(result)
	}

	lastBlock := *allBlocks[len(allBlocks)-1]

	// do some asserting to check our sanity before we start the enclave
	if !blockNumberStrictlyIncreasing(allBlocks) {
		panic("We expect this list of blocks to be correctly ordered but block numbers were not strictly increasing.")
	}
	if _, err := a.ethClient.BlockByHash(lastBlock.ParentHash()); err != nil {
		panic("Parent not found for enclave start block - this should not happen.")
	}

	nodecommon.LogWithID(a.shortID, "Start enclave on block b_%d.", obscurocommon.ShortHash(lastBlock.Header().Hash()))
	a.EnclaveClient.Start(lastBlock)

	if a.config.IsGenesis {
		a.initialiseProtocol(&lastBlock)
	}

	// Start monitoring L1 blocks
	go a.monitorBlocks()

	// Only open the p2p connection when the node is fully initialised
	a.P2p.StartListening(a)

	// used as a signaling mechanism to stop processing the old block if a new L1 block arrives earlier
	i := int32(0)
	interrupt := &i
	atomic.StoreInt32(a.readyForWork, 1)

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
			rol, err := nodecommon.DecodeRollup(r)
			log.Trace(fmt.Sprintf(">   Agg%d: Received rollup: r_%d from A%d",
				a.shortID,
				obscurocommon.ShortHash(rol.Hash()),
				obscurocommon.ShortAddress(rol.Header.Agg),
			))
			if err != nil {
				nodecommon.LogWithID(a.shortID, "Could not check enclave initialisation. Cause: %v", err)
			}

			go a.EnclaveClient.SubmitRollup(nodecommon.ExtRollup{
				Header: rol.Header,
				Txs:    rol.Transactions,
			})

		case tx := <-a.txP2PCh:
			// Ignore gossiped transactions while the node is still initialising
			// TODO Handle this correctly with the Enclave Initialization process
			// TODO Enabling this without Request/RespondSecret will make non-genesis nodes ignore txs
			if a.EnclaveClient.IsInitialised() {
				if err := a.EnclaveClient.SubmitTx(tx); err != nil {
					log.Trace(fmt.Sprintf(">   Agg%d: Could not submit transaction: %s", a.shortID, err))
				}
			}

		case <-a.exitNodeCh:
			return
		}
	}
}

func blockNumberStrictlyIncreasing(blocks []*types.Block) bool {
	var latest int64 = -1
	for _, b := range blocks {
		if b.Number().Int64() <= latest {
			return false
		}
	}
	return true
}

// RPCNewHead receives the notification of new blocks from the ethereumNode Node
func (a *Node) RPCNewHead(b obscurocommon.EncodedBlock, p obscurocommon.EncodedBlock) {
	if atomic.LoadInt32(a.interrupt) == 1 {
		return
	}
	a.blockRPCCh <- blockAndParent{b, p}
}

// RPCNewFork receives the notification of a new fork from the ethereumNode
func (a *Node) RPCNewFork(b []obscurocommon.EncodedBlock) {
	if atomic.LoadInt32(a.interrupt) == 1 {
		return
	}
	a.forkRPCCh <- b
}

// P2PGossipRollup is called by counterparties when there is a Rollup to broadcast
// All it does is forward the rollup for processing to the enclave
func (a *Node) ReceiveRollup(r obscurocommon.EncodedRollup) {
	if atomic.LoadInt32(a.interrupt) == 1 {
		return
	}
	a.rollupsP2PCh <- r
}

// P2PReceiveTx receives a new transactions from the P2P network
func (a *Node) ReceiveTx(tx nodecommon.EncryptedTx) {
	if atomic.LoadInt32(a.interrupt) == 1 {
		return
	}
	a.txP2PCh <- tx
}

// RPCBalance allows to fetch the balance of one address
func (a *Node) RPCBalance(address common.Address) uint64 {
	return a.EnclaveClient.Balance(address)
}

// RPCCurrentBlockHead returns the current head of the blocks (l1)
func (a *Node) RPCCurrentBlockHead() *types.Header {
	return a.nodeDB.GetCurrentBlockHead()
}

// RPCCurrentRollupHead returns the current head of the rollups (l2)
func (a *Node) RPCCurrentRollupHead() *nodecommon.Header {
	return a.nodeDB.GetCurrentRollupHead()
}

// DB returns the DB of the node
func (a *Node) DB() *DB {
	return a.nodeDB
}

// Stop gracefully stops the node execution
func (a *Node) Stop() {
	// block all requests
	atomic.StoreInt32(a.interrupt, 1)

	a.P2p.StopListening()
	if a.clientServer != nil {
		a.clientServer.Stop()
	}

	if err := a.EnclaveClient.Stop(); err != nil {
		nodecommon.LogWithID(a.shortID, "Could not stop enclave server. Cause: %v", err.Error())
	}
	time.Sleep(time.Second)
	a.exitNodeCh <- true
	a.EnclaveClient.StopClient()
}

func (a *Node) ConnectToEthNode(node ethclient.EthClient) {
	a.ethClient = node
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
	var result nodecommon.BlockSubmissionResponse
	for _, block := range blocks {
		// For the genesis block the parent is nil
		if block != nil {
			// todo: implement proper protocol so only one host responds to this secret requests initially
			// 	for now we just have the genesis host respond until protocol implemented
			if a.isGenesis {
				a.checkForSharedSecretRequests(block)
			}

			// submit each block to the enclave for ingestion plus validation
			result = a.EnclaveClient.SubmitBlock(*block.DecodeBlock())
			a.storeBlockProcessingResult(result)
		}
	}

	if !result.IngestedBlock {
		b := blocks[len(blocks)-1].DecodeBlock()
		nodecommon.LogWithID(a.shortID, "Did not ingest block b_%d. Cause: %s", obscurocommon.ShortHash(b.Hash()), result.BlockNotIngestedCause)
		return
	}

	// Nodes can start before the genesis was published, and it makes no sense to enter the protocol.
	if result.ProducedRollup.Header != nil {
		a.P2p.BroadcastRollup(nodecommon.EncodeRollup(result.ProducedRollup.ToRollup()))

		obscurocommon.ScheduleInterrupt(a.config.GossipRoundDuration, interrupt, a.handleRoundWinner(result))
	}
}

func (a *Node) handleRoundWinner(result nodecommon.BlockSubmissionResponse) func() {
	return func() {
		if atomic.LoadInt32(a.interrupt) == 1 {
			return
		}
		// Request the round winner for the current head
		winnerRollup, isWinner, err := a.EnclaveClient.RoundWinner(result.ProducedRollup.Header.ParentHash)
		if err != nil {
			log.Panic("could not determine round winner. Cause: %s", err)
		}
		if isWinner {
			nodecommon.LogWithID(a.shortID, "Winner (b_%d) r_%d(%d).",
				obscurocommon.ShortHash(result.BlockHeader.Hash()),
				obscurocommon.ShortHash(winnerRollup.Header.Hash()),
				winnerRollup.Header.Number,
			)

			tx := &obscurocommon.L1RollupTx{
				Rollup: nodecommon.EncodeRollup(winnerRollup.ToRollup()),
			}

			a.broadcastTx(a.mgmtContractLib.CreateRollup(tx, a.ethWallet.GetNonceAndIncrement()))
		}
	}
}

func (a *Node) storeBlockProcessingResult(result nodecommon.BlockSubmissionResponse) {
	// only update the node rollup headers if the enclave has found a new rollup head
	if result.FoundNewHead {
		// adding a header will update the head if it has a higher height
		a.DB().AddRollupHeader(result.RollupHead)
	}

	// adding a header will update the head if it has a higher height
	if result.IngestedBlock {
		a.DB().AddBlockHeader(result.BlockHeader)
	}
}

// Called only by the first enclave to bootstrap the network
func (a *Node) initialiseProtocol(block *types.Block) obscurocommon.L2RootHash {
	// Create the genesis rollup and submit it to the MC
	genesisResponse := a.EnclaveClient.ProduceGenesis(block.Hash())
	nodecommon.LogWithID(a.shortID, "Initialising network. Genesis rollup r_%d.", obscurocommon.ShortHash(genesisResponse.ProducedRollup.Header.Hash()))
	l1tx := &obscurocommon.L1RollupTx{
		Rollup: nodecommon.EncodeRollup(genesisResponse.ProducedRollup.ToRollup()),
	}

	a.broadcastTx(a.mgmtContractLib.CreateRollup(l1tx, a.ethWallet.GetNonceAndIncrement()))

	return genesisResponse.ProducedRollup.Header.ParentHash
}

func (a *Node) broadcastTx(tx types.TxData) {
	// TODO add retry and deal with failures
	signedTx, err := a.ethWallet.SignTransaction(tx)
	if err != nil {
		panic(err)
	}

	err = a.ethClient.SendTransaction(signedTx)
	if err != nil {
		panic(err)
	}
}

// This method implements the procedure by which a node obtains the secret
func (a *Node) requestSecret() {
	nodecommon.LogWithID(a.shortID, "Requesting secret.")
	att := a.EnclaveClient.Attestation()
	encodedAttestation := nodecommon.EncodeAttestation(att)
	l1tx := &obscurocommon.L1RequestSecretTx{
		Attestation: encodedAttestation,
	}
	a.broadcastTx(a.mgmtContractLib.CreateRequestSecret(l1tx, a.ethWallet.GetNonceAndIncrement()))

	a.awaitSecret()
}

func (a *Node) handleStoreSecretTx(t *obscurocommon.L1StoreSecretTx) bool {
	att, err := nodecommon.DecodeAttestation(t.Attestation)
	if err != nil {
		nodecommon.LogWithID(a.shortID, "Failed to decode attestation report %s", err)
		return false
	}
	if att.Owner != a.ID {
		// this secret is encrypted for somebody else
		return false
	}
	// someone has replied for us
	err = a.EnclaveClient.InitEnclave(t.Secret)
	if err != nil {
		nodecommon.LogWithID(a.shortID, "Failed to initialise enclave with received secret. Err: %s", err)
		return false
	}
	return true
}

func (a *Node) checkForSharedSecretRequests(block obscurocommon.EncodedBlock) {
	b := block.DecodeBlock()
	for _, tx := range b.Transactions() {
		t := a.mgmtContractLib.DecodeTx(tx)
		if t == nil {
			continue
		}
		if scrtReqTx, ok := t.(*obscurocommon.L1RequestSecretTx); ok {
			att, err := nodecommon.DecodeAttestation(scrtReqTx.Attestation)
			if err != nil {
				nodecommon.LogWithID(a.shortID, "Failed to decode attestation. %s", err)
				continue
			}
			secret, err := a.EnclaveClient.ShareSecret(att)
			if err != nil {
				nodecommon.LogWithID(a.shortID, "Secret request failed, no response will be published. %s", err)
				continue
			}
			l1tx := &obscurocommon.L1StoreSecretTx{
				Secret:      secret,
				Attestation: scrtReqTx.Attestation,
			}
			a.broadcastTx(a.mgmtContractLib.CreateStoreSecret(l1tx, a.ethWallet.GetNonceAndIncrement()))
		}
	}
}

func (a *Node) monitorBlocks() {
	listener := a.ethClient.BlockListener()
	nodecommon.LogWithID(a.shortID, "Start monitoring Ethereum blocks..")

	for atomic.LoadInt32(a.interrupt) == 0 {
		select {
		case latestBlkHeader := <-listener:
			block, err := a.ethClient.BlockByHash(latestBlkHeader.Hash())
			if err != nil {
				log.Panic("could not fetch block for hash %s. Cause: %s", latestBlkHeader.Hash().String(), err)
			}
			blockParent, err := a.ethClient.BlockByHash(block.ParentHash())
			if err != nil {
				log.Panic("could not fetch block's parent with hash %s. Cause: %s", block.ParentHash().String(), err)
			}

			nodecommon.LogWithID(a.shortID, "Received a new block b_%d(%d)",
				obscurocommon.ShortHash(latestBlkHeader.Hash()),
				latestBlkHeader.Number.Uint64())
			a.RPCNewHead(obscurocommon.EncodeBlock(block), obscurocommon.EncodeBlock(blockParent))

		// this timeout ensures we don't leak the goroutine
		case <-time.After(1 * time.Second):
			// break out of select and check for interrupt on the for loop
		}
	}
}

func (a *Node) IsReady() bool {
	return atomic.LoadInt32(a.readyForWork) == 1
}

func (a *Node) awaitSecret() {
	// start listening for l1 blocks that contain the response to the request
	for {
		select {
		// todo: find a way to get rid of this case and only listen for blocks on the expected channels
		case header := <-a.ethClient.BlockListener():
			block, err := a.ethClient.BlockByHash(header.Hash())
			if err != nil {
				log.Panic("failed to retrieve block. Cause: %s:", err)
			}
			if a.checkBlockForSecretResponse(block) {
				return
			}

		case b := <-a.blockRPCCh:
			if a.checkBlockForSecretResponse(b.b.DecodeBlock()) {
				return
			}

		case <-a.forkRPCCh:
			// todo

		case <-a.rollupsP2PCh:
			// ignore rolllups from peers as we're not part of the network just yet

		case <-time.After(time.Minute):
			// This will provide useful feedback if things are stuck (and in tests if any goroutines got stranded on this select
			nodecommon.LogWithID(a.shortID, "Still waiting for secret from the L1...")

		case <-a.exitNodeCh:
			return
		}
	}
}

func (a *Node) checkBlockForSecretResponse(block *types.Block) bool {
	for _, tx := range block.Transactions() {
		t := a.mgmtContractLib.DecodeTx(tx)
		if t == nil {
			continue
		}
		if scrtTx, ok := t.(*obscurocommon.L1StoreSecretTx); ok {
			ok := a.handleStoreSecretTx(scrtTx)
			if ok {
				return true
			}
		}
	}
	// response not found
	return false
}
