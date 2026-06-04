package conf

import (
	"path/filepath"
	"strings"
)

// ParseByExt parses content according to path extension.
func ParseByExt(path string, content []byte) (*Conf, error) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".yaml", ".yml":
		return ParseYAMLFull(string(content))
	case ".toml":
		return ParseTOML(string(content))
	default:
		return ParseBytes(content)
	}
}
