package vnum

import (
	"reflect"
	"testing"
)

func TestNumRangeFacades(t *testing.T) {
	if got := RangeClosed(1, 3, 0); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("RangeClosed = %v", got)
	}
	if got := AppendRange(3, 1, 1, []int{0}); !reflect.DeepEqual(got, []int{0, 3, 2, 1}) {
		t.Fatalf("AppendRange = %v", got)
	}
}

func TestNumRandomGenerationFacades(t *testing.T) {
	if got := GenRandomNumberWithSeed(0, 4, 2, []int{1, 2, 3}); len(got) != 2 {
		t.Fatalf("GenRandomNumberWithSeed = %v", got)
	}
	if got := GenRandomNumber(0, 3, 2); len(got) != 2 {
		t.Fatalf("GenRandomNumber = %v", got)
	}
	if got := GenBySet(0, 3, 2); len(got) != 2 {
		t.Fatalf("GenBySet = %v", got)
	}
}
