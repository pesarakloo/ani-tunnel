package transport

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pesarakloo/ani-tunnel/internal/utils"
)

// Losing only the control socket must retire idle data sockets and return
// their quota while the generation and its listeners remain available.
func TestLostTCPClientReleasesActiveDataQuota(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := NewTCPServer(ctx, &TcpConfig{ChannelSize: 4, Heartbeat: time.Hour, MaxConnections: 1}, quietLogger())
	g := s.newGen(ctx)
	control, peer := pipePair(t)
	drain(peer)
	s.seatClient(g, controlCandidate{conn: control, nonce: "first"})
	go s.handleLoop(g)
	data, backend := pipePair(t)
	s.deliverTunnelConn(g, data, "first")
	local, user := pipePair(t)
	if !s.limits.acquire() {
		t.Fatal("first connection refused")
	}
	g.localChannel <- LocalTCPConn{conn: local, remoteAddr: "127.0.0.1:443", timeCreated: nowMillis()}
	_ = backend.SetReadDeadline(time.Now().Add(3 * time.Second))
	if _, _, err := utils.ReceiveBinaryTransportString(backend); err != nil {
		t.Fatal(err)
	}
	_ = backend.SetReadDeadline(time.Time{})
	// The data peers deliberately stay open and idle, as on a blackholed path.
	peer.Close()
	eventually(t, 2*time.Second, func() bool { return !s.controlChannel.IsSet() }, "control channel did not retire")
	eventually(t, 2*time.Second, func() bool { return s.limits.active.Load() == 0 }, "lost client leaked an active data connection and its quota")
	if ctx.Err() != nil {
		t.Fatal("client loss ended the listening generation")
	}
	_ = user.SetReadDeadline(time.Now().Add(time.Second))
	if _, err := user.Read(make([]byte, 1)); err == nil {
		t.Fatal("idle user socket stayed open")
	} else if e, ok := err.(interface{ Timeout() bool }); ok && e.Timeout() {
		t.Fatal("user socket only timed out; it was not closed")
	}
	// Reconnect and transfer through the very same generation and limiter.
	control2, peer2 := pipePair(t)
	drain(peer2)
	s.seatClient(g, controlCandidate{conn: control2, nonce: "second"})
	data2, backend2 := pipePair(t)
	s.deliverTunnelConn(g, data2, "second")
	local2, user2 := pipePair(t)
	if !s.limits.acquire() {
		t.Fatal("reconnected client has no free quota")
	}
	g.localChannel <- LocalTCPConn{conn: local2, remoteAddr: "127.0.0.1:443", timeCreated: nowMillis()}
	_ = backend2.SetReadDeadline(time.Now().Add(3 * time.Second))
	if _, _, err := utils.ReceiveBinaryTransportString(backend2); err != nil {
		t.Fatal(err)
	}
	go func() { _, _ = user2.Write([]byte("ani")) }()
	b := make([]byte, 3)
	if _, err := io.ReadFull(backend2, b); err != nil || string(b) != "ani" {
		t.Fatalf("traffic after reconnect: %q, %v", b, err)
	}
}

func dataWSPair(t *testing.T) (*websocket.Conn, *websocket.Conn) {
	t.Helper()
	accepted := make(chan *websocket.Conn, 1)
	up := websocket.Upgrader{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := up.Upgrade(w, r, nil)
		if err == nil {
			accepted <- c
		}
	}))
	t.Cleanup(srv.Close)
	peer, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(srv.URL, "http"), nil)
	if err != nil {
		t.Fatal(err)
	}
	server := <-accepted
	t.Cleanup(func() { server.Close(); peer.Close() })
	return server, peer
}

func drainWS(c *websocket.Conn) {
	go func() {
		for {
			kind, b, err := c.ReadMessage()
			if err != nil {
				return
			}
			if kind == websocket.BinaryMessage && len(b) == 1 && b[0] == utils.SG_RTT {
				if c.WriteMessage(kind, b) != nil {
					return
				}
			}
		}
	}()
}

func TestLostWSClientReleasesActiveDataQuota(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := NewWSServer(ctx, &WsConfig{ChannelSize: 4, Heartbeat: time.Hour, MaxConnections: 1}, quietLogger())
	g := s.newGen(ctx)
	control, peer := dataWSPair(t)
	drainWS(peer)
	s.seatClient(g, control)
	go s.handleLoop(g)
	data, backend := dataWSPair(t)
	s.deliverWebSocket(g, data, g.seat.client())
	local, user := pipePair(t)
	if !s.limits.acquire() {
		t.Fatal("first connection refused")
	}
	g.localChannel <- LocalTCPConn{conn: local, remoteAddr: "127.0.0.1:443", timeCreated: nowMillis()}
	_ = backend.SetReadDeadline(time.Now().Add(3 * time.Second))
	if _, _, err := backend.ReadMessage(); err != nil {
		t.Fatal(err)
	}
	peer.Close()
	eventually(t, 2*time.Second, func() bool { return s.limits.active.Load() == 0 }, "lost WS client leaked an active data connection and its quota")
	_ = user.SetReadDeadline(time.Now().Add(time.Second))
	if _, err := user.Read(make([]byte, 1)); err == nil {
		t.Fatal("WS user stayed open")
	} else if e, ok := err.(interface{ Timeout() bool }); ok && e.Timeout() {
		t.Fatal("user only timed out")
	}
	if ctx.Err() != nil {
		t.Fatal("client loss ended the listening generation")
	}
	control2, peer2 := dataWSPair(t)
	drainWS(peer2)
	s.seatClient(g, control2)
	data2, backend2 := dataWSPair(t)
	s.deliverWebSocket(g, data2, g.seat.client())
	local2, user2 := pipePair(t)
	if !s.limits.acquire() {
		t.Fatal("reconnected WS client has no free quota")
	}
	g.localChannel <- LocalTCPConn{conn: local2, remoteAddr: "127.0.0.1:443", timeCreated: nowMillis()}
	_ = backend2.SetReadDeadline(time.Now().Add(3 * time.Second))
	if _, _, err := backend2.ReadMessage(); err != nil {
		t.Fatal(err)
	}
	go func() { _, _ = user2.Write([]byte("ani")) }()
	if _, b, err := backend2.ReadMessage(); err != nil || string(b) != "ani" {
		t.Fatalf("WS traffic after reconnect: %q, %v", b, err)
	}
}

func TestReplacedWSPoolRequestCannotJoinTheNewClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := NewWSServer(ctx, &WsConfig{ChannelSize: 4, Heartbeat: time.Hour}, quietLogger())
	g := s.newGen(ctx)
	first, peer := dataWSPair(t)
	drainWS(peer)
	s.seatClient(g, first)
	old := g.seat.client()
	second, peer2 := dataWSPair(t)
	drainWS(peer2)
	s.seatClient(g, second)
	data, _ := dataWSPair(t)
	s.deliverWebSocket(g, data, old)
	if len(g.tunnelChannel) != 0 {
		t.Fatal("late old pool request was adopted by the new client")
	}
}
