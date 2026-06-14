package url

import (
	"context"
	"net"
	"net/http"
	"testing"
)

func TestOpenSafeRevalidatesHostAtRoundTrip(t *testing.T) {
	lookups := [][]net.IP{{net.ParseIP("93.184.216.34")}, {net.ParseIP("127.0.0.1")}}
	lookupCount := 0
	client := &http.Client{Transport: urlRoundTripperFunc(func(*http.Request) (*http.Response, error) {
		t.Fatal("unsafe request reached base transport")
		return nil, nil
	})}
	_, err := OpenSafeWithOptions("http://example.com/config.yaml",
		WithHTTPClient(client),
		WithLookupIP(func(context.Context, string) ([]net.IP, error) {
			if lookupCount >= len(lookups) {
				return lookups[len(lookups)-1], nil
			}
			ips := lookups[lookupCount]
			lookupCount++
			return ips, nil
		}),
	)
	if err == nil {
		t.Fatal("OpenSafeWithOptions should reject a host that resolves private during RoundTrip")
	}
	if lookupCount != 2 {
		t.Fatalf("lookup count = %d, want 2", lookupCount)
	}
}
