package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponseReadOptions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("abcdef"))
	}))
	defer srv.Close()

	limited := Get(srv.URL, WithMaxResponseBytes(3)).Execute()
	if got := limited.Bytes(); len(got) != 0 || limited.Err() == nil {
		t.Fatalf("limited Bytes() = %q err=%v, want max bytes error", string(got), limited.Err())
	}

	readAllCalled := false
	resp := Get(srv.URL, WithResponseReadAllFunc(func(r io.Reader) ([]byte, error) {
		readAllCalled = true
		return []byte("provided"), nil
	})).Execute()
	if got := resp.Body(); got != "provided" || !readAllCalled || resp.Err() != nil {
		t.Fatalf("custom readAll body=%q called=%v err=%v", got, readAllCalled, resp.Err())
	}
}

func TestResponseReadLimitFollowsGlobalConfigSnapshot(t *testing.T) {
	previous := SnapshotGlobalConfig()
	defer ConfigureGlobalConfig(previous)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("abcdef"))
	}))
	defer srv.Close()

	SetGlobalMaxResponseBytes(3)
	resp := Get(srv.URL).Execute()
	SetGlobalMaxResponseBytes(0)
	if got := resp.Bytes(); len(got) != 0 || resp.Err() == nil {
		t.Fatalf("global limited Bytes() = %q err=%v, want max bytes error", string(got), resp.Err())
	}

	unlimited := Get(srv.URL, WithMaxResponseBytes(0)).Execute()
	if got := unlimited.Body(); got != "abcdef" || unlimited.Err() != nil {
		t.Fatalf("unlimited override Body() = %q err=%v", got, unlimited.Err())
	}
}

func TestResponseIgnoreEOFFollowsRequestSnapshot(t *testing.T) {
	old := IsIgnoreEOFError()
	defer SetIgnoreEOFError(old)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("abcdef"))
	}))
	defer srv.Close()

	readUnexpectedEOF := func(io.Reader) ([]byte, error) {
		return []byte("partial"), io.ErrUnexpectedEOF
	}

	ignoreResp := NewRequestWithConfig(MethodGet, srv.URL, GlobalConfig{FollowRedirects: true, MaxRedirects: 10, MaxResponseBytes: defaultGlobalMaxResponseBytes, IgnoreEOFError: true}, WithResponseReadAllFunc(readUnexpectedEOF)).Execute()
	SetIgnoreEOFError(false)
	if got := ignoreResp.Body(); got != "partial" || ignoreResp.Err() != nil {
		t.Fatalf("ignore snapshot body=%q err=%v, want partial without error", got, ignoreResp.Err())
	}

	strictResp := NewRequestWithConfig(MethodGet, srv.URL, GlobalConfig{FollowRedirects: true, MaxRedirects: 10, MaxResponseBytes: defaultGlobalMaxResponseBytes, IgnoreEOFError: false}, WithResponseReadAllFunc(readUnexpectedEOF)).Execute()
	SetIgnoreEOFError(true)
	if got := strictResp.Body(); got != "" || strictResp.Err() == nil {
		t.Fatalf("strict snapshot body=%q err=%v, want read error", got, strictResp.Err())
	}
}
