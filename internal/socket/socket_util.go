package socket

import (
	"net"
	"time"
)

// Connect creates a socket and connects to the specified address.
// When timeout <= 0, the default connection behavior without timeout is used.
func Connect(hostname string, port int, timeout time.Duration) (net.Conn, error) {
	address := net.JoinHostPort(hostname, itoa(port))
	if timeout <= 0 {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			return nil, NewSocketError(err)
		}
		return conn, nil
	}
	conn, err := net.DialTimeout("tcp", address, timeout)
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
	if conn == nil {
		return nil
	}
	return conn.RemoteAddr()
}

// IsConnected reports whether the connection is established and has a remote address.
func IsConnected(conn net.Conn) bool {
	return GetRemoteAddress(conn) != nil
}

// itoa converts an integer to a string without extra dependencies.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	negative := n < 0
	if negative {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if negative {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
