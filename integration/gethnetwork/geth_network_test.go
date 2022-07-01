package gethnetwork

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/obscuronet/obscuro-playground/go/ethadapter"

	"github.com/obscuronet/obscuro-playground/go/config"

	"github.com/obscuronet/obscuro-playground/integration"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/integration/datagenerator"
	"gopkg.in/yaml.v3"
)

const (
	numNodes        = 3
	expectedChainID = "1337"
	genesisChainID  = 1337

	peerCountCmd = "net.peerCount"
	chainIDCmd   = "admin.nodeInfo.protocols.eth.config.chainId"

	defaultWsPortOffset        = 100 // The default offset between a Geth node's HTTP and websocket ports.
	defaultL1ConnectionTimeout = 15 * time.Second

	localhost = "127.0.0.1"
)

var timeout = 15 * time.Second

func TestGethAllNodesJoinSameNetwork(t *testing.T) {
	gethBinaryPath, err := EnsureBinariesExist(LatestVersion)
	if err != nil {
		t.Fatal(err)
	}

	startPort := int(integration.StartPortGethNetworkTest)
	network := NewGethNetwork(startPort, startPort+defaultWsPortOffset, gethBinaryPath, numNodes, 1, nil)
	defer network.StopNodes()

	peerCountStr := network.IssueCommand(0, peerCountCmd)
	peerCount, _ := strconv.Atoi(peerCountStr)

	// nodes don't consider themselves a peer
	if expectedPeers := numNodes - 1; peerCount != expectedPeers {
		t.Fatalf("Wrong number of peers on the network. Found %d, expected %d.", peerCount, expectedPeers)
	}
}

func TestGethGenesisParamsAreUsed(t *testing.T) {
	gethBinaryPath, err := EnsureBinariesExist(LatestVersion)
	if err != nil {
		t.Fatal(err)
	}

	startPort := int(integration.StartPortGethNetworkTest) + numNodes
	network := NewGethNetwork(startPort, startPort+defaultWsPortOffset, gethBinaryPath, numNodes, 1, nil)
	defer network.StopNodes()

	chainID := network.IssueCommand(0, chainIDCmd)
	if chainID != expectedChainID {
		t.Fatalf("Network not using chain ID specified in the genesis file. Found %s, expected %s.", chainID, expectedChainID)
	}
}

func TestGethTransactionCanBeSubmitted(t *testing.T) {
	gethBinaryPath, err := EnsureBinariesExist(LatestVersion)
	if err != nil {
		t.Fatal(err)
	}

	startPort := int(integration.StartPortGethNetworkTest) + numNodes*2
	network := NewGethNetwork(startPort, startPort+defaultWsPortOffset, gethBinaryPath, numNodes, 1, nil)
	defer network.StopNodes()

	account := network.addresses[0]
	tx := fmt.Sprintf("{from: \"%s\", to: \"%s\", value: web3.toWei(0.001, \"ether\")}", account, account)
	txHash := network.IssueCommand(0, fmt.Sprintf("personal.sendTransaction(%s, \"%s\")", tx, password))
	issuedTxStr := network.IssueCommand(0, fmt.Sprintf("eth.getTransaction(%s)", txHash))

	// check the transaction has expected values
	issuedTx := map[string]interface{}{}
	if err = yaml.Unmarshal([]byte(issuedTxStr), issuedTx); err != nil {
		t.Fatalf("unable to unmarshall getTransaction response to YAML. Cause: %s.\nsendTransaction response was: %s\ngetTransaction response was %s", err, txHash, issuedTxStr)
	}

	if issuedTx["value"].(int) != 1000000000000000 ||
		issuedTx["to"].(string) != fmt.Sprintf("0x%s", account) {
		t.Errorf("unable to confirm values")
	}
}

func TestGethTransactionIsMintedOverRPC(t *testing.T) {
	gethBinaryPath, err := EnsureBinariesExist(LatestVersion)
	if err != nil {
		t.Fatal(err)
	}

	// wallet should be prefunded
	w := datagenerator.RandomWallet(genesisChainID)
	startPort := int(integration.StartPortGethNetworkTest) + numNodes*3
	network := NewGethNetwork(startPort, startPort+defaultWsPortOffset, gethBinaryPath, numNodes, 1, []string{w.Address().String()})
	defer network.StopNodes()

	hostConfig := config.HostConfig{
		L1NodeHost:          localhost,
		L1NodeWebsocketPort: network.WebSocketPorts[0],
		L1ConnectionTimeout: defaultL1ConnectionTimeout,
	}
	ethClient, err := ethadapter.NewEthClientFromConfig(hostConfig)
	if err != nil {
		panic(err)
	}

	// pick the first address in the network and send some funds to it
	toAddr := common.HexToAddress(fmt.Sprintf("0x%s", network.addresses[0]))
	tx := &types.LegacyTx{
		Nonce:    w.GetNonceAndIncrement(),
		GasPrice: big.NewInt(20000000000),
		Gas:      uint64(1024_000_000),
		To:       &toAddr,
		Value:    big.NewInt(100),
	}
	signedTx, err := w.SignTransaction(tx)
	if err != nil {
		t.Fatal(err)
	}
	err = ethClient.SendTransaction(signedTx)
	if err != nil {
		t.Fatal(err)
	}

	// make sure it's mined into a block within an acceptable time
	var receipt *types.Receipt
	for start := time.Now(); time.Since(start) < timeout; time.Sleep(time.Second) {
		receipt, err = ethClient.TransactionReceipt(signedTx.Hash())
		if err == nil {
			break
		}
		if !errors.Is(err, ethereum.NotFound) {
			t.Fatal(err)
		}
	}

	if receipt == nil {
		t.Fatalf("Did not mine the transaction after %s seconds - receipt: %+v", timeout, receipt)
	}

	if receipt.BlockNumber == big.NewInt(0) || bytes.Equal(receipt.BlockHash.Bytes(), common.HexToHash("").Bytes()) {
		t.Fatalf("Did not minted/mined the block - receipt: %+v", receipt)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		t.Fatalf("Did not minted/mined the tx correctly - receipt: %+v", receipt)
	}
}

func TestGethGoerliConnection(t *testing.T) {
	gethBinaryPath, err := EnsureBinariesExist(LatestVersion)
	if err != nil {
		t.Fatal(err)
	}

	// wallet should be prefunded
	startPort := int(integration.StartPortGethNetworkTest) + numNodes*3
	network := NewGethNetworkGoerli(startPort, startPort+defaultWsPortOffset, gethBinaryPath, numNodes)

	client, err := ethclient.NewEthClient(config.HostConfig{
		L1NodeHost:          "127.0.0.1",
		L1NodeWebsocketPort: network.WebSocketPorts[0],
		L1ConnectionTimeout: 30 * time.Second,
	})

	if err != nil {
		panic(err)
	}

	timeout := 20 * time.Minute
	var blk *types.Block
	for startTime := time.Now(); time.Since(startTime) < timeout; time.Sleep(time.Second) {
		blk, err = client.BlockByNumber(nil)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Block %s at Number: %d - after %s\n", blk.Hash(), blk.NumberU64(), time.Since(startTime))
	}
	if blk == nil || blk.NumberU64() == 0 {
		t.Fatal("unable to connect to goerli")
	}
}
