package regex

import (
	"reflect"
	"regexp"
	"testing"
)

func TestCompiledRegexpHelpers(t *testing.T) {
	re := regexp.MustCompile(`(\w+)=(\d+)`)
	var seen []string
	Each(re, "a=1 b=22", func(m MatchResult) { seen = append(seen, m.Groups[2]) })
	if !reflect.DeepEqual(seen, []string{"1", "22"}) {
		t.Fatalf("Each = %#v", seen)
	}
	if got := ReplaceFirstRe(re, "a=1 b=22", `$1:x`); got != "a:x b=22" {
		t.Fatalf("ReplaceFirstRe = %q", got)
	}
}
