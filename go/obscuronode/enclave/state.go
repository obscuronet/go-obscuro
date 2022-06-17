package enclave

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/trie"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/evm"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/db"

	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethclient/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// Determine the new canonical L2 head and calculate the State
func updateState(
	b *types.Block,
	blockResolver db.BlockResolver,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	erc20ContractLib erc20contractlib.ERC20ContractLib,
	rollupResolver db.RollupResolver,
	bss db.BlockStateStorage,
	nodeID uint64,
	chainID int64,
	transactionBlobCrypto core.TransactionBlobCrypto,
) *core.BlockState {
	// This method is called recursively in case of Re-orgs. Stop when state was calculated already.
	val, found := bss.FetchBlockState(b.Hash())
	if found {
		return val
	}

	rollups := extractRollups(b, blockResolver, mgmtContractLib, nodeID, transactionBlobCrypto)
	genesisRollup := rollupResolver.FetchGenesisRollup()

	// processing blocks before genesis, so there is nothing to do
	if genesisRollup == nil && len(rollups) == 0 {
		return nil
	}

	// Detect if the incoming block contains the genesis rollup, and generate an updated state.
	// Handle the case of the block containing the genesis being processed multiple times.
	genesisState, isGenesis := handleGenesisRollup(b, rollups, genesisRollup, rollupResolver, bss, nodeID)
	if isGenesis {
		return genesisState
	}

	// To calculate the state after the current block, we need the state after the parent.
	// If this point is reached, there is a parent state guaranteed, because the genesis is handled above
	parentState, parentFound := bss.FetchBlockState(b.ParentHash())
	if !parentFound {
		// go back and calculate the Root of the Parent
		p, f := blockResolver.FetchBlock(b.ParentHash())
		if !f {
			nodecommon.LogWithID(nodeID, "Could not find block parent. This should not happen.")
			return nil
		}
		parentState = updateState(p, blockResolver, mgmtContractLib, erc20ContractLib, rollupResolver, bss, nodeID, chainID, transactionBlobCrypto)
	}

	if parentState == nil {
		nodecommon.LogWithID(nodeID, "Something went wrong. There should be a parent here, blockNum=%d. \n Block: %d, Block Parent: %d ",
			b.Number(),
			obscurocommon.ShortHash(b.Hash()),
			obscurocommon.ShortHash(b.Header().ParentHash),
		)
		return nil
	}

	bs, stateDB, head, receipts := calculateBlockState(b, parentState, blockResolver, rollups, erc20ContractLib, rollupResolver, bss, chainID)
	log.Trace(fmt.Sprintf(">   Agg%d: Calc block state b_%d: Found: %t - r_%d, ",
		nodeID,
		obscurocommon.ShortHash(b.Hash()),
		bs.FoundNewRollup,
		obscurocommon.ShortHash(bs.HeadRollup)))

	if bs.FoundNewRollup {
		// todo - root
		_, err := stateDB.Commit(true)
		if err != nil {
			log.Panic("could not commit new rollup to state DB. Cause: %s", err)
		}
	}
	bss.SaveNewHead(bs, head, receipts)

	return bs
}

func handleGenesisRollup(b *types.Block, rollups []*core.Rollup, genesisRollup *core.Rollup, resolver db.RollupResolver, bss db.BlockStateStorage, nodeID uint64) (genesisState *core.BlockState, isGenesis bool) {
	// the incoming block holds the genesis rollup
	// calculate and return the new block state
	// todo change this to an hardcoded hash on testnet/mainnet
	if genesisRollup == nil && len(rollups) == 1 {
		nodecommon.LogWithID(nodeID, "Found genesis rollup")

		genesis := rollups[0]
		resolver.StoreGenesisRollup(genesis)

		// The genesis rollup is part of the canonical chain and will be included in an L1 block by the first Aggregator.
		bs := core.BlockState{
			Block:          b.Hash(),
			HeadRollup:     genesis.Hash(),
			FoundNewRollup: true,
		}
		bss.SaveNewHead(&bs, genesis, nil)
		s := bss.GenesisStateDB()
		_, err := s.Commit(true)
		if err != nil {
			return nil, false
		}
		return &bs, true
	}

	// Re-processing the block that contains the rollup. This can happen as blocks can be fed to the enclave multiple times.
	// In this case we don't update the state and move on.
	if genesisRollup != nil && len(rollups) == 1 && rollups[0].Header.Hash() == genesisRollup.Hash() {
		return nil, true
	}
	return nil, false
}

