package logging

import (
	"github.com/zeim839/go-metrics++"
	"strings"
	"testing"
	"time"
)

func BenchmarkEncode(b *testing.B) {
	// Timer is worst-case scenario (most verbose).
	timer := metrics.NewTimer(metrics.Label{"key1", "value1"},
		metrics.Label{"key2", "value2"})
	timer.Update(time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Encode("foo", "bar", timer)
	}
}

func TestEncodeCounter(t *testing.T) {
	counter := metrics.NewCounter(metrics.Label{"key1", "value1"})
	counter.Inc(500)
	expect := "bar_foo{key1:\"value1\"} 500"
	if str := Encode("foo", "bar", counter); str[:len(str)-12] != expect {
		t.Errorf("Encode(): %s != %s", str[:len(str)-12], expect)
	}

	// Without namespace.
	expect = "foo{key1:\"value1\"} 500"
	if str := Encode("foo", "", counter); str[:len(str)-12] != expect {
		t.Errorf("Encode(): %s != %s", str[:len(str)-12], expect)
	}

	// Without labels.
	counter = metrics.NewCounter()
	expect = "bar_foo 0"
	if str := Encode("foo", "bar", counter); str[:len(str)-12] != expect {
		t.Errorf("Encode(): %s != %s", str[:len(str)-12], expect)
	}
}

func TestEncodeGauge(t *testing.T) {
	gauge := metrics.NewGauge(metrics.Label{"foo", "bar"})
	gauge.Update(10)
	expect := "bar_foo{foo:\"bar\"} 10"
	if str := Encode("foo", "bar", gauge); str[:len(str)-12] != expect {
		t.Errorf("Encode(): %s != %s", str[:len(str)-12], expect)
	}

	// Without namespace.
	expect = "foo{foo:\"bar\"} 10"
	if str := Encode("foo", "", gauge); str[:len(str)-12] != expect {
		t.Errorf("Encode(): %s != %s", str[:len(str)-12], expect)
	}

	// Without labels.
	gauge = metrics.NewGauge()
	expect = "bar_foo 0"
	if str := Encode("foo", "bar", gauge); str[:len(str)-12] != expect {
		t.Errorf("Encode(): %s != %s", str[:len(str)-12], expect)
	}
}

func TestEncodeGaugeFloat64(t *testing.T) {
	gauge := metrics.NewGaugeFloat64(metrics.Label{"foo", "bar"})
	gauge.Update(10)
	expect := "bar_foo{foo:\"bar\"} 10.000000"
	if str := Encode("foo", "bar", gauge); str[:len(str)-12] != expect {
		t.Errorf("Encode(): %s != %s", str[:len(str)-12], expect)
	}

	// Without namespace.
	expect = "foo{foo:\"bar\"} 10.000000"
	if str := Encode("foo", "", gauge); str[:len(str)-12] != expect {
		t.Errorf("Encode(): %s != %s", str[:len(str)-12], expect)
	}

	// Without labels.
	gauge = metrics.NewGaugeFloat64()
	expect = "bar_foo 0.000000"
	if str := Encode("foo", "bar", gauge); str[:len(str)-12] != expect {
		t.Errorf("Encode(): %s != %s", str[:len(str)-12], expect)
	}
}

func TestEncodeHealthcheck(t *testing.T) {
	check := metrics.NewHealthcheck(func(metrics.Healthcheck) {})
	if str := Encode("foo", "bar", check); "" != str {
		t.Errorf("Encode(): Healthcheck returned non-empty string: %s", str)
	}
}

