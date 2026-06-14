package resty

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSafeRequestUtilityWrappers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_, _ = w.Write([]byte(r.Method + ":" + r.URL.Query().Get("q") + ":" + r.Header.Get("Authorization") + ":" + string(body)))
	}))
	defer srv.Close()

	if got, err := GetWithTimeoutE(srv.URL, time.Second); err != nil || !strings.HasPrefix(got, "GET:") {
		t.Fatalf("GetWithTimeoutE = %q, %v", got, err)
	}
	if got, err := GetWithTimeoutEWithOptions(srv.URL, time.Second, WithHeader("X-T", "v")); err != nil || !strings.HasPrefix(got, "GET:") {
		t.Fatalf("GetWithTimeoutEWithOptions = %q, %v", got, err)
	}
	policy := WithURLPolicy(URLPolicy{AllowedSchemes: []string{"http", "https"}, RejectPrivate: false})
	if got, err := PostStringSafeE(srv.URL, "safe", policy); err != nil || !strings.Contains(got, "POST:::safe") {
		t.Fatalf("PostStringSafeE = %q, %v", got, err)
	}
	if got, err := PostFormSafeE(srv.URL, map[string]any{"a": "b"}, policy); err != nil || !strings.HasPrefix(got, "POST:") {
		t.Fatalf("PostFormSafeE = %q, %v", got, err)
	}
	if got, err := PostJSONSafeE(srv.URL, `{"ok":true}`, policy); err != nil || !strings.Contains(got, `{"ok":true}`) {
		t.Fatalf("PostJSONSafeE = %q, %v", got, err)
	}
}
