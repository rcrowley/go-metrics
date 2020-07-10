package metrics

import (
	"sync"
	"sync/atomic"

	"github.com/bountylabs/go-metrics/clock"
)

//go:generate counterfeiter . RateCounter
type RateCounter interface {
	Mark(int64)
	// Rate1() returns the rate up to the last full sampling period, so at time 5.3s it will only return the rate on [0, 5].
	// Count() returns the total count ever, including the current sampling period.
	Count() int64
	Rate1() float64
	Clear()
	Snapshot() RateCounter
}

// A port of https://cgit.twitter.biz/source/tree/src/java/com/twitter/search/common/metrics/SearchRateCounter.java

// A counter that tells you the rate per second that something happened during the past 60 seconds
// (excluding the most recent fractional second).
type StandardRateCounter struct {
	clock clock.Clock

	// lastCount "lags behind" by a sample period by design, this one really counts all events so far.
	// Note that lastCount is at par with lastRate, so that in MonViz rate(meter.lastCount) = meter.lastRate
	counter        int64
	samplePeriodMs int64
	windowSizeMs   int64

	lock sync.RWMutex

	// These values should only be used while holding the lock
	timestampsMs     []int64
	counts           []int64
	headIndex        int // Array index to most recent written value.
	tailIndex        int // Array index to oldest written value.
	lastSampleTimeMs int64
	lastRate         float64
	lastCount        int64
}

// GetOrRegisterRateCounter returns an existing RateCounter or constructs and registers a
// new StandardRateCounter.
func GetOrRegisterRateCounter(name string, r Registry) RateCounter {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, func() RateCounter { return NewStandardRateCounter(60, 1000, clock.New()) }).(RateCounter)
}

func NewRateCounter() RateCounter {
	if UseNilMetrics {
		return NilRateCounter{}
	}
	return NewStandardRateCounter(60, 1000, clock.New())
}

