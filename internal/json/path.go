package json

import (
	"strings"
)

// pathSegment 描述路径中一段：要么是 key，要么是数组下标。
type pathSegment struct {
	key     string
	index   int
	isIndex bool
}

// parsePath 将 "$.a.b[0].c" / "a.b[0]" 解析为段序列。
func parsePath(path string) ([]pathSegment, error) {
	p := strings.TrimSpace(path)
	if p == "" {
		return nil, NewJSONError("empty path")
	}
	if strings.HasPrefix(p, "$") {
		p = strings.TrimPrefix(p, "$")
		p = strings.TrimPrefix(p, ".")
	}
	if p == "" {
		return nil, nil
	}

	var segs []pathSegment
	i := 0
	for i < len(p) {
		c := p[i]
		switch c {
		case '.':
			i++
			continue
		case '[':
			end := strings.IndexByte(p[i:], ']')
			if end < 0 {
				return nil, NewJSONError("unmatched '[' in path %q", path)
			}
			body := p[i+1 : i+end]
			n, ok := parseIndex(body)
			if !ok {
				return nil, NewJSONError("invalid index %q in path %q", body, path)
			}
			segs = append(segs, pathSegment{index: n, isIndex: true})
			i += end + 1
		default:
			// 收集到下一个 . 或 [
			end := i
			for end < len(p) && p[end] != '.' && p[end] != '[' {
				end++
			}
			segs = append(segs, pathSegment{key: p[i:end]})
			i = end
		}
	}
	return segs, nil
}

// getByPath 沿路径读取值。
func getByPath(root any, path string) any {
	segs, err := parsePath(path)
	if err != nil {
		return nil
	}
	cur := root
	for _, seg := range segs {
		if cur == nil || IsNull(cur) {
			return nil
		}
		if seg.isIndex {
			arr, ok := cur.(*JSONArray)
			if !ok {
				return nil
			}
			v, ok := arr.Get(seg.index)
			if !ok {
				return nil
			}
			cur = v
			continue
		}
		switch x := cur.(type) {
		case *JSONObject:
			v, ok := x.Get(seg.key)
			if !ok {
				return nil
			}
			cur = v
		case *JSONArray:
			// 若 key 是数字也支持
			if n, ok := parseIndex(seg.key); ok {
				v, ok := x.Get(n)
				if !ok {
					return nil
				}
				cur = v
				continue
			}
			return nil
		default:
			return nil
		}
	}
	return cur
}

// putByPath 沿路径写入值，必要时创建中间节点。
func putByPath(root any, path string, value any) error {
	segs, err := parsePath(path)
	if err != nil {
		return err
	}
	if len(segs) == 0 {
		return NewJSONError("empty path")
	}
	cur := root
	for i := 0; i < len(segs)-1; i++ {
		seg := segs[i]
		next := segs[i+1]
		switch x := cur.(type) {
		case *JSONObject:
			if seg.isIndex {
				return NewJSONError("cannot index object with [%d]", seg.index)
			}
			child, ok := x.Get(seg.key)
			if !ok || IsNull(child) {
				if next.isIndex {
					child = NewJSONArrayWithConfig(x.cfg)
				} else {
					child = NewJSONObjectWithConfig(x.cfg)
				}
				x.Set(seg.key, child)
			}
			cur = child
		case *JSONArray:
			if !seg.isIndex {
				return NewJSONError("cannot access array with key %q", seg.key)
			}
			child, _ := x.Get(seg.index)
			if child == nil || IsNull(child) {
				if next.isIndex {
					child = NewJSONArrayWithConfig(x.cfg)
				} else {
					child = NewJSONObjectWithConfig(x.cfg)
				}
				x.Set(seg.index, child)
			}
			cur = child
		default:
			return NewJSONError("path goes through non-container value")
		}
	}
	last := segs[len(segs)-1]
	switch x := cur.(type) {
	case *JSONObject:
		if last.isIndex {
			return NewJSONError("cannot index object with [%d]", last.index)
		}
		x.Set(last.key, value)
		return nil
	case *JSONArray:
		if !last.isIndex {
			return NewJSONError("cannot access array with key %q", last.key)
		}
		x.Set(last.index, value)
		return nil
	}
	return NewJSONError("invalid root container")
}
