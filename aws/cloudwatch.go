// Metrics output to StatHat.
package cloudwatch

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/rcrowley/go-metrics"
)

//EmitMetrics emits the metrics in a metrics registry to cloudwatch
//Param type *metrics.Registry (metrics registry from which to extract metrics)
//Param type time.Duration (how often to emit metrics)
//Param type string (Cloudwatch namespace under which metrics will be stored)
func EmitMetrics(r *metrics.Registry, d time.Duration, namespace string) {
	svc := cloudwatch.New(session.New())
	for {
		if errs := emit(svc, *r, namespace); 0 != len(errs) {
			log.Println(errs)
		}
		time.Sleep(d)
	}
}

func emit(svc *cloudwatch.CloudWatch, r metrics.Registry, s string) []string {
	awsErr := []string{}
	metricData := []*cloudwatch.MetricDatum{}
	now := aws.Time(time.Now())
	params := &cloudwatch.PutMetricDataInput{}

	r.Each(func(name string, i interface{}) {
		metricData = nil
		params = nil

		switch metric := i.(type) {
		case metrics.Counter:
			metricData = append(metricData, &cloudwatch.MetricDatum{
				MetricName: aws.String(name),
				Timestamp:  now,
				Unit:       aws.String("Count"),
				Value:      aws.Float64(float64(metric.Count())),
			})
		case metrics.Gauge:
			metricData = append(metricData, &cloudwatch.MetricDatum{
				MetricName: aws.String(name),
				Timestamp:  now,
				Unit:       aws.String("None"),
				Value:      aws.Float64(float64(metric.Value())),
			})
		case metrics.GaugeFloat64:
			metricData = append(metricData,
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name),
					Timestamp:  now,
					Unit:       aws.String("None"),
					Value:      aws.Float64(metric.Value()),
				})
		case metrics.Histogram:
			m := metric.Snapshot()
			ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			metricData = append(metricData,
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".count"),
					Timestamp:  now,
					Unit:       aws.String("Count"),
					Value:      aws.Float64(float64(m.Count())),
				},
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".min"),
					Timestamp:  now,
					Unit:       aws.String("Microseconds"),
					Value:      aws.Float64(float64(m.Min())),
				},
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".max"),
					Timestamp:  now,
					Unit:       aws.String("Microseconds"),
					Value:      aws.Float64(float64(m.Max())),
				},
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".mean"),
					Timestamp:  now,
					Unit:       aws.String("Microseconds"),
					Value:      aws.Float64(float64(m.Mean())),
				},
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".std-dev"),
					Timestamp:  now,
					Unit:       aws.String("None"),
					Value:      aws.Float64(float64(m.StdDev())),
				},
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".50-percentile"),
					Timestamp:  now,
					Unit:       aws.String("Microseconds"),
					Value:      aws.Float64(float64(ps[0])),
				},
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".75-percentile"),
					Timestamp:  now,
					Unit:       aws.String("Microseconds"),
					Value:      aws.Float64(float64(ps[1])),
				},
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".95-percentile"),
					Timestamp:  now,
					Unit:       aws.String("Microseconds"),
					Value:      aws.Float64(float64(ps[2])),
				},
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".99-percentile"),
					Timestamp:  now,
					Unit:       aws.String("Microseconds"),
					Value:      aws.Float64(float64(ps[3])),
				},
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".999-percentile"),
					Timestamp:  now,
					Unit:       aws.String("Microseconds"),
					Value:      aws.Float64(float64(ps[4])),
				},
			)
		case metrics.Meter:
			m := metric.Snapshot()
			metricData = append(metricData,
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".count"),
					Timestamp:  now,
					Unit:       aws.String("Count"),
					Value:      aws.Float64(float64(m.Count())),
				},
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".one-minute"),
					Timestamp:  now,
					Unit:       aws.String("Count/Second"),
					Value:      aws.Float64(float64(m.Rate1())),
				},
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".five-minute"),
					Timestamp:  now,
					Unit:       aws.String("Count/Second"),
					Value:      aws.Float64(float64(m.Rate5())),
				},
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".fifteen-minute"),
					Timestamp:  now,
					Unit:       aws.String("Count/Second"),
					Value:      aws.Float64(float64(m.Rate15())),
				},
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".mean"),
					Timestamp:  now,
					Unit:       aws.String("Count/Second"),
					Value:      aws.Float64(float64(m.RateMean())),
				})
		case metrics.Timer:
			m := metric.Snapshot()
			ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			metricData = append(metricData,
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".count"),
					Timestamp:  now,
					Unit:       aws.String("Count"),
					Value:      aws.Float64(float64(m.Count())),
				},
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".min"),
					Timestamp:  now,
					Unit:       aws.String("Microseconds"),
					Value:      aws.Float64(float64(m.Min())),
				},
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".max"),
					Timestamp:  now,
					Unit:       aws.String("Microseconds"),
					Value:      aws.Float64(float64(m.Max())),
				},
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".mean"),
					Timestamp:  now,
					Unit:       aws.String("Microseconds"),
					Value:      aws.Float64(float64(m.Mean())),
				},
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".std-dev"),
					Timestamp:  now,
					Unit:       aws.String("None"),
					Value:      aws.Float64(float64(m.StdDev())),
				},
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".one-minute"),
					Timestamp:  now,
					Unit:       aws.String("Count/Second"),
					Value:      aws.Float64(float64(m.Rate1())),
				},
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".five-minute"),
					Timestamp:  now,
					Unit:       aws.String("Count/Second"),
					Value:      aws.Float64(float64(m.Rate5())),
				},
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".fifteen-minute"),
					Timestamp:  now,
					Unit:       aws.String("Count/Second"),
					Value:      aws.Float64(float64(m.Rate15())),
				},
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".mean-rate"),
					Timestamp:  now,
					Unit:       aws.String("Count/Second"),
					Value:      aws.Float64(float64(m.RateMean())),
				},
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".50-percentile"),
					Timestamp:  now,
					Unit:       aws.String("Microseconds"),
					Value:      aws.Float64(float64(ps[0])),
				},
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".75-percentile"),
					Timestamp:  now,
					Unit:       aws.String("Microseconds"),
					Value:      aws.Float64(float64(ps[1])),
				},
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".95-percentile"),
					Timestamp:  now,
					Unit:       aws.String("Microseconds"),
					Value:      aws.Float64(float64(ps[2])),
				},
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".99-percentile"),
					Timestamp:  now,
					Unit:       aws.String("Microseconds"),
					Value:      aws.Float64(float64(ps[3])),
				},
				&cloudwatch.MetricDatum{
					MetricName: aws.String(name + ".999-percentile"),
					Timestamp:  now,
					Unit:       aws.String("Microseconds"),
					Value:      aws.Float64(float64(ps[4])),
				},
			)
		}

		if len(metricData) > 0 {
			params = &cloudwatch.PutMetricDataInput{
				MetricData: metricData,
				Namespace:  aws.String(s),
			}

			_, err := svc.PutMetricData(params)

			if err != nil {
				awsErr = append(awsErr, err.Error())
			}
		}

	})

	return awsErr

}
