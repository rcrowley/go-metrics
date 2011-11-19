package metrics

import (
	"rand"
)

type Sample interface {
	Clear()
	Count() int
	Size() int
	Update(int64)
	Values() []int64
}

type expDecaySample struct {
	reservoirSize int
	alpha float64
	count int
	values []int64
}

func NewExpDecaySample(reservoirSize int, alpha float64) Sample {
	return &expDecaySample{
		reservoirSize, alpha,
		0,
		make([]int64, reservoirSize),
	}
}

func (s *expDecaySample) Clear() {
}

func (s *expDecaySample) Count() int {
	return s.count
}

func (s *expDecaySample) Size() int {
	return 0
}

func (s *expDecaySample) Update(v int64) {
}

func (s *expDecaySample) Values() []int64 {
	return s.values // It might be worth copying this before returning it.
}

type uniformSample struct {
	reservoirSize int
	in chan int64
	out chan sampleV
	reset chan bool
}

func NewUniformSample(reservoirSize int) Sample {
	s := &uniformSample{
		reservoirSize,
		make(chan int64),
		make(chan sampleV),
		make(chan bool),
	}
	go s.arbiter()
	return s
}

func (s *uniformSample) Clear() {
	s.reset <- true
}

func (s *uniformSample) Count() int {
	return (<-s.out).count
}

func (s *uniformSample) Size() int {
	return (<-s.out).size()
}

func (s *uniformSample) Update(v int64) {
	s.in <- v
}

func (s *uniformSample) Values() []int64 {
	return (<-s.out).values
}

func (s *uniformSample) arbiter() {
	sv := newSampleV(s.reservoirSize)
	for {
		select {
		case v := <-s.in:
			sv.count++
			if sv.count < s.reservoirSize {
				sv.values[sv.count - 1] = v
			} else {
				sv.values[rand.Intn(s.reservoirSize)] = v
			}
		case s.out <- sv.dup():
		case <-s.reset:
			for i, _ := range sv.values { sv.values[i] = 0 }
		}
	}
}

type sampleV struct {
	count int
	values []int64
}

func newSampleV(reservoirSize int) sampleV {
	return sampleV{0, make([]int64, reservoirSize)}
}

func (sv sampleV) dup() sampleV {
	values := make([]int64, sv.size())
	for i := 0; i < sv.size(); i++ { values[i] = sv.values[i] }
	return sampleV{sv.count, values}
}

func (sv sampleV) size() int {
	if sv.count < len(sv.values) {
		return sv.count
	}
	return len(sv.values)
}
