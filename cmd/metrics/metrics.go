package main

import (
//	"log"
	"metrics"
//	"os"
//	"syslog"
	"time"
)

const fanout = 10

func main() {

	r := metrics.NewRegistry()

	c := metrics.NewCounter()
	r.RegisterCounter("foo", c)
	for i := 0; i < fanout; i++ {
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

	g := metrics.NewGauge()
	r.RegisterGauge("bar", g)
	for i := 0; i < fanout; i++ {
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

	s := metrics.NewExpDecaySample(1028, 0.015)
//	s := metrics.NewUniformSample(1028)
	h := metrics.NewHistogram(s)
	r.RegisterHistogram("baz", h)
	for i := 0; i < fanout; i++ {
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

	m := metrics.NewMeter()
	r.RegisterMeter("bang", m)
	for i := 0; i < fanout; i++ {
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

//	metrics.Log(r, 60, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

/*
	w, err := syslog.Dial("unixgram", "/dev/log", syslog.LOG_INFO, "metrics")
	if nil != err { log.Fatalln(err) }
	metrics.Syslog(r, 60, w)
*/

}
