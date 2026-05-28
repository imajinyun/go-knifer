package json

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

// 默认在 XML 序列化时使用的 content 键，与 hutool 保持一致。
const xmlContentKey = "content"

// XMLToJSON 将 XML 字符串解析为 JSONObject，对应 hutool XML.toJSONObject。
// - 元素属性写入对应 key；
// - 元素文本写入 content key；
// - 同名元素自动合并为数组。
func XMLToJSON(xmlStr string) (*JSONObject, error) {
	dec := xml.NewDecoder(strings.NewReader(xmlStr))
	root := NewJSONObject()
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		if start, ok := tok.(xml.StartElement); ok {
			child, err := readXMLElement(dec, start)
			if err != nil {
				return nil, err
			}
			addToObject(root, start.Name.Local, child)
		}
	}
	return root, nil
}

// readXMLElement 解析单个元素及其子元素，返回 *JSONObject 或 字符串值。
func readXMLElement(dec *xml.Decoder, start xml.StartElement) (any, error) {
	obj := NewJSONObject()
	// 属性
	for _, attr := range start.Attr {
		obj.Set(attr.Name.Local, parseXMLScalar(attr.Value))
	}
	var textBuf strings.Builder
	for {
		tok, err := dec.Token()
		if err != nil {
			return nil, WrapJSONError(err, "xml: read failed")
		}
		switch t := tok.(type) {
		case xml.StartElement:
			child, err := readXMLElement(dec, t)
			if err != nil {
				return nil, err
			}
			addToObject(obj, t.Name.Local, child)
		case xml.EndElement:
			text := strings.TrimSpace(textBuf.String())
			if obj.Len() == 0 {
				if text == "" {
					return "", nil
				}
				return parseXMLScalar(text), nil
			}
			if text != "" {
				obj.Set(xmlContentKey, parseXMLScalar(text))
			}
			return obj, nil
		case xml.CharData:
			textBuf.Write(t)
		}
	}
}

// addToObject 同 key 自动合并为数组。
func addToObject(obj *JSONObject, key string, value any) {
	if exist, ok := obj.Get(key); ok {
		if arr, ok := exist.(*JSONArray); ok {
			arr.Add(value)
			return
		}
		arr := NewJSONArray()
		arr.Add(exist)
		arr.Add(value)
		obj.Set(key, arr)
		return
	}
	obj.Set(key, value)
}

// parseXMLScalar 尝试将 XML 文本转为 bool/int/float，否则保留字符串。
func parseXMLScalar(s string) any {
	switch strings.ToLower(s) {
	case "true":
		return true
	case "false":
		return false
	case "null":
		return Null
	}
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return s
}

// JSONToXML 将 JSONObject/JSONArray/基础值序列化为 XML。
// 当 root 是 *JSONObject 时，根标签可由 rootTag 指定（为空时不包裹根标签，
// 直接以对象的每个 key 作为顶层标签，匹配 hutool 默认行为）。
func JSONToXML(root any, rootTag string) (string, error) {
	var buf bytes.Buffer
	if rootTag != "" {
		buf.WriteString("<")
		buf.WriteString(rootTag)
		buf.WriteString(">")
	}
	if err := writeXMLValue(&buf, root); err != nil {
		return "", err
	}
	if rootTag != "" {
		buf.WriteString("</")
		buf.WriteString(rootTag)
		buf.WriteString(">")
	}
	return buf.String(), nil
}

// writeXMLValue 将 JSON 值写入 XML，依据 hutool 行为：
// - JSONObject：每个 key 作为标签包裹。
// - JSONArray：每个元素重复同一标签（由调用方传入）。
// 这里对顶层数组用 <element> 作为默认标签。
func writeXMLValue(buf *bytes.Buffer, v any) error {
	switch x := v.(type) {
	case *JSONObject:
		x.ForEach(func(key string, val any) bool {
			writeXMLEntry(buf, key, val)
			return true
		})
	case *JSONArray:
		for _, val := range x.values {
			writeXMLEntry(buf, "element", val)
		}
	default:
		// 基础值：转字符串
		writeXMLText(buf, toString(wrap(v, NewConfig()), ""))
	}
	return nil
}

// writeXMLEntry 写一个 <key>...</key>。若 val 是 JSONArray，则以 key 作为标签重复。
func writeXMLEntry(buf *bytes.Buffer, key string, val any) {
	if arr, ok := val.(*JSONArray); ok {
		for _, item := range arr.values {
			writeXMLEntry(buf, key, item)
		}
		return
	}
	buf.WriteString("<")
	buf.WriteString(key)
	buf.WriteString(">")
	switch x := val.(type) {
	case *JSONObject:
		_ = writeXMLValue(buf, x)
	case jsonNull:
		// 空标签
	default:
		writeXMLText(buf, toString(wrap(val, NewConfig()), ""))
	}
	buf.WriteString("</")
	buf.WriteString(key)
	buf.WriteString(">")
}

// writeXMLText 写带转义的文本。
func writeXMLText(buf *bytes.Buffer, s string) {
	if s == "" {
		return
	}
	if err := xml.EscapeText(buf, []byte(s)); err != nil {
		// 极少触发，回退原文
		buf.WriteString(fmt.Sprint(s))
	}
}
