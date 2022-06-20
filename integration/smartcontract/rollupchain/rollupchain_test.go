package rollupchain

import (
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/contracts/compiledcontracts/generatedRollupChainTestContract"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/config"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"
	"github.com/obscuronet/obscuro-playground/integration"
	"github.com/obscuronet/obscuro-playground/integration/datagenerator"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"
	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
	"github.com/obscuronet/obscuro-playground/integration/smartcontract/debugwallet"
)

// netInfo is a bag holder struct for output data from the execution/run of a network
type netInfo struct {
	ethClients  []ethclient.EthClient
	wallets     []wallet.Wallet
	gethNetwork *gethnetwork.GethNetwork
}

// runGethNetwork runs a geth network with one prefunded wallet
func runGethNetwork(t *testing.T) *netInfo {
	// make sure the geth network binaries exist
	path, err := gethnetwork.EnsureBinariesExist(gethnetwork.LatestVersion)
	if err != nil {
		t.Fatal(err)
	}

	// prefund one wallet as the worker wallet
	workerWallet := datagenerator.RandomWallet(integration.EthereumChainID)

	// define + run the network
	gethNetwork := gethnetwork.NewGethNetwork(
		integration.StartPortRollupChainContractTests,
		integration.StartPortRollupChainContractTests+100,
		path,
		3,
		1,
		[]string{workerWallet.Address().String()},
	)

	// create a client that is connected to node 0 of the network
	client, err := ethclient.NewEthClient(config.HostConfig{
		ID:                  common.Address{1},
		L1NodeHost:          "127.0.0.1",
		L1NodeWebsocketPort: gethNetwork.WebSocketPorts[0],
		L1ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		return nil
	}

	return &netInfo{
		ethClients:  []ethclient.EthClient{client},
		wallets:     []wallet.Wallet{workerWallet},
		gethNetwork: gethNetwork,
	}
}

func TestRollupChainContract(t *testing.T) {
	// run tests on one network
	sim := runGethNetwork(t)
	defer sim.gethNetwork.StopNodes()

	// setup the client and the (debug) wallet
	client := sim.ethClients[0]
	w := debugwallet.NewDebugWallet(sim.wallets[0])

	for name, test := range map[string]func(*testing.T, *debugRollupChainTestContractLib, *debugwallet.DebugWallet, ethclient.EthClient){
		"executeSmartContractTests": executeSmartContractTests,
	} {
		t.Run(name, func(t *testing.T) {

			// deploy the library first - use the address of the library in the smart contract test compilation
			libAddress, err := network.DeployContract(client, w, common.Hex2Bytes(generatedRollupChainTestContract.RollupChainMetaData.Bin[2:]))
			generatedRollupChainTestContract.RollupChainTestContractMetaData.Bin = strings.Replace(generatedRollupChainTestContract.RollupChainTestContractMetaData.Bin, "__$7204b1ba8a254ced74f31676d70e6726eb$__", libAddress.String()[2:], -1)

			// deploy the smart contract test to a new address
			contractAddr, err := network.DeployContract(client, w, common.Hex2Bytes(generatedRollupChainTestContract.RollupChainTestContractMetaData.Bin[2:]))
			if err != nil {
				t.Error(err)
			}

			// run the test using the new contract, but same wallet
			test(t,
				newdebugRollupChainContractLib(*contractAddr, client.EthClient()),
				w,
				client,
			)
		})
	}
}

// executeSmartContractTests runs the smartcontract tests
func executeSmartContractTests(t *testing.T, contractLib *debugRollupChainTestContractLib, w *debugwallet.DebugWallet, client ethclient.EthClient) {
	// run the tests in the smartcontract ( that are expected to succeed)
	for name, scTest := range map[string]func(opts *bind.TransactOpts) (*types.Transaction, error){
		"AppendRollupTest": contractLib.genContract.AppendRollupTest,
		"ScrollTreeTest":   contractLib.genContract.ScrollTreeTest,
		"NoForkDetection":  contractLib.genContract.NoForkDetection,
	} {
		issuedTx, err := scTest(
			&bind.TransactOpts{
				From:  w.Address(),
				Nonce: big.NewInt(int64(w.GetNonceAndIncrement())),
				Signer: func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
					return types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1337)), w.PrivateKey())
				},
				GasPrice: big.NewInt(225000000),
				GasLimit: 1_000_000_000,
			},
		)
		if err != nil {
			t.Errorf("%s: %s", name, err)
		}

		receipt, err := debugwallet.WaitTxResult(client, issuedTx)
		if err != nil {
			t.Errorf("%s: %s", name, err)
		}

		if receipt.Status != types.ReceiptStatusSuccessful {
			t.Errorf("%s: transaction should have succeeded, expected %d got %d", name, types.ReceiptStatusSuccessful, receipt.Status)
			msg, err := w.DebugTransaction(client, issuedTx)
			if err != nil {
				t.Errorf("debug transaction failure reason: %s", err)
			}
			t.Log(string(msg))
		}
	}

	// runs smartcontract tests that are expected to revert
	for name, scTest := range map[string]func(opts *bind.TransactOpts) (*types.Transaction, error){
		"RevertsNoDoubleInitTest": contractLib.genContract.RevertsNoDoubleInitTest,
	} {
		issuedTx, err := scTest(
			&bind.TransactOpts{
				From:  w.Address(),
				Nonce: big.NewInt(int64(w.GetNonceAndIncrement())),
				Signer: func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
					return types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1337)), w.PrivateKey())
				},
				GasPrice: big.NewInt(225000000),
				GasLimit: 1_000_000_000,
			},
		)
		if err != nil {
			t.Errorf("%s: %s", name, err)
		}

		receipt, err := debugwallet.WaitTxResult(client, issuedTx)
		if err != nil {
			t.Errorf("%s: %s", name, err)
		}

		if receipt.Status != types.ReceiptStatusFailed {
			t.Errorf("%s: transaction should have failed, expected %d got %d", name, types.ReceiptStatusFailed, receipt.Status)
		}
	}
}
