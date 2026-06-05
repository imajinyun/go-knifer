package vskt

import (
	"bytes"
	"context"
	"net"
	"time"

	netimpl "github.com/imajinyun/go-knifer/internal/net"
	socketx "github.com/imajinyun/go-knifer/internal/socket"
)

// SocketConfig configures socket clients and servers.
type SocketConfig = socketx.SocketConfig

// ConfigOption customizes SocketConfig creation.
type ConfigOption = socketx.ConfigOption

// ConnectOption customizes socket connection helpers.
type ConnectOption = socketx.ConnectOption

// Dialer is the minimal interface used by socket connection helpers.
type Dialer = netimpl.Dialer

// SocketRuntimeError is the socket runtime error type.
type SocketRuntimeError = socketx.SocketRuntimeError

// ChannelHandler handles NIO-style connection events.
type ChannelHandler = socketx.ChannelHandler

// ChannelHandlerFunc adapts a function into ChannelHandler.
type ChannelHandlerFunc = socketx.ChannelHandlerFunc

// Operation represents NIO-style selectable operations.
type Operation = socketx.Operation

// NioServer is a NIO-style TCP server.
type NioServer = socketx.NioServer

// NioClient is a NIO-style TCP client.
type NioClient = socketx.NioClient

// AioServer is an AIO-style TCP server.
type AioServer = socketx.AioServer

// AioClient is an AIO-style TCP client.
type AioClient = socketx.AioClient

// AioSession wraps an AIO connection session.
type AioSession = socketx.AioSession

// IoAction is the AIO action callback interface.
type IoAction[T any] interface {
	socketx.IoAction[T]
}

// SimpleIoAction is a simple AIO callback implementation.
type SimpleIoAction = socketx.SimpleIoAction

// MsgDecoder decodes messages from a socket buffer.
type MsgDecoder[T any] interface {
	socketx.MsgDecoder[T]
}

// MsgEncoder encodes messages into a socket buffer.
type MsgEncoder[T any] interface {
	socketx.MsgEncoder[T]
}

// Protocol combines message encoder and decoder.
type Protocol[T any] interface {
	MsgEncoder[T]
	MsgDecoder[T]
}

// FuncDecoder adapts a function into MsgDecoder.
type FuncDecoder[T any] func(session *AioSession, readBuffer *bytes.Buffer) (T, bool)

// Decode decodes a message from readBuffer.
func (f FuncDecoder[T]) Decode(session *AioSession, readBuffer *bytes.Buffer) (T, bool) {
	return f(session, readBuffer)
}

// FuncEncoder adapts a function into MsgEncoder.
type FuncEncoder[T any] func(session *AioSession, writeBuffer *bytes.Buffer, data T)

// Encode encodes data into writeBuffer.
func (f FuncEncoder[T]) Encode(session *AioSession, writeBuffer *bytes.Buffer, data T) {
	f(session, writeBuffer, data)
}

const (
	// SocketDefaultBufferSize is the default socket buffer size.
	SocketDefaultBufferSize = socketx.DefaultBufferSize
	// OpRead represents a read operation.
	OpRead Operation = socketx.OpRead
	// OpWrite represents a write operation.
	OpWrite Operation = socketx.OpWrite
	// OpConnect represents a connect operation.
	OpConnect Operation = socketx.OpConnect
	// OpAccept represents an accept operation.
	OpAccept Operation = socketx.OpAccept
)

// NewSocketConfig creates a default socket config.
func NewSocketConfig() *SocketConfig { return socketx.NewSocketConfig() }

func WithThreadPoolSize(n int) ConfigOption { return socketx.WithThreadPoolSize(n) }

func WithReadTimeout(ms int64) ConfigOption { return socketx.WithReadTimeout(ms) }

func WithWriteTimeout(ms int64) ConfigOption { return socketx.WithWriteTimeout(ms) }

func WithReadBufferSize(n int) ConfigOption { return socketx.WithReadBufferSize(n) }

func WithWriteBufferSize(n int) ConfigOption { return socketx.WithWriteBufferSize(n) }

func NewSocketConfigWithOptions(opts ...ConfigOption) *SocketConfig {
	return socketx.NewSocketConfigWithOptions(opts...)
}

