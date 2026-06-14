package vnet_test

import (
	"context"
	"io"
	stdnet "net"
	"testing"
	"time"

	"github.com/imajinyun/go-knifer/vnet"
)

type recordingDialer struct {
	network string
	address string
	data    chan []byte
}

func (d *recordingDialer) DialContext(_ context.Context, network, address string) (stdnet.Conn, error) {
	d.network = network
	d.address = address
	client, server := stdnet.Pipe()
	go func() {
		defer func() { _ = server.Close() }()
		payload, _ := io.ReadAll(server)
		d.data <- payload
	}()
	return client, nil
}

func TestVNetConnectOptionsFacade(t *testing.T) {
	dialer := &recordingDialer{data: make(chan []byte, 1)}
	conn, err := vnet.ConnectWithOptions(
		"example.com", 8080,
		vnet.WithConnectNetwork("tcp4"),
		vnet.WithConnectTimeout(time.Second),
		vnet.WithConnectDialer(dialer),
	)
	if err != nil {
		t.Fatalf("ConnectWithOptions: %v", err)
	}
	_ = conn.Close()
	if dialer.network != "tcp4" || dialer.address != "example.com:8080" {
		t.Fatalf("dial target = %s %s", dialer.network, dialer.address)
	}

	dialer = &recordingDialer{data: make(chan []byte, 1)}
	if err := vnet.NetCatWithOptions("127.0.0.1", 1234, []byte("hello"), vnet.WithConnectDialer(dialer)); err != nil {
		t.Fatalf("NetCatWithOptions: %v", err)
	}
	if got := string(<-dialer.data); got != "hello" {
		t.Fatalf("NetCatWithOptions wrote %q", got)
	}

	addr := &stdnet.TCPAddr{IP: stdnet.ParseIP("127.0.0.1"), Port: 4321}
	dialer = &recordingDialer{data: make(chan []byte, 1)}
	if !vnet.IsOpenWithOptions(addr, vnet.WithConnectDialer(dialer)) {
		t.Fatal("IsOpenWithOptions should report true")
	}
}
