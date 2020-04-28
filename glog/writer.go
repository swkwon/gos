package glog

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/swkwon/gos/gtime"
)

const (
	rotationHour = 0
	rotationDay  = 1
)

// IWriter ...
type IWriter interface {
	Close() error
	Write(v []byte, offset int) (int, error)
}

type tcpWriter struct {
	host string
	wc   io.WriteCloser
}

type udpWriter struct {
	host string
	wc   io.WriteCloser
}

type stdoutWriter struct {
	wc io.WriteCloser
}

type fileWriter struct {
	path         string
	file         string
	rotation     time.Duration
	rotationType int
	generate     time.Time
	wc           io.WriteCloser
}

func makeTCPWriter(host string) (IWriter, error) {
	if len(host) <= 0 {
		return nil, errors.New("tcp host info is zero")
	}

	c, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	return &tcpWriter{
		host: host,
		wc:   c,
	}, nil
}

func makeUDPWriter(host string) (IWriter, error) {
	if len(host) <= 0 {
		return nil, errors.New("udp host info is zero")
	}
	c, err := net.Dial("udp", host)
	if err != nil {
		return nil, err
	}

	return &udpWriter{
		host: host,
		wc:   c,
	}, nil
}

func makeSTDOutWriter() (IWriter, error) {
	return &stdoutWriter{
		wc: os.Stdout,
	}, nil
}

func makeFileWriter(file *FileConfig) (IWriter, error) {
	fullpath := filepath.Join(file.Path, file.FileName)
	f, e := os.OpenFile(fullpath, os.O_CREATE|os.O_WRONLY, os.FileMode(0644))
	if e != nil {
		return nil, e
	}
	var rot time.Duration
	var rt int
	if file.Rotation == "hour" {
		rot = 60 * time.Minute
		rt = rotationHour
	} else {
		rot = 24 * time.Hour
		rt = rotationDay
	}
	return &fileWriter{
		file:         file.FileName,
		path:         file.Path,
		rotation:     rot,
		rotationType: rt,
		generate:     time.Now(),
		wc:           f,
	}, nil
}

func (w *tcpWriter) Write(v []byte, offset int) (int, error) {
	if w.wc == nil {
		return 0, errors.New("writer is nil")
	}

	return w.wc.Write(v[offset:])
}

func (w *tcpWriter) Close() error {
	if w.wc == nil {
		return nil
	}

	return w.wc.Close()
}

func (w *udpWriter) Write(v []byte, offset int) (int, error) {
	if w.wc == nil {
		return 0, errors.New("writer is nil")
	}

	return w.wc.Write(v[offset:])
}

func (w *udpWriter) Close() error {
	if w.wc == nil {
		return nil
	}

	return w.wc.Close()
}

func (w *stdoutWriter) Write(v []byte, offset int) (int, error) {
	if w.wc == nil {
		return 0, errors.New("writer is nil")
	}

	return w.wc.Write(v[offset:])
}

func (w *stdoutWriter) Close() error {
	return nil
}

func (w *fileWriter) Write(v []byte, offset int) (int, error) {
	if w.wc == nil {
		return 0, errors.New("writer is nil")
	}
	if offset == 0 {
		w.checkRotation()
	}
	return w.wc.Write(v[offset:])
}

func (w *fileWriter) Close() error {
	if w.wc == nil {
		return nil
	}

	return w.wc.Close()
}

func (w *fileWriter) isRotation() bool {
	var sub time.Duration
	if w.rotationType == rotationHour {
		sub = gtime.Now().Sub(gtime.CuttingMinutes(w.generate))
	} else {
		sub = gtime.Now().Sub(gtime.CuttingHours(w.generate))
	}
	if sub > w.rotation {
		return true
	}
	return false
}

func (w *fileWriter) checkRotation() {
	if w.isRotation() {
		// make daily folder
		folderName := fmt.Sprintf("%d%02d%02d", w.generate.Year(), w.generate.Month(), w.generate.Day())
		target := filepath.Join(w.path, folderName)
		_ = os.Mkdir(target, os.FileMode(0777))

		// copy file
		_ = w.wc.Close()
		srcFile := filepath.Join(w.path, w.file)
		ext := filepath.Ext(w.file)
		fileName := strings.TrimSuffix(w.file, ext)
		backup := fmt.Sprintf("%s_%d%02d%02d%02d%02d%02d.log", fileName, w.generate.Year(), w.generate.Month(), w.generate.Day(), w.generate.Hour(), w.generate.Minute(), w.generate.Second())
		dstFile := filepath.Join(target, backup)
		fileCopy(srcFile, dstFile)

		// new file
		fullpath := filepath.Join(w.path, w.file)
		f, e := os.OpenFile(fullpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0644))
		if e == nil {
			w.wc = f
		} else {
			w.wc = nil
		}
	}
}

func fileCopy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
