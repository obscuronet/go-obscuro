package noderunner

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/config"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/obscuro-playground/integration"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/enclaverunner"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host/hostrunner"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"
)

// A smoke test to check that we can stand up a standalone Obscuro host and enclave.
func TestConnectStandaloneToGoerli(t *testing.T) {
	setupTestLog()

	startPort := integration.StartPortNodeRunnerTest
	enclaveAddr := fmt.Sprintf("127.0.0.1:%d", startPort)
	clientServerAddr := fmt.Sprintf("127.0.0.1:%d", startPort+1)
	gethPort := startPort + 2
	gethWebsocketPort := gethPort + defaultWsPortOffset

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	//address := crypto.PubkeyToAddress(privateKey.PublicKey)

	hostConfig := config.DefaultHostConfig()
	hostConfig.PrivateKeyString = hex.EncodeToString(crypto.FromECDSA(privateKey))
	hostConfig.EnclaveRPCAddress = enclaveAddr
	hostConfig.ClientRPCAddress = clientServerAddr
	hostConfig.L1NodeWebsocketPort = uint(gethWebsocketPort)

	enclaveConfig := config.DefaultEnclaveConfig()
	enclaveConfig.Address = enclaveAddr

	gethBinaryPath, err := gethnetwork.EnsureBinariesExist(gethnetwork.LatestVersion)
	if err != nil {
		t.Fatal(err)
	}
	network := gethnetwork.NewGethNetworkGoerli(int(gethPort), int(gethWebsocketPort), gethBinaryPath, 1)
	defer network.StopNodes()

	go enclaverunner.RunEnclave(enclaveConfig)
	go hostrunner.RunHost(hostConfig)
	obscuroClient := obscuroclient.NewClient(clientServerAddr)
	defer teardown(obscuroClient, clientServerAddr)

	// We sleep to give the network time to produce some blocks.
	time.Sleep(3 * time.Second)

	// we wait to ensure the RPC endpoint is up
	wait := 60 // max wait in seconds
	for !tcpConnectionAvailable(clientServerAddr) {
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
