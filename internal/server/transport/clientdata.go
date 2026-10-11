package transport

import (
	"context"
	"net"
)

// clientDataConn carries the owner captured under the seat's admission lock.
// Reading the current seat when the relay starts is too late: a replacement
// may already have arrived, letting old sockets outlive the client they served.
// The wrapper is removed before the relay to preserve TCP half-close/splice.
type clientDataConn struct {
	net.Conn
	owner context.Context
}

func dataConnection(c net.Conn, fallback context.Context) (context.Context, net.Conn) {
	if owned, ok := c.(*clientDataConn); ok {
		return owned.owner, owned.Conn
	}
	return fallback, c
}
