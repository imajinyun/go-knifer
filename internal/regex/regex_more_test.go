package regex

import (
	"reflect"
	"regexp"
	"testing"
)

func TestWithDotAll(t *testing.T) {
	// WithDotAll(false) should make '.' not match newlines.
	cfg := applyOptions([]Option{WithDotAll(false)})
	if cfg.dotAll {
		t.Error("WithDotAll(false) should set dotAll to false")
	}
	// WithDotAll(true) should make '.' match newlines (default behavior).
	cfg2 := applyOptions([]Option{WithDotAll(true)})
	if !cfg2.dotAll {
		t.Error("WithDotAll(true) should set dotAll to true")
	}
}

func TestWithNamedGroupRegexp(t *testing.T) {
	re := regexp.MustCompile(`\(\?P<(\w+)>`)
	cfg := applyOptions([]Option{WithNamedGroupRegexp(re)})
	if cfg.namedGroupRegexp != re {
		t.Error("WithNamedGroupRegexp should set the custom regexp")
	}
}

func TestGet(t *testing.T) {
	if got := Get(`(\d+)`, "abc123def", 1); got != "123" {
		t.Errorf("Get = %q, want %q", got, "123")
	}
	if got := Get(`(\d+)`, "abc", 1); got != "" {
		t.Errorf("Get no match = %q, want %q", got, "")
	}
}

func TestGetOK(t *testing.T) {
	got, ok := GetOK(`(\d+)`, "abc123def", 1)
	if !ok || got != "123" {
		t.Errorf("GetOK = %q %v, want %q true", got, ok, "123")
	}
	got, ok = GetOK(`(\d+)`, "abc", 1)
	if ok || got != "" {
		t.Errorf("GetOK no match = %q %v, want \"\" false", got, ok)
	}
}

func TestFirst(t *testing.T) {
	re := regexp.MustCompile(`(\d+)`)
	var called bool
	First(re, "abc123def", func(result MatchResult) {
		called = true
		if result.Text != "123" {
			t.Errorf("First consumer result.Text = %q, want %q", result.Text, "123")
		}
	})
	if !called {
		t.Error("First should call the consumer")
	}

	// No match.
	called = false
	First(re, "abc", func(MatchResult) { called = true })
	if called {
		t.Error("First should not call the consumer when there is no match")
	}

	// Nil consumer should not panic.
	First(re, "abc", nil)

	// Nil regexp should not panic.
	First(nil, "abc", func(MatchResult) {})
}

func TestReplaceFirst(t *testing.T) {
	if got := ReplaceFirst(`\d+`, "a123b456", "X"); got != "aXb456" {
		t.Errorf("ReplaceFirst = %q, want %q", got, "aXb456")
	}
	// No match should return content unchanged.
	if got := ReplaceFirst(`\d+`, "abc", "X"); got != "abc" {
		t.Errorf("ReplaceFirst no match = %q, want %q", got, "abc")
	}
}

func TestDelAllRe(t *testing.T) {
	re := regexp.MustCompile(`\d+`)
	if got := DelAllRe(re, "a1b22c333"); got != "abc" {
		t.Errorf("DelAllRe = %q, want %q", got, "abc")
	}
	// No match should return content unchanged.
	if got := DelAllRe(re, "abc"); got != "abc" {
		t.Errorf("DelAllRe no match = %q, want %q", got, "abc")
	}
	// Nil regexp should return content unchanged.
	if got := DelAllRe(nil, "abc123"); got != "abc123" {
		t.Errorf("DelAllRe nil re = %q, want %q", got, "abc123")
	}
	// Empty content.
	if got := DelAllRe(re, ""); got != "" {
		t.Errorf("DelAllRe empty = %q, want %q", got, "")
	}
}

func TestFindAllGroup0(t *testing.T) {
	got := FindAllGroup0(`\d+`, "a1b22c333")
	want := []string{"1", "22", "333"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FindAllGroup0 = %#v, want %#v", got, want)
	}
	// No match.
	if got := FindAllGroup0(`\d+`, "abc"); len(got) != 0 {
		t.Errorf("FindAllGroup0 no match = %#v, want empty", got)
	}
}

func TestFindAll(t *testing.T) {
	got := FindAll(`x(\d+)`, "x1x22x333", 1)
	want := []string{"1", "22", "333"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FindAll = %#v, want %#v", got, want)
	}
	// Group 0 returns full matches.
	got0 := FindAll(`\d+`, "a1b22", 0)
	want0 := []string{"1", "22"}
	if !reflect.DeepEqual(got0, want0) {
		t.Errorf("FindAll group0 = %#v, want %#v", got0, want0)
	}
}

func TestContainsWithOptions(t *testing.T) {
	// With dotAll default, '.' matches newlines.
	if !ContainsWithOptions("a.b", "a\nb") {
		t.Error("ContainsWithOptions should match with dotAll by default")
	}
	// With DotAll(false), '.' should not match newlines.
	if ContainsWithOptions("a.b", "a\nb", WithDotAll(false)) {
		t.Error("ContainsWithOptions with dotAll(false) should not match cross-line")
	}
}

func TestContainsRe(t *testing.T) {
	re := regexp.MustCompile(`\d+`)
	if !ContainsRe(re, "abc123") {
		t.Error("ContainsRe should return true when content has a match")
	}
	if ContainsRe(re, "abc") {
		t.Error("ContainsRe should return false when content has no match")
	}
	if ContainsRe(nil, "abc") {
		t.Error("ContainsRe should return false for nil regexp")
	}
}

func TestTemplateVars(t *testing.T) {
	got := TemplateVars("$1 and $2 and $1")
	want := []int{2, 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("TemplateVars = %#v, want %#v", got, want)
	}
	// No placeholders.
	if got := TemplateVars("no dollars here"); len(got) != 0 {
		t.Errorf("TemplateVars no matches = %#v, want empty", got)
	}
}
