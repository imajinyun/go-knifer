package net

import (
	"context"
	"io"
	stdnet "net"
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
