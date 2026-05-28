package vbase

import (
	"io"
	"time"

	"github.com/imajinyun/go-knifer/internal/base"
)

// Ordered is the ordered constraint used by numeric helpers.
type Ordered = base.Ordered

// Number is the numeric constraint used by numeric helpers.
type Number = base.Number

const (
	// NormPattern is the standard datetime layout.
	NormPattern = base.NormPattern
	// DatePattern is the standard date layout.
	DatePattern = base.NormDatePattern
	// TimePattern is the standard time layout.
	TimePattern = base.NormTimePattern
	// BaseNumber is the digit charset.
	BaseNumber = base.BaseNumber
	// BaseChar is the lowercase alphabet charset.
	BaseChar = base.BaseChar
	// BaseCharNumber is the lowercase alphanumeric charset.
	BaseCharNumber = base.BaseCharNumber
	// BaseCharNumberUC is the uppercase alphanumeric charset.
	BaseCharNumberUC = base.BaseCharNumberUC
)

// Format replaces Hutool-style {} placeholders with args.
func Format(template string, args ...any) string { return base.Format(template, args...) }

// IsBlank reports whether s is empty or whitespace.
func IsBlank(s string) bool { return base.IsBlank(s) }

// IsNotBlank reports whether s is not blank.
func IsNotBlank(s string) bool { return base.IsNotBlank(s) }

// IsEmpty reports whether s is empty.
func IsEmpty(s string) bool { return base.IsEmpty(s) }

// IsNotEmpty reports whether s is not empty.
func IsNotEmpty(s string) bool { return base.IsNotEmpty(s) }

// DefaultIfBlank returns def when s is blank.
func DefaultIfBlank(s, def string) string { return base.DefaultIfBlank(s, def) }

// DefaultIfEmpty returns def when s is empty.
func DefaultIfEmpty(s, def string) string { return base.DefaultIfEmpty(s, def) }

// Contains reports whether s contains sub.
func Contains(s, sub string) bool { return base.Contains(s, sub) }

// ContainsIgnoreCase reports whether s contains sub ignoring case.
func ContainsIgnoreCase(s, sub string) bool { return base.ContainsIgnoreCase(s, sub) }

// EqualsIgnoreCase reports whether a equals b ignoring case.
func EqualsIgnoreCase(a, b string) bool { return base.EqualsIgnoreCase(a, b) }

// AddPrefixIfNot adds prefix when s does not already have it.
func AddPrefixIfNot(s, prefix string) string { return base.AddPrefixIfNot(s, prefix) }

// AddSuffixIfNot adds suffix when s does not already have it.
func AddSuffixIfNot(s, suffix string) string { return base.AddSuffixIfNot(s, suffix) }

// Length returns the rune length of s.
func Length(s string) int { return base.Length(s) }

// IsEmail reports whether s is an email address.
func IsEmail(s string) bool { return base.IsEmail(s) }

// IsURL reports whether s is a URL.
func IsURL(s string) bool { return base.IsURL(s) }

// IsIPv4 reports whether s is an IPv4 address.
func IsIPv4(s string) bool { return base.IsIPv4(s) }

// IsNumber reports whether s is numeric.
func IsNumber(s string) bool { return base.IsNumber(s) }

// RandomInt returns a random int in [0, max).
func RandomInt(max int) int { return base.RandomInt(max) }

// RandomString returns a random string with n characters.
func RandomString(n int) string { return base.RandomString(n) }

// RandomStringFromBase returns a random string using charset characters.
func RandomStringFromBase(charset string, n int) string { return base.RandomStringFrom(charset, n) }

// UUID returns a random UUID.
func UUID() string { return base.SimpleUUID() }

// FastUUID returns a random UUID without dashes.
func FastUUID() string { return base.FastUUID() }

// ObjectId returns an ObjectId-like identifier.
func ObjectId() string { return base.ObjectId() }

// NanoId returns a NanoID string.
func NanoId() string { return base.NanoId() }

// MD5Hex computes the MD5 hex digest of s.
func MD5Hex(s string) string { return base.MD5Hex(s) }

// HexEncode encodes data as hex.
func HexEncode(data []byte) string { return base.HexEncode(data) }

// HexDecode decodes a hex string.
func HexDecode(s string) ([]byte, error) { return base.HexDecode(s) }

// Base64Encode encodes data as base64.
func Base64Encode(data []byte) string { return base.Base64Encode(data) }

// Base64Decode decodes a base64 string.
func Base64Decode(s string) ([]byte, error) { return base.Base64Decode(s) }

// Now returns current local time.
func Now() time.Time { return base.Now() }

// FormatDateNorm formats t with NormPattern.
func FormatDateNorm(t time.Time) string { return base.FormatDateNorm(t) }

// ParseDate parses s with layout.
func ParseDate(s string) (time.Time, error) { return base.ParseDate(s) }

// ParseDateLayout parses s with layout.
func ParseDateLayout(s, layout string) (time.Time, error) { return base.ParseDateLayout(s, layout) }

// BeginOfDay returns the beginning of the day.
func BeginOfDay(t time.Time) time.Time { return base.BeginOfDay(t) }

// EndOfDay returns the end of the day.
func EndOfDay(t time.Time) time.Time { return base.EndOfDay(t) }

// BetweenDays returns absolute days between two times.
func BetweenDays(a, b time.Time) int { return base.BetweenDays(a, b) }

// Min returns the minimum value.
func Min[T Ordered](nums ...T) T { return base.Min(nums...) }

// Max returns the maximum value.
func Max[T Ordered](nums ...T) T { return base.Max(nums...) }

// Sum returns the sum of numbers.
func Sum[T Number](nums ...T) T { return base.Sum(nums...) }

// Avg returns the average of numbers.
func Avg[T Number](nums ...T) float64 { return base.Avg(nums...) }

// FileExists reports whether path exists.
func FileExists(path string) bool { return base.FileExists(path) }

// FileReadString reads a whole file as string.
func FileReadString(path string) (string, error) { return base.FileReadString(path) }

// FileWriteString writes content to path.
func FileWriteString(path, content string) error { return base.FileWriteString(path, content) }

// IoCopy copies from src to dst.
func IoCopy(dst io.Writer, src io.Reader) (int64, error) { return base.IoCopy(dst, src) }

// CloseQuietly closes c and ignores errors.
func CloseQuietly(c io.Closer) { base.CloseQuietly(c) }

// SliceContains reports whether item exists in s.
func SliceContains[T comparable](s []T, item T) bool { return base.SliceContains(s, item) }

// SliceDistinct returns deduplicated items.
func SliceDistinct[T comparable](s []T) []T { return base.SliceDistinct(s) }

// Union returns the union of two slices.
func Union[T comparable](a, b []T) []T { return base.Union(a, b) }

// Intersection returns the intersection of two slices.
func Intersection[T comparable](a, b []T) []T { return base.Intersection(a, b) }

// MapKeys returns all keys from m.
func MapKeys[K comparable, V any](m map[K]V) []K { return base.MapKeys(m) }

// MapValues returns all values from m.
func MapValues[K comparable, V any](m map[K]V) []V { return base.MapValues(m) }
