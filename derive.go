package metrics

import (
	"sync"
	"time"
)

// Derive count events to produce exponentially-weighted moving average rates
// at one-, five-, and fifteen-minutes and a mean rate. It expects always
// increasing values. The first call to Mark() initializes the base value. Any
// following calls update the value as difference to the base and update the
// base to the currently passed value. When the new value is smaller than the
// current one, the values are Clear()'d.
type Derive interface {
	Count() int64
	Base() int64
	Mark(int64)
	Rate1() float64
	Rate5() float64
	Rate15() float64
	RateMean() float64
	Snapshot() Derive
	Clear()
}

// UseNilDerive is set to true to use NilDerive type for a standard Derive
var UseNilDerive bool

// GetOrRegisterDerive returns an existing Derive or constructs and registers a
// new StandardDerive.
func GetOrRegisterDerive(name string, r Registry) Derive {
	if r == nil {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewDerive).(Derive)
}

// NewDerive constructs a new StandardDerive and launches a goroutine.
func NewDerive() Derive {
	if UseNilDerive {
		return NilDerive{}
	}
	m := newStandardDerive()
	derives.Lock()
	defer derives.Unlock()

	derives.meters = append(derives.meters, m)
	if !derives.started {
		derives.started = true
		go derives.tick()
	}
	return m
}

// NewRegisteredDerive constructs and registers a new StandardDerive
// and launches a goroutine.
func NewRegisteredDerive(name string, r Registry) Derive {
	c := NewDerive()
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// DeriveSnapshot is a read-only copy of another Derive.
type DeriveSnapshot struct {
	count, base                    int64
	rate1, rate5, rate15, rateMean float64
	initialized                    bool
}

// Count returns the count of events at the time the snapshot was taken.
func (m *DeriveSnapshot) Count() int64 { return m.count }

// Mark panics.
func (*DeriveSnapshot) Mark(n int64) {
	panic("Mark called on a DeriveSnapshot")
}

// Base returns the current base value for Mark()s
func (m *DeriveSnapshot) Base() int64 { return m.base }

// Rate1 returns the one-minute moving average rate of events per second at the
// time the snapshot was taken.
func (m *DeriveSnapshot) Rate1() float64 { return m.rate1 }

// Rate5 returns the five-minute moving average rate of events per second at
// the time the snapshot was taken.
func (m *DeriveSnapshot) Rate5() float64 { return m.rate5 }

// Rate15 returns the fifteen-minute moving average rate of events per second
// at the time the snapshot was taken.
func (m *DeriveSnapshot) Rate15() float64 { return m.rate15 }

// RateMean returns the meter's mean rate of events per second at the time the
// snapshot was taken.
func (m *DeriveSnapshot) RateMean() float64 { return m.rateMean }

// Snapshot returns the snapshot.
func (m *DeriveSnapshot) Snapshot() Derive { return m }

// Clear is a no-op.
func (m *DeriveSnapshot) Clear() {}

// NilDerive is a no-op Derive.
type NilDerive struct{}

// Base is a no-op.
func (NilDerive) Base() int64 { return 0 }

// Count is a no-op.
func (NilDerive) Count() int64 { return 0 }

// Mark is a no-op.
func (NilDerive) Mark(n int64) {}

// Rate1 is a no-op.
func (NilDerive) Rate1() float64 { return 0.0 }

// Rate5 is a no-op.
func (NilDerive) Rate5() float64 { return 0.0 }

// Rate15 is a no-op.
func (NilDerive) Rate15() float64 { return 0.0 }

// RateMean is a no-op.
func (NilDerive) RateMean() float64 { return 0.0 }

// Snapshot is a no-op.
func (NilDerive) Snapshot() Derive { return NilDerive{} }

// Clear is a no-op.
func (NilDerive) Clear() {}

// StandardDerive is the standard implementation of a Meter.
type StandardDerive struct {
	lock        sync.RWMutex
	snapshot    *DeriveSnapshot
	a1, a5, a15 EWMA
	startTime   time.Time
}

func newStandardDerive() *StandardDerive {
	return &StandardDerive{
		snapshot:  &DeriveSnapshot{},
		a1:        NewEWMA1(),
		a5:        NewEWMA5(),
		a15:       NewEWMA15(),
		startTime: time.Now(),
	}
}

// Clear clears the Derive
func (m *StandardDerive) Clear() {
	derives.Lock()
	defer derives.Unlock()
	m.clear()
}

func (m *StandardDerive) clear() {
	m.snapshot = &DeriveSnapshot{}
	m.a1 = NewEWMA1()
	m.a5 = NewEWMA5()
	m.a15 = NewEWMA15()
	m.startTime = time.Now()
}

// Base returns the current base value used for Mark()
func (m *StandardDerive) Base() int64 {
	m.lock.RLock()
	base := m.snapshot.base
	m.lock.RUnlock()
	return base
}

// Count returns the number of events recorded.
func (m *StandardDerive) Count() int64 {
	m.lock.RLock()
	count := m.snapshot.count
	m.lock.RUnlock()
	return count
}

// Mark records the occurance of n events.
func (m *StandardDerive) Mark(n int64) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if !m.snapshot.initialized {
		m.snapshot.base = n
		m.snapshot.initialized = true
		return
	}

	switch {
	case m.snapshot.base == n:
		// nothing happened
	case m.snapshot.base < n: // default case ;-)
		diff := n - m.snapshot.base
		m.snapshot.count += diff
		m.snapshot.base = n
		m.a1.Update(diff)
		m.a5.Update(diff)
		m.a15.Update(diff)
		m.updateSnapshot()
	default: // base > n: counter reset
		m.clear()
		m.snapshot.base = n
		m.snapshot.initialized = true
	}
}

