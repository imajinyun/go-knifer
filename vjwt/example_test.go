package vjwt_test

import (
	"errors"
	"fmt"

	"github.com/imajinyun/go-knifer"
	"github.com/imajinyun/go-knifer/vjwt"
)

func ExampleNewJWTError() {
	err := vjwt.NewJWTError("token must not be blank")
	fmt.Println(errors.Is(err, knifer.ErrCodeInvalidInput))
	// Output: true
}

func ExampleCreateTokenWithOptions() {
	token, err := vjwt.CreateTokenWithOptions(
		vjwt.WithTokenHeaders(map[string]any{vjwt.JWTHeaderAlgorithm: vjwt.JWTAlgHS256}),
		vjwt.WithTokenPayload(map[string]any{vjwt.JWTPayloadSubject: "alice"}),
		vjwt.WithTokenKey([]byte("secret")),
	)
	parsed, parseErr := vjwt.ParseToken(token)
	fmt.Println(err == nil, parseErr == nil, parsed.Payload(vjwt.JWTPayloadSubject))
	// Output: true true alice
}
