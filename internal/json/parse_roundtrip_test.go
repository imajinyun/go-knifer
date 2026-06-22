package json

import (
	stdjson "encoding/json"
	"reflect"
	"testing"
)

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

func FuzzParseStringifyRoundTrip(f *testing.F) {
	for _, seed := range []string{
		`null`,
		`true`,
		`42`,
		`"hello"`,
		`[1,"two",false]`,
		`{"name":"alice","age":30,"tags":["a","b"]}`,
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, src string) {
		if !stdjson.Valid([]byte(src)) {
			t.Skip()
		}
		value, err := Parse(src)
		if err != nil {
			t.Skip()
		}
		out := fuzzJSONString(t, value)
		if !stdjson.Valid([]byte(out)) {
			t.Fatalf("String() = %q, want valid JSON", out)
		}
		reparsed, err := Parse(out)
		if err != nil {
			t.Fatalf("Parse(String()) error = %v", err)
		}
		if got := fuzzJSONString(t, reparsed); !jsonSemanticallyEqual(t, got, out) {
			t.Fatalf("Parse(String()).String() = %q, want semantic JSON equal to %q", got, out)
		}
	})
}

func jsonSemanticallyEqual(t *testing.T, a, b string) bool {
	t.Helper()
	var av any
	if err := stdjson.Unmarshal([]byte(a), &av); err != nil {
		t.Fatalf("Unmarshal(%q) error = %v", a, err)
	}
	var bv any
	if err := stdjson.Unmarshal([]byte(b), &bv); err != nil {
		t.Fatalf("Unmarshal(%q) error = %v", b, err)
	}
	return reflect.DeepEqual(av, bv)
}

func fuzzJSONString(t *testing.T, value any) string {
	t.Helper()
	switch v := value.(type) {
	case jsonNull:
		return "null"
	case *JSONObject:
		return v.String()
	case *JSONArray:
		return v.String()
	default:
		b, err := stdjson.Marshal(v)
		if err != nil {
			t.Fatalf("Marshal(%T) error = %v", value, err)
		}
		return string(b)
	}
}
