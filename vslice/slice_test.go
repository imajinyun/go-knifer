package vslice

import (
	"reflect"
	"testing"
)

func TestSliceFacade(t *testing.T) {
	if !IsEmpty([]int{}) || !IsNotEmpty([]int{1}) {
		t.Fatal("empty checks failed")
	}
	values := []int{1, 2, 2, 3}
	if !Contains(values, 2) || IndexOf(values, 2) != 1 || LastIndexOf(values, 2) != 2 {
		t.Fatal("contains/index helpers failed")
	}
	reversed := Reverse([]int{1, 2, 3})
	if !reflect.DeepEqual(reversed, []int{3, 2, 1}) {
		t.Fatalf("Reverse failed: %v", reversed)
	}
	if got := Distinct(values); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("Distinct failed: %v", got)
	}
	if Join([]int{1, 2, 3}, ",") != "1,2,3" {
		t.Fatal("Join failed")
	}
	if got := Filter(values, func(v int) bool { return v%2 == 0 }); !reflect.DeepEqual(got, []int{2, 2}) {
		t.Fatalf("Filter failed: %v", got)
	}
	if got := Map([]int{1, 2}, func(v int) string { return string(rune('a' + v - 1)) }); !reflect.DeepEqual(got, []string{"a", "b"}) {
		t.Fatalf("Map failed: %v", got)
	}
	if got := Sub([]int{1, 2, 3, 4}, -3, -1); !reflect.DeepEqual(got, []int{2, 3}) {
		t.Fatalf("Sub failed: %v", got)
	}
	if got := Concat([]int{1}, []int{2, 3}); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("Concat failed: %v", got)
	}
	if got := Union([]int{1, 2}, []int{2, 3}); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("Union failed: %v", got)
	}
	if got := Intersection([]int{1, 2, 3}, []int{2, 3, 4}); !reflect.DeepEqual(got, []int{2, 3}) {
		t.Fatalf("Intersection failed: %v", got)
	}
	if got := Subtract([]int{1, 2, 3}, []int{2}); !reflect.DeepEqual(got, []int{1, 3}) {
		t.Fatalf("Subtract failed: %v", got)
	}
	if got := Page([]int{1, 2, 3, 4, 5}, 2, 2); !reflect.DeepEqual(got, []int{3, 4}) {
		t.Fatalf("Page failed: %v", got)
	}
}
