package metrics

import (

	"log"
	"time"
	"github.com/stathat/go"
)

func Stathat(r Registry, d time.Duration, userkey string) {
	for {
		if err := sh(r, userkey); nil != err {
			log.Println(err)
		}
		time.Sleep(d)
	}
}

func sh(r Registry, userkey string) error {
	r.Each(func(name string, i interface{}) {
		switch m := i.(type) {
		case Counter:
			stathat.PostEZCount(name, userkey, int(m.Count()))
		case Gauge:
			stathat.PostEZValue(name, userkey, float64(m.Value()))
		case Histogram:
			ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			stathat.PostEZCount(name+".count", userkey, int(m.Count()))
			stathat.PostEZValue(name+".min", userkey, float64(m.Min()))
			stathat.PostEZValue(name+".max", userkey, float64(m.Max()))
			stathat.PostEZValue(name+".mean", userkey, float64(m.Mean()))
			stathat.PostEZValue(name+".std-dev", userkey, float64(m.StdDev()))
			stathat.PostEZValue(name+".50-percentile", userkey, float64(ps[0]))
			stathat.PostEZValue(name+".75-percentile", userkey, float64(ps[1]))
			stathat.PostEZValue(name+".95-percentile", userkey, float64(ps[2]))
			stathat.PostEZValue(name+".99-percentile", userkey, float64(ps[3]))
			stathat.PostEZValue(name+".999-percentile", userkey, float64(ps[4]))
		case Meter:
			stathat.PostEZCount(name+".count", userkey, int(m.Count()))
			stathat.PostEZValue(name+".one-minute", userkey, float64(m.Rate1()))
			stathat.PostEZValue(name+".five-minute", userkey, float64(m.Rate5()))
			stathat.PostEZValue(name+".fifteen-minute", userkey, float64(m.Rate15()))
			stathat.PostEZValue(name+".mean", userkey, float64(m.RateMean()))
		case Timer:
			ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			stathat.PostEZCount(name+".count", userkey, int(m.Count()))
			stathat.PostEZValue(name+".min", userkey, float64(m.Min()))
			stathat.PostEZValue(name+".max", userkey, float64(m.Max()))
			stathat.PostEZValue(name+".mean", userkey, float64(m.Mean()))
			stathat.PostEZValue(name+".std-dev", userkey, float64(m.StdDev()))
			stathat.PostEZValue(name+".50-percentile", userkey, float64(ps[0]))
			stathat.PostEZValue(name+".75-percentile", userkey, float64(ps[1]))
			stathat.PostEZValue(name+".95-percentile", userkey, float64(ps[2]))
			stathat.PostEZValue(name+".99-percentile", userkey, float64(ps[3]))
			stathat.PostEZValue(name+".999-percentile", userkey, float64(ps[4]))
			stathat.PostEZValue(name+".one-minute", userkey, float64(m.Rate1()))
			stathat.PostEZValue(name+".five-minute", userkey, float64(m.Rate5()))
			stathat.PostEZValue(name+".fifteen-minute", userkey, float64(m.Rate15()))
			stathat.PostEZValue(name+".mean", userkey, float64(m.RateMean()))
		}

	})
	return nil
}
