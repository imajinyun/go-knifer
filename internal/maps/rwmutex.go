package maps

import "sync"

// RWMutexMap is a generic, goroutine-safe map guarded by a sync.RWMutex.
// The zero value is ready to use; the underlying map is lazily initialized
// on the first write. Adapted from sync/map_reference_test.go.
type RWMutexMap[K comparable, V any] struct {
	mu sync.RWMutex
	mp map[K]V
}

// NewRWMutexMap returns a ready-to-use RWMutexMap.
func NewRWMutexMap[K comparable, V any]() *RWMutexMap[K, V] {
	return &RWMutexMap[K, V]{}
}

// lazyInit allocates the backing map if it has not been created yet.
// Callers must hold the write lock.
func (m *RWMutexMap[K, V]) lazyInit() {
	if m.mp == nil {
		m.mp = make(map[K]V)
	}
}

// Load returns the value stored for key and whether it was present.
func (m *RWMutexMap[K, V]) Load(key K) (value V, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok = m.mp[key]
	return
}

// Store sets value for key, overwriting any existing entry.
func (m *RWMutexMap[K, V]) Store(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.lazyInit()
	m.mp[key] = value
}

// LoadOrStore returns the existing value for key if present; otherwise it
// stores and returns value. loaded reports whether the value was already present.
func (m *RWMutexMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if actual, loaded = m.mp[key]; loaded {
		return actual, true
	}
	m.lazyInit()
	m.mp[key] = value
	return value, false
}

// Delete removes key from the map. It is a no-op if key is absent.
func (m *RWMutexMap[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.mp, key)
}

// Len returns the number of entries in the map.
func (m *RWMutexMap[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.mp)
}

// Range calls fn sequentially for each key-value pair in the map.
// Iteration stops early if fn returns false. The order is unspecified.
//
// The read lock is held for the entire iteration, so fn must not call any
// method on the same map, or it will deadlock.
func (m *RWMutexMap[K, V]) Range(fn func(key K, value V) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for key, value := range m.mp {
		if !fn(key, value) {
			break
		}
	}
}

// Keys returns a snapshot of all keys. The order is unspecified.
func (m *RWMutexMap[K, V]) Keys() []K {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return Keys(m.mp)
}

// Values returns a snapshot of all values. The order is unspecified.
func (m *RWMutexMap[K, V]) Values() []V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return Values(m.mp)
}
