package vcodec_test

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vcodec"
)

func ExampleBase64Encode() {
	encoded := vcodec.Base64Encode([]byte("go-knifer"))
	fmt.Println(encoded)
	// Output: Z28ta25pZmVy
}

func ExampleBase64Decode() {
	decoded, _ := vcodec.Base64Decode("Z28ta25pZmVy")
	fmt.Println(string(decoded))
	// Output: go-knifer
}

func ExampleHexEncode() {
	encoded := vcodec.HexEncode([]byte{0x47, 0x6f})
	fmt.Println(encoded)
	// Output: 476f
}
