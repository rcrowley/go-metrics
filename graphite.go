package metrics

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func Graphite(r Registry, interval int, addr string) {
	for {
		now := time.Now().Unix()
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			continue
		}
		w := bufio.NewWriter(conn)
		r.Each(func(name string, i interface{}) {
			switch m := i.(type) {
			case Counter:
				fmt.Fprintf(w, "%s.count %d %d\n", name, m.Count(), now)
			case Gauge:
				fmt.Fprintf(w, "%s.value %d %d\n", name, m.Value(), now)
			case Histogram:
				ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
				fmt.Fprintf(w, "%s.count %d %d\n", name, m.Count(), now)
				fmt.Fprintf(w, "%s.min %d %d\n", name, m.Min(), now)
				fmt.Fprintf(w, "%s.max %d %d\n", name, m.Max(), now)
				fmt.Fprintf(w, "%s.mean %.2f %d\n", name, m.Mean(), now)
				fmt.Fprintf(w, "%s.std-dev %.2f %d\n", name, m.StdDev(), now)
				fmt.Fprintf(w, "%s.50-percentile %.2f %d\n", name, ps[0], now)
				fmt.Fprintf(w, "%s.75-percentile %.2f %d\n", name, ps[1], now)
				fmt.Fprintf(w, "%s.95-percentile %.2f %d\n", name, ps[2], now)
				fmt.Fprintf(w, "%s.99-percentile %.2f %d\n", name, ps[3], now)
				fmt.Fprintf(w, "%s.999-percentile %.2f %d\n", name, ps[4], now)
			case Meter:
				fmt.Fprintf(w, "%s.count %d %d\n", name, m.Count(), now)
				fmt.Fprintf(w, "%s.one-minute %.2f %d\n", name, m.Rate1(), now)
				fmt.Fprintf(w, "%s.five-minute %.2f %d\n", name, m.Rate5(), now)
				fmt.Fprintf(w, "%s.fifteen-minute %.2f %d\n", name, m.Rate15(), now)
				fmt.Fprintf(w, "%s.mean %.2f %d\n", name, m.RateMean(), now)
			case Timer:
				ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
				fmt.Fprintf(w, "%s.count %d %d\n", name, m.Count(), now)
				fmt.Fprintf(w, "%s.min %d %d\n", name, m.Min(), now)
				fmt.Fprintf(w, "%s.max %d %d\n", name, m.Max(), now)
				fmt.Fprintf(w, "%s.mean %.2f %d\n", name, m.Mean(), now)
				fmt.Fprintf(w, "%s.std-dev %.2f %d\n", name, m.StdDev(), now)
				fmt.Fprintf(w, "%s.50-percentile %.2f %d\n", name, ps[0], now)
				fmt.Fprintf(w, "%s.75-percentile %.2f %d\n", name, ps[1], now)
				fmt.Fprintf(w, "%s.95-percentile %.2f %d\n", name, ps[2], now)
				fmt.Fprintf(w, "%s.99-percentile %.2f %d\n", name, ps[3], now)
				fmt.Fprintf(w, "%s.999-percentile %.2f %d\n", name, ps[4], now)
				fmt.Fprintf(w, "%s.one-minute %.2f %d\n", name, m.Rate1(), now)
				fmt.Fprintf(w, "%s.five-minute %.2f %d\n", name, m.Rate5(), now)
				fmt.Fprintf(w, "%s.fifteen-minute %.2f %d\n", name, m.Rate15(), now)
				fmt.Fprintf(w, "%s.mean %.2f %d\n", name, m.RateMean(), now)
			}
			w.Flush()
		})
		time.Sleep(time.Duration(int64(1e9) * int64(interval)))
	}
}
