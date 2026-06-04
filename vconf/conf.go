package vconf

import (
	"time"

	confimpl "github.com/imajinyun/go-knifer/internal/conf"
)

// Conf stores grouped key-value configuration.
type Conf = confimpl.Conf

// Error is the configuration module error type.
type Error = confimpl.ConfError

// New creates an empty Conf.
func New() *Conf { return confimpl.New() }

// Load 读取并解析 setting/properties 配置文件。Load reads and parses a setting/properties file.
func Load(path string) (*Conf, error) { return confimpl.Load(path) }

// LoadProfile loads a configuration file and applies profile-specific overrides.
func LoadProfile(path, profile string) (*Conf, error) { return confimpl.LoadProfile(path, profile) }

// Parse 解析 setting/properties 文本内容。Parse parses setting/properties content.
func Parse(content string) (*Conf, error) { return confimpl.Parse(content) }

// ParseBytes 解析 setting/properties 字节内容。ParseBytes parses setting/properties content.
func ParseBytes(content []byte) (*Conf, error) { return confimpl.ParseBytes(content) }

// ParseByExt parses content according to path extension.
func ParseByExt(path string, content []byte) (*Conf, error) {
	return confimpl.ParseByExt(path, content)
}

// ParseYAML 将简单 YAML 子集解析为分组配置。ParseYAML parses a small YAML subset into grouped configuration.
func ParseYAML(content string) (*Conf, error) { return confimpl.ParseYAML(content) }

// ParseYAMLFull parses YAML using yaml.v3 and flattens nested objects into grouped keys.
func ParseYAMLFull(content string) (*Conf, error) { return confimpl.ParseYAMLFull(content) }

// ParseTOML parses common TOML key-value and section syntax into grouped configuration.
func ParseTOML(content string) (*Conf, error) { return confimpl.ParseTOML(content) }

// Watch polls path and calls onChange after successful reloads.
func Watch(path string, interval time.Duration, onChange func(*Conf, error)) (func(), error) {
	return confimpl.Watch(path, interval, onChange)
}
