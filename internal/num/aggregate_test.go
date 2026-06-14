package num

import (
	"math"
	"testing"
)

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

func TestAggregateEmptyAndNormalEdges(t *testing.T) {
	if Min[int]() != 0 || Max[int]() != 0 || Sum[int]() != 0 || Avg[int]() != 0 {
		t.Fatal("aggregate empty cases failed")
	}
	if Min("b", "a") != "a" || Max("b", "a") != "b" || Sum(1.5, 2.5) != 4 || Avg(1.5, 2.5) != 2 {
		t.Fatal("aggregate normal cases failed")
	}
}
