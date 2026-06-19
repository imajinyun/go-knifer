package vmask_test

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vmask"
)

func ExampleMobilePhone() {
	fmt.Println(vmask.MobilePhone("13812345678"))
	// Output: 138****5678
}

func ExampleEmail() {
	fmt.Println(vmask.Email("test@example.com"))
	// Output: t***@example.com
}

func ExampleIPv4() {
	fmt.Println(vmask.IPv4("192.0.2.15"))
	// Output: 192.*.*.*
}
