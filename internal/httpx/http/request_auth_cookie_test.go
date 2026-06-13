package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestBasicAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.Header.Get("Authorization")))
	}))
	defer srv.Close()

	body := Get(srv.URL).BasicAuth("aladdin", "opensesame").Execute().Body()
	if body != "Basic YWxhZGRpbjpvcGVuc2VzYW1l" {
		t.Fatalf("body: %q", body)
	}
}

func TestRequestBearerAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.Header.Get("Authorization")))
	}))
	defer srv.Close()

	body := Get(srv.URL).BearerAuth("xyz.token").Execute().Body()
	if body != "Bearer xyz.token" {
		t.Fatalf("body: %q", body)
	}
}

func TestRequestCookie(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("k")
		if c == nil {
			_, _ = w.Write([]byte("no"))
			return
		}
		_, _ = w.Write([]byte(c.Value))
	}))
	defer srv.Close()

	body := Get(srv.URL).Cookie(&http.Cookie{Name: "k", Value: "v"}).Execute().Body()
	if body != "v" {
		t.Fatalf("body: %q", body)
	}
}
