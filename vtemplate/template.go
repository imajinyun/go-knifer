package vtemplate

import templateimpl "github.com/imajinyun/go-knifer/internal/template"

// Render renders a Go html/template string with data.
func Render(tpl string, data any) (string, error) { return templateimpl.Render(tpl, data) }

// RenderTemplate renders a Go html/template string with data.
func RenderTemplate(tpl string, data any) (string, error) {
	return templateimpl.RenderTemplate(tpl, data)
}
