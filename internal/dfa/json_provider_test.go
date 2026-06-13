package dfa

import "testing"

func TestAnyHelpersUseJSONProviders(t *testing.T) {
	tree := NewWordTree().AddWords("secret")
	marshalCalls := 0
	unmarshalCalls := 0
	marshal := func(any) ([]byte, error) {
		marshalCalls++
		return []byte(`{"text":"secret"}`), nil
	}
	unmarshal := func(_ []byte, dst any) error {
		unmarshalCalls++
		dst.(*struct {
			Text string `json:"text"`
		}).Text = "***"
		return nil
	}
	if !ContainsAnyWithOptions(struct{}{}, WithMatcher(tree), WithJSONMarshal(marshal)) {
		t.Fatal("ContainsAnyWithOptions should use marshal provider")
	}
	if got := GetFoundAllAnyWithOptions(struct{}{}, WithMatcher(tree), WithJSONMarshal(marshal)); len(got) != 1 || got[0].Word != "secret" {
		t.Fatalf("GetFoundAllAnyWithOptions = %#v", got)
	}
	got, err := FilterAnyWithOptions(struct {
		Text string `json:"text"`
	}{}, true, nil,
		WithMatcher(tree), WithJSONMarshal(marshal), WithJSONUnmarshal(unmarshal))
	if err != nil {
		t.Fatalf("FilterAnyWithOptions: %v", err)
	}
	if got.Text != "***" || marshalCalls != 3 || unmarshalCalls != 1 {
		t.Fatalf("providers got=%+v marshalCalls=%d unmarshalCalls=%d", got, marshalCalls, unmarshalCalls)
	}
}
