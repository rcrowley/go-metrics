package metrics

// Healthchecks hold an os.Error value describing an arbitrary up/down status.
//
// This is an interface so as to encourage other structs to implement
// the Healthcheck API as appropriate.
type Healthcheck interface {
	Check()
	Error() error
	Healthy()
	Unhealthy(error)
}

// Create a new Healthcheck, which will use the given function to update
// its status.
func NewHealthcheck(f func(Healthcheck)) Healthcheck {
	if UseNilMetrics {
		return NilHealthcheck{}
	}
	return &StandardHealthcheck{nil, f}
}

// No-op Healthcheck.
type NilHealthcheck struct{}

// No-op.
func (h NilHealthcheck) Check() {}

// No-op.
func (h NilHealthcheck) Error() error { return nil }

// No-op.
func (h NilHealthcheck) Healthy() {}

// No-op.
func (h NilHealthcheck) Unhealthy(err error) {}

// The standard implementation of a Healthcheck stores the status and a
// function to call to update the status.
type StandardHealthcheck struct {
	err error
	f   func(Healthcheck)
}

// Update the healthcheck's status.
func (h *StandardHealthcheck) Check() {
	h.f(h)
}

// Return the healthcheck's status, which will be nil if it is healthy.
func (h *StandardHealthcheck) Error() error {
	return h.err
}

// Mark the healthcheck as healthy.
func (h *StandardHealthcheck) Healthy() {
	h.err = nil
}

// Mark the healthcheck as unhealthy.  The error should provide details.
func (h *StandardHealthcheck) Unhealthy(err error) {
	h.err = err
}
