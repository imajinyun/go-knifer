package conf

import (
	"strconv"
	"testing"
)

func TestSchemaFromStructAndValidateStruct(t *testing.T) {
	type appConfig struct {
		Name string `conf:"name,required"`
		Port int    `conf:"port,required,int"`
		Mode string `conf:"mode,default=dev,choices=dev|prod"`
	}
	c := New()
	c.Set("name", "demo")
	c.Set("port", "8080")
	c.Set("mode", "dev")
	if err := c.ValidateStruct(appConfig{}); err != nil {
		t.Fatalf("ValidateStruct() error = %v", err)
	}
	schema, err := SchemaFromStruct(appConfig{})
	if err != nil {
		t.Fatal(err)
	}
	if len(schema.Fields) != 3 {
		t.Fatalf("SchemaFromStruct fields = %d", len(schema.Fields))
	}
}

func TestValidateSchemaWithOptionsUsesParsers(t *testing.T) {
	c := New()
	c.Set("debug", "custom-bool")
	c.Set("port", "custom-int")
	c.Set("ratio", "custom-float")

	var boolCalled, intCalled, floatCalled int
	err := c.ValidateSchemaWithOptions(Schema{Fields: []FieldRule{
		{Key: "debug", Required: true, Type: TypeBool},
		{Key: "port", Required: true, Type: TypeInt},
		{Key: "ratio", Required: true, Type: TypeFloat},
	}},
		WithSchemaBoolParser(func(text string) (bool, error) {
			boolCalled++
			if text == "custom-bool" {
				return true, nil
			}
			return strconv.ParseBool(text)
		}),
		WithSchemaIntParser(func(text string, base, bitSize int) (int64, error) {
			intCalled++
			if text == "custom-int" {
				return 8080, nil
			}
			return strconv.ParseInt(text, base, bitSize)
		}),
		WithSchemaFloatParser(func(text string, bitSize int) (float64, error) {
			floatCalled++
			if text == "custom-float" {
				return 0.75, nil
			}
			return strconv.ParseFloat(text, bitSize)
		}),
	)
	if err != nil {
		t.Fatalf("ValidateSchemaWithOptions() error = %v", err)
	}
	if boolCalled != 1 || intCalled != 1 || floatCalled != 1 {
		t.Fatalf("schema parser calls bool=%d int=%d float=%d", boolCalled, intCalled, floatCalled)
	}
}

func TestValidateStructWithOptionsUsesParsers(t *testing.T) {
	type appConfig struct {
		Port int `conf:"port,required,int"`
	}
	c := New()
	c.Set("port", "custom-int")

	called := false
	if err := c.ValidateStructWithOptions(appConfig{}, WithSchemaIntParser(func(text string, base, bitSize int) (int64, error) {
		called = true
		if text == "custom-int" {
			return 8080, nil
		}
		return strconv.ParseInt(text, base, bitSize)
	})); err != nil {
		t.Fatalf("ValidateStructWithOptions() error = %v", err)
	}
	if !called {
		t.Fatal("schema int parser was not called")
	}
}
