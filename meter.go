package metrics

import (
	"sync/atomic"
	"time"
)

// Meters count events to produce exponentially-weighted moving average rates
// at one-, five-, and fifteen-minutes and a mean rate.
type Meter interface {
	Count() int64
	Mark(int64)
	Rate1() float64
	Rate5() float64
	Rate15() float64
	RateMean() float64
	Snapshot() Meter
	Labels() []Label
	WithLabels(...Label) Meter
}

// GetOrRegisterMeter returns an existing Meter or constructs and registers a
// new StandardMeter.
func GetOrRegisterMeter(name string, r Registry, labels ...Label) Meter {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, func() Meter {
		return NewMeter(labels...)
	}).(Meter)
}

// NewMeter constructs a new StandardMeter.
func NewMeter(labels ...Label) Meter {
	if UseNilMetrics {
		return NilMeter{}
	}
	return newStandardMeter(labels...)
}

// NewMeter constructs and registers a new StandardMeter.
func NewRegisteredMeter(name string, r Registry, labels ...Label) Meter {
	c := NewMeter(labels...)
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// MeterSnapshot is a read-only copy of another Meter.
type MeterSnapshot struct {
	count                          int64
	rate1, rate5, rate15, rateMean float64
	labels                         []Label
}

// Count returns the count of events at the time the snapshot was taken.
func (m *MeterSnapshot) Count() int64 { return m.count }

// Mark panics.
func (*MeterSnapshot) Mark(n int64) {
	panic("Mark called on a MeterSnapshot")
}

// Rate1 returns the one-minute moving average rate of events per second at the
// time the snapshot was taken.
func (m *MeterSnapshot) Rate1() float64 { return m.rate1 }

// Rate5 returns the five-minute moving average rate of events per second at
// the time the snapshot was taken.
func (m *MeterSnapshot) Rate5() float64 { return m.rate5 }

// Rate15 returns the fifteen-minute moving average rate of events per second
// at the time the snapshot was taken.
func (m *MeterSnapshot) Rate15() float64 { return m.rate15 }

// RateMean returns the meter's mean rate of events per second at the time the
// snapshot was taken.
func (m *MeterSnapshot) RateMean() float64 { return m.rateMean }

// Snapshot returns the snapshot.
func (m *MeterSnapshot) Snapshot() Meter { return m }

// Labels returns the snapshot's labels.
func (m *MeterSnapshot) Labels() []Label { return deepCopyLabels(m.labels) }

// WithLabels returns a copy of the snapshot with the given labels appended to
// the current list of labels.
func (m *MeterSnapshot) WithLabels(labels ...Label) Meter {
	return &MeterSnapshot{
		count:    m.Count(),
		rate1:    m.Rate1(),
		rate5:    m.Rate5(),
		rate15:   m.Rate15(),
		rateMean: m.RateMean(),
		labels:   append(m.Labels(), deepCopyLabels(labels)...),
	}
}

// NilMeter is a no-op Meter.
type NilMeter struct{}

// Count is a no-op.
func (NilMeter) Count() int64 { return 0 }

// Mark is a no-op.
func (NilMeter) Mark(n int64) {}

// Rate1 is a no-op.
func (NilMeter) Rate1() float64 { return 0.0 }

// Rate5 is a no-op.
func (NilMeter) Rate5() float64 { return 0.0 }

// Rate15is a no-op.
func (NilMeter) Rate15() float64 { return 0.0 }

// RateMean is a no-op.
func (NilMeter) RateMean() float64 { return 0.0 }

// Snapshot is a no-op.
func (NilMeter) Snapshot() Meter { return NilMeter{} }

// Labels is a no-op.
func (NilMeter) Labels() []Label { return []Label{} }

// WithLabels is a no-op.
func (NilMeter) WithLabels(...Label) Meter { return NilMeter{} }

// StandardMeter is the standard implementation of a Meter.
type StandardMeter struct {
	count       atomic.Int64
	a1, a5, a15 EWMA
	startTime   time.Time
	labels      []Label
}

func newStandardMeter(labels ...Label) *StandardMeter {
	return &StandardMeter{
		a1:        NewEWMA1(),
		a5:        NewEWMA5(),
		a15:       NewEWMA15(),
		startTime: time.Now(),
		labels:    deepCopyLabels(labels),
	}
}

// Count returns the number of events recorded.
func (m *StandardMeter) Count() int64 {
	return m.count.Load()
}

// Mark records the occurance of n events.
func (m *StandardMeter) Mark(n int64) {
	m.count.Add(n)
	m.a1.Update(n)
	m.a5.Update(n)
	m.a15.Update(n)
}

// Rate1 returns the one-minute moving average rate of events per second.
func (m *StandardMeter) Rate1() float64 {
	return m.a1.Rate()
}

// Rate5 returns the five-minute moving average rate of events per second.
func (m *StandardMeter) Rate5() float64 {
	return m.a5.Rate()
}

// Rate15 returns the fifteen-minute moving average rate of events per second.
func (m *StandardMeter) Rate15() float64 {
	return m.a15.Rate()
}

// RateMean returns the meter's mean rate of events per second.
func (m *StandardMeter) RateMean() float64 {
	return float64(m.Count()) / time.Since(m.startTime).Seconds()
}

// Snapshot returns a read-only copy of the meter.
func (m *StandardMeter) Snapshot() Meter {
	return &MeterSnapshot{
		count:    m.Count(),
		rate1:    m.Rate1(),
		rate5:    m.Rate5(),
		rate15:   m.Rate15(),
		rateMean: m.RateMean(),
		labels:   m.Labels(),
	}
}

// Labels returns a deep copy of the meter's labels.
func (m *StandardMeter) Labels() []Label { return deepCopyLabels(m.labels) }

// WithLabels returns a snapshot of the Meter with the given labels appended to
// the current list of labels.
func (m *StandardMeter) WithLabels(labels ...Label) Meter {
	return m.Snapshot().WithLabels(labels...)
}
