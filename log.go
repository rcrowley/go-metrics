package metrics

import (
	"log"
	"time"
)

func Log(r Registry, interval int, l *log.Logger) {
	for {
		for name, c := range r.Counters() {
			l.Printf("counter %s\n", name)
			l.Printf("  count:       %9d\n", c.Count())
		}
		for name, g := range r.Gauges() {
			l.Printf("gauge %s\n", name)
			l.Printf("  value:       %9d\n", g.Value())
		}
		r.RunHealthchecks()
		for name, h := range r.Healthchecks() {
			l.Printf("healthcheck %s\n", name)
			l.Printf("  error:       %v\n", h.Error())
		}
		for name, h := range r.Histograms() {
			ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			l.Printf("histogram %s\n", name)
			l.Printf("  count:       %9d\n", h.Count())
			l.Printf("  min:         %9d\n", h.Min())
			l.Printf("  max:         %9d\n", h.Max())
			l.Printf("  mean:        %12.2f\n", h.Mean())
			l.Printf("  stddev:      %12.2f\n", h.StdDev())
			l.Printf("  median:      %12.2f\n", ps[0])
			l.Printf("  75%%:         %12.2f\n", ps[1])
			l.Printf("  95%%:         %12.2f\n", ps[2])
			l.Printf("  99%%:         %12.2f\n", ps[3])
			l.Printf("  99.9%%:       %12.2f\n", ps[4])
		}
		for name, m := range r.Meters() {
			l.Printf("meter %s\n", name)
			l.Printf("  count:       %9d\n", m.Count())
			l.Printf("  1-min rate:  %12.2f\n", m.Rate1())
			l.Printf("  5-min rate:  %12.2f\n", m.Rate5())
			l.Printf("  15-min rate: %12.2f\n", m.Rate15())
			l.Printf("  mean rate:   %12.2f\n", m.RateMean())
		}
		for name, t := range r.Timers() {
			ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			l.Printf("timer %s\n", name)
			l.Printf("  count:       %9d\n", t.Count())
			l.Printf("  min:         %9d\n", t.Min())
			l.Printf("  max:         %9d\n", t.Max())
			l.Printf("  mean:        %12.2f\n", t.Mean())
			l.Printf("  stddev:      %12.2f\n", t.StdDev())
			l.Printf("  median:      %12.2f\n", ps[0])
			l.Printf("  75%%:         %12.2f\n", ps[1])
			l.Printf("  95%%:         %12.2f\n", ps[2])
			l.Printf("  99%%:         %12.2f\n", ps[3])
			l.Printf("  99.9%%:       %12.2f\n", ps[4])
			l.Printf("  1-min rate:  %12.2f\n", t.Rate1())
			l.Printf("  5-min rate:  %12.2f\n", t.Rate5())
			l.Printf("  15-min rate: %12.2f\n", t.Rate15())
			l.Printf("  mean rate:   %12.2f\n", t.RateMean())
		}
		time.Sleep(int64(1e9) * int64(interval))
	}
}
