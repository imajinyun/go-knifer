package socket

import (
	"net"
	"time"

	netimpl "github.com/imajinyun/go-knifer/internal/net"
)

// Connect creates a socket and connects to the specified address.
// When timeout <= 0, the default connection behavior without timeout is used.
func Connect(hostname string, port int, timeout time.Duration) (net.Conn, error) {
	conn, err := netimpl.Connect(hostname, port, timeout)
	if err != nil {
		return nil, NewSocketError(err)
	}
	return conn, nil
}

// ConnectAddr creates a connection from net.TCPAddr.
func ConnectAddr(addr *net.TCPAddr, timeout time.Duration) (net.Conn, error) {
	if addr == nil {
		return nil, NewSocketErrorMsg("address must not be nil")
	}
	return Connect(addr.IP.String(), addr.Port, timeout)
}

// GetRemoteAddress returns the remote address, or nil when conn is nil or disconnected.
func GetRemoteAddress(conn net.Conn) net.Addr {
	if !netimpl.IsConnected(conn) {
		return nil
	}
	return conn.RemoteAddr()
}

// IsConnected reports whether the connection is established and has a remote address.
func IsConnected(conn net.Conn) bool {
	return netimpl.IsConnected(conn)
}
