package metrics

type Registry interface{
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
