package http

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)

func TestRequestOptionCookieJar(t *testing.T) {
	oldJar := GetCookieJar()
	CloseCookie()
	defer SetCookieJar(oldJar)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/set" {
			http.SetCookie(w, &http.Cookie{Name: "sid", Value: "abc", Path: "/"})
			_, _ = w.Write([]byte("set"))
			return
		}
		c, err := r.Cookie("sid")
		if err != nil {
			_, _ = w.Write([]byte("missing"))
			return
		}
		_, _ = w.Write([]byte(c.Value))
	}))
	defer srv.Close()

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("cookiejar.New() error = %v", err)
	}
	if resp := Get(srv.URL+"/set", WithCookieJar(jar)).Execute(); resp.Err() != nil {
		t.Fatalf("set cookie request error = %v", resp.Err())
	}
	if got := Get(srv.URL+"/get", WithCookieJar(jar)).Execute().Body(); got != "abc" {
		t.Fatalf("cookie jar body = %q, want abc", got)
	}
}
