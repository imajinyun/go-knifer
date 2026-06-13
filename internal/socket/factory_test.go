package socket

import (
	"errors"
	"net"
	"testing"
)

func TestSocketListenerAndConnFactories(t *testing.T) {
	addr := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9999}
	listener := &fakeListener{addr: factoryFakeAddr("aio-listener")}
	aio, err := NewAioServerAddrWithOptions(addr, nil, WithListenerFactory(func(got *net.TCPAddr) (net.Listener, error) {
		if got != addr {
			return nil, errors.New("unexpected aio addr")
		}
		return listener, nil
	}))
	if err != nil {
		t.Fatalf("NewAioServerAddrWithOptions: %v", err)
	}
	if aio.Listener() != listener || aio.LocalAddr().String() != "aio-listener" {
		t.Fatalf("aio listener = %#v addr=%v", aio.Listener(), aio.LocalAddr())
	}
	closeAndReport(t, aio.Close)

	listener = &fakeListener{addr: factoryFakeAddr("nio-listener")}
	nio, err := NewNioServerAddrWithOptions(addr, nil, WithListenerFactory(func(got *net.TCPAddr) (net.Listener, error) {
		if got != addr {
			return nil, errors.New("unexpected nio addr")
		}
		return listener, nil
	}))
	if err != nil {
		t.Fatalf("NewNioServerAddrWithOptions: %v", err)
	}
	if nio.Listener() != listener || nio.LocalAddr().String() != "nio-listener" {
		t.Fatalf("nio listener = %#v addr=%v", nio.Listener(), nio.LocalAddr())
	}
	closeAndReport(t, nio.Close)

	client, server := net.Pipe()
	defer closeAndReport(t, server.Close)
	nioClient, err := NewNioClientAddrWithOptions(addr, WithConnFactory(func(got *net.TCPAddr) (net.Conn, error) {
		if got != addr {
			return nil, errors.New("unexpected client addr")
		}
		return client, nil
	}))
	if err != nil {
		t.Fatalf("NewNioClientAddrWithOptions: %v", err)
	}
	if nioClient.Channel() != client {
		t.Fatalf("nio client channel = %#v", nioClient.Channel())
	}
	closeAndReport(t, nioClient.Close)
}
