package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// SimpleServer 简单 HTTP 服务器（对应 hutool-http SimpleServer）。
type SimpleServer struct {
	addr   string
	mux    *http.ServeMux
	server *http.Server
}

// NewSimpleServer 在指定端口创建简单服务器。
func NewSimpleServer(port int) *SimpleServer {
	return NewSimpleServerAddr(fmt.Sprintf(":%d", port))
}

// NewSimpleServerAddr 使用监听地址创建简单服务器。
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

// AddAction 注册路径处理器。
func (s *SimpleServer) AddAction(path string, handler http.HandlerFunc) *SimpleServer {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	s.mux.HandleFunc(path, handler)
	return s
}

// AddHandler 使用 http.Handler 注册。
func (s *SimpleServer) AddHandler(path string, handler http.Handler) *SimpleServer {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	s.mux.Handle(path, handler)
	return s
}

// SetRoot 设置静态文件根目录。
func (s *SimpleServer) SetRoot(dir string) *SimpleServer {
	s.mux.Handle("/", http.FileServer(http.Dir(dir)))
	return s
}

// Start 同步阻塞启动服务。
func (s *SimpleServer) Start() error {
	return s.server.ListenAndServe()
}

// StartAsync 异步启动服务，返回错误 channel。
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

// Stop 优雅关闭服务。
func (s *SimpleServer) Stop(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}

// CreateServer 创建简单 HTTP 服务器（对应 HttpUtil.createServer）。
func CreateServer(port int) *SimpleServer { return NewSimpleServer(port) }
