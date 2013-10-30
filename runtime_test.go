package metrics

import (
	"runtime"
	"testing"
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
		}
	}()
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	t.Log(<-ch)
}
