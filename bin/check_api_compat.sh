#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SNAPSHOT="${ROOT_DIR}/docs/api/exports.txt"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "${TMP_DIR}"' EXIT

CURRENT="${TMP_DIR}/exports.txt"

(cd "${ROOT_DIR}" && go run ./bin/api_snapshot.go) >"${CURRENT}"

if [[ "${UPDATE_API:-}" == "1" ]]; then
	mkdir -p "$(dirname "${SNAPSHOT}")"
	cp "${CURRENT}" "${SNAPSHOT}"
	echo "updated ${SNAPSHOT#"${ROOT_DIR}/"}"
	exit 0
fi

if [[ ! -f "${SNAPSHOT}" ]]; then
	echo "missing ${SNAPSHOT#"${ROOT_DIR}/"}; run UPDATE_API=1 make api-check" >&2
	exit 1
fi

if ! diff -u "${SNAPSHOT}" "${CURRENT}"; then
	echo "exported API snapshot is stale; run UPDATE_API=1 make api-check after intentional API changes" >&2
	exit 1
fi

echo "exported API snapshot is current"
