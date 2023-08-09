package metrics

import "testing"

func BenchmarkCounter(b *testing.B) {
	c := NewCounter()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Inc(1)
	}
}

func TestCounterClear(t *testing.T) {
	c := NewCounter()
	c.Inc(1)
	c.Clear()
	if count := c.Count(); 0 != count {
		t.Errorf("c.Count(): 0 != %v\n", count)
	}
}

func TestCounterDec1(t *testing.T) {
	c := NewCounter()
	c.Dec(1)
	if count := c.Count(); -1 != count {
		t.Errorf("c.Count(): -1 != %v\n", count)
	}
}

func TestCounterDec2(t *testing.T) {
	c := NewCounter()
	c.Dec(2)
	if count := c.Count(); -2 != count {
		t.Errorf("c.Count(): -2 != %v\n", count)
	}
}

func TestCounterInc1(t *testing.T) {
	c := NewCounter()
	c.Inc(1)
	if count := c.Count(); 1 != count {
		t.Errorf("c.Count(): 1 != %v\n", count)
	}
}

func TestCounterInc2(t *testing.T) {
	c := NewCounter()
	c.Inc(2)
	if count := c.Count(); 2 != count {
		t.Errorf("c.Count(): 2 != %v\n", count)
	}
}

func TestCounterSnapshot(t *testing.T) {
	c := NewCounter()
	c.Inc(1)
	snapshot := c.Snapshot()
	c.Inc(1)
	if count := snapshot.Count(); 1 != count {
		t.Errorf("c.Count(): 1 != %v\n", count)
	}
}

func TestCounterZero(t *testing.T) {
	c := NewCounter()
	if count := c.Count(); 0 != count {
		t.Errorf("c.Count(): 0 != %v\n", count)
	}
}

func TestGetOrRegisterCounter(t *testing.T) {
	r := NewRegistry()
	NewRegisteredCounter("foo", r).Inc(47)
	if c := GetOrRegisterCounter("foo", r); 47 != c.Count() {
		t.Fatal(c)
	}
}

func TestCounterLabels(t *testing.T) {
	labels := []Label{Label{"key1", "value1"}}
	c := NewCounter(labels...)
	if len(c.Labels()) != 1 {
		t.Fatalf("Labels(): %v != 1", len(c.Labels()))
	}
	if lbls := c.Labels()[0]; lbls.Key != "key1" || lbls.Value != "value1" {
		t.Errorf("Labels(): %v != key1; %v != value1", lbls.Key, lbls.Value)
	}

	// Labels passed by value.
	labels[0] = Label{"key3", "value3"}
	if lbls := c.Labels()[0]; lbls.Key != "key1" || lbls.Value != "value1" {
		t.Error("Labels(): labels passed by reference")
	}

	// Labels in snapshot.
	ss := c.Snapshot()
	if len(ss.Labels()) != 1 {
		t.Fatalf("Labels(): %v != 1", len(c.Labels()))
	}
	if lbls := ss.Labels()[0]; lbls.Key != "key1" || lbls.Value != "value1" {
		t.Errorf("Labels(): %v != key1; %v != value1", lbls.Key, lbls.Value)
	}
}

func TestCounterWithLabels(t *testing.T) {
	c := NewCounter(Label{"foo", "bar"})
	new := c.WithLabels(Label{"bar", "foo"})
	if len(new.Labels()) != 2 {
		t.Fatalf("WithLabels() len: %v != 2", len(new.Labels()))
	}
	if lbls:=new.Labels()[0]; lbls.Key != "foo" || lbls.Value != "bar" {
		t.Errorf("WithLabels(): %v != foo; %v != bar", lbls.Key, lbls.Value)
	}
	if lbls:=new.Labels()[1]; lbls.Key != "bar" || lbls.Value != "foo" {
		t.Errorf("WithLabels(): %v != bar; %v != foo", lbls.Key, lbls.Value)
	}
}
