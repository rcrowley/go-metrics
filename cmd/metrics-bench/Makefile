include $(GOROOT)/src/Make.inc

TARG=metrics-bench
GOFILES=\
	main.go\

include $(GOROOT)/src/Make.cmd

uninstall:
	rm -f $(GOROOT)/bin/$(TARG)

.PHONY: uninstall
