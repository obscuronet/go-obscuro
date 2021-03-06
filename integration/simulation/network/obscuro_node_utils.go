package network

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/obscuronet/go-obscuro/integration/common/viewkey"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/enclave"
	"github.com/obscuronet/go-obscuro/integration"

	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/host"
	"github.com/obscuronet/go-obscuro/go/rpcclientlib"
	"github.com/obscuronet/go-obscuro/integration/simulation/p2p"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
	"github.com/obscuronet/go-obscuro/integration/simulation/stats"
)

func startInMemoryObscuroNodes(params *params.SimParams, stats *stats.Stats, genesisJSON []byte, l1Clients []ethadapter.EthClient) ([]rpcclientlib.Client, map[string]rpcclientlib.Client) {
	// Create the in memory obscuro nodes, each connect each to a geth node
	obscuroNodes := make([]*host.Node, params.NumberOfNodes)
	for i := 0; i < params.NumberOfNodes; i++ {
		isGenesis := i == 0
		obscuroNodes[i] = createInMemObscuroNode(
			int64(i),
			isGenesis,
			params.MgmtContractLib,
			params.ERC20ContractLib,
			params.AvgGossipPeriod,
			params.AvgBlockDuration,
			params.AvgNetworkLatency,
			stats,
			true,
			genesisJSON,
			params.Wallets.NodeWallets[i],
			l1Clients[i],
			params.ViewingKeysEnabled,
			params.Wallets,
		)
	}
	// make sure the aggregators can talk to each other
	for _, m := range obscuroNodes {
		mockP2P := m.P2p.(*p2p.MockP2P)
		mockP2P.Nodes = obscuroNodes
	}

	// start each obscuro node
	for _, m := range obscuroNodes {
		t := m
		go t.Start()
	}

	// Create a handle to each node
	obscuroClients := make([]rpcclientlib.Client, params.NumberOfNodes)
	for i, node := range obscuroNodes {
		obscuroClients[i] = host.NewInMemObscuroClient(node)
	}
	time.Sleep(100 * time.Millisecond)

	walletClients := setupWalletClientsWithoutViewingKeys(params, obscuroClients)

	return obscuroClients, walletClients
}

// setupWalletClientsWithoutViewingKeys will configure the existing obscuro clients to be used as wallet clients,
// (we typically keep a client per wallet so their viewing keys are available and registered but they can share clients
// 	if no viewing keys are in use)
func setupWalletClientsWithoutViewingKeys(params *params.SimParams, obscuroClients []rpcclientlib.Client) map[string]rpcclientlib.Client {
	walletClients := make(map[string]rpcclientlib.Client)
	var i int
	// loop through all the L2 wallets we're using and round-robin allocate them the rpc clients we have for each host
	for _, w := range params.Wallets.SimObsWallets {
		walletClients[w.Address().String()] = obscuroClients[i%len(obscuroClients)]
		i++
	}
	for _, t := range params.Wallets.Tokens {
		w := t.L2Owner
		walletClients[w.Address().String()] = obscuroClients[i%len(obscuroClients)]
		i++
	}
	return walletClients
}

