package dfa

import (
	"reflect"
	"testing"
)

func TestStopRunes(t *testing.T) {
	tree := NewWordTree().AddWord("tio")
	got := tree.MatchAll("AAAAAAAt-ioBBBBBBB")
	want := []string{"t-io"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("MatchAll() = %#v, want %#v", got, want)
	}
}

func TestFilterDoesNotMatchInsideDigits(t *testing.T) {
	Init([]string{"12宝宝龙", "34皮卡丘"})
	text := "creator_user_id=2000907612345839744"
	if got := Filter(text); got != text {
		t.Fatalf("Filter() = %q, want %q", got, text)
	}
}
