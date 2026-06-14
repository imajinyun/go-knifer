package zip

import (
	"bytes"
	"testing"
)

func TestDecompressHelpersEnforceConfiguredLimit(t *testing.T) {
	data := bytes.Repeat([]byte("x"), 16)
	gz, err := Gzip(data)
	if err != nil {
		t.Fatalf("Gzip: %v", err)
	}
	zl, err := Zlib(data)
	if err != nil {
		t.Fatalf("Zlib: %v", err)
	}

	if got := applyDecompressOptions(nil).maxBytes; got != DefaultUnzipMaxBytes {
		t.Fatalf("default decompression max bytes = %d, want %d", got, DefaultUnzipMaxBytes)
	}
	if _, err := UnGzipWithOptions(gz, WithMaxBytes(8)); err == nil {
		t.Fatal("UnGzip should enforce configured decompression limit")
	}
	if _, err := UnZlibWithOptions(zl, WithMaxBytes(8)); err == nil {
		t.Fatal("UnZlib should enforce configured decompression limit")
	}
	if out, err := UnGzipWithOptions(gz, WithMaxBytes(0)); err != nil || !bytes.Equal(out, data) {
		t.Fatalf("UnGzip explicit unlimited = %q, %v", out, err)
	}
	if out, err := UnZlibWithOptions(zl, WithMaxBytes(0)); err != nil || !bytes.Equal(out, data) {
		t.Fatalf("UnZlib explicit unlimited = %q, %v", out, err)
	}
}
