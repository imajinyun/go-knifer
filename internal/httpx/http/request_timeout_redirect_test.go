package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRequestTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	resp := Get(srv.URL).Timeout(50 * time.Millisecond).Execute()
	if resp.Err() == nil {
		t.Fatal("expected timeout error")
	}
}

func TestRequestNoFollowRedirects(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/start" {
			http.Redirect(w, r, "/end", http.StatusFound)
			return
		}
		_, _ = w.Write([]byte("end"))
	}))
	defer srv.Close()

	resp := Get(srv.URL + "/start").FollowRedirects(false).Execute()
	if resp.Status() != http.StatusFound {
		t.Fatalf("expected 302, got %d", resp.Status())
	}

	body := Get(srv.URL + "/start").FollowRedirects(true).Execute().Body()
	if body != "end" {
		t.Fatalf("redirected body: %q", body)
	}
}
