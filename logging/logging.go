package logging

import (
	"github.com/zeim839/go-metrics++"
	"os"
	"time"
)

// Logger is a block exporter function which flushes metrics in r
// to stdout in prometheus exposition format, sinking them every
// d duration and prepending metric names with prefix.
func Logger(r metrics.Registry, d time.Duration, prefix string) {
	for _ = range time.Tick(d) {
		r.Each(func(name string, i interface{}) {
			os.Stdout.Write([]byte(Encode(name, prefix, i)))
		})
	}
}
