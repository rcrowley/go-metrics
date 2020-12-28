package metrics

import "testing"

func BenchmarkVectorVectorCounter(b *testing.B) {
	c := NewVectorCounter()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.IncAll(1)
	}
}

func TestVectorVectorCounterClear(t *testing.T) {
	c := NewVectorCounter()
	c.Inc("a", 1)
	c.Clear("a")
	if count := c.Get("a").Count(); 0 != count {
		t.Errorf("c.Count(): 0 != %v\n", count)
	}
}

func TestVectorCounterDec1(t *testing.T) {
	c := NewVectorCounter()
	c.Dec("a", 1)
	if count := c.Get("a").Count(); -1 != count {
		t.Errorf("c.Count(): -1 != %v\n", count)
	}
}

func TestVectorCounterDec2(t *testing.T) {
	c := NewVectorCounter()
	c.Dec("a", 2)
	if count := c.Get("a").Count(); -2 != count {
		t.Errorf("c.Count(): -2 != %v\n", count)
	}
}

func TestVectorCounterInc1(t *testing.T) {
	c := NewVectorCounter()
	c.Inc("a", 1)
	if count := c.Get("a").Count(); 1 != count {
		t.Errorf("c.Count(): 1 != %v\n", count)
	}
}

func TestVectorCounterInc2(t *testing.T) {
	c := NewVectorCounter()
	c.Inc("a",2)
	if count := c.Get("a").Count(); 2 != count {
		t.Errorf("c.Count(): 2 != %v\n", count)
	}
}

func TestVectorCounterSnapshot(t *testing.T) {
	c := NewVectorCounter()
	c.Inc("a",1)
	snapshot := c.Snapshot()
	c.Inc("a",1)
	if count := snapshot.Get("a").Count(); 1 != count {
		t.Errorf("c.Count(): 1 != %v\n", count)
	}
}

func TestVectorCounterZero(t *testing.T) {
	c := NewVectorCounter()
	if count := c.Get("a").Count(); 0 != count {
		t.Errorf("c.Count(): 0 != %v\n", count)
	}
}

func TestGetOrRegisterVectorCounter(t *testing.T) {
	r := NewRegistry()
	NewRegisteredVectorCounter("foo", r).Inc("a",47)
	if c := GetOrRegisterVectorCounter("foo", r).Get("a"); 47 != c.Count() {
		t.Fatal(c)
	}
}

func TestIncAllVectorCounter(t *testing.T) {
	r := NewRegistry()
	counter := NewRegisteredVectorCounter("foo", r)
	counter.Get("a")
	counter.Get("b")
	counter.Get("c")
	counter.Get("d")
	counter.IncAll(1)
	if counter.Get("a").Count() != 1 {
		t.Errorf("c.Count(): 1 != %v\n", counter.Get("a").Count())
	}
	if counter.Get("b").Count() != 1 {
		t.Errorf("c.Count(): 1 != %v\n", counter.Get("b").Count())
	}
	if counter.Get("c").Count() != 1 {
		t.Errorf("c.Count(): 1 != %v\n", counter.Get("c").Count())
	}
	if counter.Get("d").Count() != 1 {
		t.Errorf("c.Count(): 1 != %v\n", counter.Get("d").Count())
	}
}

func TestDecAllVectorCounter(t *testing.T) {
	r := NewRegistry()
	counter := NewRegisteredVectorCounter("foo", r)
	counter.Get("a")
	counter.Get("b")
	counter.Get("c")
	counter.Get("d")
	counter.IncAll(2)
	counter.DecAll(1)
	if counter.Get("a").Count() != 1 {
		t.Errorf("c.Count(): 1 != %v\n", counter.Get("a").Count())
	}
	if counter.Get("b").Count() != 1 {
		t.Errorf("c.Count(): 1 != %v\n", counter.Get("b").Count())
	}
	if counter.Get("c").Count() != 1 {
		t.Errorf("c.Count(): 1 != %v\n", counter.Get("c").Count())
	}
	if counter.Get("d").Count() != 1 {
		t.Errorf("c.Count(): 1 != %v\n", counter.Get("d").Count())
	}
}

func TestClearAllVectorCounter(t *testing.T) {
	r := NewRegistry()
	counter := NewRegisteredVectorCounter("foo", r)
	counter.Get("a")
	counter.Get("b")
	counter.Get("c")
	counter.Get("d")
	counter.IncAll(1)
	counter.ClearAll()
	if counter.Get("a").Count() != 0 {
		t.Errorf("c.Count(): 0 != %v\n", counter.Get("a").Count())
	}
	if counter.Get("b").Count() != 0 {
		t.Errorf("c.Count(): 0 != %v\n", counter.Get("b").Count())
	}
	if counter.Get("c").Count() != 0 {
		t.Errorf("c.Count(): 0 != %v\n", counter.Get("c").Count())
	}
	if counter.Get("d").Count() != 0 {
		t.Errorf("c.Count(): 0 != %v\n", counter.Get("d").Count())
	}
}