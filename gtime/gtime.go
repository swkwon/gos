package gtime

import (
	"time"

	"github.com/swkwon/gos/gatomic"
)

var timeOffset gatomic.Int64

const dateTimeFormat = "2006-01-02 15:04:05"

func init() {
	timeOffset.Store(0)
}

// Now ...
func Now() time.Time {
	return time.Now().Add(time.Duration(timeOffset.Load()))
}

// AddOffset ...
func AddOffset(t time.Duration) {
	timeOffset.Store(int64(t))
}

// GetTimeOffset ...
func GetTimeOffset() time.Duration {
	return time.Duration(timeOffset.Load())
}

// ResetOffset ...
func ResetOffset() {
	timeOffset.Store(0)
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
