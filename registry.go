package metrics

import "sync"

type Registry interface{

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

type registry struct {
	mutex *sync.Mutex
	counters map[string]Counter
	gauges map[string]Gauge
	healthchecks map[string]Healthcheck
	histograms map[string]Histogram
	meters map[string]Meter
	timers map[string]Timer
}

func NewRegistry() Registry {
	return &registry {
		&sync.Mutex{},
		make(map[string]Counter),
		make(map[string]Gauge),
		make(map[string]Healthcheck),
		make(map[string]Histogram),
		make(map[string]Meter),
		make(map[string]Timer),
	}
}

func (r *registry) EachCounter(f func(string, Counter)) {
	r.mutex.Lock()
	for name, c := range r.counters { f(name, c) }
	r.mutex.Unlock()
}

func (r *registry) EachGauge(f func(string, Gauge)) {
	r.mutex.Lock()
	for name, g := range r.gauges { f(name, g) }
	r.mutex.Unlock()
}

func (r *registry) EachHealthcheck(f func(string, Healthcheck)) {
	r.mutex.Lock()
	for name, h := range r.healthchecks { f(name, h) }
	r.mutex.Unlock()
}

func (r *registry) EachHistogram(f func(string, Histogram)) {
	r.mutex.Lock()
	for name, h := range r.histograms { f(name, h) }
	r.mutex.Unlock()
}

func (r *registry) EachMeter(f func(string, Meter)) {
	r.mutex.Lock()
	for name, m := range r.meters { f(name, m) }
	r.mutex.Unlock()
}

func (r *registry) EachTimer(f func(string, Timer)) {
	r.mutex.Lock()
	for name, t := range r.timers { f(name, t) }
	r.mutex.Unlock()
}

func (r *registry) RunHealthchecks() {
	r.mutex.Lock()
	for _, h := range r.healthchecks { h.Check() }
	r.mutex.Unlock()
}

func (r *registry) GetCounter(name string) Counter {
	r.mutex.Lock()
	c := r.counters[name]
	r.mutex.Unlock()
	return c
}

func (r *registry) GetGauge(name string) Gauge {
	r.mutex.Lock()
	g := r.gauges[name]
	r.mutex.Unlock()
	return g
}

func (r *registry) GetHealthcheck(name string) Healthcheck {
	r.mutex.Lock()
	h := r.healthchecks[name]
	r.mutex.Unlock()
	return h
}

func (r *registry) GetHistogram(name string) Histogram {
	r.mutex.Lock()
	h := r.histograms[name]
	r.mutex.Unlock()
	return h
}

func (r *registry) GetMeter(name string) Meter {
	r.mutex.Lock()
	m := r.meters[name]
	r.mutex.Unlock()
	return m
}

func (r *registry) GetTimer(name string) Timer {
	r.mutex.Lock()
	t := r.timers[name]
	r.mutex.Unlock()
	return t
}

func (r *registry) RegisterCounter(name string, c Counter) {
	r.mutex.Lock()
	r.counters[name] = c
	r.mutex.Unlock()
}

func (r *registry) RegisterGauge(name string, g Gauge) {
	r.mutex.Lock()
	r.gauges[name] = g
	r.mutex.Unlock()
}

func (r *registry) RegisterHealthcheck(name string, h Healthcheck) {
	r.mutex.Lock()
	r.healthchecks[name] = h
	r.mutex.Unlock()
}

func (r *registry) RegisterHistogram(name string, h Histogram) {
	r.mutex.Lock()
	r.histograms[name] = h
	r.mutex.Unlock()
}

func (r *registry) RegisterMeter(name string, m Meter) {
	r.mutex.Lock()
	r.meters[name] = m
	r.mutex.Unlock()
}

func (r *registry) RegisterTimer(name string, t Timer) {
	r.mutex.Lock()
	r.timers[name] = t
	r.mutex.Unlock()
}

func (r *registry) UnregisterCounter(name string) {
	r.mutex.Lock()
	r.counters[name] = nil, false
	r.mutex.Unlock()
}

func (r *registry) UnregisterGauge(name string) {
	r.mutex.Lock()
	r.gauges[name] = nil, false
	r.mutex.Unlock()
}

func (r *registry) UnregisterHealthcheck(name string) {
	r.mutex.Lock()
	r.healthchecks[name] = nil, false
	r.mutex.Unlock()
}

func (r *registry) UnregisterHistogram(name string) {
	r.mutex.Lock()
	r.histograms[name] = nil, false
	r.mutex.Unlock()
}

func (r *registry) UnregisterMeter(name string) {
	r.mutex.Lock()
	r.meters[name] = nil, false
	r.mutex.Unlock()
}

func (r *registry) UnregisterTimer(name string) {
	r.mutex.Lock()
	r.timers[name] = nil, false
	r.mutex.Unlock()
}
