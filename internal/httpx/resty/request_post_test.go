package resty

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostForm(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm() error = %v", err)
		}
		if got := r.Form.Get("name"); got != "resty" {
			t.Fatalf("form name = %q, want resty", got)
		}
		_, _ = w.Write([]byte("posted"))
	}))
	defer srv.Close()

	got, err := PostFormE(srv.URL, map[string]any{"name": "resty"})
	if err != nil {
		t.Fatalf("PostFormE() error = %v", err)
	}
	if got != "posted" {
		t.Fatalf("PostFormE() = %q, want posted", got)
	}
}

func TestPostJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Content-Type"), string(ContentTypeJSON)) {
			t.Fatalf("Content-Type = %q, want application/json", r.Header.Get("Content-Type"))
		}
		_, _ = w.Write([]byte("json"))
	}))
	defer srv.Close()

	got, err := PostJSONE(srv.URL, `{"ok":true}`)
	if err != nil {
		t.Fatalf("PostJSONE() error = %v", err)
	}
	if got != "json" {
		t.Fatalf("PostJSONE() = %q, want json", got)
	}
}
