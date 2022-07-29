package rpcclientlib

import (
	"crypto/rand"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common/log"
)

type RPCMethod uint8

const (
	http           = "http://"
	reqJSONKeyFrom = "from"

	RPCGetID                   = "obscuro_getID"
	RPCGetCurrentBlockHead     = "obscuro_getCurrentBlockHead"
	RPCGetBlockHeaderByHash    = "obscuro_getBlockHeaderByHash"
	RPCGetCurrentRollupHead    = "obscuro_getCurrentRollupHead"
	RPCGetRollupHeader         = "obscuro_getRollupHeader"
	RPCGetRollupHeaderByNumber = "obscuro_getRollupHeaderByNumber"
	RPCGetRollup               = "obscuro_getRollup"
	RPCNonce                   = "obscuro_nonce"
	RPCAddViewingKey           = "obscuro_addViewingKey"
	RPCGetRollupForTx          = "obscuro_getRollupForTx"
	RPCStopHost                = "obscuro_stopHost"
	RPCCall                    = "eth_call"
	RPCChainID                 = "eth_chainId"
	RPCGetBalance              = "eth_getBalance"
	RPCGetTransactionByHash    = "eth_getTransactionByHash"
	RPCGetTxReceipt            = "eth_getTransactionReceipt"
	RPCSendRawTransaction      = "eth_sendRawTransaction"

	// todo: this is a convenience for testnet testing and will eventually be retrieved from the L1
	enclavePublicKeyHex = "034d3b7e63a8bcd532ee3d1d6ecad9d67fca7821981a044551f0f0cbec74d0bc5e"
)

// for these methods, the RPC method's requests and responses should be encrypted
var sensitiveMethods = []string{RPCCall, RPCGetBalance, RPCGetTxReceipt, RPCSendRawTransaction, RPCGetTransactionByHash}

// Client is used by client applications to interact with the Obscuro node.
type Client interface {
	// Call executes the named method via RPC.
	Call(result interface{}, method string, args ...interface{}) error
	// Stop closes the client.
	Stop()

	// SetViewingKey sets the current viewing key on the client to be used for all sensitive request decryption
	SetViewingKey(viewingKey *ecies.PrivateKey, pubKeyBytes []byte)

	// RegisterViewingKey takes a signature for the public key, it verifies the signed public key matches the currently set private viewing key
	//	and then submits it to the enclave
	RegisterViewingKey(signerAddr common.Address, signature []byte) error
}

// RPCClient is a Client implementation that wraps rpc.Client to make calls.
type networkClient struct {
	rpcClient        *rpc.Client
	enclavePublicKey *ecies.PublicKey
	// todo: add support for multiple keys on the same client?
	viewingPrivKey *ecies.PrivateKey // private viewing key to use for decrypting sensitive requests
	viewingPubKey  []byte            // public viewing key, submitted to the enclave
	viewingKeyAddr common.Address    // address that generated the public key
}

func NewClient(address string) Client {
	// todo: this is a convenience for testnet but needs to replaced by a parameter and/or retrieved from the target host
	enclPubECDSA, err := crypto.DecompressPubkey(common.Hex2Bytes(enclavePublicKeyHex))
	if err != nil {
		panic(err)
	}
	enclavePublicKey := ecies.ImportECDSAPublic(enclPubECDSA)

	rpcClient, err := rpc.Dial(http + address)
	if err != nil {
		log.Panic("could not create RPC client on %s. Cause: %s", http+address, err)
	}

	return &networkClient{
		rpcClient:        rpcClient,
		enclavePublicKey: enclavePublicKey,
	}
}

// Call handles JSON rpc requests - if the method is sensitive it will encrypt the args before sending the request and
//	then decrypts the response before returning.
// The result must be a pointer so that package json can unmarshal into it. You can also pass nil, in which case the result is ignored.
func (c *networkClient) Call(result interface{}, method string, args ...interface{}) error {
	if !isSensitive(method) {
		// for non-sensitive methods or when viewing keys are disabled we just delegate directly to the geth RPC client
		return c.rpcClient.Call(&result, method, args...)
	}

	// we setup a generic rawResult to receive the response (then we can decrypt it as necessary into the requested result type)
	var rawResult interface{}
	if result == nil {
		// we set result receiver to nil if caller set nil since no return needs to be set
		rawResult = nil
	}

	var err error
	if method == RPCCall {
		// RPCCall is a sensitive method that requires a viewing key lookup but the 'from' field is not mandatory in geth
		//	and is often not included from metamask etc. So we ensure it is populated here.
		args, err = c.ensureCallParamsHaveFromAddress(method, args)
		if err != nil {
			return err
		}
	}

	// encode the params into a json blob and encrypt them
	var encryptedParams []byte
	encryptedParams, err = c.encryptArgs(args...)
	if err != nil {
		return fmt.Errorf("failed to encrypt args for %s call - %w", method, err)
	}
	err = c.rpcClient.Call(&rawResult, method, encryptedParams)
	if err != nil {
		return fmt.Errorf("%s rpc call failed - %w", method, err)
	}

	// if caller not interested in response, we're done
	if result == nil {
		return nil
	}

	// method is sensitive, so we decrypt it before unmarshalling the result
	decrypted, err := c.decryptResponse(rawResult)
	if err != nil {
		return fmt.Errorf("failed to decrypt args for %s call - %w", method, err)
	}

	// process the decrypted result to get the desired type and set it on the result pointer
	err = c.setResult(decrypted, result)
	if err != nil {
		return fmt.Errorf("failed to extract result from %s response: %w", method, err)
	}

	return nil
}

