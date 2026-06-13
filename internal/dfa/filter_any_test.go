package dfa

import "testing"

func TestFilterAnyWithMatcherOptions(t *testing.T) {
	type payload struct {
		Text string `json:"text"`
	}
	Init([]string{"global"})
	got, err := FilterAnyWithOptions(payload{Text: "a local"}, true, nil, WithMatcherWords([]string{"local"}))
	if err != nil {
		t.Fatalf("FilterAnyWithOptions: %v", err)
	}
	if got.Text != "a *****" || Contains("local") {
		t.Fatalf("FilterAnyWithOptions = %#v globalContainsLocal=%v", got, Contains("local"))
	}
}

func TestFilterAny(t *testing.T) {
	type payload struct {
		Text string `json:"text"`
		Num  int    `json:"num"`
	}
	Init([]string{"大", "大土豆", "土豆", "刚出锅", "出锅"})
	got, err := FilterAny(payload{Text: sampleText, Num: 100}, true, nil)
	if err != nil {
		t.Fatalf("FilterAny() error = %v", err)
	}
	if got.Text != "我有一颗$****，***的" || got.Num != 100 {
		t.Fatalf("FilterAny() = %#v", got)
	}
}
