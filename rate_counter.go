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

func (this *StandardRateCounter) Clear() {
	this.lock.Lock()
	defer this.lock.Unlock()

	atomic.StoreInt64(&this.counter, 0)

	resetTimeMs := this.clock.Now().UnixNano() / 1e6
	for i, _ := range this.timestampsMs {
		this.timestampsMs[i] = resetTimeMs
		this.counts[i] = 0
	}

	this.lastSampleTimeMs = resetTimeMs
	this.lastRate = 0.0
	this.lastCount = 0

	// Head and tail never point to the same index.
	this.headIndex = 0
	this.tailIndex = len(this.timestampsMs) - 1
}

func (this *StandardRateCounter) Mark(n int64) {
	atomic.AddInt64(&this.counter, n)
	this.maybeSampleCount()
}

func (this *StandardRateCounter) Rate1() float64 {
	this.maybeSampleCount()
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.lastRate
}

func (this *StandardRateCounter) Count() int64 {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.counter
}

func (this *StandardRateCounter) Snapshot() RateCounter {
	this.maybeSampleCount()
	this.lock.RLock()
	defer this.lock.RUnlock()

	return &RateCounterSnapshot{
		count: this.counter,
		rate:  this.lastRate,
	}
}

func (this *StandardRateCounter) roundTime(timeMs int64) int64 {
	return timeMs - (timeMs % this.samplePeriodMs)
}

func (this *StandardRateCounter) advance(index int) int {
	return (index + 1) % len(this.counts)
}

/**
 * May sample the current count and timestamp.  Note that this is not an unbiased sampling
 * algorithm, but given that we are computing a rate over a ring buffer of 60 samples, it
 * should not matter in practice.
 */
func (this *StandardRateCounter) maybeSampleCount() {
	currentTimeMs := this.clock.Now().UnixNano() / 1e6
	currentSampleTimeMs := this.roundTime(currentTimeMs)

	this.lock.RLock()
	toSample := currentSampleTimeMs > this.lastSampleTimeMs
	this.lock.RUnlock()

	if !toSample {
		return
	}

	this.lock.Lock()
	defer this.lock.Unlock()

	if currentSampleTimeMs > this.lastSampleTimeMs {
		this.sampleCountAndUpdateRate(currentSampleTimeMs)
	}
}

/**
 * Records a new sample to the ring buffer, advances head and tail if necessary, and
 * recomputes the rate.
 */
func (this *StandardRateCounter) sampleCountAndUpdateRate(currentSampleTimeMs int64) {
	// Record newest up to date second sample time.  Clear rate.
	this.lastSampleTimeMs = currentSampleTimeMs

	// Advance head and write values.
	this.headIndex = this.advance(this.headIndex)
	this.timestampsMs[this.headIndex] = currentSampleTimeMs

	this.lastCount = atomic.LoadInt64(&this.counter)
	this.counts[this.headIndex] = this.lastCount

	// Ensure tail is always ahead of head.
	if this.tailIndex == this.headIndex {
		this.tailIndex = this.advance(this.tailIndex)
	}

	// Advance the 'tail' to the newest sample which is at least windowTimeMs old.
	for {
		nextWindowStart := this.advance(this.tailIndex)
		if nextWindowStart == this.headIndex ||
			this.timestampsMs[this.headIndex]-this.timestampsMs[nextWindowStart] < this.windowSizeMs {
			break
		}
		this.tailIndex = nextWindowStart
	}

	timeDeltaMs := this.timestampsMs[this.headIndex] - this.timestampsMs[this.tailIndex]
	if timeDeltaMs == 0 {
		this.lastRate = 0.0
	} else {
		if timeDeltaMs > this.windowSizeMs {
			timeDeltaMs = this.windowSizeMs
		}

		deltaTimeSecs := timeDeltaMs / 1000.0
		deltaCount := this.counts[this.headIndex] - this.counts[this.tailIndex]
		if deltaTimeSecs <= 0.0 {
			this.lastRate = 0
		} else {
			this.lastRate = float64(deltaCount) / float64(deltaTimeSecs)
		}
	}
}

type RateCounterSnapshot struct {
	rate  float64
	count int64
}

func (this *RateCounterSnapshot) Mark(n int64) {
	panic("Mark called on RateCounterSnapshot")
}

func (this *RateCounterSnapshot) Count() int64 {
	return this.count
}

func (this *RateCounterSnapshot) Rate1() float64 {
	return this.rate
}

func (this *RateCounterSnapshot) Clear() {
	panic("Clear called on RateCounterSnapshot")
}

func (this *RateCounterSnapshot) Snapshot() RateCounter {
	return this
}

type NilRateCounter struct{}

func (this NilRateCounter) Mark(int64) {
}

func (this NilRateCounter) Count() int64 {
	return 0
}

func (this NilRateCounter) Rate1() float64 {
	return 0
}

func (this NilRateCounter) Clear() {
}

func (this NilRateCounter) Snapshot() RateCounter {
	return NilRateCounter{}
}
