package main

import (
	"os"
	"time"

	"github.com/obscuronet/obscuro-playground/go/ethclient/wallet"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host/p2p"

	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"
)

func main() {
	setLogs()
	config := parseCLIArgs()

	nodeID := common.BytesToAddress([]byte(*config.nodeID))
	hostCfg := host.AggregatorCfg{GossipRoundDuration: time.Duration(*config.gossipRoundNanos), ClientRPCTimeoutSecs: *config.rpcTimeoutSecs}
	enclaveClient := host.NewEnclaveRPCClient(*config.enclaveAddr, host.ClientRPCTimeoutSecs*time.Second, nodeID)
	aggP2P := p2p.NewSocketP2PLayer(*config.ourP2PAddr, config.peerP2PAddrs)
	w := wallet.NewInMemoryWallet(*config.privateKeyString)
	contractAddr := common.HexToAddress(*config.contractAddress)
	l1Client, err := ethclient.NewEthClient(nodeID, "127.0.0.1", 7545, w, contractAddr)
	if err != nil {
		panic(err)
	}
	// todo - joel - change clientServer param to non-nil
	agg := host.NewObscuroAggregator(nodeID, hostCfg, nil, *config.isGenesis, aggP2P, l1Client, enclaveClient, nil, ethereum_mock.NewMockTxHandler())

	agg.Start()
}

// Sets the log file.
func setLogs() {
	logFile, err := os.Create("host_logs.txt")
	if err != nil {
		panic(err)
	}
	log.SetLog(logFile)
}
