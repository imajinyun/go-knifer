#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "${ROOT_DIR}"

if [ ! -f "aiflow.yaml" ]; then
	echo "AIFLOW LAYOUT ERROR: aiflow.yaml must live at the repository root" >&2
	exit 1
fi

for path in ".aiflow/aiflow.yaml" ".aiflow/aiflow.yml" ".aiflow/aiflow.json"; do
	if [ -e "${path}" ]; then
		echo "AIFLOW LAYOUT ERROR: ${path} is a project config inside the local runtime directory" >&2
		echo "Move project config to ./aiflow.yaml; keep .aiflow/ for generated evidence and temporary state only." >&2
		exit 1
	fi
done

if ! git check-ignore -q .aiflow/; then
	echo "AIFLOW LAYOUT ERROR: .aiflow/ must be ignored by Git" >&2
	exit 1
fi

tracked_aiflow="$(git ls-files '.aiflow' '.aiflow/*')"
if [ -n "${tracked_aiflow}" ]; then
	echo "AIFLOW LAYOUT ERROR: .aiflow/ contains tracked files:" >&2
	printf '%s\n' "${tracked_aiflow}" | while IFS= read -r path; do
		echo "  - ${path}" >&2
	done
	echo "Only root aiflow.yaml should be committed; .aiflow/ is local runtime evidence/state." >&2
	exit 1
fi

echo "aiflow layout is valid"
