package vref

import (
	"reflect"
	"testing"
)

func TestFacadeMethodLookupAndClassifierHelpers(t *testing.T) {
	target := facadeMethodSample{}
	names := GetPublicMethodNames(target)
	if !reflect.DeepEqual(names, []string{"Equal", "HashCode", "SetName", "String"}) {
		t.Fatalf("GetPublicMethodNames = %#v", names)
	}
	methods := GetPublicMethods(target, func(method reflect.Method) bool { return method.Name == "String" })
	if len(methods) != 1 || methods[0].Name != "String" {
		t.Fatalf("GetPublicMethods filtered = %#v", methods)
	}
	if method, ok := GetPublicMethod(target, "SetName", reflect.TypeOf("name")); !ok || method.Name != "SetName" {
		t.Fatalf("GetPublicMethod SetName = %q, %v", method.Name, ok)
	}
	if method, ok := GetMethodIgnoreCase(target, "hashcode"); !ok || method.Name != "HashCode" {
		t.Fatalf("GetMethodIgnoreCase = %q, %v", method.Name, ok)
	}
	if method, ok := GetMethodByName(target, "String"); !ok || !IsToStringMethod(method) || !IsEmptyParam(method) {
		t.Fatalf("GetMethodByName String = %q, %v", method.Name, ok)
	}
	if method, ok := GetMethodByNameIgnoreCase(target, "equal"); !ok || !IsEqualsMethod(method) {
		t.Fatalf("GetMethodByNameIgnoreCase Equal = %q, %v", method.Name, ok)
	}
	if method, ok := GetMethodOfObj(target, "SetName", "bob"); !ok || !IsGetterOrSetter(method, false) || !IsGetterOrSetterIgnoreCase(method) {
		t.Fatalf("GetMethodOfObj SetName = %q, %v", method.Name, ok)
	}
	if method, ok := GetMethodByName(target, "HashCode"); !ok || !IsHashCodeMethod(method) {
		t.Fatalf("GetMethodByName HashCode = %q, %v", method.Name, ok)
	}
	if got := GetMethodNames(target); !reflect.DeepEqual(got, names) {
		t.Fatalf("GetMethodNames = %#v", got)
	}
	if got := GetMethodsDirectly(target, true, true); len(got) != len(names) {
		t.Fatalf("GetMethodsDirectly len = %d, want %d", len(got), len(names))
	}
}
