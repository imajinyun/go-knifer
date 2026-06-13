package ref

import (
	"reflect"
	"testing"
)

func TestNewInstanceAndConstructorHelpers(t *testing.T) {
	ctor := GetConstructor(newSample)
	if !ctor.IsValid() || len(GetConstructors(newSample)) != 1 || len(GetConstructorsDirectly(newSample)) != 1 {
		t.Fatal("constructor helpers failed")
	}
	created, err := NewInstance(newSample, "alice", 20)
	if err != nil || created.(sample).Name != "alice" || created.(sample).Age != 20 {
		t.Fatalf("NewInstance constructor = %#v, %v", created, err)
	}
	zero, err := NewInstance(reflect.TypeOf(sample{}))
	if err != nil || zero.(sample).Name != "" {
		t.Fatalf("NewInstance type = %#v, %v", zero, err)
	}
	ptr, err := NewInstance(reflect.TypeOf(&sample{}))
	if err != nil || reflect.TypeOf(ptr).String() != "*ref.sample" {
		t.Fatalf("NewInstance pointer type = %#v, %v", ptr, err)
	}
	if NewInstanceIfPossible(nil) != nil {
		t.Fatal("NewInstanceIfPossible nil should be nil")
	}
	if SetAccessible(1) != 1 {
		t.Fatal("SetAccessible should return input")
	}
	RemoveFinalModify(nil)
}
