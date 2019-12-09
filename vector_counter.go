package metrics

import "sync"

// Counters holds multiple counters that can be Incremented and Decremented
type VectorCounter interface {
	DecAll(int64)
	IncAll(int64)
	GetAll() map[string]Counter
	ClearAll()

	Dec(label string, v int64)
	Inc(label string, v int64)
	Get(label string) Counter
	Clear(label string)

	Snapshot() VectorCounter
}

// GetOrRegisterCounter returns an existing Counter or constructs and registers
// a new StandardCounter.
func GetOrRegisterVectorCounter(name string, r Registry) VectorCounter {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewVectorCounter).(VectorCounter)
}

// NewCounter constructs a new StandardCounter.
func NewVectorCounter() VectorCounter {
	if UseNilMetrics {
		return &NilVectorCounter{}
	}
	return &StandardVectorCounter{vec: make(map[string]Counter)}
}

// NewRegisteredCounter constructs and registers a new StandardCounter.
func NewRegisteredVectorCounter(name string, r Registry) VectorCounter {
	c := NewVectorCounter()
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// Standard implementation

type StandardVectorCounter struct {
	vec map[string]Counter

	mtx sync.Mutex
}

func (vc *StandardVectorCounter) DecAll(v int64) {
	vc.mtx.Lock()
	for _, counter := range vc.vec {
		counter.Dec(v)
	}
	vc.mtx.Unlock()
}

func (vc *StandardVectorCounter) IncAll(v int64) {
	vc.mtx.Lock()
	for _, counter := range vc.vec {
		counter.Inc(v)
	}
	vc.mtx.Unlock()
}

func (vc *StandardVectorCounter) GetAll() map[string]Counter {
	return vc.vec
}

func (vc *StandardVectorCounter) ClearAll() {
	vc.mtx.Lock()
	for _, counter := range vc.vec {
		counter.Clear()
	}
	vc.mtx.Unlock()
}

func (vc *StandardVectorCounter) Dec(label string, v int64) {
	vc.getOrCreate(label).Dec(v)
}

func (vc *StandardVectorCounter) Inc(label string, v int64) {
	vc.getOrCreate(label).Inc(v)
}

func (vc *StandardVectorCounter) Get(label string) Counter {
	return vc.getOrCreate(label)
}

func (vc *StandardVectorCounter) Clear(label string) {
	vc.getOrCreate(label).Clear()
}

func (vc *StandardVectorCounter) Snapshot() VectorCounter {
	c := make(map[string]Counter, len(vc.vec))
	for s, counter := range vc.vec {
		c[s] = counter.Snapshot()
	}
	return &VectorCounterSnapshot{vec: c}
}

func (vc *StandardVectorCounter) getOrCreate(label string) Counter {
	counter, ok := vc.vec[label]
	if !ok {
		counter = NewCounter()
		vc.vec[label] = counter
	}
	return counter
}

// Nil implementation

type NilVectorCounter struct{}

func (vc *NilVectorCounter) DecAll(int64) {}

func (vc *NilVectorCounter) IncAll(int64) {}

func (vc *NilVectorCounter) GetAll() map[string]Counter { return nil }

func (vc *NilVectorCounter) ClearAll() {}

func (vc *NilVectorCounter) Dec(label string, v int64) {}

func (vc *NilVectorCounter) Inc(label string, v int64) {}

func (vc *NilVectorCounter) Get(label string) Counter { return NilCounter{} }

func (vc *NilVectorCounter) Clear(label string) {}

func (vc *NilVectorCounter) Snapshot() VectorCounter { return vc }

// Snapshot
type VectorCounterSnapshot struct {
	vec map[string]Counter
}

func (vc *VectorCounterSnapshot) DecAll(int64) {
	panic("decrement all on snapshot")
}

func (vc *VectorCounterSnapshot) IncAll(int64) {
	panic("increment all on snapshot")
}

func (vc *VectorCounterSnapshot) GetAll() map[string]Counter {
	return vc.vec
}

func (vc *VectorCounterSnapshot) ClearAll() {
	panic("clear all on snapshot")
}

func (vc *VectorCounterSnapshot) Dec(label string, v int64) {
	panic("decrement on snapshot")
}

func (vc *VectorCounterSnapshot) Inc(label string, v int64) {
	panic("decrement on snapshot")
}

func (vc *VectorCounterSnapshot) Get(label string) Counter {
	return vc.getOrCreate(label)
}

func (vc *VectorCounterSnapshot) Clear(label string) {
	panic("clear on snapshot")
}

func (vc *VectorCounterSnapshot) Snapshot() VectorCounter {
	return vc
}

func (vc *VectorCounterSnapshot) getOrCreate(label string) Counter {
	counter, ok := vc.vec[label]
	if !ok {
		counter = NewCounter()
		vc.vec[label] = counter
	}
	return counter
}
