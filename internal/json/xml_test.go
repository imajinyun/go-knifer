package json

import "testing"

func TestXMLToJSON(t *testing.T) {
	xmlStr := `<root><name>alice</name><age>30</age><tags>a</tags><tags>b</tags></root>`
	obj, err := XMLToJSON(xmlStr)
	if err != nil {
		t.Fatalf("xml->json: %v", err)
	}
	root := obj.GetJSONObject("root")
	if root == nil {
		t.Fatalf("missing root: %s", obj.String())
	}
	if root.GetString("name") != "alice" {
		t.Fatalf("name: %v", root.GetString("name"))
	}
	if root.GetInt("age") != 30 {
		t.Fatalf("age: %v", root.GetInt("age"))
	}
	tags := root.GetJSONArray("tags")
	if tags == nil || tags.Len() != 2 {
		t.Fatalf("tags: %v", tags)
	}
}

func TestJSONToXML(t *testing.T) {
	root := NewJSONObject().
		Set("name", "alice").
		Set("tags", NewJSONArray().Add("a").Add("b"))
	x, err := JSONToXML(root, "user")
	if err != nil {
		t.Fatalf("json->xml: %v", err)
	}
	want := "<user><name>alice</name><tags>a</tags><tags>b</tags></user>"
	if x != want {
		t.Fatalf("xml mismatch:\n got: %s\nwant: %s", x, want)
	}
}

func TestMapToJSONObjectOrdersKeysDeterministically(t *testing.T) {
	obj := mapToJSONObject(map[string]any{
		"zeta": map[string]any{"b": "bee", "a": "aye"},
		"alpha": []any{
			map[string]any{"d": "dee", "c": "see"},
		},
	})
	if got, want := obj.String(), `{"alpha":[{"c":"see","d":"dee"}],"zeta":{"a":"aye","b":"bee"}}`; got != want {
		t.Fatalf("mapToJSONObject output = %s, want %s", got, want)
	}
}

func TestJSONXMLBridgePlainValueBoundaries(t *testing.T) {
	if got := jsonValueToMap("scalar"); got["element"] != "scalar" {
		t.Fatalf("jsonValueToMap scalar = %#v", got)
	}
	if got := jsonObjectToMap(nil); got != nil {
		t.Fatalf("jsonObjectToMap(nil) = %#v", got)
	}
	arr := NewJSONArray().Add(NewJSONObject().Set("name", "alice")).Add(nil)
	plain := jsonValueToPlain(arr).([]any)
	if len(plain) != 2 || plain[1] != nil {
		t.Fatalf("jsonValueToPlain array = %#v", plain)
	}
	obj := mapXMLValue(map[string]any{"nested": map[string]any{"value": "ok"}}).(*JSONObject)
	if obj.GetJSONObject("nested").GetString("value") != "ok" {
		t.Fatalf("mapXMLValue nested = %s", obj.String())
	}
}
