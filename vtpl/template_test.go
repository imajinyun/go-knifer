package vtpl_test

import (
	"testing"

	"github.com/imajinyun/go-knifer/vtpl"
)

func TestRenderTemplateFacade(t *testing.T) {
	out, err := vtpl.RenderTemplate("hi {{.Name}}", map[string]string{"Name": "gokit"})
	if err != nil {
		t.Fatal(err)
	}
	if out != "hi gokit" {
		t.Fatalf("RenderTemplate() = %q", out)
	}
}
