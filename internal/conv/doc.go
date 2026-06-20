// Package conv provides permissive type conversion helpers.
//
// It converts values to string, int, int64, float64, bool, and []byte with
// lenient fallback semantics. Zero-value helpers return zero on failure,
// Default helpers return caller-provided fallbacks, and E helpers return a
// classified error instead of swallowing invalid scalar input.
package conv