func TestEncodeHistogram(t *testing.T) {
	hist := metrics.NewHistogram(metrics.NewUniformSample(100),
		metrics.Label{"foo", "bar"})
	hist.Update(100.0)
	lines := strings.Split(Encode("foo", "bar", hist), "\n")
	if len(lines) != 11 {
		t.Fatal("Encode(): Did not produce 11 lines for timer")
	}
	expect := "bar_foo_count{foo:\"bar\"} 1"
	if line := lines[0][:len(lines[0])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_min{foo:\"bar\"} 100"
	if line := lines[1][:len(lines[1])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_max{foo:\"bar\"} 100"
	if line := lines[2][:len(lines[2])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_mean{foo:\"bar\"} 100.000000"
	if line := lines[3][:len(lines[3])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_stddev{foo:\"bar\"} 0.000000"
	if line := lines[4][:len(lines[4])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_median{foo:\"bar\"} 100.000000"
	if line := lines[5][:len(lines[5])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_percentile_75{foo:\"bar\"} 100.000000"
	if line := lines[6][:len(lines[6])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_percentile_95{foo:\"bar\"} 100.000000"
	if line := lines[7][:len(lines[7])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_percentile_99{foo:\"bar\"} 100.000000"
	if line := lines[8][:len(lines[8])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_percentile_99_9{foo:\"bar\"} 100.000000"
	if line := lines[9][:len(lines[9])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}

	// Without namespace
	lines = strings.Split(Encode("foo", "", hist), "\n")
	if len(lines) != 11 {
		t.Fatal("Encode(): Did not produce 11 lines for timer")
	}
	expect = "foo_count{foo:\"bar\"} 1"
	if line := lines[0][:len(lines[0])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_min{foo:\"bar\"} 100"
	if line := lines[1][:len(lines[1])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_max{foo:\"bar\"} 100"
	if line := lines[2][:len(lines[2])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_mean{foo:\"bar\"} 100.000000"
	if line := lines[3][:len(lines[3])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_stddev{foo:\"bar\"} 0.000000"
	if line := lines[4][:len(lines[4])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_median{foo:\"bar\"} 100.000000"
	if line := lines[5][:len(lines[5])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_percentile_75{foo:\"bar\"} 100.000000"
	if line := lines[6][:len(lines[6])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_percentile_95{foo:\"bar\"} 100.000000"
	if line := lines[7][:len(lines[7])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_percentile_99{foo:\"bar\"} 100.000000"
	if line := lines[8][:len(lines[8])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_percentile_99_9{foo:\"bar\"} 100.000000"
	if line := lines[9][:len(lines[9])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}

	// Without labels.
	hist = metrics.NewHistogram(metrics.NewUniformSample(100))
	hist.Update(100)
	lines = strings.Split(Encode("foo", "bar", hist), "\n")
	if len(lines) != 11 {
		t.Fatal("Encode(): Did not produce 11 lines for timer")
	}
	expect = "bar_foo_count 1"
	if line := lines[0][:len(lines[0])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_min 100"
	if line := lines[1][:len(lines[1])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_max 100"
	if line := lines[2][:len(lines[2])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_mean 100.000000"
	if line := lines[3][:len(lines[3])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_stddev 0.000000"
	if line := lines[4][:len(lines[4])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_median 100.000000"
	if line := lines[5][:len(lines[5])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_percentile_75 100.000000"
	if line := lines[6][:len(lines[6])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_percentile_95 100.000000"
	if line := lines[7][:len(lines[7])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_percentile_99 100.000000"
	if line := lines[8][:len(lines[8])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_percentile_99_9 100.000000"
	if line := lines[9][:len(lines[9])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
}

func TestEncodeMeter(t *testing.T) {
	meter := metrics.NewMeter(metrics.Label{"foo", "bar"})
	meter.Mark(20)
	lines := strings.Split(Encode("foo", "bar", meter), "\n")
	if len(lines) != 6 {
		t.Fatal("Encode(): Did not produce six lines for meter")
	}
	expect := "bar_foo_count{foo:\"bar\"} 20"
	if line := lines[0][:len(lines[0])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_rate_1min{foo:\"bar\"} 20.000000"
	if line := lines[1][:len(lines[1])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_rate_5min{foo:\"bar\"} 20.000000"
	if line := lines[2][:len(lines[2])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_rate_15min{foo:\"bar\"} 20.000000"
	if line := lines[3][:len(lines[3])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_rate_mean{foo:\"bar\"}"
	if line := lines[4][:len(expect)]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}

	// Without namespace.
	lines = strings.Split(Encode("foo", "", meter), "\n")
	if len(lines) != 6 {
		t.Fatal("Encode(): Did not produce six lines for meter")
	}
	expect = "foo_count{foo:\"bar\"} 20"
	if line := lines[0][:len(lines[0])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_rate_1min{foo:\"bar\"} 20.000000"
	if line := lines[1][:len(lines[1])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_rate_5min{foo:\"bar\"} 20.000000"
	if line := lines[2][:len(lines[2])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_rate_15min{foo:\"bar\"} 20.000000"
	if line := lines[3][:len(lines[3])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_rate_mean{foo:\"bar\"}"
	if line := lines[4][:len(expect)]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}

	// Without labels.
	meter = metrics.NewMeter()
	lines = strings.Split(Encode("foo", "bar", meter), "\n")
	if len(lines) != 6 {
		t.Fatal("Encode(): Did not produce six lines for meter")
	}
	expect = "bar_foo_count 0"
	if line := lines[0][:len(lines[0])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_rate_1min 0.000000"
	if line := lines[1][:len(lines[1])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_rate_5min 0.000000"
	if line := lines[2][:len(lines[2])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_rate_15min 0.000000"
	if line := lines[3][:len(lines[3])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_rate_mean"
	if line := lines[4][:len(expect)]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
}

func TestEncodeTimer(t *testing.T) {
	// Do not timer.Update() without some time.Sleep, results are erratic.
	timer := metrics.NewTimer(metrics.Label{"foo", "bar"})
	lines := strings.Split(Encode("foo", "bar", timer), "\n")
	if len(lines) != 15 {
		t.Fatal("Encode(): Did not produce 15 lines for timer")
	}
	expect := "bar_foo_count{foo:\"bar\"} 0"
	if line := lines[0][:len(lines[0])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_min{foo:\"bar\"} 0"
	if line := lines[1][:len(lines[1])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_max{foo:\"bar\"} 0"
	if line := lines[2][:len(lines[2])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_mean{foo:\"bar\"} 0.000000"
	if line := lines[3][:len(lines[3])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_stddev{foo:\"bar\"} 0.000000"
	if line := lines[4][:len(lines[4])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_median{foo:\"bar\"} 0.000000"
	if line := lines[5][:len(lines[5])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_percentile_75{foo:\"bar\"} 0.000000"
	if line := lines[6][:len(lines[6])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_percentile_95{foo:\"bar\"} 0.000000"
	if line := lines[7][:len(lines[7])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_percentile_99{foo:\"bar\"} 0.000000"
	if line := lines[8][:len(lines[8])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_percentile_99_9{foo:\"bar\"} 0.000000"
	if line := lines[9][:len(lines[9])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_rate_1min{foo:\"bar\"} 0.000000"
	if line := lines[10][:len(lines[10])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_rate_5min{foo:\"bar\"} 0.000000"
	if line := lines[11][:len(lines[11])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_rate_15min{foo:\"bar\"} 0.000000"
	if line := lines[12][:len(lines[12])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "bar_foo_rate_mean{foo:\"bar\"} 0.000000"
	if line := lines[13][:len(lines[13])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}

	// Without namespace.
	lines = strings.Split(Encode("foo", "", timer), "\n")
	if len(lines) != 15 {
		t.Error("Encode(): Did not produce 15 lines for timer")
	}
	expect = "foo_count{foo:\"bar\"} 0"
	if line := lines[0][:len(lines[0])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_min{foo:\"bar\"} 0"
	if line := lines[1][:len(lines[1])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_max{foo:\"bar\"} 0"
	if line := lines[2][:len(lines[2])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_mean{foo:\"bar\"} 0.000000"
	if line := lines[3][:len(lines[3])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_stddev{foo:\"bar\"} 0.000000"
	if line := lines[4][:len(lines[4])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_median{foo:\"bar\"} 0.000000"
	if line := lines[5][:len(lines[5])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_percentile_75{foo:\"bar\"} 0.000000"
	if line := lines[6][:len(lines[6])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_percentile_95{foo:\"bar\"} 0.000000"
	if line := lines[7][:len(lines[7])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_percentile_99{foo:\"bar\"} 0.000000"
	if line := lines[8][:len(lines[8])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_percentile_99_9{foo:\"bar\"} 0.000000"
	if line := lines[9][:len(lines[9])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_rate_1min{foo:\"bar\"} 0.000000"
	if line := lines[10][:len(lines[10])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_rate_5min{foo:\"bar\"} 0.000000"
	if line := lines[11][:len(lines[11])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_rate_15min{foo:\"bar\"} 0.000000"
	if line := lines[12][:len(lines[12])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_rate_mean{foo:\"bar\"} 0.000000"
	if line := lines[13][:len(lines[13])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}

	// Without labels.
	timer = metrics.NewTimer()
	lines = strings.Split(Encode("foo", "", timer), "\n")
	if len(lines) != 15 {
		t.Fatal("Encode(): Did not produce 15 lines for timer")
	}
	expect = "foo_count 0"
	if line := lines[0][:len(lines[0])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_min 0"
	if line := lines[1][:len(lines[1])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_max 0"
	if line := lines[2][:len(lines[2])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_mean 0.000000"
	if line := lines[3][:len(lines[3])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_stddev 0.000000"
	if line := lines[4][:len(lines[4])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_median 0.000000"
	if line := lines[5][:len(lines[5])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_percentile_75 0.000000"
	if line := lines[6][:len(lines[6])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_percentile_95 0.000000"
	if line := lines[7][:len(lines[7])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_percentile_99 0.000000"
	if line := lines[8][:len(lines[8])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_percentile_99_9 0.000000"
	if line := lines[9][:len(lines[9])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_rate_1min 0.000000"
	if line := lines[10][:len(lines[10])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_rate_5min 0.000000"
	if line := lines[11][:len(lines[11])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_rate_15min 0.000000"
	if line := lines[12][:len(lines[12])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
	expect = "foo_rate_mean 0.000000"
	if line := lines[13][:len(lines[13])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", expect, line)
	}
}

func TestEncodeUnknown(t *testing.T) {
	srt := struct {
		a string
		b int16
	}{"asd", 123}

	if str := Encode("foo", "bar", srt); str != "" {
		t.Errorf("Encode(): Unknown struct returned non-empty string: %s", str)
	}
}

func TestEncodeLabels(t *testing.T) {
	a := metrics.Label{"key1", "value1"}
	b := metrics.Label{"key2", "value2"}
	slice := []metrics.Label{}

	// Empty slice returns empty string.
	if str := EncodeLabels(slice); str != "" {
		t.Errorf("EncodeLabels(): Empty slice returned %s", str)
	}

	slice = []metrics.Label{a, b}
	expect := "{key1:\"value1\",key2:\"value2\"}"
	if str := EncodeLabels(slice); str != expect {
		t.Errorf("EncodeLabels(): %s != %s", str, expect)
	}
}

func BenchmarkEncodeLabels(b *testing.B) {
	a := metrics.Label{"key1", "value1"}
	slice := []metrics.Label{a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EncodeLabels(slice)
	}
}
