package walletextension

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/obscuronet/obscuro-playground/go/enclave/rollupchain"
	"github.com/obscuronet/obscuro-playground/go/ethadapter/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/go/wallet"
	"github.com/obscuronet/obscuro-playground/integration/erc20contract"
	"github.com/obscuronet/obscuro-playground/integration/simulation"

	"github.com/obscuronet/obscuro-playground/go/common/log/logutil"

	"github.com/obscuronet/obscuro-playground/go/enclave/bridge"

	"github.com/obscuronet/obscuro-playground/go/enclave/rpcencryptionmanager"

	"github.com/obscuronet/obscuro-playground/tools/walletextension"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/obscuro-playground/integration"
	"github.com/obscuronet/obscuro-playground/integration/ethereummock"
	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
	"github.com/obscuronet/obscuro-playground/integration/simulation/params"
	"github.com/obscuronet/obscuro-playground/integration/simulation/stats"
)

const (
	testLogs     = "../.build/wallet_extension/"
	l2ChainIDHex = "0x309"

	reqJSONMethodChainID = "eth_chainId"
	reqJSONKeyTo         = "to"
	reqJSONKeyFrom       = "from"
	reqJSONKeyData       = "data"
	latestBlock          = "latest"
	errInsecure          = "enclave could not respond securely to %s request"

	networkStartPort = integration.StartPortWalletExtensionTest + 1
	nodeRPCHTTPPort  = integration.StartPortWalletExtensionTest + 1 + network.DefaultHostRPCHTTPOffset
	nodeRPCWSPort    = integration.StartPortWalletExtensionTest + 1 + network.DefaultHostRPCWSOffset
	httpProtocol     = "http://"
)

var (
	walletExtensionAddr   = fmt.Sprintf("%s:%d", network.Localhost, integration.StartPortWalletExtensionTest)
	walletExtensionConfig = walletextension.Config{
		WalletExtensionPort:     int(integration.StartPortWalletExtensionTest),
		NodeRPCHTTPAddress:      fmt.Sprintf("%s:%d", network.Localhost, nodeRPCHTTPPort),
		NodeRPCWebsocketAddress: fmt.Sprintf("%s:%d", network.Localhost, nodeRPCWSPort),
	}
	dummyAccountAddress = common.HexToAddress("0x8D97689C9818892B700e27F316cc3E41e17fBeb9")
)

func TestCanMakeNonSensitiveRequestWithoutSubmittingViewingKey(t *testing.T) {
	setupWalletTestLog("req-no-viewing-key")

	walletExtension := walletextension.NewWalletExtension(walletExtensionConfig)
	defer walletExtension.Shutdown()
	go walletExtension.Serve(walletExtensionAddr)
	waitForWalletExtension(t, walletExtensionAddr)

	stopHandle, err := createObscuroNetwork(t)
	defer stopHandle()
	if err != nil {
		t.Fatalf("failed to create test Obscuro network. Cause: %s", err)
	}

	respJSON := makeEthJSONReqAsJSON(t, walletExtensionAddr, reqJSONMethodChainID, []string{})

	if respJSON[walletextension.RespJSONKeyResult] != l2ChainIDHex {
		t.Fatalf("Expected chainId of %s, got %s", l2ChainIDHex, respJSON[walletextension.RespJSONKeyResult])
	}
}

func TestCannotGetBalanceWithoutSubmittingViewingKey(t *testing.T) {
	setupWalletTestLog("bal-no-viewing-key")

	walletExtension := walletextension.NewWalletExtension(walletExtensionConfig)
	defer walletExtension.Shutdown()
	go walletExtension.Serve(walletExtensionAddr)
	waitForWalletExtension(t, walletExtensionAddr)

	stopHandle, err := createObscuroNetwork(t)
	defer stopHandle()
	if err != nil {
		t.Fatalf("failed to create test Obscuro network. Cause: %s", err)
	}

	respBody := makeEthJSONReq(t, walletExtensionAddr, walletextension.ReqJSONMethodGetBalance, []string{dummyAccountAddress.Hex(), latestBlock})

	expectedErr := fmt.Sprintf(errInsecure, walletextension.ReqJSONMethodGetBalance)
	if !strings.Contains(string(respBody), expectedErr) {
		t.Fatalf("Expected error message to contain \"%s\", got \"%s\"", expectedErr, respBody)
	}
}

