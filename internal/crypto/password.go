package crypto

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	knifer "github.com/imajinyun/knifer-go"
	"golang.org/x/crypto/argon2"
)

const (
	passwordHashAlgorithmArgon2id = "argon2id"

	defaultArgon2idMemory      uint32 = 64 * 1024
	defaultArgon2idIterations  uint32 = 3
	defaultArgon2idParallelism uint8  = 4
	defaultArgon2idSaltLength  uint32 = 16
	defaultArgon2idKeyLength   uint32 = 32
)

// PasswordHashInfo describes a parsed encoded password hash.
type PasswordHashInfo struct {
	Algorithm   string
	Version     int
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  int
	KeyLength   int
}

type passwordHashConfig struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
	random      []RandomOption
}

// PasswordHashOption customizes password hashing helpers.
type PasswordHashOption func(*passwordHashConfig)

// WithArgon2idMemory sets the Argon2id memory cost in KiB.
func WithArgon2idMemory(memory uint32) PasswordHashOption {
	return func(c *passwordHashConfig) { c.memory = memory }
}

// WithArgon2idIterations sets the Argon2id iteration count.
func WithArgon2idIterations(iterations uint32) PasswordHashOption {
	return func(c *passwordHashConfig) { c.iterations = iterations }
}

// WithArgon2idParallelism sets the Argon2id parallelism value.
func WithArgon2idParallelism(parallelism uint8) PasswordHashOption {
	return func(c *passwordHashConfig) { c.parallelism = parallelism }
}

// WithArgon2idSaltLength sets the generated salt length.
func WithArgon2idSaltLength(length uint32) PasswordHashOption {
	return func(c *passwordHashConfig) { c.saltLength = length }
}

// WithArgon2idKeyLength sets the derived key length.
func WithArgon2idKeyLength(length uint32) PasswordHashOption {
	return func(c *passwordHashConfig) { c.keyLength = length }
}

// WithPasswordHashRandomOptions sets entropy source options for generated salts.
func WithPasswordHashRandomOptions(opts ...RandomOption) PasswordHashOption {
	return func(c *passwordHashConfig) { c.random = append([]RandomOption(nil), opts...) }
}

func applyPasswordHashOptions(opts []PasswordHashOption) passwordHashConfig {
	cfg := passwordHashConfig{
		memory:      defaultArgon2idMemory,
		iterations:  defaultArgon2idIterations,
		parallelism: defaultArgon2idParallelism,
		saltLength:  defaultArgon2idSaltLength,
		keyLength:   defaultArgon2idKeyLength,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return cfg
}

func validatePasswordHashConfig(cfg passwordHashConfig) error {
	if cfg.memory < 8 {
		return knifer.WrapError(knifer.ErrCodeInvalidInput, "argon2id memory must be at least 8 KiB", ErrInvalidPasswordHash)
	}
	if cfg.iterations == 0 {
		return knifer.WrapError(knifer.ErrCodeInvalidInput, "argon2id iterations must be positive", ErrInvalidPasswordHash)
	}
	if cfg.parallelism == 0 {
		return knifer.WrapError(knifer.ErrCodeInvalidInput, "argon2id parallelism must be positive", ErrInvalidPasswordHash)
	}
	if cfg.saltLength < 8 {
		return knifer.WrapError(knifer.ErrCodeInvalidInput, "argon2id salt length must be at least 8 bytes", ErrInvalidPasswordHash)
	}
	if cfg.keyLength < 16 {
		return knifer.WrapError(knifer.ErrCodeInvalidInput, "argon2id key length must be at least 16 bytes", ErrInvalidPasswordHash)
	}
	return nil
}

// HashPasswordArgon2id hashes password using Argon2id and returns an encoded hash envelope.
func HashPasswordArgon2id(password []byte, opts ...PasswordHashOption) (string, error) {
	if len(password) == 0 {
		return "", knifer.WrapError(knifer.ErrCodeInvalidInput, "password must not be empty", ErrInvalidPasswordHash)
	}
	cfg := applyPasswordHashOptions(opts)
	if err := validatePasswordHashConfig(cfg); err != nil {
		return "", err
	}
	salt, err := RandomBytesWithOptions(int(cfg.saltLength), cfg.random...)
	if err != nil {
		return "", knifer.WrapError(knifer.ErrCodeProviderFailure, "generate password hash salt", err)
	}
	hash := argon2.IDKey(password, salt, cfg.iterations, cfg.memory, cfg.parallelism, cfg.keyLength)
	return encodeArgon2idHash(cfg, salt, hash), nil
}

// VerifyPasswordArgon2id verifies password against an Argon2id encoded hash.
func VerifyPasswordArgon2id(encoded string, password []byte) (bool, error) {
	parsed, salt, expected, err := parseArgon2id(encoded)
	if err != nil {
		return false, err
	}
	if len(password) == 0 {
		return false, nil
	}
	actual := argon2.IDKey(password, salt, parsed.Iterations, parsed.Memory, parsed.Parallelism, uint32(parsed.KeyLength))
	return subtle.ConstantTimeCompare(actual, expected) == 1, nil
}

// ParsePasswordHash parses a supported encoded password hash without verifying a password.
func ParsePasswordHash(encoded string) (PasswordHashInfo, error) {
	info, _, _, err := parseArgon2id(encoded)
	return info, err
}

func encodeArgon2idHash(cfg passwordHashConfig, salt, hash []byte) string {
	enc := base64.RawStdEncoding
	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		cfg.memory,
		cfg.iterations,
		cfg.parallelism,
		enc.EncodeToString(salt),
		enc.EncodeToString(hash),
	)
}

