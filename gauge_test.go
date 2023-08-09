package metrics

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func BenchmarkGuage(b *testing.B) {
	g := NewGauge()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Update(int64(i))
	}
}

// exercise race detector
func TestGaugeConcurrency(t *testing.T) {
	rand.Seed(time.Now().Unix())
	g := NewGauge()
	wg := &sync.WaitGroup{}
	reps := 100
	for i := 0; i < reps; i++ {
		wg.Add(1)
		go func(g Gauge, wg *sync.WaitGroup) {
			g.Update(rand.Int63())
			wg.Done()
		}(g, wg)
	}
	wg.Wait()
}

func TestGauge(t *testing.T) {
	g := NewGauge()
	g.Update(int64(47))
	if v := g.Value(); 47 != v {
		t.Errorf("g.Value(): 47 != %v\n", v)
	}
}

func TestGaugeSnapshot(t *testing.T) {
	g := NewGauge()
	g.Update(int64(47))
	snapshot := g.Snapshot()
	g.Update(int64(0))
	if v := snapshot.Value(); 47 != v {
		t.Errorf("g.Value(): 47 != %v\n", v)
	}
}

func TestGetOrRegisterGauge(t *testing.T) {
	r := NewRegistry()
	NewRegisteredGauge("foo", r).Update(47)
	if g := GetOrRegisterGauge("foo", r); 47 != g.Value() {
		t.Fatal(g)
	}
}

func TestFunctionalGauge(t *testing.T) {
	var counter int64
	fg := NewFunctionalGauge(func() int64 {
		counter++
		return counter
	})
	fg.Value()
	fg.Value()
	if counter != 2 {
		t.Error("counter != 2")
	}
}

func TestGetOrRegisterFunctionalGauge(t *testing.T) {
	r := NewRegistry()
	NewRegisteredFunctionalGauge("foo", r, func() int64 { return 47 })
	if g := GetOrRegisterGauge("foo", r); 47 != g.Value() {
		t.Fatal(g)
	}
}

func TestGaugeLabels(t *testing.T) {
	labels := []Label{Label{"key1", "value1"}}
	g := NewGauge(labels...)
	if len(g.Labels()) != 1 {
		t.Fatalf("Labels(): %v != 1", len(g.Labels()))
	}
	if lbls := g.Labels()[0]; lbls.Key != "key1" || lbls.Value != "value1" {
		t.Errorf("Labels(): %v != key1; %v != value1", lbls.Key, lbls.Value)
	}

	// Labels passed by value.
	labels[0] = Label{"key3", "value3"}
	if lbls := g.Labels()[0]; lbls.Key != "key1" || lbls.Value != "value1" {
		t.Error("Labels(): labels passed by reference")
	}

	// Labels in snapshot.
	ss := g.Snapshot()
	if len(ss.Labels()) != 1 {
		t.Fatalf("Labels(): %v != 1", len(g.Labels()))
	}
	if lbls := ss.Labels()[0]; lbls.Key != "key1" || lbls.Value != "value1" {
		t.Errorf("Labels(): %v != key1; %v != value1", lbls.Key, lbls.Value)
	}
}

func TestGaugeWithLabels(t *testing.T) {
	g := NewGauge(Label{"foo", "bar"})
	new := g.WithLabels(Label{"bar", "foo"})
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
