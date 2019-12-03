package metrics

// Histograms calculate distribution statistics from a series of int64 values.
type Distribution interface {
	Clear()
	Count() int64
	Max() int64
	Mean() float64
	Min() int64
	Sum() int64
	Update(int64)
	Snapshot() Distribution
	Buckets() map[float64]int64
}

// GetOrRegisterHistogram returns an existing Histogram or constructs and
// registers a new StandardHistogram.
func GetOrRegisterDistribution(name string, r Registry, buckets []float64) Distribution {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, func() Distribution { return NewDistribution(buckets) }).(Distribution)
}

// NewHistogram constructs a new StandardHistogram from a Sample.
func NewDistribution(buckets []float64) Distribution {
	if UseNilMetrics {
		return NilDistribution{}
	}

	return &StandardDistribution{
		buckets:      buckets,
		bucketsCount: make(map[float64]int64, len(buckets)),
	}
}

// NewRegisteredDistribution constructs and registers a new StandardDistribution from
// a Sample.
func NewRegisteredDistribution(name string, r Registry, buckets []float64) Distribution {
	c := NewDistribution(buckets)
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// StandardDistribution is the default implementation of a distribution
type StandardDistribution struct {
	sum          int64
	min          int64
	max          int64
	count        int64
	buckets      []float64
	bucketsCount map[float64]int64
}

func (s *StandardDistribution) Clear() {
	s.bucketsCount = make(map[float64]int64, len(s.buckets))
	s.sum = 0
	s.min = 0
	s.max = 0
	s.count = 0
}

func (s *StandardDistribution) Count() int64 {
	return s.count
}

func (s *StandardDistribution) Max() int64 {
	return s.max
}

func (s *StandardDistribution) Mean() float64 {
	if s.count == 0 {
		return 0
	}
	return float64(s.sum) / float64(s.count)
}

func (s *StandardDistribution) Min() int64 {
	return s.min
}

func (s *StandardDistribution) Sum() int64 {
	return s.sum
}

func (s *StandardDistribution) Update(v int64) {
	for i := len(s.buckets) - 1; i > 0; i-- {
		bucket := s.buckets[i]
		if float64(v) <= bucket {
			s.bucketsCount[bucket]++
		} else {
			break
		}
	}
	s.sum += v
	s.count++
	if v < s.min || s.min == 0 {
		s.min = v
	}
	if v > s.max {
		s.max = v
	}
}

func (s *StandardDistribution) Snapshot() Distribution {
	return DistributionSnapshot{
		StandardDistribution: StandardDistribution{
			sum:          s.sum,
			min:          s.min,
			max:          s.max,
			count:        s.count,
			buckets:      s.buckets,
			bucketsCount: s.bucketsCount,
		},
	}
}

func (s *StandardDistribution) Buckets() map[float64]int64 {
	return s.bucketsCount
}

// Nil implementation of Distribution
type NilDistribution struct{}

func (n NilDistribution) Clear() {}

func (n NilDistribution) Count() int64 { return 0 }

func (n NilDistribution) Max() int64 { return 0 }

func (n NilDistribution) Mean() float64 { return 0 }

func (n NilDistribution) Min() int64 { return 0 }

func (n NilDistribution) Sum() int64 { return 0 }

func (n NilDistribution) Update(int64) {}

func (n NilDistribution) Snapshot() Distribution {
	return n
}

func (n NilDistribution) Buckets() map[float64]int64 {
	return make(map[float64]int64)
}

// Distribution snapshot
type DistributionSnapshot struct {
	StandardDistribution
}

func (d DistributionSnapshot) Clear() {
	panic("called clear on snapshot")
}

func (d DistributionSnapshot) Count() int64 {
	return d.count
}

func (d DistributionSnapshot) Max() int64 {
	return d.max
}

func (d DistributionSnapshot) Mean() float64 {
	return float64(d.sum) / float64(d.count)
}

func (d DistributionSnapshot) Min() int64 {
	return d.min
}

func (d DistributionSnapshot) Sum() int64 {
	return d.sum
}

func (d DistributionSnapshot) Update(int64) {
	panic("called update on snapshot")
}

func (d DistributionSnapshot) Snapshot() Distribution {
	return d
}

func (d DistributionSnapshot) Buckets() map[float64]int64 {
	return d.bucketsCount
}

// Utility functions

// From github.com/prometheus/client_golang/prometheus/histogram
func LinearBuckets(start, width float64, count int) []float64 {
	if count < 1 {
		panic("LinearBuckets needs a positive count")
	}
	buckets := make([]float64, count)
	for i := range buckets {
		buckets[i] = start
		start += width
	}
	return buckets
}

// From github.com/prometheus/client_golang/prometheus/histogram
func ExponentialBuckets(start, factor float64, count int) []float64 {
	if count < 1 {
		panic("ExponentialBuckets needs a positive count")
	}
	if start <= 0 {
		panic("ExponentialBuckets needs a positive start value")
	}
	if factor <= 1 {
		panic("ExponentialBuckets needs a factor greater than 1")
	}
	buckets := make([]float64, count)
	for i := range buckets {
		buckets[i] = start
		start *= factor
	}
	return buckets
}
