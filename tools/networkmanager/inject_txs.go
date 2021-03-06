package networkmanager

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/obscuronet/go-obscuro/integration/common/viewkey"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/integration/simulation/network"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"

	"github.com/obscuronet/go-obscuro/integration/simulation/stats"

	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/erc20contractlib"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/go/rpcclientlib"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration/simulation"
)

func InjectTransactions(cfg Config, args []string) {
	println("Connecting to L1 node...")
	l1Client, err := ethadapter.NewEthClient(cfg.l1NodeHost, cfg.l1NodeWebsocketPort, cfg.l1ConnectionTimeout, common.HexToAddress("0x0"))
	if err != nil {
		panic(fmt.Sprintf("could not create L1 client. Cause: %s", err))
	}
	println("Connecting to Obscuro node...")
	l2Client := rpcclientlib.NewClient(cfg.obscuroClientAddress)

	// We store the block at which we start injecting transactions.
	startBlock := l1Client.FetchHeadBlock()

	simStats := stats.NewStats(0)
	mgmtContractLib := mgmtcontractlib.NewMgmtContractLib(&cfg.mgmtContractAddress)
	erc20ContractLib := erc20contractlib.NewERC20ContractLib(&cfg.mgmtContractAddress, &cfg.erc20ContractAddress)
	avgBlockDuration := time.Second

	wallets := createWallets(cfg, l1Client, l2Client)
	walletClients := createWalletRPCClients(wallets, cfg.obscuroClientAddress)

	netwClients := &network.RPCHandles{
		EthClients:                    []ethadapter.EthClient{l1Client},
		ObscuroClients:                []rpcclientlib.Client{l2Client},
		VirtualWalletExtensionClients: walletClients,
	}

	txInjector := simulation.NewTransactionInjector(
		avgBlockDuration,
		simStats,
		netwClients,
		wallets,
		&cfg.mgmtContractAddress,
		mgmtContractLib,
		erc20ContractLib,
		parseNumOfTxs(args),
	)

	println("Injecting transactions into network...")
	txInjector.Start()
	time.Sleep(5 * avgBlockDuration) // We sleep to allow transactions to propagate.

	println(fmt.Sprintf("Stopped injecting transactions into network.\n"+
		"Attempted to inject %d L1 transactions, %d L2 transfer transactions, and %d L2 withdrawal transactions.",
		len(txInjector.TxTracker.L1Transactions), len(txInjector.TxTracker.TransferL2Transactions), len(txInjector.TxTracker.WithdrawalL2Transactions),
	))

	checkDepositsSuccessful(txInjector, l1Client, simStats, erc20ContractLib, mgmtContractLib, startBlock)
	checkL2TxsSuccessful(l2Client, txInjector)

	os.Exit(0)
}

func createWalletRPCClients(wallets *params.SimWallets, obscuroNodeAddr string) map[string]rpcclientlib.Client {
	clients := make(map[string]rpcclientlib.Client)

	for _, w := range wallets.SimObsWallets {
		clients[w.Address().String()] = rpcclientlib.NewClient(obscuroNodeAddr)
		viewkey.GenerateAndRegisterViewingKey(clients[w.Address().String()], w)
	}
	for _, t := range wallets.Tokens {
		w := t.L2Owner
		clients[w.Address().String()] = rpcclientlib.NewClient(obscuroNodeAddr)
		viewkey.GenerateAndRegisterViewingKey(clients[w.Address().String()], w)
	}

	return clients
}

