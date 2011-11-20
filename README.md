go-metrics
==========

TODO

Memory usage
------------

(Highly unscientific.)

Command used to gather static memory usage:

```sh
grep ^Vm "/proc/$(ps fax | grep [m]etrics-bench | awk '{print $1}')/status"
```

Program used to gather baseline memory usage:

```go
package main

import "time"

func main() {
	time.Sleep(600e9)
}
```

Baseline:

```
VmPeak:	  792544 kB
VmSize:	  792544 kB
VmLck:	       0 kB
VmHWM:	     496 kB
VmRSS:	     496 kB
VmData:	  792024 kB
VmStk:	     136 kB
VmExe:	     376 kB
VmLib:	       0 kB
VmPTE:	      28 kB
VmSwap:	       0 kB
```

Program used to gather metric memory usage (with other metrics being similar):

```go
package main

import (
	"fmt"
	"metrics"
	"time"
)

func main() {
	r := metrics.NewRegistry()
	for i := 0; i < 1000; i++ {
		r.RegisterCounter(fmt.Sprintf("%d", i), metrics.NewCounter())
	}
	time.Sleep(600e9)
}
```

1000 counters registered:

```
VmPeak:   807740 kB
VmSize:   807740 kB
VmLck:         0 kB
VmHWM:      5896 kB
VmRSS:      5896 kB
VmData:   805068 kB
VmStk:       136 kB
VmExe:       912 kB
VmLib:      1580 kB
VmPTE:        48 kB
VmSwap:        0 kB
```

**15 kB virtual, 5 kB resident per counter.**

100000 counters registered:

```
VmPeak:  1204156 kB
VmSize:  1204156 kB
VmLck:         0 kB
VmHWM:    450944 kB
VmRSS:    394756 kB
VmData:  1201484 kB
VmStk:       136 kB
VmExe:       912 kB
VmLib:      1580 kB
VmPTE:       928 kB
VmSwap:    56596 kB
```

**4 kB virtual, 4 kB resident per counter.**

1000 and 100000 gauges registered: negligibly different than counters.

1000 histograms with a uniform sample size of 1028:

```
VmPeak:   811724 kB
VmSize:   811724 kB
VmLck:         0 kB
VmHWM:     15568 kB
VmRSS:     15568 kB
VmData:   809036 kB
VmStk:       136 kB
VmExe:       928 kB
VmLib:      1580 kB
VmPTE:        80 kB
VmSwap:        0 kB
```

**19 kB virtual, 15 kB resident per histogram.**

10000 histograms with a uniform sample size of 1028:

```
VmPeak:   883916 kB
VmSize:   883916 kB
VmLck:         0 kB
VmHWM:    144576 kB
VmRSS:    144576 kB
VmData:   881228 kB
VmStk:       136 kB
VmExe:       928 kB
VmLib:      1580 kB
VmPTE:       432 kB
VmSwap:        0 kB
```

**9 kB virtual, 14 kB resident per histogram.**

50000 histograms with a uniform sample size of 1028:

```
VmPeak:  1204300 kB
VmSize:  1204300 kB
VmLck:         0 kB
VmHWM:    480288 kB
VmRSS:    462244 kB
VmData:  1201612 kB
VmStk:       136 kB
VmExe:       928 kB
VmLib:      1580 kB
VmPTE:      1296 kB
VmSwap:    76464 kB
```

**8 kB virtual, 9 kB resident per histogram.  WTF?**

1000 histograms with an exponentially-decaying sample size of 1028 and alpha of 0.015:

```
VmPeak:   811724 kB
VmSize:   811724 kB
VmLck:         0 kB
VmHWM:     10564 kB
VmRSS:     10564 kB
VmData:   809036 kB
VmStk:       136 kB
VmExe:       928 kB
VmLib:      1580 kB
VmPTE:        52 kB
VmSwap:        0 kB
```

**19 kB virtual, 10 kB resident per histogram.**

10000 histograms with an exponentially-decaying sample size of 1028 and alpha of 0.015:

```
VmPeak:   883788 kB
VmSize:   883788 kB
VmLck:         0 kB
VmHWM:     93484 kB
VmRSS:     93484 kB
VmData:   881100 kB
VmStk:       136 kB
VmExe:       928 kB
VmLib:      1580 kB
VmPTE:       220 kB
VmSwap:        0 kB
```

**9 kB virtual, 9 kB resident per histogram.**

50000 histograms with an exponentially-decaying sample size of 1028 and alpha of 0.015:

```
VmPeak:  1204172 kB
VmSize:  1204172 kB
VmLck:         0 kB
VmHWM:    460360 kB
VmRSS:    460248 kB
VmData:  1201484 kB
VmStk:       136 kB
VmExe:       928 kB
VmLib:      1580 kB
VmPTE:       944 kB
VmSwap:      112 kB
```

**8 kB virtual, 9 kB resident per histogram.  WTF?**

250 meters:

```
VmPeak:  2887084 kB
VmSize:  2887084 kB
VmLck:         0 kB
VmHWM:      3380 kB
VmRSS:      3380 kB
VmData:  2884404 kB
VmStk:       136 kB
VmExe:       920 kB
VmLib:      1580 kB
VmPTE:      1072 kB
VmSwap:        0 kB
```

**8378 kB virtual, 11 kB resident per meter.**
