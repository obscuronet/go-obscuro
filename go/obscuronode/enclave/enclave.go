package enclave

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/sql"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/config"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/evm"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethclient/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/db"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/mempool"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	obscurocore "github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"
)

const msgNoRollup = "could not fetch rollup"

type StatsCollector interface {
	// L2Recalc registers when a node has to discard the speculative work built on top of the winner of the gossip round.
	L2Recalc(id common.Address)
	RollupWithMoreRecentProof()
}

type enclaveImpl struct {
	config         config.EnclaveConfig
	nodeShortID    uint64
	storage        db.Storage
	blockResolver  db.BlockResolver
	mempool        mempool.Manager
	statsCollector StatsCollector
	l1Blockchain   *core.BlockChain

	txCh                 chan nodecommon.L2Tx
	roundWinnerCh        chan *obscurocore.Rollup
	exitCh               chan bool
	speculativeWorkInCh  chan bool
	speculativeWorkOutCh chan speculativeWork

	mgmtContractLib     mgmtcontractlib.MgmtContractLib
	erc20ContractLib    erc20contractlib.ERC20ContractLib
	attestationProvider AttestationProvider // interface for producing attestation reports and verifying them
	publicKeySerialized []byte
	privateKey          *rsa.PrivateKey

	blockProcessingMutex sync.Mutex
}

func (e *enclaveImpl) IsReady() error {
	return nil // The enclave is local so it is always ready
}

func (e *enclaveImpl) StopClient() error {
	return nil // The enclave is local so there is no client to stop
}

func (e *enclaveImpl) Start(block types.Block) {
	if e.config.SpeculativeExecution {
		// start the speculative rollup execution loop on its own go routine
		go e.start(block)
	}
}

func (e *enclaveImpl) start(block types.Block) {
	env := processingEnvironment{processedTxsMap: make(map[common.Hash]nodecommon.L2Tx)}
	// determine whether the block where the speculative execution will start already contains Obscuro state
	blockState, f := e.storage.FetchBlockState(block.Hash())
	if f {
		env.headRollup, _ = e.storage.FetchRollup(blockState.HeadRollup)
		if env.headRollup != nil {
			env.state = e.storage.CreateStateDB(env.headRollup.Hash())
		}
	}

	for {
		select {
		// A new winner was found after gossiping. Start speculatively executing incoming transactions to already have a rollup ready when the next round starts.
		case winnerRollup := <-e.roundWinnerCh:
			hash := winnerRollup.Hash()
			env.header = obscurocore.NewHeader(&hash, winnerRollup.Header.Number+1, e.config.HostID)
			env.headRollup = winnerRollup
			env.state = e.storage.CreateStateDB(winnerRollup.Hash())
			log.Trace(fmt.Sprintf(">   Agg%d: Create new speculative env  r_%d(%d).",
				e.nodeShortID,
				obscurocommon.ShortHash(winnerRollup.Header.Hash()),
				winnerRollup.Header.Number,
			))

			// determine the transactions that were not yet included
			env.processedTxs = currentTxs(winnerRollup, e.mempool.FetchMempoolTxs(), e.storage)
			env.processedTxsMap = makeMap(env.processedTxs)

			// calculate the State after executing them
			evm.ExecuteTransactions(env.processedTxs, env.state, env.headRollup.Header, e.storage, e.config.ObscuroChainID)

		case tx := <-e.txCh:
			// only process transactions if there is already a rollup to use as parent
			if env.headRollup != nil {
				_, found := env.processedTxsMap[tx.Hash()]
				if !found {
					env.processedTxsMap[tx.Hash()] = tx
					env.processedTxs = append(env.processedTxs, tx)
					evm.ExecuteTransactions([]nodecommon.L2Tx{tx}, env.state, env.header, e.storage, e.config.ObscuroChainID)
				}
			}

		case <-e.speculativeWorkInCh:
			if env.header == nil {
				e.speculativeWorkOutCh <- speculativeWork{found: false}
			} else {
				b := make([]nodecommon.L2Tx, 0, len(env.processedTxs))
				b = append(b, env.processedTxs...)
				e.speculativeWorkOutCh <- speculativeWork{
					found: true,
					r:     env.headRollup,
					s:     env.state,
					h:     env.header,
					txs:   b,
				}
			}

		case <-e.exitCh:
			return
		}
	}
}

