package metrics

import "os"

type Healthcheck interface {
	Check()
	Healthy()
	Unhealthy(os.Error)
}
