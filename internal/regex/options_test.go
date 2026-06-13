package regex

import (
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestRegexHelpersWithOptions(t *testing.T) {
	compiler := func(pattern string) (*regexp.Regexp, error) {
		return regexp.Compile(strings.ReplaceAll(pattern, "TOKEN", `\d+`))
	}
	opt := WithCompileFunc(compiler)

	if got := GetGroup0WithOptions(`TOKEN`, "a123", opt); got != "123" {
		t.Fatalf("GetGroup0WithOptions = %q", got)
	}
	if got := GetGroup1WithOptions(`x(TOKEN)`, "x123", opt); got != "123" {
		t.Fatalf("GetGroup1WithOptions = %q", got)
	}
	if got := GetWithOptions(`x(TOKEN)`, "x123", 1, opt); got != "123" {
		t.Fatalf("GetWithOptions = %q", got)
	}
	if got, ok := GetOKWithOptions(`x(TOKEN)`, "x123", 1, opt); !ok || got != "123" {
		t.Fatalf("GetOKWithOptions = %q %v", got, ok)
	}
	if got := GetByNameWithOptions(`x(?<num>TOKEN)`, "x123", "num", opt); got != "123" {
		t.Fatalf("GetByNameWithOptions = %q", got)
	}
	if got := GetAllGroupsWithOptions(`x(TOKEN)`, "x123", true, false, opt); !reflect.DeepEqual(got, []string{"x123", "123"}) {
		t.Fatalf("GetAllGroupsWithOptions = %#v", got)
	}
	if got := GetAllGroupNamesWithOptions(`x(?<num>TOKEN)`, "x123", opt); got["num"] != "123" {
		t.Fatalf("GetAllGroupNamesWithOptions = %#v", got)
	}
	if got := ExtractMultiWithOptions(`x(TOKEN)`, "x123", "$1", opt); got != "123" {
		t.Fatalf("ExtractMultiWithOptions = %q", got)
	}
	holder := "x123y"
	if got := ExtractMultiAndDelPreWithOptions(`x(TOKEN)`, &holder, "$1", opt); got != "123" || holder != "y" {
		t.Fatalf("ExtractMultiAndDelPreWithOptions = %q holder=%q", got, holder)
	}
	if got := DelFirstWithOptions(`TOKEN`, "a123b456", opt); got != "ab456" {
		t.Fatalf("DelFirstWithOptions = %q", got)
	}
	if got := ReplaceFirstWithOptions(`TOKEN`, "a123b456", "X", opt); got != "aXb456" {
		t.Fatalf("ReplaceFirstWithOptions = %q", got)
	}
	if got := DelLastWithOptions(`TOKEN`, "a123b456", opt); got != "a123b" {
		t.Fatalf("DelLastWithOptions = %q", got)
	}
	if got := DelAllWithOptions(`TOKEN`, "a123b456", opt); got != "ab" {
		t.Fatalf("DelAllWithOptions = %q", got)
	}
	if got := DelPreWithOptions(`TOKEN`, "a123b", opt); got != "b" {
		t.Fatalf("DelPreWithOptions = %q", got)
	}
	if got := FindAllGroup0WithOptions(`TOKEN`, "a1b22", opt); !reflect.DeepEqual(got, []string{"1", "22"}) {
		t.Fatalf("FindAllGroup0WithOptions = %#v", got)
	}
	if got := FindAllGroup1WithOptions(`x(TOKEN)`, "x1x22", opt); !reflect.DeepEqual(got, []string{"1", "22"}) {
		t.Fatalf("FindAllGroup1WithOptions = %#v", got)
	}
	if got := FindAllWithOptions(`x(TOKEN)`, "x1x22", 1, opt); !reflect.DeepEqual(got, []string{"1", "22"}) {
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
}

func TestSpecializedRegexOptions(t *testing.T) {
	if n, ok := GetFirstNumberWithOptions("v12.34", WithNumbersRegexp(regexp.MustCompile(`[3-9]\d`))); !ok || n != 34 {
		t.Fatalf("GetFirstNumberWithOptions = %d %v", n, ok)
	}
	if got := TemplateVarsWithOptions("${3} $1", WithGroupVarRegexp(regexp.MustCompile(`\$\{(\d+)\}`))); !reflect.DeepEqual(got, []int{3}) {
		t.Fatalf("TemplateVarsWithOptions = %#v", got)
	}
	if got := GetByNameWithOptions(`(?<word>\w+)`, "abc", "word", WithNamedGroupNormalizer(func(pattern string) string {
		return strings.ReplaceAll(pattern, `(?<`, `(?P<`)
	})); got != "abc" {
		t.Fatalf("GetByNameWithOptions custom normalizer = %q", got)
	}
}
