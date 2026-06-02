package slice

import "testing"

// Tests cover the utility toolkit-core ArrayUtilTest.

func TestSliceBasic(t *testing.T) {
	if !IsEmpty([]int{}) || IsEmpty([]int{1}) {
		t.Fatalf("IsEmpty failed")
	}
	if !IsNotEmpty([]int{1}) {
		t.Fatalf("IsNotEmpty failed")
	}
	if !Contains([]string{"a", "b"}, "b") || Contains([]string{"a"}, "x") {
		t.Fatalf("Contains failed")
	}
	if IndexOf([]int{1, 2, 3}, 2) != 1 {
		t.Fatalf("IndexOf failed")
	}
	if LastIndexOf([]int{1, 2, 1}, 1) != 2 {
		t.Fatalf("LastIndexOf failed")
	}
}

func TestSliceMutation(t *testing.T) {
	a := []int{1, 2, 3, 4}
	Reverse(a)
	if a[0] != 4 || a[3] != 1 {
		t.Fatalf("Reverse failed: %v", a)
	}
	d := Distinct([]int{1, 2, 1, 3, 2})
	if len(d) != 3 || d[0] != 1 || d[1] != 2 || d[2] != 3 {
		t.Fatalf("Distinct failed: %v", d)
	}
	if Join([]int{1, 2, 3}, ",") != "1,2,3" {
		t.Fatalf("Join failed")
	}
	f := Filter([]int{1, 2, 3, 4}, func(v int) bool { return v%2 == 0 })
	if len(f) != 2 || f[0] != 2 || f[1] != 4 {
		t.Fatalf("Filter failed: %v", f)
	}
	m := Map([]int{1, 2, 3}, func(v int) int { return v * v })
	if m[0] != 1 || m[1] != 4 || m[2] != 9 {
		t.Fatalf("Map failed: %v", m)
	}
}

func TestSliceSubAndConcat(t *testing.T) {
	got := Sub([]int{1, 2, 3, 4, 5}, 1, 4)
	if len(got) != 3 || got[0] != 2 || got[2] != 4 {
		t.Fatalf("Sub failed: %v", got)
	}
	got = Sub([]int{1, 2, 3, 4, 5}, -3, -1)
	if len(got) != 2 || got[0] != 3 || got[1] != 4 {
		t.Fatalf("Sub neg failed: %v", got)
	}
	c := Concat([]int{1}, []int{2, 3}, []int{4})
	if len(c) != 4 || c[3] != 4 {
		t.Fatalf("Concat failed: %v", c)
	}
}
