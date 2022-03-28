package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/integration/simulation"
)

// DefaultAverageLatencyToBlockRatio is relative to the block time
// Average eth Block duration=12s, and average eth block latency = 1s
// Determines the broadcast powTime. The lower, the more powTime.
const DefaultAverageLatencyToBlockRatio = 12

// DefaultAverageGossipPeriodToBlockRatio - how long to wait for gossip in L2.
const DefaultAverageGossipPeriodToBlockRatio = 3

func main() {
	//f, err := os.Create("cpu.prof")
	//if err != nil {
	//	log.Fatal("could not create CPU profile: ", err)
	//}
	//defer f.Close() // error handling omitted for example
	//if err := pprof.StartCPUProfile(f); err != nil {
	//	log.Fatal("could not start CPU profile: ", err)
	//}
	//defer pprof.StopCPUProfile()
	rand.Seed(time.Now().UnixNano())
	uuid.EnableRandPool()

	f1, err := os.Create("simulation_result.txt")
	if err != nil {
		panic(err)
	}
	defer f1.Close()
	log.SetLog(f1)

	// define core test parameters
	params := simulation.SimParams{
		NumberOfNodes:         10,
		NumberOfWallets:       5,
		AvgBlockDurationUSecs: uint64(25_000),
		SimulationTimeSecs:    15,
	}
	params.AvgNetworkLatency = params.AvgBlockDurationUSecs / DefaultAverageLatencyToBlockRatio
	params.AvgGossipPeriod = params.AvgBlockDurationUSecs / DefaultAverageGossipPeriodToBlockRatio
	params.SimulationTimeUSecs = params.SimulationTimeSecs * 1000 * 1000

	// define network params
	stats := simulation.NewStats(params.NumberOfNodes)

	mockEthNodes, obscuroInMemNodes := simulation.CreateBasicNetworkOfInMemoryNodes(params, stats)

	txInjector := simulation.NewTransactionInjector(params.NumberOfWallets, params.AvgBlockDurationUSecs, stats, params.SimulationTimeUSecs, mockEthNodes, obscuroInMemNodes)

	sim := simulation.Simulation{
		MockEthNodes:       mockEthNodes,      // the list of mock ethereum nodes
		ObscuroNodes:       obscuroInMemNodes, //  the list of in memory obscuro nodes
		AvgBlockDuration:   params.AvgBlockDurationUSecs,
		TxInjector:         txInjector,
		SimulationTimeSecs: params.SimulationTimeSecs,
		Stats:              stats,
	}

	// execute the simulation
	sim.Start()
	fmt.Printf("%s\n", simulation.NewOutputStats(&sim))
	sim.Stop()
}
