package metrics

import (
	"sort"
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

type sortedMap struct {
	m map[string]int64
	s []string
}

func newSortedMap() *sortedMap {
	this := new(sortedMap)
	this.m = make(map[string]int64)
	return this
}

func (this *sortedMap) set(key string, val int64) {
	this.m[key] = val
}

func (this *sortedMap) get(key string) int64 {
	return this.m[key]
}

func (this *sortedMap) inc(key string, delta int64) int64 {
	v, present := this.m[key]
	if !present {
		v = 0
	}
	v += delta
	this.m[key] = v
	return v
}

func (this *sortedMap) Len() int {
	return len(this.m)
}

func (this *sortedMap) Less(i, j int) bool {
	return this.m[this.s[i]] > this.m[this.s[j]]
}

func (this *sortedMap) Swap(i, j int) {
	this.s[i], this.s[j] = this.s[j], this.s[i]
}

func (this *sortedMap) sortedKeys() []string {
	this.s = make([]string, len(this.m))
	i := 0
	for key, _ := range this.m {
		this.s[i] = key
		i++
	}
	sort.Sort(this)
	return this.s
}

func NewPercentCounter() PercentCounter {
	return &StandardPercentCounter{items: newSortedMap()}
}

type StandardPercentCounter struct {
	mutex sync.Mutex
	items *sortedMap
}

func (pc *StandardPercentCounter) Keys() []string {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()
	r := make([]string, 0, pc.items.Len())
	for _, key := range pc.items.sortedKeys() {
		r = append(r, key)
	}
	return r
}

func (pc *StandardPercentCounter) Count(key string) int64 {
	return pc.items.get(key)
}

func (pc *StandardPercentCounter) Total() int64 {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()
	t := int64(0)
	for key, _ := range pc.items.m {
		t += pc.items.get(key)
	}
	return t
}

func (pc *StandardPercentCounter) Percent(key string) float64 {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()
	total := int64(0)
	for key, _ := range pc.items.m {
		total += pc.items.get(key)
	}
	return float64(pc.items.get(key)) * 100 / float64(total)
}

func (pc *StandardPercentCounter) Inc(key string, delta int64) {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()
	pc.items.inc(key, delta)
}

func (pc *StandardPercentCounter) Dec(key string, delta int64) {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()
	pc.items.inc(key, -delta)
}

// FIXME
func (pc *StandardPercentCounter) Snapshot() PercentCounter {
	return pc
}
