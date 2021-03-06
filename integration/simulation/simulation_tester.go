package simulation

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/integration/simulation/network"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"

	stats2 "github.com/obscuronet/go-obscuro/integration/simulation/stats"

	"github.com/google/uuid"
)

// testSimulation encapsulates the shared logic for simulating and testing various types of nodes.
func testSimulation(t *testing.T, netw network.Network, params *params.SimParams) {
	defer func() {
		// wait until clean up is complete before we log the lingering goroutine count
		log.Info("goroutine leak monitor - simulation end - %d goroutines currently running", runtime.NumGoroutine())
	}()
	log.Info("goroutine leak monitor - simulation start - %d goroutines currently running", runtime.NumGoroutine())
	rand.Seed(time.Now().UnixNano())
	uuid.EnableRandPool()

	stats := stats2.NewStats(params.NumberOfNodes) // todo - temporary object used to collect metrics. Needs to be replaced with something better

	defer netw.TearDown()
	networkClients, err := netw.Create(params, stats)
	// Return early if the network was not created
	if err != nil {
		fmt.Printf("Could not run test: %s\n", err)
		return
	}

	txInjector := NewTransactionInjector(
		params.AvgBlockDuration,
		stats,
		networkClients,
		params.Wallets,
		params.MgmtContractAddr,
		params.MgmtContractLib,
		params.ERC20ContractLib,
		0,
	)

	simulation := Simulation{
		RPCHandles:       networkClients,
		AvgBlockDuration: uint64(params.AvgBlockDuration),
		TxInjector:       txInjector,
		SimulationTime:   params.SimulationTime,
		Stats:            stats,
		Params:           params,
	}

	// execute the simulation
	simulation.Start()

	// run tests
	checkNetworkValidity(t, &simulation)

	simulation.Stop()

	// generate and print the final stats
	t.Logf("Simulation results:%+v", NewOutputStats(&simulation))
}
