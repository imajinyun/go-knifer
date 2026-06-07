package conf

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

const defaultGroup = ""

type expandConfig struct {
	envLookup func(string) string
}

type valueConfig struct {
	parseInt  func(string) (int, error)
	parseBool func(string) (bool, error)
}

type bindConfig struct {
	parseBool  func(string) (bool, error)
	parseInt   func(string, int, int) (int64, error)
	parseUint  func(string, int, int) (uint64, error)
	parseFloat func(string, int) (float64, error)
}

// ExpandOption customizes configuration variable expansion per call.
type ExpandOption func(*expandConfig)

// ValueOption customizes typed value getters per call.
type ValueOption func(*valueConfig)

// BindOption customizes struct binding per call.
type BindOption func(*bindConfig)

// WithEnvLookup sets the environment lookup function used for ${ENV:NAME} placeholders.
func WithEnvLookup(lookup func(string) string) ExpandOption {
	return func(c *expandConfig) {
		if lookup != nil {
			c.envLookup = lookup
		}
	}
}

// WithIntParser sets the parser used by GetIntWithOptions.
func WithIntParser(parser func(string) (int, error)) ValueOption {
	return func(c *valueConfig) {
		if parser != nil {
			c.parseInt = parser
		}
	}
}

// WithBoolParser sets the parser used by GetBoolWithOptions.
func WithBoolParser(parser func(string) (bool, error)) ValueOption {
	return func(c *valueConfig) {
		if parser != nil {
			c.parseBool = parser
		}
	}
}

// WithBindBoolParser sets the bool parser used by BindWithOptions and BindGroupWithOptions.
func WithBindBoolParser(parser func(string) (bool, error)) BindOption {
	return func(c *bindConfig) {
		if parser != nil {
			c.parseBool = parser
		}
	}
}

// WithBindIntParser sets the signed integer parser used by BindWithOptions and BindGroupWithOptions.
func WithBindIntParser(parser func(string, int, int) (int64, error)) BindOption {
	return func(c *bindConfig) {
		if parser != nil {
			c.parseInt = parser
		}
	}
}

// WithBindUintParser sets the unsigned integer parser used by BindWithOptions and BindGroupWithOptions.
func WithBindUintParser(parser func(string, int, int) (uint64, error)) BindOption {
	return func(c *bindConfig) {
		if parser != nil {
			c.parseUint = parser
		}
	}
}

// WithBindFloatParser sets the floating-point parser used by BindWithOptions and BindGroupWithOptions.
func WithBindFloatParser(parser func(string, int) (float64, error)) BindOption {
	return func(c *bindConfig) {
		if parser != nil {
			c.parseFloat = parser
		}
	}
}

func applyExpandOptions(opts []ExpandOption) expandConfig {
	cfg := expandConfig{envLookup: os.Getenv}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	if cfg.envLookup == nil {
		cfg.envLookup = os.Getenv
	}
	return cfg
}

func applyValueOptions(opts []ValueOption) valueConfig {
	cfg := valueConfig{parseInt: strconv.Atoi, parseBool: strconv.ParseBool}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	if cfg.parseInt == nil {
		cfg.parseInt = strconv.Atoi
	}
	if cfg.parseBool == nil {
		cfg.parseBool = strconv.ParseBool
	}
	return cfg
}

