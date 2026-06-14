package net

import (
	stdnet "net"
	"strconv"
	"testing"
	"time"
)

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
