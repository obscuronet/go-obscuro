package networkmanager

import (
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	obscuroconfig "github.com/obscuronet/obscuro-playground/go/obscuronode/config"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"
	"github.com/obscuronet/obscuro-playground/integration/erc20contract"
	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
)

var (
	mgmtContractBytes  = common.Hex2Bytes(mgmtcontractlib.MgmtContractByteCode)
	erc20ContractBytes = common.Hex2Bytes(erc20contract.ContractByteCode)
)

// DeployContract deploys a management contract or ERC20 contract to the L1 network, and prints its address.
func DeployContract(config Config) {
	var contractBytes []byte
	switch config.Command { //nolint:exhaustive
	case DeployMgmtContract:
		contractBytes = mgmtContractBytes
	case DeployERC20Contract:
		contractBytes = erc20ContractBytes
	default:
		panic("unrecognised command type")
	}

	hostConfig := obscuroconfig.HostConfig{
		L1NodeHost:          config.l1NodeHost,
		L1NodeWebsocketPort: config.l1NodeWebsocketPort,
		L1ConnectionTimeout: config.l1ConnectionTimeout,
		PrivateKeyString:    config.privateKeyString,
		ChainID:             config.chainID,
	}

	l1Client, err := ethclient.NewEthClient(hostConfig)
	if err != nil {
		panic(err)
	}

	l1Wallet := wallet.NewInMemoryWalletFromString(hostConfig)
	nonce, err := l1Client.Nonce(l1Wallet.Address())
	if err != nil {
		panic(err)
	}
	l1Wallet.SetNonce(nonce)

	var contractAddress *common.Address
	_, contractAddress, err = network.DeployContract(l1Client, l1Wallet, contractBytes)
	if err != nil {
		panic(err)
	}

	println(contractAddress.Hex())
	os.Exit(0)
}
