package http

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestSafeRequestRejectsPrivateAndUnsafeRedirects(t *testing.T) {
	if err := GetSafe("file:///tmp/secret.txt").Execute().Err(); err == nil {
		t.Fatal("GetSafe should reject non-HTTP schemes")
	}
	if err := GetSafe("http://127.0.0.1/config.yaml").Execute().Err(); err == nil {
		t.Fatal("GetSafe should reject loopback hosts by default")
	}
	if err := GetSafe("http://224.0.0.1/config.yaml").Execute().Err(); err == nil {
		t.Fatal("GetSafe should reject multicast hosts by default")
	}
	if err := GetSafe("http://0.0.0.0/config.yaml").Execute().Err(); err == nil {
		t.Fatal("GetSafe should reject unspecified hosts by default")
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/redirect":
			http.Redirect(w, r, "http://127.0.0.1/private", http.StatusFound)
		default:
			_, _ = w.Write([]byte("ok"))
		}
	}))
	defer srv.Close()

	serverURL, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatalf("parse server url: %v", err)
	}
	resp := GetSafe(srv.URL,
		WithURLPolicy(URLPolicy{AllowedSchemes: []string{"http", "https"}, AllowedHosts: []string{serverURL.Hostname()}}),
	).Execute()
	if body := resp.Body(); body != "ok" || resp.Err() != nil {
		t.Fatalf("GetSafe allowed public policy host body=%q err=%v", body, resp.Err())
	}
	if err := GetSafe(srv.URL+"/redirect",
		WithURLPolicy(URLPolicy{AllowedSchemes: []string{"http", "https"}, AllowedHosts: []string{serverURL.Hostname()}}),
	).Execute().Err(); err == nil {
		t.Fatal("GetSafe should reject unsafe redirect targets")
	}
}

func TestSafeRequestAllowedHostsDoesNotBypassPrivateRejection(t *testing.T) {
	if err := GetSafe("http://127.0.0.1/config.yaml", WithAllowedHosts("127.0.0.1")).Execute().Err(); err == nil {
		t.Fatal("GetSafe should reject allowlisted private hosts when RejectPrivate is enabled")
	}

	req := GetSafe("http://example.com/config.yaml",
		WithAllowedHosts("example.com"),
		WithLookupIP(func(context.Context, string) ([]net.IP, error) {
			return []net.IP{net.ParseIP("127.0.0.1")}, nil
		}),
		WithTransport(roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("ok")), Header: make(http.Header)}, nil
		})),
	)
	if err := req.Execute().Err(); err == nil {
		t.Fatal("GetSafe should reject allowlisted hosts that resolve private during RoundTrip")
	}
}

func TestSafeRequestRevalidatesHostAtRoundTrip(t *testing.T) {
	req := GetSafe("http://example.com/config.yaml",
		WithLookupIP(func(context.Context, string) ([]net.IP, error) {
			return []net.IP{net.ParseIP("127.0.0.1")}, nil
		}),
		WithTransport(roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("ok")), Header: make(http.Header)}, nil
		})),
	)

	if err := req.Execute().Err(); err == nil {
		t.Fatal("GetSafe should reject a host that resolves private during RoundTrip")
	}
}
