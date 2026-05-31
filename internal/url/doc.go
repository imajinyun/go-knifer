// Package url provides internal URL and URI helpers.
//
// The package centralizes URL parsing, normalization, query handling,
// Data URI building, resource sizing, and scheme checks for public vurl
// facades and other internal packages.
//
// Encoding algorithms such as Base64 and Hex belong in internal/codec. URL query
// escaping may be wrapped by codec for convenience, while URL/URI semantics stay
// in this package.
package url