// todo: this method is quite heavy, should refactor to separate out the creation of the nodes, starting of the nodes, setup of the RPC clients etc.
func startStandaloneObscuroNodes(params *params.SimParams, stats *stats.Stats, gethClients []ethadapter.EthClient, enclaveAddresses []string) ([]rpcclientlib.Client, map[string]rpcclientlib.Client) {
	// handle to the obscuro clients
	nodeRPCAddresses := make([]string, params.NumberOfNodes)
	obscuroClients := make([]rpcclientlib.Client, params.NumberOfNodes)
	obscuroNodes := make([]*host.Node, params.NumberOfNodes)

	for i := 0; i < params.NumberOfNodes; i++ {
		isGenesis := i == 0

		// We use the convention to determine the rpc ports of the node
		nodeRPCPortHTTP := params.StartPort + DefaultHostRPCHTTPOffset + i
		nodeRPCPortWS := params.StartPort + DefaultHostRPCWSOffset + i

		// create a remote enclave server
		obscuroNodes[i] = createSocketObscuroNode(
			int64(i),
			isGenesis,
			params.AvgGossipPeriod,
			stats,
			fmt.Sprintf("%s:%d", Localhost, params.StartPort+DefaultHostP2pOffset+i),
			enclaveAddresses[i],
			Localhost,
			uint64(nodeRPCPortHTTP),
			uint64(nodeRPCPortWS),
			params.Wallets.NodeWallets[i],
			params.MgmtContractLib,
			gethClients[i],
		)

		nodeRPCAddresses[i] = fmt.Sprintf("%s:%d", Localhost, nodeRPCPortHTTP)
		obscuroClients[i] = rpcclientlib.NewClient(nodeRPCAddresses[i])
	}

	// start each obscuro node
	for _, m := range obscuroNodes {
		t := m
		go t.Start()
		time.Sleep(params.AvgBlockDuration / 3)
	}

	// wait for the clients to be connected
	for i, client := range obscuroClients {
		started := false
		for !started {
			err := client.Call(nil, rpcclientlib.RPCGetID)
			started = err == nil
			if !started {
				fmt.Printf("Could not connect to client %d. Err %s. Retrying..\n", i, err)
			}
			time.Sleep(500 * time.Millisecond)
		}
	}

	// round-robin the wallets onto the different obscuro nodes, register them each a viewing key
	walletClients := make(map[string]rpcclientlib.Client)
	var i int
	for _, w := range params.Wallets.SimObsWallets {
		walletClients[w.Address().String()] = rpcclientlib.NewClient(nodeRPCAddresses[i%len(nodeRPCAddresses)])
		viewkey.GenerateAndRegisterViewingKey(walletClients[w.Address().String()], w)
		i++
	}
	for _, t := range params.Wallets.Tokens {
		w := t.L2Owner
		walletClients[w.Address().String()] = rpcclientlib.NewClient(nodeRPCAddresses[i%len(nodeRPCAddresses)])
		viewkey.GenerateAndRegisterViewingKey(walletClients[w.Address().String()], w)
		i++
	}

	return obscuroClients, walletClients
}

func startRemoteEnclaveServers(startAt int, params *params.SimParams, stats *stats.Stats) {
	for i := startAt; i < params.NumberOfNodes; i++ {
		// create a remote enclave server
		enclaveAddr := fmt.Sprintf("%s:%d", Localhost, params.StartPort+DefaultEnclaveOffset+i)
		hostAddr := fmt.Sprintf("%s:%d", Localhost, params.StartPort+DefaultHostP2pOffset+i)
		enclaveConfig := config.EnclaveConfig{
			HostID:                 common.BigToAddress(big.NewInt(int64(i))),
			HostAddress:            hostAddr,
			Address:                enclaveAddr,
			L1ChainID:              integration.EthereumChainID,
			ObscuroChainID:         integration.ObscuroChainID,
			ValidateL1Blocks:       false,
			WillAttest:             false,
			GenesisJSON:            nil,
			UseInMemoryDB:          false,
			ERC20ContractAddresses: params.Wallets.AllEthAddresses(),
			ViewingKeysEnabled:     params.ViewingKeysEnabled,
		}
		_, err := enclave.StartServer(enclaveConfig, params.MgmtContractLib, params.ERC20ContractLib, stats)
		if err != nil {
			panic(fmt.Sprintf("failed to create enclave server: %v", err))
		}
	}
}

func StopObscuroNodes(clients []rpcclientlib.Client) {
	var wg sync.WaitGroup
	for _, client := range clients {
		wg.Add(1)
		go func(c rpcclientlib.Client) {
			defer wg.Done()
			err := c.Call(nil, rpcclientlib.RPCStopHost)
			if err != nil {
				log.Error("Failed to stop client %s", err)
			}
			c.Stop()
		}(client)
	}
	if waitTimeout(&wg, 2*time.Second) {
		log.Error("Timed out waiting for the obscuro nodes to stop")
	} else {
		log.Info("Obscuro nodes stopped")
	}
	// Wait a bit for the nodes to shut down.
	time.Sleep(2 * time.Second)
}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
