package transactioninjector

import (
	"fmt"
	"time"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/integration/simulation/stats"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/ethclient/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/config"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"
	"github.com/obscuronet/obscuro-playground/integration/simulation"
)

func InjectTransactions(cfg Config) {
	hostConfig := config.HostConfig{
		ID:                  common.HexToAddress(""),
		L1NodeHost:          cfg.l1NodeHost,
		L1NodeWebsocketPort: cfg.l1NodeWebsocketPort,
		L1ConnectionTimeout: cfg.l1ConnectionTimeout,
		PrivateKeyString:    cfg.privateKeyString,
		ChainID:             cfg.chainID,
	}

	// TODO - Consider extending this tool to support multiple L1 clients and L2 clients.
	l1Client, err := ethclient.NewEthClient(hostConfig)
	if err != nil {
		panic(fmt.Sprintf("could not create L1 client. Cause: %s", err))
	}
	l2Client := obscuroclient.NewClient(cfg.obscuroClientAddress)

	l1Wallet := wallet.NewInMemoryWalletFromConfig(hostConfig)
	nonce, err := l1Client.Nonce(l1Wallet.Address())
	if err != nil {
		panic(err)
	}
	l1Wallet.SetNonce(nonce)

	txInjector := simulation.NewTransactionInjector(
		1*time.Second,
		stats.NewStats(1),
		[]ethclient.EthClient{l1Client},
		&params.SimWallets{
			// todo
		},
		&cfg.mgmtContractAddress,
		[]obscuroclient.Client{l2Client},
		mgmtcontractlib.NewMgmtContractLib(&cfg.mgmtContractAddress),
		erc20contractlib.NewERC20ContractLib(&cfg.mgmtContractAddress, &cfg.erc20ContractAddress),
	)

	println("Injecting transactions into network...")
	txInjector.Start()
}