package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestMultipartUpload(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		f, fh, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		defer func() {
			if err := f.Close(); err != nil {
				t.Errorf("close multipart file: %v", err)
			}
		}()
		data, err := io.ReadAll(f)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		_, _ = w.Write([]byte(fh.Filename + ":" + string(data) + ":" + r.FormValue("k")))
	}))
	defer srv.Close()

	body := Post(srv.URL).
		Form(map[string]any{"k": "v"}).
		FormFile("file", "hello.txt", []byte("hi")).
		Execute().Body()
	if body != "hello.txt:hi:v" {
		t.Fatalf("body: %q", body)
	}
}

func TestRequestFormFileReader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		f, fh, err := r.FormFile("f")
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		defer func() {
			if err := f.Close(); err != nil {
				t.Errorf("close multipart file: %v", err)
			}
		}()
		data, err := io.ReadAll(f)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		_, _ = w.Write([]byte(fh.Filename + ":" + string(data)))
	}))
	defer srv.Close()

	body := Post(srv.URL).
		FormFileReader("f", "in.txt", strings.NewReader("hello reader")).
		Execute().Body()
	if body != "in.txt:hello reader" {
		t.Fatalf("body: %q", body)
	}
}

func TestRequestReaderBackedBodyIsSingleUse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		_, _ = w.Write(b)
	}))
	defer srv.Close()

	req := Post(srv.URL).BodyReader(strings.NewReader("hello"))
	if got := req.Execute().Body(); got != "hello" {
		t.Fatalf("first body = %q", got)
	}
	resp := req.Execute()
	if resp.Err() == nil {
		t.Fatal("second Execute() should reject reader-backed body reuse")
	}
}
