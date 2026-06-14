package vcaptcha_test

import (
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/imajinyun/go-knifer/vcaptcha"
)

func TestFacadeCaptchaWriteToFileOptions(t *testing.T) {
	c := vcaptcha.NewLineCaptchaWithOptions(100, 40, vcaptcha.WithGenerator(fixedGenerator{code: "ABCD"}))
	c.CreateCode()
	path := filepath.Join(t.TempDir(), "nested", "captcha.png")
	if err := c.WriteToFileWithOptions(path, vcaptcha.WithFilePerm(0o600), vcaptcha.WithDirPerm(0o700)); err != nil {
		t.Fatalf("WriteToFileWithOptions: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat captcha file: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("captcha file perm = %o, want 600", got)
	}
	if err := c.WriteToFileWithOptions(path, vcaptcha.WithOverwrite(false)); err == nil {
		t.Fatal("WriteToFileWithOptions should reject overwrite=false for existing file")
	}
}

func TestFacadeCaptchaWriteProviderOptions(t *testing.T) {
	c := vcaptcha.NewLineCaptchaWithOptions(100, 40, vcaptcha.WithGenerator(fixedGenerator{code: "ABCD"}))
	c.CreateCode()

	var mkdirPath string
	var mkdirPerm fs.FileMode
	var openPath string
	var openFlag int
	var openPerm fs.FileMode
	var written bytes.Buffer
	if err := c.WriteToFileWithOptions("/virtual/captcha.png",
		vcaptcha.WithMkdirAll(func(path string, perm fs.FileMode) error {
			mkdirPath, mkdirPerm = path, perm
			return nil
		}),
		vcaptcha.WithOpenFile(func(path string, flag int, perm fs.FileMode) (io.WriteCloser, error) {
			openPath, openFlag, openPerm = path, flag, perm
			return nopWriteCloser{Writer: &written}, nil
		}),
		vcaptcha.WithDirPerm(0o700), vcaptcha.WithFilePerm(0o600),
	); err != nil {
		t.Fatalf("WriteToFileWithOptions provider: %v", err)
	}
	if mkdirPath != "/virtual" || mkdirPerm != 0o700 || openPath != "/virtual/captcha.png" || openPerm != 0o600 || openFlag&os.O_CREATE == 0 || written.Len() == 0 {
		t.Fatalf("providers mkdir=%q/%v open=%q flag=%#x perm=%v bytes=%d", mkdirPath, mkdirPerm, openPath, openFlag, openPerm, written.Len())
	}
}
