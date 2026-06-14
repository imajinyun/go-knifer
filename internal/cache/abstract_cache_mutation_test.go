package cache

import "testing"

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
