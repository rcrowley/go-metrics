go-metrics
==========

Go port of Coda Hale's Metrics library: <https://github.com/codahale/metrics>.

**This code is not safe on 32-bit architectures.  It will be as soon as `atomic.LoadInt64` lands in a release tag.**

Usage
-----

Create and update metrics:

```go
r := metrics.NewRegistry()

c := metrics.NewCounter()
r.RegisterCounter("foo", c)
c.Inc(47)

g := metrics.NewGauge()
r.RegisterGauge("bar", g)
g.Update(47)

s := metrics.NewExpDecaySample(1028, 0.015)
//s := metrics.NewUniformSample(1028)
h := metrics.NewHistogram(s)
r.RegisterHistogram("baz", h)
h.Update(47)

m := metrics.NewMeter()
r.RegisterMeter("bang", m)
m.Mark(47)
```

Periodically log every metric in human-readable form to standard error:

```go
metrics.Log(r, 60, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
```

Periodically log every metric in slightly-more-parseable form to syslog:

```go
w, err := syslog.Dial("unixgram", "/dev/log", syslog.LOG_INFO, "metrics")
if nil != err { log.Fatalln(err) }
metrics.Syslog(r, 60, w)
```

Installation
------------

```sh
goinstall github.com/rcrowley/go-metrics
```
