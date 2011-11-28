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

// The standard implementation of a Healthcheck stores the status and a
// function to call to update the status.
type StandardHealthcheck struct {
	err error
	f   func(Healthcheck)
}

// Create a new healthcheck, which will use the given function to update
// its status.
func NewHealthcheck(f func(Healthcheck)) *StandardHealthcheck {
	return &StandardHealthcheck{nil, f}
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
