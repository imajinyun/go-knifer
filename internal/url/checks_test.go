package url

import "testing"

func TestURLChecksAndDataURI(t *testing.T) {
	if !IsHTTP("Http://example.com") || !IsHTTPS("HTTPS://example.com") {
		t.Fatal("scheme prefix checks failed")
	}
	if !IsWebURL("https://example.com/a") || IsWebURL("ftp://example.com/a") {
		t.Fatal("web URL checks failed")
	}
	if !IsAbsoluteURL("ftp://example.com/a") || IsAbsoluteURL("/relative") {
		t.Fatal("absolute URL checks failed")
	}
	if got := DataURIBase64("image/png", "AAAA"); got != "data:image/png;base64,AAAA" {
		t.Fatalf("DataURIBase64: %q", got)
	}
}