func applyBindOptions(opts []BindOption) bindConfig {
	cfg := bindConfig{
		parseBool:  strconv.ParseBool,
		parseInt:   strconv.ParseInt,
		parseUint:  strconv.ParseUint,
		parseFloat: strconv.ParseFloat,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	if cfg.parseBool == nil {
		cfg.parseBool = strconv.ParseBool
	}
	if cfg.parseInt == nil {
		cfg.parseInt = strconv.ParseInt
	}
	if cfg.parseUint == nil {
		cfg.parseUint = strconv.ParseUint
	}
	if cfg.parseFloat == nil {
		cfg.parseFloat = strconv.ParseFloat
	}
	return cfg
}

// Conf stores grouped key-value configuration.
type Conf struct {
	data map[string]map[string]string
}

// New creates an empty Conf.
func New() *Conf {
	return &Conf{data: map[string]map[string]string{defaultGroup: {}}}
}

// Load 读取并解析配置文件。Load reads and parses a configuration file.
func Load(path string) (*Conf, error) {
	return LoadWithOptions(path, LoadOptions{})
}

// LoadProfile loads a configuration file and applies profile-specific overrides.
func LoadProfile(path, profile string) (*Conf, error) {
	return LoadProfileWithOptions(path, profile, LoadOptions{})
}

// LoadProfileWithOptions loads a configuration file with options and applies profile-specific overrides.
func LoadProfileWithOptions(path, profile string, opts LoadOptions) (*Conf, error) {
	c, err := LoadWithOptions(path, opts)
	if err != nil {
		return nil, err
	}
	return c.ApplyProfile(profile), nil
}

// Parse 解析 setting/properties 文本内容。Parse parses setting/properties content.
func Parse(content string) (*Conf, error) { return ParseBytes([]byte(content)) }

// ParseBytes 解析 setting/properties 字节内容。ParseBytes parses setting/properties content.
func ParseBytes(content []byte) (*Conf, error) {
	s := New()
	group := defaultGroup
	scanner := bufio.NewScanner(bytes.NewReader(content))
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			group = strings.TrimSpace(line[1 : len(line)-1])
			s.ensureGroup(group)
			continue
		}
		idx := strings.IndexAny(line, "=:")
		if idx < 0 {
			return nil, invalidInputf("invalid setting line %d: %s", lineNo, line)
		}
		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])
		if key == "" {
			return nil, invalidInputf("empty setting key at line %d", lineNo)
		}
		s.SetByGroup(group, key, unquote(value))
	}
	if err := scanner.Err(); err != nil {
		return nil, wrapConfigParse("scan setting content", err)
	}
	return s, nil
}

// Get 从默认分组获取配置值。Get returns a value from the default group.
func (s *Conf) Get(key string) string { return s.GetByGroup(defaultGroup, key) }

// GetExpanded returns a value from default group after variable expansion.
func (s *Conf) GetExpanded(key string) string { return s.GetByGroupExpanded(defaultGroup, key) }

// GetExpandedWithOptions returns a value from default group after variable expansion with per-call options.
func (s *Conf) GetExpandedWithOptions(key string, opts ...ExpandOption) string {
	return s.GetByGroupExpandedWithOptions(defaultGroup, key, opts...)
}

// GetOrDefault 从默认分组获取配置值，不存在时返回 def。GetOrDefault returns a value from the default group or def when absent.
func (s *Conf) GetOrDefault(key, def string) string {
	if v, ok := s.Lookup(defaultGroup, key); ok {
		return v
	}
	return def
}

// GetByGroup 获取指定分组中的配置值。GetByGroup returns a grouped value.
func (s *Conf) GetByGroup(group, key string) string {
	v, _ := s.Lookup(group, key)
	return v
}

// GetByGroupExpanded returns a grouped value after variable expansion.
func (s *Conf) GetByGroupExpanded(group, key string) string {
	return s.GetByGroupExpandedWithOptions(group, key)
}

// GetByGroupExpandedWithOptions returns a grouped value after variable expansion with per-call options.
func (s *Conf) GetByGroupExpandedWithOptions(group, key string, opts ...ExpandOption) string {
	v, ok := s.Lookup(group, key)
	if !ok {
		return ""
	}
	cfg := applyExpandOptions(opts)
	return s.expandValue(group, v, map[string]bool{}, cfg)
}

// Lookup 获取指定分组中的配置值并返回是否存在。Lookup returns a grouped value and whether it exists.
func (s *Conf) Lookup(group, key string) (string, bool) {
	if s == nil || s.data == nil {
		return "", false
	}
	m, ok := s.data[group]
	if !ok {
		return "", false
	}
	v, ok := m[key]
	return v, ok
}

// GetInt 从默认分组获取 int 值，不存在或格式非法时返回 def。GetInt returns an int value from the default group or def when absent/invalid.
func (s *Conf) GetInt(key string, def int) int {
	return s.GetIntWithOptions(key, def)
}

// GetIntWithOptions returns an int value from the default group using per-call parser options.
func (s *Conf) GetIntWithOptions(key string, def int, opts ...ValueOption) int {
	v, ok := s.Lookup(defaultGroup, key)
	if !ok {
		return def
	}
	cfg := applyValueOptions(opts)
	n, err := cfg.parseInt(v)
	if err != nil {
		return def
	}
	return n
}

// GetBool 从默认分组获取 bool 值，不存在或格式非法时返回 def。GetBool returns a bool value from the default group or def when absent/invalid.
func (s *Conf) GetBool(key string, def bool) bool {
	return s.GetBoolWithOptions(key, def)
}

