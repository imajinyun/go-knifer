package resty

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAdditionalClientFactoriesSafeWrappersAndMethods(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Method", r.Method)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Add("Set-Cookie", "sid=abc; Path=/")
		if r.Method != http.MethodHead {
			_, _ = w.Write([]byte(r.Method + ":" + r.Header.Get("X-Client")))
		}
	}))
	defer srv.Close()

	cfg := SnapshotGlobalConfig()
	cfg.Headers["X-Client"] = []string{"cfg"}
	client := NewClientWithConfig(cfg, WithHeader("X-Client", "opt"))
	if got := client.Get(srv.URL).Execute().Body(); got != "GET:opt" {
		t.Fatalf("client.Get body = %q", got)
	}
	if got := client.Post(srv.URL).Execute().Body(); got != "POST:opt" {
		t.Fatalf("client.Post body = %q", got)
	}
	if got := NewIsolatedClient(WithClientGlobalConfig(cfg), WithClientRequestOptions(WithHeader("X-Client", "isolated"))).NewRequest(MethodPut, srv.URL).Execute().Body(); got != "PUT:isolated" {
		t.Fatalf("NewIsolatedClient body = %q", got)
	}
	if got := (*Client)(nil).NewRequest(MethodDelete, srv.URL).Execute().Header("X-Method"); got != string(MethodDelete) {
		t.Fatalf("nil client NewRequest method = %q", got)
	}

	allowLocal := WithURLPolicy(URLPolicy{AllowedSchemes: []string{"http", "https"}, RejectPrivate: false})
	requests := []*HTTPRequest{
		PostSafe(srv.URL, allowLocal),
		Put(srv.URL),
		PutSafe(srv.URL, allowLocal),
		Delete(srv.URL),
		DeleteSafe(srv.URL, allowLocal),
		Patch(srv.URL),
		PatchSafe(srv.URL, allowLocal),
		Head(srv.URL),
		HeadSafe(srv.URL, allowLocal),
		Options(srv.URL),
		OptionsSafe(srv.URL, allowLocal),
		NewSafeRequest(MethodTrace, srv.URL, allowLocal),
		client.NewSafeRequest(MethodOptions, srv.URL, allowLocal),
	}
	for _, req := range requests {
		resp := req.Execute()
		if resp.Err() != nil {
			t.Fatalf("safe wrapper Execute: %v", resp.Err())
		}
		if resp.Status() == 0 {
			t.Fatal("safe wrapper status = 0")
		}
	}

	resp := Get(srv.URL).Cookie("k", "v").Execute()
	if resp.Err() != nil {
		t.Fatalf("cookie Execute: %v", resp.Err())
	}
	if got := resp.Headers()["X-Method"]; len(got) != 1 || got[0] != http.MethodGet {
		t.Fatalf("Headers()[X-Method] = %v", got)
	}
	if cookies := resp.Cookies(); len(cookies) != 1 || cookies[0].Name != "sid" {
		t.Fatalf("Cookies = %#v", cookies)
	}
	if resp.ContentType() == "" || resp.ContentLength() == 0 || resp.RestyRaw() == nil {
		t.Fatalf("response metadata type=%q length=%d raw=%v", resp.ContentType(), resp.ContentLength(), resp.RestyRaw())
	}
	var out bytes.Buffer
	if n, err := resp.WriteTo(&out); err != nil || n != int64(out.Len()) || !strings.Contains(out.String(), "GET") {
		t.Fatalf("WriteTo n=%d body=%q err=%v", n, out.String(), err)
	}
	if err := resp.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}
