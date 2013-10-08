package metrics

import "sync/atomic"

// Counters hold an int64 value that can be incremented and decremented.
//
// This is an interface so as to encourage other structs to implement
// the Counter API as appropriate.
type Counter interface {
	Clear()
	Count() int64
	Dec(int64)
	Inc(int64)
}

// Create a new Counter.
func NewCounter() Counter {
	if UseNilMetrics {
		return NilCounter{}
	}
	return &StandardCounter{0}
}

// No-op Counter.
type NilCounter struct{}

// No-op.
func (c NilCounter) Clear() {}

// No-op.
func (c NilCounter) Count() int64 { return 0 }

// No-op.
func (c NilCounter) Dec(i int64) {}

// No-op.
func (c NilCounter) Inc(i int64) {}

// The standard implementation of a Counter uses the sync/atomic package
// to manage a single int64 value.
type StandardCounter struct {
	count int64
}

// Clear the counter: set it to zero.
func (c *StandardCounter) Clear() {
	atomic.StoreInt64(&c.count, 0)
}

// Return the current count.
func (c *StandardCounter) Count() int64 {
	return atomic.LoadInt64(&c.count)
}

// Decrement the counter by the given amount.
func (c *StandardCounter) Dec(i int64) {
	atomic.AddInt64(&c.count, -i)
}

// Increment the counter by the given amount.
func (c *StandardCounter) Inc(i int64) {
	atomic.AddInt64(&c.count, i)
}
