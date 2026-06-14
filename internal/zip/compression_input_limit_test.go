package zip

import (
	"bytes"
	"compress/flate"
	"testing"
)

func TestCompressionHelpersEnforceInputLimit(t *testing.T) {
	data := []byte("abcd")
	if _, err := GzipReaderWithOptions(bytes.NewReader(data), len(data), WithMaxBytes(3)); err == nil {
		t.Fatal("GzipReaderWithOptions should reject compression input over max bytes")
	}
	if _, err := ZlibReaderWithOptions(bytes.NewReader(data), flate.DefaultCompression, len(data), WithMaxBytes(3)); err == nil {
		t.Fatal("ZlibReaderWithOptions should reject compression input over max bytes")
	}
	if out, err := GzipReaderWithOptions(bytes.NewReader(data), len(data), WithMaxBytes(4)); err != nil || len(out) == 0 {
		t.Fatalf("GzipReaderWithOptions exact limit = %d bytes, %v", len(out), err)
	}
	if out, err := ZlibLevelWithOptions(data, flate.DefaultCompression, WithMaxBytes(4)); err != nil || len(out) == 0 {
		t.Fatalf("ZlibLevelWithOptions exact limit = %d bytes, %v", len(out), err)
	}
}
