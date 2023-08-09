package logging

import (
	"fmt"
	"github.com/zeim839/go-metrics++"
	"time"
)

// Encode encodes a metric into a string in prometheus expositional format.
// Since the format does not support all go-metrics++ metric types, some
// interfaces are encoded over multiple lines, e.g: a Timer will have one line
// for its count, another line for its mean, etc. Returns "" if the interface
// is not a metric. Healthchecks are not supported.
func Encode(name, prefix string, i interface{}) string {
	if prefix != "" {
		prefix = prefix + "_"
	}
	head := prefix + name
	ts := time.Now().UTC().Unix()

	switch metric := i.(type) {
	case metrics.Counter:
		return fmt.Sprintf("%s%s %d %v\n", head, EncodeLabels(metric.Labels()),
			metric.Count(), ts)
	case metrics.Gauge:
		return fmt.Sprintf("%s%s %d %v\n", head, EncodeLabels(metric.Labels()),
			metric.Value(), ts)
	case metrics.GaugeFloat64:
		return fmt.Sprintf("%s%s %f %v\n", head, EncodeLabels(metric.Labels()),
			metric.Value(), ts)
	case metrics.Meter:
		m := metric.Snapshot()
		labels := EncodeLabels(m.Labels())
		str := fmt.Sprintf("%s_count%s %d %v\n", head, labels, m.Count(), ts)
		str += fmt.Sprintf("%s_rate_1min%s %f %v\n", head, labels, m.Rate1(), ts)
		str += fmt.Sprintf("%s_rate_5min%s %f %v\n", head, labels, m.Rate5(), ts)
		str += fmt.Sprintf("%s_rate_15min%s %f %v\n", head, labels, m.Rate15(), ts)
		str += fmt.Sprintf("%s_rate_mean%s %f %v\n", head, labels, m.RateMean(), ts)
		return str
	case metrics.Timer:
		t := metric.Snapshot()
		labels := EncodeLabels(t.Labels())
		ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
		str := fmt.Sprintf("%s_count%s %d %v\n", head, labels, t.Count(), ts)
		str += fmt.Sprintf("%s_min%s %d %v\n", head, labels, t.Min(), ts)
		str += fmt.Sprintf("%s_max%s %d %v\n", head, labels, t.Max(), ts)
		str += fmt.Sprintf("%s_mean%s %f %v\n", head, labels, t.Mean(), ts)
		str += fmt.Sprintf("%s_stddev%s %f %v\n", head, labels, t.StdDev(), ts)
		str += fmt.Sprintf("%s_median%s %f %v\n", head, labels, ps[0], ts)
		str += fmt.Sprintf("%s_percentile_75%s %f %v\n", head, labels, ps[1], ts)
		str += fmt.Sprintf("%s_percentile_95%s %f %v\n", head, labels, ps[2], ts)
		str += fmt.Sprintf("%s_percentile_99%s %f %v\n", head, labels, ps[3], ts)
		str += fmt.Sprintf("%s_percentile_99_9%s %f %v\n", head, labels, ps[4], ts)
		str += fmt.Sprintf("%s_rate_1min%s %f %v\n", head, labels, t.Rate1(), ts)
		str += fmt.Sprintf("%s_rate_5min%s %f %v\n", head, labels, t.Rate5(), ts)
		str += fmt.Sprintf("%s_rate_15min%s %f %v\n", head, labels, t.Rate15(), ts)
		str += fmt.Sprintf("%s_rate_mean%s %f %v\n", head, labels, t.RateMean(), ts)
		return str
	case metrics.Histogram:
		h := metric.Snapshot()
		labels := EncodeLabels(h.Labels())
		ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
		str := fmt.Sprintf("%s_count%s %d %v\n", head, labels, h.Count(), ts)
		str += fmt.Sprintf("%s_min%s %d %v\n", head, labels, h.Min(), ts)
		str += fmt.Sprintf("%s_max%s %d %v\n", head, labels, h.Max(), ts)
		str += fmt.Sprintf("%s_mean%s %f %v\n", head, labels, h.Mean(), ts)
		str += fmt.Sprintf("%s_stddev%s %f %v\n", head, labels, h.StdDev(), ts)
		str += fmt.Sprintf("%s_median%s %f %v\n", head, labels, ps[0], ts)
		str += fmt.Sprintf("%s_percentile_75%s %f %v\n", head, labels, ps[1], ts)
		str += fmt.Sprintf("%s_percentile_95%s %f %v\n", head, labels, ps[2], ts)
		str += fmt.Sprintf("%s_percentile_99%s %f %v\n", head, labels, ps[3], ts)
		str += fmt.Sprintf("%s_percentile_99_9%s %f %v\n", head, labels, ps[4], ts)
		return str
	}

	return ""
}

// EncodeLabels encodes labels into JSON format. Returns "" if the
// slice is empty.
func EncodeLabels(labels []metrics.Label) string {
	if labels == nil || len(labels) < 1 {
		return ""
	}
	str := "{"
	for _, label := range labels {
		str += label.Key + ":\"" + label.Value + "\","
	}
	// Remove last comma character and add closing brace.
	return str[:len(str)-1] + "}"
}
