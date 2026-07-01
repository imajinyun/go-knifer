package vcrypto

import cryptoimpl "github.com/imajinyun/knifer-go/internal/crypto"

// PasswordHashInfo describes a parsed encoded password hash.
type PasswordHashInfo = cryptoimpl.PasswordHashInfo

// PasswordHashOption customizes password hashing helpers.
type PasswordHashOption = cryptoimpl.PasswordHashOption

// WithArgon2idMemory sets the Argon2id memory cost in KiB.
func WithArgon2idMemory(memory uint32) PasswordHashOption {
	return cryptoimpl.WithArgon2idMemory(memory)
}

// WithArgon2idIterations sets the Argon2id iteration count.
func WithArgon2idIterations(iterations uint32) PasswordHashOption {
	return cryptoimpl.WithArgon2idIterations(iterations)
}

// WithArgon2idParallelism sets the Argon2id parallelism value.
func WithArgon2idParallelism(parallelism uint8) PasswordHashOption {
	return cryptoimpl.WithArgon2idParallelism(parallelism)
}

// WithArgon2idSaltLength sets the generated salt length.
func WithArgon2idSaltLength(length uint32) PasswordHashOption {
	return cryptoimpl.WithArgon2idSaltLength(length)
}

// WithArgon2idKeyLength sets the derived key length.
func WithArgon2idKeyLength(length uint32) PasswordHashOption {
	return cryptoimpl.WithArgon2idKeyLength(length)
}

// WithPasswordHashRandomOptions sets entropy source options for generated salts.
func WithPasswordHashRandomOptions(opts ...RandomOption) PasswordHashOption {
	return cryptoimpl.WithPasswordHashRandomOptions(opts...)
}

// HashPasswordArgon2id hashes password using Argon2id and returns an encoded hash envelope.
func HashPasswordArgon2id(password []byte, opts ...PasswordHashOption) (string, error) {
	return cryptoimpl.HashPasswordArgon2id(password, opts...)
}

// VerifyPasswordArgon2id verifies password against an Argon2id encoded hash.
func VerifyPasswordArgon2id(encoded string, password []byte) (bool, error) {
	return cryptoimpl.VerifyPasswordArgon2id(encoded, password)
}

// ParsePasswordHash parses a supported encoded password hash without verifying a password.
func ParsePasswordHash(encoded string) (PasswordHashInfo, error) {
	return cryptoimpl.ParsePasswordHash(encoded)
}
