package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/obscuronet/go-obscuro/go/common/log"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
)

// obscuroPrivateKeyHex is the private key used for sensitive communication with the enclave.
// TODO - Replace this fixed key with a key derived from the master seed.
const obscuroPrivateKeyHex = "81acce9620f0adf1728cb8df7f6b8b8df857955eb9e8b7aed6ef8390c09fc207"

func GetObscuroKey() *ecdsa.PrivateKey {
	key, err := crypto.HexToECDSA(obscuroPrivateKeyHex)
	if err != nil {
		log.Panic("failed to create enclave private key. Cause: %s", err)
	}
	return key
}

func GenerateEntropy() core.SharedEnclaveSecret {
	secret := make([]byte, core.SharedSecretLen)
	// todo - check if there is a better way to do this in ego.
	n, err := rand.Read(secret)
	if n != core.SharedSecretLen || err != nil {
		log.Panic("could not generate secret. Cause: %s", err)
	}
	var temp [core.SharedSecretLen]byte
	copy(temp[:], secret)
	return temp
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *ecdsa.PublicKey) ([]byte, error) {
	ciphertext, err := ecies.Encrypt(rand.Reader, ecies.ImportECDSAPublic(pub), msg, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt with public key. %w", err)
	}
	return ciphertext, nil
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *ecdsa.PrivateKey) ([]byte, error) {
	plaintext, err := ecies.ImportECDSA(priv).Decrypt(ciphertext, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt with private key. %w", err)
	}
	return plaintext, nil
}

func DecryptSecret(secret common.EncryptedSharedEnclaveSecret, privateKey *ecdsa.PrivateKey) (*core.SharedEnclaveSecret, error) {
	if privateKey == nil {
		return nil, errors.New("private key not found - shouldn't happen")
	}
	value, err := DecryptWithPrivateKey(secret, privateKey)
	if err != nil {
		return nil, err
	}
	var temp core.SharedEnclaveSecret
	copy(temp[:], value)
	return &temp, nil
}

func EncryptSecret(pubKeyEncoded []byte, secret core.SharedEnclaveSecret, nodeShortID uint64) (common.EncryptedSharedEnclaveSecret, error) {
	common.LogWithID(nodeShortID, "Encrypting secret with public key %s", gethcommon.Bytes2Hex(pubKeyEncoded))
	key, err := crypto.DecompressPubkey(pubKeyEncoded)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key %w", err)
	}

	encKey, err := EncryptWithPublicKey(secret[:], key)
	if err != nil {
		common.LogWithID(nodeShortID, "Failed to encrypt key, err: %s\nsecret: %v\npubkey: %v\nencKey:%v", err, secret, pubKeyEncoded, encKey)
	}
	return encKey, err
}
