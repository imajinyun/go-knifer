package httpx_test

import (
	stdhttp "net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHTTPContractInvalidURL(t *testing.T) {
	for _, backend := range httpContractBackends() {
		backend := backend
		t.Run(backend.name, func(t *testing.T) {
			backend.reset(t)

			if err := backend.invalidURL(); err == nil {
				t.Fatal("invalid URL error = nil")
			}
		})
	}
}

func TestHTTPContractSafeURLRejectsUnsafeTargets(t *testing.T) {
	for _, backend := range httpContractBackends() {
		backend := backend
		t.Run(backend.name, func(t *testing.T) {
			backend.reset(t)

			for _, rawURL := range []string{"file:///tmp/secret.txt", "http://127.0.0.1/config.yaml"} {
				if err := backend.safeURL(rawURL); err == nil {
					t.Fatalf("safe request to %q error = nil", rawURL)
				}
			}
		})
	}
}

func TestHTTPContractSafeURLAllowedHost(t *testing.T) {
	for _, backend := range httpContractBackends() {
		backend := backend
		t.Run(backend.name, func(t *testing.T) {
			backend.reset(t)

			srv := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, _ *stdhttp.Request) {
				_, _ = w.Write([]byte("safe"))
			}))
			defer srv.Close()

			parsed, err := url.Parse(srv.URL)
			if err != nil {
				t.Fatalf("parse server URL: %v", err)
			}
			resp := backend.safeAllowed(srv.URL, parsed.Hostname())
			if resp.err != nil || resp.status != stdhttp.StatusOK || resp.body != "safe" {
				t.Fatalf("safe allowed status=%d body=%q err=%v", resp.status, resp.body, resp.err)
			}
		})
	}
}
