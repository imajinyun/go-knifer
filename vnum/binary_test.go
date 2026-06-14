package vnum

import (
	"math/big"
	"reflect"
	"testing"
)

func TestNumBinaryByteFacades(t *testing.T) {
	if got, err := ToUnsignedByteArrayLen(2, big.NewInt(255)); err != nil || !reflect.DeepEqual(got, []byte{0, 255}) {
		t.Fatalf("ToUnsignedByteArrayLen = %v, %v", got, err)
	}
	if _, err := ToUnsignedByteArrayLen(1, big.NewInt(256)); err == nil {
		t.Fatal("ToUnsignedByteArrayLen overflow error = nil")
	}
	if got := FromUnsignedByteArray([]byte{1, 0}); got.Int64() != 256 {
		t.Fatalf("FromUnsignedByteArray = %s", got)
	}
	if got := FromUnsignedByteArrayRange([]byte{0, 1, 0}, 1, 2); got.Int64() != 256 {
		t.Fatalf("FromUnsignedByteArrayRange = %s", got)
	}
	if ToInt(ToBytes(12345)) != 12345 {
		t.Fatal("ToBytes/ToInt round trip failed")
	}
}
