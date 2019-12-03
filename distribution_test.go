package metrics

import "testing"

func BenchmarkDistribution(b *testing.B) {
	d := NewDistribution(LinearBuckets(0, 10, 100))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Update(int64(i))
	}
}

func TestLinearDistribution(t *testing.T) {
	d := NewDistribution(LinearBuckets(10, 10, 100))
	for i := 1; i <= 1000; i++ {
		d.Update(int64(i))
	}
	for i, v := range d.Buckets() {
		if v != int64(i) {
			t.Errorf("d.Buckets(): %v != %v\n", i, v)
		}
	}
}

func TestGetOrRegisterDistribution(t *testing.T) {
	buckets := LinearBuckets(0, 10, 1000)
	r := NewRegistry()
	NewRegisteredDistribution("foo", r, buckets).Update(47)
	if d := GetOrRegisterDistribution("foo", r, buckets); 1 != d.Count() {
		t.Fatal(d)
	}
}

func TestDistribution10000(t *testing.T) {
	d := NewDistribution(LinearBuckets(0, 10, 100))
	for i := 1; i <= 10000; i++ {
		d.Update(int64(i))
	}
	testDistribution10000(t, d)
}

func TestDistributionEmpty(t *testing.T) {
	d := NewDistribution(LinearBuckets(0, 10, 100))
	if count := d.Count(); 0 != count {
		t.Errorf("d.Count(): 0 != %v\n", count)
	}
	if min := d.Min(); 0 != min {
		t.Errorf("d.Min(): 0 != %v\n", min)
	}
	if max := d.Max(); 0 != max {
		t.Errorf("d.Max(): 0 != %v\n", max)
	}
	if mean := d.Mean(); 0.0 != mean {
		t.Errorf("d.Mean(): 0.0 != %v\n", mean)
	}
}

func TestDistributionSnapshot(t *testing.T) {
	d := NewDistribution(LinearBuckets(0, 10, 100))
	for i := 1; i <= 10000; i++ {
		d.Update(int64(i))
	}
	snapshot := d.Snapshot()
	d.Update(0)
	testDistribution10000(t, snapshot)
}

func testDistribution10000(t *testing.T, d Distribution) {
	if count := d.Count(); 10000 != count {
		t.Errorf("d.Count(): 10000 != %v\n", count)
	}
	if min := d.Min(); 1 != min {
		t.Errorf("d.Min(): 1 != %v\n", min)
	}
	if max := d.Max(); 10000 != max {
		t.Errorf("d.Max(): 10000 != %v\n", max)
	}
	if mean := d.Mean(); 5000.5 != mean {
		t.Errorf("d.Mean(): 5000.5 != %v\n", mean)
	}
}
