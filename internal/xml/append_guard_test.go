package xml

import (
	"strings"
	"testing"
)

func TestAppendAndGuards(t *testing.T) {
	if CreateXML().Root != nil || IsElement(nil) || GetRootElement(nil) != nil || GetOwnerDocument(nil) != nil {
		t.Fatal("guard helpers failed")
	}
	doc := CreateXMLWithRoot("root")
	Append(doc.Root, map[string]any{"items": []int{1, 2}, "nested": map[string]any{"ok": true}})
	got, err := MarshalString(doc, WithOmitDeclaration(true))
	if err != nil || !strings.Contains(got, `<items>1</items><items>2</items>`) || !strings.Contains(got, `<ok>true</ok>`) {
		t.Fatalf("Append map/slice serialized = %q, %v", got, err)
	}
	sliceDoc := CreateXMLWithRoot("root")
	Append(sliceDoc.Root, []string{"a", "b"})
	sliceStr, err := MarshalString(sliceDoc, WithOmitDeclaration(true))
	if err != nil || sliceStr != `<root><element>a</element><element>b</element></root>` {
		t.Fatalf("Append slice serialized = %q, %v", sliceStr, err)
	}
	nestedStructDoc := CreateXMLWithRoot("root")
	Append(nestedStructDoc.Root, map[string]any{"user": sampleBean{Name: "struct", Age: 3}})
	nestedStructStr, err := MarshalString(nestedStructDoc, WithOmitDeclaration(true))
	if err != nil || !strings.Contains(nestedStructStr, `<name>struct</name>`) || !strings.Contains(nestedStructStr, `<age>3</age>`) {
		t.Fatalf("Append nested struct serialized = %q, %v", nestedStructStr, err)
	}
	structDoc := CreateXMLWithRoot("root")
	Append(structDoc.Root, sampleBean{Name: "struct", Age: 3})
	structStr, err := MarshalString(structDoc, WithOmitDeclaration(true))
	if err != nil || structStr != `<root>{struct 3 &lt;nil&gt;}</root>` {
		t.Fatalf("Append root scalar struct serialized = %q, %v", structStr, err)
	}
	if AppendChild(nil, "x") != nil || AppendText(nil, "x") != nil {
		t.Fatal("nil append helpers should return nil")
	}
	if got := AppendText(CreateXMLWithRoot("r").Root, nil); got == nil || got.Text != "" {
		t.Fatalf("AppendText nil = %#v", got)
	}
	if child := AppendChild(doc.Root, "withNS", "urn:child"); child == nil || child.Attr[0].Value != "urn:child" {
		t.Fatalf("AppendChild namespace = %#v", child)
	}
	if got := ElementText(doc.Root, "missing", "default"); got != "default" {
		t.Fatalf("ElementText default = %q", got)
	}
	if got := ElementText(doc.Root, "missing"); got != "" {
		t.Fatalf("ElementText missing = %q", got)
	}
	if got := TransElements([]*Element{nil, doc.Root}); len(got) != 1 || got[0] != doc.Root {
		t.Fatalf("TransElements = %#v", got)
	}
}
