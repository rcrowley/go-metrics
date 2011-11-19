include $(GOROOT)/src/Make.inc

TARG=metrics
GOFILES=\
	counter.go\
	ewma.go\
	gauge.go\
	healthcheck.go\
	histogram.go\
	meter.go\
	metrics.go\
	registry.go\
	sample.go\
	timer.go\

include $(GOROOT)/src/Make.pkg

all: uninstall clean install
	make -C cmd/metrics uninstall clean install

uninstall:
	rm -f $(GOROOT)/pkg/$(GOOS)_$(GOARCH)/$(TARG).a
	rm -f $(GOROOT)/pkg/$(GOOS)_$(GOARCH)/github.com/rcrowley/go-$(TARG).a
	rm -rf $(GOROOT)/src/pkg/github.com/rcrowley/go-$(TARG)
	make -C cmd/metrics uninstall

.PHONY: all uninstall
