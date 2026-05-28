package cache

import "time"

// CacheListener 缓存监听器，对应 hutool-cache CacheListener。
type CacheListener[K comparable, V any] interface {
	OnRemove(key K, value V)
}

// CacheListenerFunc 函数式 CacheListener。
type CacheListenerFunc[K comparable, V any] func(key K, value V)

// OnRemove 实现 CacheListener。
func (f CacheListenerFunc[K, V]) OnRemove(key K, value V) { f(key, value) }

// Supplier 在缓存缺失时用于生成新值。
type Supplier[V any] func() (V, error)

// Cache 缓存接口（对应 hutool-cache Cache）。
type Cache[K comparable, V any] interface {
	// Capacity 容量，0 表示不限。
	Capacity() int
	// Timeout 默认过期时长，0 表示无限制。
	Timeout() time.Duration
	// Put 加入元素，使用默认过期时长。
	Put(key K, value V)
	// PutWithTimeout 加入元素，指定过期时长。
	PutWithTimeout(key K, value V, timeout time.Duration)
	// Get 取值，未命中或已过期返回零值与 false。默认会刷新最近访问时间。
	Get(key K) (V, bool)
	// GetWithUpdate 取值，可指定是否刷新访问时间。
	GetWithUpdate(key K, updateLastAccess bool) (V, bool)
	// GetOrLoad 缺失时调用 supplier 生成并存入缓存。
	GetOrLoad(key K, supplier Supplier[V]) (V, error)
	// GetOrLoadWith 缺失时调用 supplier，并指定是否刷新访问时间与超时。
	GetOrLoadWith(key K, updateLastAccess bool, timeout time.Duration, supplier Supplier[V]) (V, error)
	// Remove 移除一个 key。
	Remove(key K)
	// ContainsKey 是否包含 key（命中检查也会触发过期清理）。
	ContainsKey(key K) bool
	// Size 当前缓存条目数。
	Size() int
	// IsEmpty 是否为空。
	IsEmpty() bool
	// IsFull 是否已满。
	IsFull() bool
	// Prune 清理过期对象，返回清理数量。
	Prune() int
	// Clear 清空缓存。
	Clear()
	// Keys 返回所有键的快照。
	Keys() []K
	// Values 返回所有值的快照（已过期会过滤）。
	Values() []V
	// SetListener 设置移除监听器。
	SetListener(listener CacheListener[K, V]) Cache[K, V]
	// HitCount 命中数。
	HitCount() int64
	// MissCount 未命中数。
	MissCount() int64
}
