package glog

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type level int

const (
	infoLevel level = iota
	warningLevel
	errorLevel
	debugLevel
)

type message struct {
	logLevel    level
	Message     string `json:"message"`
	DateTime    string `json:"date_time,omitempty"`
	LogLevelStr string `json:"log_level"`
	FileName    string `json:"file_name"`
	Line        int    `json:"line"`
	FuncName    string `json:"func_name"`
	Param       Fields `json:"_fields_,omitempty"`
}

var logLevelMap = map[level]string{
	infoLevel:    "INFO",
	warningLevel: "WARN",
	errorLevel:   "ERROR",
	debugLevel:   "DEBUG",
}

func getLogLevel(logLevel string) level {
	switch {
	case strings.EqualFold(logLevel, "info"):
		return infoLevel
	case strings.EqualFold(logLevel, "warning"):
		return warningLevel
	case strings.EqualFold(logLevel, "error"):
		return errorLevel
	}
	return debugLevel
}

func makeJSON(m *message, timeFormat string) string {
	m.DateTime = time.Now().Format(timeFormat)
	if i, ok := logLevelMap[m.logLevel]; ok {
		m.LogLevelStr = i
	}

	b, e := json.Marshal(m)
	if e == nil {
		return string(b) + "\n"
	}
	return "{\"error\":\"make json\"}\n"
}

func makeText(m *message, timeFormat string) string {
	m.DateTime = time.Now().Format(timeFormat)
	if i, ok := logLevelMap[m.logLevel]; ok {
		m.LogLevelStr = i
	}

	var str []string
	if m.Param != nil {
		for k, v := range m.Param {
			str = append(str, fmt.Sprintf("%v=%v", k, v))
		}

		return fmt.Sprintf("[%s] [%s:%d:%s] [%5v] %s %s\n", m.DateTime, m.FileName, m.Line, m.FuncName, m.LogLevelStr, m.Message, strings.Join(str, " "))
	}

	return fmt.Sprintf("[%s] [%s:%d:%s] [%5v] %s\n", m.DateTime, m.FileName, m.Line, m.FuncName, m.LogLevelStr, m.Message)
}
