package vcrypto_test

import (
	"bytes"
	"errors"
	"testing"

	knifer "github.com/imajinyun/knifer-go"
	"github.com/imajinyun/knifer-go/vcrypto"
)

func testPasswordHashOptions() []vcrypto.PasswordHashOption {
	return []vcrypto.PasswordHashOption{
		vcrypto.WithArgon2idMemory(8),
		vcrypto.WithArgon2idIterations(1),
		vcrypto.WithArgon2idParallelism(1),
		vcrypto.WithArgon2idSaltLength(8),
		vcrypto.WithArgon2idKeyLength(16),
		vcrypto.WithPasswordHashRandomOptions(vcrypto.WithRandomReader(bytes.NewReader([]byte("12345678")))),
	}
}

func TestFacadePasswordHashArgon2id(t *testing.T) {
	encoded, err := vcrypto.HashPasswordArgon2id([]byte("correct horse battery staple"), testPasswordHashOptions()...)
	if err != nil {
		t.Fatalf("HashPasswordArgon2id error = %v", err)
	}
	ok, err := vcrypto.VerifyPasswordArgon2id(encoded, []byte("correct horse battery staple"))
	if err != nil || !ok {
		t.Fatalf("VerifyPasswordArgon2id correct = %v, %v", ok, err)
	}
	ok, err = vcrypto.VerifyPasswordArgon2id(encoded, []byte("wrong password"))
	if err != nil || ok {
		t.Fatalf("VerifyPasswordArgon2id mismatch = %v, %v", ok, err)
	}
	info, err := vcrypto.ParsePasswordHash(encoded)
	if err != nil {
		t.Fatalf("ParsePasswordHash error = %v", err)
	}
	if info.Algorithm != "argon2id" || info.Memory != 8 || info.Iterations != 1 || info.Parallelism != 1 {
		t.Fatalf("ParsePasswordHash = %#v", info)
	}
}

func TestFacadePasswordHashErrors(t *testing.T) {
	if _, err := vcrypto.HashPasswordArgon2id(nil, testPasswordHashOptions()...); !errors.Is(err, vcrypto.ErrInvalidPasswordHash) || !errors.Is(err, knifer.ErrCodeInvalidInput) {
		t.Fatalf("HashPasswordArgon2id empty password error = %v", err)
	}
	if ok, err := vcrypto.VerifyPasswordArgon2id("not encoded", []byte("secret")); ok || !errors.Is(err, vcrypto.ErrInvalidPasswordHash) {
		t.Fatalf("VerifyPasswordArgon2id malformed = %v, %v", ok, err)
	}
}
