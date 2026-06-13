package id

import (
	"bytes"
	"encoding/hex"
	"testing"
	"time"
)

func TestObjectId(t *testing.T) {
	o := ObjectId()
	if len(o) != 24 {
		t.Fatalf("ObjectId length: %s", o)
	}
}

func TestObjectIDOptions(t *testing.T) {
	obj := ObjectIdWithOptions(
		WithObjectIDTimeFunc(func() time.Time { return time.Unix(1, 0) }),
		WithObjectIDRandomReader(bytes.NewReader([]byte{1, 2, 3, 4, 5})),
		WithObjectIDCounter(func() uint32 { return 0xabcdef }),
	)
	if obj != "000000010102030405abcdef" {
		t.Fatalf("ObjectIdWithOptions = %s", obj)
	}
	if _, err := hex.DecodeString(obj); err != nil {
		t.Fatalf("ObjectIdWithOptions is not hex: %v", err)
	}
}
