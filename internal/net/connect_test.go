package net

import (
	"context"
	"io"
	stdnet "net"
	"strconv"
	"testing"
	"time"
)

type recordingDialer struct {
	network string
	address string
	err     error
	data    chan []byte
}

func (d *recordingDialer) DialContext(_ context.Context, network, address string) (stdnet.Conn, error) {
	d.network = network
	d.address = address
	if d.err != nil {
		return nil, d.err
	}
	client, server := stdnet.Pipe()
	go func() {
		defer func() { _ = server.Close() }()
		payload, _ := io.ReadAll(server)
		d.data <- payload
	}()
	return client, nil
}

func TestPingWithOptions(t *testing.T) {
	ln, err := stdnet.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen local port: %v", err)
	}
	defer func() { _ = ln.Close() }()
	done := make(chan struct{})
	go func() {
		defer close(done)
		conn, err := ln.Accept()
		if err == nil {
			_ = conn.Close()
		}
	}()
	_, portStr, err := stdnet.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatalf("split listener address: %v", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("parse listener port: %v", err)
	}
	if !PingWithOptions("127.0.0.1", WithPingPorts(port), WithPingTimeout(time.Second), WithPingNetwork("tcp")) {
		t.Fatal("PingWithOptions failed to reach local listener")
	}
	<-done
}

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

func TestResolveWithOptions(t *testing.T) {
	ips, err := GetIPByHostWithOptions("localhost", WithResolveNetwork("ip4"), WithResolveTimeout(time.Second))
	if err != nil {
		t.Fatalf("GetIPByHostWithOptions: %v", err)
	}
	if len(ips) == 0 {
		t.Fatal("GetIPByHostWithOptions returned no IPs")
	}
	dns, err := GetDNSInfoWithOptions("localhost", WithDNSTypes("A"), WithResolveTimeout(time.Second))
	if err != nil {
		t.Fatalf("GetDNSInfoWithOptions: %v", err)
	}
	if len(dns) == 0 {
		t.Fatal("GetDNSInfoWithOptions returned no A records")
	}
}
