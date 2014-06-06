// Go port of Coda Hale's Metrics library
//
// <https://github.com/rcrowley/go-metrics>
//
// Coda Hale's original work: <https://github.com/codahale/metrics>
package metrics

// UseNilMetrics is checked by the constructor functions for all of the
// standard metrics.  If it is true, the metric returned is a stub.
//
// This global kill-switch helps quantify the observer effect and makes
// for less cluttered pprof profiles.
var UseNilMetrics bool = false

type metric struct {
	name string
	m    interface{}
}

// metrics is a slice of metrics that implements sort.Interface
type metrics []metric

func (m metrics) Len() int { return len(m) }

func (m metrics) Swap(i, j int) { m[i], m[j] = m[j], m[i] }

func (m metrics) Less(i, j int) bool {
	return m[i].name < m[j].name
}
