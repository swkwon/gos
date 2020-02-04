package tcp

import (
	"context"
	"net"
)

// Session ...
type Session struct {
	ctx  context.Context
	sock net.Conn
	Key  string
}

func newSession(ctx context.Context, sock net.Conn) *Session {
	return &Session{
		ctx:  ctx,
		sock: sock,
	}
}
