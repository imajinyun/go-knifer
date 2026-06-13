package cache

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestLRUCache(t *testing.T) {
	c := NewLRU[string, string](3)
	c.PutWithTimeout("key1", "value1", 3*time.Second)
	c.PutWithTimeout("key2", "value2", 3*time.Second)
	c.PutWithTimeout("key3", "value3", 3*time.Second)
	c.Get("key1") // Move key1 to the tail; key2 becomes least recently used.
	c.PutWithTimeout("key4", "value4", 3*time.Second)

	if _, ok := c.Get("key1"); !ok {
		t.Fatalf("key1 should still exist")
	}
	if _, ok := c.Get("key2"); ok {
		t.Fatalf("key2 should be evicted (LRU)")
	}
}

func TestLRURemoveCount(t *testing.T) {
	var count int32
	c := NewLRUWithTimeout[string, int](3, 1*time.Millisecond)
	c.SetListener(CacheListenerFunc[string, int](func(string, int) {
		atomic.AddInt32(&count, 1)
	}))
	for i := 0; i < 10; i++ {
		c.Put("key-"+itoa(i), i)
		// Sleep between puts so the previous value expires and prune triggers onRemove.
		time.Sleep(2 * time.Millisecond)
	}
	if c.Size() != 1 {
		// With ttl=1ms and a 2ms sleep, each put-triggered prune removes old
		// expired entries, leaving only the last inserted entry.
		t.Fatalf("expected size=1, got %d", c.Size())
	}
}

func TestLRUReadWriteConcurrency(t *testing.T) {
	const N = 10
	c := NewLRU[int, int](N)
	for i := 0; i < N; i++ {
		c.Put(i, i)
	}
	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				c.Get(idx)
			}
		}(i)
	}
	wg.Wait()
	// The order should still be 0..9. Each get moves a node to the tail, so 0 is
	// the first element after iterating over all keys in ascending order.
	got := ""
	for i := 0; i < N; i++ {
		if v, ok := c.Get(i); ok {
			got += itoa(v)
		} else {
			got += "x"
		}
	}
	if got != "0123456789" {
		t.Fatalf("got: %s", got)
	}
	// Adding 11 should evict 0, which is now the least recently used entry.
	c.Put(11, 11)
	if _, ok := c.Get(0); ok {
		t.Fatalf("key 0 should be evicted")
	}
}
