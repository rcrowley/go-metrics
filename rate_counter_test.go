package metrics

import (
	"testing"
	"time"

	"github.com/bountylabs/go-metrics/clock"
)

func TestRateCounterZero(t *testing.T) {
	m := clock.NewMock()
	rc := NewStandardRateCounter(60, 1000, m)

	if v := rc.Rate1(); v != 0.0 {
		t.Errorf("rc.Rate1(): 1.0 != %v\n", v)
	}

	if v := rc.Count(); v != 0.0 {
		t.Errorf("rc.Count(): 1.0 != %v\n", v)
	}
}

func TestRateCounter(t *testing.T) {
	m := clock.NewMock()
	rc := NewStandardRateCounter(60, 1000, m)

	rc.Mark(1)
	m.Add(1 * time.Second)
	if v := rc.Rate1(); v != 1.0 {
		t.Errorf("rc.Rate1(): 1.0 != %v\n", v)
	}
	if v := rc.Count(); v != 1.0 {
		t.Errorf("rc.Count(): 1.0 != %v\n", v)
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

	if v := rc.Count(); v != 560.0 {
		t.Errorf("rc.Count(): 560.0 != %v\n", v)
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
	if v := rc.Count(); v != 0.0 {
		t.Errorf("rc.Count(): 0.0 != %v\n", v)
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
	if v := s.Count(); v != 1.0 {
		t.Errorf("s.Count(): 1.0 != %v\n", v)
	}
}

func TestRateAndLastCountInSnapshotShouldBeConsistent(t *testing.T) {
	// If new data isn't present in Snapshot.Rate1() it shouldn't be present in counter.lastCount but it should
	// in Snapshot.Count()

	m := clock.NewMock()
	rc := NewStandardRateCounter(60, 1000, m)

	rc.Mark(1)
	s := rc.Snapshot()

	// since we don't advance the time, the sampling period is still current, and there is nothing finalized for rate or .lastCount
	// but this doesn't matter for Count()
	if v := s.Rate1(); v != 0.0 {
		t.Errorf("s.Rate1(): 0.0 != %v\n", v)
	}
	if v := rc.(*StandardRateCounter).lastCount; v != 0.0 {
		t.Errorf("s.lastCount: 0.0 != %v\n", v)
	}
	if v := s.Count(); v != 1.0 {
		t.Errorf("s.Count(): 1.0 != %v\n", v)
	}

	m.Add(1 * time.Second)
	s = rc.Snapshot()

	if v := s.Rate1(); v != 1.0 {
		t.Errorf("s.Rate1(): 1.0 != %v\n", v)
	}
	if v := s.Count(); v != 1.0 {
		t.Errorf("s.Count(): 1.0 != %v\n", v)
	}
}
