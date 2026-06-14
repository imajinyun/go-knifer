package zip

import (
	"compress/flate"
	"os"
	"path/filepath"
	"testing"
)

func TestZlibFileWithOptionsPassesInputLimit(t *testing.T) {
	path := filepath.Join(t.TempDir(), "input.txt")
	if err := os.WriteFile(path, []byte("abcd"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := ZlibFileWithOptions(path, flate.DefaultCompression, WithMaxBytes(3)); err == nil {
		t.Fatal("ZlibFileWithOptions should pass WithMaxBytes to ZlibReaderWithOptions")
	}
}
