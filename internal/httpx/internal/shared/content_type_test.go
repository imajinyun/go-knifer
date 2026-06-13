package shared

import "testing"

func TestContentType(t *testing.T) {
	if got := ContentTypeJSON.String(); got != "application/json" {
		t.Fatalf("ContentTypeJSON.String = %q", got)
	}
	if got := ContentTypeJSON.WithCharset("utf-8"); got != "application/json;charset=utf-8" {
		t.Fatalf("ContentTypeJSON.WithCharset = %q", got)
	}
	if got := BuildContentType("text/plain", ""); got != "text/plain" {
		t.Fatalf("BuildContentType without charset = %q", got)
	}
	if !IsDefaultContentType("") || !IsDefaultContentType("Application/X-Www-Form-Urlencoded; charset=utf-8") {
		t.Fatal("IsDefaultContentType returned false for default values")
	}
	if IsDefaultContentType("application/json") {
		t.Fatal("IsDefaultContentType returned true for json")
	}
	if got := GuessContentType(` {"ok": true}`); got != ContentTypeJSON {
		t.Fatalf("GuessContentType json = %q", got)
	}
	if got := GuessContentType(" <xml></xml>"); got != ContentTypeXML {
		t.Fatalf("GuessContentType xml = %q", got)
	}
	if got := GuessContentType("plain"); got != "" {
		t.Fatalf("GuessContentType plain = %q", got)
	}
}
