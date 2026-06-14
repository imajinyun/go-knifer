package json

import (
	stdjson "encoding/json"
	"io"
	"strings"
	"testing"
)

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
