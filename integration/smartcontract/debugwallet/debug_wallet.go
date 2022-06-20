package debugwallet

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"
)

var _timeout = 30 * time.Second

// DebugWallet is a wrapper around the wallet that simplifies commonly used functions
type DebugWallet struct {
	wallet.Wallet
}

// NewDebugWallet returns a new debug wrapped wallet
func NewDebugWallet(w wallet.Wallet) *DebugWallet {
	return &DebugWallet{w}
}

// AwaitedSignAndSendTransaction signs a tx, issues the tx and awaits the tx to be minted into a block
func (w *DebugWallet) AwaitedSignAndSendTransaction(client ethclient.EthClient, txData types.TxData) (*types.Transaction, *types.Receipt, error) {
	signedTx, err := w.SignAndSendTransaction(client, txData)
	if err != nil {
		return nil, nil, err
	}
	receipt, err := WaitTxResult(client, signedTx)
	if err != nil {
		return nil, nil, err
	}
	return signedTx, receipt, nil
}

// SignAndSendTransaction signs and sends a tx
func (w *DebugWallet) SignAndSendTransaction(client ethclient.EthClient, txData types.TxData) (*types.Transaction, error) {
	signedTx, err := w.SignTransaction(txData)
	if err != nil {
		return nil, err
	}

	err = client.SendTransaction(signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

func (w *DebugWallet) DebugTransaction(client ethclient.EthClient, tx *types.Transaction) ([]byte, error) {
	return client.EthClient().CallContract(
		context.Background(),
		ethereum.CallMsg{
			From:       w.Address(),
			To:         tx.To(),
			Gas:        tx.Gas(),
			GasPrice:   tx.GasPrice(),
			GasFeeCap:  tx.GasFeeCap(),
			GasTipCap:  tx.GasTipCap(),
			Value:      tx.Value(),
			Data:       tx.Data(),
			AccessList: tx.AccessList(),
		},
		nil,
	)
}

// WaitTxResult waits for a tx to be minted into a block
func WaitTxResult(client ethclient.EthClient, tx *types.Transaction) (*types.Receipt, error) {
	var receipt *types.Receipt
	var err error
	for start := time.Now(); time.Since(start) < _timeout; time.Sleep(time.Second) {
		receipt, err = client.TransactionReceipt(tx.Hash())

		if err != nil {
			if errors.Is(err, ethereum.NotFound) {
				continue
			}
			return nil, err
		}

		return receipt, nil
	}
	return nil, fmt.Errorf("transaction not minted after timeout")
}
