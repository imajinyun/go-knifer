package conf

import (
	"context"
	"net"
	"net/http"
	"testing"
)

func TestLoadRemoteSafeRevalidatesHostAtRoundTrip(t *testing.T) {
	lookups := [][]net.IP{{net.ParseIP("93.184.216.34")}, {net.ParseIP("127.0.0.1")}}
	lookupCount := 0
	client := &http.Client{Transport: confRoundTripperFunc(func(*http.Request) (*http.Response, error) {
		t.Fatal("unsafe request reached base transport")
		return nil, nil
	})}
	_, err := LoadRemoteSafeWithOptions("http://example.com/app.yaml", LoadOptions{
		RemoteClient: client,
		LookupIP: func(context.Context, string) ([]net.IP, error) {
			if lookupCount >= len(lookups) {
				return lookups[len(lookups)-1], nil
			}
			ips := lookups[lookupCount]
			lookupCount++
			return ips, nil
		},
	})
	if err == nil {
		t.Fatal("LoadRemoteSafeWithOptions should reject a host that resolves private during RoundTrip")
	}
	if lookupCount != 2 {
		t.Fatalf("lookup count = %d, want 2", lookupCount)
	}
}
