package cache

import "testing"

func TestNoCache(t *testing.T) {
	c := NewNo[string, int]()
	c.Put("a", 1)
	if c.Size() != 0 || !c.IsEmpty() {
		t.Fatalf("nocache size/empty wrong")
	}
	if _, ok := c.Get("a"); ok {
		t.Fatalf("nocache get hit?")
	}
	v, err := c.GetOrLoad("k", func() (int, error) { return 7, nil })
	if err != nil || v != 7 {
		t.Fatalf("nocache loader: %v %v", v, err)
	}
}
