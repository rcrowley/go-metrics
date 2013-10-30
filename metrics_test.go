package metrics

import (
	"io/ioutil"
	"testing"
)

var _ = ioutil.Discard // Stop the compiler from complaining during debugging.

func BenchmarkMetrics(b *testing.B) {
	c := NewCounter()
	g := NewGauge()
	h := NewHistogram(NewUniformSample(100))
	m := NewMeter()
	t := NewTimer()
	r := NewRegistry()
	r.Register("counter", c)
	r.Register("gauge", g)
	r.Register("histogram", h)
	r.Register("meter", m)
	r.Register("timer", t)
	RegisterRuntimeMemStats(r)
	ch := make(chan bool)
/*
	go func() {
		for {
			select {
			case <-ch:
				return
			default:
				CaptureRuntimeMemStatsOnce(r)
			}
		}
	}()
//*/
//*
	go func() {
		for {
			select {
			case <-ch:
				return
			default:
				WriteOnce(r, ioutil.Discard)
			}
		}
	}()
//*/
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Inc(1)
		g.Update(int64(i))
		h.Update(int64(i))
		m.Mark(1)
		t.Update(1)
	}
	b.StopTimer()
	close(ch)
}
