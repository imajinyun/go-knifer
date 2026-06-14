package httpx_test

import (
	stdhttp "net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPContractClientGlobalSnapshot(t *testing.T) {
	for _, backend := range httpContractBackends() {
		backend := backend
		t.Run(backend.name, func(t *testing.T) {
			backend.reset(t)

			srv := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
				_, _ = w.Write([]byte(r.Header.Get("X-Contract-Global")))
			}))
			defer srv.Close()

			snapshot := backend.clientSnapshot(t, srv.URL)
			if snapshot.err != nil || snapshot.body != "snapshot" {
				t.Fatalf("client snapshot body=%q err=%v", snapshot.body, snapshot.err)
			}
		})
	}
}

func TestHTTPContractIsolatedClient(t *testing.T) {
	for _, backend := range httpContractBackends() {
		backend := backend
		t.Run(backend.name, func(t *testing.T) {
			backend.reset(t)

			srv := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
				_, _ = w.Write([]byte(r.Header.Get("X-Contract-Global")))
			}))
			defer srv.Close()

			isolated := backend.isolatedClient(t, srv.URL)
			if isolated.err != nil || isolated.body != "" {
				t.Fatalf("isolated client body=%q err=%v, want no global headers", isolated.body, isolated.err)
			}
		})
	}
}
