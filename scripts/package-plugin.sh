#!/usr/bin/env bash

set -euo pipefail

usage() {
  cat >&2 <<'EOF'
Usage: scripts/package-plugin.sh

Package the repo-root Codex plugin as a zip archive at
.tmp/plugin-zips/exploring-databyllms-plugin.zip and print the artifact path
to stdout.
EOF
}

if [[ $# -ne 0 ]]; then
  usage
  exit 1
fi

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
plugin_manifest="$repo_root/.codex-plugin/plugin.json"
output_dir="$repo_root/.tmp/plugin-zips"
output_path="$output_dir/exploring-databyllms-plugin.zip"

if ! command -v zip >/dev/null 2>&1; then
  echo "zip command not found in PATH" >&2
  exit 1
fi

if [[ ! -f "$plugin_manifest" ]]; then
  echo "plugin manifest not found: $plugin_manifest" >&2
  exit 1
fi

mkdir -p "$output_dir"
rm -f "$output_path"

(
  cd "$repo_root"
  zip -qr "$output_path" .codex-plugin Skills assets README.md
)

echo "$output_path"
