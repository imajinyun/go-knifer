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

func TestCompiledRegexpHelpersBoundaryPaths(t *testing.T) {
	re := regexp.MustCompile(`(?P<key>\w+)=(\d+)`)

	if got, ok := GetReOK(nil, "a=1", 0); ok || got != "" {
		t.Fatalf("GetReOK nil regexp = %q %v", got, ok)
	}
	if got, ok := GetReOK(re, "a=1", -1); ok || got != "" {
		t.Fatalf("GetReOK negative group = %q %v", got, ok)
	}
	if got, ok := GetReOK(re, "a=1", 3); ok || got != "" {
		t.Fatalf("GetReOK out of range = %q %v", got, ok)
	}

	if got := GetByNameRe(nil, "a=1", "key"); got != "" {
		t.Fatalf("GetByNameRe nil regexp = %q", got)
	}
	if got := GetByNameRe(re, "a=1", ""); got != "" {
		t.Fatalf("GetByNameRe empty name = %q", got)
	}
	if got := GetByNameRe(re, "zzz", "key"); got != "" {
		t.Fatalf("GetByNameRe no match = %q", got)
	}

	if got := GetAllGroupsRe(nil, "a=1", true, true); got != nil {
		t.Fatalf("GetAllGroupsRe nil regexp = %#v", got)
	}
	if got := GetAllGroupsRe(re, "a=1 b=22", false, true); !reflect.DeepEqual(got, []string{"a", "1", "b", "22"}) {
		t.Fatalf("GetAllGroupsRe without group0 = %#v", got)
	}

	if got := GetAllGroupNamesRe(nil, "a=1"); got != nil {
		t.Fatalf("GetAllGroupNamesRe nil regexp = %#v", got)
	}
	if got := GetAllGroupNamesRe(re, "zzz"); len(got) != 0 {
		t.Fatalf("GetAllGroupNamesRe no match = %#v", got)
	}

	if got := ExtractMultiRe(nil, "a=1", "$1"); got != "" {
		t.Fatalf("ExtractMultiRe nil regexp = %q", got)
	}
	if got := ExtractMultiRe(re, "zzz", "$1"); got != "" {
		t.Fatalf("ExtractMultiRe no match = %q", got)
	}

	holder := "zzz"
	if got := ExtractMultiAndDelPreRe(nil, &holder, "$1"); got != "" {
		t.Fatalf("ExtractMultiAndDelPreRe nil regexp = %q", got)
	}
	if got := ExtractMultiAndDelPreRe(re, nil, "$1"); got != "" {
		t.Fatalf("ExtractMultiAndDelPreRe nil holder = %q", got)
	}
	if got := ExtractMultiAndDelPreRe(re, &holder, "$1"); got != "" || holder != "zzz" {
		t.Fatalf("ExtractMultiAndDelPreRe no match = %q holder=%q", got, holder)
	}

	if got := ReplaceFirstRe(nil, "a=1", "x"); got != "a=1" {
		t.Fatalf("ReplaceFirstRe nil regexp = %q", got)
	}
	if got := ReplaceFirstRe(re, "", "x"); got != "" {
		t.Fatalf("ReplaceFirstRe empty content = %q", got)
	}
	if got := ReplaceFirstRe(re, "zzz", "x"); got != "zzz" {
		t.Fatalf("ReplaceFirstRe no match = %q", got)
	}

	if got := DelLastRe(nil, "a=1 b=22"); got != "a=1 b=22" {
		t.Fatalf("DelLastRe nil regexp = %q", got)
	}
	if got := DelLastRe(re, "zzz"); got != "zzz" {
		t.Fatalf("DelLastRe no match = %q", got)
	}

	if got := DelPreRe(nil, "a=1 b=22"); got != "a=1 b=22" {
		t.Fatalf("DelPreRe nil regexp = %q", got)
	}
	if got := DelPreRe(re, "zzz"); got != "zzz" {
		t.Fatalf("DelPreRe no match = %q", got)
	}

	if got := FindAllRe(nil, "a=1", 0); got != nil {
		t.Fatalf("FindAllRe nil regexp = %#v", got)
	}
	if got := FindAllRe(re, "a=1", -1); got != nil {
		t.Fatalf("FindAllRe negative group = %#v", got)
	}

	Each(nil, "a=1", func(MatchResult) { t.Fatal("nil regexp should not call consumer") })
	Each(re, "a=1", nil)

	if got := CountRe(nil, "a=1"); got != 0 {
		t.Fatalf("CountRe nil regexp = %d", got)
	}

	if got := IndexOfRe(nil, "a=1"); got != nil {
		t.Fatalf("IndexOfRe nil regexp = %#v", got)
	}
	if got := IndexOfRe(re, "zzz"); got != nil {
		t.Fatalf("IndexOfRe no match = %#v", got)
	}

	if got := LastIndexOfRe(nil, "a=1"); got != nil {
		t.Fatalf("LastIndexOfRe nil regexp = %#v", got)
	}
	if got := LastIndexOfRe(re, "zzz"); got != nil {
		t.Fatalf("LastIndexOfRe no match = %#v", got)
	}

	if !IsMatchRe(re, "a=1") {
		t.Fatal("IsMatchRe should match entire content")
	}
	if IsMatchRe(re, "prefix a=1") {
		t.Fatal("IsMatchRe should reject partial matches")
	}
	if IsMatchRe(nil, "a=1") {
		t.Fatal("IsMatchRe nil regexp should be false")
	}

	if got := ReplaceAllRe("a=1", nil, "x"); got != "a=1" {
		t.Fatalf("ReplaceAllRe nil regexp = %q", got)
	}
	if got := ReplaceAllRe("", re, "x"); got != "" {
		t.Fatalf("ReplaceAllRe empty content = %q", got)
	}

	if got := ReplaceAllFuncRe("a=1", nil, func(MatchResult) string { return "x" }); got != "a=1" {
		t.Fatalf("ReplaceAllFuncRe nil regexp = %q", got)
	}
	if got := ReplaceAllFuncRe("a=1", re, nil); got != "a=1" {
		t.Fatalf("ReplaceAllFuncRe nil func = %q", got)
	}
	if got := ReplaceAllFuncRe("", re, func(MatchResult) string { return "x" }); got != "" {
		t.Fatalf("ReplaceAllFuncRe empty content = %q", got)
	}
	if got := ReplaceAllFuncRe("zzz", re, func(MatchResult) string { return "x" }); got != "zzz" {
		t.Fatalf("ReplaceAllFuncRe no match = %q", got)
	}
}
