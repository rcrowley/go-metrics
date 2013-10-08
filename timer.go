package metrics

import "time"

// Timers capture the duration and rate of events.
//
// This is an interface so as to encourage other structs to implement
// the Timer API as appropriate.
type Timer interface {
	Count() int64
	Max() int64
	Mean() float64
	Min() int64
	Percentile(float64) float64
	Percentiles([]float64) []float64
	Rate1() float64
	Rate5() float64
	Rate15() float64
	RateMean() float64
	StdDev() float64
	Time(func())
	Update(time.Duration)
	UpdateSince(time.Time)
}

// Create a new timer with the given Histogram and Meter.
func NewCustomTimer(h Histogram, m Meter) Timer {
	if UseNilMetrics {
		return NilTimer{}
	}
	return &StandardTimer{h, m}
}

// Create a new timer with a standard histogram and meter.  The histogram
// will use an exponentially-decaying sample with the same reservoir size
// and alpha as UNIX load averages.
func NewTimer() Timer {
	if UseNilMetrics {
		return NilTimer{}
	}
	return &StandardTimer{
		NewHistogram(NewExpDecaySample(1028, 0.015)),
		NewMeter(),
	}
}

// No-op Timer.
type NilTimer struct {
	h Histogram
	m Meter
}

// No-op.
func (t NilTimer) Count() int64 { return 0 }

// No-op.
func (t NilTimer) Max() int64 { return 0 }

// No-op.
func (t NilTimer) Mean() float64 { return 0.0 }

// No-op.
func (t NilTimer) Min() int64 { return 0 }

// No-op.
func (t NilTimer) Percentile(p float64) float64 { return 0.0 }

// No-op.
func (t NilTimer) Percentiles(ps []float64) []float64 {
	return make([]float64, len(ps))
}

// No-op.
func (t NilTimer) Rate1() float64 { return 0.0 }

// No-op.
func (t NilTimer) Rate5() float64 { return 0.0 }

// No-op.
func (t NilTimer) Rate15() float64 { return 0.0 }

// No-op.
func (t NilTimer) RateMean() float64 { return 0.0 }

// No-op.
func (t NilTimer) StdDev() float64 { return 0.0 }

// No-op.
func (t NilTimer) Time(f func()) {}

// No-op.
func (t NilTimer) Update(d time.Duration) {}

// No-op.
func (t NilTimer) UpdateSince(ts time.Time) {}

// The standard implementation of a Timer uses a Histogram and Meter directly.
type StandardTimer struct {
	h Histogram
	m Meter
}

// Return the count of inputs.
func (t *StandardTimer) Count() int64 {
	return t.h.Count()
}

// Return the maximal value seen.
func (t *StandardTimer) Max() int64 {
	return t.h.Max()
}

// Return the mean of all values seen.
func (t *StandardTimer) Mean() float64 {
	return t.h.Mean()
}

// Return the minimal value seen.
func (t *StandardTimer) Min() int64 {
	return t.h.Min()
}

// Return an arbitrary percentile of all values seen.
func (t *StandardTimer) Percentile(p float64) float64 {
	return t.h.Percentile(p)
}

// Return a slice of arbitrary percentiles of all values seen.
func (t *StandardTimer) Percentiles(ps []float64) []float64 {
	return t.h.Percentiles(ps)
}

// Return the meter's one-minute moving average rate of events.
func (t *StandardTimer) Rate1() float64 {
	return t.m.Rate1()
}

// Return the meter's five-minute moving average rate of events.
func (t *StandardTimer) Rate5() float64 {
	return t.m.Rate5()
}

// Return the meter's fifteen-minute moving average rate of events.
func (t *StandardTimer) Rate15() float64 {
	return t.m.Rate15()
}

// Return the meter's mean rate of events.
func (t *StandardTimer) RateMean() float64 {
	return t.m.RateMean()
}

// Return the standard deviation of all values seen.
func (t *StandardTimer) StdDev() float64 {
	return t.h.StdDev()
}

// Record the duration of the execution of the given function.
func (t *StandardTimer) Time(f func()) {
	ts := time.Now()
	f()
	t.Update(time.Since(ts))
}

// Record the duration of an event.
func (t *StandardTimer) Update(d time.Duration) {
	t.h.Update(int64(d))
	t.m.Mark(1)
}

// Record the duration of an event that started at a time and ends now.
func (t *StandardTimer) UpdateSince(ts time.Time) {
	t.h.Update(int64(time.Since(ts)))
	t.m.Mark(1)
}
