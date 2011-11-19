package metrics

import (
	"log"
	"time"
)

func Log(r Registry, interval int, l *log.Logger) {
	for {
		for name, c := range r.Counters() {
			l.Printf("counter %s\n\tcount:\t%9d\n", name, c.Count())
		}
		for name, g := range r.Gauges() {
			l.Printf("gauge %s\n\tvalue:\t%9d\n", name, g.Value())
		}
		for name, h := range r.Healthchecks() {
			l.Printf("healthcheck %s TODO\n", name, h)
		}
		for name, h := range r.Histograms() {
			ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			l.Printf(
				`histogram %s
	count:	%9d
	min:	%9d
	max:	%9d
	mean:	%12.2f
	stddev:	%12.2f
	median:	%12.2f
	75%%:	%12.2f
	95%%:	%12.2f
	99%%:	%12.2f
	99.9%%:	%12.2f
`,
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
			)
		}
		for name, m := range r.Meters() {
			l.Printf(
				`meter %s
	count:	%9d
	1-min:	%12.2f
	5-min:	%12.2f
	15-min:	%12.2f
	mean:	%12.2f
`,
				name,
				m.Count(),
				m.Rate1(),
				m.Rate5(),
				m.Rate15(),
				m.RateMean(),
			)
		}
		for name, t := range r.Timers() {
			l.Printf("timer %s TODO\n", name, t)
		}
		time.Sleep(int64(1e9) * int64(interval))
	}
}
