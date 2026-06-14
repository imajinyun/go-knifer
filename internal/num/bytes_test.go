package num

import (
	"math/big"
	"reflect"
	"testing"
)

func TestByteConversionHelpers(t *testing.T) {
	b := ToBytes(0x01020304)
	if !reflect.DeepEqual(b, []byte{1, 2, 3, 4}) || ToInt(b) != 0x01020304 {
		t.Fatal("byte conversion failed")
	}
	if got := ToInt(ToBytes(-1)); got != -1 {
		t.Fatalf("signed byte round trip failed: %d", got)
	}
	if ToInt(nil) != 0 || ToInt([]byte{1, 2, 3}) != 0 {
		t.Fatal("ToInt short input should be zero")
	}
}

func TestUnsignedByteArrayHelpers(t *testing.T) {
	unsigned, err := ToUnsignedByteArrayLen(4, big.NewInt(255))
	if err != nil || !reflect.DeepEqual(unsigned, []byte{0, 0, 0, 255}) || FromUnsignedByteArray(unsigned).Int64() != 255 {
		t.Fatal("unsigned bytes failed")
	}
	if ToUnsignedByteArray(nil) != nil {
		t.Fatal("ToUnsignedByteArray nil should be nil")
	}
	if got := ToUnsignedByteArray(big.NewInt(0)); len(got) != 0 {
		t.Fatalf("ToUnsignedByteArray zero should be empty: %v", got)
	}
	if _, err := ToUnsignedByteArrayLen(1, big.NewInt(256)); err == nil {
		t.Fatal("ToUnsignedByteArrayLen should reject values that exceed requested length")
	}
	if got, err := ToUnsignedByteArrayLen(0, big.NewInt(0)); err != nil || len(got) != 0 {
		t.Fatalf("ToUnsignedByteArrayLen zero length/value = %v, %v", got, err)
	}
	if FromUnsignedByteArray(nil).Sign() != 0 || FromUnsignedByteArrayRange([]byte{1, 2, 3, 4}, 1, 2).Int64() != 0x0203 {
		t.Fatal("FromUnsignedByteArray cases failed")
	}
	if FromUnsignedByteArrayRange([]byte{1, 2}, -1, 1).Sign() != 0 || FromUnsignedByteArrayRange([]byte{1, 2}, 1, 3).Sign() != 0 {
		t.Fatal("FromUnsignedByteArrayRange invalid ranges should be zero")
	}
}
