package json

// jsonNull 对应 hutool 的 JSONNull，表示 JSON 中的 null。
type jsonNull struct{}

// Null 是单例 JSON null。
var Null = jsonNull{}

// String 实现 Stringer：输出 "null"。
func (jsonNull) String() string { return "null" }

// IsNull 判断 v 是否为 nil 或 JSON Null。
func IsNull(v any) bool {
	if v == nil {
		return true
	}
	_, ok := v.(jsonNull)
	return ok
}
