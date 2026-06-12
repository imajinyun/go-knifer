#!/usr/bin/env bash
#
# check_coverage.sh enforces the repository-wide statement coverage baseline.
# Set COVERAGE_THRESHOLD to raise or lower the required percentage.

set -euo pipefail

coverage_file="${1:-coverage.out}"
threshold="${COVERAGE_THRESHOLD:-69.5}"

if [ ! -f "${coverage_file}" ]; then
	echo "COVERAGE CHECK ERROR: ${coverage_file} does not exist" >&2
	exit 2
fi

total="$(
	go tool cover -func="${coverage_file}" |
		awk '/^total:/ { gsub("%", "", $3); print $3 }'
)"

if [ -z "${total}" ]; then
	echo "COVERAGE CHECK ERROR: cannot read total coverage from ${coverage_file}" >&2
	exit 2
fi

awk -v total="${total}" -v threshold="${threshold}" '
BEGIN {
	if (total + 0 < threshold + 0) {
		printf "coverage %.1f%% is below required %.1f%%\n", total, threshold > "/dev/stderr"
		exit 1
	}
	printf "coverage %.1f%% meets required %.1f%%\n", total, threshold
}
'
