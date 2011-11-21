package metrics

import (
	"fmt"
	"syslog"
	"time"
)

func Syslog(r Registry, interval int, w *syslog.Writer) {
	for {
		for name, c := range r.Counters() {
			w.Info(fmt.Sprintf("counter %s: count: %d", name, c.Count()))
		}
		for name, g := range r.Gauges() {
			w.Info(fmt.Sprintf("gauge %s: value: %d", name, g.Value()))
		}
		r.RunHealthchecks()
		for name, h := range r.Healthchecks() {
			w.Info(fmt.Sprintf("healthcheck %s: error: %v", name, h.Error()))
		}
		for name, h := range r.Histograms() {
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
		}
		for name, m := range r.Meters() {
			w.Info(fmt.Sprintf(
				"meter %s: count: %d 1-min: %.2f 5-min: %.2f 15-min: %.2f mean: %.2f",
				name,
				m.Count(),
				m.Rate1(),
				m.Rate5(),
				m.Rate15(),
				m.RateMean(),
			))
		}
		for name, t := range r.Timers() {
			w.Info(fmt.Sprintf(
				"timer %s: count: %d min: %d max: %d mean: %.2f stddev: %.2f median: %.2f 75%%: %.2f 95%%: %.2f 99%%: %.2f 99.9%%: %.2f 1-min: %.2f 5-min: %.2f 15-min: %.2f mean: %.2f",
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
				m.Rate1(),
				m.Rate5(),
				m.Rate15(),
				m.RateMean(),
			))
		}
		time.Sleep(int64(1e9) * int64(interval))
	}
}
