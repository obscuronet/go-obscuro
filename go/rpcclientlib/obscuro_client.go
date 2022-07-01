package rpcclientlib

import (
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/obscuro-playground/go/common/log"
)

type RPCMethod uint8

const (
	http = "http://"

	RPCGetID                    = "obscuro_getID"
	RPCSendTransactionEncrypted = "obscuro_sendTransactionEncrypted"
	RPCGetCurrentBlockHead      = "obscuro_getCurrentBlockHead"
	RPCGetCurrentRollupHead     = "obscuro_getCurrentRollupHead"
	RPCGetRollupHeader          = "obscuro_getRollupHeader"
	RPCGetRollup                = "obscuro_getRollup"
	RPCGetTransaction           = "obscuro_getTransaction"
	RPCAddViewingKey            = "obscuro_addViewingKey"
	RPCStopHost                 = "obscuro_stopHost"
	RPCCall                     = "eth_call"
	RPCGetTxReceipt             = "eth_getTransactionReceipt"
	RPCGetTxCount               = "eth_getTransactionCount"
)

// Client is used by client applications to interact with the Obscuro node.
type Client interface {
	// Call executes the named method via RPC.
	Call(result interface{}, method string, args ...interface{}) error
	// Stop closes the client.
	Stop()
}

// A Client implementation that wraps rpc.Client to make calls.
type clientImpl struct {
	rpcClient *rpc.Client
}

func NewClient(address string) Client {
	rpcClient, err := rpc.Dial(http + address)
	if err != nil {
		log.Panic("could not create RPC client on %s. Cause: %s", http+address, err)
	}

	return &clientImpl{
		rpcClient: rpcClient,
	}
}

func (c *clientImpl) Call(result interface{}, method string, args ...interface{}) error {
	return c.rpcClient.Call(&result, method, args...)
}

func (c *clientImpl) Stop() {
	c.rpcClient.Close()
}
