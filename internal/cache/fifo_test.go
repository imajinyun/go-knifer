package cache

import (
	"testing"
	"time"
)

func TestFIFOCache(t *testing.T) {
	var removedKey, removedValue string
	c := NewFIFO[string, string](3)
	c.SetListener(CacheListenerFunc[string, string](func(k, v string) {
		removedKey, removedValue = k, v
	}))
	c.PutWithTimeout("key1", "value1", 3*time.Second)
	c.PutWithTimeout("key2", "value2", 3*time.Second)
	c.PutWithTimeout("key3", "value3", 3*time.Second)
	c.PutWithTimeout("key4", "value4", 3*time.Second)

	// Adding the 4th entry should evict the oldest key, key1.
	if v, ok := c.Get("key1"); ok {
		t.Fatalf("key1 should be evicted, got %q", v)
	}
	if removedKey != "key1" || removedValue != "value1" {
		t.Fatalf("listener got: %s=%s", removedKey, removedValue)
	}
}

func TestFIFOCapacity(t *testing.T) {
	c := NewFIFO[string, string](100)
	for i := 0; i < 500; i++ {
		c.Put(itoa(i), "v")
	}
	if got := c.Size(); got != 100 {
		t.Fatalf("size: %d", got)
	}
}
