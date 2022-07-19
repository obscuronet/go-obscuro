package simulation

import (
	"sync"

	"github.com/obscuronet/go-obscuro/go/ethadapter"

	"github.com/obscuronet/go-obscuro/go/ethadapter/erc20contractlib"

	"github.com/obscuronet/go-obscuro/go/common"
)

type txInjectorCounter struct {
	l1TransactionsLock       sync.RWMutex
	L1Transactions           []ethadapter.L1Transaction
	l2TransactionsLock       sync.RWMutex
	TransferL2Transactions   []*common.L2Tx
	WithdrawalL2Transactions []*common.L2Tx
}

func newCounter() *txInjectorCounter {
	return &txInjectorCounter{
		l1TransactionsLock:       sync.RWMutex{},
		L1Transactions:           []ethadapter.L1Transaction{},
		l2TransactionsLock:       sync.RWMutex{},
		TransferL2Transactions:   []*common.L2Tx{},
		WithdrawalL2Transactions: []*common.L2Tx{},
	}
}

// trackL1Tx adds an L1Tx to the internal list
func (m *txInjectorCounter) trackL1Tx(tx ethadapter.L1Transaction) {
	m.l1TransactionsLock.Lock()
	defer m.l1TransactionsLock.Unlock()
	m.L1Transactions = append(m.L1Transactions, tx)
}

func (m *txInjectorCounter) trackWithdrawalL2Tx(tx *common.L2Tx) {
	m.l2TransactionsLock.Lock()
	defer m.l2TransactionsLock.Unlock()
	m.WithdrawalL2Transactions = append(m.WithdrawalL2Transactions, tx)
}

func (m *txInjectorCounter) trackTransferL2Tx(tx *common.L2Tx) {
	m.l2TransactionsLock.Lock()
	defer m.l2TransactionsLock.Unlock()
	m.TransferL2Transactions = append(m.TransferL2Transactions, tx)
}

// GetL1Transactions returns all generated L1 L2Txs
func (m *txInjectorCounter) GetL1Transactions() []ethadapter.L1Transaction {
	return m.L1Transactions
}

// GetL2Transactions returns all generated non-WithdrawalTx transactions
func (m *txInjectorCounter) GetL2Transactions() ([]*common.L2Tx, []*common.L2Tx) {
	return m.TransferL2Transactions, m.WithdrawalL2Transactions
}

// GetL2WithdrawalRequests returns generated stored WithdrawalTx transactions
func (m *txInjectorCounter) GetL2WithdrawalRequests() []common.Withdrawal {
	withdrawals := make([]common.Withdrawal, 0)
	for _, req := range m.WithdrawalL2Transactions {
		found, address, amount := erc20contractlib.DecodeTransferTx(req)
		if !found {
			panic("Should not happen")
		}
		withdrawals = append(withdrawals, common.Withdrawal{Amount: amount.Uint64(), Recipient: *address})
	}
	return withdrawals
}
