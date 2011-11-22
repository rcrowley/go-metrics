package metrics

import "time"

// Timers capture the duration and rate of events.
//
// This is an interface so as to encourage other structs to implement
// the Histogram API as appropriate.
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
	Update(uint64)
}

// The standard implementation of a Timer uses a Histogram and Meter directly.
type timer struct {
	h Histogram
	m Meter
}

// Create a new timer with the given Histogram and Meter.
func NewCustomTimer(h Histogram, m Meter) Timer {
	return &timer{h, m}
}

// Create a new timer with a standard histogram and meter.  The histogram
// will use an exponentially-decaying sample with the same reservoir size
// and alpha as UNIX load averages.
func NewTimer() Timer {
	return &timer{NewHistogram(NewExpDecaySample(1028, 0.015)), NewMeter()}
}

// Return the count of inputs.
func (t *timer) Count() int64 {
	return t.h.Count()
}

// Return the maximal value seen.
func (t *timer) Max() int64 {
	return t.h.Max()
}

// Return the mean of all values seen.
func (t *timer) Mean() float64 {
	return t.h.Mean()
}

// Return the minimal value seen.
func (t *timer) Min() int64 {
	return t.h.Min()
}

// Return an arbitrary percentile of all values seen.
func (t *timer) Percentile(p float64) float64 {
	return t.h.Percentile(p)
}

// Return a slice of arbitrary percentiles of all values seen.
func (t *timer) Percentiles(ps []float64) []float64 {
	return t.h.Percentiles(ps)
}

// Return the meter's one-minute moving average rate of events.
func (t *timer) Rate1() float64 {
	return t.m.Rate1()
}

// Return the meter's five-minute moving average rate of events.
func (t *timer) Rate5() float64 {
	return t.m.Rate5()
}

// Return the meter's fifteen-minute moving average rate of events.
func (t *timer) Rate15() float64 {
	return t.m.Rate15()
}

// Return the meter's mean rate of events.
func (t *timer) RateMean() float64 {
	return t.m.RateMean()
}

// Return the standard deviation of all values seen.
func (t *timer) StdDev() float64 {
	return t.h.StdDev()
}

// Record the duration of the execution of the given function.
func (t *timer) Time(f func()) {
	ts := time.Nanoseconds()
	f()
	t.Update(uint64(time.Nanoseconds() - ts))
}

// Record the duration of an event.
func (t *timer) Update(duration uint64) {
	t.h.Update(int64(duration))
	t.m.Mark(1)
}
