package conf

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	TypeString = "string"
	TypeBool   = "bool"
	TypeInt    = "int"
	TypeFloat  = "float"
	TypeList   = "list"
)

// FieldRule describes a schema rule for a configuration key.
type FieldRule struct {
	Group    string
	Key      string
	Required bool
	Type     string
	Default  string
	Choices  []string
}

// Schema contains configuration validation rules.
type Schema struct {
	Fields []FieldRule
}

// ValidateSchema validates s against schema.
func (s *Conf) ValidateSchema(schema Schema) error {
	for _, rule := range schema.Fields {
		group := rule.Group
		key := rule.Key
		if key == "" {
			return invalidInputf("schema key must not be empty")
		}
		value, ok := s.Lookup(group, key)
		if !ok || value == "" {
			if rule.Required {
				return invalidInputf("required config %s is missing", schemaPath(group, key))
			}
			continue
		}
		if err := validateType(value, rule.Type); err != nil {
			return invalidInputf("config %s type mismatch: %s", schemaPath(group, key), err.Error())
		}
		if len(rule.Choices) > 0 && !containsString(rule.Choices, value) {
			return invalidInputf("config %s value %q is not in choices %v", schemaPath(group, key), value, rule.Choices)
		}
	}
	return nil
}

// ApplyDefaults returns a copy with schema defaults applied.
func (s *Conf) ApplyDefaults(schema Schema) *Conf {
	out := Merge(s)
	for _, rule := range schema.Fields {
		if rule.Key == "" || rule.Default == "" {
			continue
		}
		if _, ok := out.Lookup(rule.Group, rule.Key); !ok {
			out.SetByGroup(rule.Group, rule.Key, rule.Default)
		}
	}
	return out
}

// ValidateStruct validates required conf tags on dst against s.
func (s *Conf) ValidateStruct(dst any) error {
	rules, err := SchemaFromStruct(dst)
	if err != nil {
		return err
	}
	return s.ValidateSchema(rules)
}

// SchemaFromStruct builds a schema from conf tags. A tag like `conf:"port,required,int"` marks a required int.
func SchemaFromStruct(dst any) (Schema, error) {
	return schemaFromStruct(dst)
}

func schemaFromStruct(dst any) (Schema, error) {
	rv := reflect.ValueOf(dst)
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			rv = reflect.New(rv.Type().Elem())
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return Schema{}, invalidInputf("schema target must be a struct or struct pointer")
	}
	var schema Schema
	collectStructRules(rv.Type(), "", &schema)
	return schema, nil
}

func collectStructRules(rt reflect.Type, prefix string, schema *Schema) {
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if field.PkgPath != "" {
			continue
		}
		name, options, skip := parseConfTag(field)
		if skip {
			continue
		}
		key := name
		if prefix != "" {
			key = prefix + "." + name
		}
		if field.Type.Kind() == reflect.Struct && !field.Anonymous && field.Type != reflect.TypeOf(time.Time{}) {
			collectStructRules(field.Type, key, schema)
			continue
		}
		rule := FieldRule{Key: key, Type: inferRuleType(field.Type)}
		for _, option := range options {
			option = strings.TrimSpace(option)
			switch {
			case option == "required":
				rule.Required = true
			case strings.HasPrefix(option, "type="):
				rule.Type = strings.TrimPrefix(option, "type=")
			case option == TypeString || option == TypeBool || option == TypeInt || option == TypeFloat || option == TypeList:
				rule.Type = option
			case strings.HasPrefix(option, "default="):
				rule.Default = strings.TrimPrefix(option, "default=")
			case strings.HasPrefix(option, "choices="):
				rule.Choices = strings.Split(strings.TrimPrefix(option, "choices="), "|")
			}
		}
		schema.Fields = append(schema.Fields, rule)
	}
}

func parseConfTag(field reflect.StructField) (string, []string, bool) {
	tag := field.Tag.Get("conf")
	if tag == "-" {
		return "", nil, true
	}
	parts := splitTag(tag)
	name := strings.ToLower(field.Name)
	if len(parts) > 0 && parts[0] != "" {
		name = parts[0]
	}
	if len(parts) <= 1 {
		return name, nil, false
	}
	return name, parts[1:], false
}

func splitTag(tag string) []string {
	if tag == "" {
		return nil
	}
	parts := strings.Split(tag, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func inferRuleType(rt reflect.Type) string {
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}
	switch rt.Kind() {
	case reflect.Bool:
		return TypeBool
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return TypeInt
	case reflect.Float32, reflect.Float64:
		return TypeFloat
	case reflect.Slice, reflect.Array:
		return TypeList
	default:
		return TypeString
	}
}

func validateType(value, typ string) error {
	switch strings.ToLower(strings.TrimSpace(typ)) {
	case "", TypeString:
		return nil
	case TypeBool:
		_, err := strconv.ParseBool(value)
		return err
	case TypeInt:
		_, err := strconv.ParseInt(value, 10, 64)
		return err
	case TypeFloat:
		_, err := strconv.ParseFloat(value, 64)
		return err
	case TypeList:
		return nil
	default:
		return fmt.Errorf("unsupported type %s", typ)
	}
}

func schemaPath(group, key string) string {
	if group == "" {
		return key
	}
	return group + "." + key
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
