package metrics

import (
	"math"
	"sync/atomic"
)

type EWMA interface {
	Rate() float64
	Tick()
	Update(int64)
}

type ewma struct {
	alpha float64
	uncounted int64
	rate float64
	initialized bool
	tick chan bool
}

func NewEWMA(alpha float64) EWMA {
	a := &ewma{alpha, 0, 0.0, false, make(chan bool)}
	go a.ticker()
	return a
}

func NewEWMA1() EWMA {
	return NewEWMA(1 - math.Exp(-5.0 / 60.0 / 1))
}

func NewEWMA5() EWMA {
	return NewEWMA(1 - math.Exp(-5.0 / 60.0 / 5))
}

func NewEWMA15() EWMA {
	return NewEWMA(1 - math.Exp(-5.0 / 60.0 / 15))
}

func (a *ewma) Rate() float64 {
	return a.rate * float64(1e9)
}

func (a *ewma) Tick() {
	a.tick <- true
}

func (a *ewma) Update(n int64) {
	atomic.AddInt64(&a.uncounted, n)
}

func (a *ewma) ticker() {
	for <-a.tick {
		count := a.uncounted
		atomic.AddInt64(&a.uncounted, -count)
		instantRate := float64(count) / float64(5e9)
		if a.initialized {
			a.rate += a.alpha * (instantRate - a.rate)
		} else {
			a.initialized = true
			a.rate = instantRate
		}
	}
}
