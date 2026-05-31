package template

import (
	"bytes"
	"html/template"
)

// Render renders a Go html/template string with data.
func Render(tpl string, data any) (string, error) {
	t, err := template.New("go-knifer-template").Parse(tpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderTemplate renders a Go html/template string with data.
func RenderTemplate(tpl string, data any) (string, error) { return Render(tpl, data) }
