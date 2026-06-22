package vmap_test

import (
	"fmt"
	"sort"

	"github.com/imajinyun/go-knifer/vmap"
)

func ExampleIsEmpty() {
	fmt.Println(vmap.IsEmpty(map[string]int{}))
	fmt.Println(vmap.IsEmpty(map[string]int{"a": 1}))
	// Output:
	// true
	// false
}

func ExampleInverse() {
	inv := vmap.Inverse(map[string]int{"a": 1})
	fmt.Println(inv[1])
	// Output: a
}

func ExampleMerge() {
	merged := vmap.Merge(map[string]int{"a": 1}, map[string]int{"b": 2})
	fmt.Println(merged["a"], merged["b"])
	// Output: 1 2
}

func ExampleMapErr() {
	mapped, err := vmap.MapErr(map[string]int{"a": 1}, func(k string, v int) (string, string, error) {
		return k + k, fmt.Sprint(v), nil
	})
	fmt.Println(mapped["aa"], err)
	// Output: 1 <nil>
}

func ExampleMapKeysErr() {
	mapped, err := vmap.MapKeysErr(map[string]int{"a": 1}, func(k string, v int) (string, error) {
		return k + "!", nil
	})
	fmt.Println(mapped["a!"], err)
	// Output: 1 <nil>
}

func ExampleMapValuesErr() {
	mapped, err := vmap.MapValuesErr(map[string]int{"a": 1}, func(k string, v int) (string, error) {
		return fmt.Sprintf("%s=%d", k, v), nil
	})
	fmt.Println(mapped["a"], err)
	// Output: a=1 <nil>
}

func ExampleFilterErr() {
	filtered, err := vmap.FilterErr(map[string]int{"a": 1, "b": 2}, func(k string, v int) (bool, error) {
		return v%2 == 1, nil
	})
	fmt.Println(filtered["a"], filtered["b"], err)
	// Output: 1 0 <nil>
}

func ExampleReduceErr() {
	sum, err := vmap.ReduceErr(map[string]int{"a": 1}, 0, func(acc int, k string, v int) (int, error) {
		return acc + v, nil
	})
	fmt.Println(sum, err)
	// Output: 1 <nil>
}

func ExampleIter() {
	m := map[string]int{"b": 2, "a": 1}
	items := make([]string, 0, len(m))
	for key, value := range vmap.Iter(m) {
		items = append(items, fmt.Sprintf("%s=%d", key, value))
	}
	sort.Strings(items)
	fmt.Println(items)
	// Output: [a=1 b=2]
}

func ExampleIterKeys() {
	m := map[string]int{"b": 2, "a": 1}
	keys := make([]string, 0, len(m))
	for key := range vmap.IterKeys(m) {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	fmt.Println(keys)
	// Output: [a b]
}

func ExampleIterValues() {
	m := map[string]int{"b": 2, "a": 1}
	values := make([]int, 0, len(m))
	for value := range vmap.IterValues(m) {
		values = append(values, value)
	}
	sort.Ints(values)
	fmt.Println(values)
	// Output: [1 2]
}

func ExampleNew() {
	m := vmap.New[string, int]()
	m["count"] = 2

	fmt.Println(m["count"])
	// Output: 2
}

func ExampleOf() {
	m := vmap.Of[string, int]("a", 1, "b", 2)
	fmt.Println(m["a"], m["b"], len(m))
	// Output: 1 2 2
}

func ExampleContainsKey() {
	m := map[string]int{"a": 1}
	fmt.Println(vmap.ContainsKey(m, "a"))
	fmt.Println(vmap.ContainsKey(m, "b"))
	// Output:
	// true
	// false
}

func ExampleKeys() {
	keys := vmap.Keys(map[string]int{"b": 2, "a": 1})
	sort.Strings(keys)
	fmt.Println(keys)
	// Output: [a b]
}

func ExampleValues() {
	values := vmap.Values(map[string]int{"b": 2, "a": 1})
	sort.Ints(values)
	fmt.Println(values)
	// Output: [1 2]
}

func ExampleFilter() {
	filtered := vmap.Filter(map[string]int{"a": 1, "b": 2}, func(_ string, value int) bool {
		return value%2 == 0
	})
	fmt.Println(filtered)
	// Output: map[b:2]
}

func ExampleGroupBy() {
	grouped := vmap.GroupBy([]string{"go", "js", "ts"}, func(value string) int {
		return len(value)
	})
	sort.Strings(grouped[2])
	fmt.Println(grouped[2])
	// Output: [go js ts]
}

func ExampleCountBy() {
	counts := vmap.CountBy([]string{"go", "js", "rust"}, func(value string) int {
		return len(value)
	})
	fmt.Println(counts[2], counts[4])
	// Output: 2 1
}

func ExampleClone() {
	original := map[string][]int{"a": {1, 2}}
	cloned := vmap.Clone(original)
	cloned["a"] = append(cloned["a"], 3)

	fmt.Println(original["a"])
	fmt.Println(cloned["a"])
	// Output:
	// [1 2]
	// [1 2 3]
}

func ExampleDiff() {
	diff := vmap.Diff(map[string]int{"a": 1, "b": 2}, map[string]int{"b": 20})
	fmt.Println(diff)
	// Output: map[a:1]
}