// NewRegisteredRateCounter constructs and registers a new StandardRateCounter.
func NewRegisteredRateCounter(name string, r Registry, clock clock.Clock) RateCounter {
	c := NewStandardRateCounter(60, 1000, clock)

	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

func NewStandardRateCounter(numSamples int64, samplePeriodMs int64, clock clock.Clock) RateCounter {
	rc := &StandardRateCounter{
		samplePeriodMs: samplePeriodMs,
		windowSizeMs:   numSamples * samplePeriodMs,
		timestampsMs:   make([]int64, numSamples+1),
		counts:         make([]int64, numSamples+1),
		clock:          clock,
	}

	rc.Clear()

	return rc
}

func (src *StandardRateCounter) Clear() {
	src.lock.Lock()
	defer src.lock.Unlock()

	atomic.StoreInt64(&src.counter, 0)

	resetTimeMs := src.clock.Now().UnixNano() / 1e6
	for i, _ := range src.timestampsMs {
		src.timestampsMs[i] = resetTimeMs
		src.counts[i] = 0
	}

	src.lastSampleTimeMs = resetTimeMs
	src.lastRate = 0.0
	src.lastCount = 0

	// Head and tail never point to the same index.
	src.headIndex = 0
	src.tailIndex = len(src.timestampsMs) - 1
}

func (src *StandardRateCounter) Mark(n int64) {
	atomic.AddInt64(&src.counter, n)
	src.maybeSampleCount()
}

func (src *StandardRateCounter) Rate1() float64 {
	src.maybeSampleCount()
	src.lock.RLock()
	defer src.lock.RUnlock()
	return src.lastRate
}

func (src *StandardRateCounter) Count() int64 {
	return atomic.LoadInt64(&src.counter)
}

func (src *StandardRateCounter) Snapshot() RateCounter {
	src.maybeSampleCount()
	src.lock.RLock()
	defer src.lock.RUnlock()

	return &RateCounterSnapshot{
		count: atomic.LoadInt64(&src.counter),
		rate:  src.lastRate,
	}
}

func (src *StandardRateCounter) roundTime(timeMs int64) int64 {
	return timeMs - (timeMs % src.samplePeriodMs)
}

func (src *StandardRateCounter) advance(index int) int {
	return (index + 1) % len(src.counts)
}

/**
 * May sample the current count and timestamp.  Note that this is not an unbiased sampling
 * algorithm, but given that we are computing a rate over a ring buffer of 60 samples, it
 * should not matter in practice.
 */
func (src *StandardRateCounter) maybeSampleCount() {
	currentTimeMs := src.clock.Now().UnixNano() / 1e6
	currentSampleTimeMs := src.roundTime(currentTimeMs)

	src.lock.RLock()
	toSample := currentSampleTimeMs > src.lastSampleTimeMs
	src.lock.RUnlock()

	if !toSample {
		return
	}

	src.lock.Lock()
	defer src.lock.Unlock()

	if currentSampleTimeMs > src.lastSampleTimeMs {
		src.sampleCountAndUpdateRate(currentSampleTimeMs)
	}
}

/**
 * Records a new sample to the ring buffer, advances head and tail if necessary, and
 * recomputes the rate.
 */
func (src *StandardRateCounter) sampleCountAndUpdateRate(currentSampleTimeMs int64) {
	// Record newest up to date second sample time.  Clear rate.
	src.lastSampleTimeMs = currentSampleTimeMs

	// Advance head and write values.
	src.headIndex = src.advance(src.headIndex)
	src.timestampsMs[src.headIndex] = currentSampleTimeMs

	src.lastCount = atomic.LoadInt64(&src.counter)
	src.counts[src.headIndex] = src.lastCount

	// Ensure tail is always ahead of head.
	if src.tailIndex == src.headIndex {
		src.tailIndex = src.advance(src.tailIndex)
	}

	// Advance the 'tail' to the newest sample which is at least windowTimeMs old.
	for {
		nextWindowStart := src.advance(src.tailIndex)
		if nextWindowStart == src.headIndex ||
			src.timestampsMs[src.headIndex]-src.timestampsMs[nextWindowStart] < src.windowSizeMs {
			break
		}
		src.tailIndex = nextWindowStart
	}

	timeDeltaMs := src.timestampsMs[src.headIndex] - src.timestampsMs[src.tailIndex]
	if timeDeltaMs == 0 {
		src.lastRate = 0.0
	} else {
		if timeDeltaMs > src.windowSizeMs {
			timeDeltaMs = src.windowSizeMs
		}

		deltaTimeSecs := timeDeltaMs / 1000.0
		deltaCount := src.counts[src.headIndex] - src.counts[src.tailIndex]
		if deltaTimeSecs <= 0.0 {
			src.lastRate = 0
		} else {
			src.lastRate = float64(deltaCount) / float64(deltaTimeSecs)
		}
	}
}

type RateCounterSnapshot struct {
	rate  float64
	count int64
}

func (rcSnapshot *RateCounterSnapshot) Mark(n int64) {
	panic("Mark called on RateCounterSnapshot")
}

func (rcSnapshot *RateCounterSnapshot) Count() int64 {
	return rcSnapshot.count
}

func (rcSnapshot *RateCounterSnapshot) Rate1() float64 {
	return rcSnapshot.rate
}

func (rcSnapshot *RateCounterSnapshot) Clear() {
	panic("Clear called on RateCounterSnapshot")
}

func (rcSnapshot *RateCounterSnapshot) Snapshot() RateCounter {
	return rcSnapshot
}

type NilRateCounter struct{}

func (NilRateCounter) Mark(int64) {
}

func (NilRateCounter) Count() int64 {
	return 0
}

func (NilRateCounter) Rate1() float64 {
	return 0
}

func (NilRateCounter) Clear() {
}

func (NilRateCounter) Snapshot() RateCounter {
	return NilRateCounter{}
}
