package cache

import "testing"

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
