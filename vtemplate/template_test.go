package vtemplate_test

import (
	"testing"

	"github.com/imajinyun/go-knifer/vtemplate"
)

func TestRenderTemplateFacade(t *testing.T) {
	out, err := vtemplate.RenderTemplate("hi {{.Name}}", map[string]string{"Name": "gokit"})
	if err != nil {
		t.Fatal(err)
	}
	if out != "hi gokit" {
		t.Fatalf("RenderTemplate() = %q", out)
	}
}
