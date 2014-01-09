package metrics

import (
	"net"
	"time"
)

func ExampleGraphite() {
	addr, _ := net.ResolveTCPAddr("net", ":2003")
	go Graphite(DefaultRegistry, 1*time.Second, "some.prefix", addr)
}

func ExampleGraphiteWithConfig() {
	addr, _ := net.ResolveTCPAddr("net", ":2003")
	go GraphiteWithConfig(GraphiteConfig{
		Addr:          addr,
		Registry:      DefaultRegistry,
		FlushInterval: 1 * time.Second,
		DurationUnit:  time.Millisecond,
		RateUnit:      time.Second,
	})
}
