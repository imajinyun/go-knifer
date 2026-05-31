package xml

import (
	"bytes"
	"encoding/json"
	stdxml "encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	NBSP          = "&nbsp;"
	AMP           = "&amp;"
	QUOTE         = "&quot;"
	APOS          = "&apos;"
	LT            = "&lt;"
	GT            = "&gt;"
	InvalidRegex  = "[\\x00-\\x08\\x0b-\\x0c\\x0e-\\x1f]"
	CommentRegex  = "(?s)<!--.+?-->"
	IndentDefault = 2
	ContentKey    = "content"
)

var (
	namespaceAware = true
	invalidRe      = regexp.MustCompile(InvalidRegex)
	commentRe      = regexp.MustCompile(CommentRegex)
)

// Document is a lightweight XML document tree.
type Document struct {
	Root *Element
}

// Element is a lightweight XML element node.
type Element struct {
	Name     stdxml.Name
	Attr     []stdxml.Attr
	Text     string
	Children []*Element
	Parent   *Element
}

// TokenHandler consumes streaming XML tokens.
type TokenHandler func(stdxml.Token) error

// NamespaceCache stores prefix to namespace URI mappings discovered from a document.
type NamespaceCache struct {
	Default string
	Prefix  map[string]string
	URI     map[string]string
}

// DisableDefaultDocumentBuilderFactory is a no-op compatibility hook.
func DisableDefaultDocumentBuilderFactory() {}

// SetNamespaceAware records whether parsed element names should keep namespace URIs.
func SetNamespaceAware(isNamespaceAware bool) { namespaceAware = isNamespaceAware }

// ReadXML parses XML content directly, or treats the input as a file path when it does not start with '<'.
func ReadXML(pathOrContent string) (*Document, error) {
	if strings.HasPrefix(strings.TrimSpace(pathOrContent), "<") {
		return ParseXML(pathOrContent)
	}
	return ReadXMLFile(pathOrContent)
}

// ReadXMLFile parses an XML file.
func ReadXMLFile(path string) (*Document, error) {
	// #nosec G304 -- SDK file helper intentionally reads the caller-provided XML path.
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ReadXMLBytes(data)
}

// ReadXMLBytes parses XML bytes.
func ReadXMLBytes(data []byte) (*Document, error) { return ReadXMLReader(bytes.NewReader(data)) }

// ReadXMLReader parses XML from reader.
func ReadXMLReader(r io.Reader) (*Document, error) {
	dec := stdxml.NewDecoder(r)
	var stack []*Element
	var root *Element
	for {
		tok, err := dec.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case stdxml.StartElement:
			name := t.Name
			if !namespaceAware {
				name.Space = ""
			}
			ele := &Element{Name: name, Attr: append([]stdxml.Attr(nil), t.Attr...)}
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				ele.Parent = parent
				parent.Children = append(parent.Children, ele)
			} else if root == nil {
				root = ele
			}
			stack = append(stack, ele)
		case stdxml.EndElement:
			if len(stack) == 0 {
				return nil, fmt.Errorf("unexpected closing tag: %s", t.Name.Local)
			}
			stack = stack[:len(stack)-1]
		case stdxml.CharData:
			if len(stack) > 0 {
				stack[len(stack)-1].Text += string([]byte(t))
			}
		}
	}
	if root == nil {
		return nil, errors.New("xml: root element not found")
	}
	return &Document{Root: root}, nil
}

// ParseXML parses an XML string.
func ParseXML(xmlStr string) (*Document, error) { return ReadXMLReader(strings.NewReader(xmlStr)) }

// ReadBySAX streams XML tokens from reader to handler.
func ReadBySAX(r io.Reader, handler TokenHandler) error {
	if handler == nil {
		return nil
	}
	dec := stdxml.NewDecoder(r)
	for {
		tok, err := dec.Token()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if err := handler(tok); err != nil {
			return err
		}
	}
}

// ReadBySAXFile streams XML tokens from file.
func ReadBySAXFile(path string, handler TokenHandler) (err error) {
	// #nosec G304 -- SDK file helper intentionally opens the caller-provided XML path.
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := f.Close(); err == nil {
			err = closeErr
		}
	}()
	return ReadBySAX(f, handler)
}

