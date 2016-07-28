package metrics

import "testing"

func BenchmarkGauge(b *testing.B) {
	g := NewGauge()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Update(int64(i))
	}
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

func TestFuncGauge(t *testing.T) {
	i := 0
	g := NewFuncGauge(func() int64 {
		i++
		return int64(i)
	})
	for j := 1; j <= 2; j++ {
		if v := g.Value(); int64(j) != v {
			t.Errorf("g.Value(): %v != %v\n", j, v)
		}
	}
}

func TestFuncGaugeUpdate(t *testing.T) {
	g := NewFuncGauge(func() int64 {
		return int64(47)
	})
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Update() did not panic")
		}
	}()
	g.Update(int64(48))
}

func TestFuncGaugeSnapshot(t *testing.T) {
	i := 0
	g := NewFuncGauge(func() int64 {
		i++
		return int64(i)
	})
	snapshot := g.Snapshot()
	if v := g.Value(); 2 != v {
		t.Errorf("g.Value(): 2 != %v\n", v)
	}
	if v := snapshot.Value(); 1 != v {
		t.Errorf("g.Value(): 1 != %v\n", v)
	}
}

func TestGetOrRegisterFuncGauge(t *testing.T) {
	r := NewRegistry()
	NewRegisteredFuncGauge("foo", r, func() int64 {
		return int64(47)
	})
	if g := GetOrRegisterFuncGauge("foo", r, func() int64 {
		return int64(48)
	}); 47 != g.Value() {
		t.Fatal(g)
	}
}
