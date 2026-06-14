package httpx_test

import (
	stdhttp "net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPContractRedirectControls(t *testing.T) {
	for _, backend := range httpContractBackends() {
		backend := backend
		t.Run(backend.name, func(t *testing.T) {
			backend.reset(t)

			srv := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
				if r.URL.Path == "/start" {
					stdhttp.Redirect(w, r, "/end", stdhttp.StatusFound)
					return
				}
				_, _ = w.Write([]byte("end"))
			}))
			defer srv.Close()

			noFollow := backend.getNoFollow(srv.URL + "/start")
			if noFollow.err != nil || noFollow.status != stdhttp.StatusFound {
				t.Fatalf("no-follow status=%d body=%q err=%v", noFollow.status, noFollow.body, noFollow.err)
			}

			follow := backend.getFollow(srv.URL + "/start")
			if follow.err != nil || follow.status != stdhttp.StatusOK || follow.body != "end" {
				t.Fatalf("follow status=%d body=%q err=%v", follow.status, follow.body, follow.err)
			}
		})
	}
}

func TestHTTPContractMaxResponseBytes(t *testing.T) {
	for _, backend := range httpContractBackends() {
		backend := backend
		t.Run(backend.name, func(t *testing.T) {
			backend.reset(t)

			srv := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, _ *stdhttp.Request) {
				_, _ = w.Write([]byte("abcdef"))
			}))
			defer srv.Close()

			resp := backend.getWithMaxBytes(srv.URL, 3)
			if resp.err == nil || resp.body != "" {
				t.Fatalf("limited body=%q err=%v, want max bytes error", resp.body, resp.err)
			}
		})
	}
}
