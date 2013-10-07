// +build !windows

package metrics

import (
	"fmt"
	"log/syslog"
	"time"
)

// Output each metric in the given registry to syslog periodically using
// the given syslogger.
func Syslog(r Registry, d time.Duration, w *syslog.Writer) {
	for {
		r.Each(func(name string, i interface{}) {
			switch m := i.(type) {
			case Counter:
				w.Info(fmt.Sprintf("counter %s: count: %d", name, m.Count()))
			case Gauge:
				w.Info(fmt.Sprintf("gauge %s: value: %d", name, m.Value()))
			case Healthcheck:
				m.Check()
				w.Info(fmt.Sprintf("healthcheck %s: error: %v", name, m.Error()))
			case Histogram:
				ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
				w.Info(fmt.Sprintf(
					"histogram %s: count: %d min: %d max: %d mean: %.2f stddev: %.2f median: %.2f 75%%: %.2f 95%%: %.2f 99%%: %.2f 99.9%%: %.2f",
					name,
					m.Count(),
					m.Min(),
					m.Max(),
					m.Mean(),
					m.StdDev(),
					ps[0],
					ps[1],
					ps[2],
					ps[3],
					ps[4],
				))
			case Meter:
				w.Info(fmt.Sprintf(
					"meter %s: count: %d 1-min: %.2f 5-min: %.2f 15-min: %.2f mean: %.2f",
					name,
					m.Count(),
					m.Rate1(),
					m.Rate5(),
					m.Rate15(),
					m.RateMean(),
				))
			case Timer:
				ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
				w.Info(fmt.Sprintf(
					"timer %s: count: %d min: %d max: %d mean: %.2f stddev: %.2f median: %.2f 75%%: %.2f 95%%: %.2f 99%%: %.2f 99.9%%: %.2f 1-min: %.2f 5-min: %.2f 15-min: %.2f mean-rate: %.2f",
					name,
					m.Count(),
					m.Min(),
					m.Max(),
					m.Mean(),
					m.StdDev(),
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
		})
		time.Sleep(d)
	}
}
