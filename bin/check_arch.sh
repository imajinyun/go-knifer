#!/usr/bin/env bash
#
# check_arch.sh enforces go-knifer's architectural conventions in CI.
#
# Rules checked:
#   1. Every public v* package directory has a doc.go.
#   2. Public v* packages do not import each other (production code only).
#   3. Every public v* package imports at least one internal/ implementation
#      path, and every imported internal/ path actually exists.
#
# It relies on the Go toolchain (go list) for accurate import analysis instead
# of fragile text matching, so it transparently handles abbreviated package
# names (vblf -> internal/bloomfilter), pluralized ones (vmap -> internal/maps),
# and subtrees (vhttp -> internal/httpx/http).
#
# Exit code is non-zero when any rule is violated.

set -euo pipefail

cd "$(dirname "$0")/.."

# Resolve this module's path. In some environments (e.g. a go.work workspace)
# `go list -m` prints multiple modules; pick the one for this directory.
MODULE="$(go list -m 2>/dev/null | grep 'go-knifer' | head -n1)"
if [ -z "${MODULE}" ]; then
	echo "ARCH CHECK ERROR: cannot resolve module path via 'go list -m'" >&2
	exit 2
fi
fail=0

err() {
	echo "ARCH VIOLATION: $*" >&2
	fail=1
}

# Collect public package directories (top-level v* dirs containing .go files).
for dir in v*/; do
	pkg="${dir%/}"
	# Skip directories without Go files.
	if ! ls "${pkg}"/*.go >/dev/null 2>&1; then
		continue
	fi

	# Rule 1: doc.go must exist.
	if [ ! -f "${pkg}/doc.go" ]; then
		err "${pkg}: missing doc.go"
	fi

	# Gather this package's production (non-test) imports via the Go toolchain.
	imports="$(go list -f '{{range .Imports}}{{println .}}{{end}}' "./${pkg}")"

	# Rule 2: must not import another public v* package.
	while IFS= read -r imp; do
		[ -z "${imp}" ] && continue
		case "${imp}" in
		"${MODULE}"/v*)
			err "${pkg}: imports another public package ${imp} (v* packages must not depend on each other)"
			;;
		esac
	done <<<"${imports}"

	# Rule 3: must import at least one existing internal/ implementation.
	internal_count=0
	while IFS= read -r imp; do
		[ -z "${imp}" ] && continue
		case "${imp}" in
		"${MODULE}"/internal/*)
			internal_count=$((internal_count + 1))
			rel="${imp#"${MODULE}"/}"
			if [ ! -d "${rel}" ]; then
				err "${pkg}: imports non-existent internal path ${imp}"
			fi
			;;
		esac
	done <<<"${imports}"

	if [ "${internal_count}" -eq 0 ]; then
		err "${pkg}: does not import any internal/ implementation (facade must delegate to internal)"
	fi
done

if [ "${fail}" -ne 0 ]; then
	echo "" >&2
	echo "Architecture check failed. See violations above." >&2
	exit 1
fi

echo "Architecture check passed."
