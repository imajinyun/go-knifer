package obj

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"reflect"
)

type codecConfig struct {
	newEncoder func(io.Writer) Encoder
	newDecoder func(io.Reader) Decoder
}

// Encoder is the serialization encoder contract used by object helpers.
type Encoder interface {
	Encode(any) error
}

// Decoder is the serialization decoder contract used by object helpers.
type Decoder interface {
	Decode(any) error
}

// CodecOption customizes object serialization helpers per call.
type CodecOption func(*codecConfig)

// WithEncoderFactory sets the encoder factory used by SerializeWithOptions and CloneWithOptions.
func WithEncoderFactory(factory func(io.Writer) Encoder) CodecOption {
	return func(c *codecConfig) {
		if factory != nil {
			c.newEncoder = factory
		}
	}
}

// WithDecoderFactory sets the decoder factory used by DeserializeWithOptions and CloneWithOptions.
func WithDecoderFactory(factory func(io.Reader) Decoder) CodecOption {
	return func(c *codecConfig) {
		if factory != nil {
			c.newDecoder = factory
		}
	}
}

func applyCodecOptions(opts []CodecOption) codecConfig {
	cfg := codecConfig{
		newEncoder: func(w io.Writer) Encoder { return gob.NewEncoder(w) },
		newDecoder: func(r io.Reader) Decoder { return gob.NewDecoder(r) },
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	if cfg.newEncoder == nil {
		cfg.newEncoder = func(w io.Writer) Encoder { return gob.NewEncoder(w) }
	}
	if cfg.newDecoder == nil {
		cfg.newDecoder = func(r io.Reader) Decoder { return gob.NewDecoder(r) }
	}
	return cfg
}

// Clone creates a deep copy through gob serialization.
func Clone[T any](src T) (T, error) {
	return CloneWithOptions(src)
}

// CloneWithOptions creates a deep copy using per-call codec options.
func CloneWithOptions[T any](src T, opts ...CodecOption) (T, error) {
	data, err := SerializeWithOptions(src, opts...)
	if err != nil {
		var zero T
		return zero, err
	}
	return DeserializeToWithOptions[T](data, nil, opts...)
}

// CloneIfPossible returns a cloned value when cloning succeeds, otherwise src.
func CloneIfPossible[T any](src T) T {
	return CloneIfPossibleWithOptions(src)
}

// CloneIfPossibleWithOptions returns a cloned value using per-call codec options when cloning succeeds, otherwise src.
func CloneIfPossibleWithOptions[T any](src T, opts ...CodecOption) T {
	clone, err := CloneWithOptions(src, opts...)
	if err != nil {
		return src
	}
	return clone
}

// CloneByStream creates a deep copy through gob serialization.
func CloneByStream[T any](src T) (T, error) { return CloneByStreamWithOptions(src) }

// CloneByStreamWithOptions creates a deep copy using per-call codec options.
func CloneByStreamWithOptions[T any](src T, opts ...CodecOption) (T, error) {
	return CloneWithOptions(src, opts...)
}

// Serialize encodes obj with gob.
func Serialize[T any](obj T) ([]byte, error) {
	return SerializeWithOptions(obj)
}

// SerializeWithOptions encodes obj using per-call codec options.
func SerializeWithOptions[T any](obj T, opts ...CodecOption) ([]byte, error) {
	var buf bytes.Buffer
	cfg := applyCodecOptions(opts)
	err := cfg.newEncoder(&buf).Encode(obj)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// SerializeOrNil encodes obj with gob and returns nil when encoding fails.
func SerializeOrNil[T any](obj T) []byte {
	return SerializeOrNilWithOptions(obj)
}

// SerializeOrNilWithOptions encodes obj using per-call codec options and returns nil when encoding fails.
func SerializeOrNilWithOptions[T any](obj T, opts ...CodecOption) []byte {
	data, err := SerializeWithOptions(obj, opts...)
	if err != nil {
		return nil
	}
	return data
}

// Deserialize decodes gob data into out, which must be a pointer.
//
// When acceptedTypes is not empty, the decoded object graph must contain only
// built-in container/scalar types plus values assignable to one of the accepted
// types. Accepted entries may be concrete values, pointers, or reflect.Type.
func Deserialize(data []byte, out any, acceptedTypes ...any) error {
	return DeserializeWithOptions(data, out, acceptedTypes)
}

// DeserializeWithOptions decodes data using per-call codec options.
func DeserializeWithOptions(data []byte, out any, acceptedTypes []any, opts ...CodecOption) error {
	cfg := applyCodecOptions(opts)
	if err := cfg.newDecoder(bytes.NewReader(data)).Decode(out); err != nil {
		return err
	}
	if len(acceptedTypes) == 0 {
		return nil
	}
	return ValidateAcceptedTypes(out, acceptedTypes...)
}

// DeserializeTo decodes gob data into a new value.
func DeserializeTo[T any](data []byte, acceptedTypes ...any) (T, error) {
	return DeserializeToWithOptions[T](data, acceptedTypes)
}

// DeserializeToWithOptions decodes data into a new value using per-call codec options.
func DeserializeToWithOptions[T any](data []byte, acceptedTypes []any, opts ...CodecOption) (T, error) {
	var out T
	if err := DeserializeWithOptions(data, &out, acceptedTypes, opts...); err != nil {
		var zero T
		return zero, err
	}
	return out, nil
}

// MustDeserialize decodes gob data into a new value and panics on failure.
func MustDeserialize[T any](data []byte, acceptedTypes ...any) T {
	out, err := DeserializeTo[T](data, acceptedTypes...)
	if err != nil {
		panic(err)
	}
	return out
}

// Register records a concrete type for gob interface encoding.
// It delegates to encoding/gob's process-global registry; callers that need
// isolated codecs should use Encoder/Decoder values directly instead.
func Register(value any) { gob.Register(value) }

// RegisterName records a concrete type with a custom gob name.
// It delegates to encoding/gob's process-global registry; callers that need
// isolated codecs should use Encoder/Decoder values directly instead.
func RegisterName(name string, value any) { gob.RegisterName(name, value) }

// ValidateAcceptedTypes checks whether value only contains built-in safe types
// plus values assignable to one of acceptedTypes.
func ValidateAcceptedTypes(value any, acceptedTypes ...any) error {
	if len(acceptedTypes) == 0 {
		return nil
	}
	allowed := make([]reflect.Type, 0, len(acceptedTypes))
	for _, accepted := range acceptedTypes {
		if accepted == nil {
			continue
		}
		if t, ok := accepted.(reflect.Type); ok {
			allowed = append(allowed, t)
			continue
		}
		allowed = append(allowed, reflect.TypeOf(accepted))
	}
	return validateValue(reflect.ValueOf(value), allowed, map[visit]bool{})
}

type visit struct {
	typ reflect.Type
	ptr uintptr
}

func validateValue(v reflect.Value, allowed []reflect.Type, seen map[visit]bool) error {
	if !v.IsValid() {
		return nil
	}
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		if v.Kind() == reflect.Pointer {
			key := visit{typ: v.Type(), ptr: v.Pointer()}
			if seen[key] {
				return nil
			}
			seen[key] = true
		}
		v = v.Elem()
	}
	t := v.Type()
	if isAllowedType(t, allowed) || isBuiltInAllowedKind(v.Kind()) {
		return validateChildren(v, allowed, seen)
	}
	return fmt.Errorf("serialize: decoded type %s is not accepted", t)
}

func validateChildren(v reflect.Value, allowed []reflect.Type, seen map[visit]bool) error {
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			if err := validateValue(v.Index(i), allowed, seen); err != nil {
				return err
			}
		}
	case reflect.Map:
		iter := v.MapRange()
		for iter.Next() {
			if err := validateValue(iter.Key(), allowed, seen); err != nil {
				return err
			}
			if err := validateValue(iter.Value(), allowed, seen); err != nil {
				return err
			}
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if field.CanInterface() {
				if err := validateValue(field, allowed, seen); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func isAllowedType(t reflect.Type, allowed []reflect.Type) bool {
	for _, a := range allowed {
		if t.AssignableTo(a) || reflect.PointerTo(t).AssignableTo(a) || t.AssignableTo(indirectType(a)) {
			return true
		}
	}
	return false
}

func indirectType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}

func isBuiltInAllowedKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.String,
		reflect.Array, reflect.Slice, reflect.Map:
		return true
	default:
		return false
	}
}
