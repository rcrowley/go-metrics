go-metrics
==========

Go port of Coda Hale's Metrics library: <https://github.com/codahale/metrics>.

Usage
-----

Create and update metrics:

```go
c := metrics.NewCounter()
metrics.Register("foo", c)
c.Inc(47)

g := metrics.NewGauge()
metrics.Register("bar", g)
g.Update(47)

s := metrics.NewExpDecaySample(1028, 0.015) // or metrics.NewUniformSample(1028)
h := metrics.NewHistogram(s)
metrics.Register("baz", h)
h.Update(47)

m := metrics.NewMeter()
metrics.RegisterMeter("quux", m)
m.Mark(47)

t := metrics.NewTimer()
metrics.RegisterTimer("bang", t)
t.Time(func() {})
t.Update(47)
```

Periodically log every metric in human-readable form to standard error:

```go
metrics.Log(metrics.DefaultRegistry, 60, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
```

Periodically log every metric in slightly-more-parseable form to syslog:

```go
w, _ := syslog.Dial("unixgram", "/dev/log", syslog.LOG_INFO, "metrics")
metrics.Syslog(metrics.DefaultRegistry, 60, w)
```

Periodically emit every metric to Graphite:

```go
addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:2003")
metrics.Graphite(metrics.DefaultRegistry, 10, "metrics", addr)
```

Installation
------------

```sh
go get github.com/rcrowley/go-metrics
```
