package enclave

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

const ChainID = 777 // The unique ID for the Obscuro chain. Required for Geth signing.
var genesisParentHash = common.Hash{}

// todo - this should become an elaborate data structure
type SharedEnclaveSecret []byte

type StatsCollector interface {
	// Register when a node has to discard the speculative work built on top of the winner of the gossip round.
	L2Recalc(id common.Address)
	RollupWithMoreRecentProof()
}

type enclaveImpl struct {
	node           common.Address
	mining         bool
	storage        Storage
	blockResolver  BlockResolver
	statsCollector StatsCollector
	l1Blockchain   *core.BlockChain

	txCh                 chan nodecommon.L2Tx
	roundWinnerCh        chan *Rollup
	exitCh               chan bool
	speculativeWorkInCh  chan bool
	speculativeWorkOutCh chan speculativeWork
}

func (e *enclaveImpl) IsReady() error {
	return nil // The enclave is local so it is always ready
}

func (e *enclaveImpl) StopClient() {
	// The enclave is local so there is no client to stop
}

func (e *enclaveImpl) Start(block types.Block) {
	// start the speculative rollup execution loop on its own go routine
	go e.start(block)
}

func (e *enclaveImpl) start(block types.Block) {
	var currentHead *Rollup
	var currentState *State
	var currentProcessedTxs []nodecommon.L2Tx
	currentProcessedTxsMap := make(map[common.Hash]nodecommon.L2Tx)
	// determine whether the block where the speculative execution will start already contains Obscuro state
	blockState, f := e.storage.FetchBlockState(block.Hash())
	if f {
		currentHead = blockState.head
		if currentHead != nil {
			currentState = copyState(e.storage.FetchRollupState(currentHead.Hash()))
		}
	}

	for {
		select {
		// A new winner was found after gossiping. Start speculatively executing incoming transactions to already have a rollup ready when the next round starts.
		case winnerRollup := <-e.roundWinnerCh:

			currentHead = winnerRollup
			currentState = copyState(e.storage.FetchRollupState(winnerRollup.Hash()))

			// determine the transactions that were not yet included
			currentProcessedTxs = currentTxs(winnerRollup, e.storage.FetchMempoolTxs(), e.storage)
			currentProcessedTxsMap = makeMap(currentProcessedTxs)

			// calculate the State after executing them
			currentState = executeTransactions(currentProcessedTxs, currentState)

		case tx := <-e.txCh:
			// only process transactions if there is already a rollup to use as parent
			if currentHead != nil {
				_, found := currentProcessedTxsMap[tx.Hash()]
				if !found {
					currentProcessedTxsMap[tx.Hash()] = tx
					currentProcessedTxs = append(currentProcessedTxs, tx)
					executeTx(currentState, tx)
				}
			}

		case <-e.speculativeWorkInCh:
			b := make([]nodecommon.L2Tx, 0, len(currentProcessedTxs))
			b = append(b, currentProcessedTxs...)
			state := copyState(currentState)
			e.speculativeWorkOutCh <- speculativeWork{
				r:   currentHead,
				s:   state,
				txs: b,
			}

		case <-e.exitCh:
			return
		}
	}
}

func (e *enclaveImpl) ProduceGenesis(blkHash common.Hash) nodecommon.BlockSubmissionResponse {
	rolGenesis := NewRollup(blkHash, nil, obscurocommon.L2GenesisHeight, common.HexToAddress("0x0"), []nodecommon.L2Tx{}, []nodecommon.Withdrawal{}, obscurocommon.GenerateNonce(), "")
	return nodecommon.BlockSubmissionResponse{
		L2Hash:         rolGenesis.Header.Hash(),
		L1Hash:         blkHash,
		ProducedRollup: rolGenesis.ToExtRollup(),
		IngestedBlock:  true,
	}
}

