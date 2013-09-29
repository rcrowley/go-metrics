package metrics

import (
	"testing"
	"time"
)

func TestMeterZero(t *testing.T) {
	m := NewMeter()
	if count := m.Count(); 0 != count {
		t.Errorf("m.Count(): 0 != %v\n", count)
	}
}

func TestMeterNonzero(t *testing.T) {
	m := NewMeter()
	m.Mark(3)
	if count := m.Count(); 3 != count {
		t.Errorf("m.Count(): 3 != %v\n", count)
	}
}

func TestMeterDecay(t *testing.T) {
	m := &StandardMeter{
		make(chan int64),
		make(chan meterV),
		time.NewTicker(1),
	}
	go m.arbiter()
	m.Mark(1)
	rateMean := m.RateMean()
	time.Sleep(1)
	if m.RateMean() >= rateMean {
		t.Error("m.RateMean() didn't decrease")
	}
}
