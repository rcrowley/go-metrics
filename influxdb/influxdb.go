package influxdb

import (
	"fmt"
	influxClient "github.com/influxdb/influxdb/client"
	"github.com/rcrowley/go-metrics"
	"log"
	"net/url"
	"time"
)

type Config struct {
	Host     string
	Database string
	Username string
	Password string
}

func Influxdb(r metrics.Registry, d time.Duration, config *Config) {
	client, err := NewClient(config)
	if err != nil {
		log.Println(err)
		return
	}

	for _ = range time.Tick(d) {
		if err := Send(r, client, config.Database); err != nil {
			log.Println(err)
		}
	}
}

func NewClient(config *Config) (*influxClient.Client, error) {
	URL, err := url.Parse(config.Host)
	if err != nil {
		return nil, err
	}
	return influxClient.NewClient(influxClient.Config{
		URL:      *URL,
		Username: config.Username,
		Password: config.Password,
	})
}

func Send(r metrics.Registry, client *influxClient.Client, database string) error {
	series := []influxClient.Point{}

	r.Each(func(name string, i interface{}) {
		now := time.Now() // getCurrentTime()
		switch metric := i.(type) {
		case metrics.Counter:
			series = append(series, influxClient.Point{
				Name: fmt.Sprintf("%s.count", name),
				Time: now,
				Fields: map[string]interface{}{
					"count": metric.Count(),
				},
			})
		case metrics.Gauge:
			series = append(series, influxClient.Point{
				Name: fmt.Sprintf("%s.value", name),
				Time: now,
				Fields: map[string]interface{}{
					"value": metric.Value(),
				},
			})
		case metrics.GaugeFloat64:
			series = append(series, influxClient.Point{
				Name: fmt.Sprintf("%s.value", name),
				Time: now,
				Fields: map[string]interface{}{
					"value": metric.Value(),
				},
			})
		case metrics.Histogram:
			h := metric.Snapshot()
			ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			series = append(series, influxClient.Point{
				Name: fmt.Sprintf("%s.histogram", name),
				Time: now,
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
			})
		case metrics.Meter:
			m := metric.Snapshot()
			series = append(series, influxClient.Point{
				Name: fmt.Sprintf("%s.meter", name),
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
			series = append(series, influxClient.Point{
				Name: fmt.Sprintf("%s.timer", name),
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
			})
		}
	})
	if _, err := client.Write(influxClient.BatchPoints{
		Database: database,
		Points:   series,
	}); err != nil {
		log.Println(err)
	}
	return nil
}

func getCurrentTime() int64 {
	return time.Now().UnixNano() / 1000000
}
