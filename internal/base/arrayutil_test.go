package base

import "testing"

// 对应 hutool-core ArrayUtilTest。

func TestSliceBasic(t *testing.T) {
	if !SliceIsEmpty([]int{}) || SliceIsEmpty([]int{1}) {
		t.Fatalf("SliceIsEmpty failed")
	}
	if !SliceIsNotEmpty([]int{1}) {
		t.Fatalf("SliceIsNotEmpty failed")
	}
	if !SliceContains([]string{"a", "b"}, "b") || SliceContains([]string{"a"}, "x") {
		t.Fatalf("SliceContains failed")
	}
	if SliceIndexOf([]int{1, 2, 3}, 2) != 1 {
		t.Fatalf("SliceIndexOf failed")
	}
	if SliceLastIndexOf([]int{1, 2, 1}, 1) != 2 {
		t.Fatalf("SliceLastIndexOf failed")
	}
}

func TestSliceMutation(t *testing.T) {
	a := []int{1, 2, 3, 4}
	SliceReverse(a)
	if a[0] != 4 || a[3] != 1 {
		t.Fatalf("SliceReverse failed: %v", a)
	}
	d := SliceDistinct([]int{1, 2, 1, 3, 2})
	if len(d) != 3 || d[0] != 1 || d[1] != 2 || d[2] != 3 {
		t.Fatalf("SliceDistinct failed: %v", d)
	}
	if SliceJoin([]int{1, 2, 3}, ",") != "1,2,3" {
		t.Fatalf("SliceJoin failed")
	}
	f := SliceFilter([]int{1, 2, 3, 4}, func(v int) bool { return v%2 == 0 })
	if len(f) != 2 || f[0] != 2 || f[1] != 4 {
		t.Fatalf("SliceFilter failed: %v", f)
	}
	m := SliceMap([]int{1, 2, 3}, func(v int) int { return v * v })
	if m[0] != 1 || m[1] != 4 || m[2] != 9 {
		t.Fatalf("SliceMap failed: %v", m)
	}
}

func TestSliceSubAndConcat(t *testing.T) {
	got := SliceSub([]int{1, 2, 3, 4, 5}, 1, 4)
	if len(got) != 3 || got[0] != 2 || got[2] != 4 {
		t.Fatalf("SliceSub failed: %v", got)
	}
	got = SliceSub([]int{1, 2, 3, 4, 5}, -3, -1)
	if len(got) != 2 || got[0] != 3 || got[1] != 4 {
		t.Fatalf("SliceSub neg failed: %v", got)
	}
	c := SliceConcat([]int{1}, []int{2, 3}, []int{4})
	if len(c) != 4 || c[3] != 4 {
		t.Fatalf("SliceConcat failed: %v", c)
	}
}
