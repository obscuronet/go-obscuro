package common

import (
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/log"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

type (
	StateRoot             = gethcommon.Hash
	L1RootHash            = gethcommon.Hash
	L2RootHash            = gethcommon.Hash
	TxHash                = gethcommon.Hash
	L2Tx                  = types.Transaction
	EncryptedTx           []byte // A single transaction encrypted using the enclave's public key
	EncryptedTransactions []byte // A blob of encrypted transactions, as they're stored in the rollup.

	EncryptedParamsGetBalance     []byte // The params for an RPC getBalance request, as a JSON object encrypted with the public key of the enclave.
	EncryptedParamsCall           []byte // As above, but for an RPC call request.
	EncryptedParamsGetTxReceipt   []byte // As above, but for an RPC getTransactionReceipt request.
	EncryptedParamsGetTxCount     []byte // As above, but for an RPC getTransactionCount request.
	EncryptedResponseGetBalance   []byte // The response for an RPC getBalance request, as a JSON object encrypted with the viewing key of the user.
	EncryptedResponseCall         []byte // As above, but for an RPC call request.
	EncryptedResponseGetTxReceipt []byte // As above, but for an RPC getTransactionReceipt request.
	EncryptedResponseGetTxCount   []byte // As above, but for an RPC getTransactionCount request.

	Nonce         = uint64
	EncodedRollup []byte
)

const (
	L2GenesisHeight = uint64(0)
	L1GenesisHeight = uint64(0)
	// HeightCommittedBlocks is the number of blocks deep a transaction must be to be considered safe from reorganisations.
	HeightCommittedBlocks = 20
)

// Header is a public / plaintext struct that holds common properties of the Rollup
// Making changes to this struct will require GRPC + GRPC Converters regen
type Header struct {
	ParentHash  L2RootHash
	Agg         gethcommon.Address
	Nonce       Nonce
	L1Proof     L1RootHash // the L1 block where the Parent was published
	Root        StateRoot
	TxHash      gethcommon.Hash // todo - include the synthetic deposits
	Number      *big.Int        // the rollup height
	Withdrawals []Withdrawal
	Bloom       types.Bloom
	ReceiptHash gethcommon.Hash
	Extra       []byte
}

// Withdrawal - this is the withdrawal instruction that is included in the rollup header
type Withdrawal struct {
	// Type      uint8 // the type of withdrawal. For now only ERC20. Todo - add this once more ERCs are supported
	Amount    uint64
	Recipient gethcommon.Address // the user account that will receive the money
	Contract  gethcommon.Address // the contract
}

// ExtRollup is used for communication between the enclave and the outside world.
type ExtRollup struct {
	Header          *Header
	EncryptedTxBlob EncryptedTransactions
}

// EncryptedRollup extends ExtRollup with additional fields.
// This parallels the Block/extblock split in Go Ethereum.
type EncryptedRollup struct {
	Header *Header

	hash atomic.Value

	Transactions EncryptedTransactions
}

// AttestationReport represents a signed attestation report from a TEE and some metadata about the source of it to verify it
type AttestationReport struct {
	Report      []byte             // the signed bytes of the report which includes some encrypted identifying data
	PubKey      []byte             // a public key that can be used to send encrypted data back to the TEE securely (should only be used once Report has been verified)
	Owner       gethcommon.Address // address identifying the owner of the TEE which signed this report, can also be verified from the encrypted Report data
	HostAddress string             // the IP address on which the host can be contacted by other Obscuro hosts for peer-to-peer communication
}

func (er ExtRollup) ToRollup() *EncryptedRollup {
	return &EncryptedRollup{
		Header:       er.Header,
		Transactions: er.EncryptedTxBlob,
	}
}

type (
	EncryptedSharedEnclaveSecret []byte
	EncodedAttestationReport     []byte
)

var hasherPool = sync.Pool{
	New: func() interface{} { return sha3.NewLegacyKeccak256() },
}

// RLPHash encodes value, hashes the encoded bytes and returns the hash.
func RLPHash(value interface{}) (gethcommon.Hash, error) {
	var hash gethcommon.Hash

	sha := hasherPool.Get().(crypto.KeccakState)
	defer hasherPool.Put(sha)
	sha.Reset()

	err := rlp.Encode(sha, value)
	if err != nil {
		return hash, fmt.Errorf("unable to encode Value. %w", err)
	}

	_, err = sha.Read(hash[:])
	if err != nil {
		return hash, fmt.Errorf("unable to read encoded value. %w", err)
	}

	return hash, nil
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (r *EncryptedRollup) Hash() L2RootHash {
	if hash := r.hash.Load(); hash != nil {
		return hash.(L2RootHash)
	}
	v := r.Header.Hash()
	r.hash.Store(v)

	return v
}

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (h *Header) Hash() L2RootHash {
	hash, err := RLPHash(h)
	if err != nil {
		log.Error("err hashing the l2roothash")
	}
	return hash
}
