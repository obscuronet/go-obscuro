package simulation

import (
	cryptorand "crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/obscuronet/obscuro-playground/go/common/log"

	"github.com/obscuronet/obscuro-playground/go/enclave/bridge"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/common"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/integration/erc20contract"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"golang.org/x/sync/errgroup"

	"github.com/obscuronet/obscuro-playground/go/rpcclientlib"
	"github.com/obscuronet/obscuro-playground/integration"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethadapter"
	"github.com/obscuronet/obscuro-playground/go/ethadapter/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/wallet"
	simstats "github.com/obscuronet/obscuro-playground/integration/simulation/stats"
)

const (
	timeoutMillis = 30000 // The timeout in millis to wait for an updated nonce for a wallet.
	// EnclavePublicKeyHex is the public key of the enclave.
	// TODO - Retrieve this key from the management contract instead.
	EnclavePublicKeyHex = "034d3b7e63a8bcd532ee3d1d6ecad9d67fca7821981a044551f0f0cbec74d0bc5e"
)

// TransactionInjector is a structure that generates, issues and tracks transactions
type TransactionInjector struct {
	// counters
	Counter *txInjectorCounter
	stats   *simstats.Stats

	// settings
	avgBlockDuration time.Duration

	// wallets
	wallets *params.SimWallets

	// connections
	l1Clients []ethadapter.EthClient
	l2Clients []rpcclientlib.Client

	// addrs and libs
	mgmtContractAddr *gethcommon.Address
	mgmtContractLib  mgmtcontractlib.MgmtContractLib
	erc20ContractLib erc20contractlib.ERC20ContractLib

	// controls
	interruptRun     *int32
	fullyStoppedChan chan bool

	enclavePublicKey *ecies.PublicKey

	// The number of transactions of each type to issue, or 0 for unlimited transactions
	txsToIssue int
}

// NewTransactionInjector returns a transaction manager with a given number of obsWallets
// todo Add methods that generate deterministic scenarios
func NewTransactionInjector(
	avgBlockDuration time.Duration,
	stats *simstats.Stats,
	l1Nodes []ethadapter.EthClient,
	wallets *params.SimWallets,
	mgmtContractAddr *gethcommon.Address,
	l2NodeClients []rpcclientlib.Client,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	erc20ContractLib erc20contractlib.ERC20ContractLib,
	txsToIssue int,
) *TransactionInjector {
	interrupt := int32(0)

	// We retrieve the enclave public key to encrypt transactions.
	enclavePublicKey, err := crypto.DecompressPubkey(gethcommon.Hex2Bytes(EnclavePublicKeyHex))
	if err != nil {
		panic(fmt.Errorf("could not decompress enclave public key from hex. Cause: %w", err))
	}
	enclavePublicKeyEcies := ecies.ImportECDSAPublic(enclavePublicKey)

	return &TransactionInjector{
		avgBlockDuration: avgBlockDuration,
		stats:            stats,
		l1Clients:        l1Nodes,
		l2Clients:        l2NodeClients,
		interruptRun:     &interrupt,
		fullyStoppedChan: make(chan bool, 1),
		mgmtContractAddr: mgmtContractAddr,
		mgmtContractLib:  mgmtContractLib,
		erc20ContractLib: erc20ContractLib,
		wallets:          wallets,
		Counter:          newCounter(),
		enclavePublicKey: enclavePublicKeyEcies,
		txsToIssue:       txsToIssue,
	}
}

