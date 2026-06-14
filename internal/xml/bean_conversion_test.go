package xml

import (
	"strings"
	"testing"
)

func TestBeanConversions(t *testing.T) {
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
}
