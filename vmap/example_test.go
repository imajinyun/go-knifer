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
