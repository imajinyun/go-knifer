package vzip_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/imajinyun/go-knifer/vzip"
)

func TestFacadeZipFilesUsingOptions(t *testing.T) {
	tmp, src := newZipArchiveSource(t)
	archive := filepath.Join(tmp, "filtered.zip")
	if err := vzip.ZipFilesUsingOptions(archive, []string{src}, vzip.WithSourceDir(true), vzip.WithFileFilter(zipArchiveTextFilter)); err != nil {
		t.Fatalf("ZipFilesUsingOptions: %v", err)
	}
	data, err := vzip.GetBytes(archive, "src/keep.txt")
	if err != nil || string(data) != "keep" {
		t.Fatalf("GetBytes keep = %q, %v", data, err)
	}
	if _, err := vzip.GetBytes(archive, "src/skip.log"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("skip.log err = %v, want not exist", err)
	}
}
