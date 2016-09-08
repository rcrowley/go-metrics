package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"testing"
)

func TestPrometheusRegistration(t *testing.T) {
	defaultRegistry := prometheus.DefaultRegisterer
	pClient := NewPrometheusProvider(DefaultRegistry, "test", "subsys", defaultRegistry)
	if pClient.promRegistry != defaultRegistry {
		t.Fatalf("Failed to pass prometheus registry to go-metrics provider")
	}
}

func TestMetricsGetAddedToPromRegistry(t *testing.T) {
	prometheusRegistry := prometheus.NewRegistry()
	metricsRegistry := NewRegistry()
	pClient := NewPrometheusProvider(metricsRegistry, "test", "subsys", prometheusRegistry)
	metricsRegistry.Register("counter", NewCounter())
	pClient.UpdatePrometheusMetrics()
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "test",
		Subsystem: "subsys",
		Name:      "counter",
		Help:      "counter",
	})
	err := prometheusRegistry.Register(gauge)
	if err == nil {
		t.Fatalf("Go-metrics registry didn't get registered to prometheus registry")
	}

}