// IngestBlocks is used to update the enclave with the full history of the L1 chain to date.
func (e *enclaveImpl) IngestBlocks(blocks []*types.Block) []nodecommon.BlockSubmissionResponse {
	result := make([]nodecommon.BlockSubmissionResponse, len(blocks))
	for i, block := range blocks {
		e.storage.StoreBlock(block)

		// If configured to do so, we check that the block is a valid Ethereum block.
		if e.l1Blockchain != nil && block.ParentHash() != genesisParentHash {
			_, err := e.l1Blockchain.InsertChain(types.Blocks{block})
			if err != nil {
				causeMsg := fmt.Sprintf("Block was invalid: %v", err)
				result[i] = nodecommon.BlockSubmissionResponse{IngestedBlock: false, BlockNotIngestedCause: causeMsg}
				continue
			}
		}

		bs := updateState(block, e.storage, e.blockResolver)
		if bs == nil {
			result[i] = e.noBlockStateBlockSubmissionResponse(block)
		} else {
			var rollup nodecommon.ExtRollup
			if bs.foundNewRollup {
				rollup = bs.head.ToExtRollup()
			}
			result[i] = e.blockStateBlockSubmissionResponse(bs, rollup)
		}
	}

	return result
}

// SubmitBlock is used to update the enclave with an additional block.
func (e *enclaveImpl) SubmitBlock(block types.Block) nodecommon.BlockSubmissionResponse {
	_, foundBlock := e.storage.FetchBlock(block.Hash())
	if foundBlock {
		return nodecommon.BlockSubmissionResponse{IngestedBlock: false, BlockNotIngestedCause: "Block already ingested."}
	}

	stored := e.storage.StoreBlock(&block)
	if !stored {
		return nodecommon.BlockSubmissionResponse{IngestedBlock: false}
	}

	_, f := e.storage.FetchBlock(block.Header().ParentHash)
	if !f && e.storage.HeightBlock(&block) > obscurocommon.L1GenesisHeight {
		return nodecommon.BlockSubmissionResponse{IngestedBlock: false, BlockNotIngestedCause: "Block parent not stored."}
	}

	// If configured to do so, we check that the block is a valid Ethereum block.
	if e.l1Blockchain != nil && block.ParentHash() != genesisParentHash {
		_, err := e.l1Blockchain.InsertChain(types.Blocks{&block})
		if err != nil {
			causeMsg := fmt.Sprintf("Block was invalid: %v", err)
			return nodecommon.BlockSubmissionResponse{IngestedBlock: false, BlockNotIngestedCause: causeMsg}
		}
	}

	blockState := updateState(&block, e.storage, e.blockResolver)
	if blockState == nil {
		return e.noBlockStateBlockSubmissionResponse(&block)
	}

	// todo - A verifier node will not produce rollups, we can check the e.mining to get the node behaviour
	e.storage.RemoveMempoolTxs(historicTxs(blockState.head, e.storage))
	r := e.produceRollup(&block, blockState)
	// todo - should store proposal rollups in a different storage as they are ephemeral (round based)
	e.storage.StoreRollup(r)

	log.Log(fmt.Sprintf("Agg%d:> Processed block: b_%d", obscurocommon.ShortAddress(e.node), obscurocommon.ShortHash(block.Hash())))

	return e.blockStateBlockSubmissionResponse(blockState, r.ToExtRollup())
}

func (e *enclaveImpl) SubmitRollup(rollup nodecommon.ExtRollup) {
	r := Rollup{
		Header:       rollup.Header,
		Transactions: decryptTransactions(rollup.Txs),
	}

	// only store if the parent exists
	_, found := e.storage.FetchRollup(r.Header.ParentHash)
	if found {
		e.storage.StoreRollup(&r)
	} else {
		log.Log(fmt.Sprintf("Agg%d:> Received rollup with no parent: r_%d", obscurocommon.ShortAddress(e.node), obscurocommon.ShortHash(r.Hash())))
	}
}

func (e *enclaveImpl) SubmitTx(tx nodecommon.EncryptedTx) error {
	decryptedTx := DecryptTx(tx)
	err := verifySignature(&decryptedTx)
	if err != nil {
		return err
	}
	e.storage.AddMempoolTx(decryptedTx)
	e.txCh <- decryptedTx
	return nil
}

// Checks that the L2Tx has a valid signature.
func verifySignature(decryptedTx *nodecommon.L2Tx) error {
	signer := types.NewLondonSigner(big.NewInt(ChainID))
	_, err := types.Sender(signer, decryptedTx)
	return err
}

