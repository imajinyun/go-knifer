package num

import (
	"reflect"
	"testing"
)

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

func TestRangeEdges(t *testing.T) {
	if got := Range(1, 5, 0); !reflect.DeepEqual(got, []int{1, 2, 3, 4}) {
		t.Fatalf("Range zero positive step: %v", got)
	}
	if got := Range(5, 1, 0); !reflect.DeepEqual(got, []int{5, 4, 3, 2}) {
		t.Fatalf("Range zero negative step: %v", got)
	}
	if got := RangeClosed(3, 3, 0); !reflect.DeepEqual(got, []int{3}) {
		t.Fatalf("RangeClosed equal endpoints: %v", got)
	}
	if got := RangeClosed(1, 5, -2); !reflect.DeepEqual(got, []int{1, 3, 5}) {
		t.Fatalf("RangeClosed should normalize step sign: %v", got)
	}
	if got := RangeClosed(5, 1, 0); !reflect.DeepEqual(got, []int{5, 4, 3, 2, 1}) {
		t.Fatalf("RangeClosed descending zero step: %v", got)
	}
}

func TestRangeClosedAndAppendRange(t *testing.T) {
	if got := RangeClosed(1, 5, 2); !reflect.DeepEqual(got, []int{1, 3, 5}) {
		t.Fatalf("RangeClosed asc: %v", got)
	}
	if got := RangeClosed(5, 1, 2); !reflect.DeepEqual(got, []int{5, 3, 1}) {
		t.Fatalf("RangeClosed desc: %v", got)
	}
	if got := AppendRange(1, 3, 1, []int{0}); !reflect.DeepEqual(got, []int{0, 1, 2, 3}) {
		t.Fatalf("AppendRange: %v", got)
	}
}
