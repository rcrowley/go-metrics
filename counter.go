package metrics

import "sync/atomic"

type Counter interface {
	Clear()
	Count() int64
	Dec(int64)
	Inc(int64)
}

type counter struct {
	count int64
}

func NewCounter() Counter {
	return &counter{0}
}

func (c *counter) Clear() {
	c.count = 0
}

func (c *counter) Count() int64 {
	return c.count
}

func (c *counter) Dec(i int64) {
	atomic.AddInt64(&c.count, -i)
}

func (c *counter) Inc(i int64) {
	atomic.AddInt64(&c.count, i)
}
