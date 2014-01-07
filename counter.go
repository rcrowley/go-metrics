package metrics

import "sync/atomic"

// Counters hold an int64 value that can be incremented and decremented.
type Counter interface {
	Clear()
	Count() int64
	Dec(int64)
	Inc(int64)
	Snapshot() Counter
}

// GetOrRegisterCounter returns an existing Counter or constructs and registers
// a new StandardCounter.
func GetOrRegisterCounter(name string, r Registry) Counter {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewCounter()).(Counter)
}

// NewCounter constructs a new StandardCounter.
func NewCounter() Counter {
	if UseNilMetrics {
		return NilCounter{}
	}
	return &StandardCounter{0}
}

// NewRegisteredCounter constructs and registers a new StandardCounter.
func NewRegisteredCounter(name string, r Registry) Counter {
	c := NewCounter()
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// No-op Counter.
type NilCounter struct{}

// Clear is a no-op.
func (NilCounter) Clear() {}

// Count is a no-op.
func (NilCounter) Count() int64 { return 0 }

// Dec is a no-op.
func (NilCounter) Dec(i int64) {}

// Inc is a no-op.
func (NilCounter) Inc(i int64) {}

// The standard implementation of a Counter uses the sync/atomic package
// to manage a single int64 value.
type StandardCounter struct {
	count int64
}

// Clear sets the counter to zero.
func (c *StandardCounter) Clear() {
	atomic.StoreInt64(&c.count, 0)
}

// Count returns the current count.
func (c *StandardCounter) Count() int64 {
	return atomic.LoadInt64(&c.count)
}

// Dec decrements the counter by the given amount.
func (c *StandardCounter) Dec(i int64) {
	atomic.AddInt64(&c.count, -i)
}

// Inc increments the counter by the given amount.
func (c *StandardCounter) Inc(i int64) {
	atomic.AddInt64(&c.count, i)
}