// SortByNonce a very primitive way to implement mempool logic that
// adds transactions sorted by the nonce in the rollup
// which is what the EVM expects
type SortByNonce core.L2Txs

func (c SortByNonce) Len() int           { return len(c) }
func (c SortByNonce) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c SortByNonce) Less(i, j int) bool { return c[i].Nonce() < c[j].Nonce() }

// Calculate transactions to be included in the current rollup
func currentTxs(head *core.Rollup, mempool core.L2Txs, resolver db.RollupResolver) core.L2Txs {
	txs := findTxsNotIncluded(head, mempool, resolver)
	sort.Sort(SortByNonce(txs))
	return txs
}

func FindWinner(parent *core.Rollup, rollups []*core.Rollup, blockResolver db.BlockResolver) (*core.Rollup, bool) {
	win := -1
	// todo - add statistics to determine why there are conflicts.
	for i, r := range rollups {
		switch {
		case r.Header.ParentHash != parent.Hash(): // ignore rollups from L2 forks
		case r.Header.Number.Int64() <= parent.Header.Number.Int64(): // ignore rollups that are older than the parent
		case win == -1:
			win = i
		case blockResolver.ProofHeight(r) < blockResolver.ProofHeight(rollups[win]): // ignore rollups generated with an older proof
		case blockResolver.ProofHeight(r) > blockResolver.ProofHeight(rollups[win]): // newer rollups win
			win = i
		case r.Header.Nonce < rollups[win].Header.Nonce: // for rollups with the same proof, the one with the lowest nonce wins
			win = i
		}
	}
	if win == -1 {
		return nil, false
	}
	return rollups[win], true
}

func (e *enclaveImpl) findRoundWinner(receivedRollups []*core.Rollup, parent *core.Rollup, stateDB *state.StateDB, blockResolver db.BlockResolver, rollupResolver db.RollupResolver) (*core.Rollup, *state.StateDB) {
	headRollup, found := FindWinner(parent, receivedRollups, blockResolver)
	if !found {
		log.Panic("could not find winner. This should not happen for gossip rounds")
	}
	// calculate the state to compare with what is in the Rollup
	p := blockResolver.Proof(rollupResolver.ParentRollup(headRollup))
	depositTxs := extractDeposits(p, blockResolver.Proof(headRollup), blockResolver, e.erc20ContractLib, stateDB, e.config.ObscuroChainID)
	log.Info(fmt.Sprintf(">   Agg%d: Deposits:%d", obscurocommon.ShortAddress(e.config.HostID), len(depositTxs)))

	rootHash, txReceipts, depositReceipts := e.process(headRollup, stateDB, depositTxs)

	e.validateRollup(headRollup, rootHash, txReceipts, depositReceipts, stateDB)

	return headRollup, stateDB
}

func (e *enclaveImpl) process(rollup *core.Rollup, stateDB *state.StateDB, depositTxs []*nodecommon.L2Tx) (common.Hash, []*types.Receipt, []*types.Receipt) {
	txReceipts := evm.ExecuteTransactions(rollup.Transactions, stateDB, rollup.Header, e.storage, e.config.ObscuroChainID, 0)
	depositReceipts := evm.ExecuteTransactions(depositTxs, stateDB, rollup.Header, e.storage, e.config.ObscuroChainID, len(rollup.Transactions))
	rootHash, err := stateDB.Commit(true)
	if err != nil {
		log.Panic("could not commit to state DB. Cause: %s", err)
	}
	return rootHash, txReceipts, depositReceipts
}