func (e *enclaveImpl) RoundWinner(parent obscurocommon.L2RootHash) (nodecommon.ExtRollup, bool, error) {
	head, found := e.storage.FetchRollup(parent)
	if !found {
		return nodecommon.ExtRollup{}, false, fmt.Errorf("rollup not found: r_%s", parent) //nolint
	}

	rollupsReceivedFromPeers := e.storage.FetchRollups(head.Header.Height + 1)
	// filter out rollups with a different Parent
	var usefulRollups []*Rollup
	for _, rol := range rollupsReceivedFromPeers {
		p := e.storage.ParentRollup(rol)
		if p.Hash() == head.Hash() {
			usefulRollups = append(usefulRollups, rol)
		}
	}

	parentState := e.storage.FetchRollupState(head.Hash())
	// determine the winner of the round
	winnerRollup, s := findRoundWinner(usefulRollups, head, parentState, e.storage, e.blockResolver)

	e.storage.SetRollupState(winnerRollup.Hash(), s)
	go e.notifySpeculative(winnerRollup)

	// we are the winner
	if winnerRollup.Header.Agg == e.node {
		v := winnerRollup.Proof(e.blockResolver)
		w := e.storage.ParentRollup(winnerRollup)
		log.Log(fmt.Sprintf(">   Agg%d: create rollup=r_%d(%d)[r_%d]{proof=b_%d}. Txs: %v. State=%v.",
			obscurocommon.ShortAddress(e.node),
			obscurocommon.ShortHash(winnerRollup.Hash()), winnerRollup.Header.Height,
			obscurocommon.ShortHash(w.Hash()),
			obscurocommon.ShortHash(v.Hash()),
			printTxs(winnerRollup.Transactions),
			winnerRollup.Header.State),
		)
		return winnerRollup.ToExtRollup(), true, nil
	}
	return nodecommon.ExtRollup{}, false, nil
}

func (e *enclaveImpl) notifySpeculative(winnerRollup *Rollup) {
	e.roundWinnerCh <- winnerRollup
}

func (e *enclaveImpl) Balance(address common.Address) uint64 {
	// todo user encryption
	return e.storage.FetchHeadState().state.balances[address]
}

func (e *enclaveImpl) produceRollup(b *types.Block, bs *blockState) *Rollup {
	// retrieve the speculatively calculated State based on the previous winner and the incoming transactions
	e.speculativeWorkInCh <- true
	speculativeRollup := <-e.speculativeWorkOutCh

	newRollupTxs := speculativeRollup.txs
	newRollupState := speculativeRollup.s

	// the speculative execution has been processing on top of the wrong parent - due to failure in gossip or publishing to L1
	if (speculativeRollup.r == nil) || (speculativeRollup.r.Hash() != bs.head.Hash()) {
		if speculativeRollup.r != nil {
			log.Log(fmt.Sprintf(">   Agg%d: Recalculate. speculative=r_%d(%d), published=r_%d(%d)",
				obscurocommon.ShortAddress(e.node),
				obscurocommon.ShortHash(speculativeRollup.r.Hash()),
				speculativeRollup.r.Header.Height,
				obscurocommon.ShortHash(bs.head.Hash()),
				bs.head.Header.Height),
			)
			if e.statsCollector != nil {
				e.statsCollector.L2Recalc(e.node)
			}
		}

		// determine transactions to include in new rollup and process them
		newRollupTxs = currentTxs(bs.head, e.storage.FetchMempoolTxs(), e.storage)
		newRollupState = executeTransactions(newRollupTxs, bs.state)
	}

	// always process deposits last
	// process deposits from the proof of the parent to the current block (which is the proof of the new rollup)
	proof := bs.head.Proof(e.blockResolver)
	depositTxs := processDeposits(proof, b, e.blockResolver)
	newRollupState = executeTransactions(depositTxs, newRollupState)

	// Postprocessing - withdrawals
	withdrawals := rollupPostProcessingWithdrawals(bs.head, newRollupState)

	// Create a new rollup based on the proof of inclusion of the previous, including all new transactions
	r := NewRollup(b.Hash(), bs.head, bs.head.Header.Height+1, e.node, newRollupTxs, withdrawals, obscurocommon.GenerateNonce(), serialize(newRollupState))
	return &r
}

