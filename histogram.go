package metrics

import (
	"math"
	"sort"
	"sync"
	"sync/atomic"
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

// Get an existing or create and register a new Histogram.
func GetOrRegisterHistogram(name string, r Registry, s Sample) Histogram {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewHistogram(s)).(Histogram)
}

// Create a new Histogram with the given Sample.  The initial values compare
// so that the first value will be both min and max and the variance is flagged
// for special treatment on its first iteration.
func NewHistogram(s Sample) Histogram {
	if UseNilMetrics {
		return NilHistogram{}
	}
	return &StandardHistogram{
		max:      math.MinInt64,
		min:      math.MaxInt64,
		s:        s,
		variance: [2]float64{-1.0, 0.0},
	}
}

// Create and register a new Histogram.
func NewRegisteredHistogram(name string, r Registry, s Sample) Histogram {
	c := NewHistogram(s)
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// No-op Histogram.
type NilHistogram struct{}

// No-op.
func (h NilHistogram) Clear() {}

// No-op.
func (h NilHistogram) Count() int64 { return 0 }

// No-op.
func (h NilHistogram) Max() int64 { return 0 }

// No-op.
func (h NilHistogram) Mean() float64 { return 0.0 }

// No-op.
func (h NilHistogram) Min() int64 { return 0 }

// No-op.
func (h NilHistogram) Percentile(p float64) float64 { return 0.0 }

// No-op.
func (h NilHistogram) Percentiles(ps []float64) []float64 {
	return make([]float64, len(ps))
}

// No-op.
func (h NilHistogram) StdDev() float64 { return 0.0 }

// No-op.
func (h NilHistogram) Update(v int64) {}

// No-op.
func (h NilHistogram) Variance() float64 { return 0.0 }

// The standard implementation of a Histogram uses a Sample and a goroutine
// to synchronize its calculations.
type StandardHistogram struct {
	count, sum, min, max int64
	mutex                sync.Mutex
	s                    Sample
	variance             [2]float64
}

// Clear the histogram.
func (h *StandardHistogram) Clear() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.count = 0
	h.max = math.MinInt64
	h.min = math.MaxInt64
	h.s.Clear()
	h.sum = 0
	h.variance = [...]float64{-1.0, 0.0}
}

// Return the count of inputs since the histogram was last cleared.
func (h *StandardHistogram) Count() int64 {
	return atomic.LoadInt64(&h.count)
}

// Return the maximal value seen since the histogram was last cleared.
func (h *StandardHistogram) Max() int64 {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if 0 == h.count {
		return 0
	}
	return h.max
}

// Return the mean of all values seen since the histogram was last cleared.
func (h *StandardHistogram) Mean() float64 {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if 0 == h.count {
		return 0
	}
	return float64(h.sum) / float64(h.count)
}

// Return the minimal value seen since the histogram was last cleared.
func (h *StandardHistogram) Min() int64 {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if 0 == h.count {
		return 0
	}
	return h.min
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
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.s.Update(v)
	h.count++
	if v < h.min {
		h.min = v
	}
	if v > h.max {
		h.max = v
	}
	h.sum += v
	fv := float64(v)
	if -1.0 == h.variance[0] {
		h.variance[0] = fv
		h.variance[1] = 0.0
	} else {
		m := h.variance[0]
		s := h.variance[1]
		h.variance[0] = m + (fv-m)/float64(h.count)
		h.variance[1] = s + (fv-m)*(fv-h.variance[0])
	}
}

// Return the variance of all values seen since the histogram was last cleared.
func (h *StandardHistogram) Variance() float64 {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if 1 >= h.count {
		return 0.0
	}
	return h.variance[1] / float64(h.count-1)
}

// Cribbed from the standard library's `sort` package.
type int64Slice []int64

func (p int64Slice) Len() int           { return len(p) }
func (p int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
