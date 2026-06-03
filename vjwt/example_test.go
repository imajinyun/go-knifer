package vjwt_test

import (
	"errors"
	"fmt"

	"github.com/imajinyun/go-knifer"
	"github.com/imajinyun/go-knifer/vjwt"
)

func ExampleNewJWTError() {
	err := vjwt.NewJWTError("token must not be blank")
	fmt.Println(errors.Is(err, knifer.ErrCodeInvalidInput))
	// Output: true
}
