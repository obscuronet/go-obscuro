package simulation

import (
	"sync/atomic"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// MockP2P - models a full network of in memory nodes including artificial random latencies
// Implements the P2p interface
// Will be plugged into each node
type MockP2P struct {
	currentNode *host.Node
	Nodes       []*host.Node

	avgLatency       uint64
	avgBlockDuration uint64

	listenerInterrupt *int32
}

// NewMockP2P returns an instance of a configured L2 Network (no nodes)
func NewMockP2P(avgBlockDuration uint64, avgLatency uint64) *MockP2P {
	i := int32(0)
	return &MockP2P{
		avgLatency:        avgLatency,
		avgBlockDuration:  avgBlockDuration,
		listenerInterrupt: &i,
	}
}

func (netw *MockP2P) StartListening(host.P2PCallback) {
	// nothing to do here, since communication is direct through the in memory objects
}

func (netw *MockP2P) StopListening() {
	atomic.StoreInt32(netw.listenerInterrupt, 1)
}

// BroadcastRollup Broadcasts the rollup to all L2 peers
func (netw *MockP2P) BroadcastRollup(r obscurocommon.EncodedRollup) {
	if atomic.LoadInt32(netw.listenerInterrupt) == 1 {
		return
	}

	for _, a := range netw.Nodes {
		if a.ID != netw.currentNode.ID {
			t := a
			obscurocommon.Schedule(netw.delay(), func() { t.ReceiveRollup(r) })
		}
	}
}

func (netw *MockP2P) BroadcastTx(tx nodecommon.EncryptedTx) {
	if atomic.LoadInt32(netw.listenerInterrupt) == 1 {
		return
	}

	for _, a := range netw.Nodes {
		if a.ID != netw.currentNode.ID {
			t := a
			obscurocommon.Schedule(netw.delay()/2, func() { t.ReceiveTx(tx) })
		}
	}
}

// delay returns an expected delay on the l2
func (netw *MockP2P) delay() uint64 {
	return obscurocommon.RndBtw(netw.avgLatency/10, 2*netw.avgLatency)
}
