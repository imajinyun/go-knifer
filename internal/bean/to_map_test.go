package bean

import "testing"

func TestToMapUsesPrimaryTagAndOmit(t *testing.T) {
	got, err := ToMap(sourceProfile{Name: "alice", Age: "18", Skip: "hidden"})
	if err != nil {
		t.Fatalf("ToMap() error = %v", err)
	}
	if got["name"] != "alice" || got["age"] != "18" {
		t.Fatalf("map = %#v", got)
	}
	if _, ok := got["Skip"]; ok {
		t.Fatalf("omit field leaked: %#v", got)
	}
}
