package vjson_test

import (
	"strings"
	"testing"

	"github.com/imajinyun/go-knifer/vjson"
)

func TestFacadeConversionNamesWithoutJSONPrefix(t *testing.T) {
	type user struct {
		Name string `json:"name"`
	}

	var u user
	if err := vjson.ToBean(`{"name":"go-knifer"}`, &u); err != nil {
		t.Fatal(err)
	}
	if u.Name != "go-knifer" {
		t.Fatalf("ToBean() name = %q", u.Name)
	}

	var list []user
	if err := vjson.ToList(`[{"name":"go"},{"name":"tool"}]`, &list); err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 || list[1].Name != "tool" {
		t.Fatalf("ToList() = %#v", list)
	}

	xmlObj, err := vjson.XMLToJSON(`<user><name>go-knifer</name></user>`)
	if err != nil {
		t.Fatal(err)
	}
	if got := objectString(xmlObj, "user", "name"); got != "go-knifer" {
		t.Fatalf("XMLToJSON() user.name = %q", got)
	}

	xmlStr, err := vjson.ToXML(vjson.NewObject().Set("name", "go-knifer"), "user")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(xmlStr, "<name>go-knifer</name>") {
		t.Fatalf("ToXML() = %q", xmlStr)
	}
}

func objectString(obj *vjson.Object, objectKey, valueKey string) string {
	return obj.GetJSONObject(objectKey).GetString(valueKey)
}
