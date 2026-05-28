package cache

import (
	"runtime"
	"sync"
	"time"
)

// WeakCache 弱引用缓存。
//
// Go 没有 Java 的 WeakReference，这里使用 runtime.SetFinalizer 模拟：
// 当外部对持有的引用全部释放后，元素会在下一次 GC 时被回收。
// 由于 Go 的语义不同，这里弱引用语义为：缓存只持有 *V（指针）的弱视图，
// 当外部引用消失后，对应缓存条目将自动被清理。
type WeakCache[K comparable, V any] struct {
	mu       sync.Mutex
	entries  map[K]*weakEntry[V]
	timeout  time.Duration
	listener CacheListener[K, *V]
	hits     int64
	misses   int64
}

type weakEntry[V any] struct {
	ref        *V
	lastAccess int64
	ttl        time.Duration
}

// NewWeakCache 创建弱引用缓存，timeout 为默认过期时长，0 表示无限制。
func NewWeakCache[K comparable, V any](timeout time.Duration) *WeakCache[K, V] {
	return &WeakCache[K, V]{
		entries: make(map[K]*weakEntry[V]),
		timeout: timeout,
	}
}

// SetListener 设置移除监听。
func (c *WeakCache[K, V]) SetListener(l CacheListener[K, *V]) *WeakCache[K, V] {
	c.listener = l
	return c
}

// Put 放入键值对，使用默认 timeout。
func (c *WeakCache[K, V]) Put(key K, value *V) {
	c.PutWithTimeout(key, value, c.timeout)
}

// PutWithTimeout 放入键值对，自定义超时。
func (c *WeakCache[K, V]) PutWithTimeout(key K, value *V, timeout time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if old, ok := c.entries[key]; ok {
		c.notifyRemove(key, old.ref)
	}
	if value == nil {
		delete(c.entries, key)
		return
	}
	c.entries[key] = &weakEntry[V]{
		ref:        value,
		lastAccess: time.Now().UnixNano(),
		ttl:        timeout,
	}
	// 使用 finalizer 检测对象被 GC 之后清理对应缓存条目
	keyCopy := key
	cache := c
	runtime.SetFinalizer(value, func(v *V) {
		cache.removeIfRefIs(keyCopy, v)
	})
}

// Get 取值，未命中或已过期返回 nil 与 false。
func (c *WeakCache[K, V]) Get(key K) (*V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[key]
	if !ok {
		c.misses++
		return nil, false
	}
	if e.ttl > 0 && time.Now().UnixNano()-e.lastAccess > int64(e.ttl) {
		delete(c.entries, key)
		c.notifyRemove(key, e.ref)
		c.misses++
		return nil, false
	}
	e.lastAccess = time.Now().UnixNano()
	c.hits++
	return e.ref, true
}

// Remove 移除一个 key。
func (c *WeakCache[K, V]) Remove(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e, ok := c.entries[key]; ok {
		delete(c.entries, key)
		c.notifyRemove(key, e.ref)
	}
}

// Size 当前条目数（仅基于内部结构，被 GC 但还未触发 finalizer 的可能仍计入）。
func (c *WeakCache[K, V]) Size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.entries)
}

// Clear 清空缓存。
func (c *WeakCache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, e := range c.entries {
		c.notifyRemove(k, e.ref)
	}
	c.entries = make(map[K]*weakEntry[V])
}

// Prune 清理过期条目。
func (c *WeakCache[K, V]) Prune() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	count := 0
	now := time.Now().UnixNano()
	for k, e := range c.entries {
		if e.ttl > 0 && now-e.lastAccess > int64(e.ttl) {
			delete(c.entries, k)
			c.notifyRemove(k, e.ref)
			count++
		}
	}
	return count
}

// HitCount 命中数。
func (c *WeakCache[K, V]) HitCount() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.hits
}

// MissCount 未命中数。
func (c *WeakCache[K, V]) MissCount() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.misses
}

// removeIfRefIs 在对象被 GC 时由 finalizer 调用。
func (c *WeakCache[K, V]) removeIfRefIs(key K, ref *V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[key]
	if !ok {
		return
	}
	if e.ref == ref {
		delete(c.entries, key)
		c.notifyRemove(key, ref)
	}
}

func (c *WeakCache[K, V]) notifyRemove(key K, value *V) {
	if c.listener != nil {
		c.listener.OnRemove(key, value)
	}
}