// ToStr serializes a document or element without pretty indentation.
func ToStr(v any) string { return ToStrPretty(v, false) }

// ToStrPretty serializes a document or element and optionally pretty prints it.
func ToStrPretty(v any, pretty bool) string {
	s, _ := ToStrCharset(v, "UTF-8", pretty, false)
	return s
}

// ToStrCharset serializes a document or element with charset declaration options.
func ToStrCharset(v any, charset string, pretty bool, omitXMLDeclaration bool) (string, error) {
	var buf bytes.Buffer
	if err := Write(v, &buf, charset, indent(pretty), omitXMLDeclaration); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Format pretty prints XML content.
func Format(xmlStr string) (string, error) {
	doc, err := ParseXML(xmlStr)
	if err != nil {
		return "", err
	}
	return ToStrCharset(doc, "UTF-8", true, false)
}

// ToFile writes a document or element to path.
func ToFile(v any, path string, charset string) (err error) {
	if charset == "" {
		charset = "UTF-8"
	}
	f, err := os.Create(path) // #nosec G304 -- SDK file helper intentionally creates the caller-provided XML path.
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := f.Close(); err == nil {
			err = closeErr
		}
	}()
	return Write(v, f, charset, IndentDefault, false)
}

// WriteObjectAsXML writes a struct, map, or scalar value as XML to path.
func WriteObjectAsXML(path string, bean any) error { return ToFile(BeanToXML(bean), path, "UTF-8") }

// Write serializes a document or element to writer.
func Write(v any, w io.Writer, charset string, indentSize int, omitXMLDeclaration bool) error {
	if w == nil {
		return errors.New("xml: nil writer")
	}
	if charset == "" {
		charset = "UTF-8"
	}
	if !omitXMLDeclaration {
		if _, err := fmt.Fprintf(w, `<?xml version="1.0" encoding="%s"?>`, charset); err != nil {
			return err
		}
		if indentSize > 0 {
			if _, err := io.WriteString(w, "\n"); err != nil {
				return err
			}
		}
	}
	ele := elementOf(v)
	if ele == nil {
		return errors.New("xml: unsupported node")
	}
	return writeElement(w, ele, indentSize, 0)
}

// Transform copies XML from source to result with optional pretty formatting.
func Transform(source io.Reader, result io.Writer, charset string, indentSize int, omitXMLDeclaration bool) error {
	doc, err := ReadXMLReader(source)
	if err != nil {
		return err
	}
	return Write(doc, result, charset, indentSize, omitXMLDeclaration)
}

// CreateXML creates an empty XML document.
func CreateXML() *Document { return &Document{} }

// CreateXMLWithRoot creates an XML document with root element.
func CreateXMLWithRoot(rootElementName string) *Document {
	return CreateXMLWithRootNS(rootElementName, "")
}

// CreateXMLWithRootNS creates an XML document with root element and namespace URI.
func CreateXMLWithRootNS(rootElementName, namespace string) *Document {
	root := &Element{Name: stdxml.Name{Local: rootElementName}}
	if namespace != "" {
		root.Attr = append(root.Attr, stdxml.Attr{Name: stdxml.Name{Local: "xmlns"}, Value: namespace})
	}
	return &Document{Root: root}
}

// CreateDocumentBuilder is a no-op compatibility hook.
func CreateDocumentBuilder() struct{} { return struct{}{} }

// CreateDocumentBuilderFactory is a no-op compatibility hook.
func CreateDocumentBuilderFactory() struct{} { return struct{}{} }

// GetRootElement returns the document root element.
func GetRootElement(doc *Document) *Element {
	if doc == nil {
		return nil
	}
	return doc.Root
}

// GetOwnerDocument returns the document that owns node by walking to the root.
func GetOwnerDocument(node *Element) *Document {
	if node == nil {
		return nil
	}
	for node.Parent != nil {
		node = node.Parent
	}
	return &Document{Root: node}
}

