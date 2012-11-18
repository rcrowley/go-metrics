package metrics

import (
	"container/heap"
	"math"
	"math/rand"
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
	alpha         float64
	in            chan int64
	out           chan chan []int64
	reset         chan bool
}

// Create a new exponentially-decaying sample with the given reservoir size
// and alpha.
func NewExpDecaySample(reservoirSize int, alpha float64) *ExpDecaySample {
	s := &ExpDecaySample{
		reservoirSize,
		alpha,
		make(chan int64),
		make(chan chan []int64),
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
	return len(s.Values())
}

// Update the sample with a new value.
func (s *ExpDecaySample) Update(v int64) {
	s.in <- v
}

// Return all the values in the sample.
func (s *ExpDecaySample) Values() []int64 {
	c := make(chan []int64)
	s.out <- c
	return <-c
}

// An individual sample.
type expDecaySample struct {
	k float64
	v int64
}

// A min-heap of samples.
type expDecaySampleHeap []expDecaySample

func (q expDecaySampleHeap) Len() int {
	return len(q)
}

func (q expDecaySampleHeap) Less(i, j int) bool {
	return q[i].k < q[j].k
}

func (q expDecaySampleHeap) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *expDecaySampleHeap) Push(x interface{}) {
	q_ := *q
	n := len(q_)
	q_ = q_[0 : n+1]
	q_[n] = x.(expDecaySample)
	*q = q_
}

func (q *expDecaySampleHeap) Pop() interface{} {
	q_ := *q
	n := len(q_)
	i := q_[n-1]
	q_ = q_[0 : n-1]
	*q = q_
	return i
}

// Receive inputs and send outputs.  Count and save each input value,
// rescaling the sample if enough time has elapsed since the last rescaling.
// Send a copy of the values as output.
func (s *ExpDecaySample) arbiter() {
	values := make(expDecaySampleHeap, 0, s.reservoirSize)
	start := time.Now()
	next := time.Now().Add(rescaleThreshold)
	for {
		select {
		case v := <-s.in:
			if len(values) == s.reservoirSize {
				heap.Pop(&values)
			}
			now := time.Now()
			k := math.Exp(now.Sub(start).Seconds()*s.alpha) / rand.Float64()
			heap.Push(&values, expDecaySample{k: k, v: v})
			if now.After(next) {
				oldValues := values
				oldStart := start
				values = make(expDecaySampleHeap, 0, s.reservoirSize)
				start = time.Now()
				next = start.Add(rescaleThreshold)
				for _, e := range oldValues {
					e.k = e.k * math.Exp(-s.alpha*float64(start.Sub(oldStart)))
					heap.Push(&values, e)
				}
			}
		case c := <-s.out:
			valuesCopy := make([]int64, len(values))
			for i, e := range values {
				valuesCopy[i] = e.v
			}
			c <- valuesCopy
		case <-s.reset:
			values = make(expDecaySampleHeap, 0, s.reservoirSize)
			start = time.Now()
			next = start.Add(rescaleThreshold)
		}
	}
}

// A uniform sample using Vitter's Algorithm R.
//
// <http://www.cs.umd.edu/~samir/498/vitter.pdf>
type UniformSample struct {
	reservoirSize int
	in            chan int64
	out           chan chan []int64
	reset         chan bool
}

// Create a new uniform sample with the given reservoir size.
func NewUniformSample(reservoirSize int) *UniformSample {
	s := &UniformSample{
		reservoirSize,
		make(chan int64),
		make(chan chan []int64),
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
	return len(s.Values())
}

// Update the sample with a new value.
func (s *UniformSample) Update(v int64) {
	s.in <- v
}

// Return all the values in the sample.
func (s *UniformSample) Values() []int64 {
	c := make(chan []int64)
	s.out <- c
	return <-c
}

// Receive inputs and send outputs.  Count and save each input value at a
// random index.  Send a copy of the values as output.
func (s *UniformSample) arbiter() {
	values := make([]int64, 0, s.reservoirSize)
	for {
		n := len(values)
		select {
		case v := <-s.in:
			if n < s.reservoirSize {
				values = values[0 : n+1]
				values[n] = v
			} else {
				values[rand.Intn(s.reservoirSize)] = v
			}
		case c := <-s.out:
			valuesCopy := make([]int64, n)
			for i := 0; i < n; i++ {
				valuesCopy[i] = values[i]
			}
			c <- valuesCopy
		case <-s.reset:
			values = make([]int64, 0, s.reservoirSize)
		}
	}
}
