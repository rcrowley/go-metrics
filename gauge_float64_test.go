package metrics

import "testing"

func BenchmarkGaugeFloat64(b *testing.B) {
	g := NewGaugeFloat64()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Update(float64(i))
	}
}

func TestGaugeFloat64(t *testing.T) {
	g := NewGaugeFloat64()
	g.Update(float64(47.0))
	if v := g.Value(); float64(47.0) != v {
		t.Errorf("g.Value(): 47.0 != %v\n", v)
	}
}

func TestGaugeFloat64Snapshot(t *testing.T) {
	g := NewGaugeFloat64()
	g.Update(float64(47.0))
	snapshot := g.Snapshot()
	g.Update(float64(0))
	if v := snapshot.Value(); float64(47.0) != v {
		t.Errorf("g.Value(): 47.0 != %v\n", v)
	}
}

func TestGetOrRegisterGaugeFloat64(t *testing.T) {
	r := NewRegistry()
	NewRegisteredGaugeFloat64("foo", r).Update(float64(47.0))
	t.Logf("registry: %v", r)
	if g := GetOrRegisterGaugeFloat64("foo", r); float64(47.0) != g.Value() {
		t.Fatal(g)
	}
}

func TestFuncGaugeFloat64(t *testing.T) {
	i := 0
	g := NewFuncGaugeFloat64(func() float64 {
		i++
		return float64(i)
	})
	for j := 1; j <= 2; j++ {
		if v := g.Value(); float64(j) != v {
			t.Errorf("g.Value(): %v != %v\n", j, v)
		}
	}
}

func TestFuncGaugeFloat64Update(t *testing.T) {
	g := NewFuncGaugeFloat64(func() float64 {
		return float64(47.0)
	})
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Update() did not panic")
		}
	}()
	g.Update(float64(48.0))
}

func TestFuncGaugeFloat64Snapshot(t *testing.T) {
	i := 0
	g := NewFuncGaugeFloat64(func() float64 {
		i++
		return float64(i)
	})
	snapshot := g.Snapshot()
	if v := g.Value(); float64(2.0) != v {
		t.Errorf("g.Value(): 2.0 != %v\n", v)
	}
	if v := snapshot.Value(); float64(1.0) != v {
		t.Errorf("g.Value(): 1.0 != %v\n", v)
	}
}

func TestGetOrRegisterFuncGaugeFloat64(t *testing.T) {
	r := NewRegistry()
	NewRegisteredFuncGaugeFloat64("foo", r, func() float64 {
		return float64(47.0)
	})
	if g := GetOrRegisterFuncGaugeFloat64("foo", r, func() float64 {
		return float64(48.0)
	}); float64(47.0) != g.Value() {
		t.Fatal(g)
	}
}
