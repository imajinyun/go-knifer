package regex

import "testing"

func TestExtractAndDelete(t *testing.T) {
	if got := ExtractMulti(`(.*?)年(.*?)月`, "2026年5月", "$1-$2"); got != "2026-5" {
		t.Fatalf("ExtractMulti = %q", got)
	}
	content := "prefix 2026年5月 suffix"
	if got := ExtractMultiAndDelPre(`(\d+)年(\d+)月`, &content, "$1-$2"); got != "2026-5" {
		t.Fatalf("ExtractMultiAndDelPre = %q", got)
	}
	if content != " suffix" {
		t.Fatalf("content after delete = %q", content)
	}

	if got := DelFirst(`\d+`, "a123b456"); got != "ab456" {
		t.Fatalf("DelFirst = %q", got)
	}
	if got := DelLast(`\d+`, "a123b456"); got != "a123b" {
		t.Fatalf("DelLast = %q", got)
	}
	if got := DelAll(`\d+`, "a123b456"); got != "ab" {
		t.Fatalf("DelAll = %q", got)
	}
	if got := DelPre(`\d+`, "abc123xyz"); got != "xyz" {
		t.Fatalf("DelPre = %q", got)
	}
}
