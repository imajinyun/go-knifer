package cache

import (
	"sync"
	"time"
)

// TimedCache 定时缓存（对应 hutool-cache TimedCache）。
// 没有容量限制，对象只在过期后被移除。
type TimedCache[K comparable, V any] struct {
	abstractCache[K, V]

	pruneStop chan struct{}
	pruneOnce sync.Once
	pruneWG   sync.WaitGroup
}

// NewTimedCache 创建定时缓存，timeout 为默认过期时长。
func NewTimedCache[K comparable, V any](timeout time.Duration) *TimedCache[K, V] {
	c := &TimedCache[K, V]{}
	c.init(0, timeout, timedPrune[K, V])
	return c
}

// SetListener 设置监听。
func (c *TimedCache[K, V]) SetListener(l CacheListener[K, V]) Cache[K, V] {
	c.listener = l
	return c
}

func timedPrune[K comparable, V any](c *abstractCache[K, V]) int {
	count := 0
	for _, key := range c.cacheMap.keysInOrder() {
		co, _ := c.cacheMap.get(key)
		if co.isExpired() {
			c.removeWithoutLock(key)
			count++
		}
	}
	return count
}

// SchedulePrune 启动定时清理任务，delay 为间隔时长。
// 多次调用会在第一次启动；若已停止则会重新启动。
func (c *TimedCache[K, V]) SchedulePrune(delay time.Duration) {
	c.pruneStop = make(chan struct{})
	c.pruneWG.Add(1)
	go func() {
		defer c.pruneWG.Done()
		ticker := time.NewTicker(delay)
		defer ticker.Stop()
		for {
			select {
			case <-c.pruneStop:
				return
			case <-ticker.C:
				c.Prune()
			}
		}
	}()
}

// CancelPruneSchedule 取消定时清理任务。
func (c *TimedCache[K, V]) CancelPruneSchedule() {
	if c.pruneStop != nil {
		select {
		case <-c.pruneStop:
			// 已经关闭
		default:
			close(c.pruneStop)
		}
		c.pruneWG.Wait()
		c.pruneStop = nil
	}
}
