package metrics

import (
	"math"
)

type EWMA interface {
	Clear()
	Rate() float64
	Tick()
	Update(int64)
}

type ewma struct {
	alpha float64
	in chan int64
	out chan float64
	reset, tick chan bool
}

func NewEWMA(alpha float64) EWMA {
	a := &ewma{
		alpha,
		make(chan int64),
		make(chan float64),
		make(chan bool), make(chan bool),
	}
	go a.arbiter()
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

func (a *ewma) Clear() {
	a.reset <- true
}

func (a *ewma) Rate() float64 {
	return <-a.out * float64(1e9)
}

func (a *ewma) Tick() {
	a.tick <- true
}

func (a *ewma) Update(n int64) {
	a.in <- n
}

func (a *ewma) arbiter() {
	var uncounted int64
	var rate float64
	var initialized bool
	for {
		select {
		case n := <-a.in: uncounted += n
		case a.out <- rate:
		case <-a.reset:
			uncounted = 0
			rate = 0.0
		case <-a.tick:
			instantRate := float64(uncounted) / float64(5e9)
			if initialized {
				rate += a.alpha * (instantRate - rate)
			} else {
				initialized = true
				rate = instantRate
			}
			uncounted = 0
		}
	}
}
