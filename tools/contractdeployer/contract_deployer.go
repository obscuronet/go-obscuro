package contractdeployer

// TODO we might merge this with the network manager package
import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/go-obscuro/go/ethadapter/erc20contractlib"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/wallet"
)

const (
	mgmtContract  = "MGMT"
	erc20Contract = "ERC20"
)

type ContractDeployer struct {
	client       ethadapter.EthClient
	wallet       wallet.Wallet
	contractCode []byte
}

func NewContractDeployer(config *Config) (*ContractDeployer, error) {
	wallet, err := setupWallet(config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup wallet - %w", err)
	}
	client, err := setupClient(config, wallet)
	if err != nil {
		return nil, fmt.Errorf("failed to setup client - %w", err)
	}
	contractCode, err := getContractCode(config)
	if err != nil {
		return nil, fmt.Errorf("failed to find contract bytecode to deploy - %w", err)
	}

	return &ContractDeployer{
		client:       client,
		wallet:       wallet,
		contractCode: contractCode,
	}, nil
}

func (cd *ContractDeployer) Run() error {
	nonce, err := cd.client.Nonce(cd.wallet.Address())
	if err != nil {
		return err
	}
	cd.wallet.SetNonce(nonce)

	deployContractTx := types.LegacyTx{
		Nonce:    cd.wallet.GetNonceAndIncrement(),
		GasPrice: big.NewInt(2000000000),
		Gas:      1025_000_000,
		Data:     cd.contractCode,
	}

	signedTx, err := cd.wallet.SignTransaction(&deployContractTx)
	if err != nil {
		return err
	}

	err = cd.client.SendTransaction(signedTx)
	if err != nil {
		return err
	}

	var start time.Time
	var receipt *types.Receipt
	var contractAddr *common.Address
	for start = time.Now(); time.Since(start) < 80*time.Second; time.Sleep(2 * time.Second) {
		receipt, err = cd.client.TransactionReceipt(signedTx.Hash())
		if err == nil && receipt != nil {
			if receipt.Status != types.ReceiptStatusSuccessful {
				return fmt.Errorf("unable to deploy contract, receipt status unsuccessful: %v", receipt)
			}
			log.Info("Contract successfully deployed to %s", receipt.ContractAddress)
			contractAddr = &receipt.ContractAddress
			break
		}

		log.Info("Contract deploy tx has not been mined into a block after %s...", time.Since(start))
	}
	if contractAddr == nil {
		return fmt.Errorf("failed to mine contract deploy tx into a block after %s. Aborting", time.Since(start))
	}

	// print the contract address, to be read if necessary by the caller (important: this must be the last message output by the script)
	fmt.Print(contractAddr.Hex())

	// this is a safety sleep to make sure the output is printed
	time.Sleep(5 * time.Second)
	return nil
}

func setupWallet(cfg *Config) (wallet.Wallet, error) {
	privateKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("could not recover private key from hex. Cause: %w", err)
	}

	// load the wallet
	return wallet.NewInMemoryWalletFromPK(cfg.ChainID, privateKey), nil
}

func setupClient(cfg *Config, _ wallet.Wallet) (ethadapter.EthClient, error) {
	// connect to the l1
	client, err := ethadapter.NewEthClient(cfg.NodeHost, cfg.NodePort, 30*time.Second, common.HexToAddress("0x0"))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func getContractCode(cfg *Config) ([]byte, error) {
	switch cfg.ContractName {
	case mgmtContract:
		return common.Hex2Bytes(mgmtcontractlib.MgmtContractByteCode), nil

	case erc20Contract:
		return common.Hex2Bytes(erc20contractlib.ERC20ContractABI), nil

	default:
		return nil, fmt.Errorf("unrecognised contract %s - no bytecode configured for that contract name", cfg.ContractName)
	}
}
