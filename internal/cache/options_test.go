package cache

import (
	"testing"
	"time"
)

func TestCacheOptions(t *testing.T) {
	removed := make([]string, 0)
	fifo := NewFIFOWithOptions[string, string](
		WithCapacity[string, string](1),
		WithTimeout[string, string](time.Second),
		WithListener[string, string](CacheListenerFunc[string, string](func(key, value string) {
			removed = append(removed, key+"="+value)
		})),
	)
	if fifo.Capacity() != 1 || fifo.Timeout() != time.Second {
		t.Fatalf("fifo options not applied: capacity=%d timeout=%s", fifo.Capacity(), fifo.Timeout())
	}
	fifo.Put("a", "1")
	fifo.Put("b", "2")
	if len(removed) != 1 || removed[0] != "a=1" {
		t.Fatalf("fifo listener removed = %v", removed)
	}

	lfu := NewLFUWithOptions[string, int](WithCapacity[string, int](2), WithTimeout[string, int](time.Second))
	if lfu.Capacity() != 2 || lfu.Timeout() != time.Second {
		t.Fatalf("lfu options not applied: capacity=%d timeout=%s", lfu.Capacity(), lfu.Timeout())
	}
	lru := NewLRUWithOptions[string, int](WithCapacity[string, int](2), WithTimeout[string, int](time.Second))
	if lru.Capacity() != 2 || lru.Timeout() != time.Second {
		t.Fatalf("lru options not applied: capacity=%d timeout=%s", lru.Capacity(), lru.Timeout())
	}
	timed := NewTimedWithOptions[string, int](WithTimeout[string, int](time.Second))
	if timed.Capacity() != 0 || timed.Timeout() != time.Second {
		t.Fatalf("timed options not applied: capacity=%d timeout=%s", timed.Capacity(), timed.Timeout())
	}
}

func TestCacheWithClock(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	now := base
	c := NewTimedWithOptions[string, int](
		WithTimeout[string, int](time.Second),
		WithClock[string, int](func() time.Time { return now }),
	)
	c.Put("a", 1)
	now = base.Add(500 * time.Millisecond)
	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Fatalf("expected value before custom-clock expiry, got %d ok=%v", v, ok)
	}
	now = base.Add(2 * time.Second)
	if _, ok := c.GetWithUpdate("a", false); ok {
		t.Fatalf("expected custom clock to expire entry")
	}
}
