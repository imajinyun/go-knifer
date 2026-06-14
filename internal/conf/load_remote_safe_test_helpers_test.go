package conf

import "net/http"

type confRoundTripperFunc func(*http.Request) (*http.Response, error)

func (f confRoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
