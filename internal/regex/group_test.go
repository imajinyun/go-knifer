package regex

import (
	"reflect"
	"testing"
)

func TestGetGroupsAndNamedGroups(t *testing.T) {
	pattern := `(?<year>\d{4})-(?<month>\d{2})-(?<day>\d{2})`
	content := "date=2026-05-31"

	if got := GetGroup0(pattern, content); got != "2026-05-31" {
		t.Fatalf("GetGroup0 = %q", got)
	}
	if got := GetGroup1(pattern, content); got != "2026" {
		t.Fatalf("GetGroup1 = %q", got)
	}
	if got := GetByName(pattern, content, "month"); got != "05" {
		t.Fatalf("GetByName = %q", got)
	}
	groups := GetAllGroups(pattern, content, true, false)
	if !reflect.DeepEqual(groups, []string{"2026-05-31", "2026", "05", "31"}) {
		t.Fatalf("GetAllGroups = %#v", groups)
	}
	names := GetAllGroupNames(pattern, content)
	if names["year"] != "2026" || names["month"] != "05" || names["day"] != "31" {
		t.Fatalf("GetAllGroupNames = %#v", names)
	}
}
