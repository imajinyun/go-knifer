package zip

import (
	"bytes"
	"compress/flate"
	"os"
	"path/filepath"
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

func TestZlibFileWithOptionsPassesInputLimit(t *testing.T) {
	path := filepath.Join(t.TempDir(), "input.txt")
	if err := os.WriteFile(path, []byte("abcd"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := ZlibFileWithOptions(path, flate.DefaultCompression, WithMaxBytes(3)); err == nil {
		t.Fatal("ZlibFileWithOptions should pass WithMaxBytes to ZlibReaderWithOptions")
	}
}

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

func TestGzipAndZlib(t *testing.T) {
	data := []byte("hello compression")
	gz, err := Gzip(data)
	if err != nil {
		t.Fatalf("Gzip: %v", err)
	}
	out, err := UnGzip(gz)
	if err != nil || !bytes.Equal(out, data) {
		t.Fatalf("UnGzip: %q %v", out, err)
	}
	z, err := ZlibLevel(data, 6)
	if err != nil {
		t.Fatalf("ZlibLevel: %v", err)
	}
	out, err = UnZlib(z)
	if err != nil || !bytes.Equal(out, data) {
		t.Fatalf("UnZlib: %q %v", out, err)
	}
	if _, err := UnGzipWithOptions(gz, WithMaxBytes(3)); err == nil {
		t.Fatal("UnGzipWithOptions over limit error = nil")
	}
	if _, err := UnZlibWithOptions(z, WithMaxBytes(3)); err == nil {
		t.Fatal("UnZlibWithOptions over limit error = nil")
	}
}
