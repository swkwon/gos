package tcp

import "context"

// Session ...
type Session struct {
	ctx context.Context
}

func newSession(ctx context.Context) *Session {
	return &Session{
		ctx: ctx,
	}
}
