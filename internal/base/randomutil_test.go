package base

import (
	"strings"
	"testing"
)

// 对应 hutool-core RandomUtilTest / IdUtilTest。

func TestRandomIntRange(t *testing.T) {
	for i := 0; i < 100; i++ {
		n := RandomIntRange(10, 20)
		if n < 10 || n >= 20 {
			t.Fatalf("RandomIntRange out of bounds: %d", n)
		}
	}
}

func TestRandomString(t *testing.T) {
	s := RandomString(10)
	if len(s) != 10 {
		t.Fatalf("RandomString len: %d", len(s))
	}
	for _, r := range s {
		if !strings.ContainsRune(BaseCharNumber, r) {
			t.Fatalf("RandomString out of charset: %q", s)
		}
	}
	if len(RandomNumbers(8)) != 8 {
		t.Fatalf("RandomNumbers len wrong")
	}
}

func TestRandomBytes(t *testing.T) {
	b := RandomBytes(16)
	if len(b) != 16 {
		t.Fatalf("RandomBytes len: %d", len(b))
	}
}

func TestFillRandomBytesFallbackKeepsLength(t *testing.T) {
	buf := make([]byte, 8)
	fillRandomBytes(buf)
	if len(buf) != 8 {
		t.Fatalf("fillRandomBytes changed len: %d", len(buf))
	}
}

func TestRandomEle(t *testing.T) {
	a := []string{"x", "y", "z"}
	got := RandomEle(a)
	found := false
	for _, v := range a {
		if got == v {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("RandomEle returned non-existing: %q", got)
	}
}

func TestSimpleUUID(t *testing.T) {
	u1 := SimpleUUID()
	u2 := SimpleUUID()
	if len(u1) != 32 || len(u2) != 32 {
		t.Fatalf("UUID length wrong")
	}
	if u1 == u2 {
		t.Fatalf("UUID collision")
	}
	// version 4 标志：第 13 位是 '4'
	if u1[12] != '4' {
		t.Fatalf("UUID version: %s", u1)
	}
}

func TestFastUUID(t *testing.T) {
	u := FastUUID()
	if len(u) != 36 || strings.Count(u, "-") != 4 {
		t.Fatalf("FastUUID format: %s", u)
	}
}

func TestObjectId(t *testing.T) {
	o := ObjectId()
	if len(o) != 24 {
		t.Fatalf("ObjectId length: %s", o)
	}
}

func TestNanoId(t *testing.T) {
	id := NanoId()
	if len(id) != 21 {
		t.Fatalf("NanoId default len: %s", id)
	}
	id = NanoIdN(10)
	if len(id) != 10 {
		t.Fatalf("NanoIdN len: %s", id)
	}
}
