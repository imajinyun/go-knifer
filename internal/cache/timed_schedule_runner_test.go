package cache

import (
	"testing"
	"time"
)

func TestTimedCacheWithRunner(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	now := base
	ticks := make(chan time.Time)
	ticker := &testTicker{stopped: make(chan struct{})}
	runnerCalls := make(chan struct{}, 1)
	c := NewTimedWithOptions[string, int](
		WithTimeout[string, int](time.Second),
		WithClock[string, int](func() time.Time { return now }),
		WithTickerFactory[string, int](func(time.Duration) (<-chan time.Time, Ticker) {
			return ticks, ticker
		}),
		WithRunner[string, int](func(fn func()) {
			runnerCalls <- struct{}{}
			go fn()
		}),
	)
	c.Put("a", 1)
	c.SchedulePrune(time.Second)
	select {
	case <-runnerCalls:
	case <-time.After(time.Second):
		t.Fatal("scheduled prune runner was not used")
	}
	now = base.Add(2 * time.Second)
	ticks <- now
	deadline := time.After(time.Second)
	for c.Size() != 0 {
		select {
		case <-deadline:
			t.Fatal("scheduled prune did not run")
		default:
			time.Sleep(time.Millisecond)
		}
	}
	c.CancelPruneSchedule()
}
