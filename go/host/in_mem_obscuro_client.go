package host

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/obscuro-playground/go/common"
	"github.com/obscuronet/obscuro-playground/go/rpcclientlib"
)

// An in-memory implementation of `rpcclientlib.Client` that speaks directly to the node.
type inMemObscuroClient struct {
	obscuroAPI ObscuroAPI
	ethAPI     EthereumAPI
}

func NewInMemObscuroClient(host *Node) rpcclientlib.Client {
	return &inMemObscuroClient{
		obscuroAPI: *NewObscuroAPI(host),
		ethAPI:     *NewEthereumAPI(host),
	}
}

// Call bypasses RPC, and invokes methods on the node directly.
func (c *inMemObscuroClient) Call(result interface{}, method string, args ...interface{}) error { //nolint:gocognit
	switch method {
	case rpcclientlib.RPCGetID:
		*result.(*gethcommon.Address) = c.obscuroAPI.GetID()

	case rpcclientlib.RPCSendTransactionEncrypted:
		if len(args) != 1 {
			return fmt.Errorf("expected 1 arg to %s, got %d", rpcclientlib.RPCSendTransactionEncrypted, len(args))
		}
		tx, ok := args[0].(common.EncryptedTx)
		if !ok {
			return fmt.Errorf("arg to %s was not of expected type EncryptedTx", rpcclientlib.RPCSendTransactionEncrypted)
		}

		c.obscuroAPI.SendTransactionEncrypted(tx)

	case rpcclientlib.RPCGetCurrentBlockHead:
		*result.(**types.Header) = c.obscuroAPI.GetCurrentBlockHead()

	case rpcclientlib.RPCGetCurrentRollupHead:
		*result.(**common.Header) = c.obscuroAPI.GetCurrentRollupHead()

	case rpcclientlib.RPCGetRollupHeader:
		if len(args) != 1 {
			return fmt.Errorf("expected 1 arg to %s, got %d", rpcclientlib.RPCGetRollupHeader, len(args))
		}
		hash, ok := args[0].(gethcommon.Hash)
		if !ok {
			return fmt.Errorf("arg to %s was not of expected type common.Hash", rpcclientlib.RPCGetRollupHeader)
		}

		*result.(**common.Header) = c.obscuroAPI.GetRollupHeader(hash)

	case rpcclientlib.RPCGetRollup:
		if len(args) != 1 {
			return fmt.Errorf("expected 1 arg to %s, got %d", rpcclientlib.RPCGetRollup, len(args))
		}
		hash, ok := args[0].(gethcommon.Hash)
		if !ok {
			return fmt.Errorf("arg to %s was not of expected type common.Hash", rpcclientlib.RPCGetRollup)
		}

		*result.(**common.ExtRollup) = c.obscuroAPI.GetRollup(hash)

	case rpcclientlib.RPCGetTransaction:
		if len(args) != 1 {
			return fmt.Errorf("expected 1 arg to %s, got %d", rpcclientlib.RPCGetTransaction, len(args))
		}
		hash, ok := args[0].(gethcommon.Hash)
		if !ok {
			return fmt.Errorf("arg to %s was not of expected type common.Hash", rpcclientlib.RPCGetTransaction)
		}

		*result.(**common.L2Tx) = c.obscuroAPI.GetTransaction(hash)

	case rpcclientlib.RPCCall:
		if len(args) != 1 {
			return fmt.Errorf("expected 1 arg to %s, got %d", rpcclientlib.RPCCall, len(args))
		}
		params, ok := args[0].([]byte)
		if !ok {
			return fmt.Errorf("arg 1 to %s was not of expected type []byte]", rpcclientlib.RPCCall)
		}

		encryptedResponse, err := c.ethAPI.Call(context.Background(), params)
		if err != nil {
			return fmt.Errorf("off-chain call failed. Cause: %w", err)
		}
		*result.(*string) = encryptedResponse

	case rpcclientlib.RPCGetTxCount:
		if len(args) != 1 {
			return fmt.Errorf("expected 1 arg to %s, got %d", rpcclientlib.RPCGetTxCount, len(args))
		}
		params, ok := args[0].([]byte)
		if !ok {
			return fmt.Errorf("arg to %s was not of expected type common.Address", rpcclientlib.RPCGetTxCount)
		}

		encryptedResponse, err := c.ethAPI.GetTransactionCount(context.Background(), params)
		if err != nil {
			return fmt.Errorf("GetTransactionCount call failed. Cause: %w", err)
		}

		*result.(*string) = encryptedResponse

	case rpcclientlib.RPCGetTxReceipt:
		if len(args) != 1 {
			return fmt.Errorf("expected 1 arg to %s, got %d", rpcclientlib.RPCGetTxReceipt, len(args))
		}
		params, ok := args[0].([]byte)
		if !ok {
			return fmt.Errorf("arg 1 to %s was not of expected type []byte]", rpcclientlib.RPCGetTxReceipt)
		}

		encryptedResponse, err := c.ethAPI.GetTransactionReceipt(context.Background(), params)
		if err != nil {
			return fmt.Errorf("getTransactionReceipt call failed. Cause: %w", err)
		}
		*result.(*string) = encryptedResponse

	case rpcclientlib.RPCStopHost:
		c.obscuroAPI.StopHost()

	default:
		return fmt.Errorf("RPC method %s is unknown", method)
	}

	return nil
}

func (c *inMemObscuroClient) Stop() {
	// There is no RPC connection to close.
}