// GetBoolWithOptions returns a bool value from the default group using per-call parser options.
func (s *Conf) GetBoolWithOptions(key string, def bool, opts ...ValueOption) bool {
	v, ok := s.Lookup(defaultGroup, key)
	if !ok {
		return def
	}
	cfg := applyValueOptions(opts)
	b, err := cfg.parseBool(v)
	if err != nil {
		return def
	}
	return b
}

// Set 将配置值写入默认分组。Set stores a value in the default group.
func (s *Conf) Set(key, value string) { s.SetByGroup(defaultGroup, key, value) }

// SetByGroup 将配置值写入指定分组。SetByGroup stores a grouped value.
func (s *Conf) SetByGroup(group, key, value string) {
	s.ensureGroup(group)
	s.data[group][key] = value
}

// Delete removes a value from the default group.
func (s *Conf) Delete(key string) { s.DeleteByGroup(defaultGroup, key) }

// DeleteByGroup removes a value from a group.
func (s *Conf) DeleteByGroup(group, key string) {
	if s == nil || s.data == nil {
		return
	}
	if m, ok := s.data[group]; ok {
		delete(m, key)
	}
}

// Groups 返回全部分组名称。Groups returns all group names.
func (s *Conf) Groups() []string {
	if s == nil || s.data == nil {
		return []string{}
	}
	groups := make([]string, 0, len(s.data))
	for g := range s.data {
		groups = append(groups, g)
	}
	sort.Strings(groups)
	return groups
}

// Keys 返回指定分组中的全部键。Keys returns keys from group.
func (s *Conf) Keys(group string) []string {
	if s == nil || s.data == nil {
		return []string{}
	}
	m := s.data[group]
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// ToMap 返回所有分组配置的深拷贝。ToMap returns a deep copy of all groups.
func (s *Conf) ToMap() map[string]map[string]string {
	if s == nil || s.data == nil {
		return map[string]map[string]string{}
	}
	out := make(map[string]map[string]string, len(s.data))
	for g, m := range s.data {
		out[g] = make(map[string]string, len(m))
		for k, v := range m {
			out[g][k] = v
		}
	}
	return out
}

// Expand returns a copy with ${key}, ${group.key}, ${ENV:NAME}, and ${key:default} placeholders resolved.
func (s *Conf) Expand() *Conf {
	return s.ExpandWithOptions()
}

// ExpandWithOptions returns a copy with placeholders resolved using per-call options.
func (s *Conf) ExpandWithOptions(opts ...ExpandOption) *Conf {
	out := New()
	if s == nil || s.data == nil {
		return out
	}
	cfg := applyExpandOptions(opts)
	for group, m := range s.data {
		for key, value := range m {
			out.SetByGroup(group, key, s.expandValue(group, value, map[string]bool{group + "." + key: true}, cfg))
		}
	}
	return out
}

// ApplyProfile overlays groups named profile.<profile> and profile.<profile>.<group>.
func (s *Conf) ApplyProfile(profile string) *Conf {
	out := New()
	if s == nil || s.data == nil {
		return out
	}
	profile = strings.TrimSpace(profile)
	prefix := "profile." + profile
	for group, m := range s.data {
		if profile != "" && group == "profile" {
			continue
		}
		if profile != "" && (group == prefix || strings.HasPrefix(group, prefix+".")) {
			continue
		}
		for k, v := range m {
			out.SetByGroup(group, k, v)
		}
	}
	if profile == "" {
		return out
	}
	for group, m := range s.data {
		if group != prefix && !strings.HasPrefix(group, prefix+".") {
			continue
		}
		targetGroup := defaultGroup
		if strings.HasPrefix(group, prefix+".") {
			targetGroup = strings.TrimPrefix(group, prefix+".")
		}
		for k, v := range m {
			out.SetByGroup(targetGroup, k, v)
		}
	}
	if profileGroup := s.data["profile"]; len(profileGroup) > 0 {
		keyPrefix := profile + "."
		for k, v := range profileGroup {
			if !strings.HasPrefix(k, keyPrefix) {
				continue
			}
			rest := strings.TrimPrefix(k, keyPrefix)
			if rest == "" {
				continue
			}
			targetGroup := defaultGroup
			targetKey := rest
			if idx := strings.LastIndex(rest, "."); idx >= 0 {
				targetGroup = rest[:idx]
				targetKey = rest[idx+1:]
			}
			out.SetByGroup(targetGroup, targetKey, v)
		}
	}
	return out
}

// Bind fills dst from the default group using conf tags or field names.
func (s *Conf) Bind(dst any) error { return s.BindWithOptions(dst) }

// BindWithOptions fills dst from the default group using per-call parser options.
func (s *Conf) BindWithOptions(dst any, opts ...BindOption) error {
	return s.BindGroupWithOptions(defaultGroup, dst, opts...)
}

// BindGroup fills dst from a group using conf tags or field names.
func (s *Conf) BindGroup(group string, dst any) error {
	return s.BindGroupWithOptions(group, dst)
}

// BindGroupWithOptions fills dst from a group using conf tags or field names and per-call parser options.
func (s *Conf) BindGroupWithOptions(group string, dst any, opts ...BindOption) error {
	rv := reflect.ValueOf(dst)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return invalidInputf("bind target must be a non-nil pointer")
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return invalidInputf("bind target must point to a struct")
	}
	cfg := applyBindOptions(opts)
	return s.bindStruct(group, "", rv, cfg)
}

