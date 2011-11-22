package metrics

import (
	"math"
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

// The standard implementation of an EWMA tracks the number of uncounted
// events and processes them on each tick.  It uses the sync/atomic package
// to manage uncounted events.  When the latest weeklies land in a release,
// atomic.LoadInt64 will be available and this code will become safe on
// 32-bit architectures.
type StandardEWMA struct {
	alpha float64
	uncounted int64
	rate float64
	initialized bool
	tick chan bool
}

// Create a new EWMA with the given alpha.  Create the clock channel and
// start the ticker goroutine.
func NewEWMA(alpha float64) EWMA {
	a := &StandardEWMA{alpha, 0, 0.0, false, make(chan bool)}
	go a.ticker()
	return a
}

// Create a new EWMA with alpha set for a one-minute moving average.
func NewEWMA1() EWMA {
	return NewEWMA(1 - math.Exp(-5.0 / 60.0 / 1))
}

// Create a new EWMA with alpha set for a five-minute moving average.
func NewEWMA5() EWMA {
	return NewEWMA(1 - math.Exp(-5.0 / 60.0 / 5))
}

// Create a new EWMA with alpha set for a fifteen-minute moving average.
func NewEWMA15() EWMA {
	return NewEWMA(1 - math.Exp(-5.0 / 60.0 / 15))
}

// Return the moving average rate of events per second.
func (a *StandardEWMA) Rate() float64 {
	return a.rate * float64(1e9)
}

// Tick the clock to update the moving average.
func (a *StandardEWMA) Tick() {
	a.tick <- true
}

// Add n uncounted events.
func (a *StandardEWMA) Update(n int64) {
	atomic.AddInt64(&a.uncounted, n)
}

// On each clock tick, update the moving average to reflect the number of
// events seen since the last tick.
func (a *StandardEWMA) ticker() {
	for <-a.tick {
		count := a.uncounted
		atomic.AddInt64(&a.uncounted, -count)
		instantRate := float64(count) / float64(5e9)
		if a.initialized {
			a.rate += a.alpha * (instantRate - a.rate)
		} else {
			a.initialized = true
			a.rate = instantRate
		}
	}
}
