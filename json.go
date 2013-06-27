package metrics

import (
	"encoding/json"
)

// MarshalJSON returns a byte slice containing a JSON representation of all
// the metrics in the Registry.
func (r StandardRegistry) MarshalJSON() ([]byte, error) {
	data := make(map[string]map[string]interface{})
	r.Each(func(name string, i interface{}) {
		values := make(map[string]interface{})
		switch m := i.(type) {
		case Counter:
			values["count"] = m.Count()
		case Gauge:
			values["value"] = m.Value()
		case Healthcheck:
			m.Check()
			values["error"] = m.Error()
		case Histogram:
			ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			values["count"] = m.Count()
			values["min"] = m.Min()
			values["max"] = m.Max()
			values["mean"] = m.Mean()
			values["stddev"] = m.StdDev()
			values["median"] = ps[0]
			values["75%%"] = ps[1]
			values["95%%"] = ps[2]
			values["99%%"] = ps[3]
			values["99.9%%"] = ps[4]
		case Meter:
			values["count"] = m.Count()
			values["1m.rate"] = m.Rate1()
			values["5m.rate"] = m.Rate5()
			values["15m.rate"] = m.Rate15()
			values["mean.rate"] = m.RateMean()
		case Timer:
			ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			values["count"] = m.Count()
			values["min"] = m.Min()
			values["max"] = m.Max()
			values["mean"] = m.Mean()
			values["stddev"] = m.StdDev()
			values["median"] = ps[0]
			values["75%%"] = ps[1]
			values["95%%"] = ps[2]
			values["99%%"] = ps[3]
			values["99.9%%"] = ps[4]
			values["1m.rate"] = m.Rate1()
			values["5m.rate"] = m.Rate5()
			values["15m.rate"] = m.Rate15()
			values["mean.rate"] = m.RateMean()
		}
		data[name] = values
	})
	return json.Marshal(data)
}