// Start begins the execution on the TransactionInjector
// Deposits an initial balance in to each wallet
// Generates and issues L1 and L2 transactions to the network
func (ti *TransactionInjector) Start() {
	ti.deployObscuroERC20(ti.wallets.Tokens[bridge.BTC].L2Owner)
	ti.deployObscuroERC20(ti.wallets.Tokens[bridge.ETH].L2Owner)

	// enough time to process everywhere
	time.Sleep(ti.avgBlockDuration * 6)

	// deposit some initial amount into every simulation wallet
	for _, w := range ti.wallets.SimEthWallets {
		addr := w.Address()
		txData := &ethadapter.L1DepositTx{
			Amount:        initialBalance,
			To:            ti.mgmtContractAddr,
			TokenContract: ti.wallets.Tokens[bridge.BTC].L1ContractAddress,
			Sender:        &addr,
		}
		tx := ti.erc20ContractLib.CreateDepositTx(txData, w.GetNonceAndIncrement())
		signedTx, err := w.SignTransaction(tx)
		if err != nil {
			panic(err)
		}
		err = ti.rndL1NodeClient().SendTransaction(signedTx)
		if err != nil {
			panic(err)
		}

		ti.stats.Deposit(initialBalance)
		go ti.Counter.trackL1Tx(txData)
	}

	// start transactions issuance
	var wg errgroup.Group
	wg.Go(func() error {
		ti.issueRandomDeposits()
		return nil
	})

	wg.Go(func() error {
		ti.issueRandomWithdrawals()
		return nil
	})

	wg.Go(func() error {
		ti.issueRandomTransfers()
		return nil
	})

	wg.Go(func() error {
		ti.issueInvalidL2Txs()
		return nil
	})

	_ = wg.Wait() // future proofing to return errors
	ti.fullyStoppedChan <- true
}

// This deploys an ERC20 contract on Obscuro, which is used for token arithmetic.
func (ti *TransactionInjector) deployObscuroERC20(owner wallet.Wallet) {
	// deploy the ERC20
	contractBytes := gethcommon.Hex2Bytes(erc20contract.ContractByteCode)
	deployContractTx := types.LegacyTx{
		Nonce:    NextNonce(ti.l2Clients[0], owner),
		Gas:      1025_000_000,
		GasPrice: gethcommon.Big0,
		Data:     contractBytes,
	}
	signedTx, err := owner.SignTransaction(&deployContractTx)
	if err != nil {
		panic(err)
	}
	encryptedTx, err := EncryptTx(signedTx, ti.enclavePublicKey)
	if err != nil {
		panic(err)
	}
	err = ti.rndL2NodeClient().Call(nil, rpcclientlib.RPCSendTransactionEncrypted, encryptedTx)
	if err != nil {
		panic(err)
	}
}

func (ti *TransactionInjector) Stop() {
	atomic.StoreInt32(ti.interruptRun, 1)
	for range ti.fullyStoppedChan {
		log.Info("TransactionInjector stopped successfully")
		return
	}
}

// issueRandomTransfers creates and issues a number of L2 transfer transactions proportional to the simulation time, such that they can be processed
func (ti *TransactionInjector) issueRandomTransfers() {
	for txCounter := 0; ti.shouldKeepIssuing(txCounter); txCounter++ {
		fromWallet := ti.rndObsWallet()
		toWallet := ti.rndObsWallet()
		// We avoid transfers to self, unless there is only a single L2 wallet.
		for len(ti.wallets.SimObsWallets) > 1 && fromWallet.Address().Hex() == toWallet.Address().Hex() {
			toWallet = ti.rndObsWallet()
		}
		tx := ti.newObscuroTransferTx(fromWallet, toWallet.Address(), common.RndBtw(1, 500), ti.l2Clients[0])
		signedTx, err := fromWallet.SignTransaction(tx)
		if err != nil {
			panic(err)
		}
		log.Info(
			"Transfer transaction injected into L2. Hash: %d. From address: %d. To address: %d",
			common.ShortHash(signedTx.Hash()),
			common.ShortAddress(fromWallet.Address()),
			common.ShortAddress(toWallet.Address()),
		)

		encryptedTx, err := EncryptTx(signedTx, ti.enclavePublicKey)
		if err != nil {
			panic(err)
		}
		ti.stats.Transfer()

		err = ti.rndL2NodeClient().Call(nil, rpcclientlib.RPCSendTransactionEncrypted, encryptedTx)
		if err != nil {
			log.Info("Failed to issue transfer via RPC. Cause: %s", err)
			continue
		}

		go ti.Counter.trackTransferL2Tx(signedTx)
		common.SleepRndBtw(ti.avgBlockDuration/4, ti.avgBlockDuration)
	}
}

