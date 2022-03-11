package simulation

import (
	"sync/atomic"
	"time"

	"github.com/obscuronet/obscuro-playground/go/common"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"
)

// L1NetworkCfg - models a full network including artificial random latencies
type L1NetworkCfg struct {
	nodes []*ethereum_mock.Node
	delay common.Latency // the latency
	Stats *Stats
	// used as a signal to stop all network communication.
	// This helps prevent deadlocks when stopping nodes
	interrupt *int32
}

// BroadcastBlock broadcast a block to the l1 nodes
func (n *L1NetworkCfg) BroadcastBlock(b common.EncodedBlock, p common.EncodedBlock) {
	if atomic.LoadInt32(n.interrupt) == 1 {
		return
	}

	bl, _ := b.Decode()
	for _, m := range n.nodes {
		if m.ID != bl.Header().Coinbase {
			t := m
			common.Schedule(n.delay(), func() { t.P2PReceiveBlock(b, p) })
		} else {
			common.Log(printBlock(bl, *m))
		}
	}

	n.Stats.NewBlock(bl)
}

// BroadcastTx Broadcasts the L1 tx containing the rollup to the L1 network
func (n *L1NetworkCfg) BroadcastTx(tx common.EncodedL1Tx) {
	if atomic.LoadInt32(n.interrupt) == 1 {
		return
	}

	for _, m := range n.nodes {
		t := m
		// the time to broadcast a tx is half that of a L1 block, because it is smaller.
		// todo - find a better way to express this
		d := common.Max(n.delay()/2, 1)
		common.Schedule(d, func() { t.P2PGossipTx(tx) })
	}
}

func (n *L1NetworkCfg) Start(delay time.Duration) {
	// Start l1 nodes
	for _, m := range n.nodes {
		t := m
		go t.Start()
		// don't start everything at once
		time.Sleep(delay)
	}
}

func (n *L1NetworkCfg) Stop() {
	atomic.StoreInt32(n.interrupt, 1)
	for _, m := range n.nodes {
		t := m
		go t.Stop()
		// fmt.Printf("Stopped L1 node: %d.\n", m.ID)
	}
}
