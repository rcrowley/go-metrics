package metrics

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func BenchmarkMeter(b *testing.B) {
	m := NewMeter()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Mark(1)
	}
}

func BenchmarkMeterParallel(b *testing.B) {
	m := NewMeter()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Mark(1)
		}
	})
}

// exercise race detector
func TestMeterConcurrency(t *testing.T) {
	m := newStandardMeter()
	wg := &sync.WaitGroup{}
	reps := 100
	for i := 0; i < reps; i++ {
		wg.Add(1)
		go func(m Meter, wg *sync.WaitGroup) {
			m.Mark(1)
			wg.Done()
		}(m, wg)

		// Test reading from EWMA concurrently.
		wg.Add(1)
		go func(m Meter, wg *sync.WaitGroup) {
			m.Snapshot()
			wg.Done()
		}(m, wg)
	}
	wg.Wait()
}

func TestGetOrRegisterMeter(t *testing.T) {
	r := NewRegistry()
	NewRegisteredMeter("foo", r).Mark(47)
	if m := GetOrRegisterMeter("foo", r); 47 != m.Count() {
		t.Fatal(m)
	}
}

func TestMeterDecay(t *testing.T) {
	m := newStandardMeter()
	m.Mark(1)
	rateMean := m.RateMean()
	time.Sleep(100 * time.Millisecond)
	if m.RateMean() >= rateMean {
		t.Error("m.RateMean() didn't decrease")
	}
}

func TestMeterNonzero(t *testing.T) {
	m := NewMeter()
	m.Mark(3)
	if count := m.Count(); 3 != count {
		t.Errorf("m.Count(): 3 != %v\n", count)
	}
}

func TestMeterSnapshot(t *testing.T) {
	rand.Seed(time.Now().Unix())
	m := NewMeter()
	m.Mark(rand.Int63())
	if snapshot := m.Snapshot(); m.Count() != snapshot.Count() {
		t.Fatal(snapshot)
	}
}

func TestMeterZero(t *testing.T) {
	m := NewMeter()
	if count := m.Count(); 0 != count {
		t.Errorf("m.Count(): 0 != %v\n", count)
	}
}

func TestMeterLabels(t *testing.T) {
	labels := []Label{Label{"key1", "value1"}}
	m := NewMeter(labels...)
	if len(m.Labels()) != 1 {
		t.Fatalf("Labels(): %v != 1", len(m.Labels()))
	}
	if lbls := m.Labels()[0]; lbls.Key != "key1" || lbls.Value != "value1" {
		t.Errorf("Labels(): %v != key1; %v != value1", lbls.Key, lbls.Value)
	}

	// Labels passed by value.
	labels[0] = Label{"key3", "value3"}
	if lbls := m.Labels()[0]; lbls.Key != "key1" || lbls.Value != "value1" {
		t.Error("Labels(): labels passed by reference")
	}

	// Labels in snapshot.
	ss := m.Snapshot()
	if len(ss.Labels()) != 1 {
		t.Fatalf("Labels(): %v != 1", len(m.Labels()))
	}
	if lbls := ss.Labels()[0]; lbls.Key != "key1" || lbls.Value != "value1" {
		t.Errorf("Labels(): %v != key1; %v != value1", lbls.Key, lbls.Value)
	}
}

func TestMeterWithLabels(t *testing.T) {
	m := NewMeter(Label{"foo", "bar"})
	new := m.WithLabels(Label{"bar", "foo"})
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