// issueRandomDeposits creates and issues a number of transactions proportional to the simulation time, such that they can be processed
func (ti *TransactionInjector) issueRandomDeposits() {
	for txCounter := 0; ti.shouldKeepIssuing(txCounter); txCounter++ {
		v := common.RndBtw(1, 100)
		ethWallet := ti.rndEthWallet()
		addr := ethWallet.Address()
		txData := &ethadapter.L1DepositTx{
			Amount:        v,
			To:            ti.mgmtContractAddr,
			TokenContract: ti.wallets.Tokens[bridge.BTC].L1ContractAddress,
			Sender:        &addr,
		}
		tx := ti.erc20ContractLib.CreateDepositTx(txData, ethWallet.GetNonceAndIncrement())
		signedTx, err := ethWallet.SignTransaction(tx)
		if err != nil {
			panic(err)
		}
		log.Info(
			"Deposit transaction injected into L1. Hash: %d. From address: %d",
			common.ShortHash(signedTx.Hash()),
			common.ShortAddress(ethWallet.Address()),
		)
		err = ti.rndL1NodeClient().SendTransaction(signedTx)
		if err != nil {
			panic(err)
		}

		ti.stats.Deposit(v)
		go ti.Counter.trackL1Tx(txData)
		common.SleepRndBtw(ti.avgBlockDuration, ti.avgBlockDuration*2)
	}
}

// issueRandomWithdrawals creates and issues a number of transactions proportional to the simulation time, such that they can be processed
func (ti *TransactionInjector) issueRandomWithdrawals() {
	for txCounter := 0; ti.shouldKeepIssuing(txCounter); txCounter++ {
		v := common.RndBtw(1, 100)
		obsWallet := ti.rndObsWallet()
		// todo - random client
		tx := ti.newObscuroWithdrawalTx(obsWallet, v, ti.l2Clients[0])
		signedTx, err := obsWallet.SignTransaction(tx)
		if err != nil {
			panic(err)
		}
		log.Info(
			"Withdrawal transaction injected into L2. Hash: %d. From address: %d",
			common.ShortHash(signedTx.Hash()),
			common.ShortAddress(obsWallet.Address()),
		)
		encryptedTx, err := EncryptTx(signedTx, ti.enclavePublicKey)
		if err != nil {
			panic(err)
		}

		err = ti.rndL2NodeClient().Call(nil, rpcclientlib.RPCSendTransactionEncrypted, encryptedTx)
		if err != nil {
			log.Info("Failed to issue withdrawal via RPC. Cause: %s", err)
			continue
		}

		ti.stats.Withdrawal(v)
		go ti.Counter.trackWithdrawalL2Tx(signedTx)
		common.SleepRndBtw(ti.avgBlockDuration, ti.avgBlockDuration*2)
	}
}

// issueInvalidL2Txs creates and issues invalidly-signed L2 transactions proportional to the simulation time.
// These transactions should be rejected by the nodes, and thus we expect them to not affect the simulation
func (ti *TransactionInjector) issueInvalidL2Txs() {
	for txCounter := 0; ti.shouldKeepIssuing(txCounter); txCounter++ {
		fromWallet := ti.rndObsWallet()
		toWallet := ti.rndObsWallet()
		// We avoid transfers to self, unless there is only a single L2 wallet.
		for len(ti.wallets.SimObsWallets) > 1 && fromWallet.Address().Hex() == toWallet.Address().Hex() {
			toWallet = ti.rndObsWallet()
		}
		tx := ti.newCustomObscuroWithdrawalTx(common.RndBtw(1, 100))

		signedTx := ti.createInvalidSignage(tx, fromWallet)
		encryptedTx, err := EncryptTx(signedTx, ti.enclavePublicKey)
		if err != nil {
			panic(err)
		}

		err = ti.rndL2NodeClient().Call(nil, rpcclientlib.RPCSendTransactionEncrypted, encryptedTx)
		if err != nil {
			log.Info("Failed to issue withdrawal via RPC. Cause: %s", err)
		}
		time.Sleep(common.RndBtwTime(ti.avgBlockDuration/4, ti.avgBlockDuration))
	}
}

// Uses one of the approaches to create an invalidly-signed transaction.
func (ti *TransactionInjector) createInvalidSignage(tx types.TxData, w wallet.Wallet) *types.Transaction {
	switch rand.Intn(2) { //nolint:gosec
	case 0: // We sign the transaction with a bad signer.
		incorrectChainID := int64(integration.EthereumChainID + 1)
		signer := types.NewLondonSigner(big.NewInt(incorrectChainID))
		signedTx, _ := types.SignNewTx(w.PrivateKey(), signer, tx)
		return signedTx

	case 1: // We do not sign the transaction.
		return types.NewTx(tx)
	}
	return nil
}

