package http

import (
	"crypto/tls"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
)

func TestAdditionalRequestOptionsAndAccessors(t *testing.T) {
	requestFactoryCalled := false
	readAllCalled := false
	transport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode:    http.StatusOK,
			ContentLength: int64(len(req.Method + ":" + req.Header.Get("X-A"))),
			Header:        http.Header{"Content-Type": []string{"text/plain"}},
			Body:          io.NopCloser(strings.NewReader(req.Method + ":" + req.Header.Get("X-A"))),
			Request:       req,
		}, nil
	})
	cfg := SnapshotGlobalConfig()
	cfg.Headers.Set("X-A", "cfg")
	req := NewIsolatedRequest(MethodPost, "https://example.com/upload",
		WithGlobalConfig(cfg),
		WithHeaders(map[string]string{"X-A": "one", "X-B": "two"}),
		WithClient(&http.Client{Transport: transport}),
		WithTransport(transport),
		WithTLSConfig(&tls.Config{MinVersion: tls.VersionTLS12}),
		WithResponseReadAllFunc(func(r io.Reader) ([]byte, error) {
			readAllCalled = true
			return io.ReadAll(r)
		}),
		WithRequestFactory(func(method, rawURL string, body io.Reader) (*http.Request, error) {
			requestFactoryCalled = true
			return http.NewRequest(method, rawURL, body)
		}),
		WithMultipartWriterFactory(func(w io.Writer) MultipartWriter {
			return multipart.NewWriter(w)
		}),
	)
	req.Method(MethodPatch).URL("https://example.com/changed").AddHeader("X-A", "extra").Headers(map[string]string{"X-B": "two"}).CookieString("raw=cookie")
	req.Client(&http.Client{Transport: transport}).URLPolicy(URLPolicy{AllowedSchemes: []string{"https"}, RejectPrivate: false})
	resp := req.FormFileReader("file", "a.txt", strings.NewReader("file-data")).Execute()
	if resp.Err() != nil {
		t.Fatalf("multipart Execute: %v", resp.Err())
	}
	if got := resp.Body(); got != "PATCH:one" {
		t.Fatalf("response body = %q", got)
	}
	if !requestFactoryCalled || !readAllCalled {
		t.Fatalf("providers called request=%v readAll=%v", requestFactoryCalled, readAllCalled)
	}
	if got := NewRequestWithConfig(MethodGet, "https://example.com", cfg, WithTransport(transport)).Execute().Body(); got != "GET:cfg" {
		t.Fatalf("NewRequestWithConfig body = %q", got)
	}
}
