package glog

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Logger ...
type Logger interface {
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warning(v ...interface{})
	Warningf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Close() []error
}

type makeFuncType func(*message, string) string

type logger struct {
	c          chan *message
	ctx        context.Context
	cancel     context.CancelFunc
	writer     IWriter
	wg         *sync.WaitGroup
	makeFunc   makeFuncType
	timeFormat string
	logLevel   level
	subs       []*logger
}

// New ...
func New(config *Config) (Logger, error) {
	return makeLogger(context.Background(), config)
}

func (l *logger) Info(v ...interface{}) {
	l.into(infoLevel, fmt.Sprint(v...))
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.into(infoLevel, fmt.Sprintf(format, v...))
}

func (l *logger) Warning(v ...interface{}) {
	l.into(warningLevel, fmt.Sprint(v...))
}

func (l *logger) Warningf(format string, v ...interface{}) {
	l.into(warningLevel, fmt.Sprintf(format, v...))
}

func (l *logger) Error(v ...interface{}) {
	l.into(errorLevel, fmt.Sprint(v...))
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.into(errorLevel, fmt.Sprintf(format, v...))
}

func (l *logger) Debug(v ...interface{}) {
	l.into(debugLevel, fmt.Sprint(v...))
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.into(debugLevel, fmt.Sprintf(format, v...))
}

func (l *logger) Close() []error {
	if l.cancel != nil {
		l.cancel()
	}
	if l.wg != nil {
		l.wg.Wait()
	}

	for _, v := range l.subs {
		v.wg.Wait()
	}

	var errs []error
	if l.writer != nil {
		if e := l.writer.Close(); e != nil {
			errs = append(errs, e)
		}
	}
	for _, v := range l.subs {
		if e := v.Close(); e != nil {
			errs = append(errs, e...)
		}
	}
	return errs
}

func (l *logger) writeAll() {
	for elem := range l.c {
		l.print(l.makeFunc(elem, l.timeFormat))
	}
}

func (l *logger) print(v string) {
	bytes := []byte(v)
	length := len(bytes)
	for offset := 0; offset < length; {
		if n, err := l.writer.Write(bytes, offset); err != nil {
			break
		} else {
			offset += n
		}
	}
}

func (l *logger) logging(logLevel level) bool {
	return l.logLevel >= logLevel
}

func (l *logger) into(logLevel level, msg string) {
	needLogging := false
	for _, v := range l.subs {
		if v.logging(logLevel) {
			needLogging = true
		}
	}

	if l.logging(logLevel) == true {
		needLogging = true
	}

	if needLogging == false {
		return
	}

	pc, file, line, _ := runtime.Caller(2)
	_, fileName := path.Split(file)
	parts1 := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	parts2 := strings.SplitN(parts1[len(parts1)-1], ".", 2)
	funcName := parts2[1]

	if l.logging(logLevel) {
		l.c <- &message{
			logLevel: logLevel,
			Message:  msg,
			FileName: fileName,
			Line:     line,
			FuncName: funcName,
		}
	}

	for _, v := range l.subs {
		if v.logging(logLevel) {
			v.c <- &message{
				logLevel: logLevel,
				Message:  msg,
				FileName: fileName,
				Line:     line,
				FuncName: funcName,
			}
		}
	}
}

func makeLogger(ctx context.Context, config *Config) (*logger, error) {
	checkConfig(config)
	writer, err := makeWriter(config)
	if err != nil {
		log.Println("failed make logger >", err)
		return nil, err
	}
	newCtx, cancelFunc := context.WithCancel(ctx)
	l := &logger{
		c:          make(chan *message, 1000),
		ctx:        newCtx,
		cancel:     cancelFunc,
		writer:     writer,
		wg:         &sync.WaitGroup{},
		makeFunc:   getMakeFunc(config.Format),
		timeFormat: config.DateTimeFormat,
		logLevel:   getLogLevel(config.LogLevel),
	}

	l.wg.Add(1)
	go job(l)

	for _, v := range config.Sub {
		checkConfig(v)
		subWriter, err := makeWriter(v)
		if err != nil {
			log.Println("failed make sub logger >", err)
			continue
		}
		sl := &logger{
			c:          make(chan *message, 1000),
			ctx:        newCtx,
			logLevel:   getLogLevel(v.LogLevel),
			wg:         &sync.WaitGroup{},
			makeFunc:   getMakeFunc(v.Format),
			timeFormat: v.DateTimeFormat,
			writer:     subWriter,
		}
		sl.wg.Add(1)
		l.subs = append(l.subs, sl)
		go job(sl)
	}

	return l, nil
}

func makeWriter(config *Config) (IWriter, error) {

	switch strings.ToLower(config.Type) {
	case "console":
		return makeSTDOutWriter()
	case "tcp":
		if config.TCP == nil {
			return nil, errors.New("tcp info is nil")
		}
		return makeTCPWriter(config.TCP.Host)
	case "udp":
		if config.UDP == nil {
			return nil, errors.New("udp info is nil")
		}
		return makeUDPWriter(config.UDP.Host)
	case "file":
		if config.File == nil {
			return nil, errors.New("file info is nil")
		}
		return makeFileWriter(config.File)
	}

	return nil, fmt.Errorf("invalid config type: %s", config.Type)
}

func checkConfig(config *Config) {
	if config.Format == "" {
		config.Format = "text"
	}
	if config.DateTimeFormat == "" {
		config.DateTimeFormat = time.RFC3339
	} else {
		config.DateTimeFormat = dateTimeFormatParsing(config.DateTimeFormat)
	}
}

func getMakeFunc(textFormat string) makeFuncType {
	switch {
	case strings.EqualFold(textFormat, "json"):
		return makeJSON
	default:
		return makeText
	}
}

func job(l *logger) {
	for {
		select {
		case v := <-l.c:
			l.print(l.makeFunc(v, l.timeFormat))
		case <-l.ctx.Done():
			close(l.c)
			l.writeAll()
			l.wg.Done()
			return
		}
	}
}

var timeFormat = map[string]string{
	"ansic":       time.ANSIC,
	"unixdate":    time.UnixDate,
	"rubydate":    time.RubyDate,
	"rfc822":      time.RFC822,
	"rfc822z":     time.RFC822Z,
	"rfc850":      time.RFC850,
	"rfc1123":     time.RFC1123,
	"rfc1123z":    time.RFC1123Z,
	"rfc3339":     time.RFC3339,
	"rfc3339nano": time.RFC3339Nano,
	"kitchen":     time.Kitchen,
	"stamp":       time.Stamp,
	"stampmilli":  time.StampMilli,
	"stampmicro":  time.StampMicro,
	"stampnano":   time.StampNano,
}

func dateTimeFormatParsing(s string) string {
	if format, ok := timeFormat[strings.ToLower(s)]; ok {
		return format
	}
	return s
}
