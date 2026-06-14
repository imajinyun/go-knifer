package zip

import (
	"bytes"
	"testing"
)

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
