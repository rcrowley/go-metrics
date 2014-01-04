package metrics

import (
	"testing"
)

func TestNewSnapshot(t *testing.T) {
	snapshot := createSnapshot([]int64{2, 3, 1})

	if !sameInt64Slide(snapshot.SortedValues(), []int64{1, 2, 3}) {
		t.Errorf("Wrong sortedValues:\n got: %v\nwant: %v", snapshot.SortedValues(), []int64{1, 2, 3})
	}
}

func TestSnapshot_Max(t *testing.T) {
	snapshot := createSnapshot([]int64{10, 0, -10})
	if snapshot.Max() != 10 {
		t.Errorf("Wrong Max value:\n got: %v\nwant: %v", snapshot.Max(), 10)
	}
}

func TestSnapshot_Min(t *testing.T) {
	snapshot := createSnapshot([]int64{10, 0, -10})
	if snapshot.Min() != -10 {
		t.Errorf("Wrong Min value:\n got: %v\nwant: %v", snapshot.Min(), -10)
	}
}

func TestSnapshot_Count(t *testing.T) {
	snapshot := createSnapshot([]int64{10, 0, -10})
	if snapshot.Count() != 3 {
		t.Errorf("Wrong Count value:\n got: %v\nwant: %v", snapshot.Count(), 3)
	}
}

func TestSnapshot_Sum(t *testing.T) {
	snapshot := createSnapshot([]int64{20, 0, -10})
	if snapshot.Sum() != 10 {
		t.Errorf("Wrong Sum value:\n got: %v\nwant: %v", snapshot.Sum(), 10)
	}
}

func TestSnapshot_Mean(t *testing.T) {
	snapshot := createSnapshot([]int64{20, 0, -5})
	if snapshot.Mean() != 5 {
		t.Errorf("Wrong Mean value:\n got: %v\nwant: %v", snapshot.Mean(), 5)
	}
}

func TestSnapshot_Variance(t *testing.T) {
	snapshot := createSnapshot([]int64{10, 0, -10})
	if snapshot.Variance() != 100 {
		t.Errorf("Wrong Variance value:\n got: %v\nwant: %v", snapshot.Variance(), 100)
	}
}

func TestSnapshot_StdDev(t *testing.T) {
	snapshot := createSnapshot([]int64{10, 0, -10})
	if snapshot.StdDev() != 10 {
		t.Errorf("Wrong StdDev value:\n got: %v\nwant: %v", snapshot.StdDev(), 10)
	}
}

func TestSnapshot_Percentile(t *testing.T) {
	snapshot := createSnapshot([]int64{1, 2, 3})

	var examples = []struct {
		quantile   float64
		percentile float64
	}{
		{0, 1},
		{0.25, 1},
		{0.30, 1.2},
		{0.50, 2},
		{0.60, 2.4},
		{0.75, 3},
		{1, 3},
	}

	for _, example := range examples {
		percentile := snapshot.Percentile(example.quantile)
		if percentile != example.percentile {
			t.Errorf("Wrong Percentile for %v:\n got: %v\nwant: %v", example.quantile, percentile, example.percentile)
		}
	}
}

func TestSnapshot_Empty(t *testing.T) {
	snapshot := createSnapshot([]int64{})

	if snapshot.Max() != 0 {
		t.Errorf("Wrong Max value:\n got: %v\nwant: %v", snapshot.Max(), 0)
	}
	if snapshot.Min() != 0 {
		t.Errorf("Wrong Min value:\n got: %v\nwant: %v", snapshot.Min(), 0)
	}
	if snapshot.Count() != 0 {
		t.Errorf("Wrong Count value:\n got: %v\nwant: %v", snapshot.Count(), 0)
	}
	if snapshot.Sum() != 0 {
		t.Errorf("Wrong Sum value:\n got: %v\nwant: %v", snapshot.Sum(), 0)
	}
	if snapshot.Mean() != 0 {
		t.Errorf("Wrong Mean value:\n got: %v\nwant: %v", snapshot.Mean(), 0)
	}
	if snapshot.Variance() != 0 {
		t.Errorf("Wrong Variance value:\n got: %v\nwant: %v", snapshot.Variance(), 0)
	}
	if snapshot.StdDev() != 0 {
		t.Errorf("Wrong StdDev value:\n got: %v\nwant: %v", snapshot.StdDev(), 0)
	}
	if snapshot.Percentile(0.5) != 0 {
		t.Errorf("Wrong Percentile value for 0.5:\n got: %v\nwant: %v", snapshot.Percentile(0.5), 0)
	}
}

func createSnapshot(values []int64) *Snapshot {
	sample := NewUniformSample(len(values))
	for _, value := range values {
		sample.Update(value)
	}

	return NewSnapshot(sample)
}

func sameInt64Slide(slice1 []int64, slice2 []int64) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	for index, value := range slice1 {
		if slice2[index] != value {
			return false
		}
	}

	return true
}
