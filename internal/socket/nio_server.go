package socket

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
)

// NioServer is an event-driven TCP server aligned with hutool nio.NioServer.
// In Go, goroutines plus blocking Accept/Read calls provide equivalent semantics.
type NioServer struct {
	listener net.Listener
	handler  ChannelHandler
	addr     *net.TCPAddr

	closed atomic.Bool
	wg     sync.WaitGroup
	mu     sync.Mutex
}

// NewNioServer creates and initializes a server on the given port.
func NewNioServer(port int) (*NioServer, error) {
	return NewNioServerAddr(&net.TCPAddr{Port: port})
}

// NewNioServerAddr creates a server from the specified address.
func NewNioServerAddr(addr *net.TCPAddr) (*NioServer, error) {
	if addr == nil {
		return nil, NewSocketErrorMsg("address must not be nil")
	}
	s := &NioServer{addr: addr}
	if err := s.init(addr); err != nil {
		return nil, err
	}
	return s, nil
}

// init initializes the listener.
func (s *NioServer) init(addr *net.TCPAddr) error {
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return NewSocketError(err)
	}
	s.listener = ln
	return nil
}

// SetChannelHandler sets the data handler.
func (s *NioServer) SetChannelHandler(handler ChannelHandler) *NioServer {
	s.handler = handler
	return s
}

// Listener returns the underlying net.Listener.
func (s *NioServer) Listener() net.Listener {
	return s.listener
}

// LocalAddr returns the local listen address, useful for dynamic port tests.
func (s *NioServer) LocalAddr() net.Addr {
	if s.listener == nil {
		return nil
	}
	return s.listener.Addr()
}

// Start begins listening and blocks the current goroutine.
func (s *NioServer) Start() {
	s.Listen()
}

// Listen starts synchronous blocking listening.
func (s *NioServer) Listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			if errors.Is(err, net.ErrClosed) {
				return
			}
			continue
		}
		s.handleAccept(conn)
	}
}

// ListenAsync starts listening asynchronously and closes the returned channel when done.
func (s *NioServer) ListenAsync() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		s.Listen()
	}()
	return done
}

// handleAccept handles read events from a connection in a new goroutine.
func (s *NioServer) handleAccept(conn net.Conn) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer func() { _ = conn.Close() }()
		if s.handler == nil {
			return
		}
		for {
			if s.closed.Load() {
				return
			}
			// Simulate NIO read events by invoking the handler when the connection is readable.
			// The handler usually calls conn.Read once to consume data.
			if err := s.handler.Handle(conn); err != nil {
				return
			}
		}
	}()
}

// Close closes the server.
func (s *NioServer) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed.Swap(true) {
		return nil
	}
	if s.listener != nil {
		_ = s.listener.Close()
	}
	s.wg.Wait()
	return nil
}

// IsOpen reports whether the server is still running.
func (s *NioServer) IsOpen() bool {
	return !s.closed.Load()
}
