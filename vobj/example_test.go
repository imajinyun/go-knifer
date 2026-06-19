package vobj_test

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vobj"
)

func ExampleEqual() {
	fmt.Println(vobj.Equal(42, 42))
	fmt.Println(vobj.Equal(42, 43))
	// Output:
	// true
	// false
}

func ExampleIsNil() {
	fmt.Println(vobj.IsNil(nil))
	fmt.Println(vobj.IsNil(0))
	// Output:
	// true
	// false
}

func ExampleContains() {
	fmt.Println(vobj.Contains([]string{"go", "knifer"}, "go"))
	fmt.Println(vobj.Contains([]string{"go", "knifer"}, "java"))
	// Output:
	// true
	// false
}
