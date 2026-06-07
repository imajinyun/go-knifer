package conf

import (
	"path/filepath"
	"strings"
)

type parseConfig struct {
	yamlUnmarshal func([]byte, any) error
	parsers       map[string]func([]byte) (*Conf, error)
}

// ParseOption customizes ParseByExt and full YAML parsing helpers per call.
type ParseOption func(*parseConfig)

// WithYAMLUnmarshalFunc sets the YAML unmarshal provider used by ParseYAMLFullWithOptions.
func WithYAMLUnmarshalFunc(unmarshal func([]byte, any) error) ParseOption {
	return func(c *parseConfig) {
		if unmarshal != nil {
			c.yamlUnmarshal = unmarshal
		}
	}
}

// WithParserForExt sets the parser used by ParseByExtWithOptions for an extension.
func WithParserForExt(ext string, parser func([]byte) (*Conf, error)) ParseOption {
	return func(c *parseConfig) {
		if parser == nil {
			return
		}
		ext = strings.ToLower(strings.TrimSpace(ext))
		if ext == "" {
			return
		}
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		if c.parsers == nil {
			c.parsers = map[string]func([]byte) (*Conf, error){}
		}
		c.parsers[ext] = parser
	}
}

func applyParseOptions(opts []ParseOption) parseConfig {
	cfg := parseConfig{yamlUnmarshal: defaultYAMLUnmarshal}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	if cfg.yamlUnmarshal == nil {
		cfg.yamlUnmarshal = defaultYAMLUnmarshal
	}
	return cfg
}

// ParseByExt parses content according to path extension.
func ParseByExt(path string, content []byte) (*Conf, error) {
	return ParseByExtWithOptions(path, content)
}

// ParseByExtWithOptions parses content according to path extension with custom providers.
func ParseByExtWithOptions(path string, content []byte, opts ...ParseOption) (*Conf, error) {
	cfg := applyParseOptions(opts)
	ext := strings.ToLower(filepath.Ext(path))
	if parser := cfg.parsers[ext]; parser != nil {
		return parser(content)
	}
	switch ext {
	case ".yaml", ".yml":
		return ParseYAMLFullWithOptions(string(content), opts...)
	case ".toml":
		return ParseTOML(string(content))
	default:
		return ParseBytes(content)
	}
}
