#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
AI_CONTEXT="${ROOT_DIR}/ai-context.json"

python3 - "${ROOT_DIR}" "${AI_CONTEXT}" <<'PY'
import json
import os
import re
import sys

root_dir, ai_context = sys.argv[1], sys.argv[2]
errors = []


def add_error(message):
    errors.append(message)


def require_mapping(value, path):
    if not isinstance(value, dict):
        add_error(f"{path} must be an object")
        return {}
    return value


def require_string(value, path):
    if not isinstance(value, str) or not value.strip():
        add_error(f"{path} must be a non-empty string")
        return ""
    return value


def require_string_list(value, path):
    if not isinstance(value, list):
        add_error(f"{path} must be a list")
        return []
    result = []
    for index, item in enumerate(value):
        item = require_string(item, f"{path}[{index}]")
        if item:
            result.append(item)
    return result


def workflow_job_names(workflow_text):
    in_jobs = False
    names = set()
    for line in workflow_text.splitlines():
        if re.match(r"^jobs:\s*$", line):
            in_jobs = True
            continue
        if in_jobs and line and not line.startswith(" "):
            break
        if in_jobs:
            match = re.match(r"^  ([A-Za-z0-9_-]+):\s*$", line)
            if match:
                names.add(match.group(1))
    return names


try:
    with open(ai_context, "r", encoding="utf-8") as f:
        data = json.load(f)
except FileNotFoundError:
    print("missing ai-context.json", file=sys.stderr)
    sys.exit(1)
except json.JSONDecodeError as exc:
    print(f"invalid ai-context.json: {exc}", file=sys.stderr)
    sys.exit(1)

ci_workflows = require_mapping(data.get("ci_workflows"), "ci_workflows")
tool_versions = require_mapping(ci_workflows.get("tool_versions"), "ci_workflows.tool_versions")
go_1_25_patch = require_string(tool_versions.get("go_1_25_patch"), "ci_workflows.tool_versions.go_1_25_patch")
golangci_lint_version = require_string(tool_versions.get("golangci_lint"), "ci_workflows.tool_versions.golangci_lint")
github_actions = require_mapping(ci_workflows.get("github_actions"), "ci_workflows.github_actions")

workflow_dir = os.path.join(root_dir, ".github", "workflows")
declared_paths = set()
actual_paths = set()
if os.path.isdir(workflow_dir):
    for name in os.listdir(workflow_dir):
        if name.endswith((".yml", ".yaml")):
            actual_paths.add(os.path.join(".github", "workflows", name))

for name, workflow in sorted(github_actions.items()):
    workflow = require_mapping(workflow, f"ci_workflows.github_actions.{name}")
    workflow_path = require_string(workflow.get("path"), f"ci_workflows.github_actions.{name}.path")
    required_jobs = require_string_list(
        workflow.get("required_jobs"),
        f"ci_workflows.github_actions.{name}.required_jobs",
    )
    agent_governance = require_mapping(
        workflow.get("agent_governance"),
        f"ci_workflows.github_actions.{name}.agent_governance",
    )
    required_commands = require_string_list(
        agent_governance.get("required_commands"),
        f"ci_workflows.github_actions.{name}.agent_governance.required_commands",
    )
    required_env = require_string_list(
        agent_governance.get("required_env"),
        f"ci_workflows.github_actions.{name}.agent_governance.required_env",
    )
    required_artifacts = require_string_list(
        agent_governance.get("required_artifacts"),
        f"ci_workflows.github_actions.{name}.agent_governance.required_artifacts",
    )

    if not workflow_path:
        continue
    declared_paths.add(workflow_path)
    absolute_workflow_path = os.path.join(root_dir, workflow_path)
    if not os.path.exists(absolute_workflow_path):
        add_error(f"ci_workflows.github_actions.{name}.path references missing workflow {workflow_path!r}")
        continue

    with open(absolute_workflow_path, "r", encoding="utf-8") as f:
        workflow_text = f.read()

    jobs = workflow_job_names(workflow_text)
    for job in required_jobs:
        if job not in jobs:
            add_error(f"{workflow_path} is missing required job {job!r}")

    for required_text in required_commands:
        if required_text not in workflow_text:
            add_error(f"{workflow_path} must contain required command {required_text!r}")
    for env_name in required_env:
        if env_name not in workflow_text:
            add_error(f"{workflow_path} must contain required env {env_name!r}")
    for artifact in required_artifacts:
        if artifact not in workflow_text:
            add_error(f"{workflow_path} must upload required artifact {artifact!r}")

    if name in {"go", "release"}:
        if "GO_1_25_PATCH_VERSION" not in workflow_text or go_1_25_patch not in workflow_text:
            add_error(f"{workflow_path} must use declared Go patch version {go_1_25_patch!r}")
    if name == "go":
        if "GOLANGCI_LINT_VERSION" not in workflow_text or golangci_lint_version not in workflow_text:
            add_error(f"{workflow_path} must use declared golangci-lint version {golangci_lint_version!r}")
        if go_1_25_patch and f'go-version: ["{go_1_25_patch}", "1.26"]' not in workflow_text:
            add_error(f"{workflow_path} test matrix must include {go_1_25_patch!r} and '1.26'")
        for duplicate_step in ("make race-test", "make shuffle-test", "make mod-check"):
            if duplicate_step in workflow_text:
                add_error(f"{workflow_path} should not duplicate ci-test sub-step {duplicate_step!r}")

undeclared = sorted(actual_paths - declared_paths)
missing_files = sorted(declared_paths - actual_paths)
if undeclared:
    add_error("undeclared GitHub workflow file(s): " + ", ".join(undeclared))
if missing_files:
    add_error("declared GitHub workflow path(s) not present on disk: " + ", ".join(missing_files))

if errors:
    print("CI WORKFLOW CHECK FAILED:", file=sys.stderr)
    for error in errors:
        print(f"- {error}", file=sys.stderr)
    sys.exit(1)

print(f"CI workflow governance is valid ({len(github_actions)} workflows, {len(actual_paths)} files)")
PY
