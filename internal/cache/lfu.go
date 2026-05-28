package cache

import "time"

// LFUCache 最少使用率缓存（对应 hutool-cache LFUCache）。
type LFUCache[K comparable, V any] struct {
	abstractCache[K, V]
}

// NewLFUCache 创建 LFU 缓存。
func NewLFUCache[K comparable, V any](capacity int) *LFUCache[K, V] {
	return NewLFUCacheWithTimeout[K, V](capacity, 0)
}

// NewLFUCacheWithTimeout 创建带默认超时的 LFU 缓存。
func NewLFUCacheWithTimeout[K comparable, V any](capacity int, timeout time.Duration) *LFUCache[K, V] {
	c := &LFUCache[K, V]{}
	c.init(capacity, timeout, lfuPrune[K, V])
	return c
}

// SetListener 设置监听。
func (c *LFUCache[K, V]) SetListener(l CacheListener[K, V]) Cache[K, V] {
	c.listener = l
	return c
}

func lfuPrune[K comparable, V any](c *abstractCache[K, V]) int {
	count := 0
	var minObj *CacheObj[K, V]
	for _, key := range c.cacheMap.keysInOrder() {
		co, _ := c.cacheMap.get(key)
		if co.isExpired() {
			c.removeWithoutLock(key)
			count++
			continue
		}
		if minObj == nil || co.AccessCount() < minObj.AccessCount() {
			minObj = co
		}
	}
	if c.isFullLocked() && minObj != nil {
		minAccess := minObj.AccessCount()
		for _, key := range c.cacheMap.keysInOrder() {
			co, ok := c.cacheMap.get(key)
			if !ok {
				continue
			}
			if co.addAccessCount(-minAccess) <= 0 {
				c.removeWithoutLock(key)
				count++
			}
		}
	}
	return count
}
