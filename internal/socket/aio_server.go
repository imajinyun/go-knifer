package socket

import (
	"bytes"
	"errors"
	"net"
	"sync"
	"sync/atomic"
)

// AioServer is an AIO-style socket server aligned with hutool aio.AioServer.
// In Go, goroutines plus blocking reads are used to simulate AIO callback semantics.
type AioServer struct {
	listener net.Listener
	ioAction IoAction[*bytes.Buffer]
	config   *SocketConfig

	closed atomic.Bool
	wg     sync.WaitGroup
	mu     sync.Mutex
}

// NewAioServer creates a server on the given port.
func NewAioServer(port int) (*AioServer, error) {
	return NewAioServerAddr(&net.TCPAddr{Port: port}, NewSocketConfig())
}

// NewAioServerAddr creates a server from an address and configuration.
func NewAioServerAddr(addr *net.TCPAddr, cfg *SocketConfig) (*AioServer, error) {
	if addr == nil {
		return nil, NewSocketErrorMsg("address must not be nil")
	}
	if cfg == nil {
		cfg = NewSocketConfig()
	}
	s := &AioServer{config: cfg}
	if err := s.init(addr); err != nil {
		return nil, err
	}
	return s, nil
}

// init initializes the listener.
func (s *AioServer) init(addr *net.TCPAddr) error {
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return NewSocketError(err)
	}
	s.listener = ln
	return nil
}

// SetIoAction sets the IO action.
func (s *AioServer) SetIoAction(action IoAction[*bytes.Buffer]) *AioServer {
	s.ioAction = action
	return s
}

// IoAction returns the IO action.
func (s *AioServer) IoAction() IoAction[*bytes.Buffer] {
	return s.ioAction
}

// Listener returns the underlying listener.
func (s *AioServer) Listener() net.Listener {
	return s.listener
}

// LocalAddr returns the local listen address.
func (s *AioServer) LocalAddr() net.Addr {
	if s.listener == nil {
		return nil
	}
	return s.listener.Addr()
}

// Config returns the configuration.
func (s *AioServer) Config() *SocketConfig { return s.config }

// IsOpen reports whether the server is still running.
func (s *AioServer) IsOpen() bool { return !s.closed.Load() }

// Start starts the server; sync controls whether it blocks the current goroutine.
func (s *AioServer) Start(sync bool) {
	if sync {
		s.acceptLoop()
		return
	}
	go s.acceptLoop()
}

// acceptLoop keeps accepting new connections.
func (s *AioServer) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() || errors.Is(err, net.ErrClosed) {
				return
			}
			if s.ioAction != nil {
				s.ioAction.Failed(NewSocketError(err), nil)
			}
			continue
		}
		s.handleAccept(conn)
	}
}

// handleAccept creates an AioSession for each connection and triggers callbacks.
func (s *AioServer) handleAccept(conn net.Conn) {
	if s.ioAction == nil {
		_ = conn.Close()
		return
	}
	session := NewAioSession(conn, s.ioAction, s.config)
	// Trigger Accept synchronously.
	s.ioAction.Accept(session)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer func() { _ = session.Close() }()
		// Keep reading to simulate chained AIO callbacks.
		for session.IsOpen() && !s.closed.Load() {
			if !session.doRead() {
				return
			}
		}
	}()
}

// Close closes the server.
func (s *AioServer) Close() error {
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
