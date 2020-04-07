package gtime

import (
	"sync"
	"time"
)

var timeOffset time.Duration
var mutex sync.Mutex

const dateTimeFormat = "2006-01-02 15:04:05"

func init() {
	timeOffset = 0
}

// Now ...
func Now() time.Time {
	return time.Now().Add(timeOffset)
}

// AddOffset ...
func AddOffset(t time.Duration) {
	mutex.Lock()
	defer mutex.Unlock()
	timeOffset = t
}

// GetTimeOffset ...
func GetTimeOffset() time.Duration {
	return timeOffset
}

// ResetOffset ...
func ResetOffset() {
	timeOffset = 0
}

// TimeToMySQLString ...
func TimeToMySQLString(t time.Time) string {
	return t.Format(dateTimeFormat)
}

// MySQLStringToTime ...
func MySQLStringToTime(dt string) (time.Time, error) {
	result, err := time.ParseInLocation(dateTimeFormat, dt, time.Local)
	if err != nil {
		return time.Time{}, err
	}
	return result, nil
}
