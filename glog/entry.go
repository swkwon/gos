package glog

import "fmt"

// IEntry ...
type IEntry interface {
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warning(v ...interface{})
	Warningf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
}

type entry struct {
	field       Fields
	logInstance *logger
}

func (e *entry) Info(v ...interface{}) {
	e.logInstance.into(infoLevel, fmt.Sprint(v...), e.field)
}

func (e *entry) Infof(format string, v ...interface{}) {
	e.logInstance.into(infoLevel, fmt.Sprintf(format, v...), e.field)
}

func (e *entry) Warning(v ...interface{}) {
	e.logInstance.into(warningLevel, fmt.Sprint(v...), e.field)
}

func (e *entry) Warningf(format string, v ...interface{}) {
	e.logInstance.into(warningLevel, fmt.Sprintf(format, v...), e.field)
}

func (e *entry) Error(v ...interface{}) {
	e.logInstance.into(errorLevel, fmt.Sprint(v...), e.field)
}

func (e *entry) Errorf(format string, v ...interface{}) {
	e.logInstance.into(errorLevel, fmt.Sprintf(format, v...), e.field)
}

func (e *entry) Debug(v ...interface{}) {
	e.logInstance.into(debugLevel, fmt.Sprint(v...), e.field)
}

func (e *entry) Debugf(format string, v ...interface{}) {
	e.logInstance.into(debugLevel, fmt.Sprintf(format, v...), e.field)
}
