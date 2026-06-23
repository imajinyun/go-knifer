#!/usr/bin/env bash
#
# check_api_freeze.sh validates v1 API freeze/deprecation governance metadata.

set -euo pipefail

cd "$(dirname "$0")/.."

python3 - <<'PY'
from __future__ import annotations

import json
import sys

errors: list[str] = []


def add_error(message: str) -> None:
    errors.append(message)


with open("docs/api/tools.json", "r", encoding="utf-8") as f:
    tools = json.load(f)

with open("ai-context.json", "r", encoding="utf-8") as f:
    ai_context = json.load(f)

api_freeze = ai_context.get("api_freeze", {})
if not isinstance(api_freeze, dict):
    add_error("ai-context.json api_freeze must be an object")
    api_freeze = {}

if not api_freeze.get("decision_card_required", False):
    add_error("api_freeze.decision_card_required must be true")
if not api_freeze.get("replacement_required_for_deprecation", False):
    add_error("api_freeze.replacement_required_for_deprecation must be true")

allowed_statuses = set(api_freeze.get("allowed_statuses", []))
expected_statuses = {"recommended", "compatibility", "experimental", "deprecated"}
if allowed_statuses != expected_statuses:
    add_error("api_freeze.allowed_statuses must match recommended, compatibility, experimental, deprecated")

decision_cards = api_freeze.get("decision_cards", [])
if not isinstance(decision_cards, list):
    add_error("api_freeze.decision_cards must be a list")
    decision_cards = []
if api_freeze.get("decision_card_required", False) and len(decision_cards) < 5:
    add_error("api_freeze.decision_cards must contain at least five v1 decision cards")

public_facades = {item.get("package") for item in ai_context.get("public_facades", []) if isinstance(item, dict)}
expected_card_ids = {
    "v1-public-api-entry-budget",
    "v1-dynamic-contract-matrix",
    "v1-heavy-dependency-isolation",
    "v1-error-taxonomy",
    "v1-security-threat-model",
}
seen_card_ids: set[str] = set()
for index, item in enumerate(decision_cards):
    if not isinstance(item, dict):
        add_error(f"api_freeze.decision_cards[{index}] must be an object")
        continue
    card_id = item.get("id")
    if not isinstance(card_id, str) or not card_id:
        add_error(f"api_freeze.decision_cards[{index}].id must be non-empty")
    elif card_id in seen_card_ids:
        add_error(f"api_freeze.decision_cards duplicate id {card_id}")
    else:
        seen_card_ids.add(card_id)
    if item.get("status") not in expected_statuses:
        add_error(f"api_freeze.decision_cards[{index}].status must be an allowed API status")
    for field in ("decision", "rationale"):
        if not isinstance(item.get(field), str) or not item[field].strip():
            add_error(f"api_freeze.decision_cards[{index}].{field} must be non-empty")
    packages = item.get("packages")
    if not isinstance(packages, list) or not packages:
        add_error(f"api_freeze.decision_cards[{index}].packages must be non-empty")
    else:
        unknown_packages = sorted(package for package in packages if package != "all" and package not in public_facades)
        if unknown_packages:
            add_error(f"api_freeze.decision_cards[{index}].packages contains unknown facade(s): " + ", ".join(unknown_packages))
    validation = item.get("validation")
    if not isinstance(validation, list) or len(validation) < 2:
        add_error(f"api_freeze.decision_cards[{index}].validation must contain at least two validation entries")

missing_card_ids = sorted(expected_card_ids - seen_card_ids)
if missing_card_ids:
    add_error("api_freeze.decision_cards missing required v1 decision card(s): " + ", ".join(missing_card_ids))

freeze_checks = api_freeze.get("freeze_checks", [])
if not isinstance(freeze_checks, list) or len(freeze_checks) < 4:
    add_error("api_freeze.freeze_checks must document at least four freeze checks")
else:
    freeze_text = " ".join(str(item).lower() for item in freeze_checks)
    for term in ("decision card", "replacement", "snapshot", "tools catalog"):
        if term not in freeze_text:
            add_error(f"api_freeze.freeze_checks must mention {term!r}")

deprecated_functions: list[str] = []
experimental_functions: list[str] = []
for package in tools.get("packages", []):
    package_name = package.get("name", "")
    for fn in package.get("functions", []):
        status = fn.get("status")
        name = f"{package_name}.{fn.get('name')}"
        if status not in expected_statuses:
            add_error(f"{name} has unknown API status {status!r}")
        if status == "deprecated":
            deprecated_functions.append(name)
            synopsis = fn.get("synopsis", "")
            if "Deprecated:" not in synopsis or "Use " not in synopsis:
                add_error(f"{name} is deprecated but synopsis must include 'Deprecated:' and a replacement using 'Use '")
        if status == "experimental":
            experimental_functions.append(name)

if api_freeze.get("v1_candidate", False) and experimental_functions:
    add_error("api_freeze.v1_candidate forbids experimental APIs: " + ", ".join(experimental_functions))

declared_deprecations = api_freeze.get("deprecations", [])
if not isinstance(declared_deprecations, list):
    add_error("api_freeze.deprecations must be a list")
    declared_deprecations = []
declared_deprecated_names = set()
for index, item in enumerate(declared_deprecations):
    if not isinstance(item, dict):
        add_error(f"api_freeze.deprecations[{index}] must be an object")
        continue
    name = item.get("name")
    replacement = item.get("replacement")
    rationale = item.get("rationale")
    if not isinstance(name, str) or not name:
        add_error(f"api_freeze.deprecations[{index}].name must be a non-empty string")
        continue
    declared_deprecated_names.add(name)
    if not isinstance(replacement, str) or not replacement:
        add_error(f"api_freeze.deprecations[{index}].replacement must be a non-empty string")
    if not isinstance(rationale, str) or not rationale:
        add_error(f"api_freeze.deprecations[{index}].rationale must be a non-empty string")

missing_deprecation_entries = sorted(set(deprecated_functions) - declared_deprecated_names)
if missing_deprecation_entries:
    add_error("api_freeze.deprecations missing deprecated function(s): " + ", ".join(missing_deprecation_entries))

if errors:
    for error in errors:
        print(f"api-freeze check error: {error}", file=sys.stderr)
    sys.exit(1)

print(
    "api freeze metadata is valid "
    f"({len(deprecated_functions)} deprecated, {len(experimental_functions)} experimental APIs)"
)
PY