// Rate1 returns the one-minute moving average rate of events per second.
func (m *StandardDerive) Rate1() float64 {
	m.lock.RLock()
	rate1 := m.snapshot.rate1
	m.lock.RUnlock()
	return rate1
}

// Rate5 returns the five-minute moving average rate of events per second.
func (m *StandardDerive) Rate5() float64 {
	m.lock.RLock()
	rate5 := m.snapshot.rate5
	m.lock.RUnlock()
	return rate5
}

// Rate15 returns the fifteen-minute moving average rate of events per second.
func (m *StandardDerive) Rate15() float64 {
	m.lock.RLock()
	rate15 := m.snapshot.rate15
	m.lock.RUnlock()
	return rate15
}

// RateMean returns the meter's mean rate of events per second.
func (m *StandardDerive) RateMean() float64 {
	m.lock.RLock()
	rateMean := m.snapshot.rateMean
	m.lock.RUnlock()
	return rateMean
}

// Snapshot returns a read-only copy of the meter.
func (m *StandardDerive) Snapshot() Derive {
	m.lock.RLock()
	snapshot := *m.snapshot
	m.lock.RUnlock()
	return &snapshot
}

func (m *StandardDerive) updateSnapshot() {
	// should run with write lock held on m.lock
	snapshot := m.snapshot
	snapshot.rate1 = m.a1.Rate()
	snapshot.rate5 = m.a5.Rate()
	snapshot.rate15 = m.a15.Rate()
	snapshot.rateMean = float64(snapshot.count) / time.Since(m.startTime).Seconds()
}

func (m *StandardDerive) tick() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.a1.Tick()
	m.a5.Tick()
	m.a15.Tick()
	m.updateSnapshot()
}

type deriveArbiter struct {
	sync.RWMutex
	started bool
	meters  []*StandardDerive
	ticker  *time.Ticker
}

var derives = deriveArbiter{ticker: time.NewTicker(5e9)}

// Ticks meters on the scheduled interval
func (ma *deriveArbiter) tick() {
	for {
		select {
		case <-ma.ticker.C:
			ma.tickMeters()
		}
	}
}

func (ma *deriveArbiter) tickMeters() {
	ma.RLock()
	defer ma.RUnlock()
	for _, meter := range ma.meters {
		meter.tick()
	}
}