func TestCanGetOwnBalanceAfterSubmittingViewingKey(t *testing.T) {
	setupWalletTestLog("bal-with-viewing-key")

	walletExtension := walletextension.NewWalletExtension(walletExtensionConfig)
	defer walletExtension.Shutdown()
	go walletExtension.Serve(walletExtensionAddr)
	waitForWalletExtension(t, walletExtensionAddr)

	stopHandle, err := createObscuroNetwork(t)
	defer stopHandle()
	if err != nil {
		t.Fatalf("failed to create test Obscuro network. Cause: %s", err)
	}

	// We submit a viewing key for a random account.
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	accountAddr := crypto.PubkeyToAddress(privateKey.PublicKey).String()

	generateAndSubmitViewingKey(t, walletExtensionAddr, privateKey)

	getBalanceJSON := makeEthJSONReqAsJSON(t, walletExtensionAddr, walletextension.ReqJSONMethodGetBalance, []string{accountAddr, latestBlock})

	if getBalanceJSON[walletextension.RespJSONKeyResult] != rollupchain.DummyBalance {
		t.Fatalf("Expected balance of %s, got %s", rollupchain.DummyBalance, getBalanceJSON[walletextension.RespJSONKeyResult])
	}
}

func TestCannotGetAnothersBalanceAfterSubmittingViewingKey(t *testing.T) {
	setupWalletTestLog("others-bal-with-viewing-key")

	walletExtension := walletextension.NewWalletExtension(walletExtensionConfig)
	defer walletExtension.Shutdown()
	go walletExtension.Serve(walletExtensionAddr)
	waitForWalletExtension(t, walletExtensionAddr)

	stopHandle, err := createObscuroNetwork(t)
	defer stopHandle()
	if err != nil {
		t.Fatalf("failed to create test Obscuro network. Cause: %s", err)
	}

	// We submit a viewing key for a random account.
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	generateAndSubmitViewingKey(t, walletExtensionAddr, privateKey)

	respBody := makeEthJSONReq(t, walletExtensionAddr, walletextension.ReqJSONMethodGetBalance, []string{dummyAccountAddress.Hex(), latestBlock})

	expectedErr := fmt.Sprintf(errInsecure, walletextension.ReqJSONMethodGetBalance)
	if !strings.Contains(string(respBody), expectedErr) {
		t.Fatalf("Expected error message to contain \"%s\", got \"%s\"", expectedErr, respBody)
	}
}

func TestCannotCallWithoutSubmittingViewingKey(t *testing.T) {
	setupWalletTestLog("tx-no-viewing-key")

	walletExtension := walletextension.NewWalletExtension(walletExtensionConfig)
	defer walletExtension.Shutdown()
	go walletExtension.Serve(walletExtensionAddr)
	waitForWalletExtension(t, walletExtensionAddr)

	stopHandle, err := createObscuroNetwork(t)
	defer stopHandle()
	if err != nil {
		t.Fatalf("failed to create test Obscuro network. Cause: %s", err)
	}

	time.Sleep(2 * time.Second) // We wait for the deployment of the ERC20 contract to the Obscuro network.

	// We submit a viewing key for a random account.
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	accountAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	// We submit a transaction to the Obscuro ERC20 contract. By transferring an amount of zero, we avoid the need to
	// deposit any funds in the ERC20 contract.
	transferTxBytes := erc20contractlib.CreateTransferTxData(accountAddress, 0)
	reqParams := map[string]interface{}{
		reqJSONKeyTo:   bridge.WBtcContract,
		reqJSONKeyFrom: accountAddress.String(),
		reqJSONKeyData: "0x" + common.Bytes2Hex(transferTxBytes),
	}
	respBody := makeEthJSONReq(t, walletExtensionAddr, walletextension.ReqJSONMethodCall, []interface{}{reqParams, latestBlock})

	expectedErr := fmt.Sprintf(errInsecure, walletextension.ReqJSONMethodCall)
	if !strings.Contains(string(respBody), expectedErr) {
		t.Fatalf("Expected error message \"%s\", got \"%s\"", expectedErr, respBody)
	}
}

