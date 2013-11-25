package librato

import (
	"fmt"
	"github.com/rcrowley/go-metrics"
	"github.com/samuel/go-librato/librato"
	"log"
	"math"
	"time"
)

type LibratoReporter struct {
	Email, Token string
	Source       string
	Interval     time.Duration
	Registry     metrics.Registry
	Percentiles  []float64 // percentiles to report on histogram metrics
}

func Librato(r metrics.Registry, d time.Duration, e string, t string, s string, p []float64) {
	reporter := &LibratoReporter{e, t, s, d, r, p}
	reporter.Run()
}

func (self *LibratoReporter) Run() {
	ticker := time.Tick(self.Interval)
	metricsApi := &librato.Metrics{self.Email, self.Token}
	for now := range ticker {
		var metrics *librato.MetricsFormat
		var err error
		if metrics, err = self.BuildRequest(now, self.Registry); err != nil {
			log.Printf("ERROR constructing librato request body %s", err)
		}
		if err := metricsApi.SendMetrics(metrics); err != nil {
			log.Printf("ERROR sending metrics to librato %s", err)
		}
	}
}

// calculate sum of squares from data provided by metrics.Histogram
// see http://en.wikipedia.org/wiki/Standard_deviation#Rapid_calculation_methods
func sumSquares(m metrics.Histogram) float64 {
	count := float64(m.Count())
	sum := m.Mean() * float64(m.Count())
	sumSquared := math.Pow(float64(sum), 2)
	sumSquares := math.Pow(count*m.StdDev(), 2) + sumSquared/float64(m.Count())
	if math.IsNaN(sumSquares) {
		return 0.0
	}
	return sumSquared
}
func sumSquaresTimer(m metrics.Timer) float64 {
	count := float64(m.Count())
	sum := m.Mean() * float64(m.Count())
	sumSquared := math.Pow(float64(sum), 2)
	sumSquares := math.Pow(count*m.StdDev(), 2) + sumSquared/float64(m.Count())
	if math.IsNaN(sumSquares) {
		return 0.0
	}
	return sumSquares
}

func (self *LibratoReporter) BuildRequest(now time.Time, r metrics.Registry) (snapshot *librato.MetricsFormat, err error) {
	snapshot = &librato.MetricsFormat{}
	snapshot.MeasureTime = now.Unix()
	snapshot.Source = self.Source
	snapshot.Gauges = make([]interface{}, 0)
	snapshot.Counters = make([]librato.Metric, 0)
	histogramGaugeCount := 1 + len(self.Percentiles)
	r.Each(func(name string, metric interface{}) {
		switch m := metric.(type) {
		case metrics.Counter:
			libratoName := fmt.Sprintf("%s.%s", name, "count")
			snapshot.Counters = append(snapshot.Counters, librato.Metric{Name: libratoName, Value: float64(m.Count())})
		case metrics.Gauge:
			snapshot.Gauges = append(snapshot.Gauges, librato.Metric{Name: name, Value: float64(m.Value())})
		case metrics.Histogram:
			if m.Count() > 0 {
				libratoName := fmt.Sprintf("%s.%s", name, "hist")
				gauges := make([]interface{}, histogramGaugeCount, histogramGaugeCount)
				gauges[0] = librato.Gauge{
					Name:       libratoName,
					Count:      uint64(m.Count()),
					Sum:        m.Mean() * float64(m.Count()),
					Max:        float64(m.Max()),
					Min:        float64(m.Min()),
					SumSquares: sumSquares(m),
				}
				for i, p := range self.Percentiles {
					gauges[i+1] = librato.Metric{Name: fmt.Sprintf("%s.%.2f", libratoName, p), Value: m.Percentile(p)}
				}
				snapshot.Gauges = append(snapshot.Gauges, gauges...)
			}
		case metrics.Meter:
			snapshot.Counters = append(snapshot.Counters, librato.Metric{Name: name, Value: float64(m.Count())})
			snapshot.Gauges = append(snapshot.Gauges,
				librato.Metric{
					Name:  fmt.Sprintf("%s.%s", name, "1min"),
					Value: m.Rate1(),
				},
				librato.Metric{
					Name:  fmt.Sprintf("%s.%s", name, "5min"),
					Value: m.Rate5(),
				},
				librato.Metric{
					Name:  fmt.Sprintf("%s.%s", name, "15min"),
					Value: m.Rate15(),
				},
			)
		case metrics.Timer:
			if m.Count() > 0 {
				libratoName := fmt.Sprintf("%s.%s", name, "timer")
				gauges := make([]interface{}, histogramGaugeCount, histogramGaugeCount)
				gauges[0] = librato.Gauge{
					Name:       libratoName,
					Count:      uint64(m.Count()),
					Sum:        m.Mean() * float64(m.Count()),
					Max:        float64(m.Max()),
					Min:        float64(m.Min()),
					SumSquares: sumSquaresTimer(m),
				}
				for i, p := range self.Percentiles {
					gauges[i+1] = librato.Metric{Name: fmt.Sprintf("%s.%2.0f", libratoName, p*100), Value: m.Percentile(p)}
				}
				snapshot.Gauges = append(snapshot.Gauges, gauges...)
				snapshot.Gauges = append(snapshot.Gauges,
					librato.Metric{
						Name:  fmt.Sprintf("%s.%s", name, "1min"),
						Value: m.Rate1(),
					},
					librato.Metric{
						Name:  fmt.Sprintf("%s.%s", name, "5min"),
						Value: m.Rate5(),
					},
					librato.Metric{
						Name:  fmt.Sprintf("%s.%s", name, "15min"),
						Value: m.Rate15(),
					},
				)
			}
		}
	})
	return
}