func (e *enclaveImpl) ProduceGenesis(blkHash common.Hash) nodecommon.BlockSubmissionResponse {
	rolGenesis := obscurocore.NewRollup(blkHash, nil, obscurocommon.L2GenesisHeight, common.HexToAddress("0x0"), []nodecommon.L2Tx{}, []nodecommon.Withdrawal{}, obscurocommon.GenerateNonce(), common.BigToHash(big.NewInt(0)))
	return nodecommon.BlockSubmissionResponse{
		ProducedRollup: rolGenesis.ToExtRollup(),
		IngestedBlock:  true,
	}
}

// IngestBlocks is used to update the enclave with the full history of the L1 chain to date.
func (e *enclaveImpl) IngestBlocks(blocks []*types.Block) []nodecommon.BlockSubmissionResponse {
	result := make([]nodecommon.BlockSubmissionResponse, len(blocks))
	for i, block := range blocks {
		// We ignore a failure on the genesis block, since insertion of the genesis also produces a failure in Geth
		// (at least with Clique, where it fails with a `vote nonce not 0x00..0 or 0xff..f`).
		if ingestionFailedResponse := e.insertBlockIntoL1Chain(block); !e.isGenesisBlock(block) && ingestionFailedResponse != nil {
			result[i] = *ingestionFailedResponse
			return result // We return early, as all descendant blocks will also fail verification.
		}

		e.storage.StoreBlock(block)
		bs := updateState(block, e.blockResolver, e.mgmtContractLib, e.erc20ContractLib, e.storage, e.storage, e.nodeShortID, e.config.ObscuroChainID)
		if bs == nil {
			result[i] = e.noBlockStateBlockSubmissionResponse(block)
		} else {
			var rollup nodecommon.ExtRollup
			if bs.FoundNewRollup {
				hr, f := e.storage.FetchRollup(bs.HeadRollup)
				if !f {
					log.Panic(msgNoRollup)
				}

				rollup = hr.ToExtRollup()
			}
			result[i] = e.blockStateBlockSubmissionResponse(bs, rollup)
		}
	}

	return result
}

// SubmitBlock is used to update the enclave with an additional L1 block.
func (e *enclaveImpl) SubmitBlock(block types.Block) nodecommon.BlockSubmissionResponse {
	e.blockProcessingMutex.Lock()
	defer e.blockProcessingMutex.Unlock()

	// The genesis block should always be ingested, not submitted, so we ignore it if it's passed in here.
	if e.isGenesisBlock(&block) {
		return nodecommon.BlockSubmissionResponse{IngestedBlock: false, BlockNotIngestedCause: "Block was genesis block."}
	}

	_, foundBlock := e.storage.FetchBlock(block.Hash())
	if foundBlock {
		return nodecommon.BlockSubmissionResponse{IngestedBlock: false, BlockNotIngestedCause: "Block already ingested."}
	}

	if ingestionFailedResponse := e.insertBlockIntoL1Chain(&block); ingestionFailedResponse != nil {
		return *ingestionFailedResponse
	}

	_, f := e.storage.FetchBlock(block.Header().ParentHash)
	if !f && block.NumberU64() > obscurocommon.L1GenesisHeight {
		return nodecommon.BlockSubmissionResponse{IngestedBlock: false, BlockNotIngestedCause: "Block parent not stored."}
	}

	// Only store the block if the parent is available.
	stored := e.storage.StoreBlock(&block)
	if !stored {
		return nodecommon.BlockSubmissionResponse{IngestedBlock: false}
	}

	nodecommon.LogWithID(e.nodeShortID, "Update state: %d", obscurocommon.ShortHash(block.Hash()))
	blockState := updateState(&block, e.blockResolver, e.mgmtContractLib, e.erc20ContractLib, e.storage, e.storage, e.nodeShortID, e.config.ObscuroChainID)
	if blockState == nil {
		return e.noBlockStateBlockSubmissionResponse(&block)
	}

	// todo - A verifier node will not produce rollups, we can check the e.mining to get the node behaviour
	hr, f := e.storage.FetchRollup(blockState.HeadRollup)
	if !f {
		log.Panic(msgNoRollup)
	}
	e.mempool.RemoveMempoolTxs(historicTxs(hr, e.storage))
	r := e.produceRollup(&block, blockState)
	// todo - should store proposal rollups in a different storage as they are ephemeral (round based)
	e.storage.StoreRollup(r)

	nodecommon.LogWithID(e.nodeShortID, "Processed block: b_%d(%d)", obscurocommon.ShortHash(block.Hash()), block.NumberU64())

	return e.blockStateBlockSubmissionResponse(blockState, r.ToExtRollup())
}

