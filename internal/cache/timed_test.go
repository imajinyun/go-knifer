package cache

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestTimedCache(t *testing.T) {
	c := NewTimed[string, string](4 * time.Millisecond)
	c.PutWithTimeout("key1", "value1", 1*time.Millisecond)
	c.PutWithTimeout("key2", "value2", 5*time.Second)
	c.Put("key3", "value3")               // Uses the default 4ms timeout.
	c.PutWithTimeout("key4", "value4", 0) // Never expires.

	c.SchedulePrune(5 * time.Millisecond)
	defer c.CancelPruneSchedule()
	time.Sleep(20 * time.Millisecond)

	if _, ok := c.Get("key1"); ok {
		t.Fatalf("key1 should expire")
	}
	if v, ok := c.Get("key2"); !ok || v != "value2" {
		t.Fatalf("key2: %v %v", v, ok)
	}
	if _, ok := c.Get("key3"); ok {
		t.Fatalf("key3 should expire")
	}
	if v, ok := c.Get("key4"); !ok || v != "value4" {
		t.Fatalf("key4: %v %v", v, ok)
	}

	v, err := c.GetOrLoad("key3", func() (string, error) { return "Default supplier", nil })
	if err != nil || v != "Default supplier" {
		t.Fatalf("GetOrLoad: %v %v", v, err)
	}
}

// Mirrors whenContainsKeyTimeout_shouldCallOnRemove.
func TestContainsKeyExpiredOnRemove(t *testing.T) {
	timeout := 50 * time.Millisecond
	c := NewTimed[int, string](timeout)
	var counter int32
	c.SetListener(CacheListenerFunc[int, string](func(int, string) {
		atomic.AddInt32(&counter, 1)
	}))
	c.Put(1, "value1")
	time.Sleep(100 * time.Millisecond)
	if c.ContainsKey(1) {
		t.Fatalf("should not contain key 1")
	}
	if got := atomic.LoadInt32(&counter); got != 1 {
		t.Fatalf("listener counter: %d", got)
	}
}

type testTicker struct {
	stopped chan struct{}
}

func (t *testTicker) Stop() { close(t.stopped) }

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
