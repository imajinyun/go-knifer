package num

import (
	cryptorand "crypto/rand"
	"testing"
)

func TestRemainingInternalHelperEdges(t *testing.T) {
	if IsLong("   ") {
		t.Fatal("IsLong blank should be false")
	}
	if Max(1, 3, 2) != 3 {
		t.Fatal("Max should update when a later value is larger")
	}
	if ParseInt("123") != 123 || ParseLong("123") != 123 {
		t.Fatal("ParseInt/ParseLong direct integer branch failed")
	}
	if stripTrailingZeros("12") != "12" || stripTrailingZeros("1.2300") != "1.23" || stripTrailingZeros("1.2300e2") != "1.2300e2" {
		t.Fatal("stripTrailingZeros cases failed")
	}
	if addThousands("123") != "123" || addThousands("1234") != "1,234" {
		t.Fatal("addThousands integer cases failed")
	}
	if got, err := Calculate("1+"); err == nil || got != 0 {
		t.Fatalf("Calculate trailing plus should fail: %v %v", got, err)
	}
	if got, err := Calculate("1*"); err == nil || got != 0 {
		t.Fatalf("Calculate trailing multiply should fail: %v %v", got, err)
	}
	oldReader := cryptorand.Reader
	cryptorand.Reader = errReader{}
	t.Cleanup(func() { cryptorand.Reader = oldReader })
	if got := secureIntn(10); got != 0 {
		t.Fatalf("secureIntn should return 0 when crypto random fails: %d", got)
	}
}
