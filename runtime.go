package metrics

import (
	"runtime"
	"time"
)

var (
	numGC    uint32
	memStats runtime.MemStats
)

// Capture new values for the Go runtime statistics exported in
// runtime.MemStats.  This is designed to be called as a goroutine.
func CaptureRuntimeMemStats(r Registry, interval int64) {
	for {
		CaptureRuntimeMemStatsOnce(r)
		time.Sleep(time.Duration(int64(1e9) * int64(interval)))
	}
}

// Capture new values for the Go runtime statistics exported in
// runtime.MemStats.  This is designed to be called in a background
// goroutine.  Giving a registry which has not been given to
// RegisterRuntimeMemStats will panic.
func CaptureRuntimeMemStatsOnce(r Registry) {
	runtime.ReadMemStats(&memStats)

	r.Get("runtime.MemStats.Alloc").(Gauge).Update(int64(memStats.Alloc))
	r.Get("runtime.MemStats.TotalAlloc").(Gauge).Update(int64(memStats.TotalAlloc))
	r.Get("runtime.MemStats.Sys").(Gauge).Update(int64(memStats.Sys))
	r.Get("runtime.MemStats.Lookups").(Gauge).Update(int64(memStats.Lookups))
	r.Get("runtime.MemStats.Mallocs").(Gauge).Update(int64(memStats.Mallocs))
	r.Get("runtime.MemStats.Frees").(Gauge).Update(int64(memStats.Frees))

	r.Get("runtime.MemStats.HeapAlloc").(Gauge).Update(int64(memStats.HeapAlloc))
	r.Get("runtime.MemStats.HeapSys").(Gauge).Update(int64(memStats.HeapSys))
	r.Get("runtime.MemStats.HeapIdle").(Gauge).Update(int64(memStats.HeapIdle))
	r.Get("runtime.MemStats.HeapInuse").(Gauge).Update(int64(memStats.HeapInuse))
	r.Get("runtime.MemStats.HeapReleased").(Gauge).Update(int64(memStats.HeapReleased))
	r.Get("runtime.MemStats.HeapObjects").(Gauge).Update(int64(memStats.HeapObjects))

	r.Get("runtime.MemStats.StackInuse").(Gauge).Update(int64(memStats.StackInuse))
	r.Get("runtime.MemStats.StackSys").(Gauge).Update(int64(memStats.StackSys))
	r.Get("runtime.MemStats.MSpanInuse").(Gauge).Update(int64(memStats.MSpanInuse))
	r.Get("runtime.MemStats.MSpanSys").(Gauge).Update(int64(memStats.MSpanSys))
	r.Get("runtime.MemStats.MCacheInuse").(Gauge).Update(int64(memStats.MCacheInuse))
	r.Get("runtime.MemStats.MCacheSys").(Gauge).Update(int64(memStats.MCacheSys))
	r.Get("runtime.MemStats.BuckHashSys").(Gauge).Update(int64(memStats.BuckHashSys))

	r.Get("runtime.MemStats.NextGC").(Gauge).Update(int64(memStats.NextGC))
	r.Get("runtime.MemStats.LastGC").(Gauge).Update(int64(memStats.LastGC))
	r.Get("runtime.MemStats.PauseTotalNs").(Gauge).Update(int64(memStats.PauseTotalNs))
	// <https://code.google.com/p/go/source/browse/src/pkg/runtime/mgc0.c>
	for i := uint32(1); i <= memStats.NumGC-numGC; i++ {
		r.Get("runtime.MemStats.PauseNs").(Histogram).Update(int64(memStats.PauseNs[(memStats.NumGC%256-i)%256]))
	}
	r.Get("runtime.MemStats.NumGC").(Gauge).Update(int64(memStats.NumGC))
	if memStats.EnableGC {
		r.Get("runtime.MemStats.EnableGC").(Gauge).Update(1)
	} else {
		r.Get("runtime.MemStats.EnableGC").(Gauge).Update(0)
	}
	if memStats.EnableGC {
		r.Get("runtime.MemStats.DebugGC").(Gauge).Update(1)
	} else {
		r.Get("runtime.MemStats.DebugGC").(Gauge).Update(0)
	}

	r.Get("runtime.NumCgoCall").(Gauge).Update(int64(runtime.NumCgoCall()))
	r.Get("runtime.NumGoroutine").(Gauge).Update(int64(runtime.NumGoroutine()))

}

// Register metrics for the Go runtime statistics exported in
// runtime.MemStats.  The metrics are named by their fully-qualified
// Go symbols, i.e. runtime.MemStatsAlloc.  In addition to
// runtime.MemStats, register the return value of runtime.Goroutines()
// as runtime.Goroutines.
func RegisterRuntimeMemStats(r Registry) {

	r.Register("runtime.MemStats.Alloc", NewGauge())
	r.Register("runtime.MemStats.TotalAlloc", NewGauge())
	r.Register("runtime.MemStats.Sys", NewGauge())
	r.Register("runtime.MemStats.Lookups", NewGauge())
	r.Register("runtime.MemStats.Mallocs", NewGauge())
	r.Register("runtime.MemStats.Frees", NewGauge())

	r.Register("runtime.MemStats.HeapAlloc", NewGauge())
	r.Register("runtime.MemStats.HeapSys", NewGauge())
	r.Register("runtime.MemStats.HeapIdle", NewGauge())
	r.Register("runtime.MemStats.HeapInuse", NewGauge())
	r.Register("runtime.MemStats.HeapReleased", NewGauge())
	r.Register("runtime.MemStats.HeapObjects", NewGauge())

	r.Register("runtime.MemStats.StackInuse", NewGauge())
	r.Register("runtime.MemStats.StackSys", NewGauge())
	r.Register("runtime.MemStats.MSpanInuse", NewGauge())
	r.Register("runtime.MemStats.MSpanSys", NewGauge())
	r.Register("runtime.MemStats.MCacheInuse", NewGauge())
	r.Register("runtime.MemStats.MCacheSys", NewGauge())
	r.Register("runtime.MemStats.BuckHashSys", NewGauge())

	r.Register("runtime.MemStats.NextGC", NewGauge())
	r.Register("runtime.MemStats.LastGC", NewGauge())
	r.Register("runtime.MemStats.PauseTotalNs", NewGauge())
	r.Register("runtime.MemStats.PauseNs", NewHistogram(NewExpDecaySample(1028, 0.015)))
	r.Register("runtime.MemStats.NumGC", NewGauge())
	r.Register("runtime.MemStats.EnableGC", NewGauge())
	r.Register("runtime.MemStats.DebugGC", NewGauge())

	r.Register("runtime.NumCgoCall", NewGauge())
	r.Register("runtime.NumGoroutine", NewGauge())

}