func TestCanCallAfterSubmittingViewingKey(t *testing.T) {
	setupWalletTestLog("tx-with-viewing-key")

	walletExtension := walletextension.NewWalletExtension(walletExtensionConfig)
	defer walletExtension.Shutdown()
	go walletExtension.Serve(walletExtensionAddr)
	waitForWalletExtension(t, walletExtensionAddr)

	stopHandle, err := createObscuroNetwork(t)
	defer stopHandle()
	if err != nil {
		t.Fatalf("failed to create test Obscuro network. Cause: %s", err)
	}

	time.Sleep(2 * time.Second) // We wait for the deployment of the ERC20 contract to the Obscuro network.

	// We submit a viewing key for a random account.
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	accountAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	generateAndSubmitViewingKey(t, walletExtensionAddr, privateKey)

	// We submit a transaction to the Obscuro ERC20 contract. By transferring an amount of zero, we avoid the need to
	// deposit any funds in the ERC20 contract.
	transferTxBytes := erc20contractlib.CreateTransferTxData(accountAddress, 0)
	reqParams := map[string]interface{}{
		reqJSONKeyTo:   bridge.WBtcContract,
		reqJSONKeyFrom: accountAddress.String(),
		reqJSONKeyData: "0x" + common.Bytes2Hex(transferTxBytes),
	}
	callJSON := makeEthJSONReqAsJSON(t, walletExtensionAddr, walletextension.ReqJSONMethodCall, []interface{}{reqParams, latestBlock})

	if callJSON[walletextension.RespJSONKeyResult] != string(rpcencryptionmanager.PlaceholderResult) {
		t.Fatalf("Expected call result of %s, got %s", rpcencryptionmanager.PlaceholderResult, callJSON[walletextension.RespJSONKeyResult])
	}
}

func TestCanCallWithoutSettingFromField(t *testing.T) {
	setupWalletTestLog("tx-no-from-field")

	walletExtension := walletextension.NewWalletExtension(walletExtensionConfig)
	defer walletExtension.Shutdown()
	go walletExtension.Serve(walletExtensionAddr)
	waitForWalletExtension(t, walletExtensionAddr)

	stopHandle, err := createObscuroNetwork(t)
	defer stopHandle()
	if err != nil {
		t.Fatalf("failed to create test Obscuro network. Cause: %s", err)
	}

	time.Sleep(2 * time.Second) // We wait for the deployment of the ERC20 contract to the Obscuro network.

	// We submit a viewing key for a random account.
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	accountAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	generateAndSubmitViewingKey(t, walletExtensionAddr, privateKey)

	// We submit a transaction to the Obscuro ERC20 contract. By transferring an amount of zero, we avoid the need to
	// deposit any funds in the ERC20 contract.
	transferTxBytes := erc20contractlib.CreateTransferTxData(accountAddress, 0)
	reqParams := map[string]interface{}{
		reqJSONKeyTo:   bridge.WBtcContract,
		reqJSONKeyData: "0x" + common.Bytes2Hex(transferTxBytes),
	}
	callJSON := makeEthJSONReqAsJSON(t, walletExtensionAddr, walletextension.ReqJSONMethodCall, []interface{}{reqParams, latestBlock})

	if callJSON[walletextension.RespJSONKeyResult] != string(rpcencryptionmanager.PlaceholderResult) {
		t.Fatalf("Expected call result of %s, got %s", rpcencryptionmanager.PlaceholderResult, callJSON[walletextension.RespJSONKeyResult])
	}
}

func TestCannotCallForAnotherAddressAfterSubmittingViewingKey(t *testing.T) {
	setupWalletTestLog("others-tx-with-viewing-key")

	walletExtension := walletextension.NewWalletExtension(walletExtensionConfig)
	defer walletExtension.Shutdown()
	go walletExtension.Serve(walletExtensionAddr)
	waitForWalletExtension(t, walletExtensionAddr)

	stopHandle, err := createObscuroNetwork(t)
	defer stopHandle()
	if err != nil {
		t.Fatalf("failed to create test Obscuro network. Cause: %s", err)
	}

	time.Sleep(2 * time.Second) // We wait for the deployment of the ERC20 contract to the Obscuro network.

	// We submit a viewing key for a random account.
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	generateAndSubmitViewingKey(t, walletExtensionAddr, privateKey)

	// We submit a transaction to the Obscuro ERC20 contract. By transferring an amount of zero, we avoid the need to
	// deposit any funds in the ERC20 contract.
	transferTxBytes := erc20contractlib.CreateTransferTxData(dummyAccountAddress, 0)
	reqParams := map[string]interface{}{
		reqJSONKeyTo: bridge.WBtcContract,
		// We send the request from a different address than the one we created a viewing key for.
		reqJSONKeyFrom: dummyAccountAddress.Hex(),
		reqJSONKeyData: "0x" + common.Bytes2Hex(transferTxBytes),
	}
	respBody := makeEthJSONReq(t, walletExtensionAddr, walletextension.ReqJSONMethodCall, []interface{}{reqParams, latestBlock})

	expectedErr := fmt.Sprintf(errInsecure, walletextension.ReqJSONMethodCall)
	if !strings.Contains(string(respBody), expectedErr) {
		t.Fatalf("Expected error message \"%s\", got \"%s\"", expectedErr, respBody)
	}
}

