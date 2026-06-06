package cache

import "time"

// FIFOCache evicts entries in first-in-first-out order.
type FIFOCache[K comparable, V any] struct {
	abstractCache[K, V]
}

// NewFIFOCache creates a FIFO cache with the given capacity and no default timeout.
func NewFIFOCache[K comparable, V any](capacity int) *FIFOCache[K, V] {
	return NewFIFOCacheWithOptions[K, V](WithCapacity[K, V](capacity))
}

// NewFIFOCacheWithOptions creates a FIFO cache customized by options.
func NewFIFOCacheWithOptions[K comparable, V any](opts ...Option[K, V]) *FIFOCache[K, V] {
	return newFIFOCacheWithConfig(applyOptions(opts))
}

func newFIFOCacheWithConfig[K comparable, V any](cfg cacheConfig[K, V]) *FIFOCache[K, V] {
	c := &FIFOCache[K, V]{}
	c.init(cfg.capacity, cfg.timeout, fifoPrune[K, V])
	applyListener(&c.abstractCache, cfg.listener)
	applyClock(&c.abstractCache, cfg.clock)
	applyTickerFactory(&c.abstractCache, cfg.tickerFactory)
	return c
}

// NewFIFOCacheWithTimeout creates a FIFO cache with a default timeout.
func NewFIFOCacheWithTimeout[K comparable, V any](capacity int, timeout time.Duration) *FIFOCache[K, V] {
	return NewFIFOCacheWithOptions[K, V](WithCapacity[K, V](capacity), WithTimeout[K, V](timeout))
}

// SetListener sets the removal listener and returns the cache for chaining.
func (c *FIFOCache[K, V]) SetListener(l CacheListener[K, V]) Cache[K, V] {
	c.listener = l
	return c
}

func fifoPrune[K comparable, V any](c *abstractCache[K, V]) int {
	count := 0
	var first *CacheObj[K, V]
	if c.isPruneExpiredActive() {
		// Remove expired entries and remember the first non-expired entry at the
		// head side of the list as the FIFO eviction candidate.
		for _, key := range c.cacheMap.keysInOrder() {
			co, _ := c.cacheMap.get(key)
			if co.isExpired(c.now()) {
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
