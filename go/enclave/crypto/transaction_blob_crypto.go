package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/go/enclave/core"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
)

const (
	// TODO - This fixed nonce is insecure, and should be removed alongside the fixed rollup encryption key.
	RollupCipherNonce = "000000000000"
	// RollupEncryptionKeyHex is the AES key used to encrypt and decrypt the transaction blob in rollups.
	// TODO - Replace this fixed key with derived, rotating keys.
	RollupEncryptionKeyHex = "bddbc0d46a0666ce57a466168d99c1830b0c65e052d77188f2cbfc3f6486588c"
)

// TransactionBlobCrypto handles the encryption and decryption of the transaction blobs stored inside a rollup.
type TransactionBlobCrypto interface {
	// ToExtRollup - Transforms an internal rollup as seen by the enclave to an external rollup with an encrypted payload
	ToExtRollup(r *core.Rollup) common.ExtRollup
	ToEnclaveRollup(r *common.EncryptedRollup) *core.Rollup
}

type TransactionBlobCryptoImpl struct {
	transactionCipher cipher.AEAD
}

func NewTransactionBlobCryptoImpl() TransactionBlobCrypto {
	key := gethcommon.Hex2Bytes(RollupEncryptionKeyHex)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(fmt.Sprintf("could not initialise AES cipher for enclave rollup key. Cause: %s", err))
	}
	transactionCipher, err := cipher.NewGCM(block)
	if err != nil {
		panic(fmt.Sprintf("could not initialise wrapper for AES cipher for enclave rollup key. Cause: %s", err))
	}
	return TransactionBlobCryptoImpl{
		transactionCipher: transactionCipher,
	}
}

func (t TransactionBlobCryptoImpl) ToExtRollup(r *core.Rollup) common.ExtRollup {
	txHashes := make([]gethcommon.Hash, len(r.Transactions))
	for idx, tx := range r.Transactions {
		txHashes[idx] = tx.Hash()
	}

	return common.ExtRollup{
		Header:          r.Header,
		TxHashes:        txHashes,
		EncryptedTxBlob: t.encrypt(r.Transactions),
	}
}

func (t TransactionBlobCryptoImpl) ToEnclaveRollup(r *common.EncryptedRollup) *core.Rollup {
	return &core.Rollup{
		Header:       r.Header,
		Transactions: t.decrypt(r.Transactions),
	}
}

// TODO - Modify this logic so that transactions with different reveal periods are in different blobs, as per the whitepaper.
func (t TransactionBlobCryptoImpl) encrypt(transactions []*common.L2Tx) common.EncryptedTransactions {
	encodedTxs, err := rlp.EncodeToBytes(transactions)
	if err != nil {
		log.Panic("could not encrypt L2 transaction. Cause: %s", err)
	}

	return t.transactionCipher.Seal(nil, []byte(RollupCipherNonce), encodedTxs, nil)
}

func (t TransactionBlobCryptoImpl) decrypt(encryptedTxs common.EncryptedTransactions) []*common.L2Tx {
	encodedTxs, err := t.transactionCipher.Open(nil, []byte(RollupCipherNonce), encryptedTxs, nil)
	if err != nil {
		log.Panic("could not decrypt encrypted L2 transactions. Cause: %s", err)
	}

	var txs []*common.L2Tx
	if err := rlp.DecodeBytes(encodedTxs, &txs); err != nil {
		log.Panic("could not decode encoded L2 transactions. Cause: %s", err)
	}

	return txs
}