func (e *enclaveImpl) SubmitRollup(rollup nodecommon.ExtRollup) {
	r := obscurocore.Rollup{
		Header:       rollup.Header,
		Transactions: obscurocore.DecryptTransactions(rollup.Txs),
	}

	// only store if the parent exists
	_, found := e.storage.FetchRollup(r.Header.ParentHash)
	if found {
		e.storage.StoreRollup(&r)
	} else {
		nodecommon.LogWithID(e.nodeShortID, "Received rollup with no parent: r_%d", obscurocommon.ShortHash(r.Hash()))
	}
}

func (e *enclaveImpl) SubmitTx(tx nodecommon.EncryptedTx) error {
	decryptedTx := obscurocore.DecryptTx(tx)
	err := verifySignature(e.config.ObscuroChainID, &decryptedTx)
	if err != nil {
		return err
	}
	e.mempool.AddMempoolTx(decryptedTx)
	if e.config.SpeculativeExecution {
		e.txCh <- decryptedTx
	}
	return nil
}

// Checks that the L2Tx has a valid signature.
func verifySignature(chainID int64, decryptedTx *nodecommon.L2Tx) error {
	signer := types.NewLondonSigner(big.NewInt(chainID))
	_, err := types.Sender(signer, decryptedTx)
	return err
}

func (e *enclaveImpl) RoundWinner(parent obscurocommon.L2RootHash) (nodecommon.ExtRollup, bool, error) {
	head, found := e.storage.FetchRollup(parent)
	if !found {
		return nodecommon.ExtRollup{}, false, fmt.Errorf("rollup not found: r_%s", parent)
	}

	nodecommon.LogWithID(e.nodeShortID, "Round winner height: %d", head.Header.Number)
	rollupsReceivedFromPeers := e.storage.FetchRollups(head.Header.Number + 1)
	// filter out rollups with a different Parent
	var usefulRollups []*obscurocore.Rollup
	for _, rol := range rollupsReceivedFromPeers {
		p := e.storage.ParentRollup(rol)
		if p == nil {
			nodecommon.LogWithID(e.nodeShortID, "Received rollup from peer but don't have parent rollup - discarding...")
			continue
		}
		if p.Hash() == head.Hash() {
			usefulRollups = append(usefulRollups, rol)
		}
	}

	parentState := e.storage.CreateStateDB(head.Hash())
	// determine the winner of the round
	winnerRollup, _ := e.findRoundWinner(usefulRollups, head, parentState, e.blockResolver, e.storage)
	if e.config.SpeculativeExecution {
		go e.notifySpeculative(winnerRollup)
	}

	// we are the winner
	if winnerRollup.Header.Agg == e.config.HostID {
		v := e.blockResolver.Proof(winnerRollup)
		w := e.storage.ParentRollup(winnerRollup)
		nodecommon.LogWithID(e.nodeShortID, "Publish rollup=r_%d(%d)[r_%d]{proof=b_%d(%d)}. Num Txs: %d. Txs: %v.  State=%v. ",
			obscurocommon.ShortHash(winnerRollup.Hash()), winnerRollup.Header.Number,
			obscurocommon.ShortHash(w.Hash()),
			obscurocommon.ShortHash(v.Hash()),
			v.NumberU64(),
			len(winnerRollup.Transactions),
			printTxs(winnerRollup.Transactions),
			winnerRollup.Header.State,
		)
		return winnerRollup.ToExtRollup(), true, nil
	}
	return nodecommon.ExtRollup{}, false, nil
}

