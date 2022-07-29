package host

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/crypto/ecies"

	"github.com/ethereum/go-ethereum/core/types"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/rpcclientlib"
)

const (
	// todo: this is a convenience for testnet testing and will eventually be retrieved from the L1
	enclavePublicKeyHex = "034d3b7e63a8bcd532ee3d1d6ecad9d67fca7821981a044551f0f0cbec74d0bc5e"
)

// An in-memory implementation of `rpcclientlib.Client` that speaks directly to the node.
type inMemObscuroClient struct {
	obscuroAPI       *ObscuroAPI
	ethAPI           *EthereumAPI
	enclavePublicKey *ecies.PublicKey
}

func NewInMemObscuroClient(host *Node) rpcclientlib.Client {
	// todo: this is a convenience for testnet but needs to replaced by a parameter and/or retrieved from the target host
	enclPubECDSA, err := crypto.DecompressPubkey(gethcommon.Hex2Bytes(enclavePublicKeyHex))
	if err != nil {
		panic(err)
	}
	enclPubKey := ecies.ImportECDSAPublic(enclPubECDSA)

	return &inMemObscuroClient{
		obscuroAPI:       NewObscuroAPI(host),
		ethAPI:           NewEthereumAPI(host),
		enclavePublicKey: enclPubKey,
	}
}

// Call bypasses RPC, and invokes methods on the node directly.
func (c *inMemObscuroClient) Call(result interface{}, method string, args ...interface{}) error {
	switch method {
	case rpcclientlib.RPCGetID:
		*result.(*gethcommon.Address) = c.obscuroAPI.GetID()

	case rpcclientlib.RPCSendRawTransaction:
		return c.sendRawTransaction(args)

	case rpcclientlib.RPCGetCurrentBlockHead:
		*result.(**types.Header) = c.obscuroAPI.GetCurrentBlockHead()

	case rpcclientlib.RPCGetCurrentRollupHead:
		*result.(**common.Header) = c.obscuroAPI.GetCurrentRollupHead()

	case rpcclientlib.RPCGetRollupHeader:
		return c.getRollupHeader(result, args)

	case rpcclientlib.RPCGetRollup:
		return c.getRollup(result, args)

	case rpcclientlib.RPCGetTransactionByHash:
		return c.getTransactionByHash(result, args)

	case rpcclientlib.RPCCall:
		return c.rpcCall(result, args)

	case rpcclientlib.RPCNonce:
		return c.getNonce(result, args)

	case rpcclientlib.RPCGetTxReceipt:
		return c.getTransactionReceipt(result, args)

	case rpcclientlib.RPCStopHost:
		c.obscuroAPI.StopHost()

	default:
		return fmt.Errorf("RPC method %s is unknown", method)
	}

	return nil
}

func (c *inMemObscuroClient) sendRawTransaction(args []interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 arg to %s, got %d", rpcclientlib.RPCSendRawTransaction, len(args))
	}
	bytes, err := json.Marshal(args)
	if err != nil {
		return fmt.Errorf("failed to marshal the rpc args for %s - %w", rpcclientlib.RPCSendRawTransaction, err)
	}
	enc, err := c.encryptParamBytes(bytes)
	if err != nil {
		return err
	}
	_, err = c.ethAPI.SendRawTransaction(context.Background(), enc)
	return err
}

func (c *inMemObscuroClient) getTransactionByHash(result interface{}, args []interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 arg to %s, got %d", rpcclientlib.RPCGetTransactionByHash, len(args))
	}
	bytes, err := json.Marshal(args)
	if err != nil {
		return fmt.Errorf("failed to marshal the rpc args for %s - %w", rpcclientlib.RPCGetTransactionByHash, err)
	}
	enc, err := c.encryptParamBytes(bytes)
	if err != nil {
		return err
	}
	encryptedTx, err := c.ethAPI.GetTransactionByHash(context.Background(), enc)
	if err != nil {
		return fmt.Errorf("`eth_getTransactionByHash` call failed. Cause: %w", err)
	}
	*result.(*string) = encryptedTx
	return nil
}

