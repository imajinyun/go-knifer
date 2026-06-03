package vstr_test

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vstr"
)

func ExampleToCamelCase() {
	fmt.Println(vstr.ToCamelCase("hello_world"))
	// Output: helloWorld
}

func ExampleToUnderlineCase() {
	fmt.Println(vstr.ToUnderlineCase("HelloWorld"))
	// Output: hello_world
}

func ExampleIsBlank() {
	fmt.Println(vstr.IsBlank("  "))
	fmt.Println(vstr.IsBlank("go"))
	// Output:
	// true
	// false
}

func ExampleReverse() {
	fmt.Println(vstr.Reverse("你好"))
	// Output: 好你
}

func ExampleSub() {
	fmt.Println(vstr.Sub("你好世界", 1, 3))
	// Output: 好世
}
