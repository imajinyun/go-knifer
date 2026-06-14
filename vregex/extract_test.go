package vregex

import (
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestRegexExtractFacade(t *testing.T) {
	if got := ExtractMulti(`(\d+)年(\d+)月`, "2026年5月", `$1-$2`); got != "2026-5" {
		t.Fatalf("ExtractMulti = %q", got)
	}
	if got := ExtractMultiRe(regexp.MustCompile(`(\d+)年(\d+)月`), "2026年5月", `$1-$2`); got != "2026-5" {
		t.Fatalf("ExtractMultiRe = %q", got)
	}
	holder := "x123y"
	if got := ExtractMultiAndDelPre(`x(\d+)`, &holder, "$1"); got != "123" || holder != "y" {
		t.Fatalf("ExtractMultiAndDelPre = %q holder=%q", got, holder)
	}
	holder = "x123y"
	if got := ExtractMultiAndDelPreRe(regexp.MustCompile(`x(\d+)`), &holder, "$1"); got != "123" || holder != "y" {
		t.Fatalf("ExtractMultiAndDelPreRe = %q holder=%q", got, holder)
	}
	if got := TemplateVars("$10-$2-$1"); !reflect.DeepEqual(got, []int{10, 2, 1}) {
		t.Fatalf("TemplateVars = %#v", got)
	}
}

func TestRegexExtractFacadeWithOptions(t *testing.T) {
	opt := WithCompileFunc(func(pattern string) (*regexp.Regexp, error) {
		return regexp.Compile(strings.ReplaceAll(pattern, "TOKEN", `\d+`))
	})

	if got := ExtractMultiWithOptions(`x(TOKEN)`, "x123", "$1", opt); got != "123" {
		t.Fatalf("ExtractMultiWithOptions = %q", got)
	}
	holder := "x123y"
	if got := ExtractMultiAndDelPreWithOptions(`x(TOKEN)`, &holder, "$1", opt); got != "123" || holder != "y" {
		t.Fatalf("ExtractMultiAndDelPreWithOptions = %q holder=%q", got, holder)
	}
	if got := TemplateVarsWithOptions("#10-#2", WithGroupVarRegexp(regexp.MustCompile(`#(\d+)`))); !reflect.DeepEqual(got, []int{10, 2}) {
		t.Fatalf("TemplateVarsWithOptions = %#v", got)
	}
}
