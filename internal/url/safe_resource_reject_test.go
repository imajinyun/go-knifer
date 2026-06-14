package url

import "testing"

func TestOpenSafeRejectsLocalAndPrivateResources(t *testing.T) {
	if _, err := OpenSafe("file:///tmp/secret.txt"); err == nil {
		t.Fatal("OpenSafe should reject file URLs")
	}
	if _, err := OpenSafe("/tmp/secret.txt"); err == nil {
		t.Fatal("OpenSafe should reject plain file paths")
	}
	if _, err := OpenSafe("http://127.0.0.1/config.yaml"); err == nil {
		t.Fatal("OpenSafe should reject loopback hosts by default")
	}
	if _, err := OpenSafe("http://224.0.0.1/config.yaml"); err == nil {
		t.Fatal("OpenSafe should reject multicast hosts by default")
	}
	if _, err := OpenSafe("http://0.0.0.0/config.yaml"); err == nil {
		t.Fatal("OpenSafe should reject unspecified hosts by default")
	}
	if _, err := OpenSafe("ftp://example.com/config.yaml"); err == nil {
		t.Fatal("OpenSafe should reject non-HTTP schemes")
	}
}
