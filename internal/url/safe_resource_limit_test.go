package url

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	neturl "net/url"
	"strings"
	"testing"
)

func TestOpenSafeMaxBytes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("abcd"))
	}))
	defer srv.Close()
	srvURL, err := neturl.Parse(srv.URL)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := OpenSafeWithOptions(srv.URL,
		WithAllowedHosts(srvURL.Hostname()),
		WithRejectPrivateHosts(false),
		WithMaxBytes(3),
	); err == nil {
		t.Fatal("OpenSafeWithOptions should reject response with ContentLength exceeding max bytes")
	}

	unknownLengthClient := &http.Client{Transport: urlRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode:       http.StatusOK,
			ContentLength:    -1,
			Body:             io.NopCloser(strings.NewReader("abcd")),
			Request:          req,
			Header:           make(http.Header),
			Uncompressed:     true,
			Close:            true,
			Proto:            "HTTP/1.1",
			ProtoMajor:       1,
			ProtoMinor:       1,
			TransferEncoding: nil,
		}, nil
	})}
	r, err := OpenSafeWithOptions("http://example.com/large",
		WithHTTPClient(unknownLengthClient),
		WithAllowedHosts("example.com"),
		WithMaxBytes(3),
		WithLookupIP(func(context.Context, string) ([]net.IP, error) {
			return []net.IP{net.ParseIP("93.184.216.34")}, nil
		}),
	)
	if err != nil {
		t.Fatalf("OpenSafeWithOptions: %v", err)
	}
	_, err = io.ReadAll(r)
	_ = r.Close()
	if err == nil {
		t.Fatal("OpenSafeWithOptions body read should fail when response exceeds max bytes")
	}
}
