package cache

import "testing"

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
