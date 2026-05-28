package cache

import (
	"sync"
	"sync/atomic"
	"time"
)

// pruneStrategy 子类提供的 prune 策略，返回清理数量。
// 调用时持有 abstractCache.mu 写锁。
type pruneStrategy[K comparable, V any] func(c *abstractCache[K, V]) int

// abstractCache 抽象缓存，对应 hutool-cache AbstractCache + ReentrantCache。
// 内部使用读写锁（put/prune 加写锁，get 加写锁以更新访问时间和命中计数）。
type abstractCache[K comparable, V any] struct {
	mu        sync.Mutex
	cacheMap  *linkedMap[K, V]
	capacity  int
	timeout   time.Duration
	pruneFn   pruneStrategy[K, V]
	listener  CacheListener[K, V]
	hitCount  int64
	missCount int64

	// existCustomTimeout 标记是否存在自定义 ttl 元素。
	existCustomTimeout bool

	// moveToBackOnGet 取值后是否将节点移动到链表尾部（LRU）。
	moveToBackOnGet bool

	// keyLocks 给 GetOrLoad 提供按 key 加锁能力。
	keyLocks sync.Map
}

func (c *abstractCache[K, V]) init(capacity int, timeout time.Duration, prune pruneStrategy[K, V]) {
	c.capacity = capacity
	c.timeout = timeout
	c.pruneFn = prune
	c.cacheMap = newLinkedMap[K, V](capacity)
}

func (c *abstractCache[K, V]) Capacity() int          { return c.capacity }
func (c *abstractCache[K, V]) Timeout() time.Duration { return c.timeout }
func (c *abstractCache[K, V]) HitCount() int64        { return atomic.LoadInt64(&c.hitCount) }
func (c *abstractCache[K, V]) MissCount() int64       { return atomic.LoadInt64(&c.missCount) }

// IsFull 是否已满（容量限制）。
func (c *abstractCache[K, V]) IsFull() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.isFullLocked()
}

func (c *abstractCache[K, V]) isFullLocked() bool {
	return c.capacity > 0 && c.cacheMap.size() >= c.capacity
}

// isPruneExpiredActive 是否需要清理过期对象。
func (c *abstractCache[K, V]) isPruneExpiredActive() bool {
	return c.timeout > 0 || c.existCustomTimeout
}

// Size 返回缓存大小。
func (c *abstractCache[K, V]) Size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cacheMap.size()
}

// IsEmpty 是否为空。
func (c *abstractCache[K, V]) IsEmpty() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cacheMap.size() == 0
}

// Put 加入元素，使用默认过期时长。
func (c *abstractCache[K, V]) Put(key K, value V) {
	c.PutWithTimeout(key, value, c.timeout)
}

// PutWithTimeout 加入元素，指定过期时长。
func (c *abstractCache[K, V]) PutWithTimeout(key K, value V, timeout time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.putLocked(key, value, timeout)
}

func (c *abstractCache[K, V]) putLocked(key K, value V, timeout time.Duration) {
	co := newCacheObj(key, value, timeout)
	if timeout > 0 && timeout != c.timeout {
		c.existCustomTimeout = true
	}
	if old, ok := c.cacheMap.get(key); ok {
		// 已存在则覆盖，不再做满队列清理，对齐 issue#3618
		c.cacheMap.putBack(key, co)
		c.notifyRemove(old.key, old.value)
		return
	}
	if c.isFullLocked() {
		c.pruneFn(c)
	}
	c.cacheMap.putBack(key, co)
}

// Get 取值，未命中或已过期返回零值与 false（默认刷新访问时间）。
func (c *abstractCache[K, V]) Get(key K) (V, bool) {
	return c.GetWithUpdate(key, true)
}

// GetWithUpdate 取值，可选是否刷新访问时间。
func (c *abstractCache[K, V]) GetWithUpdate(key K, updateLastAccess bool) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.getLocked(key, updateLastAccess)
}

