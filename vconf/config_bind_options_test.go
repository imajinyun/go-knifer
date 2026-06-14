package vconf_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/imajinyun/go-knifer/vconf"
)

func TestFacadeBindAndSchemaParserOptions(t *testing.T) {
	s, err := vconf.Parse(`
flag=yes
count=custom-int
amount=custom-uint
ratio=custom-float
items=1,2,3
schema_bool=yes
schema_float=custom-float
choice=blue
`)
	if err != nil {
		t.Fatal(err)
	}

	type bindConfig struct {
		Flag   bool    `conf:"flag"`
		Count  int     `conf:"count"`
		Amount uint    `conf:"amount"`
		Ratio  float64 `conf:"ratio"`
		Items  []int   `conf:"items"`
	}
	var cfg bindConfig
	if err := s.BindWithOptions(&cfg,
		vconf.WithBindBoolParser(func(value string) (bool, error) {
			return value == "yes", nil
		}),
		vconf.WithBindIntParser(func(value string, base int, bitSize int) (int64, error) {
			if value == "custom-int" {
				return 42, nil
			}
			return 7, nil
		}),
		vconf.WithBindUintParser(func(value string, base int, bitSize int) (uint64, error) {
			if value == "custom-uint" {
				return 9, nil
			}
			return 3, nil
		}),
		vconf.WithBindFloatParser(func(value string, bitSize int) (float64, error) {
			if value == "custom-float" {
				return 1.5, nil
			}
			return 0, errors.New("unexpected float")
		}),
	); err != nil {
		t.Fatalf("BindWithOptions() error = %v", err)
	}
	if !cfg.Flag || cfg.Count != 42 || cfg.Amount != 9 || cfg.Ratio != 1.5 || !reflect.DeepEqual(cfg.Items, []int{7, 7, 7}) {
		t.Fatalf("BindWithOptions cfg = %#v", cfg)
	}

	err = s.ValidateSchemaWithOptions(vconf.Schema{Fields: []vconf.FieldRule{
		{Key: "schema_bool", Required: true, Type: vconf.TypeBool},
		{Key: "schema_float", Required: true, Type: vconf.TypeFloat},
		{Key: "choice", Required: true, Choices: []string{"red", "blue"}},
	}},
		vconf.WithSchemaBoolParser(func(value string) (bool, error) {
			return value == "yes", nil
		}),
		vconf.WithSchemaFloatParser(func(value string, bitSize int) (float64, error) {
			if value == "custom-float" {
				return 2.5, nil
			}
			return 0, errors.New("unexpected schema float")
		}),
	)
	if err != nil {
		t.Fatalf("ValidateSchemaWithOptions() error = %v", err)
	}
}
