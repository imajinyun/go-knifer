package vxml

import (
	"os"
	"strings"
	"testing"
)

func TestFacadeParseAndReadOptions(t *testing.T) {
	doc, err := ParseXML(`<root xmlns:p="urn:p"><p:a>1</p:a></root>`, WithNamespaceAware(false))
	if err != nil {
		t.Fatalf("ParseXML facade failed: %v", err)
	}
	child := GetElement(GetRootElement(doc), "a")
	if child == nil || child.Name.Space != "" {
		t.Fatalf("namespace option facade not applied: %#v", child)
	}

	doc, err = ReadXMLBytes([]byte(`<root><b>2</b></root>`))
	if err != nil || ElementText(doc.Root, "b") != "2" {
		t.Fatalf("ReadXMLBytes facade doc=%#v err=%v", doc, err)
	}
	doc, err = ReadXMLReader(strings.NewReader(`<root><c>3</c></root>`))
	if err != nil || ElementText(doc.Root, "c") != "3" {
		t.Fatalf("ReadXMLReader facade doc=%#v err=%v", doc, err)
	}
	tmp := t.TempDir() + "/in.xml"
	if err := os.WriteFile(tmp, []byte(`<root><d>4</d></root>`), 0o600); err != nil {
		t.Fatal(err)
	}
	doc, err = ReadXMLFile(tmp)
	if err != nil || ElementText(doc.Root, "d") != "4" {
		t.Fatalf("ReadXMLFile facade doc=%#v err=%v", doc, err)
	}
	doc, err = ReadXML(tmp)
	if err != nil || ElementText(doc.Root, "d") != "4" {
		t.Fatalf("ReadXML path facade doc=%#v err=%v", doc, err)
	}
}
