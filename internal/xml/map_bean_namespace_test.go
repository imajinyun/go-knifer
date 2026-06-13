package xml

import (
	"errors"
	"strings"
	"testing"
)

func TestMapBeanAndNamespaceConversions(t *testing.T) {
	m, err := XMLToMap(`<root enabled="true"><name>alice</name><age>30</age><score>3.5</score><none>null</none><tags>a</tags><tags>b</tags></root>`)
	if err != nil {
		t.Fatalf("XMLToMap failed: %v", err)
	}
	root := m["root"].(map[string]any)
	if root["enabled"] != true || root["name"] != "alice" || root["age"] != int64(30) || root["score"] != 3.5 || root["none"] != nil {
		t.Fatalf("XMLToMap root = %#v", root)
	}
	if tags, ok := root["tags"].([]any); !ok || len(tags) != 2 {
		t.Fatalf("XMLToMap tags = %#v", root["tags"])
	}
	merged, err := XMLToMapInto(`<x><a>1</a></x>`, map[string]any{"old": true})
	if err != nil || merged["old"] != true || merged["x"] == nil {
		t.Fatalf("XMLToMapInto = %#v, %v", merged, err)
	}
	stripped, err := XMLToMapWithOptions(`<root xmlns:p="urn:p"><p:a>1</p:a></root>`, WithNamespaceAware(false))
	if err != nil || stripped["root"].(map[string]any)["a"] != int64(1) {
		t.Fatalf("XMLToMapWithOptions = %#v, %v", stripped, err)
	}
	limited, err := XMLToMapIntoWithOptions(`<root><a>1</a></root>`, nil, WithMaxBytes(6))
	if err == nil || limited != nil {
		t.Fatalf("XMLToMapIntoWithOptions should fail with max bytes: %#v, %v", limited, err)
	}
	customMap, err := XMLToMapWithOptions(`<root><n>custom-int</n><f>custom-float</f></root>`,
		WithScalarIntParser(func(s string, base, bitSize int) (int64, error) {
			if s == "custom-int" {
				return 99, nil
			}
			return 0, errors.New("not int")
		}),
		WithScalarFloatParser(func(s string, bitSize int) (float64, error) {
			if s == "custom-float" {
				return 6.25, nil
			}
			return 0, errors.New("not float")
		}),
	)
	if err != nil {
		t.Fatalf("XMLToMapWithOptions scalar providers: %v", err)
	}
	customRoot := customMap["root"].(map[string]any)
	if customRoot["n"] != int64(99) || customRoot["f"] != 6.25 {
		t.Fatalf("custom scalar parsers = %#v", customRoot)
	}
	if got := XMLNodeToMapInto(nil, nil); len(got) != 0 {
		t.Fatalf("XMLNodeToMapInto nil = %#v", got)
	}

	xmlStr, err := MarshalMap(map[string]any{"name": "alice", "tags": []string{"a", "b"}}, WithRootName("user"), WithOmitDeclaration(true))
	if err != nil || xmlStr != `<user><name>alice</name><tags>a</tags><tags>b</tags></user>` {
		t.Fatalf("MarshalMap = %q, %v", xmlStr, err)
	}
	beanStr, err := MarshalBean(sampleBean{Name: "bob", Age: 20}, WithRootName("sample"), WithIgnoreNullFields(true), WithOmitDeclaration(true))
	if err != nil || strings.Contains(beanStr, "empty") || !strings.Contains(beanStr, `<name>bob</name>`) || !strings.Contains(beanStr, `<age>20</age>`) {
		t.Fatalf("MarshalBean = %q, %v", beanStr, err)
	}
	defaultRoot, err := MarshalBean(sampleBean{Name: "typed"}, WithOmitDeclaration(true), WithIgnoreNullFields(true))
	if err != nil || !strings.HasPrefix(defaultRoot, `<samplebean>`) {
		t.Fatalf("MarshalBean default root = %q, %v", defaultRoot, err)
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
	var decodedNode struct {
		Root struct {
			Name string `json:"name"`
		} `json:"root"`
	}
	doc, err := ParseXML(`<root><name>node</name></root>`)
	if err != nil {
		t.Fatal(err)
	}
	if err := XMLNodeToBean(doc.Root, &decodedNode); err != nil || decodedNode.Root.Name != "node" {
		t.Fatalf("XMLNodeToBean decoded=%+v err=%v", decodedNode, err)
	}
	customCalled := false
	var custom struct{ Name string }
	if err := XMLNodeToBeanWithOptions(doc.Root, &custom, WithBeanUnmarshalFunc(func(_ []byte, dst any) error {
		customCalled = true
		dst.(*struct{ Name string }).Name = "custom"
		return nil
	})); err != nil || !customCalled || custom.Name != "custom" {
		t.Fatalf("XMLNodeToBeanWithOptions custom=%+v called=%v err=%v", custom, customCalled, err)
	}

	nsDoc, err := ParseXML(`<root xmlns="urn:default" xmlns:p="urn:p"><p:a>1</p:a></root>`)
	if err != nil {
		t.Fatal(err)
	}
	cache := NewNamespaceCache(nsDoc)
	if cache.NamespaceURI("") != "urn:default" || cache.NamespaceURI("DEFAULT") != "urn:default" || cache.NamespaceURI("p") != "urn:p" || cache.PrefixOf("urn:p") != "p" {
		t.Fatalf("namespace cache = %#v", cache)
	}
	if (*NamespaceCache)(nil).NamespaceURI("p") != "" || (*NamespaceCache)(nil).PrefixOf("urn:p") != "" {
		t.Fatal("nil namespace cache should return empty values")
	}
}
