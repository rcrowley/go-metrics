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
	promRegistry prometheus.Registerer //Prometheus registry
}

// NewPrometheusProvider returns a Provider that produces Prometheus metrics.
// Namespace and subsystem are applied to all produced metrics.
func NewPrometheusProvider(r Registry, namespace string, subsystem string, promRegistry prometheus.Registerer) *PrometheusConfig{
	return &PrometheusConfig{
		namespace: namespace,
		subsystem: subsystem,
		Registry: r,
		promRegistry: promRegistry,
	}
}

func (c *PrometheusConfig) flattenKey(key string) string {
	key = strings.Replace(key, " ", "_", -1)
	key = strings.Replace(key, ".", "_", -1)
	key = strings.Replace(key, "-", "_", -1)
	key = strings.Replace(key, "=", "_", -1)
	return key
}

func (c *PrometheusConfig) gaugeFromNameAndValue(name string, val float64) prometheus.Gauge{
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: c.flattenKey(c.namespace),
		Subsystem: c.flattenKey(c.subsystem),
		Name:      c.flattenKey(name),
		Help:      name,
	})
	gauge.Set(val)
	err := c.promRegistry.Register(gauge)
	if err != nil {
		return gauge
	}
	return gauge
}

func (c *PrometheusConfig) UpdatePrometheusMetrics() error {
	c.Registry.Each(func(name string, i interface{}) {
		switch metric := i.(type) {
		case Counter:
			c.gaugeFromNameAndValue(name, float64(metric.Count()))
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
