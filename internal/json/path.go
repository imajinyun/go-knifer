package json

import (
	"strings"
)

// pathSegment describes one path segment, either a key or an array index.
type pathSegment struct {
	key     string
	index   int
	isIndex bool
}

// parsePath parses "$.a.b[0].c" or "a.b[0]" into segments.
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
			// Collect until the next . or [.
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

// getByPath reads a value along a path.
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
			// Also supports numeric keys.
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

// putByPath writes a value along a path and creates intermediate nodes when needed.
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
