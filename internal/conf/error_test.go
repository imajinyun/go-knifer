package conf

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	knifer "github.com/imajinyun/go-knifer"
)

func TestConfErrorContract(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "missing.setting"))
	assertConfCode(t, err, knifer.ErrCodeNotFound)
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Load missing file should preserve os.ErrNotExist: %v", err)
	}

	_, err = Parse("invalid-line")
	assertConfCode(t, err, knifer.ErrCodeInvalidInput)

	_, err = Parse("=empty")
	assertConfCode(t, err, knifer.ErrCodeInvalidInput)

	_, err = ParseYAML("invalid-yaml-line")
	assertConfCode(t, err, knifer.ErrCodeInvalidInput)
}