func WithConnectContext(ctx context.Context) ConnectOption { return socketx.WithConnectContext(ctx) }

func WithConnectTimeout(timeout time.Duration) ConnectOption {
	return socketx.WithConnectTimeout(timeout)
}

func WithConnectNetwork(network string) ConnectOption { return socketx.WithConnectNetwork(network) }

func WithConnectDialer(dialer Dialer) ConnectOption { return socketx.WithConnectDialer(dialer) }

// SocketConnect connects to host:port with timeout.
func SocketConnect(hostname string, port int, timeout time.Duration) (net.Conn, error) {
	return socketx.Connect(hostname, port, timeout)
}

func SocketConnectWithOptions(hostname string, port int, opts ...ConnectOption) (net.Conn, error) {
	return socketx.ConnectWithOptions(hostname, port, opts...)
}

// SocketConnectAddr connects to addr with timeout.
func SocketConnectAddr(addr *net.TCPAddr, timeout time.Duration) (net.Conn, error) {
	return socketx.ConnectAddr(addr, timeout)
}

func SocketConnectAddrWithOptions(addr *net.TCPAddr, opts ...ConnectOption) (net.Conn, error) {
	return socketx.ConnectAddrWithOptions(addr, opts...)
}

// SocketRemoteAddress returns the remote address for conn.
func SocketRemoteAddress(conn net.Conn) net.Addr { return socketx.GetRemoteAddress(conn) }

// SocketIsConnected reports whether conn is non-nil.
func SocketIsConnected(conn net.Conn) bool { return socketx.IsConnected(conn) }

// ChannelDial dials addr with timeout.
func ChannelDial(addr *net.TCPAddr, timeout time.Duration) (net.Conn, error) {
	return socketx.ChannelUtilDial(addr, timeout)
}

func ChannelDialWithOptions(addr *net.TCPAddr, opts ...ConnectOption) (net.Conn, error) {
	return socketx.ChannelUtilDialWithOptions(addr, opts...)
}

// NewNioServer creates a NIO-style TCP server on port.
func NewNioServer(port int) (*NioServer, error) { return socketx.NewNioServer(port) }

// NewNioServerAddr creates a NIO-style TCP server at addr.
func NewNioServerAddr(addr *net.TCPAddr) (*NioServer, error) { return socketx.NewNioServerAddr(addr) }

// NewNioClient creates a NIO-style TCP client.
func NewNioClient(host string, port int) (*NioClient, error) { return socketx.NewNioClient(host, port) }

// NewNioClientAddr creates a NIO-style TCP client for addr.
func NewNioClientAddr(addr *net.TCPAddr) (*NioClient, error) { return socketx.NewNioClientAddr(addr) }

// NewAioServer creates an AIO-style TCP server on port.
func NewAioServer(port int) (*AioServer, error) { return socketx.NewAioServer(port) }

// NewAioServerAddr creates an AIO-style TCP server at addr.
func NewAioServerAddr(addr *net.TCPAddr, cfg *SocketConfig) (*AioServer, error) {
	return socketx.NewAioServerAddr(addr, cfg)
}

// NewAioClient creates an AIO-style TCP client.
func NewAioClient(addr *net.TCPAddr, action IoAction[*bytes.Buffer]) (*AioClient, error) {
	return socketx.NewAioClient(addr, action)
}

func NewAioClientWithOptions(addr *net.TCPAddr, action IoAction[*bytes.Buffer], opts ...ConnectOption) (*AioClient, error) {
	return socketx.NewAioClientWithOptions(addr, action, opts...)
}

// NewAioClientWithConfig creates an AIO-style TCP client with config.
func NewAioClientWithConfig(addr *net.TCPAddr, action IoAction[*bytes.Buffer], cfg *SocketConfig) (*AioClient, error) {
	return socketx.NewAioClientWithConfig(addr, action, cfg)
}

// NewAioSession creates an AIO session from conn.
func NewAioSession(conn net.Conn, action IoAction[*bytes.Buffer], cfg *SocketConfig) *AioSession {
	return socketx.NewAioSession(conn, action, cfg)
}
