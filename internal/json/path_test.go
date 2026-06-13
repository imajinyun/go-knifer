package json

import "testing"

func TestPathGetPut(t *testing.T) {
	src := `{"a":{"b":[10,20,{"c":"hit"}]}}`
	v, _ := Parse(src)
	if got := GetByPath(v, "a.b[2].c"); got != "hit" {
		t.Fatalf("path get: %v", got)
	}
	if got := GetByPath(v, "$.a.b[0]"); got != int64(10) {
		t.Fatalf("path get with $: %v", got)
	}
	obj := v.(*JSONObject)
	if err := obj.PutByPath("a.b[1]", "X"); err != nil {
		t.Fatalf("put: %v", err)
	}
	if got := obj.GetByPath("a.b[1]"); got != "X" {
		t.Fatalf("after put: %v", got)
	}
}

func TestPathCreatesIntermediate(t *testing.T) {
	obj := NewJSONObject()
	if err := obj.PutByPath("a.b.c", 42); err != nil {
		t.Fatalf("put: %v", err)
	}
	if obj.GetByPath("a.b.c") != int64(42) {
		t.Fatalf("nested put")
	}
}
