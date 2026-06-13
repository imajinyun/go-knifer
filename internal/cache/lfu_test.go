package cache

import (
	"testing"
	"time"
)

func TestLFUCache(t *testing.T) {
	c := NewLFU[string, string](3)
	c.PutWithTimeout("key1", "value1", 3*time.Second)
	c.Get("key1") // Increase the access count by 1.
	c.PutWithTimeout("key2", "value2", 3*time.Second)
	c.PutWithTimeout("key3", "value3", 3*time.Second)
	c.PutWithTimeout("key4", "value4", 3*time.Second)

	if _, ok := c.Get("key1"); !ok {
		t.Fatalf("key1 should still exist")
	}
	if _, ok := c.Get("key2"); ok {
		t.Fatalf("key2 should be evicted")
	}
	if _, ok := c.Get("key3"); ok {
		t.Fatalf("key3 should be evicted")
	}
}
