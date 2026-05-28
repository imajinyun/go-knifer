package socket

import "runtime"

// 默认缓冲区大小（与 hutool IoUtil.DEFAULT_BUFFER_SIZE 对齐为 8KB）。
const DefaultBufferSize = 8 * 1024

// SocketConfig 对应 hutool 的 SocketConfig，
// 提供 Socket 通讯过程中线程池大小、超时、缓冲区等配置。
type SocketConfig struct {
	// ThreadPoolSize 共享线程池大小，对应 Go 中接收/处理连接的并发上限。
	ThreadPoolSize int

	// ReadTimeout 读取超时（毫秒），<=0 表示默认（无超时）。
	ReadTimeout int64
	// WriteTimeout 写出超时（毫秒），<=0 表示默认（无超时）。
	WriteTimeout int64

	// ReadBufferSize 读缓冲区大小。
	ReadBufferSize int
	// WriteBufferSize 写缓冲区大小。
	WriteBufferSize int
}

// NewSocketConfig 构造默认配置。
func NewSocketConfig() *SocketConfig {
	return &SocketConfig{
		ThreadPoolSize:  runtime.NumCPU(),
		ReadBufferSize:  DefaultBufferSize,
		WriteBufferSize: DefaultBufferSize,
	}
}

// SetThreadPoolSize 设置线程池大小。
func (c *SocketConfig) SetThreadPoolSize(n int) *SocketConfig {
	c.ThreadPoolSize = n
	return c
}

// SetReadTimeout 设置读取超时（毫秒）。
func (c *SocketConfig) SetReadTimeout(ms int64) *SocketConfig {
	c.ReadTimeout = ms
	return c
}

// SetWriteTimeout 设置写出超时（毫秒）。
func (c *SocketConfig) SetWriteTimeout(ms int64) *SocketConfig {
	c.WriteTimeout = ms
	return c
}

// SetReadBufferSize 设置读取缓冲区大小。
func (c *SocketConfig) SetReadBufferSize(n int) *SocketConfig {
	c.ReadBufferSize = n
	return c
}

// SetWriteBufferSize 设置写出缓冲区大小。
func (c *SocketConfig) SetWriteBufferSize(n int) *SocketConfig {
	c.WriteBufferSize = n
	return c
}
