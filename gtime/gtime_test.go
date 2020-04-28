package gtime

import (
	"testing"
	"time"
)

func TestCuttingSeconds(t *testing.T) {
	dt := time.Date(2020, 4, 28, 11, 30, 27, 0, time.Local)
	v1 := CuttingSeconds(dt)
	v2 := CuttingMinutes(dt)
	v3 := CuttingHours(dt)
	if v1.Second() != 0 {
		t.Error("error cutting seconds")
	}
	if v2.Minute() != 0 || v2.Second() != 0 {
		t.Error("error cutting minutes")
	}
	if v3.Hour() != 0 || v3.Minute() != 0 || v3.Second() != 0 {
		t.Error("error cutting hours")
	}
	t.Log(v1)
	t.Log(v2)
	t.Log(v3)
}
