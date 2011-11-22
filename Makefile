include $(GOROOT)/src/Make.inc

TARG=github.com/rcrowley/go-metrics
GOFILES=\
	counter.go\
	ewma.go\
	gauge.go\
	healthcheck.go\
	histogram.go\
	log.go\
	meter.go\
	metrics.go\
	registry.go\
	runtime.go\
	sample.go\
	syslog.go\
	timer.go\

include $(GOROOT)/src/Make.pkg

all: uninstall clean install
	make -C cmd/metrics-bench uninstall clean install
	make -C cmd/metrics-example uninstall clean install

uninstall:
	rm -f $(GOROOT)/pkg/$(GOOS)_$(GOARCH)/$(TARG).a
	rm -rf $(GOROOT)/src/pkg/$(TARG)
	make -C cmd/metrics-bench uninstall
	make -C cmd/metrics-example uninstall

.PHONY: all uninstall
