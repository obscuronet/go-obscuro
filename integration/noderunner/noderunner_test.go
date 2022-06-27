package noderunner

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/config/defaultconfig"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/obscuro-playground/integration"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/enclaverunner"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host/hostrunner"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"
)

// TODO - Use the HostRunner/EnclaveRunner methods in the socket-based integration tests, and retire this smoketest.

const (
	testLogs            = "../.build/noderunner/"
	defaultWsPortOffset = 100
	localhost           = "127.0.0.1"
)

// A smoke test to check that we can stand up a standalone Obscuro host and enclave.
func TestCanStartStandaloneObscuroHostAndEnclave(t *testing.T) {
	setupTestLog()

	startPort := integration.StartPortNodeRunnerTest
	enclaveAddr := fmt.Sprintf("%s:%d", localhost, startPort)
	rpcServerAddr := fmt.Sprintf("%s:%d", localhost, startPort+1)
	gethPort := startPort + 2
	gethWebsocketPort := gethPort + defaultWsPortOffset

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	hostConfig := defaultconfig.DefaultHostConfig()
	hostConfig.PrivateKeyString = hex.EncodeToString(crypto.FromECDSA(privateKey))
	hostConfig.EnclaveRPCAddress = enclaveAddr
	hostConfig.ClientRPCPortHTTP = startPort + 1
	hostConfig.L1NodeWebsocketPort = uint(gethWebsocketPort)

	enclaveConfig := defaultconfig.DefaultEnclaveConfig()
	enclaveConfig.Address = enclaveAddr
	dummyContractAddress := common.BytesToAddress([]byte("AA"))
	enclaveConfig.ERC20ContractAddresses = []*common.Address{&dummyContractAddress, &dummyContractAddress}

	gethBinaryPath, err := gethnetwork.EnsureBinariesExist(gethnetwork.LatestVersion)
	if err != nil {
		t.Fatal(err)
	}
	network := gethnetwork.NewGethNetwork(int(gethPort), int(gethWebsocketPort), gethBinaryPath, 1, 1, []string{address.String()})
	defer network.StopNodes()

	go enclaverunner.RunEnclave(enclaveConfig)
	go hostrunner.RunHost(hostConfig)
	obscuroClient := obscuroclient.NewClient(rpcServerAddr)
	defer teardown(obscuroClient, rpcServerAddr)

	// We sleep to give the network time to produce some blocks.
	time.Sleep(3 * time.Second)

	// we wait to ensure the RPC endpoint is up
	wait := 60 // max wait in seconds
	for !tcpConnectionAvailable(rpcServerAddr) {
		if wait == 0 {
			t.Fatal("RPC client server never became available")
		}
		time.Sleep(time.Second)
		wait--
	}

	counter := 0
	// We retry 20 times to check if the network has produced any blocks, sleeping half a second between each attempt.
	for counter < 20 {
		counter++
		time.Sleep(500 * time.Millisecond)

		var result types.Header
		err = obscuroClient.Call(&result, obscuroclient.RPCGetCurrentBlockHead)
		if err == nil && result.Number.Uint64() > 0 {
			return
		}
	}

	t.Fatal("Zero blocks have been produced after ten seconds. Something is wrong.")
}

func teardown(obscuroClient obscuroclient.Client, rpcServerAddr string) {
	obscuroClient.Call(nil, obscuroclient.RPCStopHost) //nolint:errcheck

	// We wait for the client server port to be closed.
	wait := 0
	for tcpConnectionAvailable(rpcServerAddr) {
		if wait == 20 { // max wait in seconds
			panic(fmt.Sprintf("RPC client server had not shut down after %d seconds", wait))
		}
		time.Sleep(time.Second)
		wait++
	}
}

func tcpConnectionAvailable(addr string) bool {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return false
	}
	_ = conn.Close()
	// we don't worry about failure while closing, it connected successfully so let test proceed
	return true
}

func setupTestLog() *os.File {
	// create a folder specific for the test
	err := os.MkdirAll(testLogs, 0o700)
	if err != nil {
		panic(err)
	}
	timeFormatted := time.Now().Format("2006-01-02_15-04-05")
	f, err := os.CreateTemp(testLogs, fmt.Sprintf("noderunner-%s-*.txt", timeFormatted))
	if err != nil {
		panic(err)
	}
	log.OutputToFile(f)
	return f
}
