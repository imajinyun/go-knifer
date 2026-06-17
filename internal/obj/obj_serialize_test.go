package obj

import (
	"errors"
	"io"
	"reflect"
	"testing"
)

func TestCloneSerializeAndDeserialize(t *testing.T) {
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
}

func TestCloneConvenienceAndFallbacks(t *testing.T) {
	src := sample{Name: "n", Tags: []string{"a"}}
	clone := CloneIfPossible(src)
	clone.Tags[0] = "b"
	if src.Tags[0] != "a" {
		t.Fatal("CloneIfPossible should return an independent clone on success")
	}
	if got := CloneIfPossibleWithOptions(src, WithEncoderFactory(func(io.Writer) Encoder { return failingEncoder{} })); !reflect.DeepEqual(got, src) {
		t.Fatalf("CloneIfPossibleWithOptions fallback = %#v", got)
	}
	streamClone, err := CloneByStream(src)
	if err != nil || !reflect.DeepEqual(streamClone, src) {
		t.Fatalf("CloneByStream = %#v, %v", streamClone, err)
	}
	streamClone, err = CloneByStreamWithOptions(src)
	if err != nil || !reflect.DeepEqual(streamClone, src) {
		t.Fatalf("CloneByStreamWithOptions = %#v, %v", streamClone, err)
	}
}

func TestSerializeDeserializeConvenienceAndErrors(t *testing.T) {
	src := sample{Name: "n", Tags: []string{"a"}}
	if data := SerializeOrNil(src); len(data) == 0 {
		t.Fatal("SerializeOrNil should return gob bytes on success")
	}
	if data := SerializeOrNilWithOptions(src, WithEncoderFactory(func(io.Writer) Encoder { return failingEncoder{} })); data != nil {
		t.Fatalf("SerializeOrNilWithOptions failure = %#v", data)
	}
	data, err := Serialize(src)
	if err != nil {
		t.Fatalf("Serialize: %v", err)
	}
	out, err := DeserializeTo[sample](data)
	if err != nil || !reflect.DeepEqual(out, src) {
		t.Fatalf("DeserializeTo = %#v, %v", out, err)
	}
	out, err = DeserializeToWithOptions[sample](data, nil)
	if err != nil || !reflect.DeepEqual(out, src) {
		t.Fatalf("DeserializeToWithOptions = %#v, %v", out, err)
	}
	if got := MustDeserialize[sample](data); !reflect.DeepEqual(got, src) {
		t.Fatalf("MustDeserialize = %#v", got)
	}
	mustPanic(t, func() { _ = MustDeserialize[sample]([]byte("bad gob")) })
	if _, err := SerializeWithOptions(src, WithEncoderFactory(func(io.Writer) Encoder { return failingEncoder{} })); !errors.Is(err, errCodecFailure) {
		t.Fatalf("SerializeWithOptions failing encoder err = %v", err)
	}
	if err := DeserializeWithOptions([]byte("bad gob"), &out, nil); err == nil {
		t.Fatal("DeserializeWithOptions should return decoder errors")
	}
}

func TestValidateAcceptedTypesBoundaries(t *testing.T) {
	type trusted struct {
		Name   string
		Nested *trusted
		List   []any
		Meta   map[string]any
	}
	root := &trusted{Name: "root"}
	root.Nested = root
	root.List = []any{1, "ok", []string{"nested"}}
	root.Meta = map[string]any{"self": root, "numbers": []int{1, 2}, "time_like": timeLike{Value: "x"}}
	if err := ValidateAcceptedTypes(root, trusted{}, reflect.TypeOf(timeLike{}), nil); err != nil {
		t.Fatalf("ValidateAcceptedTypes accepted graph: %v", err)
	}
	if err := ValidateAcceptedTypes(nil, trusted{}); err != nil {
		t.Fatalf("ValidateAcceptedTypes nil: %v", err)
	}
	if err := ValidateAcceptedTypes(root); err != nil {
		t.Fatalf("ValidateAcceptedTypes no accepted types: %v", err)
	}
	bad := map[string]any{"bad": timeLike{}}
	if err := ValidateAcceptedTypes(bad, trusted{}); err == nil {
		t.Fatal("ValidateAcceptedTypes should reject unaccepted struct values")
	}
	data, err := Serialize(sample{Name: "bad"})
	if err != nil {
		t.Fatalf("Serialize bad graph: %v", err)
	}
	var out sample
	if err := Deserialize(data, &out, trusted{}); err == nil {
		t.Fatal("Deserialize should validate accepted decoded graph")
	}
}

func TestRegisterHelpers(t *testing.T) {
	Register(registeredValue{})
	RegisterName("objNamedRegisteredValue", namedRegisteredValue{})
	data, err := Serialize(map[string]any{
		"value": registeredValue{Value: "x"},
		"named": namedRegisteredValue{Value: "y"},
	})
	if err != nil {
		t.Fatalf("Serialize registered interface value: %v", err)
	}
	var out map[string]any
	if err := Deserialize(data, &out, registeredValue{}, namedRegisteredValue{}); err != nil {
		t.Fatalf("Deserialize registered interface value: %v", err)
	}
	if got, ok := out["value"].(registeredValue); !ok || got.Value != "x" {
		t.Fatalf("decoded registered value = %#v", out["value"])
	}
	if got, ok := out["named"].(namedRegisteredValue); !ok || got.Value != "y" {
		t.Fatalf("decoded named registered value = %#v", out["named"])
	}
}