func (s *Conf) ensureGroup(group string) {
	if s.data == nil {
		s.data = map[string]map[string]string{}
	}
	if _, ok := s.data[group]; !ok {
		s.data[group] = map[string]string{}
	}
}

func (s *Conf) bindStruct(group, prefix string, rv reflect.Value, cfg bindConfig) error {
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if field.PkgPath != "" {
			continue
		}
		name, skip := confFieldName(field)
		if skip {
			continue
		}
		key := name
		if prefix != "" {
			key = prefix + "." + name
		}
		fv := rv.Field(i)
		if fv.Kind() == reflect.Struct && !field.Anonymous && field.Type != reflect.TypeOf(time.Time{}) {
			if err := s.bindStruct(group, key, fv, cfg); err != nil {
				return err
			}
			continue
		}
		value, ok := s.Lookup(group, key)
		if !ok {
			continue
		}
		if err := setReflectValue(fv, value, cfg); err != nil {
			return invalidInputf("bind %s: %s", key, err.Error())
		}
	}
	return nil
}

func confFieldName(field reflect.StructField) (string, bool) {
	name, _, skip := parseConfTag(field)
	return name, skip
}

func setReflectValue(v reflect.Value, text string, cfg bindConfig) error {
	if !v.CanSet() {
		return nil
	}
	switch v.Kind() {
	case reflect.String:
		v.SetString(text)
	case reflect.Bool:
		b, err := cfg.parseBool(text)
		if err != nil {
			return err
		}
		v.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := cfg.parseInt(text, 10, v.Type().Bits())
		if err != nil {
			return err
		}
		v.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := cfg.parseUint(text, 10, v.Type().Bits())
		if err != nil {
			return err
		}
		v.SetUint(n)
	case reflect.Float32, reflect.Float64:
		f, err := cfg.parseFloat(text, v.Type().Bits())
		if err != nil {
			return err
		}
		v.SetFloat(f)
	case reflect.Slice:
		parts := splitList(text)
		slice := reflect.MakeSlice(v.Type(), 0, len(parts))
		for _, part := range parts {
			elem := reflect.New(v.Type().Elem()).Elem()
			if err := setReflectValue(elem, part, cfg); err != nil {
				return err
			}
			slice = reflect.Append(slice, elem)
		}
		v.Set(slice)
	default:
		return fmt.Errorf("unsupported field type %s", v.Type())
	}
	return nil
}

func splitList(text string) []string {
	if strings.TrimSpace(text) == "" {
		return nil
	}
	parts := strings.Split(text, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func (s *Conf) expandValue(group, value string, seen map[string]bool, cfg expandConfig) string {
	return os.Expand(value, func(name string) string {
		if strings.HasPrefix(name, "ENV:") {
			return cfg.envLookup(strings.TrimPrefix(name, "ENV:"))
		}
		key, fallback, hasFallback := strings.Cut(name, ":")
		lookupGroup := group
		lookupKey := key
		if dot := strings.Index(key, "."); dot > 0 {
			lookupGroup, lookupKey = key[:dot], key[dot+1:]
		}
		seenKey := lookupGroup + "." + lookupKey
		if seen[seenKey] {
			if hasFallback {
				return fallback
			}
			return ""
		}
		v, ok := s.Lookup(lookupGroup, lookupKey)
		if !ok {
			if hasFallback {
				return fallback
			}
			return ""
		}
		seen[seenKey] = true
		defer delete(seen, seenKey)
		return s.expandValue(lookupGroup, v, seen, cfg)
	})
}

func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