func (c *networkClient) encryptArgs(args ...interface{}) ([]byte, error) {
	if len(args) == 0 {
		return nil, nil
	}

	paramsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("could not json encode request params: %w", err)
	}

	return c.encryptParamBytes(paramsJSON)
}

func (c *networkClient) encryptParamBytes(params []byte) ([]byte, error) {
	encryptedParams, err := ecies.Encrypt(rand.Reader, c.enclavePublicKey, params, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt the following request params with enclave public key: %s. Cause: %w", params, err)
	}
	return encryptedParams, nil
}

func (c *networkClient) decryptResponse(resultBlob interface{}) ([]byte, error) {
	if c.viewingPrivKey == nil {
		return nil, fmt.Errorf("cannot decrypt response, viewing key has not been setup")
	}
	resultStr, ok := resultBlob.(string)
	if !ok {
		return nil, fmt.Errorf("expected hex string but result was of type %t instead, with value %s", resultBlob, resultBlob)
	}
	encryptedResult := common.Hex2Bytes(resultStr)
	decryptedResult, err := c.viewingPrivKey.Decrypt(encryptedResult, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt the following response with viewing key: %s. Cause: %w", resultStr, err)
	}

	return decryptedResult, nil
}

// setResult tries to cast/unmarshal data into the result pointer, based on its type
func (c *networkClient) setResult(data []byte, result interface{}) error {
	switch result := result.(type) {
	case *string:
		*result = string(data)
		return nil

	case *interface{}:
		err := json.Unmarshal(data, result)
		if err != nil {
			// if unmarshal failed with generic return we can try to send it back as a string
			*result = string(data)
		}
		return nil

	default:
		// for any other type we attempt to json unmarshal it
		return json.Unmarshal(data, result)
	}
}

func (c *networkClient) Stop() {
	c.rpcClient.Close()
}

func (c *networkClient) SetViewingKey(viewingKey *ecies.PrivateKey, viewingPubKeyBytes []byte) {
	c.viewingPrivKey = viewingKey
	c.viewingPubKey = viewingPubKeyBytes
}

func (c *networkClient) RegisterViewingKey(signerAddr common.Address, signature []byte) error {
	c.viewingKeyAddr = signerAddr

	// We encrypt the viewing key bytes
	encryptedViewingKeyBytes, err := ecies.Encrypt(rand.Reader, c.enclavePublicKey, c.viewingPubKey, nil, nil)
	if err != nil {
		return fmt.Errorf("could not encrypt viewing key with enclave public key: %w", err)
	}

	var rpcErr error
	err = c.Call(&rpcErr, RPCAddViewingKey, encryptedViewingKeyBytes, signature)
	if err != nil {
		return fmt.Errorf("could not add viewing key: %w", err)
	}
	return nil
}

// enclave requires a from address to be set for the viewing key encryption but sources like metamask often don't set it
func (c *networkClient) ensureCallParamsHaveFromAddress(method string, args []interface{}) ([]interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("expected %s params to have a 'from' field but no params found", RPCCall)
	}

	callParams, ok := args[0].(map[string]interface{})
	if !ok {
		callParamsJSON, ok := args[0].([]byte)
		if !ok {
			return nil, fmt.Errorf("expected %s first param to be a map or json encoded bytes but was %t", method, args[0])
		}

		err := json.Unmarshal(callParamsJSON, &callParams)
		if err != nil {
			return nil, fmt.Errorf("expected %s first param to be a map or json encoded bytes, failed to decode: %w", method, err)
		}
	}

	if callParams[reqJSONKeyFrom] != nil {
		// We only modify `eth_call` requests where the `from` field is not set.
		return args, nil
	}
	callParams[reqJSONKeyFrom] = c.viewingKeyAddr

	args[0] = callParams

	return args, nil
}

// isSensitive indicates whether the RPC method's requests and responses should be encrypted.
func isSensitive(method interface{}) bool {
	for _, m := range sensitiveMethods {
		if m == method {
			return true
		}
	}
	return false
}
