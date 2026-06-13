package socket

import (
	"net"
	"testing"
	"time"
)

func TestAioServerCloseClosesActiveConnections(t *testing.T) {
	client, serverConn := net.Pipe()
	defer closeAndReport(t, client.Close)
	listener := &queuedListener{addr: factoryFakeAddr("aio"), conns: make(chan net.Conn, 1)}
	listener.conns <- serverConn
	accepted := make(chan struct{})
	aio, err := NewAioServerAddrWithOptions(&net.TCPAddr{Port: 1}, nil, WithListenerFactory(func(*net.TCPAddr) (net.Listener, error) {
		return listener, nil
	}))
	if err != nil {
		t.Fatalf("NewAioServerAddrWithOptions: %v", err)
	}
	aio.SetIoAction(&SimpleIoAction{OnAccept: func(*AioSession) { close(accepted) }})
	aio.Start(false)
	select {
	case <-accepted:
	case <-time.After(time.Second):
		t.Fatal("aio server did not accept connection")
	}
	done := make(chan error, 1)
	go func() { done <- aio.Close() }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Close error = %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("AioServer.Close blocked with active connection")
	}
}