func (e *enclaveImpl) validateRollup(rollup *core.Rollup, rootHash common.Hash, txReceipts []*types.Receipt, depositReceipts []*types.Receipt, stateDB *state.StateDB) {
	h := rollup.Header
	if rootHash != h.Root {
		// dump := stateDB.Dump(&state.DumpConfig{})
		dump := ""
		log.Panic("Calculated a different state. This should not happen as there are no malicious actors yet. \nGot: %s\nExp: %s\nHeight:%d\nTxs:%v\nState: %s",
			rootHash, h.Root, h.Number, printTxs(rollup.Transactions), dump)
	}

	//  check that the withdrawals in the header match the withdrawals as calculated
	withdrawals := e.rollupPostProcessingWithdrawals(rollup, stateDB, txReceipts)
	for i, w := range withdrawals {
		hw := h.Withdrawals[i]
		if hw.Amount != w.Amount || hw.Address != w.Address {
			log.Panic("Withdrawals don't match")
		}
	}

	rec := getReceipts(txReceipts, depositReceipts)
	rbloom := types.CreateBloom(rec)
	if rbloom != h.Bloom {
		log.Panic("invalid bloom (remote: %x  local: %x)", h.Bloom, rbloom)
	}

	receiptSha := types.DeriveSha(rec, trie.NewStackTrie(nil))
	if receiptSha != h.ReceiptHash {
		log.Panic("invalid receipt root hash (remote: %x local: %x)", h.ReceiptHash, receiptSha)
	}
}

func toReceiptMap(txReceipts []*types.Receipt) map[common.Hash]*types.Receipt {
	result := make(map[common.Hash]*types.Receipt, 0)
	for _, r := range txReceipts {
		result[r.TxHash] = r
	}
	return result
}

func getReceipts(txReceipts []*types.Receipt, depositReceipts []*types.Receipt) types.Receipts {
	receipts := make([]*types.Receipt, 0)
	receipts = append(receipts, txReceipts...)
	receipts = append(receipts, depositReceipts...)
	return receipts
}

// returns a list of L2 deposit transactions generated from the L1 deposit transactions
// starting with the proof of the parent rollup(exclusive) to the proof of the current rollup
func extractDeposits(
	fromBlock *types.Block,
	toBlock *types.Block,
	blockResolver db.BlockResolver,
	erc20ContractLib erc20contractlib.ERC20ContractLib,
	rollupState *state.StateDB,
	chainID int64,
) []*nodecommon.L2Tx {
	from := obscurocommon.GenesisBlock.Hash()
	height := obscurocommon.L1GenesisHeight
	if fromBlock != nil {
		from = fromBlock.Hash()
		height = fromBlock.NumberU64()
		if !blockResolver.IsAncestor(toBlock, fromBlock) {
			log.Panic("Deposits can't be processed because the rollups are not on the same Ethereum fork. This should not happen.")
		}
	}

	allDeposits := make([]*nodecommon.L2Tx, 0)
	b := toBlock
	for {
		if b.Hash() == from {
			break
		}
		for _, tx := range b.Transactions() {
			t := erc20ContractLib.DecodeTx(tx)
			if t == nil {
				continue
			}

			if depositTx, ok := t.(*obscurocommon.L1DepositTx); ok {
				depL2Tx := newDepositTx(*depositTx.Sender, depositTx.Amount, rollupState, uint64(len(allDeposits)), chainID)
				allDeposits = append(allDeposits, depL2Tx)
			}
		}
		if b.NumberU64() < height {
			log.Panic("block height is less than genesis height")
		}
		p, f := blockResolver.ParentBlock(b)
		if !f {
			log.Panic("deposits can't be processed because the rollups are not on the same Ethereum fork")
		}
		b = p
	}

	return allDeposits
}

// given an L1 block, and the State as it was in the Parent block, calculates the State after the current block.
func calculateBlockState(b *types.Block, parentState *core.BlockState, blockResolver db.BlockResolver, rollups []*core.Rollup, erc20ContractLib erc20contractlib.ERC20ContractLib, rollupResolver db.RollupResolver, bss db.BlockStateStorage, chainID int64) (*core.BlockState, *state.StateDB, *core.Rollup, []*types.Receipt) {
	currentHead, found := rollupResolver.FetchRollup(parentState.HeadRollup)
	if !found {
		log.Panic("could not fetch parent rollup")
	}
	newHeadRollup, found := FindWinner(currentHead, rollups, blockResolver)
	stateDB := bss.CreateStateDB(parentState.HeadRollup)
	var receipts []*types.Receipt
	// only change the state if there is a new l2 HeadRollup in the current block
	if found {
		// Preprocessing before passing to the vm
		// todo transform into an eth block structure
		parentRollup := rollupResolver.ParentRollup(newHeadRollup)
		p := blockResolver.Proof(parentRollup)
		depositTxs := extractDeposits(p, blockResolver.Proof(newHeadRollup), blockResolver, erc20ContractLib, stateDB, chainID)

		// deposits have to be processed after the normal transactions were executed because during speculative execution they are not available
		txsToProcess := append(newHeadRollup.Transactions, depositTxs...)
		receipts = evm.ExecuteTransactions(txsToProcess, stateDB, newHeadRollup.Header, rollupResolver, chainID, 0)
		// todo - handle failure , which means a new winner must be selected
	} else {
		newHeadRollup = currentHead
	}

	bs := core.BlockState{
		Block:          b.Hash(),
		HeadRollup:     newHeadRollup.Hash(),
		FoundNewRollup: found,
	}
	return &bs, stateDB, newHeadRollup, receipts
}

