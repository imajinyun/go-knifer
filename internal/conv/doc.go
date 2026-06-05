// Package conv provides permissive type conversion helpers.
//
// It converts values between common Go types (string, int, float, bool,
// time, slice, map) with lenient fallback semantics: failed conversions
// return zero values or caller-provided defaults instead of panicking.
// This package is exposed through the vconv facade.
package conv
