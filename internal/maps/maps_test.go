package maps

import (
	"sort"
	"testing"
)

func TestMapUtil(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	if IsEmpty(m) || !IsNotEmpty(m) {
		t.Fatalf("isEmpty failed")
	}
	keys := Keys(m)
	sort.Strings(keys)
	if len(keys) != 2 || keys[0] != "a" || keys[1] != "b" {
		t.Fatalf("Keys failed: %v", keys)
	}
	vs := Values(m)
	sort.Ints(vs)
	if vs[0] != 1 || vs[1] != 2 {
		t.Fatalf("Values failed: %v", vs)
	}
	inv := Inverse(m)
	if inv[1] != "a" || inv[2] != "b" {
		t.Fatalf("Inverse failed: %v", inv)
	}
	merged := Merge(map[string]int{"a": 1}, map[string]int{"a": 9, "b": 2})
	if merged["a"] != 9 || merged["b"] != 2 {
		t.Fatalf("Merge failed: %v", merged)
	}
}
