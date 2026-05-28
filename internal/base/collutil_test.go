package base

import (
	"sort"
	"testing"
)

// 对应 hutool-core CollUtilTest / MapUtilTest。

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

func TestUnionIntersectionSubtract(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{2, 3, 4}
	u := Union(a, b)
	sort.Ints(u)
	if len(u) != 4 {
		t.Fatalf("Union failed: %v", u)
	}
	in := Intersection(a, b)
	sort.Ints(in)
	if len(in) != 2 || in[0] != 2 || in[1] != 3 {
		t.Fatalf("Intersection failed: %v", in)
	}
	sub := Subtract(a, b)
	if len(sub) != 1 || sub[0] != 1 {
		t.Fatalf("Subtract failed: %v", sub)
	}
}

func TestPage(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	if got := Page(a, 1, 2); len(got) != 2 || got[0] != 1 {
		t.Fatalf("Page p1: %v", got)
	}
	if got := Page(a, 3, 2); len(got) != 1 || got[0] != 5 {
		t.Fatalf("Page p3: %v", got)
	}
	if got := Page(a, 5, 2); len(got) != 0 {
		t.Fatalf("Page out of range: %v", got)
	}
}
