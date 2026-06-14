package conf

import (
	"reflect"
	"strconv"
	"testing"
)

func TestBindWithOptionsUsesParsers(t *testing.T) {
	s := New()
	s.SetByGroup("server", "port", "custom-int")
	s.SetByGroup("server", "debug", "custom-bool")
	s.SetByGroup("server", "ratio", "custom-float")
	s.SetByGroup("server", "ids", "a,b")

	type serverConf struct {
		Port  int     `conf:"port"`
		Debug bool    `conf:"debug"`
		Ratio float64 `conf:"ratio"`
		IDs   []uint  `conf:"ids"`
	}
	var cfg serverConf
	var intCalled, boolCalled, floatCalled, uintCalled int
	err := s.BindGroupWithOptions("server", &cfg,
		WithBindIntParser(func(text string, base, bitSize int) (int64, error) {
			intCalled++
			if text == "custom-int" {
				return 8080, nil
			}
			return strconv.ParseInt(text, base, bitSize)
		}),
		WithBindBoolParser(func(text string) (bool, error) {
			boolCalled++
			return text == "custom-bool", nil
		}),
		WithBindFloatParser(func(text string, bitSize int) (float64, error) {
			floatCalled++
			if text == "custom-float" {
				return 0.75, nil
			}
			return strconv.ParseFloat(text, bitSize)
		}),
		WithBindUintParser(func(text string, base, bitSize int) (uint64, error) {
			uintCalled++
			switch text {
			case "a":
				return 1, nil
			case "b":
				return 2, nil
			default:
				return strconv.ParseUint(text, base, bitSize)
			}
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(cfg, serverConf{Port: 8080, Debug: true, Ratio: 0.75, IDs: []uint{1, 2}}) {
		t.Fatalf("BindGroupWithOptions = %#v", cfg)
	}
	if intCalled != 1 || boolCalled != 1 || floatCalled != 1 || uintCalled != 2 {
		t.Fatalf("parser calls int=%d bool=%d float=%d uint=%d", intCalled, boolCalled, floatCalled, uintCalled)
	}
}
