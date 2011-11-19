package metrics

import (
	"math"
)

type Histogram interface {
	Clear()
	Count() int64
	Max() int64
	Mean() float64
	Min() int64
	Percentile(float64) float64
	StdDev() float64
	Sum() int64
	Update(int64)
	Variance() float64
}

type histogram struct {
	in chan int64
	out chan histogramV
	reset chan bool
}

type histogramV struct {
	count, sum, min, max int64
	variance [2]float64
}

func NewHistogram() Histogram {
	h := &histogram{make(chan int64), make(chan histogramV), make(chan bool)}
	go h.arbiter()
	return h
}

func (h *histogram) Clear() {
	h.reset <- true
}

func (h *histogram) Count() int64 {
	return (<-h.out).count
}

func (h *histogram) Max() int64 {
	hv := <-h.out
	if 0 < hv.count { return hv.max }
	return 0
}

func (h *histogram) Mean() float64 {
	hv := <-h.out
	if 0 < hv.count {
		return float64(hv.sum) / float64(hv.count)
	}
	return 0
}

func (h *histogram) Min() int64 {
	hv := <-h.out
	if 0 < hv.count { return hv.min }
	return 0
}

func (h *histogram) Percentile(p float64) float64 {
	// This requires sampling, which is more involved than I have time to
	// implement this afternoon.
	return 0.0
}

func (h *histogram) StdDev() float64 {
	return math.Sqrt(h.Variance())
}

func (h *histogram) Sum() int64 {
	return (<-h.out).sum
}

func (h *histogram) Update(i int64) {
	h.in <- i
}

func (h *histogram) Variance() float64 {
	hv := <-h.out
	if 1 >= hv.count { return 0.0 }
	return hv.variance[1] / float64(hv.count - 1)
}

func (h *histogram) arbiter() {
	hv := newHistogramV()
	for {
		select {
		case i := <-h.in:
			hv.count++
			if i < hv.min { hv.min = i }
			if i > hv.max { hv.max = i }
			hv.sum += i
			f := float64(i)
			if -1.0 == hv.variance[0] {
				hv.variance[0] = f
				hv.variance[1] = 0.0
			} else {
				m := hv.variance[0]
				s := hv.variance[1]
				hv.variance[0] = m + (f - m) / float64(hv.count)
				hv.variance[1] = s + (f - m) * (f - hv.variance[0])
			}
		case h.out <- hv:
		case <- h.reset: hv = newHistogramV()
		}
	}
}

func newHistogramV() histogramV {
	return histogramV{
		0, 0, math.MaxInt64, math.MinInt64,
		[2]float64{-1.0, 0.0},
	}
}
