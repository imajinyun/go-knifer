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
