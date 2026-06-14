package conf

import (
	"context"
	"net"
	"testing"
)

func TestLoadRemoteSafeAllowedHostsDoesNotBypassPrivateRejection(t *testing.T) {
	if _, err := LoadRemoteSafeWithOptions("http://127.0.0.1/app.yaml", LoadOptions{RemoteAllowedHosts: []string{"127.0.0.1"}}); err == nil {
		t.Fatal("LoadRemoteSafeWithOptions should reject allowlisted loopback host")
	}

	lookupCount := 0
	_, err := LoadRemoteSafeWithOptions("http://config.example/app.yaml", LoadOptions{
		RemoteAllowedHosts: []string{"config.example"},
		LookupIP: func(context.Context, string) ([]net.IP, error) {
			lookupCount++
			return []net.IP{net.ParseIP("10.0.0.1")}, nil
		},
	})
	if err == nil {
		t.Fatal("LoadRemoteSafeWithOptions should reject allowlisted host resolving to private address")
	}
	if lookupCount == 0 {
		t.Fatal("LoadRemoteSafeWithOptions did not resolve allowlisted host for private-address validation")
	}
}
