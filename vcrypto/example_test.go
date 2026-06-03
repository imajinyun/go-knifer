package vcrypto_test

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vcrypto"
)

func ExampleMD5Hex() {
	fmt.Println(vcrypto.MD5Hex("hello"))
	// Output: 5d41402abc4b2a76b9719d911017c592
}

func ExampleSHA256Hex() {
	fmt.Println(vcrypto.SHA256Hex("abc"))
	// Output: ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad
}

func ExampleHMACSHA256Hex() {
	fmt.Println(vcrypto.HMACSHA256Hex([]byte("key"), []byte("data")))
	// Output: 5031fe3d989c6d1537a013fa6e739da23463fdaec3b70137d828e36ace221bd0
}
