package sets

import (
	"encoding/json"
	"testing"
)

func TestSetJSONRoundTrip(t *testing.T) {
	original := NewInt(3, 1, 2)
	b, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}

	var decoded Int
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatal(err)
	}
	if !decoded.Equal(original) {
		t.Fatalf("decoded = %v, want %v", decoded.Members(), original.Members())
	}
}

func TestGenericSetJSONRoundTrip(t *testing.T) {
	original := New("go", "knifer")
	b, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}

	var decoded Set[string]
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatal(err)
	}
	if !decoded.Equal(original) {
		t.Fatalf("decoded = %v, want %v", decoded.Members(), original.Members())
	}
}

func TestSetJSONWithOptions(t *testing.T) {
	original := New("go")
	marshalCalled := false
	b, err := original.MarshalJSONWithOptions(WithSetMarshalFunc(func(v any) ([]byte, error) {
		marshalCalled = true
		return json.Marshal(v)
	}))
	if err != nil {
		t.Fatal(err)
	}
	if !marshalCalled {
		t.Fatal("custom marshal provider was not used")
	}

	unmarshalCalled := false
	var decoded Set[string]
	if err := decoded.UnmarshalJSONWithOptions(b, WithSetUnmarshalFunc(func(data []byte, v any) error {
		unmarshalCalled = true
		return json.Unmarshal(data, v)
	})); err != nil {
		t.Fatal(err)
	}
	if !unmarshalCalled || !decoded.Equal(original) {
		t.Fatalf("unmarshalCalled=%v decoded=%v", unmarshalCalled, decoded.Members())
	}
}