func (e *enclaveImpl) GetTransaction(txHash common.Hash) *nodecommon.L2Tx {
	// todo add some sort of cache
	rollup := e.storage.FetchHeadState().head

	var found bool
	for {
		txs := rollup.Transactions
		for _, tx := range txs {
			if tx.Hash() == txHash {
				return &tx
			}
		}
		rollup = e.storage.ParentRollup(rollup)
		rollup, found = e.storage.FetchRollup(rollup.Hash())
		if !found {
			panic(fmt.Sprintf("Could not find rollup: r_%s", rollup.Hash()))
		}
		if rollup.Header.Height == obscurocommon.L2GenesisHeight {
			return nil
		}
	}
}

func (e *enclaveImpl) Stop() error {
	e.exitCh <- true
	return nil
}

func (e *enclaveImpl) Attestation() obscurocommon.AttestationReport {
	// Todo
	return obscurocommon.AttestationReport{Owner: e.node}
}

// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
func (e *enclaveImpl) GenerateSecret() obscurocommon.EncryptedSharedEnclaveSecret {
	secret := make([]byte, 32)
	n, err := rand.Read(secret)
	if n != 32 || err != nil {
		panic(fmt.Sprintf("Could not generate secret: %s", err))
	}
	e.storage.StoreSecret(secret)
	return encryptSecret(secret)
}

// InitEnclave - initialise an enclave with a seed received by another enclave
func (e *enclaveImpl) InitEnclave(secret obscurocommon.EncryptedSharedEnclaveSecret) {
	e.storage.StoreSecret(decryptSecret(secret))
}

func (e *enclaveImpl) FetchSecret(obscurocommon.AttestationReport) obscurocommon.EncryptedSharedEnclaveSecret {
	return encryptSecret(e.storage.FetchSecret())
}

func (e *enclaveImpl) IsInitialised() bool {
	return e.storage.FetchSecret() != nil
}

func (e *enclaveImpl) noBlockStateBlockSubmissionResponse(block *types.Block) nodecommon.BlockSubmissionResponse {
	return nodecommon.BlockSubmissionResponse{
		L1Hash:            block.Hash(),
		L1Height:          e.blockResolver.HeightBlock(block),
		L1Parent:          block.ParentHash(),
		IngestedBlock:     true,
		IngestedNewRollup: false,
	}
}

func (e *enclaveImpl) blockStateBlockSubmissionResponse(bs *blockState, rollup nodecommon.ExtRollup) nodecommon.BlockSubmissionResponse {
	return nodecommon.BlockSubmissionResponse{
		L1Hash:            bs.block.Hash(),
		L1Height:          e.blockResolver.HeightBlock(bs.block),
		L1Parent:          bs.block.ParentHash(),
		L2Hash:            bs.head.Hash(),
		L2Height:          bs.head.Header.Height,
		L2Parent:          bs.head.Header.ParentHash,
		Withdrawals:       bs.head.Header.Withdrawals,
		ProducedRollup:    rollup,
		IngestedBlock:     true,
		IngestedNewRollup: bs.foundNewRollup,
	}
}

// Todo - implement with crypto
func decryptSecret(secret obscurocommon.EncryptedSharedEnclaveSecret) SharedEnclaveSecret {
	return SharedEnclaveSecret(secret)
}

// Todo - implement with crypto
func encryptSecret(secret SharedEnclaveSecret) obscurocommon.EncryptedSharedEnclaveSecret {
	return obscurocommon.EncryptedSharedEnclaveSecret(secret)
}

// internal structure to pass information.
type speculativeWork struct {
	r   *Rollup
	s   *State
	txs []nodecommon.L2Tx
}

// NewEnclave creates a new enclave.
// `genesisJSON` is the configuration for the corresponding L1's genesis block. This is used to validate the blocks
// received from the L1 node. If nil, incoming blocks are not validated.
func NewEnclave(id common.Address, mining bool, genesisJSON []byte, collector StatsCollector) nodecommon.Enclave {
	storage := NewStorage()

	var l1Blockchain *core.BlockChain
	if genesisJSON != nil {
		l1Blockchain = NewL1Blockchain(genesisJSON)
	}

	return &enclaveImpl{
		node:                 id,
		mining:               mining,
		storage:              storage,
		blockResolver:        storage,
		statsCollector:       collector,
		l1Blockchain:         l1Blockchain,
		txCh:                 make(chan nodecommon.L2Tx),
		roundWinnerCh:        make(chan *Rollup),
		exitCh:               make(chan bool),
		speculativeWorkInCh:  make(chan bool),
		speculativeWorkOutCh: make(chan speculativeWork),
	}
}
