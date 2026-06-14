package ref

import (
	"reflect"
	"testing"
)

func TestGetPublicFieldNames(t *testing.T) {
	if got := GetPublicFieldNames(sample{}); !reflect.DeepEqual(got, []string{"Name", "Age"}) {
		t.Fatalf("GetPublicFieldNames = %#v", got)
	}
	if got := GetPublicFieldNames((*sample)(nil)); !reflect.DeepEqual(got, []string{"Name", "Age"}) {
		t.Fatalf("GetPublicFieldNames pointer = %#v", got)
	}
	if got := GetPublicFieldNames(123); got != nil {
		t.Fatalf("GetPublicFieldNames non-struct = %#v", got)
	}
}
