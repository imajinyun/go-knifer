package log

import (
	"strings"
	"testing"
)

func TestFormatTemplatePlaceholder(t *testing.T) {
	template := strings.Join([]string{"hello {}", "age={}"}, ", ")
	args := []any{"world", 18}
	got := renderLogMessage(template, args...)
	want := "hello world, age=18"
	if got != want {
		t.Errorf("formatTemplate placeholder got=%q want=%q", got, want)
	}
}

func TestFormatTemplatePrintfFallback(t *testing.T) {
	got := renderLogMessage("a=%d, b=%s", 1, "x")
	want := "a=1, b=x"
	if got != want {
		t.Errorf("formatTemplate printf got=%q want=%q", got, want)
	}
}

func TestFormatTemplateNoArgs(t *testing.T) {
	if got := renderLogMessage("plain"); got != "plain" {
		t.Errorf("formatTemplate plain got=%q", got)
	}
}

func TestFormatTemplateConcat(t *testing.T) {
	args := []any{"a", "b", 1}
	if got := renderLogMessage("", args...); got != "ab1" {
		t.Errorf("formatTemplate concat got=%q", got)
	}
}

func TestFormatTemplateExtraPlaceholders(t *testing.T) {
	template := strings.Repeat("{}-", 2) + "{}"
	args := []any{"a"}
	got := renderLogMessage(template, args...)
	if got != "a-{}-{}" {
		t.Errorf("formatTemplate extra got=%q", got)
	}
}
