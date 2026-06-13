package json

import (
	stdjson "encoding/json"
	"io"
	"strings"
	"testing"
	"time"
)

func TestEncodeOptions(t *testing.T) {
	obj := map[string]any{"a": nil, "b": 1}
	out, err := ToJSONStr(obj, WithIgnoreNullValue(true))
	if err != nil {
		t.Fatalf("ToJSONStr: %v", err)
	}
	if out != `{"b":1}` {
		t.Fatalf("ignore null output = %s", out)
	}
	out, err = ToJSONPrettyStr(map[string]any{"a": 1}, WithIndent(2))
	if err != nil {
		t.Fatalf("ToJSONPrettyStr: %v", err)
	}
	if !strings.Contains(out, "\n  \"a\": 1") {
		t.Fatalf("pretty indent output = %s", out)
	}
	out, err = ToJSONStr(map[string]any{"t": time.Date(2026, 6, 2, 3, 4, 5, 0, time.UTC)}, WithDateFormat("2006-01-02"))
	if err != nil {
		t.Fatalf("ToJSONStr date: %v", err)
	}
	if !strings.Contains(out, "2026-06-02") {
		t.Fatalf("date output = %s", out)
	}
}

func TestEncodeOptionsUseMarshalFunc(t *testing.T) {
	type tagged struct {
		Name string `json:"name"`
	}
	called := false
	out, err := ToJSONStr(tagged{Name: "ignored"}, WithMarshalFunc(func(any) ([]byte, error) {
		called = true
		return []byte(`{"name":"provided"}`), nil
	}))
	if err != nil {
		t.Fatalf("ToJSONStr: %v", err)
	}
	if !called || out != `{"name":"provided"}` {
		t.Fatalf("marshal provider called=%v out=%s", called, out)
	}
}

func TestWrapUsesConfigDecoderFactory(t *testing.T) {
	type tagged struct {
		Name string `json:"name"`
	}
	called := false
	out, err := ToJSONStr(tagged{Name: "ignored"}, WithDecoderFactory(func(io.Reader) *stdjson.Decoder {
		called = true
		dec := stdjson.NewDecoder(strings.NewReader(`{"name":"provided"}`))
		dec.UseNumber()
		return dec
	}))
	if err != nil {
		t.Fatalf("ToJSONStr: %v", err)
	}
	if !called || out != `{"name":"provided"}` {
		t.Fatalf("decoder factory called=%v out=%s", called, out)
	}
}
