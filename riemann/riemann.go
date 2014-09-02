package riemann

import (
	"fmt"
	"github.com/amir/raidman"
	"github.com/pingles/backoff"
	"github.com/pingles/go-metrics"
	"log"
	"os"
	"path"
	"time"
)

func Riemann(r metrics.Registry, d time.Duration, riemannHost string) {
	var c *raidman.Client
	ch := establishRiemannClient(riemannHost)
	c = <-ch

	for _ = range time.Tick(d) {
		r.Each(func(name string, i interface{}) {
			switch metric := i.(type) {
			case metrics.Gauge:
				e := &raidman.Event{
					Host:    "",
					Service: name,
					Metric:  metric.Value(),
				}
				err := c.Send(e)
				if err != nil {
					log.Println("error sending riemann metric.", err)
					c.Close()
					ch := establishRiemannClient(riemannHost)
					c = <-ch
				}
			case metrics.GaugeFloat64:
				e := &raidman.Event{
					Host:    "",
					Service: name,
					Metric:  metric.Value(),
				}
				err := c.Send(e)
				if err != nil {
					fmt.Println("error sending riemann metric.", err)
					c.Close()
					ch := establishRiemannClient(riemannHost)
					c = <-ch
				}
			case metrics.Counter:
				e := counterEvent(name, metric.Count())

				err := c.Send(e)
				if err != nil {
					log.Println("error sending riemann metric.", err)
					c.Close()
					ch := establishRiemannClient(riemannHost)
					c = <-ch
				}
			case metrics.Meter:
				events := meterEvents(name, metric.Snapshot())
				for _, e := range events {
					err := c.Send(e)
					if err != nil {
						log.Println("error sending riemann metric.", err)
						c.Close()
						ch := establishRiemannClient(riemannHost)
						c = <-ch
					}
				}
			case metrics.Histogram:
				events := histogramEvents(name, metric.Snapshot())
				for _, e := range events {
					err := c.Send(e)
					if err != nil {
						log.Println("error sending riemann metric.", err)
						c.Close()
						ch := establishRiemannClient(riemannHost)
						c = <-ch
					}
				}
			case metrics.Timer:
				events := timerEvents(name, metric.Snapshot())
				for _, e := range events {
					err := c.Send(e)
					if err != nil {
						log.Println("error sending riemann metric.", err)
						c.Close()
						ch := establishRiemannClient(riemannHost)
						c = <-ch
					}
				}
			}
		})
	}
}

func establishRiemannClient(host string) chan *raidman.Client {
	connChannel := make(chan *raidman.Client)

	go func() {
		connect := func() error {
			c, err := raidman.Dial("tcp", host)
			if err != nil {
				log.Println("Error connecting to Riemann, will retry.", err)
				return err
			} else {
				log.Println("connected to riemann server", host)
				connChannel <- c
				return nil
			}
		}

		policy := &backoff.ConstantBackoff{time.Second * 5}
		backoff.Retry(connect, policy)
	}()

	return connChannel
}

func metricName(name string) string {
	return fmt.Sprintf("%s %s", path.Base(os.Args[0]), name)
}

func counterEvent(name string, count int64) *raidman.Event {
	return &raidman.Event{
		Host:    "", // empty string is converted to os.Hostname() by raidman
		Service: metricName(name),
		Metric:  int(count),
	}
}

func meterEvents(name string, metric metrics.Meter) []*raidman.Event {
	return []*raidman.Event{
		event(name, "count", int(metric.Count())),
		event(name, "mean", metric.RateMean()),
		event(name, "one-minute", metric.Rate1()),
		event(name, "five-minute", metric.Rate5()),
		event(name, "fifteen-minute", metric.Rate15()),
	}
}

func timerEvents(name string, metric metrics.Timer) []*raidman.Event {
	events := []*raidman.Event{
		event(name, "count", int(metric.Count())),
		event(name, "min", int(metric.Min())),
		event(name, "max", int(metric.Max())),
		event(name, "mean", metric.Mean()),
		event(name, "std-dev", metric.StdDev()),
		event(name, "one-minute", metric.Rate1()),
		event(name, "five-minute", metric.Rate5()),
		event(name, "fifteen-minute", metric.Rate15()),
	}
	percentiles := []float64{0.75, 0.95, 0.99, 0.999}
	percentileVals := metric.Percentiles(percentiles)
	for i, p := range percentiles {
		e := event(name, fmt.Sprintf("percentile-%.3f", p), percentileVals[i])
		events = append(events, e)
	}
	return events
}

func histogramEvents(name string, metric metrics.Histogram) []*raidman.Event {
	events := []*raidman.Event{
		event(name, "count", int(metric.Count())),
		event(name, "min", int(metric.Min())),
		event(name, "max", int(metric.Max())),
		event(name, "mean", metric.Mean()),
		event(name, "std-dev", metric.StdDev()),
	}

	percentiles := []float64{0.75, 0.95, 0.99, 0.999}
	percentileVals := metric.Percentiles(percentiles)
	for i, p := range percentiles {
		e := event(name, fmt.Sprintf("percentile-%.3f", p), percentileVals[i])
		events = append(events, e)
	}
	return events
}

func event(name string, measure string, val interface{}) *raidman.Event {
	return &raidman.Event{
		Host:    "",
		Service: metricName(fmt.Sprintf("%s.%s", name, measure)),
		Metric:  val,
	}
}
