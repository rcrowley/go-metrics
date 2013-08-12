package metrics

import "sync"

// A Registry holds references to a set of metrics by name and can iterate
// over them, calling callback functions provided by the user.
//
// This is an interface so as to encourage other structs to implement
// the Registry API as appropriate.
type Registry interface {

	// Call the given function for each registered metric.
	Each(func(string, interface{}))

	// Get the metric by the given name or nil if none is registered.
	Get(string) interface{}

	// Register the given metric under the given name.
	Register(string, interface{})

	// Run all registered healthchecks.
	RunHealthchecks()

	// Unregister the metric with the given name.
	Unregister(string)
}

// The standard implementation of a Registry is a mutex-protected map
// of names to metrics.
type StandardRegistry struct {
	metrics map[string]interface{}
	mutex   sync.Mutex
}

// Force the compiler to check that StandardRegistry implements Registry.
var _ Registry = &StandardRegistry{}

// Create a new registry.
func NewRegistry() *StandardRegistry {
	return &StandardRegistry{metrics: make(map[string]interface{})}
}

// Call the given function for each registered metric.
func (r *StandardRegistry) Each(f func(string, interface{})) {
	for name, i := range r.registered() {
		f(name, i)
	}
}

// Get the metric by the given name or nil if none is registered.
func (r *StandardRegistry) Get(name string) interface{} {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.metrics[name]
}

// Register the given metric under the given name.
func (r *StandardRegistry) Register(name string, i interface{}) {
	switch i.(type) {
	case Counter, Gauge, Healthcheck, Histogram, Meter, Timer:
		r.mutex.Lock()
		defer r.mutex.Unlock()
		r.metrics[name] = i
	}
}

// Run all registered healthchecks.
func (r *StandardRegistry) RunHealthchecks() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, i := range r.metrics {
		if h, ok := i.(Healthcheck); ok {
			h.Check()
		}
	}
}

// Unregister the metric with the given name.
func (r *StandardRegistry) Unregister(name string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.metrics, name)
}

func (r *StandardRegistry) registered() map[string]interface{} {
	metrics := make(map[string]interface{}, len(r.metrics))
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for name, i := range r.metrics {
		metrics[name] = i
	}
	return metrics
}

var DefaultRegistry *StandardRegistry = NewRegistry()

// Call the given function for each registered metric.
func Each(f func(string, interface{})) {
	DefaultRegistry.Each(f)
}

// Get the metric by the given name or nil if none is registered.
func Get(name string) interface{} {
	return DefaultRegistry.Get(name)
}

// Register the given metric under the given name.
func Register(name string, i interface{}) {
	DefaultRegistry.Register(name, i)
}

// Run all registered healthchecks.
func RunHealthchecks() {
	DefaultRegistry.RunHealthchecks()
}

// Unregister the metric with the given name.
func Unregister(name string) {
	DefaultRegistry.Unregister(name)
}