func createWallets(nmConfig Config, l1Client ethadapter.EthClient, l2Client rpcclientlib.Client) *params.SimWallets {
	wallets := params.NewSimWallets(len(nmConfig.privateKeys), 0, nmConfig.l1ChainID, nmConfig.obscuroChainID)

	// We override the autogenerated Ethereum wallets with ones using the provided private keys.
	wallets.SimEthWallets = make([]wallet.Wallet, len(nmConfig.privateKeys))
	for idx, privateKeyString := range nmConfig.privateKeys {
		privateKey, err := crypto.HexToECDSA(privateKeyString)
		if err != nil {
			panic(fmt.Errorf("could not recover private key from hex. Cause: %w", err))
		}
		l1Wallet := wallet.NewInMemoryWalletFromPK(big.NewInt(nmConfig.l1ChainID), privateKey)
		wallets.SimEthWallets[idx] = l1Wallet
	}

	// We update the L1 and L2 wallet nonces.
	for _, l1Wallet := range wallets.AllEthWallets() {
		nonce, err := l1Client.Nonce(l1Wallet.Address())
		if err != nil {
			panic(fmt.Errorf("could not set L1 wallet nonce. Cause: %w", err))
		}
		l1Wallet.SetNonce(nonce)
	}
	for _, l2Wallet := range wallets.AllObsWallets() {
		var nonce uint64
		err := l2Client.Call(&nonce, rpcclientlib.RPCNonce, l2Wallet.Address())
		if err != nil {
			panic(fmt.Errorf("could not set L2 wallet nonce. Cause: %w", err))
		}
		l2Wallet.SetNonce(nonce)
	}

	// We set the ERC20 contract for the tokens.
	for _, token := range wallets.Tokens {
		token.L1ContractAddress = &nmConfig.erc20ContractAddress
	}

	return wallets
}

// Extracts the number of transactions to inject from the command-line arguments.
func parseNumOfTxs(args []string) int {
	if len(args) != 1 {
		panic(fmt.Errorf("expected one argument to %s command, got %d", injectTxsName, len(args)))
	}
	numOfTxs, err := strconv.Atoi(args[0])
	if err != nil {
		panic(fmt.Errorf("could not parse number of transactions to inject. Cause: %w", err))
	}
	return numOfTxs
}

func checkDepositsSuccessful(txInjector *simulation.TransactionInjector, l1Client ethadapter.EthClient, stats *stats.Stats, erc20ContractLib erc20contractlib.ERC20ContractLib, mgmtContractLib mgmtcontractlib.MgmtContractLib, startBlock *types.Block) {
	currentBlock := l1Client.FetchHeadBlock()
	dummySim := simulation.Simulation{
		Stats: stats,
		Params: &params.SimParams{
			ERC20ContractLib: erc20ContractLib,
			MgmtContractLib:  mgmtContractLib,
		},
	}
	deposits, _, _, _ := simulation.ExtractDataFromEthereumChain(startBlock, currentBlock, l1Client, &dummySim) // nolint:dogsled

	if len(deposits) != len(txInjector.TxTracker.L1Transactions) {
		println(fmt.Sprintf("Injected %d deposits into the L1 but %d were missing.",
			len(deposits), len(txInjector.TxTracker.L1Transactions)))
	} else {
		println(fmt.Sprintf("Successfully injected %d deposits into the L1.", len(deposits)))
	}
}

func checkL2TxsSuccessful(l2Client rpcclientlib.Client, txInjector *simulation.TransactionInjector) {
	injectedTransfers := len(txInjector.TxTracker.TransferL2Transactions)
	injectedWithdrawals := len(txInjector.TxTracker.WithdrawalL2Transactions)
	notFoundTransfers, notFoundWithdrawals := simulation.FindNotIncludedL2Txs(l2Client, txInjector)

	if notFoundTransfers != 0 {
		println(fmt.Sprintf("Injected %d transfers into the L2 but %d were missing.", injectedTransfers, notFoundTransfers))
	} else {
		println(fmt.Sprintf("Successfully injected %d transfers into the L1.", injectedTransfers))
	}

	if notFoundWithdrawals != 0 {
		println(fmt.Sprintf("Injected %d withdrawals into the L2 but %d were missing.", injectedWithdrawals, notFoundWithdrawals))
	} else {
		println(fmt.Sprintf("Successfully injected %d withdrawals into the L1.", injectedWithdrawals))
	}
}
