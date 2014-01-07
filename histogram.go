package metrics

import (
	"math"
	"sync"
	"sync/atomic"
)

// Histograms calculate distribution statistics from a series of int64 values.
type Histogram interface {
	Clear()
	Count() int64
	Max() int64
	Mean() float64
	Min() int64
	Percentile(float64) float64
	Percentiles([]float64) []float64
	Sample() Sample
	Snapshot() Histogram
	StdDev() float64
	Update(int64)
	Variance() float64
}

// GetOrRegisterHistogram returns an existing Histogram or constructs and
// registers a new StandardHistogram.
func GetOrRegisterHistogram(name string, r Registry, s Sample) Histogram {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewHistogram(s)).(Histogram)
}

// NewHistogram constructs a new StandardHistogram from a Sample.
func NewHistogram(s Sample) Histogram {
	if UseNilMetrics {
		return NilHistogram{}
	}
	return &StandardHistogram{
		max:        math.MinInt64,
		min:        math.MaxInt64,
		sample:     s,
		sampleMean: -1.0,
	}
}

// NewRegisteredHistogram constructs and registers a new StandardHistogram from
// a Sample.
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

// Clear is a no-op.
func (NilHistogram) Clear() {}

// Count is a no-op.
func (NilHistogram) Count() int64 { return 0 }

// Max is a no-op.
func (NilHistogram) Max() int64 { return 0 }

// Mean is a no-op.
func (NilHistogram) Mean() float64 { return 0.0 }

// Min is a no-op.
func (NilHistogram) Min() int64 { return 0 }

// Percentile is a no-op.
func (NilHistogram) Percentile(p float64) float64 { return 0.0 }

// Percentiles is a no-op.
func (NilHistogram) Percentiles(ps []float64) []float64 {
	return make([]float64, len(ps))
}

// Sample is a no-op.
func (NilHistogram) Sample() Sample { return NilSample{} }

// No-op.
func (NilHistogram) StdDev() float64 { return 0.0 }

// StdDev is a no-op.
func (NilHistogram) StdDev() float64 { return 0.0 }

// Update is a no-op.
func (NilHistogram) Update(v int64) {}

// Variance is a no-op.
func (NilHistogram) Variance() float64 { return 0.0 }

// StandardHistogram is the standard implementation of a Histogram and uses a
// Sample to bound its memory use.
type StandardHistogram struct {
	count, max, min, sum int64
	mutex                sync.Mutex
	sample               Sample
	sampleMean           float64
	varianceNumerator    float64
}

// Clear clears the histogram and its sample.
func (h *StandardHistogram) Clear() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.count = 0
	h.max = math.MinInt64
	h.min = math.MaxInt64
	h.sample.Clear()
	h.sum = 0
	h.sampleMean = -1.0
	h.varianceNumerator = 0.0
}

// Count returns the count of events since the histogram was last cleared.
func (h *StandardHistogram) Count() int64 {
	return atomic.LoadInt64(&h.count)
}

// Max returns the maximum value seen since the histogram was last cleared.
func (h *StandardHistogram) Max() int64 {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if 0 == h.count {
		return 0
	}
	return h.max
}

// Mean returns the mean of all values seen since the histogram was last
// cleared.
func (h *StandardHistogram) Mean() float64 {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if 0 == h.count {
		return 0
	}
	return float64(h.sum) / float64(h.count)
}

// Min returns the minimum value seen since the histogram was last cleared.
func (h *StandardHistogram) Min() int64 {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if 0 == h.count {
		return 0
	}
	return h.min
}

// Percentile returns an arbitrary percentile of the values in the sample.
func (h *StandardHistogram) Percentile(p float64) float64 {
	return h.s.Percentile(p)
}

// Percentiles returns a slice of arbitrary percentiles of the values in the
// sample.
func (h *StandardHistogram) Percentiles(ps []float64) []float64 {
	return h.s.Percentiles(ps)
}

// Sample returns a copy of the Sample underlying the Histogram.
func (h *StandardHistogram) Sample() Sample {
	return h.s.Dup()
}

// StdDev returns the standard deviation of all values seen since the histogram
// was last cleared.
func (h *StandardHistogram) StdDev() float64 {
	return math.Sqrt(h.Variance())
}

// Update updates the histogram with a new value.
func (h *StandardHistogram) Update(v int64) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.sample.Update(v)
	h.count++
	if v < h.min {
		h.min = v
	}
	if v > h.max {
		h.max = v
	}
	h.sum += v
	fv := float64(v)
	if -1.0 == h.sampleMean {
		h.sampleMean = fv
		h.varianceNumerator = 0.0
	} else {
		sampleMean := h.sampleMean
		varianceNumerator := h.varianceNumerator
		h.sampleMean = sampleMean + (fv-sampleMean)/float64(h.count)
		h.varianceNumerator = varianceNumerator + (fv-sampleMean)*(fv-h.sampleMean)
	}
}

// Variance returns the variance of all values seen since the histogram was
// last cleared.
func (h *StandardHistogram) Variance() float64 {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	return h.variance()
}

// variance returns the variance of all the values in the sample but expects
// the lock to already be held.
func (h *StandardHistogram) variance() float64 {
	if 1 >= h.count {
		return 0.0
	}
	return h.varianceNumerator / float64(h.count-1)
}

type int64Slice []int64

func (p int64Slice) Len() int           { return len(p) }
func (p int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
