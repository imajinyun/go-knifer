package socket

import (
	"net"
	"sync"
	"sync/atomic"
)

// NioClient is an event-driven TCP client aligned with hutool nio.NioClient.
type NioClient struct {
	conn    net.Conn
	handler ChannelHandler

	closed atomic.Bool
	mu     sync.Mutex
	wg     sync.WaitGroup
}

// NewNioClient creates a client and connects to the specified host and port.
func NewNioClient(host string, port int) (*NioClient, error) {
	return NewNioClientAddr(&net.TCPAddr{IP: net.ParseIP(host), Port: port})
}

// NewNioClientAddr creates a client and connects to the specified address.
func NewNioClientAddr(addr *net.TCPAddr) (*NioClient, error) {
	if addr == nil {
		return nil, NewSocketErrorMsg("address must not be nil")
	}
	c := &NioClient{}
	if err := c.init(addr); err != nil {
		return nil, err
	}
	return c, nil
}

// init initializes the connection.
func (c *NioClient) init(addr *net.TCPAddr) error {
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return NewSocketError(err)
	}
	c.conn = conn
	return nil
}

// SetChannelHandler sets the data handler.
func (c *NioClient) SetChannelHandler(handler ChannelHandler) *NioClient {
	c.handler = handler
	return c
}

// Channel returns the underlying net.Conn.
func (c *NioClient) Channel() net.Conn {
	return c.conn
}

// Listen asynchronously listens for data pushed by the server.
func (c *NioClient) Listen() {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			if c.closed.Load() {
				return
			}
			if c.handler == nil {
				return
			}
			if err := c.handler.Handle(c.conn); err != nil {
				// Stop listening on errors such as closed connections or read failures.
				return
			}
		}
	}()
}

// Write writes multiple data fragments.
func (c *NioClient) Write(datas ...[]byte) (*NioClient, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, d := range datas {
		if len(d) == 0 {
			continue
		}
		if _, err := c.conn.Write(d); err != nil {
			return c, NewSocketError(err)
		}
	}
	return c, nil
}

// Close closes the client.
func (c *NioClient) Close() error {
	if c.closed.Swap(true) {
		return nil
	}
	if c.conn != nil {
		_ = c.conn.Close()
	}
	c.wg.Wait()
	return nil
}
