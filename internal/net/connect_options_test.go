package net

import (
	stdnet "net"
	"testing"
	"time"
)

func TestConnectHelpersWithOptionsUseDialerNetworkAndTimeout(t *testing.T) {
	dialer := &recordingDialer{data: make(chan []byte, 1)}
	conn, err := ConnectWithOptions(
		"example.com", 8080,
		WithConnectNetwork("tcp4"),
		WithConnectTimeout(time.Second),
		WithConnectDialer(dialer),
	)
	if err != nil {
		t.Fatalf("ConnectWithOptions: %v", err)
	}
	_ = conn.Close()
	if dialer.network != "tcp4" || dialer.address != "example.com:8080" {
		t.Fatalf("dial target = %s %s", dialer.network, dialer.address)
	}

	dialer = &recordingDialer{data: make(chan []byte, 1)}
	if err := NetCatWithOptions("127.0.0.1", 1234, []byte("hello"), WithConnectDialer(dialer)); err != nil {
		t.Fatalf("NetCatWithOptions: %v", err)
	}
	if got := string(<-dialer.data); got != "hello" {
		t.Fatalf("NetCatWithOptions wrote %q", got)
	}
	if dialer.network != "tcp" || dialer.address != "127.0.0.1:1234" {
		t.Fatalf("netcat dial target = %s %s", dialer.network, dialer.address)
	}

	addr := &stdnet.TCPAddr{IP: stdnet.ParseIP("127.0.0.1"), Port: 4321}
	dialer = &recordingDialer{data: make(chan []byte, 1)}
	if !IsOpenWithOptions(addr, WithConnectDialer(dialer), WithConnectNetwork("tcp4")) {
		t.Fatal("IsOpenWithOptions should report true for successful dialer")
	}
	if dialer.network != "tcp4" || dialer.address != addr.String() {
		t.Fatalf("is-open dial target = %s %s", dialer.network, dialer.address)
	}
}
