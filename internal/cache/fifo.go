package cache

import "time"

// FIFOCache 先进先出缓存（对应 hutool-cache FIFOCache）。
type FIFOCache[K comparable, V any] struct {
	abstractCache[K, V]
}

// NewFIFOCache 创建容量为 capacity 的 FIFO 缓存，默认无超时。
func NewFIFOCache[K comparable, V any](capacity int) *FIFOCache[K, V] {
	return NewFIFOCacheWithTimeout[K, V](capacity, 0)
}

// NewFIFOCacheWithTimeout 创建带超时的 FIFO 缓存。
func NewFIFOCacheWithTimeout[K comparable, V any](capacity int, timeout time.Duration) *FIFOCache[K, V] {
	c := &FIFOCache[K, V]{}
	c.init(capacity, timeout, fifoPrune[K, V])
	return c
}

// SetListener 设置移除监听。
func (c *FIFOCache[K, V]) SetListener(l CacheListener[K, V]) Cache[K, V] {
	c.listener = l
	return c
}

func fifoPrune[K comparable, V any](c *abstractCache[K, V]) int {
	count := 0
	var first *CacheObj[K, V]
	if c.isPruneExpiredActive() {
		// 清理过期对象，并找到链表头部第一个未过期对象
		for _, key := range c.cacheMap.keysInOrder() {
			co, _ := c.cacheMap.get(key)
			if co.isExpired() {
				c.removeWithoutLock(key)
				count++
				continue
			}
			if first == nil {
				first = co
			}
		}
	} else {
		if k, ok := c.cacheMap.firstKey(); ok {
			first, _ = c.cacheMap.get(k)
		}
	}
	if c.isFullLocked() && first != nil {
		c.removeWithoutLock(first.key)
		count++
	}
	return count
}
