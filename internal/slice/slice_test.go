package slice

import (
	"reflect"
	"testing"
)

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

func TestSliceLoStyleHelpers(t *testing.T) {
	values := []int{1, 2, 2, 3, 4}

	if got := Uniq(values); !reflect.DeepEqual(got, []int{1, 2, 3, 4}) {
		t.Fatalf("Uniq failed: %v", got)
	}
	if got := UniqBy([]string{"go", "js", "java"}, func(s string) int { return len(s) }); !reflect.DeepEqual(got, []string{"go", "java"}) {
		t.Fatalf("UniqBy failed: %v", got)
	}
	if got := Reject(values, func(v int) bool { return v%2 == 0 }); !reflect.DeepEqual(got, []int{1, 3}) {
		t.Fatalf("Reject failed: %v", got)
	}
	if got := FilterMap(values, func(v int) (int, bool) { return v * 10, v%2 == 0 }); !reflect.DeepEqual(got, []int{20, 20, 40}) {
		t.Fatalf("FilterMap failed: %v", got)
	}
	if got := FlatMap([]int{1, 2}, func(v int) []int { return []int{v, -v} }); !reflect.DeepEqual(got, []int{1, -1, 2, -2}) {
		t.Fatalf("FlatMap failed: %v", got)
	}
	if got := Reduce(values, 0, func(acc, v int) int { return acc + v }); got != 12 {
		t.Fatalf("Reduce failed: %v", got)
	}

	seen := []int{}
	ForEach([]int{1, 2, 3}, func(v int) { seen = append(seen, v) })
	if !reflect.DeepEqual(seen, []int{1, 2, 3}) {
		t.Fatalf("ForEach failed: %v", seen)
	}
	if got, ok := Find(values, func(v int) bool { return v > 2 }); !ok || got != 3 {
		t.Fatalf("Find = %v, %v", got, ok)
	}
	if got := FindIndex(values, func(v int) bool { return v == 3 }); got != 3 {
		t.Fatalf("FindIndex = %d", got)
	}

	words := []string{"go", "js", "rust", "java"}
	if got := GroupBy(words, func(s string) int { return len(s) }); !reflect.DeepEqual(got, map[int][]string{2: {"go", "js"}, 4: {"rust", "java"}}) {
		t.Fatalf("GroupBy failed: %v", got)
	}
	if got := CountBy(words, func(s string) int { return len(s) }); !reflect.DeepEqual(got, map[int]int{2: 2, 4: 2}) {
		t.Fatalf("CountBy failed: %v", got)
	}
	if got := KeyBy(words, func(s string) int { return len(s) }); !reflect.DeepEqual(got, map[int]string{2: "js", 4: "java"}) {
		t.Fatalf("KeyBy failed: %v", got)
	}
	if got := Associate(words, func(s string) (string, int) { return s, len(s) }); !reflect.DeepEqual(got, map[string]int{"go": 2, "js": 2, "rust": 4, "java": 4}) {
		t.Fatalf("Associate failed: %v", got)
	}
	if got := SliceToMap(words, func(s string) (int, string) { return len(s), s }); !reflect.DeepEqual(got, map[int]string{2: "js", 4: "java"}) {
		t.Fatalf("SliceToMap failed: %v", got)
	}

	if got := Chunk([]int{1, 2, 3, 4, 5}, 2); !reflect.DeepEqual(got, [][]int{{1, 2}, {3, 4}, {5}}) {
		t.Fatalf("Chunk failed: %v", got)
	}
	if got := Chunk([]int{1, 2}, 0); !reflect.DeepEqual(got, [][]int{}) {
		t.Fatalf("Chunk zero failed: %v", got)
	}
	if got := Flatten([][]int{{1}, {2, 3}}); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("Flatten failed: %v", got)
	}
	if got := Compact([]int{0, 1, 0, 2}); !reflect.DeepEqual(got, []int{1, 2}) {
		t.Fatalf("Compact failed: %v", got)
	}
	if got := PartitionBy([]int{1, 3, 2, 4, 5}, func(v int) bool { return v%2 == 0 }); !reflect.DeepEqual(got, [][]int{{1, 3}, {2, 4}, {5}}) {
		t.Fatalf("PartitionBy failed: %v", got)
	}
}
