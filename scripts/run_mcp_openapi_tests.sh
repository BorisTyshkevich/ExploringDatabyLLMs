#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

if [[ -f .env ]]; then
  # shellcheck disable=SC1091
  source .env
else
  # shellcheck disable=SC1091
  source .env.example
fi

mkdir -p results/mcp

execute_url="${MCP_BASE_URL}/${MCP_DEMO_TOKEN}/openapi/execute_query"

normalize_response() {
  local input="$1"
  local output="$2"
  if [[ "$input" =~ ^\{ ]]; then
    printf '%s\n' "$input" | jq . >"$output"
  else
    jq -n --arg error "$input" '{error: $error}' >"$output"
  fi
}

health_raw="$(curl -sS "${MCP_BASE_URL}/health")"
normalize_response "$health_raw" "results/mcp/health.json"

version_raw="$(curl -sS --get "$execute_url" \
  --data-urlencode "query=select version() as version, hostName() as host")"
normalize_response "$version_raw" "results/mcp/version_host.json"

as_written_raw="$(curl -sS --get "$execute_url" \
  --data-urlencode "query=SELECT carrier, COUNT(*) AS flights FROM ontime WHERE year = 2019 GROUP BY carrier ORDER BY flights DESC LIMIT 1")"
normalize_response "$as_written_raw" "results/mcp/part2_case1_as_written.json"

fixed_raw="$(curl -sS --get "$execute_url" \
  --data-urlencode "query=SELECT trimRight(toString(Carrier)) AS carrier, COUNT(*) AS flights FROM ontime WHERE Year = 2019 GROUP BY carrier ORDER BY flights DESC LIMIT 1")"
normalize_response "$fixed_raw" "results/mcp/part2_case1_fixed.json"

echo "=== MCP checks ==="
echo "Health:"
jq -r '.status, .version' results/mcp/health.json

echo "As-written query error:"
jq -r '.error // "<none>"' results/mcp/part2_case1_as_written.json

echo "Fixed query first row:"
jq -c '.rows[0]' results/mcp/part2_case1_fixed.json
