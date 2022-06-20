package gethnetwork

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/obscuronet/obscuro-playground/go/log"
)

const (
	nodeFolderName = "node_datadir_"
	buildDirBase   = "../.build/geth"
	keystoreDir    = "keystore"

	genesisFileName = "genesis.json"
	ipcFileName     = "geth.ipc"
	logFile         = "node_logs.txt"
	passwordFile    = "password.txt"
	password        = "password"

	accountCmd    = "account"
	accountNewCmd = "new"
	addPeerCmd    = "admin.addPeer(%s)"
	attachCmd     = "attach"
	enodeCmd      = "admin.nodeInfo.enode"
	initCmd       = "init"

	dataDirFlag        = "--datadir"
	execFlag           = "--exec"
	mineFlag           = "--mine"
	passwordFlag       = "--password"
	portFlag           = "--port"
	httpEnableFlag     = "--http"
	httpPortFlag       = "--http.port"
	httpIPFlag         = "--http.addr"
	httpEnableApis     = "--http.api"
	allowedAPIs        = "personal,eth,net,web3,debug"
	allowCORSDomain    = "--http.corsdomain"
	rpcFeeCapFlag      = "--rpc.txfeecap=0" // Disables the 1 ETH cap for RPC transactions.
	unlockFlag         = "--unlock"
	unlockInsecureFlag = "--allow-insecure-unlock"
	websocketFlag      = "--ws" // Enables websocket connections to the node.
	wsPortFlag         = "--ws.port"
	gasLimitFlag       = "--miner.gaslimit=2000000000" // Ensures the miners don't gradually reduce the block gas limit.

	// syncModeFlag defines the node block sync approach
	// snap (the default) mode does not work well for small, rapidly deployed private networks
	syncModeFlag = "--syncmode=full"

	// We pre-allocate a wallet matching the private key used in the tests, plus an account per clique member.
	genesisJSONTemplate = `{
	  "config": {
		"chainId": 1337,
		"homesteadBlock": 0,
		"eip150Block": 0,
		"eip155Block": 0,
		"eip158Block": 0,
		"byzantiumBlock": 0,
		"constantinopleBlock": 0,
		"petersburgBlock": 0,
		"istanbulBlock": 0,
		"berlinBlock": 0,
		"londonBlock": 0,
		"clique": {
		  "period": %d,
		  "epoch": 30000
		}
	  },
	  "alloc": {
		"0x323AefbFC16159655514846a9e5433C457de9389": {
		  "balance": "1000000000000000000000"
		},
%s
	  },
	  "coinbase": "0x0000000000000000000000000000000000000000",
	  "difficulty": "0x20000",
	  "extraData": "0x0000000000000000000000000000000000000000000000000000000000000000%s0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	  "gasLimit": "2000000000",
	  "nonce": "0x0000000000000042",
	  "mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	  "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	  "timestamp": "0x00"
  }`
	addrBlockTemplate = `		"%s": {
		  "balance": "1000000000000000000000"
		}`
	genesisJSONAddrKey = "address"
)

// GethNetwork is a network of Geth nodes, built using the provided Geth binary.
type GethNetwork struct {
	GenesisJSON []byte // The genesis JSON config used by the network.

	gethBinaryPath   string
	genesisFilePath  string
	dataDirs         []string
	addresses        []string      // The public keys of the nodes' accounts.
	nodesProcs       []*os.Process // The running Geth node processes.
	logFile          *os.File
	passwordFilePath string // The path to the file storing the password to unlock node accounts.
	WebSocketPorts   []uint // Ports exposed by the geth nodes for
	commStartPort    int
	wsStartPort      int
}