func (e *enclaveImpl) notifySpeculative(winnerRollup *obscurocore.Rollup) {
	e.roundWinnerCh <- winnerRollup
}

func (e *enclaveImpl) Balance(address common.Address) uint64 {
	// todo user encryption
	hs := e.storage.FetchHeadState()
	if hs == nil {
		panic("Not initialised")
	}
	r, f := e.storage.FetchRollup(hs.HeadRollup)
	if !f {
		panic("not found")
	}
	s := e.storage.CreateStateDB(hs.HeadRollup)
	return evm.BalanceOfErc20(address, s, r.Header, e.storage, e.config.ObscuroChainID)
}

func (e *enclaveImpl) Nonce(address common.Address) uint64 {
	// todo user encryption
	hs := e.storage.FetchHeadState()
	if hs == nil {
		return 0
	}
	s := e.storage.CreateStateDB(hs.HeadRollup)
	return s.GetNonce(address)
}

func (e *enclaveImpl) produceRollup(b *types.Block, bs *obscurocore.BlockState) *obscurocore.Rollup {
	headRollup, f := e.storage.FetchRollup(bs.HeadRollup)
	if !f {
		log.Panic(msgNoRollup)
	}

	// These variables will be used to create the new rollup
	var newRollupTxs obscurocore.L2Txs
	var newRollupState *state.StateDB
	var newRollupHeader *nodecommon.Header

	/*
			speculativeExecutionSucceeded := false
		   todo - reenable
			if e.speculativeExecutionEnabled {
				// retrieve the speculatively calculated State based on the previous winner and the incoming transactions
				e.speculativeWorkInCh <- true
				speculativeRollup := <-e.speculativeWorkOutCh

				newRollupTxs = speculativeRollup.txs
				newRollupState = speculativeRollup.s
				newRollupHeader = speculativeRollup.h

				// the speculative execution has been processing on top of the wrong parent - due to failure in gossip or publishing to L1
				// or speculative execution is disabled
				speculativeExecutionSucceeded = speculativeRollup.found && (speculativeRollup.r.Hash() == bs.HeadRollup)

				if !speculativeExecutionSucceeded && speculativeRollup.r != nil {
					nodecommon.LogWithID(e.nodeShortID, "Recalculate. speculative=r_%d(%d), published=r_%d(%d)",
						obscurocommon.ShortHash(speculativeRollup.r.Hash()),
						speculativeRollup.r.Header.Number,
						obscurocommon.ShortHash(bs.HeadRollup),
						headRollup.Header.Number)
					if e.statsCollector != nil {
						e.statsCollector.L2Recalc(e.nodeID)
					}
				}
			}
	*/

	successfulTransactions := make([]nodecommon.L2Tx, 0)
	// if !speculativeExecutionSucceeded {
	// In case the speculative execution thread has not succeeded in producing a valid rollup
	// we have to create a new one from the mempool transactions
	newRollupHeader = obscurocore.NewHeader(&bs.HeadRollup, headRollup.Header.Number+1, e.config.HostID)
	newRollupTxs = currentTxs(headRollup, e.mempool.FetchMempoolTxs(), e.storage)
	newRollupState = e.storage.CreateStateDB(bs.HeadRollup)
	receipts := evm.ExecuteTransactions(newRollupTxs, newRollupState, newRollupHeader, e.storage, e.config.ObscuroChainID)
	// todo - only transactions that fail because of the nonce should be excluded
	for _, tx := range newRollupTxs {
		_, f := receipts[tx.Hash()]
		if f {
			successfulTransactions = append(successfulTransactions, tx)
		} else {
			log.Info(">   Agg%d: Excluding transaction %d", obscurocommon.ShortAddress(e.config.HostID), obscurocommon.ShortHash(tx.Hash()))
		}
	}

	// always process deposits last, either on top of the rollup produced speculatively or the newly created rollup
	// process deposits from the proof of the parent to the current block (which is the proof of the new rollup)
	proof := e.blockResolver.Proof(headRollup)
	depositTxs := extractDeposits(proof, b, e.blockResolver, e.erc20ContractLib, newRollupState, e.config.ObscuroChainID)
	depositReceipts := evm.ExecuteTransactions(depositTxs, newRollupState, newRollupHeader, e.storage, e.config.ObscuroChainID)
	for _, tx := range depositTxs {
		if depositReceipts[tx.Hash()] == nil {
			panic("Should not happen")
		}
	}

	// Create a new rollup based on the proof of inclusion of the previous, including all new transactions
	rootHash, err := newRollupState.Commit(true)
	if err != nil {
		return nil
	}
	r := obscurocore.NewRollupFromHeader(newRollupHeader, b.Hash(), successfulTransactions, obscurocommon.GenerateNonce(), rootHash)

	// Postprocessing - withdrawals
	r.Header.Withdrawals = e.rollupPostProcessingWithdrawals(&r, newRollupState, receipts)

	return &r
}

