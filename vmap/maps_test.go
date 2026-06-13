package vmap

import (
	"reflect"
	"sort"
	"testing"
)

func TestMapFacade(t *testing.T) {
	fromPairs := FromPairs(
		Pair[string, int]{Key: "a", Value: 1},
		Pair[string, int]{Key: "b", Value: 2},
	)
	if fromPairs["a"] != 1 || fromPairs["b"] != 2 {
		t.Fatalf("FromPairs failed: %v", fromPairs)
	}
	fromAny, err := OfE[string, int]("a", 1, "b", 2)
	if err != nil || fromAny["a"] != 1 || fromAny["b"] != 2 {
		t.Fatalf("OfE failed: %v, %v", fromAny, err)
	}

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

func TestMapConstructionLookupAndPredicateFacades(t *testing.T) {
	if got := New[string, int](); got == nil || len(got) != 0 {
		t.Fatalf("New = %#v", got)
	}
	if got := NewWithCap[string, int](-1); got == nil || len(got) != 0 {
		t.Fatalf("NewWithCap negative = %#v", got)
	}
	if got := Of[string, int]("a", 1, "a", 2); got["a"] != 2 {
		t.Fatalf("Of duplicate = %#v", got)
	}
	if _, err := OfE[string, int]("a"); err == nil {
		t.Fatal("OfE odd args error = nil")
	}
	if _, err := OfE[string, int](1, 2); err == nil {
		t.Fatal("OfE wrong key type error = nil")
	}
	if _, err := OfE[string, int]("a", "bad"); err == nil {
		t.Fatal("OfE wrong value type error = nil")
	}

	original := map[string]int{"a": 1}
	if got := OrEmpty(original); !reflect.DeepEqual(got, original) {
		t.Fatalf("OrEmpty existing = %#v", got)
	}
	if got := OrEmpty[string, int](nil); got == nil || len(got) != 0 {
		t.Fatalf("OrEmpty nil = %#v", got)
	}

	m := map[string]int{"a": 1, "b": 2, "c": 3}
	if !ContainsKey(m, "a") || ContainsKey(m, "missing") {
		t.Fatal("ContainsKey facade returned unexpected result")
	}
	if !ContainsValue(m, 2) || ContainsValue(m, 9) {
		t.Fatal("ContainsValue facade returned unexpected result")
	}
	if !Some(m, func(_ string, v int) bool { return v > 2 }) {
		t.Fatal("Some = false, want true")
	}
	if Every(m, func(_ string, v int) bool { return v < 3 }) {
		t.Fatal("Every = true, want false")
	}
	if got := Get(m, "missing"); got != 0 {
		t.Fatalf("Get missing = %d", got)
	}
	if got := GetOr(m, "missing", 9); got != 9 {
		t.Fatalf("GetOr missing = %d", got)
	}
	if got, ok := GetAny(m, "missing", "b"); !ok || got != 2 {
		t.Fatalf("GetAny = %d, %v", got, ok)
	}
	if _, ok := GetAny(m, "missing"); ok {
		t.Fatal("GetAny missing ok = true")
	}
	if k, v, ok := Find(m, func(_ string, v int) bool { return v == 3 }); !ok || k != "c" || v != 3 {
		t.Fatalf("Find = %q, %d, %v", k, v, ok)
	}
	if k, ok := FindKey(m, func(v int) bool { return v == 2 }); !ok || k != "b" {
		t.Fatalf("FindKey = %q, %v", k, ok)
	}
}

func TestMapTransformFilterAggregateFacades(t *testing.T) {
	m := map[string]int{"b": 2, "a": 1, "c": 3}
	if got := SortedKeys(m); !reflect.DeepEqual(got, []string{"a", "b", "c"}) {
		t.Fatalf("SortedKeys = %#v", got)
	}
	if got := SortedKeysFunc(m, func(a, b string) bool { return a > b }); !reflect.DeepEqual(got, []string{"c", "b", "a"}) {
		t.Fatalf("SortedKeysFunc = %#v", got)
	}
	if got := SortedValues(m); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("SortedValues = %#v", got)
	}
	keysOf := KeysOf(map[string]int{"a": 1, "b": 2, "c": 1}, 1)
	sort.Strings(keysOf)
	if !reflect.DeepEqual(keysOf, []string{"a", "c"}) {
		t.Fatalf("KeysOf = %#v", keysOf)
	}

	mapped := Map(m, func(k string, v int) (string, string) { return k + k, string(rune('0' + v)) })
	if !reflect.DeepEqual(mapped, map[string]string{"aa": "1", "bb": "2", "cc": "3"}) {
		t.Fatalf("Map = %#v", mapped)
	}
	if got := MapKeys(m, func(k string, _ int) string { return k + "!" }); !reflect.DeepEqual(got, map[string]int{"a!": 1, "b!": 2, "c!": 3}) {
		t.Fatalf("MapKeys = %#v", got)
	}
	if got := MapValues(m, func(_ string, v int) string { return string(rune('A' + v)) }); !reflect.DeepEqual(got, map[string]string{"a": "B", "b": "C", "c": "D"}) {
		t.Fatalf("MapValues = %#v", got)
	}
	toSlice := ToSlice(m, func(k string, v int) string { return k + string(rune('0'+v)) })
	sort.Strings(toSlice)
	if !reflect.DeepEqual(toSlice, []string{"a1", "b2", "c3"}) {
		t.Fatalf("ToSlice = %#v", toSlice)
	}
	if got := ToSlice(map[string]int(nil), func(k string, v int) string { return k + string(rune('0'+v)) }); got == nil || len(got) != 0 {
		t.Fatalf("ToSlice nil = %#v", got)
	}
	if got := Filter(m, func(_ string, v int) bool { return v%2 == 1 }); !reflect.DeepEqual(got, map[string]int{"a": 1, "c": 3}) {
		t.Fatalf("Filter = %#v", got)
	}
	if got := Reject(m, func(_ string, v int) bool { return v%2 == 1 }); !reflect.DeepEqual(got, map[string]int{"b": 2}) {
		t.Fatalf("Reject = %#v", got)
	}
	if got := FilterKeys(m, func(k string) bool { return k != "b" }); !reflect.DeepEqual(got, map[string]int{"a": 1, "c": 3}) {
		t.Fatalf("FilterKeys = %#v", got)
	}
	if got := FilterValues(m, func(v int) bool { return v >= 2 }); !reflect.DeepEqual(got, map[string]int{"b": 2, "c": 3}) {
		t.Fatalf("FilterValues = %#v", got)
	}
	matched, rest := Partition(m, func(_ string, v int) bool { return v >= 2 })
	if !reflect.DeepEqual(matched, map[string]int{"b": 2, "c": 3}) || !reflect.DeepEqual(rest, map[string]int{"a": 1}) {
		t.Fatalf("Partition matched=%#v rest=%#v", matched, rest)
	}

	seen := map[string]int{}
	ForEach(m, func(k string, v int) { seen[k] = v })
	if !reflect.DeepEqual(seen, m) {
		t.Fatalf("ForEach seen = %#v", seen)
	}
	if got := Reduce(m, 0, func(acc int, _ string, v int) int { return acc + v }); got != 6 {
		t.Fatalf("Reduce = %d", got)
	}
	if got := GroupBy([]string{"go", "js", "java"}, func(s string) int { return len(s) }); !reflect.DeepEqual(got, map[int][]string{2: {"go", "js"}, 4: {"java"}}) {
		t.Fatalf("GroupBy = %#v", got)
	}
	if got := CountBy([]string{"go", "js", "java"}, func(s string) int { return len(s) }); !reflect.DeepEqual(got, map[int]int{2: 2, 4: 1}) {
		t.Fatalf("CountBy = %#v", got)
	}
}

