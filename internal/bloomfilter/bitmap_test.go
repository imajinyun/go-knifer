package bloomfilter

import "testing"

func TestIntMapAndLongMap(t *testing.T) {
	im := NewIntMap(8)
	im.Add(0)
	im.Add(31)
	im.Add(64)
	if !im.Contains(0) || !im.Contains(31) || !im.Contains(64) {
		t.Fatal("intmap contains failed")
	}
	if im.Contains(1) {
		t.Fatal("intmap should not contain 1")
	}
	im.Remove(31)
	if im.Contains(31) {
		t.Fatal("intmap remove failed")
	}

	lm := NewLongMap(4)
	lm.Add(0)
	lm.Add(63)
	lm.Add(128)
	if !lm.Contains(0) || !lm.Contains(63) || !lm.Contains(128) {
		t.Fatal("longmap contains failed")
	}
	lm.Remove(128)
	if lm.Contains(128) {
		t.Fatal("longmap remove failed")
	}
}
