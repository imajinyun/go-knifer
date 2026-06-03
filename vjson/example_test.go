package vjson_test

import (
	"errors"
	"fmt"

	"github.com/imajinyun/go-knifer"
	"github.com/imajinyun/go-knifer/vjson"
)

func ExampleToStr() {
	s, _ := vjson.ToStr(map[string]any{"name": "go"})
	fmt.Println(s)
	// Output: {"name":"go"}
}

func ExampleIsJSON() {
	fmt.Println(vjson.IsJSON(`{"a":1}`))
	fmt.Println(vjson.IsJSON(`not json`))
	// Output:
	// true
	// false
}

func ExampleGetByPath() {
	root, _ := vjson.Parse(`{"user":{"name":"go"}}`)
	fmt.Println(vjson.GetByPath(root, "user.name"))
	// Output: go
}

func ExampleParseObj_error() {
	_, err := vjson.ParseObj(`[1,2,3]`)
	fmt.Println(errors.Is(err, knifer.ErrCodeInvalidInput))
	// Output: true
}
