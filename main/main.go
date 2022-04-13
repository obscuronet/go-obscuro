package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/beacon"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
)

func main() {
	db := createDB()
	cacheConfig := createCacheConfig()
	chainConfig, _, _ := createChainConfig(db)
	engine := createEngine()
	vmConfig := createVMConfig()
	shouldPreserve := createShouldPreserve()
	txLookupLimit := uint64(2_350_000) // Default.

	blockchain, err := core.NewBlockChain(db, cacheConfig, chainConfig, engine, vmConfig, shouldPreserve, &txLookupLimit)
	if err != nil {
		panic(err)
	}

	// We print the genesis block hash.
	println(blockchain.Genesis().Hash().String())
}

// Non-ephemeral nodes use `rawdb.NewLevelDBDatabaseWithFreezer` instead.
func createDB() ethdb.Database {
	return rawdb.NewMemoryDatabase()
}

func createCacheConfig() *core.CacheConfig {
	return &core.CacheConfig{
		TrieCleanLimit:      614,           // Default.
		TrieCleanJournal:    "",            // Defaults to `geth/triecache` in the node's data directory.
		TrieCleanRejournal:  3600000000000, // Default.
		TrieCleanNoPrefetch: false,         // Default.
		TrieDirtyLimit:      1024,          // Default.
		TrieDirtyDisabled:   false,         // Default.
		TrieTimeLimit:       3600000000000, // Default.
		SnapshotLimit:       409,           // Default.
		Preimages:           false,         // Default.
	}
}

func createChainConfig(db ethdb.Database) (*params.ChainConfig, common.Hash, error) {
	return core.SetupGenesisBlockWithOverride(
		db,
		nil, // Default.
		nil, // Default.
		nil, // Default.
	)
}

// Recreates the standard path through `eth/ethconfig/config.go/CreateConsensusEngine()`.
func createEngine() consensus.Engine {
	var engine consensus.Engine //nolint
	engine = ethash.New(ethash.Config{
		PowMode:          ethash.ModeNormal, // Default.
		CacheDir:         "",                // Defaults to `geth/ethash` in the node's data directory.
		CachesInMem:      2,                 // Default.
		CachesOnDisk:     3,                 // Default.
		CachesLockMmap:   false,             // Default.
		DatasetDir:       "",                // Defaults to `~/Library/Ethash` in the node's data directory.
		DatasetsInMem:    1,                 // Default.
		DatasetsOnDisk:   2,                 // Default.
		DatasetsLockMmap: false,             // Default.
		NotifyFull:       false,             // Default.
	}, nil, false) // Defaults.
	engine.(*ethash.Ethash).SetThreads(-1)
	return beacon.New(engine)
}

func createVMConfig() vm.Config {
	return vm.Config{
		EnablePreimageRecording: false, // Default.
	}
}

// We indicate that no blocks are authored by local accounts, and thus all blocks are discarded during reorgs.
func createShouldPreserve() func(header *types.Header) bool {
	return func(header *types.Header) bool {
		return false
	}
}
