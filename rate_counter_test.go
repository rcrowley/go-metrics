package metrics

import (
	"testing"
	"time"

	"github.com/bountylabs/go-metrics/clock"
)

func TestRateCounter(t *testing.T) {
	m := clock.NewMock()
	rc := NewStandardRateCounter(60, 1000, m)

	rc.Mark(1)
	m.Add(1 * time.Second)
	if v := rc.Rate1(); v != 1.0 {
		t.Errorf("rc.Rate1(): 1.0 != %v\n", v)
	}
}

func TestShouldNotTakeIntoAccountDataFromOverAMinuteAgo(t *testing.T) {
	m := clock.NewMock()
	rc := NewStandardRateCounter(60, 1000, m)

	// Mark only takes up to one sample per second so must increment by a second so a sample is taken
	m.Add(1 * time.Second)

	rc.Mark(500)
	m.Add(60 * time.Second)

	rc.Mark(60)

	if v := rc.Rate1(); v != 1.0 {
		t.Errorf("rc.Rate1(): 1.0 != %v\n", v)
	}
}

func TestShouldClearRateCounter(t *testing.T) {
	m := clock.NewMock()
	rc := NewStandardRateCounter(60, 1000, m)

	rc.Mark(1)
	m.Add(1 * time.Second)

	rc.Clear()

	if v := rc.Rate1(); v != 0.0 {
		t.Errorf("rc.Rate1(): 0.0 != %v\n", v)
	}
}

func TestSnapshot(t *testing.T) {
	m := clock.NewMock()
	rc := NewStandardRateCounter(60, 1000, m)

	rc.Mark(1)
	m.Add(1 * time.Second)

	s := rc.Snapshot()
	rc.Mark(100)

	if v := s.Rate1(); v != 1.0 {
		t.Errorf("s.Rate1(): 1.0 != %v\n", v)
	}
}
