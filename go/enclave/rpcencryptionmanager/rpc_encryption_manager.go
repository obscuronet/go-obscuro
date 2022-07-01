package rpcencryptionmanager

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/common"

	"github.com/ethereum/go-ethereum/accounts"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
)

// ViewingKeySignedMsgPrefix is the prefix added when signing the viewing key in MetaMask using the personal_sign
// API. Why is this needed? MetaMask has a security feature whereby if you ask it to sign something that looks like
// a transaction using the personal_sign API, it modifies the data being signed. The goal is to prevent hackers
// from asking a visitor to their website to personal_sign something that is actually a malicious transaction (e.g.
// theft of funds). By adding a prefix, the viewing key bytes no longer looks like a transaction hash, and thus get
// signed as-is.
const ViewingKeySignedMsgPrefix = "vk"

// PlaceholderResult is used when the result to an eth_call is equal to nil. Attempting to encrypt then decrypt nil
// using ECIES throws an exception.
var PlaceholderResult = []byte("<nil result>")

// RPCEncryptionManager manages the decryption and encryption of sensitive RPC requests.
type RPCEncryptionManager struct {
	viewingKeysEnabled     bool
	enclavePrivateKeyECIES *ecies.PrivateKey
	// TODO - Replace with persistent storage.
	// TODO - Handle multiple viewing keys per address.
	viewingKeys map[gethcommon.Address]*ecies.PublicKey
}

func NewRPCEncryptionManager(viewingKeysEnabled bool, enclavePrivateKeyECIES *ecies.PrivateKey) RPCEncryptionManager {
	return RPCEncryptionManager{
		viewingKeysEnabled:     viewingKeysEnabled,
		enclavePrivateKeyECIES: enclavePrivateKeyECIES,
		viewingKeys:            make(map[gethcommon.Address]*ecies.PublicKey),
	}
}

// DecryptRPCCall decrypts the bytes with the enclave's private key if viewing keys are enabled.
func (rpc *RPCEncryptionManager) DecryptRPCCall(encryptedBytes []byte) ([]byte, error) {
	if !rpc.viewingKeysEnabled {
		return encryptedBytes, nil
	}
	return rpc.DecryptWithEnclavePrivateKey(encryptedBytes)
}

// DecryptWithEnclavePrivateKey the bytes with the enclave's private key.
func (rpc *RPCEncryptionManager) DecryptWithEnclavePrivateKey(encryptedBytes []byte) ([]byte, error) {
	bytes, err := rpc.enclavePrivateKeyECIES.Decrypt(encryptedBytes, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt bytes with enclave private key. Cause: %w", err)
	}

	return bytes, nil
}

// AddViewingKey - see the description of Enclave.AddViewingKey.
func (rpc *RPCEncryptionManager) AddViewingKey(viewingKeyBytes []byte, signature []byte) error {
	// We recalculate the message signed by MetaMask.
	msgToSign := ViewingKeySignedMsgPrefix + hex.EncodeToString(viewingKeyBytes)

	// We recover the key based on the signed message and the signature.
	recoveredPublicKey, err := crypto.SigToPub(accounts.TextHash([]byte(msgToSign)), signature)
	if err != nil {
		return fmt.Errorf("received viewing key but could not validate its signature. Cause: %w", err)
	}
	recoveredAddress := crypto.PubkeyToAddress(*recoveredPublicKey)

	// We decompress the viewing key and create the corresponding ECIES key.
	viewingKey, err := crypto.DecompressPubkey(viewingKeyBytes)
	if err != nil {
		return fmt.Errorf("received viewing key bytes but could not decompress them. Cause: %w", err)
	}
	eciesPublicKey := ecies.ImportECDSAPublic(viewingKey)

	rpc.viewingKeys[recoveredAddress] = eciesPublicKey

	return nil
}

// EncryptWithViewingKey encrypts the bytes with a viewing key for the address.
func (rpc *RPCEncryptionManager) EncryptWithViewingKey(address gethcommon.Address, bytes []byte) ([]byte, error) {
	if !rpc.viewingKeysEnabled {
		return bytes, nil
	}

	viewingKey := rpc.viewingKeys[address]
	if viewingKey == nil {
		return nil, fmt.Errorf("could not encrypt bytes because it does not have a viewing key for account %s", address.String())
	}

	if bytes == nil {
		bytes = PlaceholderResult
	}

	encryptedBytes, err := ecies.Encrypt(rand.Reader, viewingKey, bytes, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt bytes becauseit could not encrypt the response using a viewing key for account %s", address.String())
	}

	return encryptedBytes, nil
}

// ExtractTxHash - Returns the transaction hash from a common.EncryptedParamsGetTxReceipt object.
func (rpc *RPCEncryptionManager) ExtractTxHash(encryptedParams common.EncryptedParamsGetTxReceipt) (gethcommon.Hash, error) {
	paramBytes, err := rpc.DecryptRPCCall(encryptedParams)
	if err != nil {
		return gethcommon.Hash{}, fmt.Errorf("could not decrypt params in eth_getTransactionReceipt request. Cause: %w", err)
	}

	var paramsJSONList []string
	err = json.Unmarshal(paramBytes, &paramsJSONList)
	if err != nil {
		return gethcommon.Hash{}, fmt.Errorf("could not parse JSON params in eth_getTransactionReceipt request. Cause: %w", err)
	}
	txHash := gethcommon.HexToHash(paramsJSONList[0]) // The only argument is the transaction hash.
	return txHash, err
}

// ExtractAddress - Returns the address from a common.EncryptedParamsGetTransactionCount blob
func (rpc *RPCEncryptionManager) ExtractAddress(encryptedParams common.EncryptedParamsGetTxCount) (gethcommon.Address, error) {
	paramBytes, err := rpc.DecryptRPCCall(encryptedParams)
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("could not decrypt params in eth_getTransactionCount request. Cause: %w", err)
	}

	var paramsJSONList []string
	err = json.Unmarshal(paramBytes, &paramsJSONList)
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("could not parse JSON params in eth_getTransactionCount request. Cause: %w", err)
	}
	txHash := gethcommon.HexToAddress(paramsJSONList[0]) // The only argument is the transaction hash.
	return txHash, err
}

// Marshalls the transaction receipt to JSON, and encrypts it with a viewing key for the address.
func (rpc *RPCEncryptionManager) EncryptTxReceiptWithViewingKey(address gethcommon.Address, txReceipt *types.Receipt) ([]byte, error) {
	txReceiptBytes, err := txReceipt.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("could not marshall transaction receipt to JSON in eth_getTransactionReceipt request. Cause: %w", err)
	}
	return rpc.EncryptWithViewingKey(address, txReceiptBytes)
}

// DecryptTx decrypts an L2 transaction encrypted with the enclave's public key.
func (rpc *RPCEncryptionManager) DecryptTx(encryptedTx common.EncryptedTx) (*common.L2Tx, error) {
	txBytes, err := rpc.DecryptWithEnclavePrivateKey(encryptedTx)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt transaction with enclave private key. Cause: %w", err)
	}

	transaction := common.L2Tx{}
	if err = rlp.DecodeBytes(txBytes, &transaction); err != nil {
		return nil, fmt.Errorf("could not decrypt encrypted L2 transaction. Cause: %w", err)
	}

	return &transaction, nil
}
