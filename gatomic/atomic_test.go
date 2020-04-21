package gatomic

import (
	"sync"
	"testing"
)

func TestAddInt32(t *testing.T) {
	v := &Int32{10000}
	wg := &sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(delta int) {
			v.Add(-1)
			wg.Done()
		}(i)
	}

	wg.Wait()
	value := v.Load()
	if value != 9000 {
		t.Error("invalid value", value)
	}
}