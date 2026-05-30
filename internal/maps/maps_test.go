package maps

import (
	"sort"
	"testing"
)

func TestMapUtil(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	if MapIsEmpty(m) || !MapIsNotEmpty(m) {
		t.Fatalf("Map isEmpty failed")
	}
	keys := MapKeys(m)
	sort.Strings(keys)
	if len(keys) != 2 || keys[0] != "a" || keys[1] != "b" {
		t.Fatalf("MapKeys failed: %v", keys)
	}
	vs := MapValues(m)
	sort.Ints(vs)
	if vs[0] != 1 || vs[1] != 2 {
		t.Fatalf("MapValues failed: %v", vs)
	}
	inv := MapInverse(m)
	if inv[1] != "a" || inv[2] != "b" {
		t.Fatalf("MapInverse failed: %v", inv)
	}
	merged := MapMerge(map[string]int{"a": 1}, map[string]int{"a": 9, "b": 2})
	if merged["a"] != 9 || merged["b"] != 2 {
		t.Fatalf("MapMerge failed: %v", merged)
	}
}
