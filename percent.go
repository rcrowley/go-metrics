package metrics

import (
	"sync"
)

type PercentCounter interface {
	Keys() []string
	Count(key string) int64
	Percent(key string) float64
	Total() int64
	Inc(key string, delta int64)
	Dec(key string, delta int64)
	Snapshot() PercentCounter
}

func NewPercentCounter() PercentCounter {
	return &StandardPercentCounter{items: make(map[string]int64)}
}

type StandardPercentCounter struct {
	mutex sync.Mutex
	items map[string]int64
}

func (pc *StandardPercentCounter) Keys() []string {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()
	r := make([]string, 0, len(pc.items))
	for key, _ := range pc.items {
		r = append(r, key)
	}
	return r
}

func (pc *StandardPercentCounter) Count(key string) int64 {
	return pc.items[key]
}

func (pc *StandardPercentCounter) Total() int64 {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()
	t := int64(0)
	for key, _ := range pc.items {
		t += pc.items[key]
	}
	return t
}

func (pc *StandardPercentCounter) Percent(key string) float64 {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()
	total := int64(0)
	for key, _ := range pc.items {
		total += pc.items[key]
	}
	return float64(pc.items[key]) * 100 / float64(total)
}

func (pc *StandardPercentCounter) Inc(key string, delta int64) {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()
	pc.items[key] += delta
}

func (pc *StandardPercentCounter) Dec(key string, delta int64) {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()
	pc.items[key] -= delta
}

// FIXME
func (pc *StandardPercentCounter) Snapshot() PercentCounter {
	return pc
}
