package vregex

import (
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestRegexFindFacade(t *testing.T) {
	if Find(`\d+`, "ab123cd") != "123" || Find(`(`, "x") != "" {
		t.Fatal("Find failed")
	}
	if all := FindAll(`\d+`, "a1b22c333"); len(all) != 3 || all[2] != "333" {
		t.Fatalf("FindAll failed")
	}
	if Count(`\d+`, "a1b22") != 2 {
		t.Fatal("Count failed")
	}
	if CountRe(regexp.MustCompile(`\d+`), "a1b22") != 2 {
		t.Fatal("CountRe failed")
	}
	if match := IndexOf(`\d+`, "ab12"); match == nil || match.Start != 2 || match.Text != "12" {
		t.Fatalf("IndexOf = %#v", match)
	}
	if match := IndexOfRe(regexp.MustCompile(`\d+`), "ab12"); match == nil || match.Start != 2 || match.Text != "12" {
		t.Fatalf("IndexOfRe = %#v", match)
	}
	if match := LastIndexOf(`\d+`, "ab12cd34"); match == nil || match.Start != 6 || match.Text != "34" {
		t.Fatalf("LastIndexOf = %#v", match)
	}
	if match := LastIndexOfRe(regexp.MustCompile(`\d+`), "ab12cd34"); match == nil || match.Start != 6 || match.Text != "34" {
		t.Fatalf("LastIndexOfRe = %#v", match)
	}
	if n, ok := GetFirstNumber("a123"); !ok || n != 123 {
		t.Fatalf("GetFirstNumber = %d %v", n, ok)
	}

	var visited []string
	First(regexp.MustCompile(`\d+`), "a1b22", func(match MatchResult) {
		visited = append(visited, match.Text)
	})
	Each(regexp.MustCompile(`\d+`), "a1b22", func(match MatchResult) {
		visited = append(visited, match.Text)
	})
	if !reflect.DeepEqual(visited, []string{"1", "1", "22"}) {
		t.Fatalf("visitor matches = %#v", visited)
	}
}

func TestRegexFindFacadeWithOptions(t *testing.T) {
	opt := WithCompileFunc(func(pattern string) (*regexp.Regexp, error) {
		return regexp.Compile(strings.ReplaceAll(pattern, "TOKEN", `\d+`))
	})

	if got := FindWithOptions(`TOKEN`, "ab123", opt); got != "123" {
		t.Fatalf("FindWithOptions = %q", got)
	}
	if got := FindAllWithOptions(`TOKEN`, "a1b22", opt); !reflect.DeepEqual(got, []string{"1", "22"}) {
		t.Fatalf("FindAllWithOptions = %#v", got)
	}
	if got := CountWithOptions(`TOKEN`, "a1b22", opt); got != 2 {
		t.Fatalf("CountWithOptions = %d", got)
	}
	if got := IndexOfWithOptions(`TOKEN`, "ab12", opt); got == nil || got.Text != "12" || got.Start != 2 {
		t.Fatalf("IndexOfWithOptions = %#v", got)
	}
	if got := LastIndexOfWithOptions(`TOKEN`, "ab12cd34", opt); got == nil || got.Text != "34" || got.Start != 6 {
		t.Fatalf("LastIndexOfWithOptions = %#v", got)
	}
	if n, ok := GetFirstNumberWithOptions("a123", WithNumbersRegexp(regexp.MustCompile(`\d+`))); !ok || n != 123 {
		t.Fatalf("GetFirstNumberWithOptions = %d %v", n, ok)
	}
}