// This trio of tests covers both `eth_sendRawTransaction` and `eth_getTransactionReceipt`.

func TestCannotSubmitTxWithoutSubmittingViewingKey(t *testing.T) {
	setupWalletTestLog("submit-tx-no-viewing-key")

	walletExtension := walletextension.NewWalletExtension(walletExtensionConfig)
	defer walletExtension.Shutdown()
	go walletExtension.Serve(walletExtensionAddr)
	waitForWalletExtension(t, walletExtensionAddr)

	stopHandle, err := createObscuroNetwork(t)
	defer stopHandle()
	if err != nil {
		t.Fatalf("failed to create test Obscuro network. Cause: %s", err)
	}

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	txWallet := wallet.NewInMemoryWalletFromPK(big.NewInt(integration.ObscuroChainID), privateKey)
	tx := types.LegacyTx{
		Nonce:    0,
		Gas:      1025_000_000,
		GasPrice: common.Big0,
		Data:     common.Hex2Bytes(erc20contract.ContractByteCode),
	}
	txBinaryHex, err := formatTxForSubmission(txWallet, &tx)
	if err != nil {
		panic(err)
	}
	respBody := makeEthJSONReq(t, walletExtensionAddr, walletextension.ReqJSONMethodSendRawTx, []interface{}{txBinaryHex})

	expectedErr := fmt.Sprintf(errInsecure, walletextension.ReqJSONMethodSendRawTx)
	if !strings.Contains(string(respBody), expectedErr) {
		t.Fatalf("Expected error message \"%s\", got \"%s\"", expectedErr, respBody)
	}
}

func TestCanSubmitTxAndGetTxReceiptAfterSubmittingViewingKey(t *testing.T) {
	setupWalletTestLog("submit-tx-with-viewing-key")

	walletExtension := walletextension.NewWalletExtension(walletExtensionConfig)
	defer walletExtension.Shutdown()
	go walletExtension.Serve(walletExtensionAddr)
	waitForWalletExtension(t, walletExtensionAddr)

	stopHandle, err := createObscuroNetwork(t)
	defer stopHandle()
	if err != nil {
		t.Fatalf("failed to create test Obscuro network. Cause: %s", err)
	}

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	txWallet := wallet.NewInMemoryWalletFromPK(big.NewInt(integration.ObscuroChainID), privateKey)
	generateAndSubmitViewingKey(t, walletExtensionAddr, privateKey)
	tx := types.LegacyTx{
		Nonce:    0,
		Gas:      1025_000_000,
		GasPrice: common.Big0,
		Data:     common.Hex2Bytes(erc20contract.ContractByteCode),
	}
	txBinaryHex, err := formatTxForSubmission(txWallet, &tx)
	if err != nil {
		panic(err)
	}
	sendTxJSON := makeEthJSONReqAsJSON(t, walletExtensionAddr, walletextension.ReqJSONMethodSendRawTx, []interface{}{txBinaryHex})

	time.Sleep(6 * time.Second) // We wait for the deployment of the contract to the Obscuro network.

	// We get the transaction receipt for the Obscuro ERC20 contract.
	txHash, ok := sendTxJSON[walletextension.RespJSONKeyResult].(string)
	if !ok {
		panic("could not retrieve transaction hash from JSON result")
	}
	txReceiptJSON := makeEthJSONReqAsJSON(t, walletExtensionAddr, walletextension.ReqJSONMethodGetTxReceipt, []string{txHash})

	result := fmt.Sprintf("%s", txReceiptJSON[walletextension.RespJSONKeyResult])
	expectedTxHashJSON := fmt.Sprintf("transactionHash:%s", txHash)
	if !strings.Contains(result, expectedTxHashJSON) {
		t.Fatalf("Expected transaction receipt containing %s, got %s", expectedTxHashJSON, result)
	}
}

