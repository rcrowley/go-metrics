package metrics

import "os"

type Healthcheck interface {
	Check()
	Error() os.Error
	Healthy()
	Unhealthy(os.Error)
}

type healthcheck struct {
	err os.Error
	f func(Healthcheck)
}

func NewHealthcheck(f func(Healthcheck)) Healthcheck {
	return &healthcheck{nil, f}
}

func (h *healthcheck) Check() {
	h.f(h)
}

func (h *healthcheck) Error() os.Error {
	return h.err
}

func (h *healthcheck) Healthy() {
	h.err = nil
}

func (h *healthcheck) Unhealthy(err os.Error) {
	h.err = err
}
