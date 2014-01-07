package metrics

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func Graphite(r Registry, d time.Duration, prefix string, addr *net.TCPAddr) {
	for {
		if err := graphite(r, prefix, addr); nil != err {
			log.Println(err)
		}
		time.Sleep(d)
	}
}

func graphite(r Registry, prefix string, addr *net.TCPAddr) error {
	now := time.Now().Unix()
	conn, err := net.DialTCP("tcp", nil, addr)
	if nil != err {
		return err
	}
	defer conn.Close()
	w := bufio.NewWriter(conn)
	r.Each(func(name string, i interface{}) {
		switch metric := i.(type) {
		case Counter:
			fmt.Fprintf(w, "%s.%s.count %d %d\n", prefix, name, metric.Count(), now)
		case Gauge:
			fmt.Fprintf(w, "%s.%s.value %d %d\n", prefix, name, metric.Value(), now)
		case Histogram:
			h := metric.Snapshot()
			ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			fmt.Fprintf(w, "%s.%s.count %d %d\n", prefix, name, h.Count(), now)
			fmt.Fprintf(w, "%s.%s.min %d %d\n", prefix, name, h.Min(), now)
			fmt.Fprintf(w, "%s.%s.max %d %d\n", prefix, name, h.Max(), now)
			fmt.Fprintf(w, "%s.%s.mean %.2f %d\n", prefix, name, h.Mean(), now)
			fmt.Fprintf(w, "%s.%s.std-dev %.2f %d\n", prefix, name, h.StdDev(), now)
			fmt.Fprintf(w, "%s.%s.50-percentile %.2f %d\n", prefix, name, ps[0], now)
			fmt.Fprintf(w, "%s.%s.75-percentile %.2f %d\n", prefix, name, ps[1], now)
			fmt.Fprintf(w, "%s.%s.95-percentile %.2f %d\n", prefix, name, ps[2], now)
			fmt.Fprintf(w, "%s.%s.99-percentile %.2f %d\n", prefix, name, ps[3], now)
			fmt.Fprintf(w, "%s.%s.999-percentile %.2f %d\n", prefix, name, ps[4], now)
		case Meter:
			m := metric.Snapshot()
			fmt.Fprintf(w, "%s.%s.count %d %d\n", prefix, name, m.Count(), now)
			fmt.Fprintf(w, "%s.%s.one-minute %.2f %d\n", prefix, name, m.Rate1(), now)
			fmt.Fprintf(w, "%s.%s.five-minute %.2f %d\n", prefix, name, m.Rate5(), now)
			fmt.Fprintf(w, "%s.%s.fifteen-minute %.2f %d\n", prefix, name, m.Rate15(), now)
			fmt.Fprintf(w, "%s.%s.mean %.2f %d\n", prefix, name, m.RateMean(), now)
		case Timer:
			t := metric.Snapshot()
			ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			fmt.Fprintf(w, "%s.%s.count %d %d\n", prefix, name, t.Count(), now)
			fmt.Fprintf(w, "%s.%s.min %d %d\n", prefix, name, t.Min(), now)
			fmt.Fprintf(w, "%s.%s.max %d %d\n", prefix, name, t.Max(), now)
			fmt.Fprintf(w, "%s.%s.mean %.2f %d\n", prefix, name, t.Mean(), now)
			fmt.Fprintf(w, "%s.%s.std-dev %.2f %d\n", prefix, name, t.StdDev(), now)
			fmt.Fprintf(w, "%s.%s.50-percentile %.2f %d\n", prefix, name, ps[0], now)
			fmt.Fprintf(w, "%s.%s.75-percentile %.2f %d\n", prefix, name, ps[1], now)
			fmt.Fprintf(w, "%s.%s.95-percentile %.2f %d\n", prefix, name, ps[2], now)
			fmt.Fprintf(w, "%s.%s.99-percentile %.2f %d\n", prefix, name, ps[3], now)
			fmt.Fprintf(w, "%s.%s.999-percentile %.2f %d\n", prefix, name, ps[4], now)
			fmt.Fprintf(w, "%s.%s.one-minute %.2f %d\n", prefix, name, t.Rate1(), now)
			fmt.Fprintf(w, "%s.%s.five-minute %.2f %d\n", prefix, name, t.Rate5(), now)
			fmt.Fprintf(w, "%s.%s.fifteen-minute %.2f %d\n", prefix, name, t.Rate15(), now)
			fmt.Fprintf(w, "%s.%s.mean-rate %.2f %d\n", prefix, name, t.RateMean(), now)
		}
		w.Flush()
	})
	return nil
}
