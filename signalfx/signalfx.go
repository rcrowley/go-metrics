package signalfx

import (
	"context"
	"fmt"
	metrics "github.com/rcrowley/go-metrics"
	"github.com/signalfx/golib/datapoint"
	"github.com/signalfx/golib/sfxclient"
	"time"
)

// PublishToSignalFx publishes periodically all the metrics of the specified
// registry to SignalFX (https://signalfx.com/). This is designed to be called
// as a goroutine. Providing a logger is optional and only used to report
// publishing errors.
func PublishToSignalFx(r metrics.Registry, d time.Duration, logger metrics.Logger, authToken string) {
	publisher := publisher{authToken: authToken}
	for _ = range time.Tick(d) {
		if err := publisher.single(r); err != nil {
			publisher.client = nil
			if logger != nil {
				logger.Printf("Unable to publish to SignalFX: %s.", err)
			}
		}
	}
}

type publisher struct {
	authToken string
	client    *sfxclient.HTTPSink
}

func (p publisher) single(r metrics.Registry) error {
	if p.client == nil {
		p.client = sfxclient.NewHTTPSink()
		p.client.AuthToken = p.authToken
	}
	ctx := context.Background()
	var datapoints []*datapoint.Datapoint
	r.Each(func(name string, i interface{}) {
		ds := metricToDatapoints(name, i)
		datapoints = append(datapoints, ds...)
	})
	return p.client.AddDatapoints(ctx, datapoints)
}

func metricToDatapoints(name string, i interface{}) []*datapoint.Datapoint {
	switch metric := i.(type) {
	case metrics.Counter:
		return []*datapoint.Datapoint{
			sfxclient.Counter(name, nil, metric.Count()),
		}
	case metrics.Gauge:
		return []*datapoint.Datapoint{
			sfxclient.Gauge(name, nil, metric.Value()),
		}
	case metrics.GaugeFloat64:
		return []*datapoint.Datapoint{
			sfxclient.GaugeF(name, nil, metric.Value()),
		}
	case metrics.Histogram:
		h := metric.Snapshot()
		ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
		return []*datapoint.Datapoint{
			sfxclient.Counter(name+".count", nil, h.Count()),
			sfxclient.Counter(name+".min", nil, h.Min()),
			sfxclient.Counter(name+".max", nil, h.Max()),
			sfxclient.GaugeF(name+".mean", nil, h.Mean()),
			sfxclient.GaugeF(name+".std-dev", nil, h.StdDev()),
			sfxclient.GaugeF(name+".50-percentile", nil, ps[0]),
			sfxclient.GaugeF(name+".75-percentile", nil, ps[1]),
			sfxclient.GaugeF(name+".95-percentile", nil, ps[2]),
			sfxclient.GaugeF(name+".99-percentile", nil, ps[3]),
			sfxclient.GaugeF(name+".999-percentile", nil, ps[4]),
		}
	case metrics.Meter:
		m := metric.Snapshot()
		return []*datapoint.Datapoint{
			sfxclient.Counter(name+".count", nil, m.Count()),
			sfxclient.GaugeF(name+".one-minute", nil, m.Rate1()),
			sfxclient.GaugeF(name+".five-minute", nil, m.Rate5()),
			sfxclient.GaugeF(name+".fifteen-minute", nil, m.Rate15()),
			sfxclient.GaugeF(name+".mean-rate", nil, m.RateMean()),
		}
	case metrics.Timer:
		t := metric.Snapshot()
		ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
		return []*datapoint.Datapoint{
			sfxclient.Counter(name+".count", nil, t.Count()),
			sfxclient.Counter(name+".min", nil, t.Min()),
			sfxclient.Counter(name+".max", nil, t.Max()),
			sfxclient.GaugeF(name+".mean", nil, t.Mean()),
			sfxclient.GaugeF(name+".std-dev", nil, t.StdDev()),
			sfxclient.GaugeF(name+".50-percentile", nil, ps[0]),
			sfxclient.GaugeF(name+".75-percentile", nil, ps[1]),
			sfxclient.GaugeF(name+".95-percentile", nil, ps[2]),
			sfxclient.GaugeF(name+".99-percentile", nil, ps[3]),
			sfxclient.GaugeF(name+".999-percentile", nil, ps[4]),
			sfxclient.GaugeF(name+".one-minute", nil, t.Rate1()),
			sfxclient.GaugeF(name+".five-minute", nil, t.Rate5()),
			sfxclient.GaugeF(name+".fifteen-minute", nil, t.Rate15()),
			sfxclient.GaugeF(name+".mean-rate", nil, t.RateMean()),
		}
	default:
		panic(fmt.Sprintf("Unrecognized metric: %t.", i))
	}
}
