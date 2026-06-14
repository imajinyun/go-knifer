package zip

import (
	"bytes"
	"compress/flate"
	"testing"
)

func TestGzipWithOptions(t *testing.T) {
	data := bytes.Repeat([]byte("abcdef"), 32)
	gz, err := GzipWithOptions(data, WithCompressionLevel(flate.BestSpeed))
	if err != nil {
		t.Fatalf("GzipWithOptions: %v", err)
	}
	out, err := UnGzip(gz)
	if err != nil || !bytes.Equal(out, data) {
		t.Fatalf("UnGzip after GzipWithOptions = %q, %v", out, err)
	}

	gz, err = GzipReaderWithOptions(bytes.NewReader(data), 0, WithCompressionLevel(flate.NoCompression))
	if err != nil {
		t.Fatalf("GzipReaderWithOptions: %v", err)
	}
	out, err = UnGzip(gz)
	if err != nil || !bytes.Equal(out, data) {
		t.Fatalf("UnGzip after GzipReaderWithOptions = %q, %v", out, err)
	}
	if _, err := GzipWithOptions(data, WithCompressionLevel(flate.HuffmanOnly+100)); err == nil {
		t.Fatal("GzipWithOptions should reject invalid compression level")
	}
}
