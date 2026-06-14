package xml

import (
	"regexp"
	"testing"
)

func TestCleanHelpersWithOptions(t *testing.T) {
	if CleanInvalid("a\x00b\x08c") != "abc" {
		t.Fatal("CleanInvalid failed")
	}
	if CleanComment("<a><!-- hidden --><b/></a>") != "<a><b/></a>" {
		t.Fatal("CleanComment failed")
	}
	if got := CleanInvalidWithOptions("aXb", WithInvalidRegexp(regexp.MustCompile(`X`))); got != "ab" {
		t.Fatalf("CleanInvalidWithOptions = %q", got)
	}
	if got := CleanCommentWithOptions("<a>[hidden]<b/></a>", WithCommentRegexp(regexp.MustCompile(`\[[^]]+\]`))); got != "<a><b/></a>" {
		t.Fatalf("CleanCommentWithOptions = %q", got)
	}
}

func TestEscapeUnescapeHelpers(t *testing.T) {
	if Escape(`<a&"'>`) != "&lt;a&amp;&#34;&#39;&gt;" || Unescape("&lt;a&amp;&gt;") != "<a&>" {
		t.Fatal("escape/unescape failed")
	}
}
