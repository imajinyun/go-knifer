package vsetting

import settingimpl "github.com/imajinyun/go-knifer/internal/setting"

// Setting 存储分组键值配置，对应 Hutool Setting/GroupedMap 的核心能力。Setting stores grouped key-value configuration.
type Setting = settingimpl.Setting

// New 创建空的 Setting。New creates an empty Setting.
func New() *Setting { return settingimpl.New() }

// Load 读取并解析 setting/properties 配置文件。Load reads and parses a setting/properties file.
func Load(path string) (*Setting, error) { return settingimpl.Load(path) }

// Parse 解析 setting/properties 文本内容。Parse parses setting/properties content.
func Parse(content string) (*Setting, error) { return settingimpl.Parse(content) }

// ParseBytes 解析 setting/properties 字节内容。ParseBytes parses setting/properties content.
func ParseBytes(content []byte) (*Setting, error) { return settingimpl.ParseBytes(content) }

// ParseYAML 将简单 YAML 子集解析为分组配置。ParseYAML parses a small YAML subset into grouped configuration.
func ParseYAML(content string) (*Setting, error) { return settingimpl.ParseYAML(content) }
