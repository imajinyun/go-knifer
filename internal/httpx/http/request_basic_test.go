package http

import (
	"io"
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

func TestRequestHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.Header.Get("X-Token")))
	}))
	defer srv.Close()

	body := Get(srv.URL).Header("X-Token", "abc").Execute().Body()
	if body != "abc" {
		t.Fatalf("body: %q", body)
	}
}

func TestRequestPatch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.Method))
	}))
	defer srv.Close()

	body := Patch(srv.URL).Execute().Body()
	if body != http.MethodPatch {
		t.Fatalf("method: %q", body)
	}
}

func TestRequestDelete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.Method))
	}))
	defer srv.Close()

	body := Delete(srv.URL).Execute().Body()
	if body != http.MethodDelete {
		t.Fatalf("method: %q", body)
	}
}
