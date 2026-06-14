package ref

import (
	"context"
	"reflect"
	"testing"
)

func TestInterfaceImplementationPredicates(t *testing.T) {
	if !ImplementsError(reflect.TypeOf(sampleError{})) || ImplementsError(nil) || ImplementsError(reflect.TypeOf("value")) {
		t.Fatal("ImplementsError returned unexpected result")
	}
	if !ImplementsContext(reflect.TypeOf(context.Background())) || ImplementsContext(nil) || ImplementsContext(reflect.TypeOf(sampleError{})) {
		t.Fatal("ImplementsContext returned unexpected result")
	}
}
