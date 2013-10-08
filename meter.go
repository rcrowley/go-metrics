package metrics

import "time"

// Meters count events to produce exponentially-weighted moving average rates
// at one-, five-, and fifteen-minutes and a mean rate.
//
// This is an interface so as to encourage other structs to implement
// the Meter API as appropriate.
type Meter interface {
	Count() int64
	Mark(int64)
	Rate1() float64
	Rate5() float64
	Rate15() float64
	RateMean() float64
}

// Create a new Meter.  Create the communication channels and start the
// synchronizing goroutine.
func NewMeter() Meter {
	if UseNilMetrics {
		return NilMeter{}
	}
	m := &StandardMeter{
		make(chan int64),
		make(chan meterV),
		time.NewTicker(5e9),
	}
	go m.arbiter()
	return m
}

// No-op Meter.
type NilMeter struct{}

// No-op.
func (m NilMeter) Count() int64 { return 0 }

// No-op.
func (m NilMeter) Mark(n int64) {}

// No-op.
func (m NilMeter) Rate1() float64 { return 0.0 }

// No-op.
func (m NilMeter) Rate5() float64 { return 0.0 }

// No-op.
func (m NilMeter) Rate15() float64 { return 0.0 }

// No-op.
func (m NilMeter) RateMean() float64 { return 0.0 }

// The standard implementation of a Meter uses a goroutine to synchronize
// its calculations and another goroutine (via time.Ticker) to produce
// clock ticks.
type StandardMeter struct {
	in     chan int64
	out    chan meterV
	ticker *time.Ticker
}

// Return the count of events seen.
func (m *StandardMeter) Count() int64 {
	return (<-m.out).count
}

// Mark the occurance of n events.
func (m *StandardMeter) Mark(n int64) {
	m.in <- n
}

// Return the meter's one-minute moving average rate of events.
func (m *StandardMeter) Rate1() float64 {
	return (<-m.out).rate1
}

// Return the meter's five-minute moving average rate of events.
func (m *StandardMeter) Rate5() float64 {
	return (<-m.out).rate5
}

// Return the meter's fifteen-minute moving average rate of events.
func (m *StandardMeter) Rate15() float64 {
	return (<-m.out).rate15
}

// Return the meter's mean rate of events.
func (m *StandardMeter) RateMean() float64 {
	return (<-m.out).rateMean
}

// Receive inputs and send outputs.  Count each input and update the various
// moving averages and the mean rate of events.  Send a copy of the meterV
// as output.
func (m *StandardMeter) arbiter() {
	var mv meterV
	a1 := NewEWMA1()
	a5 := NewEWMA5()
	a15 := NewEWMA15()
	t := time.Now()
	for {
		select {
		case n := <-m.in:
			mv.count += n
			a1.Update(n)
			a5.Update(n)
			a15.Update(n)
			mv.rate1 = a1.Rate()
			mv.rate5 = a5.Rate()
			mv.rate15 = a15.Rate()
			mv.rateMean = float64(1e9*mv.count) / float64(time.Since(t))
		case m.out <- mv:
		case <-m.ticker.C:
			a1.Tick()
			a5.Tick()
			a15.Tick()
			mv.rate1 = a1.Rate()
			mv.rate5 = a5.Rate()
			mv.rate15 = a15.Rate()
			mv.rateMean = float64(1e9*mv.count) / float64(time.Since(t))
		}
	}
}

// A meterV contains all the values that would need to be passed back
// from the synchronizing goroutine.
type meterV struct {
	count                          int64
	rate1, rate5, rate15, rateMean float64
}
