#!/usr/bin/env bash

set -euo pipefail

usage() {
  cat <<'EOF'
Usage: scripts/regenerate_compare_push.sh QUESTION [DATE]

Regenerate compare artifacts for one question and push them to the
ExploringDatabyLLMs-runs repo.

Arguments:
  QUESTION   Question id like q002 or q003
  DATE       Optional day in YYYY-MM-DD format (defaults to today)

Environment:
  QFORGE_RUNNER   Compare runner to use (default: claude)
  QFORGE_MODEL    Optional compare model override
  QFORGE_VERBOSE  Set to 1 to pass -v to qforge compare
EOF
}

if [[ $# -lt 1 || $# -gt 2 ]]; then
  usage >&2
  exit 1
fi

question="$1"
day="${2:-$(date +%F)}"
runner="${QFORGE_RUNNER:-claude}"
model="${QFORGE_MODEL:-}"
verbose="${QFORGE_VERBOSE:-0}"

code_repo="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
runs_repo="$(cd "$code_repo/../ExploringDatabyLLMs-runs" && pwd)"

if [[ ! -x "$code_repo/scripts/qforge" ]]; then
  echo "qforge script not found at $code_repo/scripts/qforge" >&2
  exit 1
fi

question_dir_glob="$runs_repo/$day/${question}"_*
question_dir="$(find "$runs_repo/$day" -maxdepth 1 -type d -name "${question}_*" | sort | head -n 1 || true)"

if [[ -z "$question_dir" ]]; then
  echo "No question directory found for $question on $day under $runs_repo/$day" >&2
  echo "Looked for: $question_dir_glob" >&2
  exit 1
fi

compare_args=(
  "$code_repo/scripts/qforge"
  compare
  --day "$day"
  --question "$question"
  --runner "$runner"
)

if [[ -n "$model" ]]; then
  compare_args+=(--model "$model")
fi

if [[ "$verbose" == "1" ]]; then
  compare_args+=(-v)
fi

echo "[compare-push] code_repo=$code_repo"
echo "[compare-push] runs_repo=$runs_repo"
echo "[compare-push] question=$question day=$day runner=$runner"
if [[ -n "$model" ]]; then
  echo "[compare-push] model=$model"
fi

(cd "$code_repo" && "${compare_args[@]}")

artifacts=(
  "$question_dir/compare/analysis.prompt.md"
  "$question_dir/compare/analysis.raw.md"
  "$question_dir/compare/compare.json"
  "$question_dir/compare/provider.raw.md"
  "$question_dir/compare_report.md"
)

existing_artifacts=()
for path in "${artifacts[@]}"; do
  if [[ -e "$path" ]]; then
    existing_artifacts+=("${path#$runs_repo/}")
  fi
done

if [[ ${#existing_artifacts[@]} -eq 0 ]]; then
  echo "No compare artifacts found to stage under $question_dir" >&2
  exit 1
fi

cd "$runs_repo"

git add -- "${existing_artifacts[@]}"

if git diff --cached --quiet; then
  echo "[compare-push] no staged changes for $question on $day"
  exit 0
fi

git commit -m "Refresh ${question} compare artifacts for ${day}"
git pull --rebase --autostash origin main
git push origin main

echo "[compare-push] pushed:"
printf '  %s\n' "${existing_artifacts[@]}"
