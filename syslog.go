package metrics

import (
	"fmt"
	"syslog"
	"time"
)

// Output each metric in the given registry to syslog periodically using
// the given syslogger.  The interval is to be given in seconds.
func Syslog(r Registry, interval int, w *syslog.Writer) {
	for {
		r.EachCounter(func(name string, c Counter) {
			w.Info(fmt.Sprintf("counter %s: count: %d", name, c.Count()))
		})
		r.EachGauge(func(name string, g Gauge) {
			w.Info(fmt.Sprintf("gauge %s: value: %d", name, g.Value()))
		})
		r.RunHealthchecks()
		r.EachHealthcheck(func(name string, h Healthcheck) {
			w.Info(fmt.Sprintf("healthcheck %s: error: %v", name, h.Error()))
		})
		r.EachHistogram(func(name string, h Histogram) {
			ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			w.Info(fmt.Sprintf(
				"histogram %s: count: %d min: %d max: %d mean: %.2f stddev: %.2f median: %.2f 75%%: %.2f 95%%: %.2f 99%%: %.2f 99.9%%: %.2f",
				name,
				h.Count(),
				h.Min(),
				h.Max(),
				h.Mean(),
				h.StdDev(),
				ps[0],
				ps[1],
				ps[2],
				ps[3],
				ps[4],
			))
		})
		r.EachMeter(func(name string, m Meter) {
			w.Info(fmt.Sprintf(
				"meter %s: count: %d 1-min: %.2f 5-min: %.2f 15-min: %.2f mean: %.2f",
				name,
				m.Count(),
				m.Rate1(),
				m.Rate5(),
				m.Rate15(),
				m.RateMean(),
			))
		})
		r.EachTimer(func(name string, t Timer) {
			ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			w.Info(fmt.Sprintf(
				"timer %s: count: %d min: %d max: %d mean: %.2f stddev: %.2f median: %.2f 75%%: %.2f 95%%: %.2f 99%%: %.2f 99.9%%: %.2f 1-min: %.2f 5-min: %.2f 15-min: %.2f mean: %.2f",
				name,
				t.Count(),
				t.Min(),
				t.Max(),
				t.Mean(),
				t.StdDev(),
				ps[0],
				ps[1],
				ps[2],
				ps[3],
				ps[4],
				t.Rate1(),
				t.Rate5(),
				t.Rate15(),
				t.RateMean(),
			))
		})
		time.Sleep(int64(1e9) * int64(interval))
	}
}