func (c *abstractCache[K, V]) getLocked(key K, updateLastAccess bool) (V, bool) {
	var zero V
	co, ok := c.cacheMap.get(key)
	if !ok {
		atomic.AddInt64(&c.missCount, 1)
		return zero, false
	}
	if co.isExpired() {
		// 过期，移除并触发监听
		c.cacheMap.remove(key)
		c.notifyRemove(co.key, co.value)
		atomic.AddInt64(&c.missCount, 1)
		return zero, false
	}
	v := co.get(updateLastAccess)
	atomic.AddInt64(&c.hitCount, 1)
	c.afterGet(key)
	return v, true
}

// afterGet 子类钩子：用于 LRU 把 key 调到尾部。
func (c *abstractCache[K, V]) afterGet(key K) {
	if c.moveToBackOnGet {
		c.cacheMap.moveToBack(key)
	}
}

// GetOrLoad 缺失时调用 supplier 生成并存入。
func (c *abstractCache[K, V]) GetOrLoad(key K, supplier Supplier[V]) (V, error) {
	return c.GetOrLoadWith(key, true, c.timeout, supplier)
}

// GetOrLoadWith 缺失时调用 supplier，可指定是否刷新访问时间与超时。
func (c *abstractCache[K, V]) GetOrLoadWith(key K, updateLastAccess bool, timeout time.Duration, supplier Supplier[V]) (V, error) {
	if v, ok := c.GetWithUpdate(key, updateLastAccess); ok {
		return v, nil
	}
	if supplier == nil {
		var zero V
		return zero, nil
	}
	// 双重检查锁
	lockAny, _ := c.keyLocks.LoadOrStore(key, &sync.Mutex{})
	lock := lockAny.(*sync.Mutex)
	lock.Lock()
	defer func() {
		lock.Unlock()
		c.keyLocks.Delete(key)
	}()
	if v, ok := c.GetWithUpdate(key, updateLastAccess); ok {
		return v, nil
	}
	v, err := supplier()
	if err != nil {
		return v, err
	}
	c.PutWithTimeout(key, v, timeout)
	return v, nil
}

// ContainsKey 是否包含 key（命中检查时也会触发过期清理）。
func (c *abstractCache[K, V]) ContainsKey(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	co, ok := c.cacheMap.get(key)
	if !ok {
		return false
	}
	if co.isExpired() {
		c.cacheMap.remove(key)
		c.notifyRemove(co.key, co.value)
		return false
	}
	return true
}

// Remove 移除一个 key。
func (c *abstractCache[K, V]) Remove(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if old, ok := c.cacheMap.remove(key); ok {
		c.notifyRemove(old.key, old.value)
	}
}

// Clear 清空缓存，并对每个元素触发监听。
func (c *abstractCache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, co := range c.cacheMap.valuesInOrder() {
		c.notifyRemove(co.key, co.value)
	}
	c.cacheMap.clear()
}

// Prune 清理过期对象。
func (c *abstractCache[K, V]) Prune() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.pruneFn(c)
}

// Keys 返回所有 key 快照。
func (c *abstractCache[K, V]) Keys() []K {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cacheMap.keysInOrder()
}

// Values 返回所有未过期 value 快照。
func (c *abstractCache[K, V]) Values() []V {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]V, 0, c.cacheMap.size())
	for _, co := range c.cacheMap.valuesInOrder() {
		if !co.isExpired() {
			out = append(out, co.value)
		}
	}
	return out
}

func (c *abstractCache[K, V]) notifyRemove(key K, value V) {
	if c.listener != nil {
		c.listener.OnRemove(key, value)
	}
}

// removeWithoutLock 不加锁地移除，并触发监听。
func (c *abstractCache[K, V]) removeWithoutLock(key K) {
	if old, ok := c.cacheMap.remove(key); ok {
		c.notifyRemove(old.key, old.value)
	}
}
