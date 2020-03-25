package glog

import (
	"context"
	"fmt"
)

// GLogger ...
type GLogger interface {
	Info(v ...interface{})
}

type logger struct {
	Message chan []string
	Ctx     context.Context
	cancel  context.CancelFunc
}

func (l *logger) Info(v ...interface{}) {
	l.Message <- fmt.Sprintln(v...)
}

// New ...
func New(config *Config) GLogger {
	return makeLogger(config, context.Background())
}

// NewWithContext ...
func NewWithContext(ctx context.Context, config *Config) GLogger {
	return makeLogger(config, ctx)
}

func makeLogger(ctx context.Context, config *Config) *logger {
	newCtx, cancelFunc := context.WithCancel(context.Background())
	ins := &logger{
		Message: make(chan []string, 1000),
		Ctx:     newCtx,
		cancel:  cancelFunc,
	}

	go func() {
		for {
			select {
			case msg := <-ins.Message:
				ins.Print(msg)
			case <-ins.Ctx.Done():
				return
			}
		}
	}()
}
