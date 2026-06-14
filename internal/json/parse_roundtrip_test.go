package json

import "testing"

func TestParseAndStringify(t *testing.T) {
	src := `{"name":"alice","age":30,"tags":["a","b"],"meta":{"x":1}}`
	v, err := Parse(src)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	obj, ok := v.(*JSONObject)
	if !ok {
		t.Fatalf("expect *JSONObject, got %T", v)
	}
	if obj.GetString("name") != "alice" {
		t.Fatalf("name: %v", obj.GetString("name"))
	}
	if obj.GetInt("age") != 30 {
		t.Fatalf("age")
	}
	arr := obj.GetJSONArray("tags")
	if arr == nil || arr.Len() != 2 || arr.GetString(0) != "a" {
		t.Fatalf("tags %v", arr)
	}
	if obj.GetJSONObject("meta").GetInt("x") != 1 {
		t.Fatalf("meta.x")
	}
	out := obj.String()
	if out != src {
		t.Fatalf("round-trip mismatch:\n  in : %s\n  out: %s", src, out)
	}
}
