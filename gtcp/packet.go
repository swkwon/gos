package gtcp

import (
	"bytes"
	"encoding/json"
)

// Context ...
type Context interface {
	JSON(v interface{}) error
	JSONCompact(v interface{}) error
	GetReceived() []byte
}

type sessionContext struct {
	session  *Session
	received []byte
}

// GetReceived ...
func (c *sessionContext) GetReceived() []byte {
	return c.received
}

func marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// JSON ...
func (c *sessionContext) JSON(v interface{}) error {
	b, e := marshal(v)
	if e != nil {
		return e
	}
	c.session.send(b)
	return nil
}

// JSONCompact ...
func (c *sessionContext) JSONCompact(v interface{}) error {
	b, e := marshal(v)
	if e != nil {
		return e
	}
	cb := new(bytes.Buffer)
	e = json.Compact(cb, b)
	if e != nil {
		return e
	}
	c.session.send(cb.Bytes())
	return nil
}

// NewContext ...
func NewContext(s *Session, raw []byte) (Context, error) {
	return &sessionContext{
		session:  s,
		received: raw,
	}, nil
}
