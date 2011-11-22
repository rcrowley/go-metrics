package metrics

import "sync"

// A Registry holds references to a set of metrics by name and can iterate
// over them, calling callback functions provided by the user.
//
// This is an interface so as to encourage other structs to implement
// the Registry API as appropriate.
type Registry interface {

	EachCounter(func(string, Counter))
	EachGauge(func(string, Gauge))
	EachHealthcheck(func(string, Healthcheck))
	EachHistogram(func(string, Histogram))
	EachMeter(func(string, Meter))
	EachTimer(func(string, Timer))

	GetCounter(string) Counter
	GetGauge(string) Gauge
	GetHealthcheck(string) Healthcheck
	GetHistogram(string) Histogram
	GetMeter(string) Meter
	GetTimer(string) Timer

	RegisterCounter(string, Counter)
	RegisterGauge(string, Gauge)
	RegisterHealthcheck(string, Healthcheck)
	RegisterHistogram(string, Histogram)
	RegisterMeter(string, Meter)
	RegisterTimer(string, Timer)

	RunHealthchecks()

	UnregisterCounter(string)
	UnregisterGauge(string)
	UnregisterHealthcheck(string)
	UnregisterHistogram(string)
	UnregisterMeter(string)
	UnregisterTimer(string)

}

// The standard implementation of a Registry is a set of mutex-protected
// maps of names to metrics.
type StandardRegistry struct {
	mutex *sync.Mutex
	counters map[string]Counter
	gauges map[string]Gauge
	healthchecks map[string]Healthcheck
	histograms map[string]Histogram
	meters map[string]Meter
	timers map[string]Timer
}

// Create a new registry.
func NewRegistry() Registry {
	return &StandardRegistry {
		&sync.Mutex{},
		make(map[string]Counter),
		make(map[string]Gauge),
		make(map[string]Healthcheck),
		make(map[string]Histogram),
		make(map[string]Meter),
		make(map[string]Timer),
	}
}

// Call the given function for each registered counter.
func (r *StandardRegistry) EachCounter(f func(string, Counter)) {
	r.mutex.Lock()
	for name, c := range r.counters { f(name, c) }
	r.mutex.Unlock()
}

// Call the given function for each registered gauge.
func (r *StandardRegistry) EachGauge(f func(string, Gauge)) {
	r.mutex.Lock()
	for name, g := range r.gauges { f(name, g) }
	r.mutex.Unlock()
}

// Call the given function for each registered healthcheck.
func (r *StandardRegistry) EachHealthcheck(f func(string, Healthcheck)) {
	r.mutex.Lock()
	for name, h := range r.healthchecks { f(name, h) }
	r.mutex.Unlock()
}

// Call the given function for each registered histogram.
func (r *StandardRegistry) EachHistogram(f func(string, Histogram)) {
	r.mutex.Lock()
	for name, h := range r.histograms { f(name, h) }
	r.mutex.Unlock()
}

// Call the given function for each registered meter.
func (r *StandardRegistry) EachMeter(f func(string, Meter)) {
	r.mutex.Lock()
	for name, m := range r.meters { f(name, m) }
	r.mutex.Unlock()
}

// Call the given function for each registered timer.
func (r *StandardRegistry) EachTimer(f func(string, Timer)) {
	r.mutex.Lock()
	for name, t := range r.timers { f(name, t) }
	r.mutex.Unlock()
}

// Get the Counter by the given name or nil if none is registered.
func (r *StandardRegistry) GetCounter(name string) Counter {
	r.mutex.Lock()
	c := r.counters[name]
	r.mutex.Unlock()
	return c
}

// Get the Gauge by the given name or nil if none is registered.
func (r *StandardRegistry) GetGauge(name string) Gauge {
	r.mutex.Lock()
	g := r.gauges[name]
	r.mutex.Unlock()
	return g
}

// Get the Healthcheck by the given name or nil if none is registered.
func (r *StandardRegistry) GetHealthcheck(name string) Healthcheck {
	r.mutex.Lock()
	h := r.healthchecks[name]
	r.mutex.Unlock()
	return h
}

// Get the Histogram by the given name or nil if none is registered.
func (r *StandardRegistry) GetHistogram(name string) Histogram {
	r.mutex.Lock()
	h := r.histograms[name]
	r.mutex.Unlock()
	return h
}

// Get the Meter by the given name or nil if none is registered.
func (r *StandardRegistry) GetMeter(name string) Meter {
	r.mutex.Lock()
	m := r.meters[name]
	r.mutex.Unlock()
	return m
}

// Get the Timer by the given name or nil if none is registered.
func (r *StandardRegistry) GetTimer(name string) Timer {
	r.mutex.Lock()
	t := r.timers[name]
	r.mutex.Unlock()
	return t
}

// Register the given Counter under the given name.
func (r *StandardRegistry) RegisterCounter(name string, c Counter) {
	r.mutex.Lock()
	r.counters[name] = c
	r.mutex.Unlock()
}

// Register the given Gauge under the given name.
func (r *StandardRegistry) RegisterGauge(name string, g Gauge) {
	r.mutex.Lock()
	r.gauges[name] = g
	r.mutex.Unlock()
}

// Register the given Healthcheck under the given name.
func (r *StandardRegistry) RegisterHealthcheck(name string, h Healthcheck) {
	r.mutex.Lock()
	r.healthchecks[name] = h
	r.mutex.Unlock()
}

// Register the given Histogram under the given name.
func (r *StandardRegistry) RegisterHistogram(name string, h Histogram) {
	r.mutex.Lock()
	r.histograms[name] = h
	r.mutex.Unlock()
}

// Register the given Meter under the given name.
func (r *StandardRegistry) RegisterMeter(name string, m Meter) {
	r.mutex.Lock()
	r.meters[name] = m
	r.mutex.Unlock()
}

// Register the given Timer under the given name.
func (r *StandardRegistry) RegisterTimer(name string, t Timer) {
	r.mutex.Lock()
	r.timers[name] = t
	r.mutex.Unlock()
}

// Run all registered healthchecks.
func (r *StandardRegistry) RunHealthchecks() {
	r.mutex.Lock()
	for _, h := range r.healthchecks { h.Check() }
	r.mutex.Unlock()
}

// Unregister the given Counter with the given name.
func (r *StandardRegistry) UnregisterCounter(name string) {
	r.mutex.Lock()
	r.counters[name] = nil, false
	r.mutex.Unlock()
}

// Unregister the given Gauge with the given name.
func (r *StandardRegistry) UnregisterGauge(name string) {
	r.mutex.Lock()
	r.gauges[name] = nil, false
	r.mutex.Unlock()
}

// Unregister the given Healthcheck with the given name.
func (r *StandardRegistry) UnregisterHealthcheck(name string) {
	r.mutex.Lock()
	r.healthchecks[name] = nil, false
	r.mutex.Unlock()
}

// Unregister the given Histogram with the given name.
func (r *StandardRegistry) UnregisterHistogram(name string) {
	r.mutex.Lock()
	r.histograms[name] = nil, false
	r.mutex.Unlock()
}

// Unregister the given Meter with the given name.
func (r *StandardRegistry) UnregisterMeter(name string) {
	r.mutex.Lock()
	r.meters[name] = nil, false
	r.mutex.Unlock()
}

// Unregister the given Timer with the given name.
func (r *StandardRegistry) UnregisterTimer(name string) {
	r.mutex.Lock()
	r.timers[name] = nil, false
	r.mutex.Unlock()
}
