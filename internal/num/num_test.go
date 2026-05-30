package num

import (
	"math"
	"testing"
)

// Tests aligned with hutool-core NumberUtilTest.

func TestNumberArith(t *testing.T) {
	if !Equals(NumberAdd(0.1, 0.2), 0.3) {
		t.Fatalf("Add failed: %v", NumberAdd(0.1, 0.2))
	}
	if !Equals(NumberSub(1.0, 0.7), 0.3) {
		t.Fatalf("Sub failed: %v", NumberSub(1.0, 0.7))
	}
	if !Equals(NumberMul(0.1, 3), 0.3) {
		t.Fatalf("Mul failed: %v", NumberMul(0.1, 3))
	}
	if got := NumberDiv(10, 3, 2); !Equals(got, 3.33) {
		t.Fatalf("Div failed: %v", got)
	}
}

func TestRound(t *testing.T) {
	if Round(3.14159, 2) != 3.14 {
		t.Fatalf("Round 3.14")
	}
	if Round(3.145, 2) != 3.15 {
		t.Fatalf("Round half up")
	}
	if Round(-3.145, 2) != -3.15 {
		t.Fatalf("Round neg half up")
	}
}

func TestIsNumber(t *testing.T) {
	if !IsNumber("123") || !IsNumber("-3.14") || IsNumber("abc") {
		t.Fatalf("IsNumber failed")
	}
	if !IsInteger("-12") || IsInteger("12.3") {
		t.Fatalf("IsInteger failed")
	}
	if !IsDigits("12345") || IsDigits("-12") {
		t.Fatalf("IsDigits failed")
	}
}

func TestMinMaxSumAvg(t *testing.T) {
	if Min(3, 1, 2) != 1 {
		t.Fatalf("Min failed")
	}
	if Max(3, 1, 2) != 3 {
		t.Fatalf("Max failed")
	}
	if Sum(1, 2, 3, 4) != 10 {
		t.Fatalf("Sum failed")
	}
	if math.Abs(Avg(1, 2, 3, 4)-2.5) > 1e-9 {
		t.Fatalf("Avg failed")
	}
}

func TestRangeFunc(t *testing.T) {
	r := Range(0, 5, 1)
	if len(r) != 5 || r[0] != 0 || r[4] != 4 {
		t.Fatalf("Range asc: %v", r)
	}
	r = Range(5, 0, -1)
	if len(r) != 5 || r[0] != 5 || r[4] != 1 {
		t.Fatalf("Range desc: %v", r)
	}
	r = Range(0, 10, 3)
	if len(r) != 4 || r[3] != 9 {
		t.Fatalf("Range step: %v", r)
	}
}
