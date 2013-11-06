package metrics

import (
	"math"
	"sort"
)

type Snapshot struct {
	sortedValues []int64
}

func NewSnapshot(sample Sample) *Snapshot {
	sortedValues := int64Slice(sample.Values())
	sort.Sort(sortedValues)

	return &Snapshot{sortedValues: sortedValues}
}

func (snapshot *Snapshot) isEmpty() bool {
	return len(snapshot.sortedValues) == 0
}

func (snapshot *Snapshot) SortedValues() []int64 {
	return snapshot.sortedValues
}

func (snapshot *Snapshot) Max() int64 {
	if snapshot.isEmpty() {
		return 0
	}

	return snapshot.sortedValues[len(snapshot.sortedValues)-1]
}

func (snapshot *Snapshot) Min() int64 {
	if snapshot.isEmpty() {
		return 0
	}

	return snapshot.sortedValues[0]
}

func (snapshot *Snapshot) Count() int64 {
	return int64(len(snapshot.sortedValues))
}

func (snapshot *Snapshot) Sum() (sum int64) {
	for _, value := range snapshot.sortedValues {
		sum += value
	}
	return sum
}

func (snapshot *Snapshot) Mean() float64 {
	if snapshot.isEmpty() {
		return 0
	}

	return float64(snapshot.Sum()) / float64(snapshot.Count())
}

func (snapshot *Snapshot) Variance() float64 {
	mean := snapshot.Mean()

	sum := 0.0
	for _, value := range snapshot.sortedValues {
		diff := float64(value) - mean
		sum += diff * diff
	}
	return sum / float64(snapshot.Count()-1)
}

func (snapshot *Snapshot) Percentile(quantile float64) float64 {
	position := quantile * float64(snapshot.Count()+1)

	if position < 1.0 {
		return float64(snapshot.Min())
	}

	if position >= float64(snapshot.Count()) {
		return float64(snapshot.Max())
	}

	lower := float64(snapshot.sortedValues[int(position)-1])
	upper := float64(snapshot.sortedValues[int(position)])

	return lower + (position-math.Floor(position))*(upper-lower)
}

func (snapshot *Snapshot) StdDev() float64 {
	return math.Sqrt(snapshot.Variance())
}
