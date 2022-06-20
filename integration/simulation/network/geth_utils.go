package network

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/config"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"
	"github.com/obscuronet/obscuro-playground/integration/erc20contract"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"
)

func SetUpGethNetwork(wallets *params.SimWallets, StartPort int, nrNodes int, blockDurationSeconds int) (*common.Address, *common.Address, []ethclient.EthClient, *gethnetwork.GethNetwork) {
	// make sure the geth network binaries exist
	path, err := gethnetwork.EnsureBinariesExist(gethnetwork.LatestVersion)
	if err != nil {
		panic(err)
	}

	// get the node wallet addresses to prefund them with Eth, so they can submit rollups, deploy contracts, deposit to the bridge, etc
	walletAddresses := make([]string, len(wallets.AllEthWallets()))
	for i, w := range wallets.AllEthWallets() {
		walletAddresses[i] = w.Address().String()
	}

	// kickoff the network with the prefunded wallet addresses
	gethNetwork := gethnetwork.NewGethNetwork(
		StartPort,
		StartPort+DefaultWsPortOffset,
		path,
		nrNodes,
		blockDurationSeconds,
		walletAddresses,
	)

	// connect to the first host to deploy
	tmpHostConfig := config.HostConfig{
		L1NodeHost:          Localhost,
		L1NodeWebsocketPort: gethNetwork.WebSocketPorts[0],
		L1ConnectionTimeout: DefaultL1ConnectionTimeout,
	}
	tmpEthClient, err := ethclient.NewEthClient(tmpHostConfig)
	if err != nil {
		panic(err)
	}

	mgmtContractAddr, err := DeployContract(tmpEthClient, wallets.MCOwnerWallet, common.Hex2Bytes(mgmtcontractlib.MgmtContractByteCode))
	if err != nil {
		panic(fmt.Sprintf("failed to deploy management contract. Cause: %s", err))
	}
	// todo deploy multiple erc20s here and store the mappings, etc
	erc20ContractAddr, err := DeployContract(tmpEthClient, wallets.Erc20EthOwnerWallets[0], common.Hex2Bytes(erc20contract.ContractByteCode))
	if err != nil {
		panic(fmt.Sprintf("failed to deploy ERC20 contract. Cause: %s", err))
	}

	ethClients := make([]ethclient.EthClient, nrNodes)
	for i := 0; i < nrNodes; i++ {
		ethClients[i] = createEthClientConnection(int64(i), gethNetwork.WebSocketPorts[i])
	}

	return mgmtContractAddr, erc20ContractAddr, ethClients, gethNetwork
}

func StopGethNetwork(clients []ethclient.EthClient, netw *gethnetwork.GethNetwork) {
	// Stop the clients first
	for _, c := range clients {
		if c != nil {
			c.Stop()
		}
	}
	// Stop the nodes second
	netw.StopNodes()
}

// DeployContract todo -this should live somewhere else
func DeployContract(workerClient ethclient.EthClient, w wallet.Wallet, contractBytes []byte) (*common.Address, error) {
	if len(contractBytes) == 0 {
		return nil, fmt.Errorf("unable to deploy a 0 byte contract")
	}

	deployContractTx := types.LegacyTx{
		Nonce:    w.GetNonceAndIncrement(),
		GasPrice: big.NewInt(2000000000),
		Gas:      1025_000_000,
		Data:     contractBytes,
	}

	signedTx, err := w.SignTransaction(&deployContractTx)
	if err != nil {
		return nil, err
	}

	err = workerClient.SendTransaction(signedTx)
	if err != nil {
		return nil, err
	}

	var receipt *types.Receipt
	for start := time.Now(); time.Since(start) < 80*time.Second; time.Sleep(2 * time.Second) {
		receipt, err = workerClient.TransactionReceipt(signedTx.Hash())
		if err == nil && receipt != nil {
			if receipt.Status != types.ReceiptStatusSuccessful {
				return nil, errors.New("unable to deploy contract")
			}
			break
		}

		log.Info("Contract deploy tx has not been mined into a block after %s...", time.Since(start))
	}

	log.Info("Contract successfully deployed to %s", receipt.ContractAddress)
	return &receipt.ContractAddress, nil
}

func createEthClientConnection(id int64, port uint) ethclient.EthClient {
	hostConfig := config.HostConfig{
		ID:                  common.BigToAddress(big.NewInt(id)),
		L1NodeHost:          Localhost,
		L1NodeWebsocketPort: port,
		L1ConnectionTimeout: DefaultL1ConnectionTimeout,
	}
	ethnode, err := ethclient.NewEthClient(hostConfig)
	if err != nil {
		panic(err)
	}
	return ethnode
}