func (c *inMemObscuroClient) rpcCall(result interface{}, args []interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 arg to %s, got %d", rpcclientlib.RPCCall, len(args))
	}
	bytes, err := json.Marshal(args)
	if err != nil {
		return fmt.Errorf("failed to marshal the rpc args for %s - %w", rpcclientlib.RPCCall, err)
	}
	enc, err := c.encryptParamBytes(bytes)
	if err != nil {
		return err
	}
	encryptedResponse, err := c.ethAPI.Call(context.Background(), enc)
	if err != nil {
		return fmt.Errorf("`eth_call` call failed. Cause: %w", err)
	}
	*result.(*string) = encryptedResponse
	return nil
}

func (c *inMemObscuroClient) getTransactionReceipt(result interface{}, args []interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 arg to %s, got %d", rpcclientlib.RPCGetTxReceipt, len(args))
	}
	bytes, err := json.Marshal(args)
	if err != nil {
		return fmt.Errorf("failed to marshal the rpc args for %s - %w", rpcclientlib.RPCGetTxReceipt, err)
	}
	enc, err := c.encryptParamBytes(bytes)
	if err != nil {
		return err
	}
	encryptedResponse, err := c.ethAPI.GetTransactionReceipt(context.Background(), enc)
	if err != nil {
		return fmt.Errorf("`obscuro_getTransactionReceipt` call failed. Cause: %w", err)
	}
	*result.(*string) = *encryptedResponse
	return nil
}

func (c *inMemObscuroClient) getRollupHeader(result interface{}, args []interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 arg to %s, got %d", rpcclientlib.RPCGetRollupHeader, len(args))
	}
	// we expect a hex string representation of the hash, since that's what gets sent over RPC
	hashStr, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("arg to %s was not of expected type string", rpcclientlib.RPCGetRollupHeader)
	}
	hash := gethcommon.HexToHash(hashStr)

	*result.(**common.Header) = c.obscuroAPI.GetRollupHeader(hash)
	return nil
}

func (c *inMemObscuroClient) getRollup(result interface{}, args []interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 arg to %s, got %d", rpcclientlib.RPCGetRollup, len(args))
	}
	hash, ok := args[0].(gethcommon.Hash)
	if !ok {
		return fmt.Errorf("arg to %s was not of expected type common.Hash", rpcclientlib.RPCGetRollup)
	}

	extRollup, err := c.obscuroAPI.GetRollup(hash)
	if err != nil {
		return fmt.Errorf("`obscuro_getRollup` call failed. Cause: %w", err)
	}
	*result.(**common.ExtRollup) = extRollup
	return nil
}

func (c *inMemObscuroClient) getNonce(result interface{}, args []interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 arg to %s, got %d", rpcclientlib.RPCNonce, len(args))
	}
	address, ok := args[0].(gethcommon.Address)
	if !ok {
		return fmt.Errorf("arg to %s was not of expected type common.Address", rpcclientlib.RPCNonce)
	}

	*result.(*uint64) = c.obscuroAPI.Nonce(address)
	return nil
}

func (c *inMemObscuroClient) Stop() {
	// There is no RPC connection to close.
}

func (c *inMemObscuroClient) SetViewingKey(_ *ecies.PrivateKey, _ []byte) {
	panic("viewing key encryption/decryption is not currently supported by in-memory obscuro-client")
}

func (c *inMemObscuroClient) RegisterViewingKey(_ gethcommon.Address, _ []byte) error {
	panic("viewing key encryption/decryption is not currently supported by in-memory obscuro-client")
}

func (c *inMemObscuroClient) encryptParamBytes(params []byte) ([]byte, error) {
	encryptedParams, err := ecies.Encrypt(rand.Reader, c.enclavePublicKey, params, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt request params with enclave public key: %w", err)
	}
	return encryptedParams, nil
}
