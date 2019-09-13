package metrics

type Snapshot interface {
	Count() int64
	Max() int64
	Mean() float64
	Min() int64
	Percentile(float64) float64
	Percentiles([]float64) []float64
	Size() int
	StdDev() float64
	Sum() int64
	Variance() float64
	Values() []int64
}

// SimpleSampleSnapshot is a read-only copy of another Sample.
type SimpleSampleSnapshot struct {
	count  int64
	values []int64
}

func NewSimpleSampleSnapshot(count int64, values []int64) *SimpleSampleSnapshot {
	return &SimpleSampleSnapshot{
		count:  count,
		values: values,
	}
}

// Clear panics.
func (*SimpleSampleSnapshot) Clear() {
	panic("Clear called on a SimpleSampleSnapshot")
}

// Count returns the count of inputs at the time the snapshot was taken.
func (s *SimpleSampleSnapshot) Count() int64 { return s.count }

// Max returns the maximal value at the time the snapshot was taken.
func (s *SimpleSampleSnapshot) Max() int64 { return SampleMax(s.values) }

// Mean returns the mean value at the time the snapshot was taken.
func (s *SimpleSampleSnapshot) Mean() float64 { return SampleMean(s.values) }

// Min returns the minimal value at the time the snapshot was taken.
func (s *SimpleSampleSnapshot) Min() int64 { return SampleMin(s.values) }

// Percentile returns an arbitrary percentile of values at the time the
// snapshot was taken.
func (s *SimpleSampleSnapshot) Percentile(p float64) float64 {
	return SamplePercentile(s.values, p)
}

// Percentiles returns a slice of arbitrary percentiles of values at the time
// the snapshot was taken.
func (s *SimpleSampleSnapshot) Percentiles(ps []float64) []float64 {
	return SamplePercentiles(s.values, ps)
}

// Size returns the size of the sample at the time the snapshot was taken.
func (s *SimpleSampleSnapshot) Size() int { return len(s.values) }

// StdDev returns the standard deviation of values at the time the snapshot was
// taken.
func (s *SimpleSampleSnapshot) StdDev() float64 { return SampleStdDev(s.values) }

// Sum returns the sum of values at the time the snapshot was taken.
func (s *SimpleSampleSnapshot) Sum() int64 { return SampleSum(s.values) }

// Values returns a copy of the values in the sample.
func (s *SimpleSampleSnapshot) Values() []int64 {
	values := make([]int64, len(s.values))
	copy(values, s.values)
	return values
}

// Variance returns the variance of values at the time the snapshot was taken.
func (s *SimpleSampleSnapshot) Variance() float64 { return SampleVariance(s.values) }
