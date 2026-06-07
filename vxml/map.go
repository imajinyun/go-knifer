package vxml

import xmlimpl "github.com/imajinyun/go-knifer/internal/xml"

// XMLToMap parses XML into a nested map. Repeated sibling tags become []any.
func XMLToMap(xmlStr string) (map[string]any, error) { return xmlimpl.XMLToMap(xmlStr) }

// XMLToMapWithOptions parses XML into a nested map with parser options.
func XMLToMapWithOptions(xmlStr string, opts ...ParseOption) (map[string]any, error) {
	return xmlimpl.XMLToMapWithOptions(xmlStr, opts...)
}

// XMLNodeToMap converts an element into a nested map value.
func XMLNodeToMap(node *Element) map[string]any { return xmlimpl.XMLNodeToMap(node) }

// XMLToMapInto parses XML and merges values into result.
func XMLToMapInto(xmlStr string, result map[string]any) (map[string]any, error) {
	return xmlimpl.XMLToMapInto(xmlStr, result)
}

// XMLToMapIntoWithOptions parses XML and merges values into result with parser options.
func XMLToMapIntoWithOptions(xmlStr string, result map[string]any, opts ...ParseOption) (map[string]any, error) {
	return xmlimpl.XMLToMapIntoWithOptions(xmlStr, result, opts...)
}

// XMLNodeToMapInto converts an element to map and merges values into result.
func XMLNodeToMapInto(node *Element, result map[string]any) map[string]any {
	return xmlimpl.XMLNodeToMapInto(node, result)
}

// XMLToBean parses XML and decodes the generated map into dst.
func XMLToBean(xmlStr string, dst any) error { return xmlimpl.XMLToBean(xmlStr, dst) }

// XMLToBeanWithOptions parses XML and decodes the generated map into dst with parser options.
func XMLToBeanWithOptions(xmlStr string, dst any, opts ...ParseOption) error {
	return xmlimpl.XMLToBeanWithOptions(xmlStr, dst, opts...)
}

// XMLNodeToBean converts an element tree to a map and decodes it into dst.
func XMLNodeToBean(node *Element, dst any) error { return xmlimpl.XMLNodeToBean(node, dst) }

// XMLNodeToBeanWithOptions converts an element tree to a map and decodes it into dst with bean options.
func XMLNodeToBeanWithOptions(node *Element, dst any, opts ...BeanOption) error {
	return xmlimpl.XMLNodeToBeanWithOptions(node, dst, opts...)
}
