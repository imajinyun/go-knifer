package template

import "testing"

func TestRenderTemplate(t *testing.T) {
	out, err := RenderTemplate("hello {{.Name}}", map[string]string{"Name": "gokit"})
	if err != nil {
		t.Fatal(err)
	}
	if out != "hello gokit" {
		t.Fatalf("RenderTemplate() = %q", out)
	}
}
