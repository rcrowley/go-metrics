package metrics

import (
	"math"
	"rand"
	"time"
)

const rescaleThreshold = 1e9 * 60 * 60

type Sample interface {
	Clear()
	Size() int
	Update(int64)
	Values() []int64
}

type expDecaySample struct {
	reservoirSize int
	alpha float64
	in chan int64
	out chan []int64
	reset chan bool
}

func NewExpDecaySample(reservoirSize int, alpha float64) Sample {
	s := &expDecaySample{
		reservoirSize,
		alpha,
		make(chan int64),
		make(chan []int64),
		make(chan bool),
	}
	go s.arbiter()
	return s
}

func (s *expDecaySample) Clear() {
	s.reset <- true
}

func (s *expDecaySample) Size() int {
	return len(<-s.out)
}

func (s *expDecaySample) Update(v int64) {
	s.in <- v
}

func (s *expDecaySample) Values() []int64 {
	return <-s.out
}

func (s *expDecaySample) arbiter() {
	count := 0
	values := make(map[float64]int64)
	tsStart := time.Nanoseconds()
	tsNext := tsStart + rescaleThreshold
	var valuesCopy []int64
	for {
		select {
		case v := <-s.in:
			ts := time.Nanoseconds()
			k := math.Exp(float64(ts - tsStart) * s.alpha) / rand.Float64()
			count++
			values[k] = v
			if count > s.reservoirSize {
				min := math.MaxFloat64
				for k, _ := range values {
					if k < min { min = k }
				}
				values[min] = 0, false
				valuesCopy = make([]int64, s.reservoirSize)
			} else {
				valuesCopy = make([]int64, count)
			}
			if ts > tsNext {
				tsOldStart := tsStart
				tsStart = time.Nanoseconds()
				tsNext = ts + rescaleThreshold
				oldValues := values
				values = make(map[float64]int64, len(oldValues))
				for k, v := range oldValues {
					values[k * math.Exp(-s.alpha * float64(
						tsStart - tsOldStart))] = v
				}
			}
			i := 0
			for _, v := range values {
				valuesCopy[i] = v
				i++
			}
		case s.out <- valuesCopy: // TODO Might need to make another copy here.
		case <-s.reset:
			count = 0
			values = make(map[float64]int64)
			valuesCopy = make([]int64, 0)
			tsStart = time.Nanoseconds()
			tsNext = tsStart + 1e9 * 60 * 60
		}
	}
}

type uniformSample struct {
	reservoirSize int
	in chan int64
	out chan []int64
	reset chan bool
}

func NewUniformSample(reservoirSize int) Sample {
	s := &uniformSample{
		reservoirSize,
		make(chan int64),
		make(chan []int64),
		make(chan bool),
	}
	go s.arbiter()
	return s
}

func (s *uniformSample) Clear() {
	s.reset <- true
}

func (s *uniformSample) Size() int {
	return len(<-s.out)
}

func (s *uniformSample) Update(v int64) {
	s.in <- v
}

func (s *uniformSample) Values() []int64 {
	return <-s.out
}

func (s *uniformSample) arbiter() {
	count := 0
	values := make([]int64, s.reservoirSize)
	var valuesCopy []int64
	for {
		select {
		case v := <-s.in:
			count++
			if count < s.reservoirSize {
				values[count - 1] = v
				valuesCopy = make([]int64, count)
			} else {
				values[rand.Intn(s.reservoirSize)] = v
				valuesCopy = make([]int64, len(values))
			}
			for i := 0; i < len(valuesCopy); i++ { valuesCopy[i] = values[i] }
		case s.out <- valuesCopy: // TODO Might need to make another copy here.
		case <-s.reset:
			count = 0
			values = make([]int64, s.reservoirSize)
			valuesCopy = make([]int64, 0)
		}
	}
}
