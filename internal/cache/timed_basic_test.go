package cache

import (
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
