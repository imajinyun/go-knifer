package cache

import (
	"testing"
	"time"
)

func TestTimedCacheWithTickerFactory(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	now := base
	ticks := make(chan time.Time)
	ticker := &testTicker{stopped: make(chan struct{})}
	c := NewTimedWithOptions[string, int](
		WithTimeout[string, int](time.Second),
		WithClock[string, int](func() time.Time { return now }),
		WithTickerFactory[string, int](func(delay time.Duration) (<-chan time.Time, Ticker) {
			if delay != 10*time.Second {
				t.Fatalf("ticker delay = %s, want 10s", delay)
			}
			return ticks, ticker
		}),
	)
	c.Put("a", 1)
	c.SchedulePrune(10 * time.Second)
	now = base.Add(2 * time.Second)
	ticks <- now

	deadline := time.After(time.Second)
	for c.Size() != 0 {
		select {
		case <-deadline:
			t.Fatal("scheduled prune did not run from custom ticker")
		default:
			time.Sleep(time.Millisecond)
		}
	}
	c.CancelPruneSchedule()
	select {
	case <-ticker.stopped:
	case <-time.After(time.Second):
		t.Fatal("custom ticker was not stopped")
	}
}
