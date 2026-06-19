package vtpl_test

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/imajinyun/go-knifer/vtpl"
)

func ExampleRender() {
	result, err := vtpl.Render("Hello, {{.Name}}!", map[string]any{"Name": "World"})
	fmt.Println(result)
	fmt.Println(err)
	// Output:
	// Hello, World!
	// <nil>
}

func ExampleRenderTemplate() {
	result, err := vtpl.RenderTemplate("{{.Greeting}}, {{.Name}}!", map[string]string{
		"Greeting": "Hi",
		"Name":     "Gopher",
	})

	fmt.Println(result)
	fmt.Println(err)
	// Output:
	// Hi, Gopher!
	// <nil>
}

func ExampleRenderWithOptions() {
	result, err := vtpl.RenderWithOptions(
		"Hello [[upper .Name]]",
		map[string]string{"Name": "gopher"},
		vtpl.WithDelims("[[", "]]"),
		vtpl.WithFuncMap(template.FuncMap{"upper": strings.ToUpper}),
	)

	fmt.Println(result)
	fmt.Println(err)
	// Output:
	// Hello GOPHER
	// <nil>
}
