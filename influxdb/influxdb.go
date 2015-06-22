package influxdb

import (
	"fmt"
	"log"
	"net/url"
	"time"

	influxClient "github.com/influxdb/influxdb/client"
	"github.com/rcrowley/go-metrics"
)

type Config struct {
	Host     string
	Database string
	Username string
	Password string
}

func Influxdb(r metrics.Registry, d time.Duration, config *Config) {
	url, err := url.Parse(fmt.Sprintf("http://%s", config.Host))
	if err != nil {
		log.Println(err)
		return
	}

	client, err := influxClient.NewClient(influxClient.Config{
		URL:      *url,
		Username: config.Username,
		Password: config.Password,
	})

	if err != nil {
		log.Println(err)
		return
	}

	for _ = range time.Tick(d) {
		if err := send(r, client, config.Database); err != nil {
			log.Println(err)
		}
	}
}

func send(r metrics.Registry, client *influxClient.Client, database string) error {
	var points []influxClient.Point
	now := time.Now()

	r.Each(func(name string, i interface{}) {
		switch metric := i.(type) {
		case metrics.Counter:
			points = append(points, influxClient.Point{
				Measurement: fmt.Sprintf("%s.count", name),
				Fields: map[string]interface{}{
					"value": metric.Count(),
				},
				Time: now,
			})
		case metrics.Gauge:
			points = append(points, influxClient.Point{
				Measurement: fmt.Sprintf("%s.value", name),
				Fields: map[string]interface{}{
					"value": metric.Value(),
				},
				Time: now,
			})
		case metrics.GaugeFloat64:
			points = append(points, influxClient.Point{
				Measurement: fmt.Sprintf("%s.value", name),
				Fields: map[string]interface{}{
					"value": metric.Value(),
				},
				Time: now,
			})
		case metrics.Histogram:
			h := metric.Snapshot()
			ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})

			points = append(points, influxClient.Point{
				Measurement: fmt.Sprintf("%s.value", name),
				Fields: map[string]interface{}{
					"count":          h.Count(),
					"min":            h.Min(),
					"max":            h.Max(),
					"mean":           h.Mean(),
					"std-dev":        h.StdDev(),
					"50-percentile":  ps[0],
					"75-percentile":  ps[1],
					"95-percentile":  ps[2],
					"99-percentile":  ps[3],
					"999-percentile": ps[4],
				},
				Time: now,
			})
		case metrics.Meter:
			m := metric.Snapshot()

			points = append(points, influxClient.Point{
				Measurement: fmt.Sprintf("%s.meter", name),
				Fields: map[string]interface{}{
					"count":          m.Count(),
					"one-minute":     m.Rate1(),
					"five-minute":    m.Rate5(),
					"fifteen-minute": m.Rate15(),
					"mean":           m.RateMean(),
				},
			})
		case metrics.Timer:
			h := metric.Snapshot()
			ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})

			points = append(points, influxClient.Point{
				Measurement: fmt.Sprintf("%s.timer", name),
				Fields: map[string]interface{}{
					"count":          h.Count(),
					"min":            h.Min(),
					"max":            h.Max(),
					"mean":           h.Mean(),
					"std-dev":        h.StdDev(),
					"50-percentile":  ps[0],
					"75-percentile":  ps[1],
					"95-percentile":  ps[2],
					"99-percentile":  ps[3],
					"999-percentile": ps[4],
					"one-minute":     h.Rate1(),
					"five-minute":    h.Rate5(),
					"fifteen-minute": h.Rate15(),
					"mean-rate":      h.RateMean(),
				},
				Time: now,
			})
		}
	})

	bps := influxClient.BatchPoints{
		Points:   points,
		Database: database,
	}

	_, err := client.Write(bps)
	return err
}

func getCurrentTime() int64 {
	return time.Now().UnixNano() / 1000000
}