// CleanInvalid removes XML 1.0 invalid control characters.
func CleanInvalid(xmlContent string) string { return invalidRe.ReplaceAllString(xmlContent, "") }

// CleanComment removes XML comments.
func CleanComment(xmlContent string) string { return commentRe.ReplaceAllString(xmlContent, "") }

// GetElements returns child elements with tag name. Empty tagName returns all direct children.
func GetElements(element *Element, tagName string) []*Element {
	if element == nil {
		return nil
	}
	out := make([]*Element, 0)
	for _, child := range element.Children {
		if tagName == "" || child.Name.Local == tagName {
			out = append(out, child)
		}
	}
	return out
}

// GetElement returns the first child element with tag name.
func GetElement(element *Element, tagName string) *Element {
	children := GetElements(element, tagName)
	if len(children) == 0 {
		return nil
	}
	return children[0]
}

// ElementText returns child text or defaultValue when missing.
func ElementText(element *Element, tagName string, defaultValue ...string) string {
	child := GetElement(element, tagName)
	if child != nil {
		return strings.TrimSpace(child.Text)
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

// TransElements returns the input list without nil elements.
func TransElements(nodes []*Element) []*Element {
	out := make([]*Element, 0, len(nodes))
	for _, node := range nodes {
		if node != nil {
			out = append(out, node)
		}
	}
	return out
}

// CreateXPath is a compatibility hook for simple path expressions.
func CreateXPath() struct{} { return struct{}{} }

// GetElementByXPath returns the first element matched by a simple XPath-like expression.
func GetElementByXPath(expression string, source any) *Element {
	if node := GetNodeByXPath(expression, source); node != nil {
		return node
	}
	return nil
}

// GetNodeListByXPath returns all elements matched by a simple XPath-like expression.
func GetNodeListByXPath(expression string, source any) []*Element {
	root := elementOf(source)
	if root == nil {
		return nil
	}
	return findByPath(root, expression)
}

// GetNodeByXPath returns the first node matched by a simple XPath-like expression.
func GetNodeByXPath(expression string, source any) *Element {
	nodes := GetNodeListByXPath(expression, source)
	if len(nodes) == 0 {
		return nil
	}
	return nodes[0]
}

// GetByXPath returns matched text, element, or list based on returnType: string, node, or nodes.
func GetByXPath(expression string, source any, returnType string) any {
	switch strings.ToLower(returnType) {
	case "string", "text":
		if node := GetNodeByXPath(expression, source); node != nil {
			return strings.TrimSpace(node.Text)
		}
		return ""
	case "nodes", "nodelist", "list":
		return GetNodeListByXPath(expression, source)
	default:
		return GetNodeByXPath(expression, source)
	}
}

// Escape escapes XML text.
func Escape(s string) string {
	var buf bytes.Buffer
	if err := stdxml.EscapeText(&buf, []byte(s)); err != nil {
		return s
	}
	return buf.String()
}

// Unescape unescapes XML/HTML entities.
func Unescape(s string) string { return html.UnescapeString(s) }

// XMLToMap parses XML into a nested map. Repeated sibling tags become []any.
func XMLToMap(xmlStr string) (map[string]any, error) {
	doc, err := ParseXML(xmlStr)
	if err != nil {
		return nil, err
	}
	result := map[string]any{}
	if doc.Root != nil {
		addMapValue(result, doc.Root.Name.Local, elementToValue(doc.Root))
	}
	return result, nil
}

// XMLNodeToMap converts an element into a nested map value.
func XMLNodeToMap(node *Element) map[string]any {
	result := map[string]any{}
	if node != nil {
		addMapValue(result, node.Name.Local, elementToValue(node))
	}
	return result
}

// XMLToBean parses XML and decodes the generated map into dst.
func XMLToBean(xmlStr string, dst any) error {
	m, err := XMLToMap(xmlStr)
	if err != nil {
		return err
	}
	return mapToBean(m, dst)
}

// XMLNodeToBean converts an element tree to a map and decodes it into dst.
func XMLNodeToBean(node *Element, dst any) error { return mapToBean(XMLNodeToMap(node), dst) }

func mapToBean(m map[string]any, dst any) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

// XMLToMapInto parses XML and merges values into result.
func XMLToMapInto(xmlStr string, result map[string]any) (map[string]any, error) {
	m, err := XMLToMap(xmlStr)
	if err != nil {
		return result, err
	}
	if result == nil {
		result = map[string]any{}
	}
	for k, v := range m {
		result[k] = v
	}
	return result, nil
}

// XMLNodeToMapInto converts an element to map and merges values into result.
func XMLNodeToMapInto(node *Element, result map[string]any) map[string]any {
	if result == nil {
		result = map[string]any{}
	}
	for k, v := range XMLNodeToMap(node) {
		result[k] = v
	}
	return result
}

// MapToXMLStr serializes map data to XML string. Empty rootName emits each top-level key as a tag.
func MapToXMLStr(data map[string]any, rootName ...string) (string, error) {
	root := ""
	if len(rootName) > 0 {
		root = rootName[0]
	}
	return MapToXMLStrOptions(data, root, "", true, false, "UTF-8")
}

// MapToXMLStrOptions serializes map data with namespace, pretty, declaration, and charset options.
func MapToXMLStrOptions(data map[string]any, rootName, namespace string, pretty, omitXMLDeclaration bool, charset string) (string, error) {
	doc := MapToXML(data, rootName, namespace)
	return ToStrCharset(doc, charset, pretty, omitXMLDeclaration)
}

// MapToXML converts map data to a document.
func MapToXML(data map[string]any, rootName string, namespace ...string) *Document {
	ns := ""
	if len(namespace) > 0 {
		ns = namespace[0]
	}
	if rootName != "" {
		doc := CreateXMLWithRootNS(rootName, ns)
		Append(doc.Root, data)
		return doc
	}
	doc := CreateXMLWithRootNS("xml", ns)
	Append(doc.Root, data)
	return doc
}

// BeanToXML converts a struct or map-like value to a document.
func BeanToXML(bean any, namespace ...string) *Document {
	return BeanToXMLWithOptions(bean, first(namespace), false)
}

// BeanToXMLWithOptions converts a struct or map-like value to a document with namespace and nil-field handling.
func BeanToXMLWithOptions(bean any, namespace string, ignoreNull bool) *Document {
	name := typeName(bean)
	if name == "" {
		name = "xml"
	}
	m := structToMap(bean, false, ignoreNull)
	return MapToXML(m, name, namespace)
}

// IsElement reports whether node is not nil.
func IsElement(node *Element) bool { return node != nil }

// AppendChild appends and returns a child element.
func AppendChild(node *Element, tagName string, namespace ...string) *Element {
	if node == nil {
		return nil
	}
	ns := ""
	if len(namespace) > 0 {
		ns = namespace[0]
	}
	child := &Element{Name: stdxml.Name{Local: tagName}, Parent: node}
	if ns != "" {
		child.Attr = append(child.Attr, stdxml.Attr{Name: stdxml.Name{Local: "xmlns"}, Value: ns})
	}
	node.Children = append(node.Children, child)
	return child
}

// AppendText appends text to an element.
func AppendText(node *Element, text any) *Element {
	if node != nil && text != nil {
		node.Text += fmt.Sprint(text)
	}
	return node
}

// Append appends map, slice, struct, or scalar data to node.
func Append(node *Element, data any) {
	if node == nil || data == nil {
		return
	}
	appendValue(node, data)
}

// NewNamespaceCache collects namespace declarations from doc.
func NewNamespaceCache(doc *Document) *NamespaceCache {
	cache := &NamespaceCache{Prefix: map[string]string{}, URI: map[string]string{}}
	walk(GetRootElement(doc), func(ele *Element) {
		for _, attr := range ele.Attr {
			if attr.Name.Local == "xmlns" && attr.Name.Space == "" {
				cache.Default = attr.Value
				cache.Prefix["DEFAULT"] = attr.Value
				cache.URI[attr.Value] = "DEFAULT"
				continue
			}
			if attr.Name.Space == "xmlns" {
				cache.Prefix[attr.Name.Local] = attr.Value
				cache.URI[attr.Value] = attr.Name.Local
			}
		}
	})
	return cache
}

// NamespaceURI returns namespace URI for prefix.
func (c *NamespaceCache) NamespaceURI(prefix string) string {
	if c == nil {
		return ""
	}
	if prefix == "" {
		return c.Default
	}
	return c.Prefix[prefix]
}

// PrefixOf returns one prefix for namespace URI.
func (c *NamespaceCache) PrefixOf(uri string) string {
	if c == nil {
		return ""
	}
	return c.URI[uri]
}

func elementOf(v any) *Element {
	switch x := v.(type) {
	case *Document:
		if x == nil {
			return nil
		}
		return x.Root
	case *Element:
		return x
	default:
		return nil
	}
}

func indent(pretty bool) int {
	if pretty {
		return IndentDefault
	}
	return 0
}

func writeElement(w io.Writer, ele *Element, indentSize, level int) error {
	if ele == nil {
		return nil
	}
	if indentSize > 0 {
		if _, err := io.WriteString(w, strings.Repeat(" ", indentSize*level)); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(w, "<"); err != nil {
		return err
	}
	if err := writeName(w, ele.Name); err != nil {
		return err
	}
	for _, attr := range ele.Attr {
		if _, err := io.WriteString(w, " "); err != nil {
			return err
		}
		if err := writeName(w, attr.Name); err != nil {
			return err
		}
		if _, err := io.WriteString(w, `="`); err != nil {
			return err
		}
		if _, err := io.WriteString(w, Escape(attr.Value)); err != nil {
			return err
		}
		if _, err := io.WriteString(w, `"`); err != nil {
			return err
		}
	}
	if len(ele.Children) == 0 && strings.TrimSpace(ele.Text) == "" {
		_, err := io.WriteString(w, "/>")
		return err
	}
	if _, err := io.WriteString(w, ">"); err != nil {
		return err
	}
	text := strings.TrimSpace(ele.Text)
	if text != "" {
		if _, err := io.WriteString(w, Escape(text)); err != nil {
			return err
		}
	}
	if len(ele.Children) > 0 {
		if indentSize > 0 {
			if _, err := io.WriteString(w, "\n"); err != nil {
				return err
			}
		}
		for i, child := range ele.Children {
			if err := writeElement(w, child, indentSize, level+1); err != nil {
				return err
			}
			if indentSize > 0 && i < len(ele.Children)-1 {
				if _, err := io.WriteString(w, "\n"); err != nil {
					return err
				}
			}
		}
		if indentSize > 0 {
			if _, err := io.WriteString(w, "\n"); err != nil {
				return err
			}
			if _, err := io.WriteString(w, strings.Repeat(" ", indentSize*level)); err != nil {
				return err
			}
		}
	}
	if _, err := io.WriteString(w, "</"); err != nil {
		return err
	}
	if err := writeName(w, ele.Name); err != nil {
		return err
	}
	_, err := io.WriteString(w, ">")
	return err
}

func writeName(w io.Writer, name stdxml.Name) error {
	_, err := io.WriteString(w, name.Local)
	return err
}

func elementToValue(ele *Element) any {
	obj := map[string]any{}
	for _, attr := range ele.Attr {
		obj[attr.Name.Local] = parseScalar(attr.Value)
	}
	for _, child := range ele.Children {
		addMapValue(obj, child.Name.Local, elementToValue(child))
	}
	text := strings.TrimSpace(ele.Text)
	if len(obj) == 0 {
		if text == "" {
			return ""
		}
		return parseScalar(text)
	}
	if text != "" {
		obj[ContentKey] = parseScalar(text)
	}
	return obj
}

func addMapValue(m map[string]any, key string, value any) {
	if old, ok := m[key]; ok {
		if arr, ok := old.([]any); ok {
			m[key] = append(arr, value)
		} else {
			m[key] = []any{old, value}
		}
		return
	}
	m[key] = value
}

func parseScalar(s string) any {
	s = strings.TrimSpace(s)
	switch strings.ToLower(s) {
	case "true":
		return true
	case "false":
		return false
	case "null":
		return nil
	}
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return s
}

func appendValue(node *Element, data any) {
	if m, ok := normalizeMap(data); ok {
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, key := range keys {
			appendNamedValue(node, key, m[key])
		}
		return
	}
	rv := reflect.ValueOf(data)
	if rv.IsValid() && (rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array) {
		for i := 0; i < rv.Len(); i++ {
			appendNamedValue(node, "element", rv.Index(i).Interface())
		}
		return
	}
	AppendText(node, data)
}

func appendNamedValue(parent *Element, key string, value any) {
	rv := reflect.ValueOf(value)
	if rv.IsValid() && (rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array) {
		for i := 0; i < rv.Len(); i++ {
			appendNamedValue(parent, key, rv.Index(i).Interface())
		}
		return
	}
	child := AppendChild(parent, key)
	if value == nil {
		return
	}
	if m, ok := normalizeMap(value); ok {
		appendValue(child, m)
		return
	}
	if isStruct(value) {
		appendValue(child, structToMap(value, true, false))
		return
	}
	AppendText(child, value)
}

func normalizeMap(data any) (map[string]any, bool) {
	if data == nil {
		return nil, false
	}
	switch m := data.(type) {
	case map[string]any:
		return m, true
	case map[string]string:
		out := map[string]any{}
		for k, v := range m {
			out[k] = v
		}
		return out, true
	case map[any]any:
		out := map[string]any{}
		for k, v := range m {
			out[fmt.Sprint(k)] = v
		}
		return out, true
	}
	rv := reflect.ValueOf(data)
	if rv.Kind() == reflect.Map {
		out := map[string]any{}
		iter := rv.MapRange()
		for iter.Next() {
			out[fmt.Sprint(iter.Key().Interface())] = iter.Value().Interface()
		}
		return out, true
	}
	return nil, false
}

func isStruct(data any) bool {
	if data == nil {
		return false
	}
	rv := reflect.ValueOf(data)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	return rv.IsValid() && rv.Kind() == reflect.Struct
}

func structToMap(data any, honorXMLName bool, ignoreNull bool) map[string]any {
	out := map[string]any{}
	rv := reflect.ValueOf(data)
	if !rv.IsValid() {
		return out
	}
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return out
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return out
	}
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if field.PkgPath != "" {
			continue
		}
		name := field.Name
		if tag := field.Tag.Get("xml"); tag != "" {
			name = strings.Split(tag, ",")[0]
			if name == "-" {
				continue
			}
		}
		if honorXMLName && name == "XMLName" {
			continue
		}
		if name == "" || name == "XMLName" {
			continue
		}
		value := rv.Field(i).Interface()
		if ignoreNull && isNilValue(value) {
			continue
		}
		out[name] = value
	}
	return out
}

func isNilValue(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}

func first(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func typeName(data any) string {
	if data == nil {
		return ""
	}
	t := reflect.TypeOf(data)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() == reflect.Struct {
		return strings.ToLower(t.Name())
	}
	return "xml"
}

func findByPath(root *Element, expr string) []*Element {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return nil
	}
	if strings.HasPrefix(expr, "//") {
		name := strings.TrimPrefix(expr, "//")
		var out []*Element
		walk(root, func(ele *Element) {
			if ele.Name.Local == name {
				out = append(out, ele)
			}
		})
		return out
	}
	parts := strings.Split(strings.Trim(expr, "/"), "/")
	if len(parts) == 0 {
		return nil
	}
	current := []*Element{root}
	if parts[0] == root.Name.Local {
		parts = parts[1:]
	}
	for _, part := range parts {
		if part == "" {
			continue
		}
		next := make([]*Element, 0)
		for _, ele := range current {
			for _, child := range ele.Children {
				if child.Name.Local == part {
					next = append(next, child)
				}
			}
		}
		current = next
	}
	return current
}

func walk(ele *Element, fn func(*Element)) {
	if ele == nil {
		return
	}
	fn(ele)
	for _, child := range ele.Children {
		walk(child, fn)
	}
}
