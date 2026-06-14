package url

import (
	"io"
	"net/http"
	"net/http/httptest"
	neturl "net/url"
	"testing"
)

func TestOpenSafeAllowsExplicitHost(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("safe"))
	}))
	defer target.Close()
	targetURL, err := neturl.Parse(target.URL)
	if err != nil {
		t.Fatal(err)
	}

	r, err := OpenSafeWithOptions(target.URL, WithAllowedHosts(targetURL.Hostname()), WithRejectPrivateHosts(false))
	if err != nil {
		t.Fatalf("OpenSafeWithOptions allow host: %v", err)
	}
	data, err := io.ReadAll(r)
	_ = r.Close()
	if err != nil || string(data) != "safe" {
		t.Fatalf("safe body = %q, %v", data, err)
	}
}

func TestOpenSafeAllowedHostsDoesNotBypassPrivateRejection(t *testing.T) {
	privateHosts := []string{"127.0.0.1", "localhost"}
	for _, host := range privateHosts {
		t.Run(host, func(t *testing.T) {
			if _, err := OpenSafeWithOptions("http://"+host+"/config.yaml", WithAllowedHosts(host)); err == nil {
				t.Fatal("OpenSafeWithOptions should reject allowlisted private host")
			}
		})
	}
}
