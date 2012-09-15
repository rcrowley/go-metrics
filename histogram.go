package metrics

import (
	"math"
	"sort"
)

// Histograms calculate distribution statistics from an int64 value.
//
// This is an interface so as to encourage other structs to implement
// the Histogram API as appropriate.
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

// The standard implementation of a Histogram uses a Sample and a goroutine
// to synchronize its calculations.
type StandardHistogram struct {
	s     Sample
	in    chan int64
	out   chan histogramV
	reset chan bool
}

// A histogramV contains all the values that would need to be passed back
// from the synchronizing goroutine.
type histogramV struct {
	count, sum, min, max int64
	variance             [2]float64
}

// Create a new histogram with the given Sample.  Create the communication
// channels and start the synchronizing goroutine.
func NewHistogram(s Sample) *StandardHistogram {
	h := &StandardHistogram{
		s,
		make(chan int64),
		make(chan histogramV),
		make(chan bool),
	}
	go h.arbiter()
	return h
}

// Create a new histogramV.  The initial values compare so that the first
// value will be both min and max and the variance is flagged for special
// treatment on its first iteration.
func newHistogramV() histogramV {
	return histogramV{
		0, 0, math.MaxInt64, math.MinInt64,
		[2]float64{-1.0, 0.0},
	}
}

// Clear the histogram.
func (h *StandardHistogram) Clear() {
	h.reset <- true
}

// Return the count of inputs since the histogram was last cleared.
func (h *StandardHistogram) Count() int64 {
	return (<-h.out).count
}

// Return the maximal value seen since the histogram was last cleared.
func (h *StandardHistogram) Max() int64 {
	hv := <-h.out
	if 0 < hv.count {
		return hv.max
	}
	return 0
}

// Return the mean of all values seen since the histogram was last cleared.
func (h *StandardHistogram) Mean() float64 {
	hv := <-h.out
	if 0 < hv.count {
		return float64(hv.sum) / float64(hv.count)
	}
	return 0
}

// Return the minimal value seen since the histogram was last cleared.
func (h *StandardHistogram) Min() int64 {
	hv := <-h.out
	if 0 < hv.count {
		return hv.min
	}
	return 0
}

// Return an arbitrary percentile of all values seen since the histogram was
// last cleared.
func (h *StandardHistogram) Percentile(p float64) float64 {
	return h.Percentiles([]float64{p})[0]
}

// Return a slice of arbitrary percentiles of all values seen since the
// histogram was last cleared.
func (h *StandardHistogram) Percentiles(ps []float64) []float64 {
	scores := make([]float64, len(ps))
	values := int64Slice(h.s.Values())
	size := len(values)
	if size > 0 {
		sort.Sort(values)
		for i, p := range ps {
			pos := p * float64(size+1)
			if pos < 1.0 {
				scores[i] = float64(values[0])
			} else if pos >= float64(size) {
				scores[i] = float64(values[size-1])
			} else {
				lower := float64(values[int(pos)-1])
				upper := float64(values[int(pos)])
				scores[i] = lower + (pos-math.Floor(pos))*(upper-lower)
			}
		}
	}
	return scores
}

// Return the standard deviation of all values seen since the histogram was
// last cleared.
func (h *StandardHistogram) StdDev() float64 {
	return math.Sqrt(h.Variance())
}

// Update the histogram with a new value.
func (h *StandardHistogram) Update(v int64) {
	h.in <- v
}

// Return the variance of all values seen since the histogram was last cleared.
func (h *StandardHistogram) Variance() float64 {
	hv := <-h.out
	if 1 >= hv.count {
		return 0.0
	}
	return hv.variance[1] / float64(hv.count-1)
}

// Receive inputs and send outputs.  Sample each input and update values in
// the histogramV.  Send a copy of the histogramV as output.
func (h *StandardHistogram) arbiter() {
	hv := newHistogramV()
	for {
		select {
		case v := <-h.in:
			h.s.Update(v)
			hv.count++
			if v < hv.min {
				hv.min = v
			}
			if v > hv.max {
				hv.max = v
			}
			hv.sum += v
			fv := float64(v)
			if -1.0 == hv.variance[0] {
				hv.variance[0] = fv
				hv.variance[1] = 0.0
			} else {
				m := hv.variance[0]
				s := hv.variance[1]
				hv.variance[0] = m + (fv-m)/float64(hv.count)
				hv.variance[1] = s + (fv-m)*(fv-hv.variance[0])
			}
		case h.out <- hv:
		case <-h.reset:
			h.s.Clear()
			hv = newHistogramV()
		}
	}
}

// Cribbed from the standard library's `sort` package.
type int64Slice []int64

func (p int64Slice) Len() int           { return len(p) }
func (p int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
