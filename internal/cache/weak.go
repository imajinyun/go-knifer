package cache

import (
	"runtime"
	"sync"
	"time"
)

// WeakCache is a weak-reference-like cache for pointer values.
//
// Go does not provide Java-style WeakReference. This implementation uses
// runtime.SetFinalizer to approximate weak-reference behavior: when all strong
// references to a cached pointer disappear, a later GC cycle may run the
// finalizer and remove the corresponding entry.
//
// Because finalizer scheduling is intentionally non-deterministic in Go, callers
// should treat GC-based cleanup as eventual cleanup. TTL checks and explicit
// Prune/Remove/Clear remain deterministic.
type WeakCache[K comparable, V any] struct {
	mu        sync.Mutex
	entries   map[K]*weakEntry[V]
	timeout   time.Duration
	listener  CacheListener[K, *V]
	hits      int64
	misses    int64
	clock     func() time.Time
	finalizer func(*V, func(*V))

	// pendingRemoves stores removal notifications collected while mu is held.
	// They are delivered after unlocking to avoid listener reentry deadlocks.
	pendingRemoves []removeEvent[K, *V]
}

type weakEntry[V any] struct {
	ref        *V
	lastAccess int64
	ttl        time.Duration
}

// NewWeakCache creates a weak-reference-like cache with timeout as default TTL.
// A zero timeout means entries do not expire by time.
func NewWeakCache[K comparable, V any](timeout time.Duration) *WeakCache[K, V] {
	return NewWeakCacheWithOptions[K, V](WithTimeout[K, *V](timeout))
}

// NewWeakCacheWithOptions creates a weak-reference-like cache customized by options.
func NewWeakCacheWithOptions[K comparable, V any](opts ...Option[K, *V]) *WeakCache[K, V] {
	return newWeakCacheWithConfig[K, V](applyOptions(opts))
}

func newWeakCacheWithConfig[K comparable, V any](cfg cacheConfig[K, *V]) *WeakCache[K, V] {
	clock := cfg.clock
	if clock == nil {
		clock = time.Now
	}
	c := &WeakCache[K, V]{
		entries:  make(map[K]*weakEntry[V]),
		timeout:  cfg.timeout,
		listener: cfg.listener,
		clock:    clock,
	}
	if !cfg.finalizerOff {
		if finalizer, ok := cfg.finalizer.(func(*V, func(*V))); ok && finalizer != nil {
			c.finalizer = finalizer
		} else {
			c.finalizer = defaultWeakFinalizer[V]
		}
	}
	return c
}

func defaultWeakFinalizer[V any](value *V, finalizer func(*V)) {
	runtime.SetFinalizer(value, finalizer)
}

func (c *WeakCache[K, V]) now() time.Time {
	if c.clock != nil {
		return c.clock()
	}
	return time.Now()
}

// SetListener sets the removal listener and returns the cache for chaining.
func (c *WeakCache[K, V]) SetListener(l CacheListener[K, *V]) *WeakCache[K, V] {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.listener = l
	return c
}

// Put stores a pointer value using the default timeout.
func (c *WeakCache[K, V]) Put(key K, value *V) {
	c.PutWithTimeout(key, value, c.timeout)
}

// PutWithTimeout stores a pointer value using a custom timeout.
func (c *WeakCache[K, V]) PutWithTimeout(key K, value *V, timeout time.Duration) {
	c.mu.Lock()
	if old, ok := c.entries[key]; ok {
		c.collectRemoveLocked(key, old.ref)
	}
	if value == nil {
		delete(c.entries, key)
		listener, events := c.drainRemoveEventsLocked()
		c.mu.Unlock()
		notifyRemoveEvents(listener, events)
		return
	}
	c.entries[key] = &weakEntry[V]{
		ref:        value,
		lastAccess: c.now().UnixNano(),
		ttl:        timeout,
	}
	// Use a finalizer to remove the entry after the pointed value is collected.
	keyCopy := key
	cache := c
	if c.finalizer != nil {
		c.finalizer(value, func(v *V) {
			cache.removeIfRefIs(keyCopy, v)
		})
	}
	listener, events := c.drainRemoveEventsLocked()
	c.mu.Unlock()
	notifyRemoveEvents(listener, events)
}

