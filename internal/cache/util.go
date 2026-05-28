package cache

import "time"

// 包级便捷构造函数（对应 hutool-cache CacheUtil）。

// NewFIFO 创建 FIFO 缓存。
func NewFIFO[K comparable, V any](capacity int) *FIFOCache[K, V] {
	return NewFIFOCache[K, V](capacity)
}

// NewFIFOWithTimeout 创建带超时的 FIFO 缓存。
func NewFIFOWithTimeout[K comparable, V any](capacity int, timeout time.Duration) *FIFOCache[K, V] {
	return NewFIFOCacheWithTimeout[K, V](capacity, timeout)
}

// NewLFU 创建 LFU 缓存。
func NewLFU[K comparable, V any](capacity int) *LFUCache[K, V] {
	return NewLFUCache[K, V](capacity)
}

// NewLFUWithTimeout 创建带超时的 LFU 缓存。
func NewLFUWithTimeout[K comparable, V any](capacity int, timeout time.Duration) *LFUCache[K, V] {
	return NewLFUCacheWithTimeout[K, V](capacity, timeout)
}

// NewLRU 创建 LRU 缓存。
func NewLRU[K comparable, V any](capacity int) *LRUCache[K, V] {
	return NewLRUCache[K, V](capacity)
}

// NewLRUWithTimeout 创建带超时的 LRU 缓存。
func NewLRUWithTimeout[K comparable, V any](capacity int, timeout time.Duration) *LRUCache[K, V] {
	return NewLRUCacheWithTimeout[K, V](capacity, timeout)
}

// NewTimed 创建定时缓存。
func NewTimed[K comparable, V any](timeout time.Duration) *TimedCache[K, V] {
	return NewTimedCache[K, V](timeout)
}

// NewTimedScheduled 创建定时缓存并启动定时清理。
func NewTimedScheduled[K comparable, V any](timeout, schedulePruneDelay time.Duration) *TimedCache[K, V] {
	c := NewTimedCache[K, V](timeout)
	c.SchedulePrune(schedulePruneDelay)
	return c
}

// NewWeak 创建弱引用缓存。
func NewWeak[K comparable, V any](timeout time.Duration) *WeakCache[K, V] {
	return NewWeakCache[K, V](timeout)
}

// NewNo 创建无缓存实现。
func NewNo[K comparable, V any]() *NoCache[K, V] {
	return NewNoCache[K, V]()
}
