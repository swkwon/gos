package tcp

import (
	"context"
	"encoding/binary"
	"net"
)

const (
	defaultBufferLength = 2048
)

type HandlerFunc func(*Session, []byte)

// Session ...
type Session struct {
	ctx     context.Context
	sock    *net.TCPConn
	buffer  []byte
	offset  int
	handler HandlerFunc
}

func newSession(t *Server, sock net.Conn) *Session {
	return &Session{
		ctx:     t.serverCtx,
		sock:    sock.(*net.TCPConn),
		handler: t.handler,
	}
}

func (s *Session) receive() {
	if s.buffer == nil {
		s.buffer = make([]byte, defaultBufferLength)
	}
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
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
				s.received(received)
				if s.offset+n > int(bodyLength)+4 {
					copy(s.buffer, s.buffer[int(bodyLength)+4:s.offset+n])
				}
				s.offset = (s.offset + n) - (int(bodyLength) + 4)
			}
		}
	}
}

func (s *Session) received(data []byte) {
	if s.handler != nil {
		s.handler(s, data)
	}
}

func (s *Session) Send(data []byte) {
	sendBuffer := make([]byte, len(data)+4)
	binary.LittleEndian.PutUint32(sendBuffer, uint32(len(data)))
}
