package vmail_test

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vmail"
)

func ExampleNewMessage() {
	m, err := vmail.NewMessage(
		vmail.WithFrom("sender@example.com"),
		vmail.WithTo("recipient@example.com"),
		vmail.WithSubject("Hello"),
		vmail.WithText("World"),
	)
	fmt.Println(m != nil)
	fmt.Println(err)
	// Output:
	// true
	// <nil>
}

func ExampleNewAddress() {
	addr, err := vmail.NewAddress("Alice", "alice@example.com")

	fmt.Println(addr.String())
	fmt.Println(err)
	// Output:
	// "Alice" <alice@example.com>
	// <nil>
}

func ExampleParseAddressList() {
	list, err := vmail.ParseAddressList("bob@example.com, carol@example.com")

	fmt.Println(len(list), list[0].Email, list[1].Email)
	fmt.Println(err)
	// Output:
	// 2 bob@example.com carol@example.com
	// <nil>
}
