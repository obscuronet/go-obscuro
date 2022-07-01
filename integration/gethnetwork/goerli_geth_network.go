package gethnetwork

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"os"
	"os/exec"
	"path"
	"strconv"
	"sync"
	"time"
)

// NewGethNetworkGoerli returns an Ethereum network with numNodes nodes using the provided Geth binary.
// This networks connects to the Goerli Testnet.
func NewGethNetworkGoerli(portStart int, websocketPortStart int, gethBinaryPath string, numNodes int) *GethNetwork {
	err := ensurePortsAreAvailable(portStart, websocketPortStart, numNodes)
	if err != nil {
		panic(err)
	}

	// Build dirs are suffixed with a timestamp so multiple executions don't collide
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	buildDir := path.Join(basepath, buildDirGoerliBase, timestamp)

	// We create a data directory for each node.
	// Tries to reuse the same data directory per default
	dataDir := path.Join(basepath, buildDirGoerliBase, "data")
	fmt.Printf("Goerli Geth nodes created in: %s\n", dataDir)
	if err != nil {
		panic(err)
	}
	dataDirs := make([]string, numNodes)
	for i := 0; i < numNodes; i++ {
		nodeFolder := nodeFolderName + strconv.Itoa(i+1)
		dataDirs[i] = path.Join(dataDir, nodeFolder)
		err = os.MkdirAll(dataDirs[i], os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	// We push all the node logs to a single file.
	err = os.MkdirAll(buildDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	logFile, err := os.Create(path.Join(buildDir, logFile))
	if err != nil {
		panic(err)
	}

	// We create a password file to unlock the node accounts.
	passwordFile, _ := os.Create(path.Join(dataDir, passwordFile))
	if err != nil {
		panic(err)
	}
	_, err = passwordFile.WriteString(password)
	if err != nil {
		panic(err)
	}

	network := GethNetwork{
		gethBinaryPath:   gethBinaryPath,
		dataDirs:         dataDirs,
		addresses:        make([]string, numNodes),
		nodesProcs:       make([]*os.Process, numNodes),
		logFile:          logFile,
		passwordFilePath: passwordFile.Name(),
		WebSocketPorts:   make([]uint, numNodes),
		commStartPort:    portStart,
		wsStartPort:      websocketPortStart,
	}

	//// We create an account for each node.
	//var wg sync.WaitGroup
	//for idx, dataDir := range dataDirs {
	//	wg.Add(1)
	//	go func(idx int, dataDir string) {
	//		defer wg.Done()
	//		network.createAccount(dataDir)
	//		network.addresses[idx] = network.retrieveAccount(dataDir)
	//	}(idx, dataDir)
	//}
	//wg.Wait()

	// We start the miners.
	var wg sync.WaitGroup
	for idx, dataDir := range dataDirs {
		wg.Add(1)
		go func(idx int, dataDir string) {
			defer wg.Done()
			network.startGoerliMiner(dataDir, idx)
		}(idx, dataDir)
	}
	wg.Wait()

	// start the poller to provide info
	for idx, dataDir := range dataDirs {
		wg.Add(1)
		go func(idx int, dataDir string) {
			defer wg.Done()
			network.startEthPoller(idx)
		}(idx, dataDir)
	}
	wg.Wait()
	//
	//time.Sleep(10 * time.Minute)
	//// We retrieve the enode address for each node.
	//enodeAddrs := make([]string, len(network.dataDirs))
	//for idx, dataDir := range network.dataDirs {
	//	wg.Add(1)
	//	go func(idx int, dataDir string) {
	//		defer wg.Done()
	//		waitForIPC(dataDir) // We cannot issue RPC commands until the IPC files are available.
	//		enodeAddrs[idx] = network.IssueCommand(idx, enodeCmd)
	//	}(idx, dataDir)
	//}
	//wg.Wait()
	//
	//// We manually tell the nodes about one another.
	//for _, enodeAddr := range enodeAddrs {
	//	for idx := range network.dataDirs {
	//		wg.Add(1)
	//		go func(idx int, enodeAddr string) {
	//			defer wg.Done()
	//			// As part of this loop, we also try and peer a node with itself, but Geth ignores this.
	//			network.IssueCommand(idx, fmt.Sprintf(addPeerCmd, enodeAddr))
	//		}(idx, enodeAddr)
	//	}
	//}
	//wg.Wait()

	return &network
}

func (network *GethNetwork) startEthPoller(idx int) {
	client, err := connect("127.0.0.1", network.WebSocketPorts[idx], 2*time.Hour)
	if err != nil {
		panic(err)
	}

	isInSync := true
	for startTime := time.Now(); ; time.Sleep(30 * time.Second) {
		if isInSync {
			syncProgress, err := client.SyncProgress(context.Background())
			if err != nil {
				panic(err)
			}
			if syncProgress != nil {
				fmt.Printf("Node Bootstrapping at height %d of %d after %s\n", syncProgress.CurrentBlock, syncProgress.HighestBlock, time.Since(startTime))
				continue
			}
			isInSync = false
			fmt.Printf("Node Bootstrapping Is finished after %s\n", time.Since(startTime))
			continue
		}

		nodeHeight, err := client.BlockByNumber(context.Background(), nil)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Node at Height %d after %s\n", nodeHeight.NumberU64(), time.Since(startTime))
		continue
	}
}

func connect(ipaddress string, port uint, connectionTimeout time.Duration) (*ethclient.Client, error) {
	var err error
	var c *ethclient.Client
	for start := time.Now(); time.Since(start) < connectionTimeout; time.Sleep(5 * time.Second) {
		c, err = ethclient.Dial(fmt.Sprintf("ws://%s:%d", ipaddress, port))
		if err == nil {
			break
		}
		fmt.Printf("Unable to connect to node after %s\n", time.Since(start))
	}

	return c, err
}

func (network *GethNetwork) createMiner2(dataDir string, idx int) {
	// We delete the leftover IPC file from the previous run, if it exists.
	_ = os.Remove(path.Join(dataDir, ipcFileName))
	// The node must create its initial config based on the network's genesis file before it can be started.
	network.initNode(dataDir)
	webSocketPort := network.wsStartPort + idx
	port := network.commStartPort + idx
	httpPort := network.commStartPort + 25 + idx

	args := []string{
		websocketFlag, wsPortFlag, strconv.Itoa(webSocketPort), dataDirFlag, dataDir, portFlag,
		strconv.Itoa(port), unlockInsecureFlag, unlockFlag, network.addresses[idx], passwordFlag,
		network.passwordFilePath, mineFlag, rpcFeeCapFlag, syncModeFlag,
		httpEnableFlag, httpPortFlag, strconv.Itoa(httpPort), httpEnableApis, allowedAPIs, allowCORSDomain, "*",
		httpIPFlag, "0.0.0.0", gasLimitFlag,
	}
	cmd := exec.Command(network.gethBinaryPath, args...) // nolint

	cmd.Stdout = network.logNodeID(idx)
	cmd.Stderr = network.logNodeID(idx)

	if err := cmd.Start(); err != nil {
		panic(fmt.Errorf("could not start Geth node. Cause: %w", err))
	}
	network.nodesProcs[idx] = cmd.Process
	network.WebSocketPorts[idx] = uint(webSocketPort)
}

// Starts a Geth miner in the goerli network.
func (network *GethNetwork) startGoerliMiner(dataDirPath string, idx int) {
	webSocketPort := network.wsStartPort + idx
	port := network.commStartPort + idx
	httpPort := network.commStartPort + 25 + idx

	args := []string{
		websocketFlag, wsPortFlag, strconv.Itoa(webSocketPort), dataDirFlag, dataDirPath, portFlag,
		strconv.Itoa(port),
		unlockInsecureFlag,
		//unlockFlag, network.addresses[idx], passwordFlag, network.passwordFilePath,
		syncModeFlag,
		httpEnableFlag, httpPortFlag, strconv.Itoa(httpPort), httpEnableApis, allowedAPIs, allowCORSDomain, "*",
		httpIPFlag, "127.0.0.1", "--goerli", "--ipcdisable",
	}
	cmd := exec.Command(network.gethBinaryPath, args...) // nolint

	cmd.Stdout = network.logNodeID(idx)
	cmd.Stderr = network.logNodeID(idx)

	if err := cmd.Start(); err != nil {
		panic(fmt.Errorf("could not start Geth node. Cause: %w", err))
	}
	network.nodesProcs[idx] = cmd.Process
	network.WebSocketPorts[idx] = uint(webSocketPort)
}
