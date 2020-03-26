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
	Message     string                 `json:"message"`
	Parameter   map[string]interface{} `json:"parameter,omitempty"`
	DateTime    string                 `json:"date_time,omitempty"`
	LogLevelStr string                 `json:"log_level"`
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

	return fmt.Sprintf("[%s] [%5v] %s\n", m.DateTime, m.LogLevelStr, m.Message)
}
