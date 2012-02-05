package metrics

import "sync/atomic"

// Gauges hold an int64 value that can be set arbitrarily.
//
// This is an interface so as to encourage other structs to implement
// the Gauge API as appropriate.
type Gauge interface {
	Update(int64)
	Value() int64
}

// The standard implementation of a Gauge uses the sync/atomic package
// to manage a single int64 value.
type StandardGauge struct {
	value int64
}

// Create a new gauge.
func NewGauge() *StandardGauge {
	return &StandardGauge{0}
}

// Update the gauge's value.
func (g *StandardGauge) Update(v int64) {
	atomic.StoreInt64(&g.value, v)
}

// Return the gauge's current value.
func (g *StandardGauge) Value() int64 {
	return atomic.LoadInt64(&g.value)
}
