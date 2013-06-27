package metrics

import (
	"log"
	"time"
)

// Output each metric in the given registry periodically using the given
// logger.
func Log(r Registry, d time.Duration, l *log.Logger) {
	for {
		r.Each(func(name string, i interface{}) {
			switch m := i.(type) {
			case Counter:
				l.Printf("counter %s\n", name)
				l.Printf("  count:       %9d\n", m.Count())
			case Gauge:
				l.Printf("gauge %s\n", name)
				l.Printf("  value:       %9d\n", m.Value())
			case Healthcheck:
				m.Check()
				l.Printf("healthcheck %s\n", name)
				l.Printf("  error:       %v\n", m.Error())
			case Histogram:
				ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
				l.Printf("histogram %s\n", name)
				l.Printf("  count:       %9d\n", m.Count())
				l.Printf("  min:         %9d\n", m.Min())
				l.Printf("  max:         %9d\n", m.Max())
				l.Printf("  mean:        %12.2f\n", m.Mean())
				l.Printf("  stddev:      %12.2f\n", m.StdDev())
				l.Printf("  median:      %12.2f\n", ps[0])
				l.Printf("  75%%:         %12.2f\n", ps[1])
				l.Printf("  95%%:         %12.2f\n", ps[2])
				l.Printf("  99%%:         %12.2f\n", ps[3])
				l.Printf("  99.9%%:       %12.2f\n", ps[4])
			case Meter:
				l.Printf("meter %s\n", name)
				l.Printf("  count:       %9d\n", m.Count())
				l.Printf("  1-min rate:  %12.2f\n", m.Rate1())
				l.Printf("  5-min rate:  %12.2f\n", m.Rate5())
				l.Printf("  15-min rate: %12.2f\n", m.Rate15())
				l.Printf("  mean rate:   %12.2f\n", m.RateMean())
			case Timer:
				ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
				l.Printf("timer %s\n", name)
				l.Printf("  count:       %9d\n", m.Count())
				l.Printf("  min:         %9d\n", m.Min())
				l.Printf("  max:         %9d\n", m.Max())
				l.Printf("  mean:        %12.2f\n", m.Mean())
				l.Printf("  stddev:      %12.2f\n", m.StdDev())
				l.Printf("  median:      %12.2f\n", ps[0])
				l.Printf("  75%%:         %12.2f\n", ps[1])
				l.Printf("  95%%:         %12.2f\n", ps[2])
				l.Printf("  99%%:         %12.2f\n", ps[3])
				l.Printf("  99.9%%:       %12.2f\n", ps[4])
				l.Printf("  1-min rate:  %12.2f\n", m.Rate1())
				l.Printf("  5-min rate:  %12.2f\n", m.Rate5())
				l.Printf("  15-min rate: %12.2f\n", m.Rate15())
				l.Printf("  mean rate:   %12.2f\n", m.RateMean())
			}
		})
		time.Sleep(d)
	}
}
