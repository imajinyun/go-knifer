package conf

import (
	"fmt"
	"strconv"
	"strings"

	toml "github.com/pelletier/go-toml/v2"
)

// ParseTOML parses TOML into grouped configuration.
func ParseTOML(content string) (*Conf, error) {
	var root map[string]any
	if err := toml.Unmarshal([]byte(content), &root); err != nil {
		return nil, wrapConfigParse("parse toml content", err)
	}
	c := New()
	flattenTOMLMap(c, nil, root)
	return c, nil
}

func flattenTOMLMap(c *Conf, path []string, values map[string]any) {
	for key, child := range values {
		switch v := child.(type) {
		case map[string]any:
			flattenTOMLMap(c, appendPath(path, key), v)
		case []any:
			if flattenTOMLArrayTables(c, appendPath(path, key), v) {
				continue
			}
			c.SetByGroup(tomlGroup(path), key, joinTOMLArray(v))
		default:
			c.SetByGroup(tomlGroup(path), key, tomlScalarString(v))
		}
	}
}

func flattenTOMLArrayTables(c *Conf, path []string, values []any) bool {
	if len(values) == 0 {
		return false
	}
	for i, item := range values {
		m, ok := item.(map[string]any)
		if !ok {
			return false
		}
		flattenTOMLMap(c, appendPath(path, strconv.Itoa(i)), m)
	}
	return true
}

func appendPath(path []string, part string) []string {
	next := make([]string, 0, len(path)+1)
	next = append(next, path...)
	next = append(next, part)
	return next
}

func tomlGroup(path []string) string { return strings.Join(path, ".") }

func joinTOMLArray(values []any) string {
	parts := make([]string, 0, len(values))
	for _, value := range values {
		parts = append(parts, tomlScalarString(value))
	}
	return strings.Join(parts, ",")
}

func tomlScalarString(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprint(v)
	}
}
