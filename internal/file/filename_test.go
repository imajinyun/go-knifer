package file

import "testing"

func TestMainNameAndExtension(t *testing.T) {
	if MainName("/x/y/foo.txt") != "foo" {
		t.Fatalf("MainName failed")
	}
	if Extension("/x/y/foo.txt") != "txt" {
		t.Fatalf("Extension failed")
	}
	if MainName("foo") != "foo" || Extension("foo") != "" {
		t.Fatalf("no-ext failed")
	}
}
