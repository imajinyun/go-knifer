package vmap

import (
	"sort"
	"testing"
)

func TestMapFacade(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	if IsEmpty(m) || !IsNotEmpty(m) {
		t.Fatal("empty checks failed")
	}
	keys := Keys(m)
	sort.Strings(keys)
	if len(keys) != 2 || keys[0] != "a" || keys[1] != "b" {
		t.Fatalf("Keys failed: %v", keys)
	}
	values := Values(m)
	sort.Ints(values)
	if len(values) != 2 || values[0] != 1 || values[1] != 2 {
		t.Fatalf("Values failed: %v", values)
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