// Get returns a cached pointer, or nil and false when missing or expired.
func (c *WeakCache[K, V]) Get(key K) (*V, bool) {
	c.mu.Lock()
	e, ok := c.entries[key]
	if !ok {
		c.misses++
		c.mu.Unlock()
		return nil, false
	}
	now := c.now().UnixNano()
	if e.ttl > 0 && now-e.lastAccess > int64(e.ttl) {
		delete(c.entries, key)
		c.collectRemoveLocked(key, e.ref)
		c.misses++
		listener, events := c.drainRemoveEventsLocked()
		c.mu.Unlock()
		notifyRemoveEvents(listener, events)
		return nil, false
	}
	e.lastAccess = now
	c.hits++
	c.mu.Unlock()
	return e.ref, true
}

// Remove deletes one key and notifies the removal listener when present.
func (c *WeakCache[K, V]) Remove(key K) {
	c.mu.Lock()
	if e, ok := c.entries[key]; ok {
		delete(c.entries, key)
		c.collectRemoveLocked(key, e.ref)
	}
	listener, events := c.drainRemoveEventsLocked()
	c.mu.Unlock()
	notifyRemoveEvents(listener, events)
}

// Size returns the number of entries still tracked internally.
// Values already collected by GC may still be counted until their finalizers run.
func (c *WeakCache[K, V]) Size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.entries)
}

// Clear removes all entries and notifies the listener for each value.
func (c *WeakCache[K, V]) Clear() {
	c.mu.Lock()
	for k, e := range c.entries {
		c.collectRemoveLocked(k, e.ref)
	}
	c.entries = make(map[K]*weakEntry[V])
	listener, events := c.drainRemoveEventsLocked()
	c.mu.Unlock()
	notifyRemoveEvents(listener, events)
}

// Prune removes expired entries and returns the removed count.
func (c *WeakCache[K, V]) Prune() int {
	c.mu.Lock()
	count := 0
	now := c.now().UnixNano()
	for k, e := range c.entries {
		if e.ttl > 0 && now-e.lastAccess > int64(e.ttl) {
			delete(c.entries, k)
			c.collectRemoveLocked(k, e.ref)
			count++
		}
	}
	listener, events := c.drainRemoveEventsLocked()
	c.mu.Unlock()
	notifyRemoveEvents(listener, events)
	return count
}

// HitCount returns the number of successful lookups.
func (c *WeakCache[K, V]) HitCount() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.hits
}

// MissCount returns the number of missed or expired lookups.
func (c *WeakCache[K, V]) MissCount() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.misses
}

// removeIfRefIs is called by the finalizer and removes key only when it still
// points to the same value. This avoids deleting a newer value stored under the
// same key after the finalizer was registered.
func (c *WeakCache[K, V]) removeIfRefIs(key K, ref *V) {
	c.mu.Lock()
	e, ok := c.entries[key]
	if !ok {
		c.mu.Unlock()
		return
	}
	if e.ref == ref {
		delete(c.entries, key)
		c.collectRemoveLocked(key, ref)
	}
	listener, events := c.drainRemoveEventsLocked()
	c.mu.Unlock()
	notifyRemoveEvents(listener, events)
}

func (c *WeakCache[K, V]) collectRemoveLocked(key K, value *V) {
	if c.listener != nil {
		c.pendingRemoves = append(c.pendingRemoves, removeEvent[K, *V]{key: key, value: value})
	}
}

func (c *WeakCache[K, V]) drainRemoveEventsLocked() (CacheListener[K, *V], []removeEvent[K, *V]) {
	if c.listener == nil || len(c.pendingRemoves) == 0 {
		c.pendingRemoves = nil
		return nil, nil
	}
	events := append([]removeEvent[K, *V](nil), c.pendingRemoves...)
	c.pendingRemoves = nil
	return c.listener, events
}
