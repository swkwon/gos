package gtime

import (
	"time"
)

const (
	dateTimeFormat = "2006-01-02 15:04:05"
)

// Now ...
func Now() time.Time {
	return time.Now()
}

// UTCNow ...
func UTCNow() time.Time {
	return time.Now().UTC()
}

// CuttingSeconds ...
func CuttingSeconds(t time.Time) time.Time {
	return t.Truncate(time.Minute)
}

// CuttingMinutes ...
func CuttingMinutes(t time.Time) time.Time {
	return t.Truncate(time.Hour)
}

// CuttingHours ...
func CuttingHours(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
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
