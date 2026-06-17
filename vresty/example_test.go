package vresty_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/imajinyun/go-knifer/vresty"
)

func ExampleGetStringSafeE() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("safe response"))
	}))
	defer server.Close()

	body, err := vresty.GetStringSafeE(server.URL,
		vresty.WithURLPolicy(vresty.URLPolicy{AllowedSchemes: []string{"http", "https"}, RejectPrivate: false}),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(body)
	// Output: safe response
}
