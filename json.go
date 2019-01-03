package metrics

import (
	"encoding/json"
	"io"
	"time"
	"fmt"
)

// MarshalJSON returns a byte slice containing a JSON representation of all
// the metrics in the Registry.
func (r *StandardRegistry) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.GetAll())
}


// MarshalJSONStringified returns a byte slice containing a JSON representation of all
// the metrics in the Registry with values stringified.
func (r *StandardRegistry) MarshalJSONStringified(scale time.Duration) ([]byte, error) {
	du := float64(scale)
	duSuffix := scale.String()[1:]

	data := make(map[string]map[string]interface{})
	r.Each(func(name string, i interface{}) {
		values := make(map[string]interface{})
		switch metric := i.(type) {
		case Counter:
			values["count"] = fmt.Sprintf("%9d", metric.Count())
		case Gauge:
			values["value"] = fmt.Sprintf("%9d", metric.Value())
		case GaugeFloat64:
			values["value"] = fmt.Sprintf("%f", metric.Value())
		case Healthcheck:
			values["error"] = nil
			metric.Check()
			if err := metric.Error(); nil != err {
				values["error"] = metric.Error().Error()
			}
		case Histogram:
			h := metric.Snapshot()
			ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			values["count"] =  fmt.Sprintf("%9d", h.Count())
			values["min"] =    fmt.Sprintf("%9d", h.Min())
			values["max"] =    fmt.Sprintf("%9d", h.Max())
			values["mean"] =   fmt.Sprintf("%12.2f", h.Mean())
			values["stddev"] = fmt.Sprintf("%12.2f", h.StdDev())
			values["median"] = fmt.Sprintf("%12.2f", ps[0])
			values["75%"] =    fmt.Sprintf("%12.2f", ps[1])
			values["95%"] =    fmt.Sprintf("%12.2f", ps[2])
			values["99%"] =    fmt.Sprintf("%12.2f", ps[3])
			values["99.9%"] =  fmt.Sprintf("%12.2f", ps[4])
		case Meter:
			m := metric.Snapshot()
			values["count"] =     fmt.Sprintf("%9d", m.Count())
			values["1m.rate"] =   fmt.Sprintf("%12.2f", m.Rate1())
			values["5m.rate"] =   fmt.Sprintf("%12.2f", m.Rate5())
			values["15m.rate"] =  fmt.Sprintf("%12.2f", m.Rate15())
			values["mean.rate"] = fmt.Sprintf("%12.2f", m.RateMean())
		case Timer:
			t := metric.Snapshot()
			ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			values["count"] =     fmt.Sprintf("%9d", t.Count())
			values["min"] =       fmt.Sprintf("%12.2f%s", float64(t.Min())/du, duSuffix)
			values["max"] =       fmt.Sprintf("%12.2f%s", float64(t.Max())/du, duSuffix)
			values["mean"] =      fmt.Sprintf("%12.2f%s", t.Mean()/du, duSuffix)
			values["stddev"] =    fmt.Sprintf("%12.2f%s", t.StdDev()/du, duSuffix)
			values["median"] =    fmt.Sprintf("%12.2f%s", ps[0]/du, duSuffix)
			values["75%"] =       fmt.Sprintf("%12.2f%s", ps[1]/du, duSuffix)
			values["95%"] =       fmt.Sprintf("%12.2f%s", ps[2]/du, duSuffix)
			values["99%"] =       fmt.Sprintf("%12.2f%s", ps[3]/du, duSuffix)
			values["99.9%"] =     fmt.Sprintf("%12.2f%s", ps[4]/du, duSuffix)
			values["1m.rate"] =   fmt.Sprintf("%12.2f", t.Rate1())
			values["5m.rate"] =   fmt.Sprintf("%12.2f", t.Rate5())
			values["15m.rate"] =  fmt.Sprintf("%12.2f", t.Rate15())
			values["mean.rate"] = fmt.Sprintf("%12.2f", t.RateMean())
		}
		data[name] = values
	})
	return json.Marshal(data)
}

// WriteJSON writes metrics from the given registry  periodically to the
// specified io.Writer as JSON.
func WriteJSON(r Registry, d time.Duration, w io.Writer) {
	for _ = range time.Tick(d) {
		WriteJSONOnce(r, w)
	}
}

// WriteJSONOnce writes metrics from the given registry to the specified
// io.Writer as JSON.
func WriteJSONOnce(r Registry, w io.Writer) {
	json.NewEncoder(w).Encode(r)
}

func (p *PrefixedRegistry) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.GetAll())
}
