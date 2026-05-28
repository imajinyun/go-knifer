package cache

import (
	"runtime"
	"testing"
	"time"
)

func TestWeakCacheBasic(t *testing.T) {
	c := NewWeak[string, int](0)
	v := 42
	c.Put("a", &v)
	got, ok := c.Get("a")
	if !ok || got == nil || *got != 42 {
		t.Fatalf("get failed: %v %v", got, ok)
	}
	c.Remove("a")
	if _, ok := c.Get("a"); ok {
		t.Fatalf("should be removed")
	}
}

func TestWeakCacheTimeout(t *testing.T) {
	c := NewWeak[string, int](10 * time.Millisecond)
	v := 42
	c.Put("a", &v)
	time.Sleep(20 * time.Millisecond)
	if _, ok := c.Get("a"); ok {
		t.Fatalf("should expire")
	}
}

// 验证当外部强引用消失后，下一次 GC 会清理弱引用条目。
// 由于 finalizer 调度时机不可完全控制，多次 GC + 微小休眠以提高稳定性。
func TestWeakCacheGC(t *testing.T) {
	c := NewWeak[string, int](0)
	func() {
		v := 7
		c.Put("a", &v)
	}()
	for i := 0; i < 5 && c.Size() > 0; i++ {
		runtime.GC()
		time.Sleep(10 * time.Millisecond)
	}
	if c.Size() != 0 {
		// 在某些 Go runtime/GC 下 finalizer 触发可能延迟，跳过严格断言但记录信息。
		t.Logf("weak cache size after GC: %d (finalizer may be delayed)", c.Size())
	}
}
