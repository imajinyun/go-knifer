package vregex

import (
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestRegexGroupFacade(t *testing.T) {
	if got := GetGroup0(`(\d+)`, "abc123"); got != "123" {
		t.Fatalf("GetGroup0 = %q", got)
	}
	if got := GetGroup1(`(\d+)`, "abc123"); got != "123" {
		t.Fatalf("GetGroup1 = %q", got)
	}
	if got := Get(`(a)(b)`, "ab", 2); got != "b" {
		t.Fatalf("Get = %q", got)
	}
	if got, ok := GetOK(`(a)(b)`, "ab", 2); !ok || got != "b" {
		t.Fatalf("GetOK = %q %v", got, ok)
	}
	if got := GetByName(`(?<word>\w+)-(?<num>\d+)`, "abc-123", "num"); got != "123" {
		t.Fatalf("GetByName = %q", got)
	}
	if got := GetRe(regexp.MustCompile(`(\d+)`), "abc123", 1); got != "123" {
		t.Fatalf("GetRe = %q", got)
	}
	if got := GetByNameRe(regexp.MustCompile(`(?P<num>\d+)`), "abc123", "num"); got != "123" {
		t.Fatalf("GetByNameRe = %q", got)
	}
	if got := GetAllGroups(`(a)(b)`, "ab", true, false); !reflect.DeepEqual(got, []string{"ab", "a", "b"}) {
		t.Fatalf("GetAllGroups = %#v", got)
	}
	if got := GetAllGroupsRe(regexp.MustCompile(`(a)(b)`), "ab", true, false); !reflect.DeepEqual(got, []string{"ab", "a", "b"}) {
		t.Fatalf("GetAllGroupsRe = %#v", got)
	}
	if got := GetAllGroupNames(`(?<word>\w+)-(?<num>\d+)`, "abc-123"); got["num"] != "123" {
		t.Fatalf("GetAllGroupNames = %#v", got)
	}
	if got := GetAllGroupNamesRe(regexp.MustCompile(`(?P<num>\d+)`), "abc123"); got["num"] != "123" {
		t.Fatalf("GetAllGroupNamesRe = %#v", got)
	}
	if got := FindAllGroup0(`(\d+)`, "a1b22"); !reflect.DeepEqual(got, []string{"1", "22"}) {
		t.Fatalf("FindAllGroup0 = %#v", got)
	}
	if got := FindAllGroup1(`(\d+)`, "a1b22"); !reflect.DeepEqual(got, []string{"1", "22"}) {
		t.Fatalf("FindAllGroup1 = %#v", got)
	}
	if got := FindAllGroup(`x(\d+)`, "x1x22", 1); !reflect.DeepEqual(got, []string{"1", "22"}) {
		t.Fatalf("FindAllGroup = %#v", got)
	}
	if got := FindAllGroupRe(regexp.MustCompile(`x(\d+)`), "x1x22", 1); !reflect.DeepEqual(got, []string{"1", "22"}) {
		t.Fatalf("FindAllGroupRe = %#v", got)
	}
}

func TestRegexGroupFacadeWithOptions(t *testing.T) {
	opt := WithCompileFunc(func(pattern string) (*regexp.Regexp, error) {
		return regexp.Compile(strings.ReplaceAll(pattern, "TOKEN", `\d+`))
	})

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
	if got := FindAllGroup0WithOptions(`TOKEN`, "a1b22", opt); !reflect.DeepEqual(got, []string{"1", "22"}) {
		t.Fatalf("FindAllGroup0WithOptions = %#v", got)
	}
	if got := FindAllGroup1WithOptions(`x(TOKEN)`, "x1x22", opt); !reflect.DeepEqual(got, []string{"1", "22"}) {
		t.Fatalf("FindAllGroup1WithOptions = %#v", got)
	}
	if got := FindAllGroupWithOptions(`x(TOKEN)`, "x1x22", 1, opt); !reflect.DeepEqual(got, []string{"1", "22"}) {
		t.Fatalf("FindAllGroupWithOptions = %#v", got)
	}
}
