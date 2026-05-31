package vxml

import (
	"strings"
	"testing"
)

func TestFacadeXMLUtilities(t *testing.T) {
	doc, err := ParseXML(`<root><name>alice</name></root>`)
	if err != nil {
		t.Fatalf("ParseXML failed: %v", err)
	}
	if ElementText(GetRootElement(doc), "name") != "alice" {
		t.Fatal("ElementText facade failed")
	}
	AppendChild(GetRootElement(doc), "age")
	AppendText(GetElement(GetRootElement(doc), "age"), 30)
	out, err := ToStrCharset(doc, "UTF-8", false, true)
	if err != nil || !strings.Contains(out, `<age>30</age>`) {
		t.Fatalf("ToStrCharset facade = %q, %v", out, err)
	}
	m, err := XMLToMap(out)
	if err != nil || m["root"] == nil {
		t.Fatalf("XMLToMap facade = %#v, %v", m, err)
	}
	back, err := MapToXMLStrOptions(map[string]any{"name": "bob"}, "user", "", false, true, "UTF-8")
	if err != nil || back != `<user><name>bob</name></user>` {
		t.Fatalf("MapToXMLStrOptions facade = %q, %v", back, err)
	}
	if Escape("<x>") != "&lt;x&gt;" || Unescape("&lt;x&gt;") != "<x>" {
		t.Fatal("escape facade failed")
	}
}
