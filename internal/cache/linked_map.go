package cache

// linkedNode 双向链表节点。
type linkedNode[K comparable, V any] struct {
	key   K
	value *CacheObj[K, V]
	prev  *linkedNode[K, V]
	next  *linkedNode[K, V]
}

// linkedMap 提供 O(1) 的 get/put/remove，并支持把节点移到尾部（用于 LRU）。
// 链表头部为最早元素，尾部为最近元素。
type linkedMap[K comparable, V any] struct {
	m    map[K]*linkedNode[K, V]
	head *linkedNode[K, V]
	tail *linkedNode[K, V]
}

func newLinkedMap[K comparable, V any](initialCap int) *linkedMap[K, V] {
	if initialCap < 0 {
		initialCap = 0
	}
	return &linkedMap[K, V]{m: make(map[K]*linkedNode[K, V], initialCap)}
}

func (lm *linkedMap[K, V]) size() int { return len(lm.m) }

func (lm *linkedMap[K, V]) get(key K) (*CacheObj[K, V], bool) {
	n, ok := lm.m[key]
	if !ok {
		return nil, false
	}
	return n.value, true
}

// putBack 将 key->value 放在链表尾部；若已存在则替换并保持原位置。
func (lm *linkedMap[K, V]) putBack(key K, value *CacheObj[K, V]) (old *CacheObj[K, V], existed bool) {
	if n, ok := lm.m[key]; ok {
		old = n.value
		n.value = value
		return old, true
	}
	n := &linkedNode[K, V]{key: key, value: value}
	lm.m[key] = n
	lm.appendNode(n)
	return zeroOf[CacheObj[K, V]](), false
}

// remove 从链表与 map 中移除 key。
func (lm *linkedMap[K, V]) remove(key K) (*CacheObj[K, V], bool) {
	n, ok := lm.m[key]
	if !ok {
		return nil, false
	}
	delete(lm.m, key)
	lm.detach(n)
	return n.value, true
}

// moveToBack 把 key 对应节点移动到链表尾部。
func (lm *linkedMap[K, V]) moveToBack(key K) {
	n, ok := lm.m[key]
	if !ok {
		return
	}
	if lm.tail == n {
		return
	}
	lm.detach(n)
	lm.appendNode(n)
}

// firstKey 返回最早入队的 key（FIFO/LRU 头部），为空返回零值与 false。
func (lm *linkedMap[K, V]) firstKey() (K, bool) {
	if lm.head == nil {
		var zero K
		return zero, false
	}
	return lm.head.key, true
}

// keysInOrder 顺序返回所有 key（头到尾）。
func (lm *linkedMap[K, V]) keysInOrder() []K {
	out := make([]K, 0, len(lm.m))
	for n := lm.head; n != nil; n = n.next {
		out = append(out, n.key)
	}
	return out
}

// valuesInOrder 顺序返回所有缓存对象（头到尾）。
func (lm *linkedMap[K, V]) valuesInOrder() []*CacheObj[K, V] {
	out := make([]*CacheObj[K, V], 0, len(lm.m))
	for n := lm.head; n != nil; n = n.next {
		out = append(out, n.value)
	}
	return out
}

func (lm *linkedMap[K, V]) clear() {
	lm.m = make(map[K]*linkedNode[K, V])
	lm.head = nil
	lm.tail = nil
}

func (lm *linkedMap[K, V]) appendNode(n *linkedNode[K, V]) {
	n.prev = lm.tail
	n.next = nil
	if lm.tail != nil {
		lm.tail.next = n
	} else {
		lm.head = n
	}
	lm.tail = n
}

func (lm *linkedMap[K, V]) detach(n *linkedNode[K, V]) {
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		lm.head = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	} else {
		lm.tail = n.prev
	}
	n.prev = nil
	n.next = nil
}

// zeroOf 返回 *T 的 nil。
func zeroOf[T any]() *T { return nil }
