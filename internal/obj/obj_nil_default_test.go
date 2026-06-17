package obj

import "testing"

func TestDefaultsApplyAcceptAndAggregates(t *testing.T) {
	value := "go"
	if DefaultIfNil(&value, "x") != "go" || DefaultIfNil[string](nil, "x") != "x" {
		t.Fatal("DefaultIfNil failed")
	}
	if got := Apply(&value, func(s string) int { return len(s) }); got != 2 {
		t.Fatalf("Apply: %d", got)
	}
	called := false
	Accept(&value, func(string) { called = true })
	if !called {
		t.Fatal("Accept not called")
	}
	if EmptyCount(nil, "", []int{}, 1) != 3 || !HasNil(1, nil) || !HasEmpty(1, "") {
		t.Fatal("aggregate checks failed")
	}
	if !IsAllEmpty(nil, "") || !IsAllNotEmpty(1, "x") {
		t.Fatal("all checks failed")
	}
}

func TestNilAliasesDefaultSuppliersAndAggregateFalseBranches(t *testing.T) {
	value := 3
	supplierCalled := false
	if got := DefaultIfNilFunc(&value, func() int {
		supplierCalled = true
		return 9
	}); got != 3 || supplierCalled {
		t.Fatalf("DefaultIfNilFunc non-nil = %d called=%v", got, supplierCalled)
	}
	if got := DefaultIfNilFunc[int](nil, func() int { return 9 }); got != 9 {
		t.Fatalf("DefaultIfNilFunc nil = %d", got)
	}
	if got := DefaultIfNilApply(&value, func(v int) string { return "v" }, "default"); got != "v" {
		t.Fatalf("DefaultIfNilApply non-nil = %q", got)
	}
	if !IsNull(nil) || IsNull(value) || !IsNotNil(value) || IsNotNil(nil) || !IsNotNull(value) || IsNotNull(nil) {
		t.Fatal("nil/null aliases returned unexpected values")
	}
	if HasNil(1, "x") || HasNull(1, "x") || HasEmpty(1, "x") {
		t.Fatal("aggregate helpers should return false when no nil or empty values exist")
	}
	if IsAllEmpty(nil, "x") || IsAllNotEmpty(1, "") || !IsAllEmpty() || !IsAllNotEmpty() {
		t.Fatal("all-empty/all-not-empty boundary checks failed")
	}
}
