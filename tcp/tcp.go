package tcp

import (
	"context"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Server ...
type Server struct {
	addr       *net.TCPAddr
	serverCtx  context.Context
	cancelFunc context.CancelFunc
	sigCh      chan os.Signal
	done       chan bool
	listener   *net.TCPListener
	sessionMap sync.Map
}

// New ...
func New() *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		serverCtx:  ctx,
		cancelFunc: cancel,
		sigCh:      make(chan os.Signal, 1),
		done:       make(chan bool, 1),
		sessionMap: sync.Map{},
	}
}

// Start ...
func (t *Server) Start(address string) error {
	if addr, err := net.ResolveTCPAddr("tcp", address); err == nil {
		t.addr = addr
		if l, err := net.ListenTCP("tcp", t.addr); err == nil {
			t.listener = l
		} else {
			return err
		}
	} else {
		return err
	}
	t.setSignal()

	go t.accept()
	<-t.done
	return nil
}

func (t *Server) setSignal() {
	signal.Notify(t.sigCh, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGSEGV)
	go func() {
		select {
		case <-t.sigCh:
			t.done <- true
		}
		t.cancelFunc()
	}()
}

func (t *Server) accept() {
	for {
		select {
		case <-t.serverCtx.Done():
			return
		default:
			sock, err := t.listener.Accept()
			if err == nil {
				t.sessionMap.Store(newSession(t.serverCtx, sock), 1)
			}
		}
	}
}
