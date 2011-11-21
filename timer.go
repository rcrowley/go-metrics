package metrics

import "time"

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

type timer struct {
	h Histogram
	m Meter
}

func NewCustomTimer(h Histogram, m Meter) Timer {
	return &timer{h, m}
}

func NewTimer() Timer {
	return &timer{NewHistogram(NewExpDecaySample(1028, 0.015)), NewMeter()}
}

func (t *timer) Count() int64 {
	return t.h.Count()
}

func (t *timer) Max() int64 {
	return t.h.Max()
}

func (t *timer) Mean() float64 {
	return t.h.Mean()
}

func (t *timer) Min() int64 {
	return t.h.Min()
}

func (t *timer) Percentile(p float64) float64 {
	return t.h.Percentile(p)
}

func (t *timer) Percentiles(ps []float64) []float64 {
	return t.h.Percentiles(ps)
}

func (t *timer) Rate1() float64 {
	return t.m.Rate1()
}

func (t *timer) Rate5() float64 {
	return t.m.Rate5()
}

func (t *timer) Rate15() float64 {
	return t.m.Rate15()
}

func (t *timer) RateMean() float64 {
	return t.m.RateMean()
}

func (t *timer) StdDev() float64 {
	return t.h.StdDev()
}

func (t *timer) Time(f func()) {
	ts := time.Nanoseconds()
	f()
	t.Update(uint64(time.Nanoseconds() - ts))
}

func (t *timer) Update(duration uint64) {
	t.h.Update(int64(duration))
	t.m.Mark(1)
}
