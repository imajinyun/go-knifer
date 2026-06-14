package conf

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestLoadRemoteSafeRejectsPrivateHosts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("app:\n  name: remote"))
	}))
	defer server.Close()

	if _, err := LoadRemoteSafe(server.URL + "/app.yaml"); err == nil {
		t.Fatal("LoadRemoteSafe should reject private hosts by default")
	}
	if _, err := LoadRemoteSafe("http://224.0.0.1/app.yaml"); err == nil {
		t.Fatal("LoadRemoteSafe should reject multicast hosts by default")
	}
	if _, err := LoadRemoteSafe("http://0.0.0.0/app.yaml"); err == nil {
		t.Fatal("LoadRemoteSafe should reject unspecified hosts by default")
	}
	remoteURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := LoadRemoteSafeWithOptions(server.URL+"/app.yaml", LoadOptions{RemoteAllowedHosts: []string{remoteURL.Hostname()}}); err == nil {
		t.Fatal("LoadRemoteSafeWithOptions should reject allowlisted private hosts")
	}
}

func TestLoadRemoteSafeRejectsUnsafeRedirectTarget(t *testing.T) {
	redirect := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://127.0.0.1/private.yaml", http.StatusFound)
	}))
	defer redirect.Close()
	redirectURL, err := url.Parse(redirect.URL)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := LoadRemoteSafeWithOptions(redirect.URL+"/app.yaml", LoadOptions{RemoteAllowedHosts: []string{redirectURL.Hostname()}}); err == nil {
		t.Fatal("LoadRemoteSafeWithOptions should reject unsafe redirect target")
	}
}
