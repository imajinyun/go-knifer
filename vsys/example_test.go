package vsys_test

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vsys"
)

func ExampleGetCurrentPID() {
	pid := vsys.GetCurrentPID()
	fmt.Println(pid > 0)
	// Output: true
}

func ExampleEnvWithOptions() {
	lookup := vsys.WithEnvLookupFunc(func(key string) (string, bool) {
		if key == "APP_MODE" {
			return "test", true
		}
		return "", false
	})

	fmt.Println(vsys.EnvWithOptions("APP_MODE", lookup))
	fmt.Println(vsys.EnvWithOptions("MISSING", lookup))
	// Output:
	// test
	//
}

func ExampleEnvOrDefaultWithOptions() {
	lookup := vsys.WithEnvLookupFunc(func(key string) (string, bool) {
		if key == "PORT" {
			return "8080", true
		}
		return "", false
	})

	fmt.Println(vsys.EnvOrDefaultWithOptions("MISSING", "fallback", lookup))
	fmt.Println(vsys.EnvIntWithOptions("PORT", 0, lookup))
	// Output:
	// fallback
	// 8080
}
