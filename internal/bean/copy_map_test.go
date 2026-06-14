package bean

import "testing"

func TestCopyPropertiesMapToStruct(t *testing.T) {
	src := map[string]any{
		"displayName": "bob",
		"age":         7.9,
		"admin":       1,
		"trace_id":    "t-2",
	}
	var dst targetProfile
	if err := Copy(src, &dst); err != nil {
		t.Fatalf("Copy() error = %v", err)
	}
	if dst.Name != "bob" || dst.Age != 7 || !dst.Admin || dst.Trace != "t-2" {
		t.Fatalf("dst = %+v", dst)
	}
}
