package httpx_test

import (
	stdhttp "net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPContractDefaultTimeout(t *testing.T) {
	for _, backend := range httpContractBackends() {
		backend := backend
		t.Run(backend.name, func(t *testing.T) {
			backend.reset(t)

			if got := backend.snapshotTimeout(); got <= 0 {
				t.Fatalf("SnapshotGlobalConfig().Timeout = %v, want positive timeout", got)
			}
		})
	}
}

func TestHTTPContractRequestBasics(t *testing.T) {
	for _, backend := range httpContractBackends() {
		backend := backend
		t.Run(backend.name, func(t *testing.T) {
			backend.reset(t)

			srv := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
				_, _ = w.Write([]byte(r.Method + ":" + r.URL.Query().Get("q") + ":" + r.Header.Get("X-Contract") + ":" + r.Header.Get("User-Agent")))
			}))
			defer srv.Close()

			resp := backend.getBasic(srv.URL)
			if resp.err != nil || resp.status != stdhttp.StatusOK || resp.body != "GET:go:yes:contract-agent" {
				t.Fatalf("GET contract status=%d body=%q err=%v", resp.status, resp.body, resp.err)
			}
		})
	}
}

func TestHTTPContractPostJSONContentType(t *testing.T) {
	for _, backend := range httpContractBackends() {
		backend := backend
		t.Run(backend.name, func(t *testing.T) {
			backend.reset(t)

			srv := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
				_, _ = w.Write([]byte(r.Method + ":" + r.Header.Get("Content-Type")))
			}))
			defer srv.Close()

			resp := backend.postJSON(srv.URL)
			if resp.err != nil || resp.status != stdhttp.StatusOK || resp.body != "POST:application/json;charset=UTF-8" {
				t.Fatalf("POST JSON contract status=%d body=%q err=%v", resp.status, resp.body, resp.err)
			}
		})
	}
}
