package transport

import (
	"context"
	"testing"
	"time"

	"github.com/pesarakloo/ani-tunnel/internal/utils"
)

func TestSilentControlPeerIsRetiredEvenWhenWritesSucceed(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	peer, restarts, finished := startLoop(t, ctx, func(c *controlLoop) {
		c.beat = 20 * time.Millisecond
		c.probeRTT = true
		c.probeTimeout = 200 * time.Millisecond
	})
	drain(peer) // kernel accepts every write, but no RTT response returns
	select {
	case <-finished:
	case <-time.After(2 * time.Second):
		t.Fatal("silent control peer held the seat indefinitely")
	}
	eventually(t, time.Second, func() bool { return restarts.Load() == 1 }, "dead peer was not retired")
}

func TestRTTRepliesKeepAnIdleLegacyTCPClientAlive(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	peer, restarts, finished := startLoop(t, ctx, func(c *controlLoop) {
		c.beat = 20 * time.Millisecond
		c.probeRTT = true
		c.probeTimeout = 200 * time.Millisecond
	})
	go func() {
		for {
			b, err := utils.ReceiveBinaryByte(peer)
			if err != nil {
				return
			}
			// Legacy TCP replies to RTT, never to HB.
			if b == utils.SG_RTT {
				_ = utils.SendBinaryByteWithin(peer, b, time.Second)
			}
		}
	}()
	select {
	case <-finished:
		t.Fatal("healthy idle client was retired")
	case <-time.After(700 * time.Millisecond):
	}
	if restarts.Load() != 0 {
		t.Fatal("healthy client requested reconnect")
	}
}
