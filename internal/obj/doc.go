// Package obj provides internal object helpers.
//
// The package centralizes nil checks, emptiness checks, default-value helpers,
// comparison, cloning, serialization, type inspection, and generic container
// helpers for the public vobj facade.
//
// This package is intentionally a convenience object-level facade, not a place
// for new domain logic. Prefer implementing concrete behavior in focused
// packages such as str, slice, maps, serialize, ref, or conv first; add thin obj
// wrappers only when an object-level helper improves ergonomics.
package obj
