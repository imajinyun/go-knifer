package json

import (
	stdjson "encoding/json"
	"errors"
	"io"
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestWrapAllScalarKinds drives wrap() through every primitive branch via Set.
func TestWrapAllScalarKinds(t *testing.T) {
	obj := NewJSONObject()
	obj.Set("int", int(1))
	obj.Set("int8", int8(2))
	obj.Set("int16", int16(3))
	obj.Set("int32", int32(4))
	obj.Set("int64", int64(5))
	obj.Set("uint", uint(6))
	obj.Set("uint8", uint8(7))
	obj.Set("uint16", uint16(8))
	obj.Set("uint32", uint32(9))
	obj.Set("uint64", uint64(10))
	obj.Set("float32", float32(1.5))
	obj.Set("float64", float64(2.5))
	obj.Set("bool", true)
	obj.Set("bytes", []byte("hello"))
	obj.Set("nil", nil)

	if obj.GetInt64("int64") != 5 || obj.GetInt("uint32") != 9 {
		t.Fatalf("int wrap mismatch: %s", obj)
	}
	if obj.GetFloat64("float32") != 1.5 {
		t.Fatalf("float32 wrap = %v", obj.GetFloat64("float32"))
	}
	if obj.GetString("bytes") != "hello" {
		t.Fatalf("bytes wrap = %q", obj.GetString("bytes"))
	}
	if !obj.IsNull("nil") {
		t.Fatal("nil should wrap to Null")
	}
}

func TestWrapCompositeKinds(t *testing.T) {
	obj := NewJSONObject()
	obj.Set("map", map[string]any{"k": 1})
	obj.Set("slice", []any{1, 2, 3})
	obj.Set("time", time.UnixMilli(0).UTC())

	if obj.GetJSONObject("map") == nil {
		t.Fatal("map should wrap to *JSONObject")
	}
	arr := obj.GetJSONArray("slice")
	if arr == nil || arr.Len() != 3 {
		t.Fatalf("slice should wrap to *JSONArray: %v", arr)
	}
	// time.Time with empty DateFormat wraps to UnixMilli int.
	if obj.GetInt64("time") != 0 {
		t.Fatalf("time wrap = %v", obj.GetInt64("time"))
	}
}

func TestWrapStruct(t *testing.T) {
	type inner struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	obj := NewJSONObject()
	obj.Set("person", inner{Name: "alice", Age: 30})
	person := obj.GetJSONObject("person")
	if person == nil || person.GetString("name") != "alice" || person.GetInt("age") != 30 {
		t.Fatalf("struct wrap = %v", person)
	}
}

func TestTypedGetterDefaultsAndConversions(t *testing.T) {
	obj := NewJSONObject()
	obj.Set("s", "text")
	obj.Set("i", 42)
	obj.Set("f", 3.14)
	obj.Set("b", true)
	obj.Set("numStr", "100")

	// Absent keys return defaults.
	if obj.GetStringOr("missing", "def") != "def" {
		t.Fatal("absent string default")
	}
	if obj.GetInt64Or("missing", -1) != -1 {
		t.Fatal("absent int default")
	}
	if obj.GetFloat64Or("missing", -1) != -1 {
		t.Fatal("absent float default")
	}
	if obj.GetBoolOr("missing", true) != true {
		t.Fatal("absent bool default")
	}
	// Cross-type conversions through toX helpers.
	if obj.GetInt64("numStr") != 100 {
		t.Fatalf("string->int = %d", obj.GetInt64("numStr"))
	}
	if obj.GetString("i") != "42" {
		t.Fatalf("int->string = %q", obj.GetString("i"))
	}
	if obj.GetBool("i") != true {
		t.Fatal("non-zero int -> bool true")
	}
	// Wrong-type nested getters return nil.
	if obj.GetJSONObject("s") != nil || obj.GetJSONArray("s") != nil {
		t.Fatal("scalar should not be object/array")
	}
}

func TestObjectRemoveForEachAndKeys(t *testing.T) {
	obj := NewJSONObject()
	obj.Set("a", 1).Set("b", 2).Set("c", 3)

	if !obj.Remove("b") || obj.Remove("missing") {
		t.Fatal("Remove return values")
	}
	if len(obj.Keys()) != 2 {
		t.Fatalf("keys after remove = %v", obj.Keys())
	}

	// ForEach with early stop.
	count := 0
	obj.ForEach(func(string, any) bool {
		count++
		return false
	})
	if count != 1 {
		t.Fatalf("ForEach early stop visited %d", count)
	}

	// Empty object Keys returns non-nil empty slice.
	if got := NewJSONObject().Keys(); got == nil || len(got) != 0 {
		t.Fatalf("empty Keys = %#v", got)
	}
}

func TestPathGetAndPut(t *testing.T) {
	obj, err := ParseObj(`{"user":{"name":"bob","tags":["x","y"]}}`)
	if err != nil {
		t.Fatalf("ParseObj: %v", err)
	}
	if got := obj.GetByPath("user.name"); got != "bob" {
		t.Fatalf("GetByPath = %v", got)
	}
	if got := obj.GetByPath("user.tags[1]"); got != "y" {
		t.Fatalf("GetByPath index = %v", got)
	}
	if got := obj.GetByPath("user.missing"); !IsNull(got) {
		t.Fatalf("GetByPath missing = %v", got)
	}

	if err := obj.PutByPath("user.name", "carol"); err != nil {
		t.Fatalf("PutByPath: %v", err)
	}
	if got := obj.GetByPath("user.name"); got != "carol" {
		t.Fatalf("after PutByPath = %v", got)
	}
	if err := obj.PutByPath("user.age", 41); err != nil {
		t.Fatalf("PutByPath new key: %v", err)
	}
	if obj.GetByPath("user.age") != int64(41) {
		t.Fatalf("new key value = %v", obj.GetByPath("user.age"))
	}
}

func TestParseHelpersAndValidators(t *testing.T) {
	if v, err := Parse(`{"a":1}`); err != nil {
		t.Fatalf("Parse object: %v", err)
	} else if _, ok := v.(*JSONObject); !ok {
		t.Fatalf("Parse type = %T", v)
	}
	if v, err := Parse(nil); err != nil || !IsNull(v) {
		t.Fatalf("Parse(nil) = %v, %v", v, err)
	}
	if _, err := ParseArray(`[1,2,3]`); err != nil {
		t.Fatalf("ParseArray: %v", err)
	}
	// Type-mismatch errors.
	if _, err := ParseObj(`[1]`); err == nil {
		t.Fatal("ParseObj on array should error")
	}
	if _, err := ParseArray(`{}`); err == nil {
		t.Fatal("ParseArray on object should error")
	}

	if !IsJSON(`{"a":1}`) || IsJSON("  ") || IsJSON("nope") {
		t.Fatal("IsJSON")
	}
	if !IsJSONObj(`{"a":1}`) || IsJSONObj(`[1]`) {
		t.Fatal("IsJSONObj")
	}
	if !IsJSONArray(`[1]`) || IsJSONArray(`{}`) {
		t.Fatal("IsJSONArray")
	}
}

// TestParseWithCustomProviders ensures per-call ParseOptions install providers.
func TestParseWithCustomProviders(t *testing.T) {
	called := false
	obj, err := ParseObjWithOptions(`{"a":1}`, WithParseUnmarshalFunc(func(b []byte, v any) error {
		called = true
		return stdjson.Unmarshal(b, v)
	}))
	if err != nil || obj == nil || !called {
		t.Fatalf("ParseObjWithOptions err=%v obj=%v called=%v", err, obj, called)
	}

	arr, err := ParseArrayWithOptions(`[1,2]`, WithParseConfig(NewConfig()))
	if err != nil || arr.Len() != 2 {
		t.Fatalf("ParseArrayWithOptions err=%v arr=%v", err, arr)
	}
}

// TestEncodeWithCustomNumberProviders exercises the parser/formatter providers
// stored on Config via wrap/toX during round-trips.
func TestEncodeWithCustomNumberProviders(t *testing.T) {
	cfg := NewConfig()
	intCalled := false
	cfg.ParseIntFunc = func(s string, base, bit int) (int64, error) {
		intCalled = true
		return strconv.ParseInt(s, base, bit)
	}
	obj := NewJSONObjectWithConfig(cfg)
	obj.Set("n", "123")
	if obj.GetInt64("n") != 123 || !intCalled {
		t.Fatalf("custom ParseIntFunc not used (val=%d called=%v)", obj.GetInt64("n"), intCalled)
	}
}

func TestToBeanRoundTrip(t *testing.T) {
	type person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var p person
	if err := ToBean(`{"name":"alice","age":30}`, &p); err != nil {
		t.Fatalf("ToBean string: %v", err)
	}
	if p.Name != "alice" || p.Age != 30 {
		t.Fatalf("ToBean result = %#v", p)
	}

	// From a *JSONObject (default branch -> wrap + serialize).
	obj := NewJSONObject().Set("name", "bob").Set("age", 41)
	var p2 person
	if err := ToBean(obj, &p2); err != nil {
		t.Fatalf("ToBean object: %v", err)
	}
	if p2.Name != "bob" || p2.Age != 41 {
		t.Fatalf("ToBean object result = %#v", p2)
	}

	// From []byte and a slice via ToList.
	var nums []int
	if err := ToList([]byte(`[1,2,3]`), &nums); err != nil || len(nums) != 3 {
		t.Fatalf("ToList = %v err=%v", nums, err)
	}

	// nil dst is rejected.
	if err := ToBean(`{}`, nil); err == nil {
		t.Fatal("ToBean(nil dst) should error")
	}
}

func TestIgnoreCaseConfig(t *testing.T) {
	cfg := NewConfig()
	cfg.IgnoreCase = true
	obj := NewJSONObjectWithConfig(cfg)
	obj.Set("Name", "alice")

	if !obj.Has("name") || !obj.Has("NAME") {
		t.Fatal("IgnoreCase Has should match any casing")
	}
	if obj.GetString("nAmE") != "alice" {
		t.Fatalf("IgnoreCase Get = %q", obj.GetString("nAmE"))
	}
	// Overwrite via different casing keeps a single key.
	obj.Set("NAME", "bob")
	if obj.Len() != 1 || obj.GetString("name") != "bob" {
		t.Fatalf("IgnoreCase overwrite: len=%d val=%q", obj.Len(), obj.GetString("name"))
	}
	if !obj.Remove("name") || obj.Len() != 0 {
		t.Fatal("IgnoreCase Remove should drop the canonical key")
	}
}

func TestIgnoreNullValueConfig(t *testing.T) {
	cfg := NewConfig()
	cfg.IgnoreNullValue = true
	obj := NewJSONObjectWithConfig(cfg)
	obj.Set("a", 1).Set("b", nil)
	if obj.Has("b") {
		t.Fatal("IgnoreNullValue should drop null entries")
	}
	if obj.Len() != 1 {
		t.Fatalf("len = %d", obj.Len())
	}
}

func TestWriteQuotedEscapes(t *testing.T) {
	obj := NewJSONObject()
	obj.Set("s", "line1\nline2\t\"quoted\"\\back/slash")
	out := obj.String()
	// Ensure the special characters are escaped in the serialized output.
	for _, want := range []string{`\n`, `\t`, `\"`, `\\`} {
		if !strings.Contains(out, want) {
			t.Fatalf("output %q missing escape %q", out, want)
		}
	}
}

func TestPrettyOutput(t *testing.T) {
	obj := NewJSONObject().Set("a", 1).Set("b", NewJSONObject().Set("c", 2))
	pretty := obj.ToStringPretty()
	if !strings.Contains(pretty, "\n") {
		t.Fatalf("pretty output should contain newlines: %q", pretty)
	}
	if s, err := ToJSONPrettyStr(map[string]any{"x": 1}); err != nil || !strings.Contains(s, "\n") {
		t.Fatalf("ToJSONPrettyStr = %q err=%v", s, err)
	}
}

func TestMapSerializationIsDeterministic(t *testing.T) {
	input := map[string]any{
		"zeta":  1,
		"alpha": 2,
		"mid":   map[string]any{"b": true, "a": false},
	}
	const want = `{"alpha":2,"mid":{"a":false,"b":true},"zeta":1}`
	for i := 0; i < 20; i++ {
		got, err := ToJSONStr(input)
		if err != nil {
			t.Fatalf("ToJSONStr: %v", err)
		}
		if got != want {
			t.Fatalf("iteration %d ToJSONStr = %s, want %s", i, got, want)
		}
	}
}

func TestJSONRealUserPathBeanRoundTrip(t *testing.T) {
	obj, err := ParseObj(`{"user":{"name":"alice","tags":["go"]}}`)
	if err != nil {
		t.Fatalf("ParseObj: %v", err)
	}
	if err := obj.PutByPath("user.tags[1]", "json"); err != nil {
		t.Fatalf("PutByPath tags: %v", err)
	}
	if err := obj.PutByPath("user.profile.active", true); err != nil {
		t.Fatalf("PutByPath profile: %v", err)
	}

	type profile struct {
		Active bool `json:"active"`
	}
	type user struct {
		Name    string   `json:"name"`
		Tags    []string `json:"tags"`
		Profile profile  `json:"profile"`
	}
	type payload struct {
		User user `json:"user"`
	}
	var out payload
	if err := ToBean(obj, &out); err != nil {
		t.Fatalf("ToBean: %v", err)
	}
	if out.User.Name != "alice" || len(out.User.Tags) != 2 || out.User.Tags[1] != "json" || !out.User.Profile.Active {
		t.Fatalf("round trip payload = %#v", out)
	}
	encoded, err := ToJSONStr(obj)
	if err != nil {
		t.Fatalf("ToJSONStr: %v", err)
	}
	for _, want := range []string{`"name":"alice"`, `"tags":["go","json"]`, `"active":true`} {
		if !strings.Contains(encoded, want) {
			t.Fatalf("encoded %s missing %s", encoded, want)
		}
	}
}

func TestArrayAdditionalBoundaryMethods(t *testing.T) {
	cfg := NewConfig()
	cfg.IgnoreNullValue = true
	arr := NewJSONArrayWithConfig(cfg)
	arr.Add(nil).Add("kept")
	if arr.Len() != 1 || arr.GetString(0) != "kept" {
		t.Fatalf("IgnoreNullValue Add = %s", arr.String())
	}
	arr.Insert(-10, "front")
	arr.Insert(10, "back")
	if arr.GetString(0) != "front" || arr.GetString(2) != "back" {
		t.Fatalf("Insert boundaries = %s", arr.String())
	}
	if arr.Remove(-1) || arr.Remove(99) {
		t.Fatal("Remove out of range should return false")
	}
	if got := NewJSONArray().ToSlice(); got == nil || len(got) != 0 {
		t.Fatalf("empty ToSlice = %#v", got)
	}
	if arr.GetStringOr(99, "def") != "def" || arr.GetBoolOr(99, true) != true {
		t.Fatal("array getter defaults not applied")
	}
	if arr.GetJSONObject(0) != nil || arr.GetJSONArray(0) != nil {
		t.Fatal("scalar array item should not be object/array")
	}
}

func TestPathErrorAndContainerBoundaries(t *testing.T) {
	obj := NewJSONObject()
	for _, path := range []string{"", "items[abc]", "items[1"} {
		if err := obj.PutByPath(path, "x"); err == nil {
			t.Fatalf("PutByPath(%q) should fail", path)
		}
	}
	if err := PutByPath(obj, "[0]", "x"); err == nil {
		t.Fatal("object should reject root index")
	}
	arr := NewJSONArray().Add(NewJSONObject().Set("name", "first"))
	if got := arr.GetByPath("0.name"); got != "first" {
		t.Fatalf("array numeric key path = %v", got)
	}
	if err := arr.PutByPath("name", "x"); err == nil {
		t.Fatal("array should reject key path")
	}
	if err := PutByPath("not-container", "a.b", "x"); err == nil {
		t.Fatal("non-container root should fail")
	}
	if got := GetByPath(NewJSONObject().Set("leaf", "value"), "leaf.name"); got != nil {
		t.Fatalf("path through scalar = %v", got)
	}
}

func TestScalarConvertersAndConfigProviders(t *testing.T) {
	if NewConfig().Clone() == nil || (*Config)(nil).Clone() == nil {
		t.Fatal("Clone should always return config")
	}
	cfg := NewConfig()
	cfg.ParseFloatFunc = func(s string, bitSize int) (float64, error) {
		if s == "custom-float" {
			return 12.5, nil
		}
		return 0, errors.New("bad float")
	}
	cfg.ParseBoolFunc = func(s string) (bool, error) {
		if s == "custom-bool" {
			return true, nil
		}
		return false, errors.New("bad bool")
	}
	cfg.FormatFloatFunc = func(float64, byte, int, int) string { return "formatted-float" }
	obj := NewJSONObjectWithConfig(cfg).
		Set("float", "custom-float").
		Set("bool", "custom-bool").
		Set("badFloat", "nope").
		Set("badBool", "nope").
		Set("rawFloat", 1.25)
	if obj.GetFloat64("float") != 12.5 || !obj.GetBool("bool") {
		t.Fatalf("custom scalar providers failed: %s", obj.String())
	}
	if obj.GetFloat64Or("badFloat", 7.5) != 7.5 || obj.GetBoolOr("badBool", true) != true {
		t.Fatal("failed scalar parsing should return defaults")
	}
	if obj.GetString("rawFloat") != "formatted-float" {
		t.Fatalf("custom format float = %q", obj.GetString("rawFloat"))
	}
	if NewJSONObject().Set("arr", NewJSONArray().Add(1)).GetString("arr") != "[1]" {
		t.Fatal("toString should serialize nested array")
	}
}

func TestPathRootAndInvalidContainerBoundaries(t *testing.T) {
	obj := NewJSONObject().Set("name", "root")
	if got := obj.GetByPath("$"); got != obj {
		t.Fatalf("GetByPath($) = %#v, want root", got)
	}
	if err := obj.PutByPath("$", "bad"); err == nil {
		t.Fatal("PutByPath($) should reject empty resolved path")
	}
	if got := GetByPath(NewJSONArray().Add("x"), "name"); got != nil {
		t.Fatalf("array key lookup should fail, got %#v", got)
	}
	if err := obj.PutByPath("name.first", "bad"); err == nil {
		t.Fatal("path through scalar should fail")
	}
	if err := obj.PutByPath("name[0]", "bad"); err == nil {
		t.Fatal("object scalar should reject indexed continuation")
	}
}

func TestMalformedNestedJSONErrorContracts(t *testing.T) {
	for _, src := range []string{
		`{"a":`,
		`{"a":1`,
		`[1`,
		`[1,`,
		`1 true`,
	} {
		t.Run(src, func(t *testing.T) {
			if _, err := Parse(src); err == nil {
				t.Fatalf("Parse(%q) should fail", src)
			}
		})
	}
}

func TestEncodeProviderFallbacksForStructWrapping(t *testing.T) {
	type tagged struct {
		Name string `json:"name"`
	}
	if got, err := ToJSONStr(tagged{Name: "alice"}, WithUnmarshalFunc(func([]byte, any) error { return errors.New("decode failed") })); err != nil || got != `"{alice}"` {
		t.Fatalf("ToJSONStr with failing unmarshal = %s err=%v", got, err)
	}
	if got, err := ToJSONStr(tagged{Name: "bob"}, WithMarshalFunc(func(any) ([]byte, error) { return nil, errors.New("marshal failed") })); err != nil || got != `"{bob}"` {
		t.Fatalf("ToJSONStr with failing marshal = %s err=%v", got, err)
	}
	if got, err := ToJSONStr(tagged{Name: "carol"}, WithDecoderFactory(func(io.Reader) *stdjson.Decoder { return nil })); err != nil || got != `"{carol}"` {
		t.Fatalf("ToJSONStr with nil decoder = %s err=%v", got, err)
	}
}

func TestQuoteEscapesControlCharacters(t *testing.T) {
	quoted := Quote("a\rb\bc\fd\x01")
	for _, want := range []string{`\r`, `\b`, `\f`, `\u0001`} {
		if !strings.Contains(quoted, want) {
			t.Fatalf("Quote output %q missing %s", quoted, want)
		}
	}
}
