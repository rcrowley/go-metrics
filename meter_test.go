package metrics

import (
	"testing"
	"time"
)

func BenchmarkMeter(b *testing.B) {
	m := NewMeter()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Mark(1)
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

func TestMeterNonzero(t *testing.T) {
	m := NewMeter()
	m.Mark(3)
	if count := m.Count(); 3 != count {
		t.Errorf("m.Count(): 3 != %v\n", count)
	}
}

func TestMeterZero(t *testing.T) {
	m := NewMeter()
	if count := m.Count(); 0 != count {
		t.Errorf("m.Count(): 0 != %v\n", count)
	}
}
