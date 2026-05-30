package socket

import "net"

// ChannelHandler is aligned with hutool nio.ChannelHandler.
// Implement this interface to handle reads and writes for a Conn.
type ChannelHandler interface {
	Handle(conn net.Conn) error
}

// ChannelHandlerFunc is a function adapter for ChannelHandler.
type ChannelHandlerFunc func(conn net.Conn) error

// Handle implements ChannelHandler.
func (f ChannelHandlerFunc) Handle(conn net.Conn) error { return f(conn) }

// Operation is aligned with hutool nio.Operation.
// Constants identify interesting event types and are mainly kept for semantic alignment.
type Operation int

const (
	// OpRead is a read operation.
	OpRead Operation = 1 << iota
	// OpWrite is a write operation.
	OpWrite
	// OpConnect is a connect operation.
	OpConnect
	// OpAccept is an accept operation.
	OpAccept
)
