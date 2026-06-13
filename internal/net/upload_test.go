package net

import (
	"bytes"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

func TestMultipartFileExts(t *testing.T) {
	req := multipartAvatarRequest(t, "a.txt")
	setting := NewUploadSetting()
	setting.FileExts = []string{".jpg"}
	setting.AllowFileExts = true
	if _, err := ParseMultipartForm(req, setting); err == nil {
		t.Fatal("ParseMultipartForm should reject extension outside allow list")
	}

	req = multipartAvatarRequest(t, "a.txt")
	setting.FileExts = []string{"txt"}
	setting.AllowFileExts = true
	if _, err := ParseMultipartForm(req, setting); err != nil {
		t.Fatalf("ParseMultipartForm should accept allowed extension: %v", err)
	}

	req = multipartAvatarRequest(t, "a.exe")
	setting.FileExts = []string{".exe"}
	setting.AllowFileExts = false
	if _, err := ParseMultipartForm(req, setting); err == nil {
		t.Fatal("ParseMultipartForm should reject extension in deny list")
	}
}

func TestSaveUploadedFileProviderOptions(t *testing.T) {
	req := multipartAvatarRequest(t, "a.txt")
	form, err := ParseMultipartForm(req, NewUploadSetting())
	if err != nil {
		t.Fatalf("ParseMultipartForm: %v", err)
	}
	file := form.GetFile("avatar")
	if file == nil {
		t.Fatal("uploaded file is nil")
	}

	var mkdirPath string
	var mkdirPerm fs.FileMode
	var openPath string
	var openFlag int
	var openPerm fs.FileMode
	var written bytes.Buffer
	err = SaveUploadedFile(file, "/virtual/upload/a.txt",
		WithUploadMkdirAll(func(path string, perm fs.FileMode) error {
			mkdirPath, mkdirPerm = path, perm
			return nil
		}),
		WithUploadOpenFile(func(path string, flag int, perm fs.FileMode) (io.WriteCloser, error) {
			openPath, openFlag, openPerm = path, flag, perm
			return nopWriteCloser{Writer: &written}, nil
		}),
		WithUploadDirPerm(0o700), WithUploadFilePerm(0o600),
	)
	if err != nil {
		t.Fatalf("SaveUploadedFile provider: %v", err)
	}
	if mkdirPath != "/virtual/upload" || mkdirPerm != 0o700 || openPath != "/virtual/upload/a.txt" || openPerm != 0o600 || openFlag&os.O_CREATE == 0 || written.String() != "hello" {
		t.Fatalf("providers mkdir=%q/%v open=%q flag=%#x perm=%v content=%q", mkdirPath, mkdirPerm, openPath, openFlag, openPerm, written.String())
	}
}

type nopWriteCloser struct{ io.Writer }

func (w nopWriteCloser) Close() error { return nil }

func multipartAvatarRequest(t *testing.T, filename string) *http.Request {
	t.Helper()
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	part, err := w.CreateFormFile("avatar", filename)
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := part.Write([]byte("hello")); err != nil {
		t.Fatalf("write form file: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, "/upload", body)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req
}
