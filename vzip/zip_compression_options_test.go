package vzip_test

import (
	"bytes"
	"compress/flate"
	"os"
	"path/filepath"
	"testing"

	"github.com/imajinyun/go-knifer/vzip"
)

func TestFacadeCompressionOptions(t *testing.T) {
	tmp := t.TempDir()
	payload := []byte("compression option payload")
	source := filepath.Join(tmp, "payload.txt")
	if err := os.WriteFile(source, payload, 0o644); err != nil {
		t.Fatal(err)
	}
	gz, err := vzip.GzipFileWithOptions(source, vzip.WithMaxBytes(int64(len(payload))))
	if err != nil {
		t.Fatalf("GzipFileWithOptions: %v", err)
	}
	if out, err := vzip.UnGzipReaderWithOptions(bytes.NewReader(gz), len(payload), vzip.WithMaxBytes(int64(len(payload)))); err != nil || !bytes.Equal(out, payload) {
		t.Fatalf("UnGzipReaderWithOptions = %q, %v", out, err)
	}
	if _, err := vzip.GzipFileWithOptions(source, vzip.WithMaxBytes(1)); err == nil {
		t.Fatal("GzipFileWithOptions max bytes error = nil")
	}

	zlibBytes, err := vzip.ZlibFileWithOptions(source, flate.BestSpeed, vzip.WithMaxBytes(int64(len(payload))))
	if err != nil {
		t.Fatalf("ZlibFileWithOptions: %v", err)
	}
	if out, err := vzip.UnZlibReaderWithOptions(bytes.NewReader(zlibBytes), len(payload), vzip.WithMaxBytes(int64(len(payload)))); err != nil || !bytes.Equal(out, payload) {
		t.Fatalf("UnZlibReaderWithOptions = %q, %v", out, err)
	}
	if got, err := vzip.ZlibLevelWithOptions(payload, flate.NoCompression, vzip.WithMaxBytes(int64(len(payload)))); err != nil || len(got) == 0 {
		t.Fatalf("ZlibLevelWithOptions len=%d err=%v", len(got), err)
	}
	if got, err := vzip.ZlibReaderWithOptions(bytes.NewReader(payload), flate.BestSpeed, len(payload), vzip.WithMaxBytes(int64(len(payload)))); err != nil || len(got) == 0 {
		t.Fatalf("ZlibReaderWithOptions len=%d err=%v", len(got), err)
	}
}
