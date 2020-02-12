package tcp

import (
	"context"
	"encoding/binary"
	"net"
)

const (
	defaultBufferLength = 2048
)

// HandlerFunc ...
type HandlerFunc func(c Context)

// Session ...
type Session struct {
	ctx     context.Context
	sock    *net.TCPConn
	buffer  []byte
	offset  int
	handler HandlerFunc
	sendCh  chan []byte
}

func newSession(t *Server, sock net.Conn) *Session {
	tc := sock.(*net.TCPConn)
	tc.SetKeepAlive(true)
	tc.SetNoDelay(true)
	return &Session{
		ctx:     t.serverCtx,
		sock:    tc,
		handler: t.handler,
		sendCh:  make(chan []byte, 100),
	}
}

func (s *Session) receiveProcess() {
	defer func() {
		if err := recover(); err != nil {
			s.close()
		}
	}()

	if s.buffer == nil {
		s.buffer = make([]byte, defaultBufferLength)
	}
	for {
		n, e := s.sock.Read(s.buffer[s.offset:])
		if e != nil {
			return
		}

		if n <= 0 {
			return
		}

		bodyLength := binary.LittleEndian.Uint32(s.buffer[:4])

		if s.offset >= 4 && int(bodyLength) > len(s.buffer) {
			temp := make([]byte, bodyLength)
			copy(temp, s.buffer)
			s.buffer = temp
		}

		if s.offset+n < int(bodyLength)+4 {
			s.offset = s.offset + n
		} else {
			received := make([]byte, bodyLength)
			copy(received, s.buffer[4:int(bodyLength)+4])
			if e := s.received(received); e != nil {
				return
			}
			if s.offset+n > int(bodyLength)+4 {
				copy(s.buffer, s.buffer[int(bodyLength)+4:s.offset+n])
			}
			s.offset = (s.offset + n) - (int(bodyLength) + 4)
		}
	}
}

func (s *Session) received(raw []byte) error {
	if s.handler != nil {
		c, e := NewContext(s, raw)
		if e != nil {
			return e
		}
		s.handler(c)
	}
	return nil
}

// Send ...
func (s *Session) send(v []byte) {
	size := len(v)
	header := make([]byte, 4)
	binary.LittleEndian.PutUint32(header, uint32(size))
	s.sendCh <- append(header, v...)
}

// Close ...
func (s *Session) close() {
	if s != nil {
		s.sock.Close()
	}
}

func (s *Session) sendProcess() {
	defer func() {
		if e := recover(); e != nil {
			s.close()
		}
	}()

	for {
		select {
		case v := <-s.sendCh:
			transferred := 0
			for {
				n, _ := s.sock.Write(v[transferred:])
				transferred += n
				if transferred >= len(v) {
					break
				}
			}
		case <-s.ctx.Done():
			break
		}
	}
}
