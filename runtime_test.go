package metrics

import (
	"runtime"
	"testing"
	"time"
)

func BenchmarkRuntimeMemStats(b *testing.B) {
	r := NewRegistry()
	RegisterRuntimeMemStats(r)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CaptureRuntimeMemStatsOnce(r)
	}
}

func TestRuntimeMemStatsBlocking(t *testing.T) {
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
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	t.Log(<-ch)
}
