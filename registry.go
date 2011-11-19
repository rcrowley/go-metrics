package metrics

type Registry interface{

	Counters() map[string]Counter
	Gauges() map[string]Gauge
	Healthchecks() map[string]Healthcheck
	Histograms() map[string]Histogram
	Meters() map[string]Meter
	Timers() map[string]Timer

	GetCounter(string) (Counter, bool)
	GetGauge(string) (Gauge, bool)
	GetHealthcheck(string) (Healthcheck, bool)
	GetHistogram(string) (Histogram, bool)
	GetMeter(string) (Meter, bool)
	GetTimer(string) (Timer, bool)

	RegisterCounter(string, Counter)
	RegisterGauge(string, Gauge)
	RegisterHealthcheck(string, Healthcheck)
	RegisterHistogram(string, Histogram)
	RegisterMeter(string, Meter)
	RegisterTimer(string, Timer)

}

type registry struct {
	counters map[string]Counter
	gauges map[string]Gauge
	healthchecks map[string]Healthcheck
	histograms map[string]Histogram
	meters map[string]Meter
	timers map[string]Timer
}

func NewRegistry() Registry {
	return &registry {
		make(map[string]Counter),
		make(map[string]Gauge),
		make(map[string]Healthcheck),
		make(map[string]Histogram),
		make(map[string]Meter),
		make(map[string]Timer),
	}
}

func (r *registry) Counters() map[string]Counter {
	return r.counters
}

func (r *registry) Gauges() map[string]Gauge {
	return r.gauges
}

func (r *registry) Healthchecks() map[string]Healthcheck {
	return r.healthchecks
}

func (r *registry) Histograms() map[string]Histogram {
	return r.histograms
}

func (r *registry) Meters() map[string]Meter {
	return r.meters
}

func (r *registry) Timers() map[string]Timer {
	return r.timers
}

func (r *registry) GetCounter(name string) (Counter, bool) {
	c, ok := r.counters[name]
	return c, ok
}

func (r *registry) GetGauge(name string) (Gauge, bool) {
	g, ok := r.gauges[name]
	return g, ok
}

func (r *registry) GetHealthcheck(name string) (Healthcheck, bool) {
	h, ok := r.healthchecks[name]
	return h, ok
}

func (r *registry) GetHistogram(name string) (Histogram, bool) {
	h, ok := r.histograms[name]
	return h, ok
}

func (r *registry) GetMeter(name string) (Meter, bool) {
	m, ok := r.meters[name]
	return m, ok
}

func (r *registry) GetTimer(name string) (Timer, bool) {
	t, ok := r.timers[name]
	return t, ok
}

func (r *registry) RegisterCounter(name string, c Counter) {
	r.counters[name] = c
}

func (r *registry) RegisterGauge(name string, g Gauge) {
	r.gauges[name] = g
}

func (r *registry) RegisterHealthcheck(name string, h Healthcheck) {
	r.healthchecks[name] = h
}

func (r *registry) RegisterHistogram(name string, h Histogram) {
	r.histograms[name] = h
}

func (r *registry) RegisterMeter(name string, m Meter) {
	r.meters[name] = m
}

func (r *registry) RegisterTimer(name string, t Timer) {
	r.timers[name] = t
}
