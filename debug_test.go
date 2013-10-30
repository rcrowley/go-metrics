package metrics

import (
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
	go func() {
		i := 0
		for {
			select {
			case ch <- i:
				return
			default:
				i++
			}
			time.Sleep(1e3)
		}
	}()
	var gcStats debug.GCStats
	debug.ReadGCStats(&gcStats)
	t.Log(<-ch)
}
