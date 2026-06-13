package vobj_test

import (
	"errors"
	"io"
	"math"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/imajinyun/go-knifer/vobj"
)

type record struct {
	Name string
	Tags []string
}

func TestFacadeObjectHelpers(t *testing.T) {
	if !vobj.Equal(1, int64(1)) || !vobj.Contains([]string{"go", "tool"}, "go") {
		t.Fatal("equality or contains failed")
	}
	if !vobj.IsEmpty([]int{}) || vobj.Length(map[string]int{"a": 1}) != 1 {
		t.Fatal("empty or length failed")
	}
	name := "go"
	if vobj.DefaultIfNil(&name, "x") != "go" || vobj.DefaultIfNil[string](nil, "x") != "x" {
		t.Fatal("defaults failed")
	}
	if got := vobj.Apply(&name, func(s string) int { return len(s) }); got != 2 {
		t.Fatal("apply failed")
	}
}

func TestFacadeCloneAndSerialize(t *testing.T) {
	src := record{Name: "go", Tags: []string{"tool"}}
	clone, err := vobj.Clone(src)
	if err != nil {
		t.Fatalf("Clone: %v", err)
	}
	clone.Tags[0] = "sdk"
	if src.Tags[0] != "tool" {
		t.Fatal("clone changed source")
	}
	data, err := vobj.Serialize(src)
	if err != nil {
		t.Fatalf("Serialize: %v", err)
	}
	var out record
	if err := vobj.Deserialize(data, &out); err != nil || out.Name != src.Name {
		t.Fatalf("Deserialize: %#v %v", out, err)
	}
}

func TestFacadeSerializeExtended(t *testing.T) {
	src := record{Name: "go", Tags: []string{"tool"}}

	data, err := vobj.Serialize(src)
	if err != nil {
		t.Fatalf("Serialize: %v", err)
	}

	out, err := vobj.DeserializeTo[record](data)
	if err != nil {
		t.Fatalf("DeserializeTo: %v", err)
	}
	if out.Name != src.Name || len(out.Tags) != 1 || out.Tags[0] != "tool" {
		t.Fatalf("DeserializeTo mismatch: %#v", out)
	}

	nilData := vobj.SerializeOrNil(src)
	if nilData == nil {
		t.Fatal("SerializeOrNil should not return nil for valid input")
	}
}

