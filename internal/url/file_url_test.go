package url

import (
	neturl "net/url"
	"testing"
)

func TestHostDecodedPathAndJar(t *testing.T) {
	u, _ := neturl.Parse("https://example.com/a%20b?q=1")
	if got := Host(u).String(); got != "https://example.com" {
		t.Fatalf("Host: %q", got)
	}
	if got := DecodedPath(u); got != "/a b" {
		t.Fatalf("DecodedPath: %q", got)
	}
	jar, _ := neturl.Parse("file:///tmp/a.jar")
	if !IsFileURL(jar) || !IsJarFileURL(jar) {
		t.Fatal("jar file checks failed")
	}
}
