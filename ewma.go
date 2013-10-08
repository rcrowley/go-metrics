package metrics

import (
	"math"
	"sync"
	"sync/atomic"
)

// EWMAs continuously calculate an exponentially-weighted moving average
// based on an outside source of clock ticks.
//
// This is an interface so as to encourage other structs to implement
// the EWMA API as appropriate.
type EWMA interface {
	Rate() float64
	Tick()
	Update(int64)
}

// Create a new EWMA with the given alpha.
func NewEWMA(alpha float64) EWMA {
	if UseNilMetrics {
		return NilEWMA{}
	}
	return &StandardEWMA{alpha: alpha}
}

// Create a new EWMA with alpha set for a one-minute moving average.
func NewEWMA1() EWMA {
	return NewEWMA(1 - math.Exp(-5.0/60.0/1))
}

// Create a new EWMA with alpha set for a five-minute moving average.
func NewEWMA5() EWMA {
	return NewEWMA(1 - math.Exp(-5.0/60.0/5))
}

// Create a new EWMA with alpha set for a fifteen-minute moving average.
func NewEWMA15() EWMA {
	return NewEWMA(1 - math.Exp(-5.0/60.0/15))
}

// No-op EWMA.
type NilEWMA struct{}

// No-op.
func (a NilEWMA) Rate() float64 { return 0.0 }

// No-op.
func (a NilEWMA) Tick() {}

// No-op.
func (a NilEWMA) Update(n int64) {}

// The standard implementation of an EWMA tracks the number of uncounted
// events and processes them on each tick.  It uses the sync/atomic package
// to manage uncounted events.
type StandardEWMA struct {
	alpha     float64
	init      bool
	mutex     sync.Mutex
	rate      float64
	uncounted int64
}

// Return the moving average rate of events per second.
func (a *StandardEWMA) Rate() float64 {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.rate * float64(1e9)
}

// Tick the clock to update the moving average.
func (a *StandardEWMA) Tick() {
	count := atomic.LoadInt64(&a.uncounted)
	atomic.AddInt64(&a.uncounted, -count)
	instantRate := float64(count) / float64(5e9)
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.init {
		a.rate += a.alpha * (instantRate - a.rate)
	} else {
		a.init = true
		a.rate = instantRate
	}
}

// Add n uncounted events.
func (a *StandardEWMA) Update(n int64) {
	atomic.AddInt64(&a.uncounted, n)
}
