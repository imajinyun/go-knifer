package obj

import (
	"encoding/gob"
	"io"
	"math"
	"reflect"
	"testing"
	"time"
)

type recordingEncoder struct {
	inner  Encoder
	called *bool
}

func (e recordingEncoder) Encode(v any) error {
	*e.called = true
	return e.inner.Encode(v)
}

type recordingDecoder struct {
	inner  Decoder
	called *bool
}

func (d recordingDecoder) Decode(v any) error {
	*d.called = true
	return d.inner.Decode(v)
}

type sample struct {
	Name string
	Tags []string
}

func TestEqualLengthContainsAndEmpty(t *testing.T) {
	if !Equal(1, int64(1)) || NotEqual("a", "a") {
		t.Fatal("numeric or string equality failed")
	}
	utc := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	sameInstant := time.Date(2024, 1, 1, 8, 0, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	if !Equals(utc, sameInstant) {
		t.Fatal("time equality should compare instants")
	}
	if Equals(utc, "2024-01-01T00:00:00Z") {
		t.Fatal("time equality should reject non-time values")
	}
	if Length([]int{1, 2, 3}) != 3 || Length(10) != -1 {
		t.Fatal("length failed")
	}
	if !Contains([]int{1, 2, 3}, int64(2)) || !Contains("hello", "ell") {
		t.Fatal("contains failed")
	}
	if !IsEmpty(map[string]int{}) || IsEmpty(1) || !IsNotEmpty([]int{1}) {
		t.Fatal("empty checks failed")
	}
}

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

func TestCloneSerializeCompareAndType(t *testing.T) {
	src := sample{Name: "n", Tags: []string{"a"}}
	clone, err := Clone(src)
	if err != nil {
		t.Fatalf("Clone: %v", err)
	}
	clone.Tags[0] = "b"
	if src.Tags[0] != "a" {
		t.Fatal("clone is not independent")
	}
	data, err := Serialize(src)
	if err != nil {
		t.Fatalf("Serialize: %v", err)
	}
	var out sample
	if err := Deserialize(data, &out); err != nil || !reflect.DeepEqual(out, src) {
		t.Fatalf("Deserialize: %#v %v", out, err)
	}
	a, b := 1, 2
	if Compare(&a, &b) >= 0 || CompareNull[int](nil, &b, true) <= 0 {
		t.Fatal("compare failed")
	}
	if TypeName(src) == "" || ToString(nil) != "null" {
		t.Fatal("type or string failed")
	}
}

func TestSerializeWithOptionsUsesCodecFactories(t *testing.T) {
	src := sample{Name: "n", Tags: []string{"a"}}
	encoderCalled := false
	data, err := SerializeWithOptions(src, WithEncoderFactory(func(w io.Writer) Encoder {
		return recordingEncoder{inner: gob.NewEncoder(w), called: &encoderCalled}
	}))
	if err != nil {
		t.Fatalf("SerializeWithOptions: %v", err)
	}
	if !encoderCalled {
		t.Fatal("custom encoder factory was not used")
	}

	decoderCalled := false
	var out sample
	err = DeserializeWithOptions(data, &out, nil, WithDecoderFactory(func(r io.Reader) Decoder {
		return recordingDecoder{inner: gob.NewDecoder(r), called: &decoderCalled}
	}))
	if err != nil {
		t.Fatalf("DeserializeWithOptions: %v", err)
	}
	if !decoderCalled || !reflect.DeepEqual(out, src) {
		t.Fatalf("decoderCalled=%v out=%#v", decoderCalled, out)
	}

	clone, err := CloneWithOptions(src,
		WithEncoderFactory(func(w io.Writer) Encoder { return gob.NewEncoder(w) }),
		WithDecoderFactory(func(r io.Reader) Decoder { return gob.NewDecoder(r) }),
	)
	if err != nil || !reflect.DeepEqual(clone, src) {
		t.Fatalf("CloneWithOptions = %#v, %v", clone, err)
	}
}

func TestBasicAndValidNumber(t *testing.T) {
	if !IsBasicType("x") || IsBasicType(sample{}) {
		t.Fatal("basic type check failed")
	}
	if !IsValidIfNumber(1) || IsValidIfNumber(math.NaN()) || IsValidIfNumber(math.Inf(1)) {
		t.Fatal("valid number check failed")
	}
}
