package resty

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestOptionContentTypeAndCharset(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.Header.Get("Content-Type")))
	}))
	defer srv.Close()

	got := Post(srv.URL, WithCharset("GBK"), WithContentType("text/custom")).BodyString("hello").Execute().Body()
	if got != "text/custom" {
		t.Fatalf("Content-Type = %q, want text/custom", got)
	}

	got = Post(srv.URL, WithCharset("GBK")).BodyJSON(`{"ok":true}`).Execute().Body()
	if got != "application/json;charset=GBK" {
		t.Fatalf("JSON Content-Type = %q", got)
	}
}

func TestRequestOptionTLSConfig(t *testing.T) {
	c := Get("https://example.com", WithTLSConfig(&tls.Config{ServerName: "example.com"})).buildClient()
	if c == nil {
		t.Fatal("client is nil")
	}
}
