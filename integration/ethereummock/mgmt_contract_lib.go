package ethereummock

import (
	"bytes"
	"encoding/gob"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

var (
	depositTxAddr       = common.HexToAddress("0x01")
	rollupTxAddr        = common.HexToAddress("0x02")
	storeSecretTxAddr   = common.HexToAddress("0x03")
	requestSecretTxAddr = common.HexToAddress("0x04")
)

// mockContractLib is an implementation of the mgmtcontractlib.MgmtContractLib
// it creates ethereum mocked transactions from obscurocommon.L1Transaction
// and converts ethereum mocked transactions to obscurocommon.L1Transaction
type mockContractLib struct{}

func NewMgmtContractLibMock() mgmtcontractlib.MgmtContractLib {
	return &mockContractLib{}
}

func (m *mockContractLib) DecodeTx(tx *types.Transaction) obscurocommon.L1Transaction {
	return decodeTx(tx)
}

func (m *mockContractLib) CreateRollup(tx *obscurocommon.L1RollupTx, nonce uint64) types.TxData {
	return encodeTx(tx, nonce, rollupTxAddr)
}

func (m *mockContractLib) CreateRequestSecret(tx *obscurocommon.L1RequestSecretTx, nonce uint64) types.TxData {
	return encodeTx(tx, nonce, requestSecretTxAddr)
}

func (m *mockContractLib) CreateStoreSecret(tx *obscurocommon.L1StoreSecretTx, nonce uint64) types.TxData {
	return encodeTx(tx, nonce, storeSecretTxAddr)
}

func decodeTx(tx *types.Transaction) obscurocommon.L1Transaction {
	if len(tx.Data()) == 0 {
		panic("Data cannot be 0 in the mock implementation")
	}

	// prepare byte buffer
	buf := bytes.NewBuffer(tx.Data())
	dec := gob.NewDecoder(buf)

	// in the mock implementation we use the To address field to specify what is the L1 operation
	// the mock implementation does not process contracts
	// so this is a way that we can differentiate different contract calls
	var t obscurocommon.L1Transaction
	switch tx.To().Hex() {
	case rollupTxAddr.Hex():
		t = &obscurocommon.L1RollupTx{}
	case storeSecretTxAddr.Hex():
		t = &obscurocommon.L1StoreSecretTx{}
	case depositTxAddr.Hex():
		t = &obscurocommon.L1DepositTx{}
	case requestSecretTxAddr.Hex():
		t = &obscurocommon.L1RequestSecretTx{}
	default:
		panic("unexpected type")
	}

	// decode to interface implementation
	if err := dec.Decode(t); err != nil {
		panic(err)
	}
	return t
}

func encodeTx(tx obscurocommon.L1Transaction, nonce uint64, opType common.Address) types.TxData {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(tx); err != nil {
		panic(err)
	}

	// the mock implementation does not process contract calls
	// this uses the To address to distinguish between different contract calls / different l1 transactions
	return &types.LegacyTx{
		Nonce: nonce,
		Data:  buf.Bytes(),
		To:    &opType,
	}
}