// NewGethNetwork returns an Ethereum network with numNodes nodes using the provided Geth binary and allows for prefunding addresses.
// The network uses the Clique consensus algorithm, producing a block every blockTimeSecs.
// A portStart is required for running multiple networks in the same host ( specially useful for unit tests )
func NewGethNetwork(portStart int, websocketPortStart int, gethBinaryPath string, numNodes int, blockTimeSecs int, preFundedAddrs []string) *GethNetwork {
	err := ensurePortsAreAvailable(portStart, websocketPortStart, numNodes)
	if err != nil {
		panic(err)
	}

	// Build dirs are suffixed with a timestamp so multiple executions don't collide
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	buildDir := path.Join(basepath, buildDirBase, timestamp)
	// We create a data directory for each node.
	nodesDir, err := ioutil.TempDir("", timestamp)
	fmt.Printf("Geth nodes created in: %s\n", nodesDir)
	fmt.Printf("Geth nodes ports: Websocket: %d - %d HTTP: %d - %d\n",
		websocketPortStart,
		websocketPortStart+numNodes,
		portStart+25,
		portStart+25+numNodes,
	)
	if err != nil {
		panic(err)
	}
	dataDirs := make([]string, numNodes)
	for i := 0; i < numNodes; i++ {
		nodeFolder := nodeFolderName + strconv.Itoa(i+1)
		dataDirs[i] = path.Join(nodesDir, nodeFolder)
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
	passwordFile, _ := os.Create(path.Join(nodesDir, passwordFile))
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

	// We create an account for each node.
	var wg sync.WaitGroup
	for idx, dataDir := range dataDirs {
		wg.Add(1)
		go func(idx int, dataDir string) {
			defer wg.Done()
			network.createAccount(dataDir)
			network.addresses[idx] = network.retrieveAccount(dataDir)
		}(idx, dataDir)
	}
	wg.Wait()

	// We generate the genesis config file based on the accounts above and the prefunded addresses.
	allocs := make([]string, numNodes+len(preFundedAddrs))
	for i, addr := range append(network.addresses, preFundedAddrs...) {
		allocs[i] = fmt.Sprintf(addrBlockTemplate, addr)
	}
	network.GenesisJSON = []byte(
		fmt.Sprintf(genesisJSONTemplate, blockTimeSecs, strings.Join(allocs, ",\r\n"), strings.Join(network.addresses, "")),
	)

	// We write out the `genesis.json` file to be used by the network.
	genesisFilePath := path.Join(buildDir, genesisFileName)
	err = os.WriteFile(genesisFilePath, network.GenesisJSON, 0o600)
	if err != nil {
		panic(err)
	}
	network.genesisFilePath = genesisFilePath

	// We start the miners.
	for idx, dataDir := range dataDirs {
		wg.Add(1)
		go func(idx int, dataDir string) {
			defer wg.Done()
			network.createMiner(dataDir, idx)
		}(idx, dataDir)
	}
	wg.Wait()

	// We retrieve the enode address for each node.
	enodeAddrs := make([]string, len(network.dataDirs))
	for idx, dataDir := range network.dataDirs {
		wg.Add(1)
		go func(idx int, dataDir string) {
			defer wg.Done()
			waitForIPC(dataDir) // We cannot issue RPC commands until the IPC files are available.
			enodeAddrs[idx] = network.IssueCommand(idx, enodeCmd)
		}(idx, dataDir)
	}
	wg.Wait()

	// We manually tell the nodes about one another.
	for _, enodeAddr := range enodeAddrs {
		for idx := range network.dataDirs {
			wg.Add(1)
			go func(idx int, enodeAddr string) {
				defer wg.Done()
				// As part of this loop, we also try and peer a node with itself, but Geth ignores this.
				network.IssueCommand(idx, fmt.Sprintf(addPeerCmd, enodeAddr))
			}(idx, enodeAddr)
		}
	}
	wg.Wait()

	return &network
}

// IssueCommand sends the command via RPC to the nodeIdx'th node in the network.
func (network *GethNetwork) IssueCommand(nodeIdx int, command string) string {
	dataDir := network.dataDirs[nodeIdx]

	args := []string{dataDirFlag, dataDir, attachCmd, path.Join(dataDir, ipcFileName), execFlag, command}
	cmd := exec.Command(network.gethBinaryPath, args...) // nolint
	cmd.Stderr = network.logNodeID(nodeIdx)

	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(string(output))
}

// StopNodes kills the Geth node processes.
func (network *GethNetwork) StopNodes() {
	var wg sync.WaitGroup
	for _, process := range network.nodesProcs {
		if process != nil {
			wg.Add(1)
			go func(process *os.Process) {
				defer wg.Done()
				err := process.Kill()
				if err != nil {
					log.Error("geth node could not be killed: %s", err)
				}
				_, err = process.Wait()
				if err != nil {
					log.Error("geth node was killed successfully but did not exit: %s", err)
				}
			}(process)
		}
	}
	wg.Wait()
}

// Initialises and starts a miner.
func (network *GethNetwork) createMiner(dataDir string, idx int) {
	// We delete the leftover IPC file from the previous run, if it exists.
	_ = os.Remove(path.Join(dataDir, ipcFileName))
	// The node must create its initial config based on the network's genesis file before it can be started.
	network.initNode(dataDir)
	network.startMiner(dataDir, idx)
}

// Creates an account for a Geth node.
func (network *GethNetwork) createAccount(dataDirPath string) {
	args := []string{dataDirFlag, dataDirPath, accountCmd, accountNewCmd, passwordFlag, network.passwordFilePath}
	cmd := exec.Command(network.gethBinaryPath, args...) // nolint
	cmd.Stdout = network.logFile
	cmd.Stderr = network.logFile

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

// Adds a Geth node's account public key to the `network` object.
func (network *GethNetwork) retrieveAccount(dataDirPath string) string {
	dir := path.Join(dataDirPath, keystoreDir)
	files, _ := ioutil.ReadDir(dir)
	for _, file := range files {
		// `ReadDir` returns the folder itself, as well as the files within it.
		if file.IsDir() {
			continue
		}

		y, err := os.ReadFile(path.Join(dir, file.Name()))
		if err != nil {
			panic(err)
		}
		contents := make(map[string]interface{})
		err = json.Unmarshal(y, &contents)
		if err != nil {
			panic(err)
		}
		return contents[genesisJSONAddrKey].(string) // We return as we only expect one account per node.
	}

	panic(fmt.Sprintf("could not find account for node at %s", dataDirPath))
}

// Initialises a Geth node based on the network genesis file.
func (network *GethNetwork) initNode(dataDirPath string) {
	args := []string{dataDirFlag, dataDirPath, initCmd, network.genesisFilePath}
	cmd := exec.Command(network.gethBinaryPath, args...) // nolint
	cmd.Stdout = network.logFile
	cmd.Stderr = network.logFile

	if err := cmd.Run(); err != nil {
		panic(fmt.Errorf("could not initialise Geth node. Cause: %w", err))
	}
}

// Starts a Geth miner.
func (network *GethNetwork) startMiner(dataDirPath string, idx int) {
	webSocketPort := network.wsStartPort + idx
	port := network.commStartPort + idx
	httpPort := network.commStartPort + 25 + idx

	args := []string{
		websocketFlag, wsPortFlag, strconv.Itoa(webSocketPort), dataDirFlag, dataDirPath, portFlag,
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

// logNodeID prepends the nodeID to the log entries
func (network *GethNetwork) logNodeID(idx int) io.Writer {
	r, w, _ := os.Pipe()
	go func() {
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			_, _ = network.logFile.WriteString(fmt.Sprintf("EthNode-%d: %s\n", idx, sc.Text()))
		}
	}()
	return w
}

// Waits for a node's IPC file to exist.
func waitForIPC(dataDir string) {
	totalCounter := 0
	counter := 0

	for {
		ipcFilePath := path.Join(dataDir, ipcFileName)
		_, err := os.Stat(ipcFilePath)
		if err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)

		if totalCounter > 300 {
			panic(fmt.Errorf("waited over 30 seconds for .ipc file of node at %s", dataDir))
		}

		if counter > 20 {
			log.Info(fmt.Sprintf("Waiting for .ipc file of node at %s", dataDir))
			totalCounter += counter
			counter = 0
		}

		counter++
	}
}

func ensurePortsAreAvailable(startPort int, websocketStartPort int, numberNodes int) error {
	var unavailablePorts []int

	for i := 0; i < numberNodes; i++ {
		commsPort := startPort + i
		if !isPortAvailable(commsPort) {
			unavailablePorts = append(unavailablePorts, commsPort)
		}
		wsPort := websocketStartPort + i
		if !isPortAvailable(wsPort) {
			unavailablePorts = append(unavailablePorts, wsPort)
		}
	}

	if len(unavailablePorts) > 0 {
		list, _ := json.Marshal(unavailablePorts)
		return fmt.Errorf("could not run geth network because test ports are unavailable for use - the following ports were unavailable: %s", list)
	}
	return nil
}

func isPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("Listen port %d. Err: %s. ", port, err)
	}
	if ln != nil {
		_ = ln.Close()
	}
	return err == nil
}