func TestCannotSubmitTxFromAnotherAddressAfterSubmittingViewingKey(t *testing.T) {
	setupWalletTestLog("others-submit-tx-with-viewing-key")

	walletExtension := walletextension.NewWalletExtension(walletExtensionConfig)
	defer walletExtension.Shutdown()
	go walletExtension.Serve(walletExtensionAddr)
	waitForWalletExtension(t, walletExtensionAddr)

	stopHandle, err := createObscuroNetwork(t)
	defer stopHandle()
	if err != nil {
		t.Fatalf("failed to create test Obscuro network. Cause: %s", err)
	}

	// We submit a viewing key for a random account.
	randomPrivateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	generateAndSubmitViewingKey(t, walletExtensionAddr, randomPrivateKey)

	// We submit a transaction using another account.
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	txWallet := wallet.NewInMemoryWalletFromPK(big.NewInt(integration.ObscuroChainID), privateKey)
	tx := types.LegacyTx{
		Nonce:    0,
		Gas:      1025_000_000,
		GasPrice: common.Big0,
		Data:     common.Hex2Bytes(erc20contract.ContractByteCode),
	}
	txBinaryHex, err := formatTxForSubmission(txWallet, &tx)
	if err != nil {
		panic(err)
	}
	respBody := makeEthJSONReq(t, walletExtensionAddr, walletextension.ReqJSONMethodSendRawTx, []interface{}{txBinaryHex})

	expectedErr := fmt.Sprintf(errInsecure, walletextension.ReqJSONMethodSendRawTx)
	if !strings.Contains(string(respBody), expectedErr) {
		t.Fatalf("Expected error message \"%s\", got \"%s\"", expectedErr, respBody)
	}
}

// Waits for wallet extension to be ready. Times out after three seconds.
func waitForWalletExtension(t *testing.T, walletExtensionAddr string) {
	retries := 30
	for i := 0; i < retries; i++ {
		resp, err := http.Get(httpProtocol + walletExtensionAddr + walletextension.PathReady) //nolint:noctx
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		if err == nil {
			return
		}
		time.Sleep(300 * time.Millisecond)
	}
	t.Fatal("could not establish connection to wallet extension")
}

// Makes an Ethereum JSON RPC request and returns the response body.
func makeEthJSONReq(t *testing.T, walletExtensionAddr string, method string, params interface{}) []byte {
	reqBodyBytes, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      "1",
	})
	if err != nil {
		t.Fatal(err)
	}
	reqBody := bytes.NewBuffer(reqBodyBytes)

	var resp *http.Response
	// We retry for three seconds to handle node start-up time.
	timeout := time.Now().Add(3 * time.Second)
	for i := time.Now(); i.Before(timeout); i = time.Now() {
		resp, err = http.Post(httpProtocol+walletExtensionAddr, "text/html", reqBody) //nolint:noctx
		if err == nil {
			break
		}
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}

	if err != nil {
		t.Fatalf("received error response from wallet extension: %s", err)
	}
	if resp == nil {
		t.Fatal("did not receive a response from the wallet extension")
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	return respBody
}

// Makes an Ethereum JSON RPC request and returns the response body as JSON.
func makeEthJSONReqAsJSON(t *testing.T, walletExtensionAddr string, method string, params interface{}) map[string]interface{} {
	respBody := makeEthJSONReq(t, walletExtensionAddr, method, params)

	if respBody[0] != '{' {
		t.Fatalf("expected JSON response but received: %s", respBody)
	}

	var respBodyJSON map[string]interface{}
	err := json.Unmarshal(respBody, &respBodyJSON)
	if err != nil {
		t.Fatal(err)
	}

	return respBodyJSON
}

