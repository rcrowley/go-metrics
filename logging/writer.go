package logging

import (
	"github.com/zeim839/go-metrics++"
	"io"
	"sort"
	"time"
)

// Write sorts & writes each metric in the given registry periodically to the
// given io.Writer, in the prometheus expositional format.
func Write(r metrics.Registry, d time.Duration, w io.Writer) {
	for _ = range time.Tick(d) {
		WriteOnce(r, w)
	}
}

// WriteOnce sorts & writes metrics in the given registry to the given
// io.Writer, in the prometheus expositional format.
func WriteOnce(r metrics.Registry, w io.Writer) {
	var namedMetrics namedMetricSlice
	r.Each(func(name string, i interface{}) {
		namedMetrics = append(namedMetrics, namedMetric{name, i})
	})

	sort.Sort(namedMetrics)
	for _, namedMetric := range namedMetrics {
		w.Write([]byte(Encode(namedMetric.name, "", namedMetric.m)))
	}
}

type namedMetric struct {
	name string
	m    interface{}
}

// namedMetricSlice is a slice of namedMetrics that implements sort.Interface.
type namedMetricSlice []namedMetric

func (nms namedMetricSlice) Len() int { return len(nms) }

func (nms namedMetricSlice) Swap(i, j int) { nms[i], nms[j] = nms[j], nms[i] }

func (nms namedMetricSlice) Less(i, j int) bool {
	return nms[i].name < nms[j].name
}
