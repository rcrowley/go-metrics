package metrics

import (
	"sync"
	"time"
)

// Timers capture the duration and rate of events.
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
	Variance() float64
}

// GetOrRegisterTimer returns an existing Timer or constructs and registers a
// new StandardTimer.
func GetOrRegisterTimer(name string, r Registry) Timer {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewTimer()).(Timer)
}

// NewCustomTimer constructs a new StandardTimer from a Histogram and a Meter.
func NewCustomTimer(h Histogram, m Meter) Timer {
	if UseNilMetrics {
		return NilTimer{}
	}
	return &StandardTimer{
		histogram: h,
		meter:     m,
	}
}

// NewRegisteredTimer constructs and registers a new StandardTimer.
func NewRegisteredTimer(name string, r Registry) Timer {
	c := NewTimer()
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// NewTimer constructs a new StandardTimer using an exponentially-decaying
// sample with the same reservoir size and alpha as UNIX load averages.
func NewTimer() Timer {
	if UseNilMetrics {
		return NilTimer{}
	}
	return &StandardTimer{
		histogram: NewHistogram(NewExpDecaySample(1028, 0.015)),
		meter:     NewMeter(),
	}
}

// NilTimer is a no-op Timer.
type NilTimer struct {
	h Histogram
	m Meter
}

// Count is a no-op.
func (NilTimer) Count() int64 { return 0 }

// Max is a no-op.
func (NilTimer) Max() int64 { return 0 }

// Mean is a no-op.
func (NilTimer) Mean() float64 { return 0.0 }

// Min is a no-op.
func (NilTimer) Min() int64 { return 0 }

// Percentile is a no-op.
func (NilTimer) Percentile(p float64) float64 { return 0.0 }

// Percentiles is a no-op.
func (NilTimer) Percentiles(ps []float64) []float64 {
	return make([]float64, len(ps))
}

// Rate1 is a no-op.
func (NilTimer) Rate1() float64 { return 0.0 }

// Rate5 is a no-op.
func (NilTimer) Rate5() float64 { return 0.0 }

// Rate15 is a no-op.
func (NilTimer) Rate15() float64 { return 0.0 }

// RateMean is a no-op.
func (NilTimer) RateMean() float64 { return 0.0 }

// No-op.
func (NilTimer) StdDev() float64 { return 0.0 }

// Time is a no-op.
func (NilTimer) Time(func()) {}

// Update is a no-op.
func (NilTimer) Update(time.Duration) {}

// UpdateSince is a no-op.
func (NilTimer) UpdateSince(time.Time) {}

// Variance is a no-op.
func (NilTimer) Variance() float64 { return 0.0 }

// StandardTimer is the standard implementation of a Timer and uses a Histogram
// and Meter.
type StandardTimer struct {
	histogram Histogram
	meter     Meter
	mutex     sync.Mutex
}

// Count returns the number of events recorded.
func (t *StandardTimer) Count() int64 {
	return t.histogram.Count()
}

// Max returns the maximum value in the sample.
func (t *StandardTimer) Max() int64 {
	return t.histogram.Max()
}

// Mean returns the mean of the values in the sample.
func (t *StandardTimer) Mean() float64 {
	return t.histogram.Mean()
}

// Min returns the minimum value in the sample.
func (t *StandardTimer) Min() int64 {
	return t.histogram.Min()
}

// Percentile returns an arbitrary percentile of the values in the sample.
func (t *StandardTimer) Percentile(p float64) float64 {
	return t.histogram.Percentile(p)
}

// Percentiles returns a slice of arbitrary percentiles of the values in the
// sample.
func (t *StandardTimer) Percentiles(ps []float64) []float64 {
	return t.histogram.Percentiles(ps)
}

// Rate1 returns the one-minute moving average rate of events per second.
func (t *StandardTimer) Rate1() float64 {
	return t.meter.Rate1()
}

// Rate5 returns the five-minute moving average rate of events per second.
func (t *StandardTimer) Rate5() float64 {
	return t.meter.Rate5()
}

// Rate15 returns the fifteen-minute moving average rate of events per second.
func (t *StandardTimer) Rate15() float64 {
	return t.meter.Rate15()
}

// RateMean returns the meter's mean rate of events per second.
func (t *StandardTimer) RateMean() float64 {
	return t.m.RateMean()
}

// StdDev returns the standard deviation of the values in the sample.
func (t *StandardTimer) StdDev() float64 {
	return t.histogram.StdDev()
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
