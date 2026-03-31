#!/usr/bin/env bash

set -euo pipefail

qforge_run_dir_wrapper() {
  if [[ $# -lt 3 ]]; then
    echo "qforge_run_dir_wrapper requires <command> <error-label> <usage-fn> [args...]" >&2
    return 2
  fi

  local command="$1"
  local error_label="$2"
  local usage_fn="$3"
  shift 3

  local script_dir repo_root runs_root
  script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  repo_root="$(cd "${script_dir}/.." && pwd)"
  runs_root="${QFORGE_RUN_ROOT:-$repo_root/runs}"

  if [[ $# -eq 0 ]]; then
    "$usage_fn" >&2
    return 1
  fi

  if [[ "${1:-}" == "--help" || "${1:-}" == "-h" ]]; then
    "$usage_fn"
    return 0
  fi

  local args=("$@")
  if [[ "${args[0]}" != "--run-dir" ]]; then
    args=(--run-dir "${args[0]}" "${args[@]:1}")
  fi

  local run_dir=""
  local i
  for ((i = 0; i < ${#args[@]}; i++)); do
    if [[ "${args[i]}" == "--run-dir" && $((i + 1)) -lt ${#args[@]} ]]; then
      run_dir="${args[i + 1]}"
      break
    fi
  done

  if [[ -z "$run_dir" ]]; then
    echo "${error_label} requires a run directory" >&2
    "$usage_fn" >&2
    return 1
  fi

  if [[ "$run_dir" != /* ]]; then
    run_dir="${runs_root}/${run_dir}"
  fi

  local args_with_path=()
  local skip_next=0
  for ((i = 0; i < ${#args[@]}; i++)); do
    if [[ $skip_next -eq 1 ]]; then
      skip_next=0
      continue
    fi
    if [[ "${args[i]}" == "--run-dir" && $((i + 1)) -lt ${#args[@]} ]]; then
      args_with_path+=(--run-dir "$run_dir")
      skip_next=1
      continue
    fi
    args_with_path+=("${args[i]}")
  done

  exec "${script_dir}/qforge" "$command" -v "${args_with_path[@]}"
}
