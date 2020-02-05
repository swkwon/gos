package tcp

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// Server ...
type Server struct {
	addr       *net.TCPAddr
	serverCtx  context.Context
	cancelFunc context.CancelFunc
	sigCh      chan os.Signal
	listener   *net.TCPListener
	handler    HandlerFunc
	done       chan bool
}

// New ...
func New() *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		serverCtx:  ctx,
		cancelFunc: cancel,
		sigCh:      make(chan os.Signal, 1),
		done:       make(chan bool, 1),
	}
}

// RegisterHandler ...
func (t *Server) RegisterHandler(handler HandlerFunc) {
	t.handler = handler
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
		case s := <-t.sigCh:
			log.Println("got signal", s)
			t.done <- true
		}
		t.cancelFunc()
	}()
}

func (t *Server) accept() {
	for {
		sock, err := t.listener.Accept()
		if err == nil {
			s := newSession(t, sock)
			go s.receive()
		}
	}
}
