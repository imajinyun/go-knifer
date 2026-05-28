package cache

import (
	"sync"
	"time"
)

// CacheObj 缓存对象，对应 hutool-cache CacheObj。
type CacheObj[K comparable, V any] struct {
	key         K
	value       V
	ttl         time.Duration // 0 表示永不过期
	lastAccess  int64         // Unix 纳秒
	accessCount int64         // 访问次数
	mu          sync.Mutex
}

// newCacheObj 创建一个 CacheObj。
func newCacheObj[K comparable, V any](key K, value V, ttl time.Duration) *CacheObj[K, V] {
	return &CacheObj[K, V]{
		key:        key,
		value:      value,
		ttl:        ttl,
		lastAccess: time.Now().UnixNano(),
	}
}

// Key 返回键。
func (c *CacheObj[K, V]) Key() K { return c.key }

// Value 返回值。
func (c *CacheObj[K, V]) Value() V { return c.value }

// TTL 返回过期时长。
func (c *CacheObj[K, V]) TTL() time.Duration { return c.ttl }

// LastAccess 上次访问时间。
func (c *CacheObj[K, V]) LastAccess() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return time.Unix(0, c.lastAccess)
}

// AccessCount 访问次数。
func (c *CacheObj[K, V]) AccessCount() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.accessCount
}

// ExpiredTime 过期时间，永不过期返回零值与 false。
func (c *CacheObj[K, V]) ExpiredTime() (time.Time, bool) {
	if c.ttl <= 0 {
		return time.Time{}, false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return time.Unix(0, c.lastAccess).Add(c.ttl), true
}

// isExpired 判断是否已经过期。
func (c *CacheObj[K, V]) isExpired() bool {
	if c.ttl <= 0 {
		return false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return time.Now().UnixNano()-c.lastAccess > int64(c.ttl)
}

// get 取值，更新访问计数与可选的访问时间。
func (c *CacheObj[K, V]) get(updateLastAccess bool) V {
	c.mu.Lock()
	defer c.mu.Unlock()
	if updateLastAccess {
		c.lastAccess = time.Now().UnixNano()
	}
	c.accessCount++
	return c.value
}

// addAccessCount 调整访问计数（用于 LFU 衰减），返回调整后的值。
func (c *CacheObj[K, V]) addAccessCount(delta int64) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.accessCount += delta
	return c.accessCount
}