func (ti *TransactionInjector) rndObsWallet() wallet.Wallet {
	return ti.wallets.SimObsWallets[rand.Intn(len(ti.wallets.SimObsWallets))] //nolint:gosec
}

func (ti *TransactionInjector) rndEthWallet() wallet.Wallet {
	return ti.wallets.SimEthWallets[rand.Intn(len(ti.wallets.SimEthWallets))] //nolint:gosec
}

func (ti *TransactionInjector) rndL1NodeClient() ethadapter.EthClient {
	return ti.l1Clients[rand.Intn(len(ti.l1Clients))] //nolint:gosec
}

func (ti *TransactionInjector) rndL2NodeClient() rpcclientlib.Client {
	return ti.l2Clients[rand.Intn(len(ti.l2Clients))] //nolint:gosec
}

func (ti *TransactionInjector) newObscuroTransferTx(from wallet.Wallet, dest gethcommon.Address, amount uint64, client rpcclientlib.Client) types.TxData {
	data := erc20contractlib.CreateTransferTxData(dest, amount)
	t := ti.newTx(data, NextNonce(client, from))
	return t
}

func (ti *TransactionInjector) newObscuroWithdrawalTx(from wallet.Wallet, amount uint64, client rpcclientlib.Client) types.TxData {
	transferERC20data := erc20contractlib.CreateTransferTxData(bridge.BridgeAddress, amount)
	t := ti.newTx(transferERC20data, NextNonce(client, from))
	return t
}

func (ti *TransactionInjector) newCustomObscuroWithdrawalTx(amount uint64) types.TxData {
	transferERC20data := erc20contractlib.CreateTransferTxData(bridge.BridgeAddress, amount)
	return ti.newTx(transferERC20data, 1)
}

func (ti *TransactionInjector) newTx(data []byte, nonce uint64) types.TxData {
	return &types.LegacyTx{
		Nonce:    nonce,
		Value:    gethcommon.Big0,
		Gas:      1_000_000,
		GasPrice: gethcommon.Big0,
		Data:     data,
		To:       ti.wallets.Tokens[bridge.BTC].L2ContractAddress,
	}
}

func readNonce(cl rpcclientlib.Client, a gethcommon.Address) uint64 {
	reqParams := []string{a.String()}

	jsonParams, err := json.Marshal(reqParams)
	if err != nil {
		panic(err)
	}

	// we don't encode the parameter here as currently the enclave supports unencrypted communication, and we haven't set up a viewing key
	var result string
	err = cl.Call(&result, rpcclientlib.RPCGetTxCount, jsonParams)
	if err != nil {
		panic(err)
	}
	nonce, err := hexutil.DecodeUint64(result)
	if err != nil {
		panic(err)
	}
	return nonce
}

func NextNonce(cl rpcclientlib.Client, w wallet.Wallet) uint64 {
	counter := 0

	// only returns the nonce when the previous transaction was recorded
	for {
		remoteNonce := readNonce(cl, w.Address())
		localNonce := w.GetNonce()
		if remoteNonce == localNonce {
			return w.GetNonceAndIncrement()
		}
		if remoteNonce > localNonce {
			panic("remote nonce exceeds local nonce")
		}

		counter++
		if counter > timeoutMillis {
			panic("transaction injector failed to retrieve nonce after thirty seconds")
		}
		time.Sleep(time.Millisecond)
	}
}

// EncryptTx encrypts a single transaction using the enclave's public key to send it privately to the enclave.
func EncryptTx(tx *common.L2Tx, enclavePublicKey *ecies.PublicKey) (common.EncryptedTx, error) {
	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, fmt.Errorf("could not encode transaction bytes with RLP. Cause: %w", err)
	}

	encryptedTxBytes, err := ecies.Encrypt(cryptorand.Reader, enclavePublicKey, txBytes, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt request params with enclave public key. Cause: %w", err)
	}

	return encryptedTxBytes, nil
}

// Indicates whether to keep issuing transactions, or halt.
func (ti *TransactionInjector) shouldKeepIssuing(txCounter int) bool {
	isInterrupted := atomic.LoadInt32(ti.interruptRun) != 0

	// 0 is a special value indicating we should only stop issuing transactions when interrupted.
	if ti.txsToIssue == 0 {
		return !isInterrupted
	}

	return !isInterrupted && txCounter < ti.txsToIssue
}
