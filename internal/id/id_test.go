package id

import (
	"strings"
	"testing"
)

func TestSimpleUUID(t *testing.T) {
	u1 := SimpleUUID()
	u2 := SimpleUUID()
	if len(u1) != 32 || len(u2) != 32 {
		t.Fatalf("UUID length wrong")
	}
	if u1 == u2 {
		t.Fatalf("UUID collision")
	}
	// Version 4 marker: the 13th character is '4'.
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
