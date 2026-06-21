# vcrypto Quickstart

`vcrypto` provides common cryptographic helpers, including digests, HMAC, AES-GCM, random bytes, PBKDF2, RSA encryption/decryption/signing, and PEM conversion.

## Recommended API families

| Need | Recommended API family | Notes |
| --- | --- | --- |
| Hash non-secret data | SHA-256/SHA-512 helpers | Do not use plain hashes as authentication codes. |
| Authenticate a message | HMAC helpers | Compare outputs using constant-time checks where exposed. |
| Encrypt bytes | AES-GCM helpers | Use fresh nonces and authenticated encryption. |
| Generate secrets | `vrand` secure helpers | Do not use `math/rand` for secrets. |

## Misuse checklist

- Do not use MD5 or SHA-1 for new security-sensitive designs.
- Do not reuse AES-GCM nonces with the same key.
- Do not generate keys, tokens, nonces, or salts with `math/rand`.
- Do not log secret bytes, private keys, raw tokens, or derived credentials.
- Do not ignore crypto errors; treat them as security-relevant failures.

## Benchmarks

Measure crypto helper overhead locally with the focused benchmark suite:

```bash
go test -bench=. -benchmem -run=^$ ./internal/crypto ./vcrypto ./internal/rand ./vrand
```

The suite covers SHA-256 digest, HMAC-SHA256 signing, AES-GCM encrypt/decrypt, and secure random byte generation. Treat benchmark output as a local baseline, not a universal performance claim.

## FAQ

### Does go-knifer replace Go's `crypto/*` standard library packages?

No. It provides focused helper entry points and documents safe defaults for common workflows. Use the standard library directly when you need low-level control.

### Which APIs should I avoid for security-sensitive code?

Avoid legacy digest algorithms such as MD5 or SHA-1 for new security-sensitive designs. Prefer documented recommended APIs.

### Are secrets ever logged?

Security-sensitive helpers must not log raw secrets, tokens, keys, nonces, or salts. Treat any such behavior as a security bug.

## SHA and HMAC

```go
package main

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vcrypto"
)

func main() {
	digest := vcrypto.SHA256Hex("hello")
	mac := vcrypto.HMACSHA256Hex([]byte("secret"), []byte("hello"))

	fmt.Println(digest)
	fmt.Println(mac)
}
```

When the hash factory is `nil`, `HMACHex` and `HMACBytes` fall back to SHA-256 instead of panicking. Prefer passing an explicit hash when interoperability matters.

## AES-GCM encryption and decryption

```go
package main

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vcrypto"
)

func main() {
	key, err := vcrypto.GenAESKey(32)
	if err != nil {
		panic(err)
	}

	nonce, cipherText, err := vcrypto.AESSealGCM([]byte("secret data"), key, []byte("aad"))
	if err != nil {
		panic(err)
	}

	plain, err := vcrypto.AESOpenGCM(cipherText, key, nonce, []byte("aad"))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(plain))
}
```

Authentication failures from `AESOpenGCM` / `AESDecryptGCM` match `vcrypto.ErrInvalidCipherText`, so callers can distinguish tampering or wrong AAD from nonce-length validation errors.

## Derive keys with PBKDF2

```go
package main

import (
	"crypto/sha256"
	"fmt"

	"github.com/imajinyun/go-knifer/vcrypto"
)

func main() {
	key, err := vcrypto.PBKDF2([]byte("password"), []byte("salt"), 10000, 32, sha256.New)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(key))
}
```

## RSA signing and verification

```go
package main

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vcrypto"
)

func main() {
	priv, err := vcrypto.GenRSAKey(2048)
	if err != nil {
		panic(err)
	}

	data := []byte("message")
	sig, err := vcrypto.SignSHA256WithRSA(data, priv)
	if err != nil {
		panic(err)
	}

	err = vcrypto.VerifySHA256WithRSA(data, sig, &priv.PublicKey)
	fmt.Println(err == nil)
}
```
