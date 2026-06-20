package vconv_test

import (
	"errors"
	"fmt"

	knifer "github.com/imajinyun/go-knifer"
	"github.com/imajinyun/go-knifer/vconv"
)

func ExampleToInt() {
	fmt.Println(vconv.ToInt("42"))
	fmt.Println(vconv.ToInt(true))
	// Output:
	// 42
	// 1
}

func ExampleToIntDefault() {
	fmt.Println(vconv.ToIntDefault("not-a-number", -1))
	// Output: -1
}

func ExampleToBool() {
	fmt.Println(vconv.ToBool("true"))
	fmt.Println(vconv.ToBool(0))
	// Output:
	// true
	// false
}

func ExampleToString() {
	fmt.Println(vconv.ToString(3.14))
	// Output: 3.14
}

func ExampleToBytes() {
	fmt.Println(string(vconv.ToBytes("go")))
	// Output: go
}

func ExampleToInt64E() {
	value, err := vconv.ToInt64E("42.9")
	fmt.Println(value)
	fmt.Println(err)
	// Output:
	// 42
	// <nil>
}

func ExampleToBoolE() {
	value, err := vconv.ToBoolE("maybe")
	fmt.Println(value)
	fmt.Println(errors.Is(err, vconv.ErrInvalidConversion))
	fmt.Println(errors.Is(err, knifer.ErrCodeInvalidInput))
	// Output:
	// false
	// true
	// true
}
