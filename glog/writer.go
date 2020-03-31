package glog

import (
	"errors"
	"io"
	"net"
	"os"
	"strings"
)

type glogWriter struct {
	writerType string
	wc         io.WriteCloser
}

func makeTCPWriter(host string) (io.WriteCloser, error) {
	if len(host) <= 0 {
		return nil, errors.New("tcp host info is zero")
	}

	c, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	return &glogWriter{
		writerType: "tcp",
		wc:         c,
	}, nil
}

func makeUDPWriter(host string) (io.WriteCloser, error) {
	if len(host) <= 0 {
		return nil, errors.New("udp host info is zero")
	}
	c, err := net.Dial("udp", host)
	if err != nil {
		return nil, err
	}

	return &glogWriter{
		writerType: "udp",
		wc:         c,
	}, nil
}

func makeSTDOutWriter() (io.WriteCloser, error) {
	return &glogWriter{
		writerType: "console",
		wc:         os.Stdout,
	}, nil
}

func makeFileWriter(file *FileConfig) (io.WriteCloser, error) {
	return nil, errors.New("not implement")
}

// Write ...
func (w *glogWriter) Write(v []byte) (int, error) {
	if w == nil || w.wc == nil {
		return 0, errors.New("writer is nil")
	}

	return w.wc.Write(v)
}

// Close ...
func (w *glogWriter) Close() error {
	if w == nil || w.wc == nil {
		return nil
	}

	if strings.EqualFold(w.writerType, "console") {
		return nil
	}

	return w.wc.Close()
}
