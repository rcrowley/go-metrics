// Hook go-metrics into expvar
// on any /debug/metrics request, load all vars from the registry into expvar, and execute regular expvar handler
package exp

import (
	"expvar"
	"fmt"
	"github.com/rcrowley/go-metrics"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type config struct {
	// If StructuredMetrics is true, each metric type will be exposed
	// as a map with the fields (rates, percentiles, etc...), rather
	// than as individual variables.
	StructuredMetrics bool

	// If RootKey is set to a non-empty string, all metrics will be
	// grouped until a root Map, rather than exposed as top-level
	// metrics. This can make it easier to separate go-metrics from
	// other expvars when processing the output.
	RootKey string

	// Percentiles determines the set of percentile metrics to be
	// exported for each histogram or timer.
	Percentiles []float64

	// TimerScale can be used to scale time-based metrics into
	// different resolutions (i.e time.Milliseconds, etc...). The
	// default TimerScale is nanoseconds (i.e. no conversion).
	TimerScale time.Duration

	percentileLabels []string
}

var (
	Config = newConfig(false, "", []float64{0.5, 0.75, 0.95, 0.99, 0.999})
)

func newConfig(structured bool, root string, ps []float64) *config {
	c := &config{
		StructuredMetrics: structured,
		RootKey:           root,
		Percentiles:       ps,
	}
	c.SetPercentiles(ps)
	return c
}

func (c *config) SetPercentiles(ps []float64) {
	labels := make([]string, len(ps), len(ps))
	for i, p := range ps {
		// Stolen from https://github.com/cyberdelia/go-metrics-graphite/blob/master/graphite.go
		labels[i] = fmt.Sprintf("%s-percentile",
			strings.Replace(strconv.FormatFloat(p*100.0, 'f', -1, 64), ".", "", 1))
	}
	c.percentileLabels = labels
}

type exp struct {
	expvarLock sync.Mutex // expvar panics if you try to register the same var twice, so we must probe it safely
	registry   metrics.Registry
}

func (exp *exp) expHandler(w http.ResponseWriter, r *http.Request) {
	// load our variables into expvar
	exp.syncToExpvar()

	// now just run the official expvar handler code (which is not publicly callable, so pasted inline)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, "{\n")
	first := true
	expvar.Do(func(kv expvar.KeyValue) {
		if !first {
			fmt.Fprintf(w, ",\n")
		}
		first = false
		fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
	})
	fmt.Fprintf(w, "\n}\n")
}

func Exp(r metrics.Registry) {
	e := exp{sync.Mutex{}, r}
	// this would cause a panic:
	// panic: http: multiple registrations for /debug/vars
	// http.HandleFunc("/debug/vars", e.expHandler)
	// haven't found an elegant way, so just use a different endpoint
	http.HandleFunc("/debug/metrics", e.expHandler)
}

func scaleFloat(value float64) float64 {
	return value / float64(Config.TimerScale)
}

func scaleInt(value int64) float64 {
	return float64(value) / float64(Config.TimerScale)
}

func (exp *exp) getInt(name string) *expvar.Int {
	if Config.RootKey != "" {
		m := exp.getMap(Config.RootKey)
		i := newInt(0)
		m.Set(name, i)
		return i
	}
	var v *expvar.Int
	exp.expvarLock.Lock()
	p := expvar.Get(name)
	if p != nil {
		v = p.(*expvar.Int)
	} else {
		v = new(expvar.Int)
		expvar.Publish(name, v)
	}
	exp.expvarLock.Unlock()
	return v
}

func (exp *exp) getFloat(name string) *expvar.Float {
	if Config.RootKey != "" {
		m := exp.getMap(Config.RootKey)
		f := newFloat(0)
		m.Set(name, f)
		return f
	}
	var v *expvar.Float
	exp.expvarLock.Lock()
	p := expvar.Get(name)
	if p != nil {
		v = p.(*expvar.Float)
	} else {
		v = new(expvar.Float)
		expvar.Publish(name, v)
	}
	exp.expvarLock.Unlock()
	return v
}

func (exp *exp) getMap(name string) *expvar.Map {
	var v *expvar.Map
	exp.expvarLock.Lock()
	p := expvar.Get(name)
	if p != nil {
		v = p.(*expvar.Map)
	} else {
		v = new(expvar.Map).Init()
		expvar.Publish(name, v)
	}
	exp.expvarLock.Unlock()
	return v
}

func newMap() *expvar.Map {
	return new(expvar.Map).Init()
}

func newInt(value int64) *expvar.Int {
	i := new(expvar.Int)
	i.Set(value)
	return i
}

func newFloat(value float64) *expvar.Float {
	i := new(expvar.Float)
	i.Set(value)
	return i
}

func (exp *exp) publish(name string, v expvar.Var) {
	exp.expvarLock.Lock()
	p := expvar.Get(name)
	if p == nil {
		expvar.Publish(name, v)
	}
	exp.expvarLock.Unlock()
}

func (exp *exp) publishStructured(name string, value expvar.Var) {
	if Config.RootKey != "" {
		exp.getMap(Config.RootKey).Set(name, value)
	} else {
		exp.publish(name, value)
	}
}

func (exp *exp) publishCounter(name string, metric metrics.Counter) {
	if Config.StructuredMetrics {
		exp.publishStructured(name, convertCounter(metric))
	} else {
		v := exp.getInt(name)
		v.Set(metric.Count())
	}
}

func convertCounter(metric metrics.Counter) *expvar.Map {
	m := newMap()
	m.Set("count", newInt(metric.Count()))
	return m
}

func (exp *exp) publishGauge(name string, metric metrics.Gauge) {
	if Config.StructuredMetrics {
		exp.publishStructured(name, convertGauge(metric))
	} else {
		v := exp.getInt(name)
		v.Set(metric.Value())
	}
}

func convertGauge(metric metrics.Gauge) *expvar.Map {
	m := newMap()
	m.Set("value", newInt(metric.Value()))
	return m
}

func (exp *exp) publishGaugeFloat64(name string, metric metrics.GaugeFloat64) {
	if Config.StructuredMetrics {
		exp.publishStructured(name, convertGaugeFloat64(metric))
	} else {
		exp.getFloat(name).Set(metric.Value())
	}
}

func convertGaugeFloat64(metric metrics.GaugeFloat64) *expvar.Map {
	m := newMap()
	m.Set("value", newFloat(metric.Value()))
	return m
}

func (exp *exp) publishHistogram(name string, metric metrics.Histogram) {
	if Config.StructuredMetrics {
		exp.publishStructured(name, convertHistogram(metric))
	} else {
		h := metric.Snapshot()
		exp.getInt(name + ".count").Set(h.Count())
		exp.getFloat(name + ".min").Set(float64(h.Min()))
		exp.getFloat(name + ".max").Set(float64(h.Max()))
		exp.getFloat(name + ".mean").Set(float64(h.Mean()))
		exp.getFloat(name + ".std-dev").Set(float64(h.StdDev()))
		for i, p := range h.Percentiles(Config.Percentiles) {
			exp.getFloat(fmt.Sprintf("%s.%s", name, Config.percentileLabels[i])).Set(float64(p))
		}
	}
}

func convertHistogram(metric metrics.Histogram) *expvar.Map {
	h := metric.Snapshot()
	m := newMap()
	m.Set("count", newInt(h.Count()))
	m.Set("min", newFloat(float64(h.Min())))
	m.Set("max", newFloat(float64(h.Max())))
	m.Set("std-dev", newFloat(h.StdDev()))
	for i, p := range h.Percentiles(Config.Percentiles) {
		m.Set(Config.percentileLabels[i], newFloat(p))
	}
	return m
}

func (exp *exp) publishMeter(name string, metric metrics.Meter) {
	if Config.StructuredMetrics {
		exp.publishStructured(name, convertMeter(metric))
	} else {
		m := metric.Snapshot()
		exp.getInt(name + ".count").Set(m.Count())
		exp.getFloat(name + ".one-minute").Set(float64(m.Rate1()))
		exp.getFloat(name + ".five-minute").Set(float64(m.Rate5()))
		exp.getFloat(name + ".fifteen-minute").Set(float64((m.Rate15())))
		exp.getFloat(name + ".mean").Set(float64(m.RateMean()))
	}
}

func convertMeter(metric metrics.Meter) *expvar.Map {
	t := metric.Snapshot()
	m := newMap()
	m.Set("count", newInt(t.Count()))
	m.Set("one-minute", newFloat(t.Rate1()))
	m.Set("five-minute", newFloat(t.Rate5()))
	m.Set("fifteen-minute", newFloat(t.Rate15()))
	m.Set("mean", newFloat(t.RateMean()))
	return m
}

func (exp *exp) publishTimer(name string, metric metrics.Timer) {
	if Config.StructuredMetrics {
		exp.publishStructured(name, convertTimer(metric))
	} else {
		t := metric.Snapshot()

		exp.getInt(name + ".count").Set(t.Count())
		exp.getFloat(name + ".min").Set(scaleInt(t.Min()))
		exp.getFloat(name + ".max").Set(scaleInt(t.Max()))
		exp.getFloat(name + ".mean").Set(scaleFloat(t.Mean()))
		exp.getFloat(name + ".std-dev").Set(scaleFloat(t.StdDev()))
		for i, p := range t.Percentiles(Config.Percentiles) {
			exp.getFloat(fmt.Sprintf("%s.%s", name, Config.percentileLabels[i])).Set(scaleFloat(p))
		}
		exp.getFloat(name + ".one-minute").Set(float64(t.Rate1()))
		exp.getFloat(name + ".five-minute").Set(float64(t.Rate5()))
		exp.getFloat(name + ".fifteen-minute").Set(float64((t.Rate15())))
		exp.getFloat(name + ".mean-rate").Set(float64(t.RateMean()))
	}
}

func convertTimer(metric metrics.Timer) *expvar.Map {
	t := metric.Snapshot()
	m := newMap()
	m.Set("count", newInt(t.Count()))
	m.Set("min", newFloat(scaleInt(t.Min())))
	m.Set("max", newFloat(scaleInt(t.Max())))
	m.Set("std-dev", newFloat(scaleFloat(t.StdDev())))
	for i, p := range t.Percentiles(Config.Percentiles) {
		m.Set(Config.percentileLabels[i], newFloat(scaleFloat(p)))
	}
	m.Set("one-minute", newFloat(t.Rate1()))
	m.Set("five-minute", newFloat(t.Rate5()))
	m.Set("fifteen-minute", newFloat(t.Rate15()))
	m.Set("mean-rate", newFloat(t.RateMean()))
	return m
}

func (exp *exp) syncToExpvar() {
	exp.registry.Each(func(name string, i interface{}) {
		switch i.(type) {
		case metrics.Counter:
			exp.publishCounter(name, i.(metrics.Counter))
		case metrics.Gauge:
			exp.publishGauge(name, i.(metrics.Gauge))
		case metrics.GaugeFloat64:
			exp.publishGaugeFloat64(name, i.(metrics.GaugeFloat64))
		case metrics.Histogram:
			exp.publishHistogram(name, i.(metrics.Histogram))
		case metrics.Meter:
			exp.publishMeter(name, i.(metrics.Meter))
		case metrics.Timer:
			exp.publishTimer(name, i.(metrics.Timer))
		default:
			panic(fmt.Sprintf("unsupported type for '%s': %T", name, i))
		}
	})
}
