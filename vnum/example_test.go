package vnum_test

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vnum"
)

func ExampleRound() {
	fmt.Println(vnum.Round(3.14159, 2))
	// Output: 3.14
}

func ExampleAddStr() {
	// AddStr keeps exact precision and avoids float rounding errors.
	fmt.Println(vnum.AddStr("0.1", "0.2").FloatString(1))
	// Output: 0.3
}

func ExampleIsPrimes() {
	fmt.Println(vnum.IsPrimes(7))
	fmt.Println(vnum.IsPrimes(8))
	// Output:
	// true
	// false
}

func ExampleMax() {
	fmt.Println(vnum.Max(3, 7, 1))
	// Output: 7
}

func ExampleCalculate() {
	result, _ := vnum.Calculate("1 + 2 * 3")
	fmt.Println(result)
	// Output: 7
}
