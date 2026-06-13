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

type urlRoundTripperFunc func(*http.Request) (*http.Response, error)

func (f urlRoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestOpenSafeRejectsLocalAndPrivateResources(t *testing.T) {
	if _, err := OpenSafe("file:///tmp/secret.txt"); err == nil {
		t.Fatal("OpenSafe should reject file URLs")
	}
	if _, err := OpenSafe("/tmp/secret.txt"); err == nil {
		t.Fatal("OpenSafe should reject plain file paths")
	}
	if _, err := OpenSafe("http://127.0.0.1/config.yaml"); err == nil {
		t.Fatal("OpenSafe should reject loopback hosts by default")
	}
	if _, err := OpenSafe("http://224.0.0.1/config.yaml"); err == nil {
		t.Fatal("OpenSafe should reject multicast hosts by default")
	}
	if _, err := OpenSafe("http://0.0.0.0/config.yaml"); err == nil {
		t.Fatal("OpenSafe should reject unspecified hosts by default")
	}
	if _, err := OpenSafe("ftp://example.com/config.yaml"); err == nil {
		t.Fatal("OpenSafe should reject non-HTTP schemes")
	}
}

func TestOpenSafeAllowsExplicitHostAndValidatesRedirects(t *testing.T) {
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