func TestFacadeNilDefaultAndCollectionHelpers(t *testing.T) {
	var values []int
	if !vobj.IsNil(values) || !vobj.IsNull(values) || vobj.IsNotNil(values) || vobj.IsNotNull(values) {
		t.Fatal("typed nil checks returned unexpected result")
	}
	if got := vobj.Length(42); got != -1 {
		t.Fatalf("Length unsupported = %d", got)
	}
	if !vobj.Contains("go-knifer", "knife") || vobj.Contains(map[string]int{"a": 1}, 2) {
		t.Fatal("Contains returned unexpected result")
	}
	if !vobj.Equals(1, uint(1)) || !vobj.NotEqual("a", "b") {
		t.Fatal("Equals/NotEqual returned unexpected result")
	}
	utc := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	sameInstant := time.Date(2024, 1, 1, 8, 0, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	if !vobj.Equals(utc, sameInstant) || vobj.Equals(utc, "2024-01-01T00:00:00Z") {
		t.Fatal("time Equals returned unexpected result")
	}

	name := "go"
	supplierCalls := 0
	if got := vobj.DefaultIfNilFunc(&name, func() string {
		supplierCalls++
		return "fallback"
	}); got != "go" || supplierCalls != 0 {
		t.Fatalf("DefaultIfNilFunc existing = %q calls=%d", got, supplierCalls)
	}
	if got := vobj.DefaultIfNilFunc[string](nil, func() string {
		supplierCalls++
		return "fallback"
	}); got != "fallback" || supplierCalls != 1 {
		t.Fatalf("DefaultIfNilFunc nil = %q calls=%d", got, supplierCalls)
	}
	if got := vobj.DefaultIfNilApply(&name, strings.ToUpper, "fallback"); got != "GO" {
		t.Fatalf("DefaultIfNilApply existing = %q", got)
	}
	if got := vobj.DefaultIfNilApply[string, string](nil, strings.ToUpper, "fallback"); got != "fallback" {
		t.Fatalf("DefaultIfNilApply nil = %q", got)
	}
	accepted := ""
	vobj.Accept(&name, func(s string) { accepted = s })
	vobj.Accept[string](nil, func(s string) { t.Fatalf("Accept nil called with %q", s) })
	if accepted != "go" {
		t.Fatalf("Accept captured = %q", accepted)
	}
}

func TestFacadeSerializationOptionsAndValidation(t *testing.T) {
	src := record{Name: "go", Tags: []string{"tool"}}
	clone, err := vobj.CloneByStream(src)
	if err != nil || clone.Name != src.Name {
		t.Fatalf("CloneByStream = %#v, %v", clone, err)
	}
	clone = vobj.CloneIfPossible(src)
	clone.Tags[0] = "copy"
	if src.Tags[0] != "tool" {
		t.Fatal("CloneIfPossible changed source")
	}

	failingOpt := vobj.WithEncoderFactory(func(io.Writer) vobj.Encoder {
		return encoderFunc(func(any) error { return errors.New("encode failed") })
	})
	if got := vobj.SerializeOrNilWithOptions(src, failingOpt); got != nil {
		t.Fatalf("SerializeOrNilWithOptions failing = %v", got)
	}
	if got := vobj.CloneIfPossibleWithOptions(src, failingOpt); !reflect.DeepEqual(got, src) {
		t.Fatalf("CloneIfPossibleWithOptions fallback = %#v", got)
	}

	var decoded record
	err = vobj.DeserializeWithOptions([]byte("ignored"), &decoded, nil,
		vobj.WithDecoderFactory(func(io.Reader) vobj.Decoder {
			return decoderFunc(func(out any) error {
				*out.(*record) = record{Name: "decoded", Tags: []string{"via-option"}}
				return nil
			})
		}),
	)
	if err != nil || decoded.Name != "decoded" || decoded.Tags[0] != "via-option" {
		t.Fatalf("DeserializeWithOptions = %#v, %v", decoded, err)
	}

	data, err := vobj.Serialize(src)
	if err != nil {
		t.Fatal(err)
	}
	out := vobj.MustDeserialize[record](data)
	if !reflect.DeepEqual(out, src) {
		t.Fatalf("MustDeserialize = %#v", out)
	}
	if _, err := vobj.DeserializeToWithOptions[record]([]byte("bad"), nil); err == nil {
		t.Fatal("DeserializeToWithOptions invalid data error = nil")
	}
	if err := vobj.ValidateAcceptedTypes(src, record{}); err != nil {
		t.Fatalf("ValidateAcceptedTypes accepted record: %v", err)
	}
	if err := vobj.ValidateAcceptedTypes(src, "not-record"); err == nil {
		t.Fatal("ValidateAcceptedTypes rejected type error = nil")
	}
}

func TestFacadeTypeNumberCompareAndEmptyHelpers(t *testing.T) {
	if !vobj.IsBasicType("go") || vobj.IsBasicType(record{}) {
		t.Fatal("IsBasicType returned unexpected result")
	}
	if !vobj.IsValidIfNumber(1.25) || vobj.IsValidIfNumber(math.NaN()) || vobj.IsValidIfNumber(math.Inf(1)) {
		t.Fatal("IsValidIfNumber returned unexpected result")
	}
	a, b := 1, 2
	if vobj.Compare(&a, &b) >= 0 || vobj.Compare[int](nil, &a) <= 0 || vobj.CompareNull[int](nil, &a, false) >= 0 {
		t.Fatal("Compare helpers returned unexpected ordering")
	}
	if typ := vobj.TypeOf(record{}); typ == nil || typ.Name() != "record" {
		t.Fatalf("TypeOf = %v", typ)
	}
	if got := vobj.TypeName(record{}); !strings.Contains(got, "record") {
		t.Fatalf("TypeName = %q", got)
	}
	if got := vobj.ToString(nil); got != "null" {
		t.Fatalf("ToString(nil) = %q", got)
	}
	if got := vobj.EmptyCount(nil, "", []int{}, "x"); got != 3 {
		t.Fatalf("EmptyCount = %d", got)
	}
	if !vobj.HasNil("x", (*int)(nil)) || !vobj.HasNull((*int)(nil)) || !vobj.HasEmpty("x", []int{}) {
		t.Fatal("HasNil/HasNull/HasEmpty returned unexpected result")
	}
	if !vobj.IsAllEmpty(nil, "", []int{}) || vobj.IsAllEmpty("", "x") {
		t.Fatal("IsAllEmpty returned unexpected result")
	}
	if !vobj.IsAllNotEmpty("x", []int{1}) || vobj.IsAllNotEmpty("x", []int{}) {
		t.Fatal("IsAllNotEmpty returned unexpected result")
	}
}

type encoderFunc func(any) error

func (f encoderFunc) Encode(v any) error { return f(v) }

type decoderFunc func(any) error

func (f decoderFunc) Decode(v any) error { return f(v) }
