package metrics

import "sync/atomic"

// Gauges hold an int64 value that can be set arbitrarily.
type Gauge interface {
	Snapshot() Gauge
	Update(int64)
	Value() int64
	Labels() []Label
	WithLabels(...Label) Gauge
}

// GetOrRegisterGauge returns an existing Gauge or constructs and registers a
// new StandardGauge.
func GetOrRegisterGauge(name string, r Registry, labels ...Label) Gauge {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, func() Gauge {
		return NewGauge(labels...)
	}).(Gauge)
}

// NewGauge constructs a new StandardGauge.
func NewGauge(labels ...Label) Gauge {
	if UseNilMetrics {
		return NilGauge{}
	}
	return &StandardGauge{labels: deepCopyLabels(labels)}
}

// NewRegisteredGauge constructs and registers a new StandardGauge.
func NewRegisteredGauge(name string, r Registry, labels ...Label) Gauge {
	c := NewGauge(labels...)
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// NewFunctionalGauge constructs a new FunctionalGauge.
func NewFunctionalGauge(f func() int64, labels ...Label) Gauge {
	if UseNilMetrics {
		return NilGauge{}
	}
	return &FunctionalGauge{value: f, labels: deepCopyLabels(labels)}
}

// NewRegisteredFunctionalGauge constructs and registers a new StandardGauge.
func NewRegisteredFunctionalGauge(name string, r Registry, f func() int64) Gauge {
	c := NewFunctionalGauge(f)
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// GaugeSnapshot is a read-only copy of another Gauge.
type GaugeSnapshot struct {
	value  int64
	labels []Label
}

// Snapshot returns the snapshot.
func (g GaugeSnapshot) Snapshot() Gauge { return g }

// Update panics.
func (GaugeSnapshot) Update(int64) {
	panic("Update called on a GaugeSnapshot")
}

// Value returns the value at the time the snapshot was taken.
func (g GaugeSnapshot) Value() int64 { return g.value }

// Returns a deep copy of the gauge's labels.
func (g GaugeSnapshot) Labels() []Label { return deepCopyLabels(g.labels) }

// WithLabels returns a copy of the snapshot with the given labels appended to
// the current list of labels.
func (g GaugeSnapshot) WithLabels(labels ...Label) Gauge {
	return GaugeSnapshot{
		value:  g.Value(),
		labels: append(g.Labels(), deepCopyLabels(labels)...),
	}
}

// NilGauge is a no-op Gauge.
type NilGauge struct{}

// Snapshot is a no-op.
func (NilGauge) Snapshot() Gauge { return NilGauge{} }

// Update is a no-op.
func (NilGauge) Update(v int64) {}

// Value is a no-op.
func (NilGauge) Value() int64 { return 0 }

// Labels is a no-op.
func (NilGauge) Labels() []Label { return []Label{} }

// WithLabels is a no-op.
func (NilGauge) WithLabels(...Label) Gauge { return NilGauge{} }

// StandardGauge is the standard implementation of a Gauge and uses the
// sync/atomic package to manage a single int64 value.
type StandardGauge struct {
	value  atomic.Int64
	labels []Label
}

// Snapshot returns a read-only copy of the gauge.
func (g *StandardGauge) Snapshot() Gauge {
	return GaugeSnapshot{
		value:  g.Value(),
		labels: g.Labels(),
	}
}

// Update updates the gauge's value.
func (g *StandardGauge) Update(v int64) {
	g.value.Store(v)
}

// Value returns the gauge's current value.
func (g *StandardGauge) Value() int64 {
	return g.value.Load()
}

// Returns a copy of the gauge's labels.
func (g *StandardGauge) Labels() []Label { return deepCopyLabels(g.labels) }

// WithLabels returns a snapshot of the gauge with the given labels appended to
// the current list of labels.
func (g *StandardGauge) WithLabels(labels ...Label) Gauge {
	return g.Snapshot().WithLabels(labels...)
}

// FunctionalGauge returns value from given function
type FunctionalGauge struct {
	value  func() int64
	labels []Label
}

// Value returns the gauge's current value.
func (g FunctionalGauge) Value() int64 {
	return g.value()
}

// Snapshot returns the snapshot.
func (g FunctionalGauge) Snapshot() Gauge {
	return &GaugeSnapshot{
		value:  g.Value(),
		labels: g.Labels(),
	}
}

// Update panics.
func (FunctionalGauge) Update(int64) {
	panic("Update called on a FunctionalGauge")
}

// Returns a copy of the gauge's labels.
func (g FunctionalGauge) Labels() []Label { return deepCopyLabels(g.labels) }

// WithLabels returns a snapshot of the gauge with the given labels appended to
// the current list of labels.
func (g FunctionalGauge) WithLabels(labels ...Label) Gauge {
	return g.Snapshot().WithLabels(labels...)
}
