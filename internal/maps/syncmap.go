package maps

import "sync"

// SyncMap is a generic, goroutine-safe map backed by sync.Map.
// It wraps the standard library's untyped sync.Map to provide type-safe
// operations over keys of type K (comparable) and values of type V (any).
// The zero value is ready to use and must not be copied after first use.
//
// Like sync.Map, it is optimized for two cases: (1) a key is written once but
// read many times, or (2) multiple goroutines operate on disjoint key sets.
// For other workloads a plain map guarded by a Mutex/RWMutex may perform better.
type SyncMap[K comparable, V any] struct {
	m sync.Map
}

// cast converts a value retrieved from the underlying sync.Map back to V.
// When loaded is false the stored slot was empty, so the zero value of V is
// returned instead of asserting on a nil any (which would panic).
func cast[V any](raw any, loaded bool) (value V, ok bool) {
	if !loaded {
		return value, false
	}
	return raw.(V), true
}

// Store sets value for key, overwriting any existing entry.
func (m *SyncMap[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}

// Load returns the value stored for key and whether it was present.
// If absent, it returns the zero value of V and false.
func (m *SyncMap[K, V]) Load(key K) (value V, ok bool) {
	return cast[V](m.m.Load(key))
}

// Delete removes key from the map. It is a no-op if key is absent.
func (m *SyncMap[K, V]) Delete(key K) {
	m.m.Delete(key)
}

// LoadAndDelete returns the value stored for key (if any) and deletes it.
// loaded reports whether key was present.
func (m *SyncMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	return cast[V](m.m.LoadAndDelete(key))
}

// LoadOrStore returns the existing value for key if present; otherwise it
// stores and returns value. loaded reports whether the value was already present.
func (m *SyncMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	raw, loaded := m.m.LoadOrStore(key, value)
	return raw.(V), loaded
}

// Range calls f sequentially for each key-value pair in the map, stopping
// early if f returns false. The iteration order is unspecified, and the set
// of pairs visited reflects a snapshot that may not include concurrent
// modifications. f may safely call other methods on the map, including Store
// and Delete (consistent with sync.Map semantics).
func (m *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	m.m.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}
