package resty

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStringHelpersReturnErrorsExplicitly(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.Method + ":" + r.URL.Query().Get("k")))
	}))
	defer srv.Close()

	body, err := GetWithParamsE(srv.URL, map[string]any{"k": "v"})
	if err != nil || body != "GET:v" {
		t.Fatalf("GetWithParamsE = %q, %v", body, err)
	}

	if body, err = PostStringE(srv.URL, "payload"); err != nil || body != "POST:" {
		t.Fatalf("PostStringE = %q, %v", body, err)
	}

	if _, err = GetStringE("http://[::1"); err == nil {
		t.Fatal("GetStringE invalid URL error = nil")
	}
	if _, err = DownloadBytesE("http://[::1"); err == nil {
		t.Fatal("DownloadBytesE invalid URL error = nil")
	}
	if _, err = GetStringSafeE(srv.URL); err == nil {
		t.Fatal("GetStringSafeE local URL error = nil, want private address rejection")
	}
}
