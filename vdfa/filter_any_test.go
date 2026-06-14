package vdfa

import "testing"

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
