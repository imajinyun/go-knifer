package xml

import (
	"strings"
	"testing"
)

func TestXMLErrorContract(t *testing.T) {
	_, err := ParseXML("")
	assertXMLInvalidInput(t, err)

	assertXMLInvalidInput(t, WriteTo(nil, CreateXMLWithRoot("root")))
	assertXMLInvalidInput(t, WriteTo(&strings.Builder{}, "unsupported"))

	var dst struct {
		Root struct {
			Value int `json:"value"`
		} `json:"root"`
	}
	assertXMLInvalidInput(t, XMLToBean(`<root><value>not-int</value></root>`, &dst))
}
