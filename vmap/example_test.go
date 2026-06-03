package vmap_test

import (
	"fmt"

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
