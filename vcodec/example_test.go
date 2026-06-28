package vcodec_test

import (
	"encoding/base64"
	"fmt"

	"github.com/imajinyun/knifer-go/vcodec"
)

func ExampleBase64Encode() {
	encoded := vcodec.Base64Encode([]byte("knifer-go"))
	fmt.Println(encoded)
	// Output: a25pZmVyLWdv
}

func ExampleBase64Decode() {
	decoded, _ := vcodec.Base64Decode("a25pZmVyLWdv")
	fmt.Println(string(decoded))
	// Output: knifer-go
}

func ExampleHexEncode() {
	encoded := vcodec.HexEncode([]byte{0x47, 0x6f})
	fmt.Println(encoded)
	// Output: 476f
}

func ExampleBase64RawURLEncode() {
	encoded := vcodec.Base64RawURLEncode([]byte("go?"))
	decoded, _ := vcodec.Base64RawURLDecode(encoded)

	fmt.Println(encoded)
	fmt.Println(string(decoded))
	// Output:
	// Z28_
	// go?
}

func ExampleHexDecodeStr() {
	decoded, err := vcodec.HexDecodeStr("676f")

	fmt.Println(decoded)
	fmt.Println(err)
	// Output:
	// go
	// <nil>
}

func ExampleBase32Encode() {
	encoded := vcodec.Base32Encode([]byte("go"))
	decoded, _ := vcodec.Base32Decode(encoded)

	fmt.Println(encoded)
	fmt.Println(string(decoded))
	// Output:
	// M5XQ====
	// go
}

func ExampleBase58Encode() {
	encoded := vcodec.Base58Encode([]byte("hello world"))
	decoded, _ := vcodec.Base58Decode(encoded)

	fmt.Println(encoded)
	fmt.Println(string(decoded))
	// Output:
	// StV1DL6CwTryKyV
	// hello world
}

func ExampleBase62Encode() {
	encoded := vcodec.Base62Encode([]byte("hello"))
	decoded, _ := vcodec.Base62Decode(encoded)

	fmt.Println(encoded)
	fmt.Println(string(decoded))
	// Output:
	// 7tQLFHz
	// hello
}

func ExampleMorseEncode() {
	encoded, _ := vcodec.MorseEncode("SOS 1")
	decoded, _ := vcodec.MorseDecode(encoded)

	fmt.Println(encoded)
	fmt.Println(decoded)
	// Output:
	// ... --- ... / .----
	// SOS 1
}

func ExampleROT13() {
	encoded := vcodec.ROT13("hello")

	fmt.Println(encoded)
	fmt.Println(vcodec.ROT13(encoded))
	// Output:
	// uryyb
	// hello
}

func ExampleBase64EncodeWithEncoding() {
	encoded := vcodec.Base64EncodeWithEncoding([]byte("go?"), base64.RawURLEncoding)
	fmt.Println(encoded)
	// Output: Z28_
}

func ExampleBase64DecodeWithEncoding() {
	decoded, _ := vcodec.Base64DecodeWithEncoding("Z28_", base64.RawURLEncoding)
	fmt.Println(string(decoded))
	// Output: go?
}

func ExampleBase64EncodeStr() {
	fmt.Println(vcodec.Base64EncodeStr("go"))
	// Output: Z28=
}

func ExampleBase64DecodeStr() {
	decoded, _ := vcodec.Base64DecodeStr("Z28=")
	fmt.Println(decoded)
	// Output: go
}

func ExampleBase64URLEncode() {
	encoded := vcodec.Base64URLEncode([]byte("a/b?c=d"))
	decoded, _ := vcodec.Base64URLDecode(encoded)

	fmt.Println(encoded)
	fmt.Println(string(decoded))
	// Output:
	// YS9iP2M9ZA==
	// a/b?c=d
}

func ExampleBase64URLDecode() {
	decoded, _ := vcodec.Base64URLDecode("YS9iP2M9ZA==")
	fmt.Println(string(decoded))
	// Output: a/b?c=d
}

func ExampleBase64RawURLDecode() {
	decoded, _ := vcodec.Base64RawURLDecode("Z28_")
	fmt.Println(string(decoded))
	// Output: go?
}

func ExampleHexEncodeStr() {
	fmt.Println(vcodec.HexEncodeStr("go"))
	// Output: 676f
}

func ExampleHexDecode() {
	decoded, _ := vcodec.HexDecode("476f")
	fmt.Println(string(decoded))
	// Output: Go
}

func ExampleBase32EncodeWithEncoding() {
	encoded := vcodec.Base32EncodeWithEncoding([]byte("go"), vcodec.Base32HexEncoding)
	fmt.Println(encoded)
	// Output: CTNG====
}

func ExampleBase32Decode() {
	decoded, _ := vcodec.Base32Decode("M5XQ====")
	fmt.Println(string(decoded))
	// Output: go
}

func ExampleBase32DecodeWithEncoding() {
	decoded, _ := vcodec.Base32DecodeWithEncoding("CTNG====", vcodec.Base32HexEncoding)
	fmt.Println(string(decoded))
	// Output: go
}

func ExampleBase58EncodeWithAlphabet() {
	encoded := vcodec.Base58EncodeWithAlphabet([]byte("hello"), vcodec.Base58FlickrAlphabet)
	decoded, _ := vcodec.Base58DecodeWithAlphabet(encoded, vcodec.Base58FlickrAlphabet)

	fmt.Println(encoded)
	fmt.Println(string(decoded))
	// Output:
	// cM8DuyF
	// hello
}

func ExampleBase58Decode() {
	decoded, _ := vcodec.Base58Decode("StV1DL6CwTryKyV")
	fmt.Println(string(decoded))
	// Output: hello world
}

func ExampleBase58DecodeWithAlphabet() {
	decoded, _ := vcodec.Base58DecodeWithAlphabet("cM8DuyF", vcodec.Base58FlickrAlphabet)
	fmt.Println(string(decoded))
	// Output: hello
}

func ExampleBase62Decode() {
	decoded, _ := vcodec.Base62Decode("7tQLFHz")
	fmt.Println(string(decoded))
	// Output: hello
}

func ExampleMorseDecode() {
	decoded, _ := vcodec.MorseDecode("... --- ... / .----")
	fmt.Println(decoded)
	// Output: SOS 1
}

func ExampleROT47() {
	encoded := vcodec.ROT47("Hello!")
	fmt.Println(encoded)
	fmt.Println(vcodec.ROT47(encoded))
	// Output:
	// w6==@P
	// Hello!
}

func ExampleROTN() {
	fmt.Println(vcodec.ROTN("abcXYZ", 2))
	fmt.Println(vcodec.ROTN("abcXYZ", -2))
	// Output:
	// cdeZAB
	// yzaVWX
}
