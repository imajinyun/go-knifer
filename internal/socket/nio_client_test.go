package socket

import (
	"net"
	"testing"
)

func TestNioClientWithOptionsUsesIPParser(t *testing.T) {
	client, server := net.Pipe()
	defer closeAndReport(t, server.Close)

	var parsedHost string
	var dialAddr *net.TCPAddr
	nioClient, err := NewNioClientWithOptions("alias-host", 4321,
		WithSocketIPParser(func(host string) net.IP {
			parsedHost = host
			return net.ParseIP("127.0.0.2")
		}),
		WithConnFactory(func(got *net.TCPAddr) (net.Conn, error) {
			dialAddr = got
			return client, nil
		}),
	)
	if err != nil {
		t.Fatalf("NewNioClientWithOptions: %v", err)
	}
	defer closeAndReport(t, nioClient.Close)
	if parsedHost != "alias-host" {
		t.Fatalf("parsed host = %q", parsedHost)
	}
	if dialAddr == nil || !dialAddr.IP.Equal(net.ParseIP("127.0.0.2")) || dialAddr.Port != 4321 {
		t.Fatalf("dial addr = %#v", dialAddr)
	}
}
