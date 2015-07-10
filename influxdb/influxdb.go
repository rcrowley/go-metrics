package influxdb

import (
	"fmt"
	"log"
	"net/url"
	"time"

	influx "github.com/influxdb/influxdb/client"
	"github.com/rcrowley/go-metrics"
)

type Config struct {
	Host      string
	Database  string
	Username  string
	Password  string
	UserAgent string
	Timeout   time.Duration
}

func Influxdb(r metrics.Registry, d time.Duration, config *Config) {
	hostURL, err := url.Parse(config.Host)
	if err != nil {
		log.Println(err)
		return
	}
	client, err := influx.NewClient(influx.Config{
		URL:       *hostURL,
		Username:  config.Username,
		Password:  config.Password,
		UserAgent: config.UserAgent,
		Timeout:   config.Timeout,
	})
	if err != nil {
		log.Println(err)
		return
	}

	for _ = range time.Tick(d) {
		if err := send(client, config.Database, r); err != nil {
			log.Println(err)
		}
	}
}

func send(client *influx.Client, database string, r metrics.Registry) error {

	var (
		points []influx.Point
		now    = time.Now()
	)

	r.Each(func(name string, i interface{}) {

		switch metric := i.(type) {
		case metrics.Counter:
			points = append(points, influx.Point{
				Measurement: fmt.Sprintf("%s.count", name),
				Time:        now,
				Fields: map[string]interface{}{
					"count": metric.Count(),
				},
			})
		case metrics.Gauge:
			points = append(points, influx.Point{
				Measurement: fmt.Sprintf("%s.value", name),
				Time:        now,
				Fields: map[string]interface{}{
					"value": metric.Value(),
				},
			})
		case metrics.GaugeFloat64:
			points = append(points, influx.Point{
				Measurement: fmt.Sprintf("%s.value", name),
				Time:        now,
				Fields: map[string]interface{}{
					"value": metric.Value(),
				},
			})
		case metrics.Histogram:
			sn := metric.Snapshot()
			ps := sn.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			points = append(points, influx.Point{
				Measurement: fmt.Sprintf("%s.histogram", name),
				Time:        now,
				Fields: map[string]interface{}{
					"count":          sn.Count,
					"min":            sn.Min(),
					"max":            sn.Max(),
					"mean":           sn.Mean(),
					"std-dev":        sn.StdDev(),
					"50-percentile":  ps[0],
					"75-percentile":  ps[1],
					"95-percentile":  ps[2],
					"99-percentile":  ps[3],
					"999-percentile": ps[4],
				},
			})
		case metrics.Meter:
			sn := metric.Snapshot()
			points = append(points, influx.Point{
				Measurement: fmt.Sprintf("%s.meter", name),
				Time:        now,
				Fields: map[string]interface{}{
					"count":          sn.Count(),
					"one-minute":     sn.Rate1(),
					"five-minute":    sn.Rate5(),
					"fifteen-minute": sn.Rate15(),
					"mean":           sn.RateMean(),
				},
			})
		case metrics.Timer:
			sn := metric.Snapshot()
			ps := sn.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			points = append(points, influx.Point{
				Measurement: fmt.Sprintf("%s.timer", name),
				Fields: map[string]interface{}{
					"count":          sn.Count(),
					"min":            sn.Min(),
					"max":            sn.Max(),
					"mean":           sn.Mean(),
					"std-dev":        sn.StdDev(),
					"50-percentile":  ps[0],
					"75-percentile":  ps[1],
					"95-percentile":  ps[2],
					"99-percentile":  ps[3],
					"999-percentile": ps[4],
					"one-minute":     sn.Rate1(),
					"five-minute":    sn.Rate5(),
					"fifteen-minute": sn.Rate15(),
					"mean-rate":      sn.RateMean(),
				},
			})
		}
	})

	batch := influx.BatchPoints{
		Points:   points,
		Database: database,
		Time:     now,
	}

	if response, err := client.Write(batch); err != nil {
		log.Println(err)
	} else if err := response.Error(); err != nil {
		log.Println(err)
	}

	return nil
}
