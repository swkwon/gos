package gatomic

import "sync/atomic"

// Int32 ...
type Int32 struct {
	value int32
}

// Int64 ...
type Int64 struct {
	value int64
}

// Add ...
func (v *Int32) Add(delta int32) (new int32) {
	return atomic.AddInt32(&v.value, delta)
}

// Add ...
func (v *Int64) Add(delta int64) (new int64) {
	return atomic.AddInt64(&v.value, delta)
}

// Load ...
func (v *Int32) Load() int32 {
	return atomic.LoadInt32(&v.value)
}

// Load ...
func (v *Int64) Load() int64 {
	return atomic.LoadInt64(&v.value)
}

// Store ...
func (v *Int32) Store(val int32) {
	atomic.StoreInt32(&v.value, val)
}

// Store ...
func (v *Int64) Store(val int64) {
	atomic.StoreInt64(&v.value, val)
}

// Swap ...
func (v *Int32) Swap(val int32) (old int32) {
	return atomic.SwapInt32(&v.value, val)
}

// Swap ...
func (v *Int64) Swap(val int64) (old int64) {
	return atomic.SwapInt64(&v.value, val)
}

// CompareAndSwap ...
func (v *Int32) CompareAndSwap(old, new int32) (swapped bool) {
	return atomic.CompareAndSwapInt32(&v.value, old, new)
}

// CompareAndSwap ...
func (v *Int64) CompareAndSwap(old, new int64) (swapped bool) {
	return atomic.CompareAndSwapInt64(&v.value, old, new)
}
