package cache

import "time"

// LFUCache evicts entries with the lowest access frequency.
type LFUCache[K comparable, V any] struct {
	abstractCache[K, V]
}

// NewLFUCache creates an LFU cache with no default timeout.
func NewLFUCache[K comparable, V any](capacity int) *LFUCache[K, V] {
	return NewLFUCacheWithOptions[K, V](WithCapacity[K, V](capacity))
}

// NewLFUCacheWithOptions creates an LFU cache customized by options.
func NewLFUCacheWithOptions[K comparable, V any](opts ...Option[K, V]) *LFUCache[K, V] {
	return newLFUCacheWithConfig(applyOptions(opts))
}

func newLFUCacheWithConfig[K comparable, V any](cfg cacheConfig[K, V]) *LFUCache[K, V] {
	c := &LFUCache[K, V]{}
	c.init(cfg.capacity, cfg.timeout, lfuPrune[K, V])
	applyListener(&c.abstractCache, cfg.listener)
	applyClock(&c.abstractCache, cfg.clock)
	applyTickerFactory(&c.abstractCache, cfg.tickerFactory)
	return c
}

// NewLFUCacheWithTimeout creates an LFU cache with a default timeout.
func NewLFUCacheWithTimeout[K comparable, V any](capacity int, timeout time.Duration) *LFUCache[K, V] {
	return NewLFUCacheWithOptions[K, V](WithCapacity[K, V](capacity), WithTimeout[K, V](timeout))
}

// SetListener sets the removal listener and returns the cache for chaining.
func (c *LFUCache[K, V]) SetListener(l CacheListener[K, V]) Cache[K, V] {
	c.listener = l
	return c
}

func lfuPrune[K comparable, V any](c *abstractCache[K, V]) int {
	count := 0
	var minObj *CacheObj[K, V]
	for _, key := range c.cacheMap.keysInOrder() {
		co, _ := c.cacheMap.get(key)
		if co.isExpired(c.now()) {
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
