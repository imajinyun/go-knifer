package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestPostForm(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		_, _ = w.Write([]byte(r.PostForm.Get("k")))
	}))
	defer srv.Close()

	body := Post(srv.URL).Form(map[string]any{"k": "v"}).Execute().Body()
	if body != "v" {
		t.Fatalf("body: %q", body)
	}
}

func TestRequestPostJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			http.Error(w, "bad ct", 400)
			return
		}
		b, _ := io.ReadAll(r.Body)
		_, _ = w.Write(b)
	}))
	defer srv.Close()

	resp, err := PostJSONE(srv.URL, `{"a":1}`)
	if err != nil {
		t.Fatalf("PostJSONE() error = %v", err)
	}
	if resp != `{"a":1}` {
		t.Fatalf("body: %q", resp)
	}
}

func TestRequestBodyStringAutoContentType(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.Header.Get("Content-Type")))
	}))
	defer srv.Close()

	ct := Post(srv.URL).BodyString(`{"x":1}`).Execute().Body()
	if !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("expected json content-type detected, got %q", ct)
	}

	ct2 := Post(srv.URL).BodyString(`<x/>`).Execute().Body()
	if !strings.HasPrefix(ct2, "application/xml") {
		t.Fatalf("expected xml content-type detected, got %q", ct2)
	}
}
