package crypto

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	knifer "github.com/imajinyun/knifer-go"
)

func testPasswordHashOptions() []PasswordHashOption {
	return []PasswordHashOption{
		WithArgon2idMemory(8),
		WithArgon2idIterations(1),
		WithArgon2idParallelism(1),
		WithArgon2idSaltLength(8),
		WithArgon2idKeyLength(16),
		WithPasswordHashRandomOptions(WithRandomReader(bytes.NewReader([]byte("12345678")))),
	}
}

func TestHashPasswordArgon2idRoundTrip(t *testing.T) {
	encoded, err := HashPasswordArgon2id([]byte("correct horse battery staple"), testPasswordHashOptions()...)
	if err != nil {
		t.Fatalf("HashPasswordArgon2id error = %v", err)
	}
	if !strings.HasPrefix(encoded, "$argon2id$v=19$m=8,t=1,p=1$") {
		t.Fatalf("encoded hash = %q", encoded)
	}
	ok, err := VerifyPasswordArgon2id(encoded, []byte("correct horse battery staple"))
	if err != nil || !ok {
		t.Fatalf("VerifyPasswordArgon2id correct = %v, %v", ok, err)
	}
	ok, err = VerifyPasswordArgon2id(encoded, []byte("wrong password"))
	if err != nil || ok {
		t.Fatalf("VerifyPasswordArgon2id mismatch = %v, %v", ok, err)
	}
}

func TestParsePasswordHash(t *testing.T) {
	encoded, err := HashPasswordArgon2id([]byte("secret"), testPasswordHashOptions()...)
	if err != nil {
		t.Fatalf("HashPasswordArgon2id error = %v", err)
	}
	info, err := ParsePasswordHash(encoded)
	if err != nil {
		t.Fatalf("ParsePasswordHash error = %v", err)
	}
	if info.Algorithm != "argon2id" || info.Version != 19 || info.Memory != 8 || info.Iterations != 1 || info.Parallelism != 1 || info.SaltLength != 8 || info.KeyLength != 16 {
		t.Fatalf("ParsePasswordHash = %#v", info)
	}
}

func TestPasswordHashErrors(t *testing.T) {
	_, err := HashPasswordArgon2id(nil, testPasswordHashOptions()...)
	if !errors.Is(err, ErrInvalidPasswordHash) || !errors.Is(err, knifer.ErrCodeInvalidInput) {
		t.Fatalf("HashPasswordArgon2id empty password error = %v", err)
	}
	_, err = HashPasswordArgon2id([]byte("secret"), WithArgon2idMemory(7))
	if !errors.Is(err, ErrInvalidPasswordHash) {
		t.Fatalf("HashPasswordArgon2id invalid memory error = %v", err)
	}
	_, err = HashPasswordArgon2id([]byte("secret"), WithArgon2idMemory(8), WithArgon2idIterations(1), WithArgon2idParallelism(1), WithArgon2idSaltLength(8), WithArgon2idKeyLength(16), WithPasswordHashRandomOptions(WithRandomReader(bytes.NewReader([]byte("short")))))
	if !errors.Is(err, knifer.ErrCodeProviderFailure) {
		t.Fatalf("HashPasswordArgon2id short salt reader error = %v", err)
	}
	if ok, err := VerifyPasswordArgon2id("$argon2id$v=19$m=8,t=1,p=1$bad$hash", []byte("secret")); ok || !errors.Is(err, ErrInvalidPasswordHash) {
		t.Fatalf("VerifyPasswordArgon2id malformed = %v, %v", ok, err)
	}
	if _, err := ParsePasswordHash("$bcrypt$v=1$bad$hash"); !errors.Is(err, ErrInvalidPasswordHash) {
		t.Fatalf("ParsePasswordHash unsupported envelope error = %v", err)
	}
}
