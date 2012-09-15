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
// to manage uncounted events.
type StandardEWMA struct {
	alpha     float64
	uncounted int64
	in        chan bool
	out       chan float64
}

// Create a new EWMA with the given alpha.  Create the clock channel and
// start the ticker goroutine.
func NewEWMA(alpha float64) *StandardEWMA {
	a := &StandardEWMA{alpha, 0, make(chan bool), make(chan float64)}
	go a.arbiter()
	return a
}

// Create a new EWMA with alpha set for a one-minute moving average.
func NewEWMA1() *StandardEWMA {
	return NewEWMA(1 - math.Exp(-5.0/60.0/1))
}

// Create a new EWMA with alpha set for a five-minute moving average.
func NewEWMA5() *StandardEWMA {
	return NewEWMA(1 - math.Exp(-5.0/60.0/5))
}

// Create a new EWMA with alpha set for a fifteen-minute moving average.
func NewEWMA15() *StandardEWMA {
	return NewEWMA(1 - math.Exp(-5.0/60.0/15))
}

// Return the moving average rate of events per second.
func (a *StandardEWMA) Rate() float64 {
	return <-a.out * float64(1e9)
}

// Tick the clock to update the moving average.
func (a *StandardEWMA) Tick() {
	a.in <- true
}

// Add n uncounted events.
func (a *StandardEWMA) Update(n int64) {
	atomic.AddInt64(&a.uncounted, n)
}

// On each clock tick, update the moving average to reflect the number of
// events seen since the last tick.
func (a *StandardEWMA) arbiter() {
	var initialized bool
	var rate float64
	for {
		select {
		case <-a.in:
			count := atomic.LoadInt64(&a.uncounted)
			atomic.AddInt64(&a.uncounted, -count)
			instantRate := float64(count) / float64(5e9)
			if initialized {
				rate += a.alpha * (instantRate - rate)
			} else {
				initialized = true
				rate = instantRate
			}
		case a.out <- rate:
		}
	}
}
