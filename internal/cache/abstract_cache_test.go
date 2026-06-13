package cache

import (
	"sync/atomic"
	"testing"
	"time"
)

// Mirrors reentrantCache_clear_Method_Test.
func TestLRUClearTriggersListener(t *testing.T) {
	var removeCount int32
	c := NewLRU[string, string](4)
	c.SetListener(CacheListenerFunc[string, string](func(string, string) {
		atomic.AddInt32(&removeCount, 1)
	}))
	c.Put("key1", "String1")
	c.Put("key2", "String2")
	c.Put("key3", "String3")
	c.Put("key1", "String4") // Replacement triggers one removal notification.
	c.Put("key4", "String5")
	c.Clear() // Clearing triggers the remaining 4 removal notifications.
	if got := atomic.LoadInt32(&removeCount); got != 5 {
		t.Fatalf("removeCount expected 5, got %d", got)
	}
}

func TestListenerCanReenterSameCache(t *testing.T) {
	c := NewLRU[string, string](2)
	done := make(chan struct{})
	c.SetListener(CacheListenerFunc[string, string](func(string, string) {
		defer close(done)
		_ = c.Size()
		c.SetListener(nil)
		c.Put("listener", "ok")
	}))
	c.Put("a", "1")
	c.Put("b", "2")
	c.Put("c", "3")
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("listener reentry deadlocked")
	}
	if v, ok := c.Get("listener"); !ok || v != "ok" {
		t.Fatalf("listener reentry value = %q ok=%v", v, ok)
	}
}

func TestGetOrLoad(t *testing.T) {
	c := NewLRU[string, int](3)
	v, err := c.GetOrLoad("a", func() (int, error) { return 42, nil })
	if err != nil || v != 42 {
		t.Fatalf("first: %v %v", v, err)
	}
	// The second call hits the cache directly and does not call supplier again.
	called := 0
	v, err = c.GetOrLoad("a", func() (int, error) {
		called++
		return 99, nil
	})
	if err != nil || v != 42 || called != 0 {
		t.Fatalf("second: v=%d err=%v called=%d", v, err, called)
	}
}

func TestRemoveAndContains(t *testing.T) {
	c := NewLRU[string, int](5)
	c.Put("a", 1)
	c.Put("b", 2)
	if !c.ContainsKey("a") {
		t.Fatal("expected contains a")
	}
	c.Remove("a")
	if c.ContainsKey("a") {
		t.Fatal("a should be removed")
	}
	if c.Size() != 1 {
		t.Fatalf("size: %d", c.Size())
	}
}

func TestHitMissCount(t *testing.T) {
	c := NewLRU[string, int](5)
	c.Put("a", 1)
	c.Get("a")
	c.Get("a")
	c.Get("b")
	if c.HitCount() != 2 || c.MissCount() != 1 {
		t.Fatalf("hit=%d miss=%d", c.HitCount(), c.MissCount())
	}
}
