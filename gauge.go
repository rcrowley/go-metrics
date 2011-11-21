package metrics

import "sync/atomic"

type Gauge interface {
	Update(int64)
	Value() int64
}

type gauge struct {
	value int64
}

func NewGauge() Gauge {
	return &gauge{0}
}

func (g *gauge) Update(v int64) {
	atomic.AddInt64(&g.value, v)
}

func (g *gauge) Value() int64 {
	return g.value
}