func (e *enclaveImpl) GetTransaction(txHash common.Hash) *nodecommon.L2Tx {
	// todo add some sort of cache
	hs := e.storage.FetchHeadState()
	if hs == nil {
		panic("should not happen")
	}
	rollup, found := e.storage.FetchRollup(hs.HeadRollup)
	if !found {
		log.Panic("could not fetch block's head rollup")
	}

	for {
		txs := rollup.Transactions
		for _, tx := range txs {
			if tx.Hash() == txHash {
				return &tx
			}
		}
		rollup = e.storage.ParentRollup(rollup)
		if rollup == nil || rollup.Header.Number == obscurocommon.L2GenesisHeight {
			return nil
		}
	}
}

func (e *enclaveImpl) Stop() error {
	if e.config.SpeculativeExecution {
		e.exitCh <- true
	}
	return nil
}

func (e *enclaveImpl) Attestation() *obscurocommon.AttestationReport {
	if e.publicKeySerialized == nil {
		panic("public key not initialized, we can't produce the attestation report")
	}
	report, err := e.attestationProvider.GetReport(e.publicKeySerialized, e.config.HostID)
	if err != nil {
		panic("Failed to produce remote report.")
	}
	return report
}

// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
func (e *enclaveImpl) GenerateSecret() obscurocommon.EncryptedSharedEnclaveSecret {
	secret := make([]byte, 32)
	n, err := rand.Read(secret)
	if n != 32 || err != nil {
		log.Panic("could not generate secret. Cause: %s", err)
	}
	e.storage.StoreSecret(secret)
	encSec, err := e.encryptSecret(e.publicKeySerialized, secret)
	if err != nil {
		log.Panic("failed to encrypt secret. Cause: %s", err)
	}
	return encSec
}

// InitEnclave - initialise an enclave with a seed received by another enclave
func (e *enclaveImpl) InitEnclave(s obscurocommon.EncryptedSharedEnclaveSecret) error {
	secret, err := e.decryptSecret(s)
	if err != nil {
		return err
	}
	e.storage.StoreSecret(secret)
	log.Trace(">   Agg%d: Secret decrypted and stored. Secret: %v", e.nodeShortID, secret)
	return nil
}

// ShareSecret verifies the request and if it trusts the report and the public key it will return the secret encrypted with that public key.
func (e *enclaveImpl) ShareSecret(att *obscurocommon.AttestationReport) (obscurocommon.EncryptedSharedEnclaveSecret, error) {
	// First we verify the attestation report has come from a valid obscuro enclave running in a verified TEE.
	data, err := e.attestationProvider.VerifyReport(att)
	if err != nil {
		return nil, err
	}
	// Then we verify the public key provided has come from the same enclave as that attestation report
	if err = verifyIdentity(data, att); err != nil {
		return nil, err
	}
	nodecommon.LogWithID(e.nodeShortID, "Successfully verified attestation and identity. Owner: %s", att.Owner)

	secret := e.storage.FetchSecret()
	if secret == nil {
		return nil, errors.New("secret was nil, no secret to share - this shouldn't happen")
	}
	return e.encryptSecret(att.PubKey, secret)
}

