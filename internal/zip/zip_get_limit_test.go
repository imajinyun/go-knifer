package zip

import (
	"bytes"
	"path/filepath"
	"testing"
)

func TestGetBytesEnforcesConfiguredLimit(t *testing.T) {
	tmp := t.TempDir()
	archive := filepath.Join(tmp, "limit.zip")
	data := bytes.Repeat([]byte("x"), 16)
	if err := ZipEntries(archive, EntryData{Name: "a.txt", Data: data}); err != nil {
		t.Fatalf("ZipEntries: %v", err)
	}

	if got := applyDecompressOptions(nil).maxBytes; got != DefaultUnzipMaxBytes {
		t.Fatalf("default entry read max bytes = %d, want %d", got, DefaultUnzipMaxBytes)
	}
	if _, err := GetBytesWithOptions(archive, "a.txt", WithMaxBytes(8)); err == nil {
		t.Fatal("GetBytes should enforce configured entry read limit")
	}
	if out, err := GetBytesWithOptions(archive, "a.txt", WithMaxBytes(0)); err != nil || !bytes.Equal(out, data) {
		t.Fatalf("GetBytes explicit unlimited = %q, %v", out, err)
	}
}
