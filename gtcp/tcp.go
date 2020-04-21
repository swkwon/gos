package gtcp

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	name = `

     __________  _____
    / ____/ __ \/ ___/
   / / __/ / / /\__ \
  / /_/ / /_/ /___/ /
  \____/\____//____/   v0.0.1
           simple, fast, easy
-----------------------------
  `
)

// Server ...
type Server struct {
	addr       *net.TCPAddr
	serverCtx  context.Context
	cancelFunc context.CancelFunc
	sigCh      chan os.Signal
	listener   *net.TCPListener
	handler    HandlerFunc
}

// New ...
func New() *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		serverCtx:  ctx,
		cancelFunc: cancel,
		sigCh:      make(chan os.Signal, 1),
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
	t.accept()
	return nil
}

func (t *Server) setSignal() {
	signal.Notify(t.sigCh, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGSEGV)
	go func() {
		select {
		case s := <-t.sigCh:
			log.Println("got signal", s)
			t.listener.Close()
		}
		t.cancelFunc()
	}()
}

func (t *Server) accept() {
	fmt.Println(name)
	log.Println("TCP server started on", t.addr.String())
	for {
		sock, err := t.listener.Accept()
		if err == nil {
			s := newSession(t, sock)
			go s.receiveProcess()
			go s.sendProcess()
		} else {
			log.Fatal(err)
		}
	}
}
