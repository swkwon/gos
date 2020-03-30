package glog

import (
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	myLogger, err := New(&Config{Type: "console", DateTimeFormat: "RubyDate"})
	if err != nil {
		t.Error(err)
	}

	myLogger.Info("hello world")
	myLogger.Info("I am swkwon")
	myLogger.Info("this is test code")
	// time.Sleep(1 * time.Second)
	myLogger.Info("hello again")
	myLogger.Info("you are my friend")
	myLogger.Info("just do it")

	myLogger.Close()

	myLogger2, err := New(&Config{
		Format:         "json",
		Type:           "console",
		DateTimeFormat: time.RFC3339Nano,
	})

	myLogger2.Info("hello world")
	myLogger2.Info("I am swkwon")
	myLogger2.Info("this is test code")
	myLogger2.Infof("%s %s", "first", "second")
	myLogger2.Close()
}

func BenchmarkSingleLog(b *testing.B) {
	myLogger, err := New(&Config{Type: "console", DateTimeFormat: time.RFC3339Nano})
	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		myLogger.Info("hello world")
		// str := makeText(&message{
		// 	logLevel: infoLevel,
		// 	Message:  "hello world123",
		// }, time.RFC3339Nano)
		// os.Stdout.Write([]byte(str))
	}

	myLogger.Close()
}

func BenchmarkPrint(b *testing.B) {
	type message struct {
		Member int
	}
	for i := 0; i < b.N; i++ {
		_ = make([]byte, 1024)
	}
}
