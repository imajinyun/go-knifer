package vhttp_test

import (
	"errors"
	"fmt"

	"github.com/imajinyun/go-knifer"
	"github.com/imajinyun/go-knifer/vhttp"
)

func ExampleNewError() {
	err := vhttp.NewError("no response", nil)
	fmt.Println(errors.Is(err, knifer.ErrCodeInternal))
	// Output: true
}
