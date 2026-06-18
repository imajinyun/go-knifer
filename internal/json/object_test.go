package json

import (
	"strconv"
	"strings"
	"testing"
)

func TestObjectOrderPreserved(t *testing.T) {
	obj := NewJSONObject()
	obj.Set("c", 3).Set("a", 1).Set("b", 2)
	if got := strings.Join(obj.Keys(), ","); got != "c,a,b" {
		t.Fatalf("expect insertion order, got %s", got)
	}
	s := obj.String()
	if s != `{"c":3,"a":1,"b":2}` {
		t.Fatalf("compact: %s", s)
	}
}

func TestObjectUnmarshalJSON(t *testing.T) {
	var obj JSONObject
	if err := obj.UnmarshalJSON([]byte(`{"a":1,"b":"x"}`)); err != nil {
		t.Fatalf("UnmarshalJSON: %v", err)
	}
	if obj.GetInt("a") != 1 || obj.GetString("b") != "x" {
		t.Fatalf("after UnmarshalJSON: %s", obj.String())
	}

	var bad JSONObject
	if err := bad.UnmarshalJSON([]byte(`[1,2]`)); err == nil {
		t.Fatal("expect error for array in object UnmarshalJSON")
	}
}

func TestQuoteAndGetByPathOr(t *testing.T) {
	if got := Quote(`a"b`); got != `"a\"b"` {
		t.Fatalf("Quote = %q", got)
	}
	if got := Quote("plain"); got != `"plain"` {
		t.Fatalf("Quote plain = %q", got)
	}

	root := NewJSONObject().Set("x", 42)
	got := GetByPathOr(root, "x", 0)
	if v, ok := got.(int64); !ok || v != 42 {
		t.Fatalf("GetByPathOr found = %v (type %T), want int64(42)", got, got)
	}
	if got := GetByPathOr(root, "missing", "def"); got != "def" {
		t.Fatalf("GetByPathOr default = %v", got)
	}
}

func TestNullHandling(t *testing.T) {
	obj := NewJSONObject().Set("a", nil)
	if !obj.IsNull("a") {
		t.Fatalf("expect a is null")
	}
	if obj.String() != `{"a":null}` {
		t.Fatalf("got %s", obj.String())
	}
}

func TestJSONScalarProviders(t *testing.T) {
	cfg := NewConfig()
	cfg.ParseIntFunc = func(s string, base, bitSize int) (int64, error) {
		if s == "custom-int" {
			return 77, nil
		}
		return strconv.ParseInt(s, base, bitSize)
	}
	cfg.ParseFloatFunc = func(s string, bitSize int) (float64, error) {
		if s == "custom-float" {
			return 8.5, nil
		}
		return strconv.ParseFloat(s, bitSize)
	}
	cfg.ParseBoolFunc = func(s string) (bool, error) {
		if s == "yep" {
			return true, nil
		}
		return false, strconv.ErrSyntax
	}
	obj := NewJSONObjectWithConfig(cfg)
	obj.Set("int", "custom-int").Set("float", "custom-float").Set("bool", "yep")
	if got := obj.GetInt64("int"); got != 77 {
		t.Fatalf("custom int = %d", got)
	}
	if got := obj.GetFloat64("float"); got != 8.5 {
		t.Fatalf("custom float = %v", got)
	}
	if !obj.GetBool("bool") {
		t.Fatal("custom bool parser not used")
	}

	out, err := ToJSONStr(map[string]any{"n": int64(7), "f": 1.25},
		WithFormatIntFunc(func(v int64, base int) string { return strconv.FormatInt(v*10, base) }),
		WithFormatFloatFunc(func(v float64, fmtByte byte, prec, bitSize int) string {
			return strconv.FormatFloat(v*2, fmtByte, prec, bitSize)
		}),
	)
	if err != nil {
		t.Fatalf("ToJSONStr with scalar providers: %v", err)
	}
	if out != `{"f":2.5,"n":70}` && out != `{"n":70,"f":2.5}` {
		t.Fatalf("formatted json = %s", out)
	}

	out, err = ToJSONStr(map[customKey]string{{name: "k"}: "v"}, WithSprintFunc(func(any) string { return "custom-key" }))
	if err != nil {
		t.Fatalf("ToJSONStr with sprint provider: %v", err)
	}
	if out != `{"custom-key":"v"}` {
		t.Fatalf("sprint json = %s", out)
	}
}

type customKey struct{ name string }
