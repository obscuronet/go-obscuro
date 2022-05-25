package simulation

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/contracts"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

type txInjectorCounter struct {
	l1TransactionsLock       sync.RWMutex
	l1Transactions           []obscurocommon.L1Transaction
	l2TransactionsLock       sync.RWMutex
	transferL2Transactions   core.L2Txs
	withdrawalL2Transactions core.L2Txs
}

func newCounter() *txInjectorCounter {
	return &txInjectorCounter{
		l1TransactionsLock:       sync.RWMutex{},
		l1Transactions:           []obscurocommon.L1Transaction{},
		l2TransactionsLock:       sync.RWMutex{},
		transferL2Transactions:   []nodecommon.L2Tx{},
		withdrawalL2Transactions: []nodecommon.L2Tx{},
	}
}

// trackL1Tx adds an L1Tx to the internal list
func (m *txInjectorCounter) trackL1Tx(tx obscurocommon.L1Transaction) {
	m.l1TransactionsLock.Lock()
	defer m.l1TransactionsLock.Unlock()
	m.l1Transactions = append(m.l1Transactions, tx)
}

func (m *txInjectorCounter) trackWithdrawalL2Tx(tx nodecommon.L2Tx) {
	m.l2TransactionsLock.Lock()
	defer m.l2TransactionsLock.Unlock()
	m.withdrawalL2Transactions = append(m.withdrawalL2Transactions, tx)
}

func (m *txInjectorCounter) trackTransferL2Tx(tx nodecommon.L2Tx) {
	m.l2TransactionsLock.Lock()
	defer m.l2TransactionsLock.Unlock()
	m.transferL2Transactions = append(m.transferL2Transactions, tx)
}

// GetL1Transactions returns all generated L1 L2Txs
func (m *txInjectorCounter) GetL1Transactions() []obscurocommon.L1Transaction {
	return m.l1Transactions
}

// GetL2Transactions returns all generated non-WithdrawalTx transactions
func (m *txInjectorCounter) GetL2Transactions() (core.L2Txs, core.L2Txs) {
	return m.transferL2Transactions, m.withdrawalL2Transactions
}

// GetL2WithdrawalRequests returns generated stored WithdrawalTx transactions
func (m *txInjectorCounter) GetL2WithdrawalRequests() []nodecommon.Withdrawal {
	withdrawals := make([]nodecommon.Withdrawal, 0)
	for _, req := range m.withdrawalL2Transactions {
		// todo - helper
		method, err := contracts.PedroERC20ContractABIJSON.MethodById(req.Data()[:4])
		if err != nil || method.Name != "transfer" {
			panic(err)
		}
		args := map[string]interface{}{}
		if err := method.Inputs.UnpackIntoMap(args, req.Data()[4:]); err != nil {
			panic(err)
		}
		withdrawals = append(withdrawals, nodecommon.Withdrawal{Amount: args["amount"].(*big.Int).Uint64(), Address: args["to"].(common.Address)})
	}
	return withdrawals
}