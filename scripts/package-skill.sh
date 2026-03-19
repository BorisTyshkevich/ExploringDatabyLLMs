#!/usr/bin/env bash

set -euo pipefail

usage() {
  cat >&2 <<'EOF'
Usage: scripts/package-skill.sh <skill-dir-name>

Package Skills/<skill-dir-name> as a zip archive at .tmp/skill-zips/<skill-dir-name>.zip
and print the artifact path to stdout.
EOF
}

if [[ $# -ne 1 ]]; then
  usage
  exit 1
fi

skill_name="$1"

if [[ -z "$skill_name" ]]; then
  echo "skill name must not be empty" >&2
  usage
  exit 1
fi

if [[ "$skill_name" == /* || "$skill_name" == *"/"* || "$skill_name" == "." || "$skill_name" == ".." ]]; then
  echo "skill name must be a directory name relative to Skills/" >&2
  exit 1
fi

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
skills_root="$repo_root/Skills"
skill_dir="$skills_root/$skill_name"
output_dir="$repo_root/.tmp/skill-zips"
output_path="$output_dir/$skill_name.zip"

if ! command -v zip >/dev/null 2>&1; then
  echo "zip command not found in PATH" >&2
  exit 1
fi

if [[ ! -d "$skill_dir" ]]; then
  echo "skill directory not found: $skill_dir" >&2
  exit 1
fi

mkdir -p "$output_dir"
rm -f "$output_path"

(
  cd "$skills_root"
  zip -qr "$output_path" "$skill_name"
)

echo "$output_path"