// Generates a signed viewing key and submits it to the wallet extension.
func generateAndSubmitViewingKey(t *testing.T, walletExtensionAddr string, accountPrivateKey *ecdsa.PrivateKey) {
	viewingKey := generateViewingKey(t, walletExtensionAddr)
	signature := signViewingKey(t, accountPrivateKey, viewingKey)

	submitViewingKeyBodyBytes, err := json.Marshal(map[string]interface{}{
		"signature": hex.EncodeToString(signature),
	})
	if err != nil {
		t.Fatal(err)
	}
	submitViewingKeyBody := bytes.NewBuffer(submitViewingKeyBodyBytes)
	resp, err := http.Post(httpProtocol+walletExtensionAddr+walletextension.PathSubmitViewingKey, "application/json", submitViewingKeyBody) //nolint:noctx
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Fatalf("request to add viewing key failed with following status: %s", resp.Status)
	}
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
}

// Generates a viewing key.
func generateViewingKey(t *testing.T, walletExtensionAddr string) []byte {
	resp, err := http.Get(httpProtocol + walletExtensionAddr + walletextension.PathGenerateViewingKey) //nolint:noctx
	if err != nil {
		t.Fatal(err)
	}
	viewingKey, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	return viewingKey
}

// Signs a viewing key.
func signViewingKey(t *testing.T, privateKey *ecdsa.PrivateKey, viewingKey []byte) []byte {
	msgToSign := rpcencryptionmanager.ViewingKeySignedMsgPrefix + string(viewingKey)
	signature, err := crypto.Sign(accounts.TextHash([]byte(msgToSign)), privateKey)
	if err != nil {
		t.Fatal(err)
	}

	// We have to transform the V from 0/1 to 27/28, and add the leading "0".
	signature[64] += 27
	signatureWithLeadBytes := append([]byte("0"), signature...)

	return signatureWithLeadBytes
}

// Creates a single-node Obscuro network for testing, and deploys an ERC20 contract to it.
func createObscuroNetwork(t *testing.T) (func(), error) {
	// Create the Obscuro network.
	numberOfNodes := 1
	wallets := params.NewSimWallets(1, numberOfNodes, integration.EthereumChainID, integration.ObscuroChainID)
	simParams := params.SimParams{
		NumberOfNodes:      numberOfNodes,
		AvgBlockDuration:   1 * time.Second,
		AvgGossipPeriod:    1 * time.Second / 3,
		MgmtContractLib:    ethereummock.NewMgmtContractLibMock(),
		ERC20ContractLib:   ethereummock.NewERC20ContractLibMock(),
		Wallets:            wallets,
		StartPort:          int(networkStartPort),
		ViewingKeysEnabled: true,
	}
	simStats := stats.NewStats(simParams.NumberOfNodes)
	obscuroNetwork := network.NewNetworkOfSocketNodes(wallets)
	_, l2Clients, err := obscuroNetwork.Create(&simParams, simStats)
	if err != nil {
		return obscuroNetwork.TearDown, err
	}

	// Deploy an ERC20 contract to the Obscuro network.
	txWallet := wallets.Tokens[bridge.BTC].L2Owner
	generateAndSubmitViewingKey(t, walletExtensionAddr, txWallet.PrivateKey())

	deployContractTx := types.LegacyTx{
		Nonce:    simulation.NextNonce(l2Clients[0], txWallet),
		Gas:      1025_000_000,
		GasPrice: common.Big0,
		Data:     common.Hex2Bytes(erc20contract.ContractByteCode),
	}
	txBinaryHex, err := formatTxForSubmission(txWallet, &deployContractTx)
	if err != nil {
		return obscuroNetwork.TearDown, err
	}
	makeEthJSONReqAsJSON(t, walletExtensionAddr, walletextension.ReqJSONMethodSendRawTx, []interface{}{txBinaryHex})

	return obscuroNetwork.TearDown, nil
}

// Formats a transaction for submission to the enclave.
func formatTxForSubmission(wallet wallet.Wallet, tx types.TxData) (string, error) {
	signedTx, err := wallet.SignTransaction(tx)
	if err != nil {
		return "", err
	}
	// We convert the transaction to the form expected for sending transactions via RPC.
	txBinary, err := signedTx.MarshalBinary()
	if err != nil {
		return "", err
	}
	txBinaryHex := "0x" + common.Bytes2Hex(txBinary)
	if err != nil {
		return "", err
	}

	return txBinaryHex, nil
}

func setupWalletTestLog(testName string) {
	// todo: creating an individual file for every test is very heavy-handed, come up with a better solution?
	logutil.SetupTestLog(&logutil.TestLogCfg{
		LogDir:      testLogs,
		TestType:    "wal-ext",
		TestSubtype: testName,
	})
}
