package metrics

import (
	"math"
	"rand"
	"time"
)

const rescaleThreshold = 1e9 * 60 * 60

// Samples maintain a statistically-significant selection of values from
// a stream.
//
// This is an interface so as to encourage other structs to implement
// the Sample API as appropriate.
type Sample interface {
	Clear()
	Size() int
	Update(int64)
	Values() []int64
}

// An exponentially-decaying sample using a forward-decaying priority
// reservoir.  See Cormode et al's "Forward Decay: A Practical Time Decay
// Model for Streaming Systems".
//
// <http://www.research.att.com/people/Cormode_Graham/library/publications/CormodeShkapenyukSrivastavaXu09.pdf>
type ExpDecaySample struct {
	reservoirSize int
	alpha float64
	in chan int64
	out chan []int64
	reset chan bool
}

// Create a new exponentially-decaying sample with the given reservoir size
// and alpha.
func NewExpDecaySample(reservoirSize int, alpha float64) *ExpDecaySample {
	s := &ExpDecaySample{
		reservoirSize,
		alpha,
		make(chan int64),
		make(chan []int64),
		make(chan bool),
	}
	go s.arbiter()
	return s
}

// Clear all samples.
func (s *ExpDecaySample) Clear() {
	s.reset <- true
}

// Return the size of the sample, which is at most the reservoir size.
func (s *ExpDecaySample) Size() int {
	return len(<-s.out)
}

// Update the sample with a new value.
func (s *ExpDecaySample) Update(v int64) {
	s.in <- v
}

// Return all the values in the sample.
func (s *ExpDecaySample) Values() []int64 {
	return <-s.out
}

// Receive inputs and send outputs.  Count and save each input value,
// rescaling the sample if enough time has elapsed since the last rescaling.
// Send a copy of the values as output.
func (s *ExpDecaySample) arbiter() {
	count := 0
	values := make(map[float64]int64)
	tsStart := time.Seconds()
	tsNext := time.Nanoseconds() + rescaleThreshold
	var valuesCopy []int64
	for {
		select {
		case v := <-s.in:
			ts := time.Seconds()
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
			tsNano := time.Nanoseconds()
			if tsNano > tsNext {
				tsOldStart := tsStart
				tsStart = time.Seconds()
				tsNext = tsNano + rescaleThreshold
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
			tsStart = time.Seconds()
			tsNext = tsStart + rescaleThreshold
		}
	}
}

// A uniform sample using Vitter's Algorithm R.
//
// <http://www.cs.umd.edu/~samir/498/vitter.pdf>
type UniformSample struct {
	reservoirSize int
	in chan int64
	out chan []int64
	reset chan bool
}

// Create a new uniform sample with the given reservoir size.
func NewUniformSample(reservoirSize int) *UniformSample {
	s := &UniformSample{
		reservoirSize,
		make(chan int64),
		make(chan []int64),
		make(chan bool),
	}
	go s.arbiter()
	return s
}

// Clear all samples.
func (s *UniformSample) Clear() {
	s.reset <- true
}

// Return the size of the sample, which is at most the reservoir size.
func (s *UniformSample) Size() int {
	return len(<-s.out)
}

// Update the sample with a new value.
func (s *UniformSample) Update(v int64) {
	s.in <- v
}

// Return all the values in the sample.
func (s *UniformSample) Values() []int64 {
	return <-s.out
}

// Receive inputs and send outputs.  Count and save each input value at a
// random index.  Send a copy of the values as output.
func (s *UniformSample) arbiter() {
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
