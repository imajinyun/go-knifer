package vid

import (
	"strings"
	"testing"
)

func TestIDFacade(t *testing.T) {
	u1 := SimpleUUID()
	u2 := UUID()
	if len(u1) != 32 || len(u2) != 32 || u1 == u2 || u1[12] != '4' {
		t.Fatalf("uuid failed: %q %q", u1, u2)
	}
	if fast := FastUUID(); len(fast) != 36 || strings.Count(fast, "-") != 4 {
		t.Fatalf("FastUUID failed: %q", fast)
	}
	if oid := ObjectId(); len(oid) != 24 {
		t.Fatalf("ObjectId failed: %q", oid)
	}
	if nid := NanoId(); len(nid) != 21 {
		t.Fatalf("NanoId failed: %q", nid)
	}
	if nid := NanoIdN(10); len(nid) != 10 {
		t.Fatalf("NanoIdN failed: %q", nid)
	}
}
