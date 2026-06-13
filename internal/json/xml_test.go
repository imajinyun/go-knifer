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
