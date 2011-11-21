package metrics

import (
	"math"
	"sort"
)

type Histogram interface {
	Clear()
	Count() int64
	Max() int64
	Mean() float64
	Min() int64
	Percentile(float64) float64
	Percentiles([]float64) []float64
	StdDev() float64
	Update(int64)
	Variance() float64
}

type histogram struct {
	s Sample
	in chan int64
	out chan histogramV
	reset chan bool
}

type histogramV struct {
	count, sum, min, max int64
	variance [2]float64
}

func NewHistogram(s Sample) Histogram {
	h := &histogram{
		s,
		make(chan int64),
		make(chan histogramV),
		make(chan bool),
	}
	go h.arbiter()
	return h
}

func newHistogramV() histogramV {
	return histogramV{
		0, 0, math.MaxInt64, math.MinInt64,
		[2]float64{-1.0, 0.0},
	}
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
	return h.Percentiles([]float64{p})[0]
}

func (h *histogram) Percentiles(ps []float64) []float64 {
	scores := make([]float64, len(ps))
	values := Int64Slice(h.s.Values())
	size := len(values)
	if size > 0 {
		sort.Sort(values)
		for i, p := range ps {
			pos := p * float64(size + 1)
			if pos < 1.0 {
				scores[i] = float64(values[0])
			} else if pos >= float64(size) {
				scores[i] = float64(values[size - 1])
			} else {
				lower := float64(values[int(pos) - 1])
				upper := float64(values[int(pos)])
				scores[i] = lower + (pos - math.Floor(pos)) * (upper - lower)
			}
		}
	}
	return scores
}

func (h *histogram) StdDev() float64 {
	return math.Sqrt(h.Variance())
}

func (h *histogram) Update(v int64) {
	h.in <- v
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
		case v := <-h.in:
			h.s.Update(v)
			hv.count++
			if v < hv.min { hv.min = v }
			if v > hv.max { hv.max = v }
			hv.sum += v
			fv := float64(v)
			if -1.0 == hv.variance[0] {
				hv.variance[0] = fv
				hv.variance[1] = 0.0
			} else {
				m := hv.variance[0]
				s := hv.variance[1]
				hv.variance[0] = m + (fv - m) / float64(hv.count)
				hv.variance[1] = s + (fv - m) * (fv - hv.variance[0])
			}
		case h.out <- hv:
		case <- h.reset:
			h.s.Clear()
			hv = newHistogramV()
		}
	}
}

// Cribbed from the standard library's `sort` package.
type Int64Slice []int64
func (p Int64Slice) Len() int           { return len(p) }
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
