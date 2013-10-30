package metrics

import (
	"io/ioutil"
	"log"
	"sync"
	"testing"
)

const FANOUT = 4

// Stop the compiler from complaining during debugging.
var (
	_ = ioutil.Discard
	_ = log.LstdFlags
)

func BenchmarkMetrics(b *testing.B) {
	r := NewRegistry()
	c := NewRegisteredCounter("counter", r)
	g := NewRegisteredGauge("gauge", r)
	h := NewRegisteredHistogram("histogram", r, NewUniformSample(100))
	m := NewRegisteredMeter("meter", r)
	t := NewRegisteredTimer("timer", r)
	RegisterRuntimeMemStats(r)
	b.ResetTimer()
	ch := make(chan bool)
	wgC := &sync.WaitGroup{}
//*
	wgC.Add(1)
	go func() {
		defer wgC.Done()
		//log.Println("go CaptureRuntimeMemStats")
		for {
			select {
			case <-ch:
				//log.Println("done CaptureRuntimeMemStats")
				return
			default:
				CaptureRuntimeMemStatsOnce(r)
			}
		}
	}()
//*/
	wgW := &sync.WaitGroup{}
/*
	wgW.Add(1)
	go func() {
		defer wgW.Done()
		//log.Println("go Write")
		for {
			select {
			case <-ch:
				//log.Println("done Write")
				return
			default:
				WriteOnce(r, ioutil.Discard)
			}
		}
	}()
//*/
	wg := &sync.WaitGroup{}
	wg.Add(FANOUT)
	for i := 0; i < FANOUT; i++ {
		go func(i int) {
			defer wg.Done()
			//log.Println("go", i)
			for i := 0; i < b.N; i++ {
				c.Inc(1)
				g.Update(int64(i))
				h.Update(int64(i))
				m.Mark(1)
				t.Update(1)
			}
			//log.Println("done", i)
		}(i)
	}
	wg.Wait()
	close(ch)
	wgC.Wait()
	wgW.Wait()
}
