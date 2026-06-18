package json

import "testing"

func TestArrayOps(t *testing.T) {
	arr := NewJSONArray()
	arr.Add(1).Add("x").Add(true).Add(nil)
	if arr.Len() != 4 || arr.GetInt(0) != 1 || arr.GetString(1) != "x" || !arr.GetBool(2) || !arr.IsNull(3) {
		t.Fatalf("array basic: %s", arr.String())
	}
	arr.Insert(1, "y")
	if arr.GetString(1) != "y" {
		t.Fatalf("insert")
	}
	arr.Remove(0)
	if arr.GetString(0) != "y" {
		t.Fatalf("remove")
	}
}

func TestArrayAddAllAndJoin(t *testing.T) {
	arr := NewJSONArray().AddAll(1, "x", true)
	if arr.Len() != 3 || arr.GetInt(0) != 1 || arr.GetString(1) != "x" || !arr.GetBool(2) {
		t.Fatalf("AddAll: %s", arr.String())
	}
	if got := arr.Join(","); got != "1,x,true" {
		t.Fatalf("Join: %q", got)
	}
	if got := NewJSONArray().Join(","); got != "" {
		t.Fatalf("Join empty: %q", got)
	}
}

func TestArrayToStringPretty(t *testing.T) {
	arr := NewJSONArray().Add(1).Add(2)
	out := arr.ToStringPretty()
	if out != "[\n    1,\n    2\n]" {
		t.Fatalf("ToStringPretty: %q", out)
	}
}

func TestArrayUnmarshalJSON(t *testing.T) {
	var arr JSONArray
	if err := arr.UnmarshalJSON([]byte(`[1, "x", true]`)); err != nil {
		t.Fatalf("UnmarshalJSON: %v", err)
	}
	if arr.Len() != 3 || arr.GetInt(0) != 1 || arr.GetString(1) != "x" || !arr.GetBool(2) {
		t.Fatalf("after UnmarshalJSON: %s", arr.String())
	}

	var bad JSONArray
	if err := bad.UnmarshalJSON([]byte(`not json`)); err == nil {
		t.Fatal("expect error on invalid JSON")
	}
}
