package file

import (
	"path/filepath"
	"testing"
)

func TestMkdirTouchDelTouch(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "x", "y")
	if err := Mkdir(sub); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	f := filepath.Join(sub, "a.txt")
	if err := Touch(f); err != nil {
		t.Fatalf("touch: %v", err)
	}
	if !IsFile(f) {
		t.Fatalf("touch did not create file")
	}
}
