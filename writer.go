package metrics

import (
	"fmt"
	"io"
	"time"
)

// Output each metric in the given registry periodically using the given
// io.Writer.
func Write(r Registry, d time.Duration, w io.Writer) {
	for {
		WriteOnce(r, w)
		time.Sleep(d)
	}
}

func WriteOnce(r Registry, w io.Writer) {
	r.Each(func(name string, i interface{}) {
		switch m := i.(type) {
		case Counter:
			fmt.Fprintf(w, "counter %s\n", name)
			fmt.Fprintf(w, "  count:       %9d\n", m.Count())
		case Gauge:
			fmt.Fprintf(w, "gauge %s\n", name)
			fmt.Fprintf(w, "  value:       %9d\n", m.Value())
		case Healthcheck:
			m.Check()
			fmt.Fprintf(w, "healthcheck %s\n", name)
			fmt.Fprintf(w, "  error:       %v\n", m.Error())
		case Histogram:
			ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			fmt.Fprintf(w, "histogram %s\n", name)
			fmt.Fprintf(w, "  count:       %9d\n", m.Count())
			fmt.Fprintf(w, "  min:         %9d\n", m.Min())
			fmt.Fprintf(w, "  max:         %9d\n", m.Max())
			fmt.Fprintf(w, "  mean:        %12.2f\n", m.Mean())
			fmt.Fprintf(w, "  stddev:      %12.2f\n", m.StdDev())
			fmt.Fprintf(w, "  median:      %12.2f\n", ps[0])
			fmt.Fprintf(w, "  75%%:         %12.2f\n", ps[1])
			fmt.Fprintf(w, "  95%%:         %12.2f\n", ps[2])
			fmt.Fprintf(w, "  99%%:         %12.2f\n", ps[3])
			fmt.Fprintf(w, "  99.9%%:       %12.2f\n", ps[4])
		case Meter:
			fmt.Fprintf(w, "meter %s\n", name)
			fmt.Fprintf(w, "  count:       %9d\n", m.Count())
			fmt.Fprintf(w, "  1-min rate:  %12.2f\n", m.Rate1())
			fmt.Fprintf(w, "  5-min rate:  %12.2f\n", m.Rate5())
			fmt.Fprintf(w, "  15-min rate: %12.2f\n", m.Rate15())
			fmt.Fprintf(w, "  mean rate:   %12.2f\n", m.RateMean())
		case Timer:
			ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			fmt.Fprintf(w, "timer %s\n", name)
			fmt.Fprintf(w, "  count:       %9d\n", m.Count())
			fmt.Fprintf(w, "  min:         %9d\n", m.Min())
			fmt.Fprintf(w, "  max:         %9d\n", m.Max())
			fmt.Fprintf(w, "  mean:        %12.2f\n", m.Mean())
			fmt.Fprintf(w, "  stddev:      %12.2f\n", m.StdDev())
			fmt.Fprintf(w, "  median:      %12.2f\n", ps[0])
			fmt.Fprintf(w, "  75%%:         %12.2f\n", ps[1])
			fmt.Fprintf(w, "  95%%:         %12.2f\n", ps[2])
			fmt.Fprintf(w, "  99%%:         %12.2f\n", ps[3])
			fmt.Fprintf(w, "  99.9%%:       %12.2f\n", ps[4])
			fmt.Fprintf(w, "  1-min rate:  %12.2f\n", m.Rate1())
			fmt.Fprintf(w, "  5-min rate:  %12.2f\n", m.Rate5())
			fmt.Fprintf(w, "  15-min rate: %12.2f\n", m.Rate15())
			fmt.Fprintf(w, "  mean rate:   %12.2f\n", m.RateMean())
		}
	})
}
