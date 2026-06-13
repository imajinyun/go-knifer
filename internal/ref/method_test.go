package ref

import (
	"reflect"
	"testing"
)

func TestMethodHelpers(t *testing.T) {
	s := &sample{Name: "alice"}
	if names := GetPublicMethodNames(s); !containsString(names, "GetName") || !containsString(names, "SetName") {
		t.Fatalf("method names = %v", names)
	}
	if methods := GetPublicMethods(s, func(m reflect.Method) bool { return m.Name == "Add" }); len(methods) != 1 {
		t.Fatalf("filtered methods = %v", methods)
	}
	if _, ok := GetPublicMethod(s, "Add", reflect.TypeOf(1), reflect.TypeOf(2)); !ok {
		t.Fatal("GetPublicMethod failed")
	}
	if _, ok := GetMethodOfObj(s, "Add", 1, 2); !ok {
		t.Fatal("GetMethodOfObj failed")
	}
	if _, ok := GetMethodIgnoreCase(s, "getname"); !ok {
		t.Fatal("GetMethodIgnoreCase failed")
	}
	if _, ok := GetMethodByName(s, "String"); !ok {
		t.Fatal("GetMethodByName failed")
	}
	if _, ok := GetMethodByNameIgnoreCase(s, "string"); !ok {
		t.Fatal("GetMethodByNameIgnoreCase failed")
	}
	if len(GetMethods(s)) == 0 || len(GetMethodsDirectly(s, true, true)) == 0 {
		t.Fatal("GetMethods failed")
	}
	stringMethod, _ := GetMethodByName(s, "String")
	equalMethod, _ := GetMethodByName(s, "Equal")
	hashMethod, _ := GetMethodByName(s, "HashCode")
	getMethod, _ := GetMethodByName(s, "GetName")
	setMethod, _ := GetMethodByName(s, "SetName")
	if !IsToStringMethod(stringMethod) || !IsEqualsMethod(equalMethod) || !IsHashCodeMethod(hashMethod) || !IsEmptyParam(getMethod) {
		t.Fatal("method classification failed")
	}
	if !IsGetterOrSetter(getMethod, false) || !IsGetterOrSetterIgnoreCase(setMethod) {
		t.Fatal("getter/setter classification failed")
	}
}
