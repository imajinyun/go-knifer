package vdfa

import "testing"

func TestFacadeWordTree(t *testing.T) {
	tree := NewWordTree().AddWords("foo", "foobar")
	got := tree.MatchAllMode("foo foobar", -1, true, true)
	if len(got) != 3 || got[0] != "foo" || got[1] != "foo" || got[2] != "foobar" {
		t.Fatalf("MatchAllMode() = %#v", got)
	}
}

func TestFacadePackageMatcher(t *testing.T) {
	Init([]string{"bad", "badword"})
	if !Contains("a badword") {
		t.Fatalf("expected package matcher to contain word")
	}
	got := Filter("a badword")
	if got != "a *******" {
		t.Fatalf("Filter() = %q", got)
	}
}

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
