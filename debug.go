package metrics

import (
	"runtime/debug"
	"time"
)

var gcStats debug.GCStats

// Capture new values for the Go garbage collector statistics exported in
// debug.GCStats.  This is designed to be called as a goroutine.
func CaptureDebugGCStats(r Registry, d time.Duration) {
	for {
		CaptureDebugGCStatsOnce(r)
		time.Sleep(d)
	}
}

// Capture new values for the Go garbage collector statistics exported in
// debug.GCStats.  This is designed to be called in a background goroutine.
// Giving a registry which has not been given to RegisterDebugGCStats will
// panic.
func CaptureDebugGCStatsOnce(r Registry) {
	lastGC := gcStats.LastGC
	debug.ReadGCStats(&gcStats)
	r.Get("debug.GCStats.LastGC").(Gauge).Update(int64(gcStats.LastGC.UnixNano()))
	r.Get("debug.GCStats.NumGC").(Gauge).Update(int64(gcStats.NumGC))
	r.Get("debug.GCStats.PauseTotal").(Gauge).Update(int64(gcStats.PauseTotal))
	if lastGC != gcStats.LastGC && 0 < len(gcStats.Pause) {
		r.Get("debug.GCStats.Pause").(Histogram).Update(int64(gcStats.Pause[0]))
	}
	//r.Get("debug.GCStats.PauseQuantiles").(Histogram).Update(gcStats.PauseQuantiles)
}

// Register metrics for the Go garbage collector statistics exported in
// debug.GCStats.  The metrics are named by their fully-qualified Go symbols,
// i.e. debug.GCStats.PauseTotal.
func RegisterDebugGCStats(r Registry) {
	r.Register("debug.GCStats.LastGC", NewGauge())
	r.Register("debug.GCStats.NumGC", NewGauge())
	r.Register("debug.GCStats.PauseTotal", NewGauge())
	r.Register("debug.GCStats.Pause", NewHistogram(NewExpDecaySample(1028, 0.015)))
	//r.Register("debug.GCStats.PauseQuantiles", NewHistogram(NewExpDecaySample(1028, 0.015)))
}