func parseArgon2id(encoded string) (PasswordHashInfo, []byte, []byte, error) {
	parts := strings.Split(strings.TrimSpace(encoded), "$")
	if len(parts) != 6 || parts[0] != "" || parts[1] != passwordHashAlgorithmArgon2id {
		return PasswordHashInfo{}, nil, nil, knifer.WrapError(knifer.ErrCodeInvalidInput, "password hash must be argon2id encoded envelope", ErrInvalidPasswordHash)
	}
	versionText := strings.TrimPrefix(parts[2], "v=")
	version, err := strconv.Atoi(versionText)
	if err != nil || version != argon2.Version {
		return PasswordHashInfo{}, nil, nil, knifer.WrapError(knifer.ErrCodeInvalidInput, "password hash has unsupported argon2 version", ErrInvalidPasswordHash)
	}
	memory, iterations, parallelism, err := parseArgon2idParams(parts[3])
	if err != nil {
		return PasswordHashInfo{}, nil, nil, err
	}
	enc := base64.RawStdEncoding
	salt, err := enc.DecodeString(parts[4])
	if err != nil || len(salt) < 8 {
		return PasswordHashInfo{}, nil, nil, knifer.WrapError(knifer.ErrCodeInvalidInput, "password hash salt is malformed", ErrInvalidPasswordHash)
	}
	hash, err := enc.DecodeString(parts[5])
	if err != nil || len(hash) < 16 {
		return PasswordHashInfo{}, nil, nil, knifer.WrapError(knifer.ErrCodeInvalidInput, "password hash value is malformed", ErrInvalidPasswordHash)
	}
	info := PasswordHashInfo{
		Algorithm:   passwordHashAlgorithmArgon2id,
		Version:     version,
		Memory:      memory,
		Iterations:  iterations,
		Parallelism: parallelism,
		SaltLength:  len(salt),
		KeyLength:   len(hash),
	}
	return info, salt, hash, nil
}

func parseArgon2idParams(raw string) (uint32, uint32, uint8, error) {
	var memory uint64
	var iterations uint64
	var parallelism uint64
	for _, part := range strings.Split(raw, ",") {
		key, value, ok := strings.Cut(part, "=")
		if !ok {
			return 0, 0, 0, knifer.WrapError(knifer.ErrCodeInvalidInput, "password hash parameters are malformed", ErrInvalidPasswordHash)
		}
		n, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return 0, 0, 0, knifer.WrapError(knifer.ErrCodeInvalidInput, "password hash parameter value is malformed", ErrInvalidPasswordHash)
		}
		switch key {
		case "m":
			memory = n
		case "t":
			iterations = n
		case "p":
			parallelism = n
		default:
			return 0, 0, 0, knifer.WrapError(knifer.ErrCodeInvalidInput, "password hash contains unknown parameter", ErrInvalidPasswordHash)
		}
	}
	if memory < 8 || iterations == 0 || parallelism == 0 || parallelism > 255 {
		return 0, 0, 0, knifer.WrapError(knifer.ErrCodeInvalidInput, "password hash parameters are out of range", ErrInvalidPasswordHash)
	}
	return uint32(memory), uint32(iterations), uint8(parallelism), nil
}
