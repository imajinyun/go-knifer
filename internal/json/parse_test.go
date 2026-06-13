package json

import (
	stdjson "encoding/json"
	"io"
	"strings"
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

func TestParseRejectsTrailingContent(t *testing.T) {
	if _, err := Parse(`{"a":1} {"b":2}`); err == nil {
		t.Fatal("expected trailing content error")
	}
}

func TestParseObjAndArrayErrors(t *testing.T) {
	if _, err := ParseObj(`[1,2]`); err == nil {
		t.Fatalf("expect error parsing array as obj")
	}
	if _, err := ParseArray(`{}`); err == nil {
		t.Fatalf("expect error parsing object as array")
	}
}

func TestParseObjAndArrayWithOptionsUseUnmarshalFunc(t *testing.T) {
	objCalled := false
	obj, err := ParseObjWithOptions(`{"ignored":true}`, WithParseUnmarshalFunc(func(_ []byte, dst any) error {
		objCalled = true
		*(dst.(*any)) = map[string]any{"provided": "yes"}
		return nil
	}))
	if err != nil {
		t.Fatalf("ParseObjWithOptions: %v", err)
	}
	if !objCalled || obj.GetString("provided") != "yes" {
		t.Fatalf("object unmarshal provider called=%v obj=%s", objCalled, obj.String())
	}

	arrCalled := false
	arr, err := ParseArrayWithOptions(`["ignored"]`, WithParseUnmarshalFunc(func(_ []byte, dst any) error {
		arrCalled = true
		*(dst.(*any)) = []any{"provided"}
		return nil
	}))
	if err != nil {
		t.Fatalf("ParseArrayWithOptions: %v", err)
	}
	if !arrCalled || arr.GetString(0) != "provided" {
		t.Fatalf("array unmarshal provider called=%v arr=%s", arrCalled, arr.String())
	}
}

func TestParseWithOptionsUsesDecoderFactory(t *testing.T) {
	called := false
	v, err := ParseWithOptions(`{"ignored":true}`, WithParseDecoderFactory(func(io.Reader) *stdjson.Decoder {
		called = true
		dec := stdjson.NewDecoder(strings.NewReader(`{"provided":"yes"}`))
		dec.UseNumber()
		return dec
	}))
	if err != nil {
		t.Fatalf("ParseWithOptions decoder factory: %v", err)
	}
	obj, ok := v.(*JSONObject)
	if !called || !ok || obj.GetString("provided") != "yes" {
		t.Fatalf("decoder factory called=%v value=%#v", called, v)
	}
	if _, err := ParseWithOptions(`{"ignored":true}`, WithParseDecoderFactory(func(io.Reader) *stdjson.Decoder { return nil })); err == nil {
		t.Fatal("nil decoder factory should fail")
	}
}
