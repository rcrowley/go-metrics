package metrics

import "sync"

// GaugeFloat64s hold a float64 value that can be set arbitrarily.
type GaugeFloat64 interface {
	Snapshot() GaugeFloat64
	Update(float64)
	Value() float64
}

// GetOrRegisterGaugeFloat64 returns an existing GaugeFloat64 or constructs and registers a
// new StandardGaugeFloat64.
func GetOrRegisterGaugeFloat64(name string, r Registry) GaugeFloat64 {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewGaugeFloat64()).(GaugeFloat64)
}

// NewGaugeFloat64 constructs a new StandardGaugeFloat64.
func NewGaugeFloat64() GaugeFloat64 {
	if UseNilMetrics {
		return NilGaugeFloat64{}
	}
	return &StandardGaugeFloat64{
		value: 0.0,
	}
}

// NewRegisteredGaugeFloat64 constructs and registers a new StandardGaugeFloat64.
func NewRegisteredGaugeFloat64(name string, r Registry) GaugeFloat64 {
	c := NewGaugeFloat64()
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// GetOrRegisterFuncGaugeFloat64 returns an existing Gauge or constructs and registers a
// new FuncGaugeFloat64.
func GetOrRegisterFuncGaugeFloat64(name string, r Registry, f func() float64) GaugeFloat64 {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewFuncGaugeFloat64(f)).(GaugeFloat64)
}

// NewFuncGaugeFloat64 constructs a new FuncGaugeFloat64.
func NewFuncGaugeFloat64(f func() float64) GaugeFloat64 {
	if UseNilMetrics {
		return NilGaugeFloat64{}
	}
	return &FuncGaugeFloat64{f}
}

// NewRegisteredFuncGaugeFloat64 constructs and registers a new FuncGaugeFloat64.
func NewRegisteredFuncGaugeFloat64(name string, r Registry, f func() float64) GaugeFloat64 {
	c := NewFuncGaugeFloat64(f)
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// GaugeFloat64Snapshot is a read-only copy of another GaugeFloat64.
type GaugeFloat64Snapshot float64

// Snapshot returns the snapshot.
func (g GaugeFloat64Snapshot) Snapshot() GaugeFloat64 { return g }

// Update panics.
func (GaugeFloat64Snapshot) Update(float64) {
	panic("Update called on a GaugeFloat64Snapshot")
}

// Value returns the value at the time the snapshot was taken.
func (g GaugeFloat64Snapshot) Value() float64 { return float64(g) }

// NilGauge is a no-op Gauge.
type NilGaugeFloat64 struct{}

// Snapshot is a no-op.
func (NilGaugeFloat64) Snapshot() GaugeFloat64 { return NilGaugeFloat64{} }

// Update is a no-op.
func (NilGaugeFloat64) Update(v float64) {}

// Value is a no-op.
func (NilGaugeFloat64) Value() float64 { return 0.0 }

// StandardGaugeFloat64 is the standard implementation of a GaugeFloat64 and uses
// sync.Mutex to manage a single float64 value.
type StandardGaugeFloat64 struct {
	mutex sync.Mutex
	value float64
}

// Snapshot returns a read-only copy of the gauge.
func (g *StandardGaugeFloat64) Snapshot() GaugeFloat64 {
	return GaugeFloat64Snapshot(g.Value())
}

// Update updates the gauge's value.
func (g *StandardGaugeFloat64) Update(v float64) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.value = v
}

// Value returns the gauge's current value.
func (g *StandardGaugeFloat64) Value() float64 {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return g.value
}

// FuncGaugeFloat64 is the dynamic implementation of a Gauge and uses a
// function to compute the current float64 value.
type FuncGaugeFloat64 struct {
	f func() float64
}

// Snapshot returns a read-only copy of the gauge.
func (g *FuncGaugeFloat64) Snapshot() GaugeFloat64 {
	return GaugeFloat64Snapshot(g.Value())
}

// Update panics.
func (g *FuncGaugeFloat64) Update(v float64) {
	panic("Update called on a FuncGaugeFloat64")
}

// Invokes the underyling function to get the current value.
func (g *FuncGaugeFloat64) Value() float64 {
	return g.f()
}
