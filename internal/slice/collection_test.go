package slice

import (
	"slices"
	"testing"
)

// Tests cover the utility toolkit-core CollUtilTest and ListUtilTest.

func TestUnionIntersectionSubtract(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{2, 3, 4}
	u := Union(a, b)
	slices.Sort(u)
	if len(u) != 4 {
		t.Fatalf("Union failed: %v", u)
	}
	in := Intersection(a, b)
	slices.Sort(in)
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