func verifyIdentity(data []byte, att *obscurocommon.AttestationReport) error {
	expectedIDHash := getIDHash(att.Owner, att.PubKey)
	// we trim the actual data because data extracted from the verified attestation is always 64 bytes long (padded with zeroes at the end)
	if !bytes.Equal(expectedIDHash, data[:len(expectedIDHash)]) {
		return fmt.Errorf("failed to verify hash for attestation report with owner: %s", att.Owner)
	}
	return nil
}

func (e *enclaveImpl) IsInitialised() bool {
	return e.storage.FetchSecret() != nil
}

func (e *enclaveImpl) isGenesisBlock(block *types.Block) bool {
	return e.l1Blockchain != nil && block.Hash() == e.l1Blockchain.Genesis().Hash()
}

// Inserts the block into the L1 chain if it exists and the block is not the genesis block. Returns a non-nil
// BlockSubmissionResponse if the insertion failed.
func (e *enclaveImpl) insertBlockIntoL1Chain(block *types.Block) *nodecommon.BlockSubmissionResponse {
	if e.l1Blockchain != nil {
		_, err := e.l1Blockchain.InsertChain(types.Blocks{block})
		if err != nil {
			causeMsg := fmt.Sprintf("Block was invalid: %v", err)
			return &nodecommon.BlockSubmissionResponse{IngestedBlock: false, BlockNotIngestedCause: causeMsg}
		}
	}
	return nil
}

func (e *enclaveImpl) noBlockStateBlockSubmissionResponse(block *types.Block) nodecommon.BlockSubmissionResponse {
	return nodecommon.BlockSubmissionResponse{
		BlockHeader:   block.Header(),
		IngestedBlock: true,
		FoundNewHead:  false,
	}
}

func (e *enclaveImpl) blockStateBlockSubmissionResponse(bs *obscurocore.BlockState, rollup nodecommon.ExtRollup) nodecommon.BlockSubmissionResponse {
	headRollup, f := e.storage.FetchRollup(bs.HeadRollup)
	if !f {
		log.Panic(msgNoRollup)
	}

	headBlock, f := e.storage.FetchBlock(bs.Block)
	if !f {
		log.Panic("could not fetch block")
	}

	var head *nodecommon.Header
	if bs.FoundNewRollup {
		head = headRollup.Header
	}
	return nodecommon.BlockSubmissionResponse{
		BlockHeader:    headBlock.Header(),
		ProducedRollup: rollup,
		IngestedBlock:  true,
		FoundNewHead:   bs.FoundNewRollup,
		RollupHead:     head,
	}
}

func generateKeyPair() *rsa.PrivateKey {
	// todo: This should be generated deterministically based on some enclave attributes if possible
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic("Failed to create RSA key")
	}
	return key
}

// Todo - implement with better crypto
func (e *enclaveImpl) decryptSecret(secret obscurocommon.EncryptedSharedEnclaveSecret) ([]byte, error) {
	if e.privateKey == nil {
		return nil, errors.New("private key not found - shouldn't happen")
	}
	return DecryptWithPrivateKey(secret, e.privateKey)
}

// Todo - implement with better crypto
func (e *enclaveImpl) encryptSecret(pubKeyEncoded []byte, secret obscurocore.SharedEnclaveSecret) (obscurocommon.EncryptedSharedEnclaveSecret, error) {
	nodecommon.LogWithID(e.nodeShortID, "Encrypting secret with public key %s", common.Bytes2Hex(pubKeyEncoded))
	key, err := x509.ParsePKCS1PublicKey(pubKeyEncoded)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key %w", err)
	}

	encKey, err := EncryptWithPublicKey(secret, key)
	if err != nil {
		nodecommon.LogWithID(e.nodeShortID, "Failed to encrypt key, err: %s\nsecret: %v\npubkey: %v\nencKey:%v", err, secret, pubKeyEncoded, encKey)
	}
	return encKey, err
}

