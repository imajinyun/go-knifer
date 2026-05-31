package xml

import (
	stdxml "encoding/xml"
	"os"
	"reflect"
	"strings"
	"testing"
)

type sampleBean struct {
	Name  string  `xml:"name" json:"name"`
	Age   int     `xml:"age" json:"age"`
	Empty *string `xml:"empty" json:"empty"`
}

func TestReadCreateWriteAndFormat(t *testing.T) {
	doc, err := ParseXML(`<root><name id="1">alice</name><tags>a</tags><tags>b</tags></root>`)
	if err != nil {
		t.Fatalf("ParseXML failed: %v", err)
	}
	root := GetRootElement(doc)
	if root == nil || root.Name.Local != "root" {
		t.Fatalf("root mismatch: %#v", root)
	}
	if got := ElementText(root, "name"); got != "alice" {
		t.Fatalf("ElementText = %q", got)
	}
	if got := GetElements(root, "tags"); len(got) != 2 {
		t.Fatalf("GetElements tags length = %d", len(got))
	}
	if GetOwnerDocument(GetElement(root, "name")).Root != root {
		t.Fatal("GetOwnerDocument should walk to root")
	}
	plain, err := ToStrCharset(doc, "UTF-8", false, true)
	if err != nil || plain != `<root><name id="1">alice</name><tags>a</tags><tags>b</tags></root>` {
		t.Fatalf("ToStrCharset = %q, %v", plain, err)
	}
	formatted, err := Format(plain)
	if err != nil || !strings.Contains(formatted, "\n  <name") {
		t.Fatalf("Format = %q, %v", formatted, err)
	}
	created := CreateXMLWithRootNS("user", "urn:test")
	AppendChild(created.Root, "name")
	AppendText(GetElement(created.Root, "name"), "bob")
	createdStr, err := ToStrCharset(created, "UTF-8", false, true)
	if err != nil || !strings.Contains(createdStr, `xmlns="urn:test"`) || !strings.Contains(createdStr, `<name>bob</name>`) {
		t.Fatalf("CreateXMLWithRootNS serialized = %q, %v", createdStr, err)
	}
}

func TestCleanEscapeSAXXPathAndFile(t *testing.T) {
	if CleanInvalid("a\x00b\x08c") != "abc" {
		t.Fatal("CleanInvalid failed")
	}
	if CleanComment("<a><!-- hidden --><b/></a>") != "<a><b/></a>" {
		t.Fatal("CleanComment failed")
	}
	if Escape(`<a&"'>`) != "&lt;a&amp;&#34;&#39;&gt;" || Unescape("&lt;a&amp;&gt;") != "<a&>" {
		t.Fatal("escape/unescape failed")
	}
	var starts []string
	if err := ReadBySAX(strings.NewReader(`<root><a>1</a></root>`), func(tok stdxml.Token) error {
		if start, ok := tok.(stdxml.StartElement); ok {
			starts = append(starts, start.Name.Local)
		}
		return nil
	}); err != nil || !reflect.DeepEqual(starts, []string{"root", "a"}) {
		t.Fatalf("ReadBySAX starts=%v err=%v", starts, err)
	}
	tmp := t.TempDir() + "/x.xml"
	if err := os.WriteFile(tmp, []byte(`<root><a>1</a><a>2</a></root>`), 0o600); err != nil {
		t.Fatal(err)
	}
	doc, err := ReadXML(tmp)
	if err != nil {
		t.Fatalf("ReadXML file failed: %v", err)
	}
	if got := GetByXPath("/root/a", doc, "string"); got != "1" {
		t.Fatalf("GetByXPath string = %v", got)
	}
	if got := GetNodeListByXPath("//a", doc); len(got) != 2 {
		t.Fatalf("GetNodeListByXPath = %d", len(got))
	}
	var out strings.Builder
	if err := Transform(strings.NewReader(`<root><a>1</a></root>`), &out, "UTF-8", 0, true); err != nil || out.String() != `<root><a>1</a></root>` {
		t.Fatalf("Transform = %q, %v", out.String(), err)
	}
}

