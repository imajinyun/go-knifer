package verr_test

import (
	"errors"
	"fmt"

	knifer "github.com/imajinyun/go-knifer"
	"github.com/imajinyun/go-knifer/verr"
)

func ExampleErrorIs() {
	inputErr := knifer.WrapError(knifer.ErrCodeInvalidInput, "bad value", errors.New("parse failed"))

	fmt.Println(verr.ErrorIs(inputErr, knifer.ErrCodeInvalidInput))
	fmt.Println(verr.ErrorIs(inputErr, knifer.ErrCodeInternal))
	// Output:
	// true
	// false
}

func ExampleGetStackWithOptions() {
	stack := verr.GetStackWithOptions(errors.New("plain"), verr.WithDebugStackFunc(func() []byte {
		return []byte("captured stack")
	}))

	fmt.Println(stack)
	// Output: captured stack
}