func TestMapSetAlgebraSelectionMutationComparisonFacades(t *testing.T) {
	a := map[string]int{"a": 1, "b": 2}
	b := map[string]int{"b": 20, "c": 3}
	c := map[string]int{"b": 200, "d": 4}

	dst := map[string]int{"a": 1, "shared": 1}
	MergeWithOverwrite(dst, map[string]int{"shared": 2}, map[string]int{"x": 3})
	if !reflect.DeepEqual(dst, map[string]int{"a": 1, "shared": 2, "x": 3}) {
		t.Fatalf("MergeWithOverwrite dst = %#v", dst)
	}
	MergeWithoutOverwrite(dst, map[string]int{"a": 9, "y": 4})
	if !reflect.DeepEqual(dst, map[string]int{"a": 1, "shared": 2, "x": 3, "y": 4}) {
		t.Fatalf("MergeWithoutOverwrite dst = %#v", dst)
	}
	if got := MergeCopyWithOverwrite(map[string]int{"k": 1}, map[string]int{"k": 2}); !reflect.DeepEqual(got, map[string]int{"k": 2}) {
		t.Fatalf("MergeCopyWithOverwrite = %#v", got)
	}
	if got := MergeCopyWithoutOverwrite(map[string]int{"k": 1}, map[string]int{"k": 2, "x": 3}); !reflect.DeepEqual(got, map[string]int{"k": 1, "x": 3}) {
		t.Fatalf("MergeCopyWithoutOverwrite = %#v", got)
	}

	if got := MergeFunc(func(old, new int) int { return old + new }, a, b, c); !reflect.DeepEqual(got, map[string]int{"a": 1, "b": 222, "c": 3, "d": 4}) {
		t.Fatalf("MergeFunc = %#v", got)
	}
	if got := Intersect(a, b, c); !reflect.DeepEqual(got, map[string]int{"b": 200}) {
		t.Fatalf("Intersect = %#v", got)
	}
	if got := Diff(a, b); !reflect.DeepEqual(got, map[string]int{"a": 1}) {
		t.Fatalf("Diff = %#v", got)
	}
	if got := SymmetricDiff(a, b); !reflect.DeepEqual(got, map[string]int{"a": 1, "c": 3}) {
		t.Fatalf("SymmetricDiff = %#v", got)
	}
	if got := Pick(a, "b", "missing"); !reflect.DeepEqual(got, map[string]int{"b": 2}) {
		t.Fatalf("Pick = %#v", got)
	}
	if got := Omit(a, "a"); !reflect.DeepEqual(got, map[string]int{"b": 2}) {
		t.Fatalf("Omit = %#v", got)
	}

	updated := Update[string, int](nil, a)
	if !reflect.DeepEqual(updated, a) {
		t.Fatalf("Update nil dst = %#v", updated)
	}
	clone := Clone(a)
	clone["a"] = 99
	if a["a"] != 1 || clone["a"] != 99 {
		t.Fatalf("Clone should not alias input: original=%#v clone=%#v", a, clone)
	}
	if got := Clone[string, int](nil); got == nil || len(got) != 0 {
		t.Fatalf("Clone nil = %#v", got)
	}
	Clear(updated)
	if len(updated) != 0 {
		t.Fatalf("Clear = %#v", updated)
	}
	if !Equal(a, map[string]int{"a": 1, "b": 2}) || Equal(a, b) {
		t.Fatal("Equal facade returned unexpected result")
	}
	if !EqualFunc(map[string]int{"a": 1}, map[string]string{"a": "1"}, func(i int, s string) bool { return string(rune('0'+i)) == s }) {
		t.Fatal("EqualFunc = false, want true")
	}
}