func TestMapBeanAndNamespaceConversions(t *testing.T) {
	m, err := XMLToMap(`<root><name>alice</name><age>30</age><tags>a</tags><tags>b</tags></root>`)
	if err != nil {
		t.Fatalf("XMLToMap failed: %v", err)
	}
	root := m["root"].(map[string]any)
	if root["name"] != "alice" || root["age"] != int64(30) {
		t.Fatalf("XMLToMap root = %#v", root)
	}
	if tags, ok := root["tags"].([]any); !ok || len(tags) != 2 {
		t.Fatalf("XMLToMap tags = %#v", root["tags"])
	}
	merged, err := XMLToMapInto(`<x><a>1</a></x>`, map[string]any{"old": true})
	if err != nil || merged["old"] != true || merged["x"] == nil {
		t.Fatalf("XMLToMapInto = %#v, %v", merged, err)
	}
	xmlStr, err := MapToXMLStrOptions(map[string]any{"name": "alice", "tags": []string{"a", "b"}}, "user", "", false, true, "UTF-8")
	if err != nil || xmlStr != `<user><name>alice</name><tags>a</tags><tags>b</tags></user>` {
		t.Fatalf("MapToXMLStrOptions = %q, %v", xmlStr, err)
	}
	doc := BeanToXMLWithOptions(sampleBean{Name: "bob", Age: 20}, "", true)
	beanStr, err := ToStrCharset(doc, "UTF-8", false, true)
	if err != nil || strings.Contains(beanStr, "empty") || !strings.Contains(beanStr, `<name>bob</name>`) {
		t.Fatalf("BeanToXMLWithOptions = %q, %v", beanStr, err)
	}
	var decoded struct {
		Root struct {
			Name string `json:"name"`
			Age  int64  `json:"age"`
		} `json:"root"`
	}
	if err := XMLToBean(`<root><name>alice</name><age>30</age></root>`, &decoded); err != nil || decoded.Root.Name != "alice" || decoded.Root.Age != 30 {
		t.Fatalf("XMLToBean decoded=%+v err=%v", decoded, err)
	}
	nsDoc, err := ParseXML(`<root xmlns="urn:default" xmlns:p="urn:p"><p:a>1</p:a></root>`)
	if err != nil {
		t.Fatal(err)
	}
	cache := NewNamespaceCache(nsDoc)
	if cache.NamespaceURI("DEFAULT") != "urn:default" || cache.NamespaceURI("p") != "urn:p" || cache.PrefixOf("urn:p") != "p" {
		t.Fatalf("namespace cache = %#v", cache)
	}
}

func TestParseXMLWithOptionsDoesNotMutateGlobalNamespaceSetting(t *testing.T) {
	namespaceAware = true
	doc, err := ParseXMLWithOptions(`<root xmlns:p="urn:p"><p:a>1</p:a></root>`, WithNamespaceAware(false))
	if err != nil {
		t.Fatalf("ParseXMLWithOptions failed: %v", err)
	}
	child := GetElement(doc.Root, "a")
	if child == nil || child.Name.Space != "" {
		t.Fatalf("per-call namespace option not applied: %#v", child)
	}

	doc, err = ParseXML(`<root xmlns:p="urn:p"><p:a>1</p:a></root>`)
	if err != nil {
		t.Fatalf("ParseXML failed: %v", err)
	}
	child = GetElement(doc.Root, "a")
	if child == nil || child.Name.Space != "urn:p" {
		t.Fatalf("global namespace setting should remain aware: %#v", child)
	}
}

func TestAppendAndGuards(t *testing.T) {
	if CreateXML().Root != nil || IsElement(nil) || GetRootElement(nil) != nil || GetOwnerDocument(nil) != nil {
		t.Fatal("guard helpers failed")
	}
	doc := CreateXMLWithRoot("root")
	Append(doc.Root, map[string]any{"items": []int{1, 2}, "nested": map[string]any{"ok": true}})
	got, err := ToStrCharset(doc, "UTF-8", false, true)
	if err != nil || !strings.Contains(got, `<items>1</items><items>2</items>`) || !strings.Contains(got, `<ok>true</ok>`) {
		t.Fatalf("Append map/slice serialized = %q, %v", got, err)
	}
	if AppendChild(nil, "x") != nil || AppendText(nil, "x") != nil {
		t.Fatal("nil append helpers should return nil")
	}
	if got := ElementText(doc.Root, "missing", "default"); got != "default" {
		t.Fatalf("ElementText default = %q", got)
	}
	if got := TransElements([]*Element{nil, doc.Root}); len(got) != 1 || got[0] != doc.Root {
		t.Fatalf("TransElements = %#v", got)
	}
}
