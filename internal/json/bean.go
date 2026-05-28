package json

import (
	"bytes"
	"encoding/json"
)

// ToBean 将 JSON 值转换到给定 dst（应为指针）。
// 当传入 *JSONObject/*JSONArray/string/[]byte/map/slice 等时，会先序列化为
// JSON 字节，再交给 encoding/json 反序列化到 dst。
func ToBean(src any, dst any) error {
	if dst == nil {
		return NewJSONError("dst is nil")
	}
	var data []byte
	switch x := src.(type) {
	case []byte:
		data = x
	case string:
		data = []byte(x)
	default:
		w := wrap(src, NewConfig())
		s, err := writeValue(w, 0, 0)
		if err != nil {
			return err
		}
		data = []byte(s)
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(dst); err != nil {
		return WrapJSONError(err, "to bean failed")
	}
	return nil
}

// ToList 将 JSON 数组转换到 dst（必须是指向 slice 的指针）。
func ToList(src any, dst any) error { return ToBean(src, dst) }
