package host

import (
	"context"
	"errors"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common"
)

// EthereumAPI implements a subset of the Ethereum JSON RPC operations. All the method signatures are copied from the
// corresponding Geth implementations.
type EthereumAPI struct {
	host *Node
}

func NewEthereumAPI(host *Node) *EthereumAPI {
	return &EthereumAPI{
		host: host,
	}
}

// ChainId returns the Obscuro chain ID.
func (api *EthereumAPI) ChainId() (*hexutil.Big, error) { //nolint:stylecheck,revive
	return (*hexutil.Big)(big.NewInt(api.host.config.ObscuroChainID)), nil
}

// BlockNumber returns the height of the current head rollup.
func (api *EthereumAPI) BlockNumber() hexutil.Uint64 {
	return hexutil.Uint64(api.host.nodeDB.GetCurrentRollupHead().Number.Uint64())
}

// GetBalance returns the address's balance on the Obscuro network, encrypted with the viewing key corresponding to the
// `address` field and encoded as hex.
func (api *EthereumAPI) GetBalance(_ context.Context, encryptedParams common.EncryptedParamsGetBalance) (string, error) {
	encryptedBalance, err := api.host.EnclaveClient.GetBalance(encryptedParams)
	if err != nil {
		return "", err
	}
	return gethcommon.Bytes2Hex(encryptedBalance), nil
}

// GetBlockByNumber returns the rollup with the given height as a block. No transactions are included.
func (api *EthereumAPI) GetBlockByNumber(_ context.Context, number rpc.BlockNumber, _ bool) (map[string]interface{}, error) {
	extRollup, err := api.host.EnclaveClient.GetRollupByHeight(number.Int64())
	return extRollupToBlock(extRollup), err
}

// GetBlockByHash returns the rollup with the given hash as a block. No transactions are included.
func (api *EthereumAPI) GetBlockByHash(_ context.Context, hash gethcommon.Hash, _ bool) (map[string]interface{}, error) {
	extRollup, err := api.host.EnclaveClient.GetRollup(hash)
	return extRollupToBlock(extRollup), err
}

// GasPrice is a placeholder for an RPC method required by MetaMask/Remix.
func (api *EthereumAPI) GasPrice(context.Context) (*hexutil.Big, error) {
	return (*hexutil.Big)(big.NewInt(0)), nil
}

// Call returns the result of executing the smart contract as a user, encrypted with the viewing key corresponding to
// the `from` field and encoded as hex.
func (api *EthereumAPI) Call(_ context.Context, encryptedParams common.EncryptedParamsCall) (string, error) {
	encryptedResponse, err := api.host.EnclaveClient.ExecuteOffChainTransaction(encryptedParams)
	if err != nil {
		return "", err
	}
	return gethcommon.Bytes2Hex(encryptedResponse), nil
}

// GetTransactionReceipt returns the transaction receipt for the given transaction hash, encrypted with the viewing key
// corresponding to the original transaction submitter and encoded as hex.
func (api *EthereumAPI) GetTransactionReceipt(_ context.Context, encryptedParams common.EncryptedParamsGetTxReceipt) (*string, error) {
	encryptedResponse, err := api.host.EnclaveClient.GetTransactionReceipt(encryptedParams)
	if err != nil {
		return nil, err
	}
	if encryptedResponse == nil {
		return nil, err
	}
	encryptedResponseHex := gethcommon.Bytes2Hex(encryptedResponse)
	return &encryptedResponseHex, nil
}

// EstimateGas is a placeholder for an RPC method required by MetaMask/Remix.
func (api *EthereumAPI) EstimateGas(_ context.Context, _ interface{}, _ *rpc.BlockNumberOrHash) (hexutil.Uint64, error) {
	// TODO - Return a non-dummy gas estimate.
	return 0, nil
}

// SendRawTransaction sends the encrypted transaction.
func (api *EthereumAPI) SendRawTransaction(_ context.Context, encryptedParams common.EncryptedParamsSendRawTx) (string, error) {
	encryptedResponse, err := api.host.SubmitAndBroadcastTx(encryptedParams)
	if err != nil {
		return "", err
	}
	return gethcommon.Bytes2Hex(encryptedResponse), nil
}

// GetCode returns the code stored at the given address in the state for the given rollup height or rollup hash.
func (api *EthereumAPI) GetCode(_ context.Context, address gethcommon.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	rollupHeight, ok := blockNrOrHash.Number()
	if ok {
		rollup, err := api.host.EnclaveClient.GetRollupByHeight(rollupHeight.Int64())
		if err != nil {
			return nil, err
		}
		rollupHash := rollup.Header.Hash()
		return api.host.EnclaveClient.GetCode(address, &rollupHash)
	}

	rollupHash, ok := blockNrOrHash.Hash()
	if ok {
		return api.host.EnclaveClient.GetCode(address, &rollupHash)
	}

	return nil, errors.New("invalid arguments; neither rollup height nor rollup hash specified")
}

// TODO - Temporary. Will be replaced by encrypted implementation.
func (api *EthereumAPI) GetTransactionCount(_ context.Context, address gethcommon.Address, _ rpc.BlockNumberOrHash) (*hexutil.Uint64, error) {
	nonce := api.host.EnclaveClient.Nonce(address)
	return (*hexutil.Uint64)(&nonce), nil
}

// GetTransactionByHash returns the transaction with the given hash, encrypted with the viewing key corresponding to the
// `from` field and encoded as hex.
func (api *EthereumAPI) GetTransactionByHash(_ context.Context, encryptedParams common.EncryptedParamsGetTxByHash) (string, error) {
	encryptedResponse, err := api.host.EnclaveClient.GetTransaction(encryptedParams)
	if err != nil {
		return "", err
	}
	return gethcommon.Bytes2Hex(encryptedResponse), nil
}

// Maps an external rollup to a block.
func extRollupToBlock(extRollup *common.ExtRollup) map[string]interface{} {
	return map[string]interface{}{
		"number":           (*hexutil.Big)(extRollup.Header.Number),
		"hash":             extRollup.Header.Hash(),
		"parentHash":       extRollup.Header.ParentHash,
		"nonce":            extRollup.Header.Nonce,
		"logsBloom":        extRollup.Header.Bloom,
		"stateRoot":        extRollup.Header.Root,
		"receiptsRoot":     extRollup.Header.ReceiptHash,
		"miner":            extRollup.Header.Agg,
		"extraData":        hexutil.Bytes(extRollup.Header.Extra),
		"transactionsRoot": extRollup.Header.TxHash,
		"transactions":     extRollup.TxHashes,

		"sha3Uncles":    extRollup.Header.UncleHash,
		"difficulty":    extRollup.Header.Difficulty,
		"gasLimit":      extRollup.Header.GasLimit,
		"gasUsed":       extRollup.Header.GasUsed,
		"timestamp":     extRollup.Header.Time,
		"mixHash":       extRollup.Header.MixDigest,
		"baseFeePerGas": extRollup.Header.BaseFee,
	}
}
