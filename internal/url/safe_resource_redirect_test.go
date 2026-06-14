package url

import (
	"context"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
)

type urlRoundTripperFunc func(*http.Request) (*http.Response, error)

func (f urlRoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestOpenSafeValidatesRedirects(t *testing.T) {
	client := &http.Client{Transport: urlRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusFound,
			Header:     http.Header{"Location": []string{"http://127.0.0.1/private"}},
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    req,
		}, nil
	})}
	if _, err := OpenSafeWithOptions("http://example.com/redirect",
		WithHTTPClient(client),
		WithAllowedHosts("example.com"),
		WithLookupIP(func(context.Context, string) ([]net.IP, error) {
			return []net.IP{net.ParseIP("93.184.216.34")}, nil
		}),
	); err == nil {
		t.Fatal("OpenSafeWithOptions should reject unsafe redirect target")
	}
}
