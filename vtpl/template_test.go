package vtpl_test

import (
	"errors"
	"html/template"
	"io"
	"strings"
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

func TestRenderFacadeAndOptions(t *testing.T) {
	out, err := vtpl.Render("hello {{.Name}}", map[string]string{"Name": "tpl"})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if out != "hello tpl" {
		t.Fatalf("Render = %q", out)
	}

	out, err = vtpl.RenderWithOptions(
		"hello [[upper .Name]]",
		map[string]string{"Name": "tpl"},
		vtpl.WithTemplateName("custom"),
		vtpl.WithDelims("[[", "]]"),
		vtpl.WithFuncMap(template.FuncMap{"upper": strings.ToUpper}),
	)
	if err != nil {
		t.Fatalf("RenderWithOptions: %v", err)
	}
	if out != "hello TPL" {
		t.Fatalf("RenderWithOptions = %q", out)
	}
}

func TestRenderWithProviderOptions(t *testing.T) {
	factoryCalled := false
	parserCalled := false
	executorCalled := false

	out, err := vtpl.RenderWithOptions(
		"ignored",
		map[string]string{"Name": "provider"},
		vtpl.WithTemplateFactory(func(name string) *template.Template {
			factoryCalled = name == "provider-name"
			return template.New(name)
		}),
		vtpl.WithTemplateName("provider-name"),
		vtpl.WithTemplateParser(func(t *template.Template, source string) (*template.Template, error) {
			parserCalled = source == "ignored"
			return t.Parse("provider {{.Name}}")
		}),
		vtpl.WithTemplateExecutor(func(t *template.Template, w io.Writer, data any) error {
			executorCalled = true
			return t.Execute(w, data)
		}),
	)
	if err != nil {
		t.Fatalf("RenderWithOptions providers: %v", err)
	}
	if out != "provider provider" {
		t.Fatalf("RenderWithOptions providers = %q", out)
	}
	if !factoryCalled || !parserCalled || !executorCalled {
		t.Fatalf("provider calls factory=%v parser=%v executor=%v", factoryCalled, parserCalled, executorCalled)
	}
}

func TestRenderWithOptionsPropagatesProviderErrors(t *testing.T) {
	if _, err := vtpl.RenderWithOptions("bad", nil, vtpl.WithTemplateParser(func(*template.Template, string) (*template.Template, error) {
		return nil, errors.New("parse failed")
	})); err == nil {
		t.Fatal("RenderWithOptions parser error = nil")
	}

	if _, err := vtpl.RenderWithOptions("ok", nil, vtpl.WithTemplateExecutor(func(*template.Template, io.Writer, any) error {
		return errors.New("execute failed")
	})); err == nil {
		t.Fatal("RenderWithOptions executor error = nil")
	}
}
