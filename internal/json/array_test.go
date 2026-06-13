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
