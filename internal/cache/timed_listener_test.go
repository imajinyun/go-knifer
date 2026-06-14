package cache

import (
	"sync/atomic"
	"testing"
	"time"
)

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
