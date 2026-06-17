package obj

import (
	"encoding/gob"
	"errors"
	"io"
	"reflect"
	"testing"
)

var errCodecFailure = errors.New("codec failure")

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

type failingEncoder struct{}

func (failingEncoder) Encode(any) error { return errCodecFailure }

type failingDecoder struct{}

func (failingDecoder) Decode(any) error { return errCodecFailure }

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

func TestCodecOptionsIgnoreNilFactoriesAndCanFail(t *testing.T) {
	src := sample{Name: "n", Tags: []string{"a"}}
	data, err := SerializeWithOptions(src, WithEncoderFactory(nil), nil)
	if err != nil || len(data) == 0 {
		t.Fatalf("SerializeWithOptions nil factories = len %d, %v", len(data), err)
	}
	var out sample
	if err := DeserializeWithOptions(data, &out, nil, WithDecoderFactory(nil), nil); err != nil || !reflect.DeepEqual(out, src) {
		t.Fatalf("DeserializeWithOptions nil factories = %#v, %v", out, err)
	}
	if _, err := CloneWithOptions(src, WithEncoderFactory(func(io.Writer) Encoder { return failingEncoder{} })); !errors.Is(err, errCodecFailure) {
		t.Fatalf("CloneWithOptions encoder failure = %v", err)
	}
	if _, err := CloneWithOptions(src, WithDecoderFactory(func(io.Reader) Decoder { return failingDecoder{} })); !errors.Is(err, errCodecFailure) {
		t.Fatalf("CloneWithOptions decoder failure = %v", err)
	}
}
