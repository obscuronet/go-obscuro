package network

import (
	"math/big"
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/config"

	"github.com/obscuronet/obscuro-playground/go/ethclient/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"github.com/obscuronet/obscuro-playground/integration"
	simp2p "github.com/obscuronet/obscuro-playground/integration/simulation/p2p"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"

	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"

	"github.com/obscuronet/obscuro-playground/integration/simulation/stats"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host/p2p"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"
)

const (
	Localhost                  = "127.0.0.1"
	DefaultWsPortOffset        = 100 // The default offset between a Geth node's port and websocket ports.
	DefaultHostP2pOffset       = 200 //  The default offset for the host P2p
	DefaultHostRPCOffset       = 400 //  The default offset for the host RPC
	DefaultEnclaveOffset       = 300 //  The default offset between a Geth nodes port and the enclave ports. Used in Socket Simulations.
	ClientRPCTimeoutSecs       = 5
	DefaultL1ConnectionTimeout = 15 * time.Second
)

func createMockEthNode(id int64, nrNodes int, avgBlockDuration time.Duration, avgNetworkLatency time.Duration, stats *stats.Stats) *ethereum_mock.Node {
	mockEthNetwork := ethereum_mock.NewMockEthNetwork(avgBlockDuration, avgNetworkLatency, stats)
	ethereumMockCfg := defaultMockEthNodeCfg(nrNodes, avgBlockDuration)
	// create an in memory mock ethereum node responsible with notifying the layer 2 node about blocks
	miner := ethereum_mock.NewMiner(common.BigToAddress(big.NewInt(id)), ethereumMockCfg, mockEthNetwork, stats)
	mockEthNetwork.CurrentNode = miner
	return miner
}

func createInMemObscuroNode(
	id int64,
	isGenesis bool,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	stableTokenContractLib erc20contractlib.ERC20ContractLib,
	avgGossipPeriod time.Duration,
	avgBlockDuration time.Duration,
	avgNetworkLatency time.Duration,
	stats *stats.Stats,
	validateBlocks bool,
	genesisJSON []byte,
	ethWallet wallet.Wallet,
	mgmtContractBlkHash *common.Hash,
) *host.Node {
	obscuroInMemNetwork := simp2p.NewMockP2P(avgBlockDuration, avgNetworkLatency)

	hostConfig := config.HostConfig{
		ID:                  common.BigToAddress(big.NewInt(id)),
		IsGenesis:           isGenesis,
		GossipRoundDuration: avgGossipPeriod,
		HasClientRPC:        false,
		ContractMgmtBlkHash: mgmtContractBlkHash,
	}

	enclaveConfig := config.EnclaveConfig{
		HostID:           hostConfig.ID,
		L1ChainID:        integration.EthereumChainID,
		ObscuroChainID:   integration.ObscuroChainID,
		WillAttest:       false,
		ValidateL1Blocks: validateBlocks,
		GenesisJSON:      genesisJSON,
		UseInMemoryDB:    true,
	}
	enclaveClient := enclave.NewEnclave(enclaveConfig, mgmtContractLib, stableTokenContractLib, stats)

	// create an in memory obscuro node
	node := host.NewHost(
		hostConfig,
		stats,
		obscuroInMemNetwork,
		nil,
		enclaveClient,
		ethWallet,
		mgmtContractLib,
	)
	obscuroInMemNetwork.CurrentNode = node
	return node
}

func createSocketObscuroNode(
	id int64,
	isGenesis bool,
	avgGossipPeriod time.Duration,
	stats *stats.Stats,
	p2pAddr string,
	peerAddrs []string,
	enclaveAddr string,
	clientServerAddr string,
	ethWallet wallet.Wallet,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	mgmtContractTxHash *common.Hash,
) *host.Node {
	hostConfig := config.HostConfig{
		ID:                  common.BigToAddress(big.NewInt(id)),
		IsGenesis:           isGenesis,
		GossipRoundDuration: avgGossipPeriod,
		HasClientRPC:        true,
		ClientRPCAddress:    clientServerAddr,
		ClientRPCTimeout:    ClientRPCTimeoutSecs * time.Second,
		EnclaveRPCTimeout:   ClientRPCTimeoutSecs * time.Second,
		EnclaveRPCAddress:   enclaveAddr,
		P2PAddress:          p2pAddr,
		AllP2PAddresses:     peerAddrs,
		ContractMgmtBlkHash: mgmtContractTxHash,
	}

	// create an enclave client
	enclaveClient := host.NewEnclaveRPCClient(hostConfig)

	// create a socket obscuro node
	nodeP2p := p2p.NewSocketP2PLayer(hostConfig)

	node := host.NewHost(
		hostConfig,
		stats,
		nodeP2p,
		nil,
		enclaveClient,
		ethWallet,
		mgmtContractLib,
	)

	return node
}

func defaultMockEthNodeCfg(nrNodes int, avgBlockDuration time.Duration) ethereum_mock.MiningConfig {
	return ethereum_mock.MiningConfig{
		PowTime: func() time.Duration {
			// This formula might feel counter-intuitive, but it is a good approximation for Proof of Work.
			// It creates a uniform distribution up to nrMiners*avgDuration
			// Which means on average, every round, the winner (miner who gets the lowest nonce) will pick a number around "avgDuration"
			// while everyone else will have higher values.
			// Over a large number of rounds, the actual average block duration will be around the desired value, while the number of miners who get very close numbers will be limited.
			return obscurocommon.RndBtwTime(avgBlockDuration/time.Duration(nrNodes), avgBlockDuration*time.Duration(nrNodes))
		},
	}
}
