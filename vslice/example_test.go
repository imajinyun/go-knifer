package vslice_test

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vslice"
)

func ExampleDistinct() {
	fmt.Println(vslice.Distinct([]int{1, 2, 2, 3, 3, 3}))
	// Output: [1 2 3]
}

func ExampleContains() {
	fmt.Println(vslice.Contains([]string{"go", "rust"}, "go"))
	// Output: true
}

func ExampleMap() {
	doubled := vslice.Map([]int{1, 2, 3}, func(n int) int { return n * 2 })
	fmt.Println(doubled)
	// Output: [2 4 6]
}

func ExampleFilter() {
	even := vslice.Filter([]int{1, 2, 3, 4}, func(n int) bool { return n%2 == 0 })
	fmt.Println(even)
	// Output: [2 4]
}

func ExampleUnion() {
	fmt.Println(vslice.Union([]int{1, 2}, []int{2, 3}))
	// Output: [1 2 3]
}
