package vdfa

import (
	"errors"
	"strings"
	"testing"
)

func TestFacadeFilterAny(t *testing.T) {
	type payload struct {
		Text string `json:"text"`
	}
	InitString("secret", DefaultSeparator)
	got, err := FilterAny(payload{Text: "a secret"}, true, nil)
	if err != nil {
		t.Fatalf("FilterAny() error = %v", err)
	}
	if got.Text != "a ******" {
		t.Fatalf("FilterAny() = %#v", got)
	}
}

func TestFacadeAnyHelpersWithJSONOptions(t *testing.T) {
	tree := NewWordTree().AddWords("secret", "token")
	marshal := func(any) ([]byte, error) {
		return []byte(`{"text":"secret token"}`), nil
	}

	if !ContainsAnyWithOptions(struct{}{}, WithMatcher(tree), WithJSONMarshal(marshal)) {
		t.Fatal("ContainsAnyWithOptions should use custom marshal output")
	}
	first, ok := GetFoundFirstAnyWithOptions(struct{}{}, WithMatcher(tree), WithJSONMarshal(marshal))
	if !ok || first.Word != "secret" || first.String() != "secret" {
		t.Fatalf("GetFoundFirstAnyWithOptions = %#v, %v", first, ok)
	}
	all := GetFoundAllAnyWithOptions(struct{}{}, WithMatcher(tree), WithJSONMarshal(marshal))
	if len(all) != 2 || all[0].Word != "secret" || all[1].Word != "token" {
		t.Fatalf("GetFoundAllAnyWithOptions = %#v", all)
	}

	var filteredJSON string
	got, err := FilterAnyWithOptions(struct{}{}, true, func(word FoundWord) string {
		return strings.ToUpper(word.FoundWord)
	}, WithMatcher(tree), WithJSONMarshal(marshal), WithJSONUnmarshal(func(data []byte, out any) error {
		filteredJSON = string(data)
		return nil
	}))
	if err != nil {
		t.Fatalf("FilterAnyWithOptions custom JSON: %v", err)
	}
	if got != (struct{}{}) || !strings.Contains(filteredJSON, "SECRET") || !strings.Contains(filteredJSON, "TOKEN") {
		t.Fatalf("filtered result = %#v json=%q", got, filteredJSON)
	}
}

func TestFacadeAnyHelpersMarshalErrors(t *testing.T) {
	marshalErr := errors.New("marshal failed")
	if ContainsAnyWithOptions(struct{}{}, WithJSONMarshal(func(any) ([]byte, error) { return nil, marshalErr })) {
		t.Fatal("ContainsAnyWithOptions should report false when custom marshal fails")
	}
	if found, ok := GetFoundFirstAnyWithOptions(struct{}{}, WithJSONMarshal(func(any) ([]byte, error) { return nil, marshalErr })); ok || found != (FoundWord{}) {
		t.Fatalf("GetFoundFirstAnyWithOptions = %#v, %v", found, ok)
	}
	if got := GetFoundAllAnyWithOptions(struct{}{}, WithJSONMarshal(func(any) ([]byte, error) { return nil, marshalErr })); len(got) != 0 {
		t.Fatalf("GetFoundAllAnyWithOptions = %#v", got)
	}
	if _, err := FilterAnyWithOptions(struct{}{}, true, nil, WithJSONMarshal(func(any) ([]byte, error) { return nil, marshalErr })); !errors.Is(err, marshalErr) {
		t.Fatalf("FilterAnyWithOptions marshal error = %v, want %v", err, marshalErr)
	}
}

func TestFacadeFilterAnyStringAndUnmarshalError(t *testing.T) {
	got, err := FilterAnyWithOptions("a secret", true, nil, WithMatcherWords([]string{"secret"}))
	if err != nil {
		t.Fatalf("FilterAnyWithOptions string: %v", err)
	}
	if got != "a ******" {
		t.Fatalf("FilterAnyWithOptions string = %q", got)
	}

	unmarshalErr := errors.New("unmarshal failed")
	_, err = FilterAnyWithOptions(struct{ Text string }{Text: "secret"}, true, nil,
		WithMatcherWords([]string{"secret"}),
		WithJSONUnmarshal(func([]byte, any) error { return unmarshalErr }),
	)
	if !errors.Is(err, unmarshalErr) {
		t.Fatalf("FilterAnyWithOptions unmarshal error = %v, want %v", err, unmarshalErr)
	}
}
