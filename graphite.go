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
		switch m := i.(type) {
		case Counter:
			fmt.Fprintf(w, "%s.%s.count %d %d\n", prefix, name, m.Count(), now)
		case Gauge:
			fmt.Fprintf(w, "%s.%s.value %d %d\n", prefix, name, m.Value(), now)
		case Histogram:
			ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			fmt.Fprintf(w, "%s.%s.count %d %d\n", prefix, name, m.Count(), now)
			fmt.Fprintf(w, "%s.%s.min %d %d\n", prefix, name, m.Min(), now)
			fmt.Fprintf(w, "%s.%s.max %d %d\n", prefix, name, m.Max(), now)
			fmt.Fprintf(w, "%s.%s.mean %.2f %d\n", prefix, name, m.Mean(), now)
			fmt.Fprintf(w, "%s.%s.std-dev %.2f %d\n", prefix, name, m.StdDev(), now)
			fmt.Fprintf(w, "%s.%s.50-percentile %.2f %d\n", prefix, name, ps[0], now)
			fmt.Fprintf(w, "%s.%s.75-percentile %.2f %d\n", prefix, name, ps[1], now)
			fmt.Fprintf(w, "%s.%s.95-percentile %.2f %d\n", prefix, name, ps[2], now)
			fmt.Fprintf(w, "%s.%s.99-percentile %.2f %d\n", prefix, name, ps[3], now)
			fmt.Fprintf(w, "%s.%s.999-percentile %.2f %d\n", prefix, name, ps[4], now)
		case Meter:
			fmt.Fprintf(w, "%s.%s.count %d %d\n", prefix, name, m.Count(), now)
			fmt.Fprintf(w, "%s.%s.one-minute %.2f %d\n", prefix, name, m.Rate1(), now)
			fmt.Fprintf(w, "%s.%s.five-minute %.2f %d\n", prefix, name, m.Rate5(), now)
			fmt.Fprintf(w, "%s.%s.fifteen-minute %.2f %d\n", prefix, name, m.Rate15(), now)
			fmt.Fprintf(w, "%s.%s.mean %.2f %d\n", prefix, name, m.RateMean(), now)
		case Timer:
			ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			fmt.Fprintf(w, "%s.%s.count %d %d\n", prefix, name, m.Count(), now)
			fmt.Fprintf(w, "%s.%s.min %d %d\n", prefix, name, m.Min(), now)
			fmt.Fprintf(w, "%s.%s.max %d %d\n", prefix, name, m.Max(), now)
			fmt.Fprintf(w, "%s.%s.mean %.2f %d\n", prefix, name, m.Mean(), now)
			fmt.Fprintf(w, "%s.%s.std-dev %.2f %d\n", prefix, name, m.StdDev(), now)
			fmt.Fprintf(w, "%s.%s.50-percentile %.2f %d\n", prefix, name, ps[0], now)
			fmt.Fprintf(w, "%s.%s.75-percentile %.2f %d\n", prefix, name, ps[1], now)
			fmt.Fprintf(w, "%s.%s.95-percentile %.2f %d\n", prefix, name, ps[2], now)
			fmt.Fprintf(w, "%s.%s.99-percentile %.2f %d\n", prefix, name, ps[3], now)
			fmt.Fprintf(w, "%s.%s.999-percentile %.2f %d\n", prefix, name, ps[4], now)
			fmt.Fprintf(w, "%s.%s.one-minute %.2f %d\n", prefix, name, m.Rate1(), now)
			fmt.Fprintf(w, "%s.%s.five-minute %.2f %d\n", prefix, name, m.Rate5(), now)
			fmt.Fprintf(w, "%s.%s.fifteen-minute %.2f %d\n", prefix, name, m.Rate15(), now)
			fmt.Fprintf(w, "%s.%s.mean-rate %.2f %d\n", prefix, name, m.RateMean(), now)
		}
		w.Flush()
	})
	return nil
}
