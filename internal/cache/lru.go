package cache

import "time"

// LRUCache 最近最久未使用缓存（对应 hutool-cache LRUCache）。
type LRUCache[K comparable, V any] struct {
	abstractCache[K, V]
}

// NewLRUCache 创建 LRU 缓存。
func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	return NewLRUCacheWithTimeout[K, V](capacity, 0)
}

// NewLRUCacheWithTimeout 创建带默认超时的 LRU 缓存。
func NewLRUCacheWithTimeout[K comparable, V any](capacity int, timeout time.Duration) *LRUCache[K, V] {
	c := &LRUCache[K, V]{}
	c.init(capacity, timeout, lruPrune[K, V])
	c.moveToBackOnGet = true
	return c
}

// SetListener 设置监听。
func (c *LRUCache[K, V]) SetListener(l CacheListener[K, V]) Cache[K, V] {
	c.listener = l
	return c
}

// putLocked 重写：满了之后淘汰链表头部（最久未使用）。
// 因为 putLocked 在 abstractCache 内部使用 prune 钩子，这里直接通过 lruPrune 完成。

func lruPrune[K comparable, V any](c *abstractCache[K, V]) int {
	count := 0
	// 先清理所有过期对象
	if c.isPruneExpiredActive() {
		for _, key := range c.cacheMap.keysInOrder() {
			co, _ := c.cacheMap.get(key)
			if co.isExpired() {
				c.removeWithoutLock(key)
				count++
			}
		}
	}
	// 容量超限时淘汰链表头部
	for c.capacity > 0 && c.cacheMap.size() >= c.capacity {
		k, ok := c.cacheMap.firstKey()
		if !ok {
			break
		}
		c.removeWithoutLock(k)
		count++
	}
	return count
}
