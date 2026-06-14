package cache

import (
	"testing"
	"time"
)

func TestTimedCacheSchedulePruneTwiceIsNoop(t *testing.T) {
	runnerCalls := make(chan struct{}, 2)
	ticks := make(chan time.Time)
	ticker := &testTicker{stopped: make(chan struct{})}
	c := NewTimedWithOptions[string, int](
		WithTickerFactory[string, int](func(time.Duration) (<-chan time.Time, Ticker) {
			return ticks, ticker
		}),
		WithRunner[string, int](func(fn func()) {
			runnerCalls <- struct{}{}
			go fn()
		}),
	)
	c.SchedulePrune(time.Second)
	c.SchedulePrune(time.Second)
	select {
	case <-runnerCalls:
	case <-time.After(time.Second):
		t.Fatal("first scheduled prune did not start")
	}
	select {
	case <-runnerCalls:
		t.Fatal("second SchedulePrune should be a no-op while pruning is already running")
	case <-time.After(20 * time.Millisecond):
	}
	c.CancelPruneSchedule()
	c.CancelPruneSchedule()
	select {
	case <-ticker.stopped:
	case <-time.After(time.Second):
		t.Fatal("ticker was not stopped")
	}
}
