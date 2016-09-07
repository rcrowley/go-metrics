package metrics

import (
	"strings"
	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusConfig provides a container with config parameters for the
// Prometheus Exporter

type PrometheusConfig struct {
	namespace string
	Registry Registry // Registry to be exported
	subsystem string
	Percentiles []float64     // Percentiles to export from timers and histograms
}

// NewPrometheusProvider returns a Provider that produces Prometheus metrics.
// Namespace and subsystem are applied to all produced metrics.
func NewPrometheusProvider(r Registry, namespace string, subsystem string) *PrometheusConfig{
	return &PrometheusConfig{
		namespace: namespace,
		subsystem: subsystem,
		Registry: r,
		Percentiles:   []float64{0.5, 0.75, 0.95, 0.99, 0.999},
	}
}

func (c *PrometheusConfig) flattenKey(key string) string {
	key = strings.Replace(key, " ", "_", -1)
	key = strings.Replace(key, ".", "_", -1)
	key = strings.Replace(key, "-", "_", -1)
	key = strings.Replace(key, "=", "_", -1)
	return key
}

func (c *PrometheusConfig) gaugeFromNameAndValue(name string, val float64) {
	gauge := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: c.flattenKey(c.namespace),
		Subsystem: c.flattenKey(c.subsystem),
		Name:      c.flattenKey(name),
		Help:      name,
	})
	prometheus.MustRegister(gauge)
	gauge.Set(val)
}

func (c *PrometheusConfig) update_prometheus_metrics() error {
	c.Registry.Each(func(name string, i interface{}) {
		switch metric := i.(type) {
		case Counter:
 			cntr := prometheus.NewCounter(prometheus.CounterOpts{
				Namespace: c.flattenKey(c.namespace),
				Subsystem: c.flattenKey(c.subsystem),
				Name:      c.flattenKey(name),
				Help:      name,
			})
			prometheus.MustRegister(cntr)
			cntr.Set(float64(metric.Count()))
		case Gauge:
		case GaugeFloat64:
			c.gaugeFromNameAndValue(name, float64(metric.Value()))
		case Histogram:
			samples := metric.Snapshot().Sample().Values()
			lastSample :=  samples[len(samples)-1]
			c.gaugeFromNameAndValue(name, float64(lastSample))
		case Meter:
		case Timer:
			lastSample := metric.Snapshot().Rate1()
			c.gaugeFromNameAndValue(name, lastSample)
		}
	})
	return nil
}
