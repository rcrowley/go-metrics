package metrics

import "os"

// Healthchecks hold an os.Error value describing an arbitrary up/down status.
//
// This is an interface so as to encourage other structs to implement
// the Healthcheck API as appropriate.
type Healthcheck interface {
	Check()
	Error() os.Error
	Healthy()
	Unhealthy(os.Error)
}

// The standard implementation of a Healthcheck stores the status and a
// function to call to update the status.
type healthcheck struct {
	err os.Error
	f func(Healthcheck)
}

// Create a new healthcheck, which will use the given function to update
// its status.
func NewHealthcheck(f func(Healthcheck)) Healthcheck {
	return &healthcheck{nil, f}
}

// Update the healthcheck's status.
func (h *healthcheck) Check() {
	h.f(h)
}

// Return the healthcheck's status, which will be nil if it is healthy.
func (h *healthcheck) Error() os.Error {
	return h.err
}

// Mark the healthcheck as healthy.
func (h *healthcheck) Healthy() {
	h.err = nil
}

// Mark the healthcheck as unhealthy.  The error should provide details.
func (h *healthcheck) Unhealthy(err os.Error) {
	h.err = err
}
