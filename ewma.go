package metrics

import (
	"math"
	"sync"
	"time"
)

// EWMAs calculate an exponentially-weighted moving average
type EWMA interface {
	Rate() float64
	Snapshot() EWMA
	Update(int64)
}

// NewEWMA constructs a new EWMA with the given alpha.
func NewEWMA(alpha float64) EWMA {
	if UseNilMetrics {
		return NilEWMA{}
	}
	return &StandardEWMA{alpha: alpha}
}

// NewEWMA1 constructs a new EWMA for a one-minute moving average.
func NewEWMA1() EWMA {
	return NewEWMA(1 - math.Exp(-1.0/60.0/1))
}

// NewEWMA5 constructs a new EWMA for a five-minute moving average.
func NewEWMA5() EWMA {
	return NewEWMA(1 - math.Exp(-1.0/60.0/5))
}

// NewEWMA15 constructs a new EWMA for a fifteen-minute moving average.
func NewEWMA15() EWMA {
	return NewEWMA(1 - math.Exp(-1.0/60.0/15))
}

// EWMASnapshot is a read-only copy of another EWMA.
type EWMASnapshot float64

// Rate returns the rate of events per second at the time the snapshot was
// taken.
func (a EWMASnapshot) Rate() float64 { return float64(a) }

// Snapshot returns the snapshot.
func (a EWMASnapshot) Snapshot() EWMA { return a }

// Updating a snapshot is a no-op.
func (EWMASnapshot) Update(int64) {}

// NilEWMA is a no-op EWMA.
type NilEWMA struct{}

// Rate is a no-op.
func (NilEWMA) Rate() float64 { return 0.0 }

// Snapshot is a no-op.
func (NilEWMA) Snapshot() EWMA { return NilEWMA{} }

// Tick is a no-op.
func (NilEWMA) Tick() {}

// Update is a no-op.
func (NilEWMA) Update(n int64) {}

// StandardEWMA is the standard implementation of an EWMA.
type StandardEWMA struct {
	alpha     float64
	ewma      float64
	uncounted int64
	timestamp int64
	init      bool
	mutex     sync.Mutex
}

// Rate returns the moving average rate of events per second.
func (a *StandardEWMA) Rate() float64 {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if !a.init {
		return 0
	}

	now := time.Now().UnixNano()
	elapsed := math.Floor(float64(now-a.timestamp) / 1e9)
	if elapsed >= 1 && a.uncounted != 0 {
		a.ewma = a.alpha*float64(a.uncounted) + (1-a.alpha)*a.ewma
		a.ewma = math.Pow(1-a.alpha, elapsed-1) * a.ewma
		a.timestamp = now
		a.uncounted = 0
		return a.ewma
	}

	a.ewma = math.Pow(1-a.alpha, float64(elapsed)) * a.ewma
	a.timestamp = now
	return a.ewma
}

// Snapshot returns a read-only copy of the EWMA.
func (a *StandardEWMA) Snapshot() EWMA {
	return EWMASnapshot(a.Rate())
}

// Used to elapse time in unit tests.
func (a *StandardEWMA) addToTimestamp(ns int64) {
	a.timestamp += ns
}

// Update registers n events that occured within the last ± 0.5 sec.
func (a *StandardEWMA) Update(n int64) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	now := time.Now().UnixNano()
	if !a.init {
		a.ewma = float64(n)
		a.timestamp = now
		a.init = true
		return
	}

	elapsed := math.Floor(float64(now-a.timestamp) / 1e9)
	if elapsed < 1 {
		a.uncounted += n
		return
	}

	a.ewma = a.alpha*float64(a.uncounted) + (1-a.alpha)*a.ewma
	a.ewma = math.Pow(1-a.alpha, elapsed-1) * a.ewma
	a.timestamp = now
	a.uncounted = n
}
