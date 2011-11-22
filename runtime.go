package metrics

import "runtime"

// Register metrics for the Go runtime statistics exported in
// runtime.MemStats.  The metrics are named by their fully-qualified
// Go symbols, i.e. runtime.MemStatsAlloc.  In addition to
// runtime.MemStats, register the return value of runtime.Goroutines()
// as runtime.Goroutines.
func RegisterRuntimeMemStats(r Registry) {

	r.RegisterGauge("runtime.Goroutines", NewGauge())

	r.RegisterGauge("runtime.MemStats.Alloc", NewGauge())
	r.RegisterGauge("runtime.MemStats.TotalAlloc", NewGauge())
	r.RegisterGauge("runtime.MemStats.Sys", NewGauge())
	r.RegisterGauge("runtime.MemStats.Lookups", NewGauge())
	r.RegisterGauge("runtime.MemStats.Mallocs", NewGauge())
	r.RegisterGauge("runtime.MemStats.Frees", NewGauge())

	r.RegisterGauge("runtime.MemStats.HeapAlloc", NewGauge())
	r.RegisterGauge("runtime.MemStats.HeapSys", NewGauge())
	r.RegisterGauge("runtime.MemStats.HeapIdle", NewGauge())
	r.RegisterGauge("runtime.MemStats.HeapInuse", NewGauge())
	r.RegisterGauge("runtime.MemStats.HeapObjects", NewGauge())

	r.RegisterGauge("runtime.MemStats.StackInuse", NewGauge())
	r.RegisterGauge("runtime.MemStats.StackSys", NewGauge())
	r.RegisterGauge("runtime.MemStats.MSpanInuse", NewGauge())
	r.RegisterGauge("runtime.MemStats.MSpanSys", NewGauge())
	r.RegisterGauge("runtime.MemStats.MCacheInuse", NewGauge())
	r.RegisterGauge("runtime.MemStats.MCacheSys", NewGauge())
	r.RegisterGauge("runtime.MemStats.BuckHashSys", NewGauge())

	r.RegisterGauge("runtime.MemStats.NextGC", NewGauge())
	r.RegisterGauge("runtime.MemStats.PauseTotalNs", NewGauge())
	r.RegisterHistogram("runtime.MemStats.PauseNs",
		NewHistogram(NewExpDecaySample(1028, 0.015)))
	r.RegisterGauge("runtime.MemStats.NumGC", NewGauge())
	r.RegisterGauge("runtime.MemStats.EnableGC", NewGauge())
	r.RegisterGauge("runtime.MemStats.DebugGC", NewGauge())

}

// Capture new values for the Go runtime statistics exported in
// runtime.MemStats.  This is designed to be called in a background
// goroutine.  Giving a registry which has not been given to
// RegisterRuntimeMemStats will panic.  If the second parameter is
// false, the counters will be left to the lazy updates provided by
// the runtime.
func CaptureRuntimeMemStats(r Registry, updateMemStats bool) {
	if updateMemStats {
		runtime.UpdateMemStats()
	}

	r.GetGauge("runtime.Goroutines").Update(int64(runtime.Goroutines()))

	r.GetGauge("runtime.MemStats.Alloc").Update(
		int64(runtime.MemStats.Alloc))
	r.GetGauge("runtime.MemStats.TotalAlloc").Update(
		int64(runtime.MemStats.TotalAlloc))
	r.GetGauge("runtime.MemStats.Sys").Update(
		int64(runtime.MemStats.Sys))
	r.GetGauge("runtime.MemStats.Lookups").Update(
		int64(runtime.MemStats.Lookups))
	r.GetGauge("runtime.MemStats.Mallocs").Update(
		int64(runtime.MemStats.Mallocs))
	r.GetGauge("runtime.MemStats.Frees").Update(
		int64(runtime.MemStats.Frees))

	r.GetGauge("runtime.MemStats.HeapAlloc").Update(
		int64(runtime.MemStats.HeapAlloc))
	r.GetGauge("runtime.MemStats.HeapSys").Update(
		int64(runtime.MemStats.HeapSys))
	r.GetGauge("runtime.MemStats.HeapIdle").Update(
		int64(runtime.MemStats.HeapIdle))
	r.GetGauge("runtime.MemStats.HeapInuse").Update(
		int64(runtime.MemStats.HeapInuse))
	r.GetGauge("runtime.MemStats.HeapObjects").Update(
		int64(runtime.MemStats.HeapObjects))

	r.GetGauge("runtime.MemStats.StackInuse").Update(
		int64(runtime.MemStats.StackInuse))
	r.GetGauge("runtime.MemStats.StackSys").Update(
		int64(runtime.MemStats.StackSys))
	r.GetGauge("runtime.MemStats.MSpanInuse").Update(
		int64(runtime.MemStats.MSpanInuse))
	r.GetGauge("runtime.MemStats.MSpanSys").Update(
		int64(runtime.MemStats.MSpanSys))
	r.GetGauge("runtime.MemStats.MCacheInuse").Update(
		int64(runtime.MemStats.MCacheInuse))
	r.GetGauge("runtime.MemStats.MCacheSys").Update(
		int64(runtime.MemStats.MCacheSys))
	r.GetGauge("runtime.MemStats.BuckHashSys").Update(
		int64(runtime.MemStats.BuckHashSys))

	r.GetGauge("runtime.MemStats.NextGC").Update(
		int64(runtime.MemStats.NextGC))
	r.GetGauge("runtime.MemStats.PauseTotalNs").Update(
		int64(runtime.MemStats.PauseTotalNs))
	r.GetHistogram("runtime.MemStats.PauseNs").Update(
		int64(runtime.MemStats.PauseNs[0]))
	r.GetGauge("runtime.MemStats.NumGC").Update(
		int64(runtime.MemStats.NumGC))
	if runtime.MemStats.EnableGC {
		r.GetGauge("runtime.MemStats.EnableGC").Update(1)
	} else {
		r.GetGauge("runtime.MemStats.EnableGC").Update(0)
	}
	if runtime.MemStats.EnableGC {
		r.GetGauge("runtime.MemStats.DebugGC").Update(1)
	} else {
		r.GetGauge("runtime.MemStats.DebugGC").Update(0)
	}

}