// internal structure to pass information.
type speculativeWork struct {
	found bool
	r     *obscurocore.Rollup
	s     *state.StateDB
	h     *nodecommon.Header
	txs   []nodecommon.L2Tx
}

// internal structure used for the speculative execution.
type processingEnvironment struct {
	headRollup      *obscurocore.Rollup             // the current head rollup, which will be the parent of the new rollup
	header          *nodecommon.Header              // the header of the new rollup
	processedTxs    []nodecommon.L2Tx               // txs that were already processed
	processedTxsMap map[common.Hash]nodecommon.L2Tx // structure used to prevent duplicates
	state           *state.StateDB                  // the state as calculated from the previous rollup and the processed transactions
}

// NewEnclave creates a new enclave.
// `genesisJSON` is the configuration for the corresponding L1's genesis block. This is used to validate the blocks
// received from the L1 node if `validateBlocks` is set to true.
func NewEnclave(
	config config.EnclaveConfig,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	erc20ContractLib erc20contractlib.ERC20ContractLib,
	collector StatsCollector,
) nodecommon.Enclave {
	nodeShortID := obscurocommon.ShortAddress(config.HostID)

	connect := getDBConnector(config)
	backingDB, err := connect(nodeShortID)
	if err != nil {
		log.Panic("Failed to connect to backing database - %s", err)
	}
	storage := db.NewStorage(backingDB, nodeShortID)

	var l1Blockchain *core.BlockChain
	if config.ValidateL1Blocks {
		if config.GenesisJSON == nil {
			log.Panic("enclave is configured to validate blocks, but genesis JSON is nil")
		}
		l1Blockchain = NewL1Blockchain(config.GenesisJSON)
	} else {
		nodecommon.LogWithID(obscurocommon.ShortAddress(config.HostID), "validateBlocks is set to false. L1 blocks will not be validated.")
	}

	var attestationProvider AttestationProvider
	if config.WillAttest {
		attestationProvider = &EgoAttestationProvider{}
	} else {
		nodecommon.LogWithID(nodeShortID, "WARNING - Attestation is not enabled, enclave will not create a verified attestation report.")
		attestationProvider = &DummyAttestationProvider{}
	}

	nodecommon.LogWithID(nodeShortID, "Generating public key")
	privKey := generateKeyPair()
	serializedPubKey := x509.MarshalPKCS1PublicKey(&privKey.PublicKey)
	nodecommon.LogWithID(nodeShortID, "Generated public key %s", common.Bytes2Hex(serializedPubKey))

	return &enclaveImpl{
		config:               config,
		nodeShortID:          nodeShortID,
		storage:              storage,
		blockResolver:        storage,
		mempool:              mempool.New(),
		statsCollector:       collector,
		l1Blockchain:         l1Blockchain,
		txCh:                 make(chan nodecommon.L2Tx),
		roundWinnerCh:        make(chan *obscurocore.Rollup),
		exitCh:               make(chan bool),
		speculativeWorkInCh:  make(chan bool),
		speculativeWorkOutCh: make(chan speculativeWork),
		mgmtContractLib:      mgmtContractLib,
		erc20ContractLib:     erc20ContractLib,
		attestationProvider:  attestationProvider,
		privateKey:           privKey,
		publicKeySerialized:  serializedPubKey,
	}
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt with public key. %w", err)
	}
	return ciphertext, nil
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt with private key. %w", err)
	}
	return plaintext, nil
}

// getDBConnector creates an appropriate EthDBConnector function based on your config
func getDBConnector(cfg config.EnclaveConfig) db.EthDBConnector {
	if cfg.UseInMemoryDB {
		// not persistent
		return func(nodeID uint64) (ethdb.Database, error) {
			nodecommon.LogWithID(nodeID, "created in-memory database")
			return rawdb.NewMemoryDatabase(), nil
		}
	}

	if !cfg.WillAttest {
		// persistent but not secure in an enclave, we'll connect to a throwaway sqlite DB and test out persistence/sql implementations
		return sql.CreateTemporarySQLiteDB
	}

	// persistent and with attestation means connecting to edgeless DB in a trusted enclave from a secure enclave
	panic("Haven't implemented edgeless DB enclave connection yet")
}
