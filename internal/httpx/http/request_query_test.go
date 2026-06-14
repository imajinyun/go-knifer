package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Covers the utility toolkit-http HttpRequestTest.

func TestRequestGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/hello" || r.URL.Query().Get("name") != "world" {
			http.Error(w, "bad", 400)
			return
		}
		w.Header().Set("Content-Type", "text/plain;charset=UTF-8")
		_, _ = w.Write([]byte("hi world"))
	}))
	defer srv.Close()

	resp := Get(srv.URL+"/hello").Query("name", "world").Execute()
	if resp.Err() != nil {
		t.Fatalf("err: %v", resp.Err())
	}
	if !resp.IsOK() {
		t.Fatalf("status: %d", resp.Status())
	}
	if got := resp.Body(); got != "hi world" {
		t.Fatalf("body: %q", got)
	}
	if cs := resp.Charset(); strings.ToUpper(cs) != "UTF-8" {
		t.Fatalf("charset: %q", cs)
	}
}

func TestRequestQueryMap(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.URL.Query().Get("a") + "," + r.URL.Query().Get("b")))
	}))
	defer srv.Close()

	body := Get(srv.URL).QueryMap(map[string]any{"a": 1, "b": "x"}).Execute().Body()
	if body != "1,x" {
		t.Fatalf("body: %q", body)
	}
}
