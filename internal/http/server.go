package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// SimpleServer is a simple HTTP server, aligned with hutool-http SimpleServer.
type SimpleServer struct {
	addr   string
	mux    *http.ServeMux
	server *http.Server
}

// NewSimpleServer creates a simple server on the specified port.
func NewSimpleServer(port int) *SimpleServer {
	return NewSimpleServerAddr(fmt.Sprintf(":%d", port))
}

// NewSimpleServerAddr creates a simple server with the specified listen address.
func NewSimpleServerAddr(addr string) *SimpleServer {
	mux := http.NewServeMux()
	return &SimpleServer{
		addr: addr,
		mux:  mux,
		server: &http.Server{
			Addr:              addr,
			Handler:           mux,
			ReadHeaderTimeout: 10 * time.Second,
		},
	}
}

// AddAction registers a path handler.
func (s *SimpleServer) AddAction(path string, handler http.HandlerFunc) *SimpleServer {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	s.mux.HandleFunc(path, handler)
	return s
}

// AddHandler registers an http.Handler.
func (s *SimpleServer) AddHandler(path string, handler http.Handler) *SimpleServer {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	s.mux.Handle(path, handler)
	return s
}

// SetRoot sets the static file root directory.
func (s *SimpleServer) SetRoot(dir string) *SimpleServer {
	s.mux.Handle("/", http.FileServer(http.Dir(dir)))
	return s
}

// Start starts the server synchronously and blocks.
func (s *SimpleServer) Start() error {
	return s.server.ListenAndServe()
}

// StartAsync starts the server asynchronously and returns an error channel.
func (s *SimpleServer) StartAsync() <-chan error {
	ch := make(chan error, 1)
	go func() {
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			ch <- err
		}
		close(ch)
	}()
	return ch
}

// Stop shuts down the server gracefully.
func (s *SimpleServer) Stop(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}

// CreateServer creates a simple HTTP server, aligned with HttpUtil.createServer.
func CreateServer(port int) *SimpleServer { return NewSimpleServer(port) }
