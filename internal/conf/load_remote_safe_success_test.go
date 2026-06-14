package conf

import (
	"context"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
)

func TestLoadRemoteSafeAllowsAllowedPublicHost(t *testing.T) {
	client := &http.Client{Transport: confRoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("app:\n  name: remote")),
			Request:    r,
		}, nil
	})}
	c, err := LoadRemoteSafeWithOptions("http://config.example/app.yaml", LoadOptions{
		RemoteClient:       client,
		RemoteAllowedHosts: []string{"config.example"},
		LookupIP: func(context.Context, string) ([]net.IP, error) {
			return []net.IP{net.ParseIP("93.184.216.34")}, nil
		},
	})
	if err != nil {
		t.Fatalf("LoadRemoteSafeWithOptions allowed public host: %v", err)
	}
	if got := c.GetByGroup("app", "name"); got != "remote" {
		t.Fatalf("remote app.name = %q", got)
	}
}
