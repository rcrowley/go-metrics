package main

import (
	"fmt"
	"metrics"
	"time"
)

func main() {

	r := metrics.NewRegistry()

/*
	c := metrics.NewCounter()
	r.RegisterCounter("foo", c)
	for i := 0; i < 1000; i++ {
		go func() {
			for {
				c.Dec(19)
				time.Sleep(300e6)
			}
		}()
		go func() {
			for {
				c.Inc(47)
				time.Sleep(400e6)
			}
		}()
	}
	for {
		fmt.Printf("c.Count(): %v\n", c.Count())
		time.Sleep(500e6)
	}
*/

/*
	g := metrics.NewGauge()
	r.RegisterGauge("bar", g)
	for i := 0; i < 1000; i++ {
		go func() {
			for {
				g.Update(19)
				time.Sleep(300e6)
			}
		}()
		go func() {
			for {
				g.Update(47)
				time.Sleep(400e6)
			}
		}()
	}
	for {
		fmt.Printf("g.Value(): %v\n", g.Value())
		time.Sleep(500e6)
	}
*/

/*
	s := metrics.NewExpDecaySample(1028, 0.015)
//	s := metrics.NewUniformSample(1028)
	h := metrics.NewHistogram(s)
	r.RegisterHistogram("baz", h)
	for i := 0; i < 1000; i++ {
		go func() {
			for {
				h.Update(19)
				time.Sleep(300e6)
			}
		}()
		go func() {
			for {
				h.Update(47)
				time.Sleep(400e6)
			}
		}()
	}
	for {
		fmt.Printf(
			"h: %v %v %v %v %v %v %v %v %v\n",
			h.Count(), h.Sum(), h.Min(), h.Max(),
			h.Percentile(95.0), h.Percentile(99.0), h.Percentile(99.9),
			h.StdDev(), h.Variance(),
		)
		time.Sleep(500e6)
	}
*/

	m := metrics.NewMeter()
	r.RegisterMeter("bang", m)
	for i := 0; i < 1000; i++ {
		go func() {
			for {
				m.Mark(19)
				time.Sleep(300e6)
			}
		}()
		go func() {
			for {
				m.Mark(47)
				time.Sleep(400e6)
			}
		}()
	}
	for {
		fmt.Printf(
			"m: %v %v %v %v %v\n",
			m.Count(), m.Rate1(), m.Rate5(), m.Rate15(), m.RateMean(),
		)
		time.Sleep(500e6)
	}

}
