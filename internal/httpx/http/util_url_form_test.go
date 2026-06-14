package http

import (
	"strings"
	"testing"
)

func TestURLWithFormFunc(t *testing.T) {
	got := URLWithForm("http://api.gokit.cn/login", map[string]any{"a": 1})
	if !strings.Contains(got, "?a=1") {
		t.Fatalf("URLWithForm: %q", got)
	}
	got2 := URLWithForm("http://api.gokit.cn/login?type=aaa", map[string]any{"x": "y"})
	if !strings.Contains(got2, "type=aaa") || !strings.Contains(got2, "x=y") || !strings.Contains(got2, "&") {
		t.Fatalf("URLWithForm2: %q", got2)
	}
}
