// Package vpass provides password strength helpers.
//
// It is a thin facade over internal/pass and exposes deterministic local
// password analysis: score, strength bucket, and rule-level signals such as
// character classes, repeated runs, sequential runs, and common weak password
// matches.
package vpass