// Todo - this has to be implemented differently based on how we define the ObsERC20
func (e *enclaveImpl) rollupPostProcessingWithdrawals(newHeadRollup *core.Rollup, state *state.StateDB, receipts []*types.Receipt) []nodecommon.Withdrawal {
	receiptsMap := toReceiptMap(receipts)
	w := make([]nodecommon.Withdrawal, 0)
	// go through each transaction and check if the withdrawal was processed correctly
	for _, t := range newHeadRollup.Transactions {
		found, address, amount := erc20contractlib.DecodeTransferTx(t)

		if found && *t.To() == evm.Erc20ContractAddress && *address == evm.WithdrawalAddress {
			receipt := receiptsMap[t.Hash()]
			if receipt != nil && receipt.Status == types.ReceiptStatusSuccessful {
				signer := types.NewLondonSigner(big.NewInt(e.config.ObscuroChainID))
				from, err := types.Sender(signer, t)
				if err != nil {
					panic(err)
				}
				state.Logs()
				w = append(w, nodecommon.Withdrawal{
					Amount:  amount.Uint64(),
					Address: from,
				})
			}
		}
	}

	// TODO - fix the withdrawals logic
	// clearWithdrawals(state, withdrawalTxs)
	return w
}

func extractRollups(b *types.Block, blockResolver db.BlockResolver, mgmtContractLib mgmtcontractlib.MgmtContractLib, nodeID uint64, transactionBlobCrypto core.TransactionBlobCrypto) []*core.Rollup {
	rollups := make([]*core.Rollup, 0)
	for _, tx := range b.Transactions() {
		// go through all rollup transactions
		t := mgmtContractLib.DecodeTx(tx)
		if t == nil {
			continue
		}

		if rolTx, ok := t.(*obscurocommon.L1RollupTx); ok {
			r := nodecommon.DecodeRollupOrPanic(rolTx.Rollup)

			// Ignore rollups created with proofs from different L1 blocks
			// In case of L1 reorgs, rollups may end published on a fork
			if blockResolver.IsBlockAncestor(b, r.Header.L1Proof) {
				rollups = append(rollups, toEnclaveRollup(r, transactionBlobCrypto))
				nodecommon.LogWithID(nodeID, "Extracted Rollup r_%d from block b_%d",
					obscurocommon.ShortHash(r.Hash()),
					obscurocommon.ShortHash(b.Hash()),
				)
			}
		}
	}
	return rollups
}

func toEnclaveRollup(r *nodecommon.Rollup, transactionBlobCrypto core.TransactionBlobCrypto) *core.Rollup {
	return &core.Rollup{
		Header:       r.Header,
		Transactions: transactionBlobCrypto.Decrypt(r.Transactions),
	}
}

func newDepositTx(address common.Address, amount uint64, rollupState *state.StateDB, adjustNonce uint64, chainID int64) *nodecommon.L2Tx {
	transferERC20data := erc20contractlib.CreateTransferTxData(address, amount)
	signer := types.NewLondonSigner(big.NewInt(chainID))

	// The nonce is adjusted with the number of deposits added to the rollup already.
	nonce := rollupState.GetNonce(evm.Erc20OwnerAddress) + adjustNonce

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		Value:    common.Big0,
		Gas:      1_000_000,
		GasPrice: common.Big0,
		Data:     transferERC20data,
		To:       &evm.Erc20ContractAddress,
	})

	newTx, err := types.SignTx(tx, signer, evm.Erc20OwnerKey)
	if err != nil {
		log.Panic("could not encode L2 transaction data. Cause: %s", err)
	}
	return newTx
}
