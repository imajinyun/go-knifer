package socket

import "runtime"

// DefaultBufferSize is aligned with hutool IoUtil.DEFAULT_BUFFER_SIZE at 8 KB.
const DefaultBufferSize = 8 * 1024

// SocketConfig is aligned with hutool SocketConfig.
// It provides thread-pool size, timeout, buffer-size, and related socket options.
type SocketConfig struct {
	// ThreadPoolSize is the shared pool size and maps to the concurrency limit for accepting and handling connections in Go.
	ThreadPoolSize int

	// ReadTimeout is the read timeout in milliseconds; <= 0 means no timeout.
	ReadTimeout int64
	// WriteTimeout is the write timeout in milliseconds; <= 0 means no timeout.
	WriteTimeout int64

	// ReadBufferSize is the read buffer size.
	ReadBufferSize int
	// WriteBufferSize is the write buffer size.
	WriteBufferSize int
}

// NewSocketConfig creates the default configuration.
func NewSocketConfig() *SocketConfig {
	return &SocketConfig{
		ThreadPoolSize:  runtime.NumCPU(),
		ReadBufferSize:  DefaultBufferSize,
		WriteBufferSize: DefaultBufferSize,
	}
}

// SetThreadPoolSize sets the thread-pool size.
func (c *SocketConfig) SetThreadPoolSize(n int) *SocketConfig {
	c.ThreadPoolSize = n
	return c
}

// SetReadTimeout sets the read timeout in milliseconds.
func (c *SocketConfig) SetReadTimeout(ms int64) *SocketConfig {
	c.ReadTimeout = ms
	return c
}

// SetWriteTimeout sets the write timeout in milliseconds.
func (c *SocketConfig) SetWriteTimeout(ms int64) *SocketConfig {
	c.WriteTimeout = ms
	return c
}

// SetReadBufferSize sets the read buffer size.
func (c *SocketConfig) SetReadBufferSize(n int) *SocketConfig {
	c.ReadBufferSize = n
	return c
}

// SetWriteBufferSize sets the write buffer size.
func (c *SocketConfig) SetWriteBufferSize(n int) *SocketConfig {
	c.WriteBufferSize = n
	return c
}
