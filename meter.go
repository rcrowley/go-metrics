package metrics

import (
	"time"
)

type Meter interface {
	Count() int64
	Mark(int64)
	Rate1() float64
	Rate5() float64
	Rate15() float64
	RateMean() float64
}

type meter struct {
	in chan int64
	out chan meterV
	reset chan bool
	ticker *time.Ticker
}

type meterV struct {
	count int64
	rate1, rate5, rate15, rateMean float64
}

func NewMeter() Meter {
	m := &meter{
		make(chan int64),
		make(chan meterV),
		make(chan bool),
		time.NewTicker(5e9),
	}
	go m.arbiter()
	return m
}

func (m *meter) Clear() {
	m.reset <- true
}

func (m *meter) Count() int64 {
	return (<-m.out).count
}

func (m *meter) Mark(n int64) {
	m.in <- n
}

func (m *meter) Rate1() float64 {
	return (<-m.out).rate1
}

func (m *meter) Rate5() float64 {
	return (<-m.out).rate5
}

func (m *meter) Rate15() float64 {
	return (<-m.out).rate15
}

func (m *meter) RateMean() float64 {
	return (<-m.out).rateMean
}

func (m *meter) arbiter() {
	var mv meterV
	a1 := NewEWMA1()
	a5 := NewEWMA5()
	a15 := NewEWMA15()
	tsStart := time.Nanoseconds()
	for {
		select {
		case n := <-m.in:
			mv.count += n
			a1.Update(n); mv.rate1 = a1.Rate()
			a5.Update(n); mv.rate5 = a5.Rate()
			a15.Update(n); mv.rate15 = a15.Rate()
			mv.rateMean = float64(1e9 * mv.count) / float64(
				time.Nanoseconds() - tsStart)
		case m.out <- mv:
		case <-m.reset:
			mv = meterV{}
			a1.Clear()
			a5.Clear()
			a15.Clear()
			tsStart = time.Nanoseconds()
		case <-m.ticker.C:
			a1.Tick()
			a5.Tick()
			a15.Tick()
		}
	}
}
