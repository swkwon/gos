package glog

import (
	"errors"
	"io"
	"net"
	"os"
)

func makeTCPWriter(host string) (io.Writer, error) {
	if len(host) <= 0 {
		return nil, errors.New("tcp host info is zero")
	}

	return net.Dial("tcp", host)
}

func makeUDPWriter(host string) (io.Writer, error) {
	if len(host) <= 0 {
		return nil, errors.New("udp host info is zero")
	}
	return net.Dial("udp", host)
}

func makeSTDOutWriter() (io.WriteCloser, error) {
	return os.Stdout, nil
}
