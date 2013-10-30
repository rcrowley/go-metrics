package metrics

import (
	// "runtime"
	"runtime/debug"
	"testing"
	"time"
)

func BenchmarkDebugGCStats(b *testing.B) {
	r := NewRegistry()
	RegisterDebugGCStats(r)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CaptureDebugGCStatsOnce(r)
	}
}

func TestDebugGCStatsBlocking(t *testing.T) {
	ch := make(chan int)
	go testDebugGCStatsBlocking(ch)
	//runtime.Gosched()
	var gcStats debug.GCStats
	t0 := time.Now()
	debug.ReadGCStats(&gcStats)
	t1 := time.Now()
	t.Log("i++ during debug.ReadGCStats:", <-ch)
	go testDebugGCStatsBlocking(ch)
	d := t1.Sub(t0)
	t.Log(d)
	time.Sleep(d)
	t.Log("i++ during time.Sleep:", <-ch)
}

func testDebugGCStatsBlocking(ch chan int) {
	i := 0
	for {
		select {
		case ch <- i:
			return
		default:
			i++
		}
	}
}
