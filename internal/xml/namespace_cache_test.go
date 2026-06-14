package xml

import "testing"

func TestNamespaceCache(t *testing.T) {
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
