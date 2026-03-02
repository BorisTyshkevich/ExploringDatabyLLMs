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

mkdir -p results/clickhouse

pass_count=0
fail_count=0

run_sql_file() {
  local label="$1"
  local host="$2"
  local port="$3"
  local user="$4"
  local password="$5"
  local database="$6"
  local sql_file="$7"
  local out_file="$8"

  if clickhouse-client \
    --secure \
    --host "$host" \
    --port "$port" \
    --user "$user" \
    --password "$password" \
    --database "$database" \
    --queries-file "$sql_file" \
    --format JSON >"$out_file" 2>"${out_file}.err"; then
    echo "PASS [$label] $sql_file"
    pass_count=$((pass_count + 1))
  else
    echo "FAIL [$label] $sql_file"
    cat "${out_file}.err"
    fail_count=$((fail_count + 1))
  fi
}

run_expect_fail() {
  local label="$1"
  local host="$2"
  local port="$3"
  local user="$4"
  local password="$5"
  local database="$6"
  local sql_file="$7"
  local out_file="$8"

  if clickhouse-client \
    --secure \
    --host "$host" \
    --port "$port" \
    --user "$user" \
    --password "$password" \
    --database "$database" \
    --queries-file "$sql_file" \
    --format JSON >"$out_file" 2>"${out_file}.err"; then
    echo "UNEXPECTED_PASS [$label] $sql_file"
    fail_count=$((fail_count + 1))
  else
    echo "EXPECTED_FAIL [$label] $sql_file"
    pass_count=$((pass_count + 1))
  fi
}

echo "=== Running Altinity queries ==="
for file in sql/altinity/*.sql; do
  base="$(basename "$file" .sql)"
  run_sql_file "altinity" "$ALTINITY_HOST" "$ALTINITY_PORT" "$ALTINITY_USER" "$ALTINITY_PASSWORD" "$ALTINITY_DATABASE" "$file" "results/clickhouse/altinity_${base}.json"
done

echo "=== Running ClickHouse Inc queries ==="
for file in sql/clickhouse_inc/*.sql; do
  base="$(basename "$file" .sql)"
  run_sql_file "clickhouse_inc" "$CLICKHOUSE_INC_HOST" "$CLICKHOUSE_INC_PORT" "$CLICKHOUSE_INC_USER" "$CLICKHOUSE_INC_PASSWORD" "$CLICKHOUSE_INC_DATABASE" "$file" "results/clickhouse/clickhouse_inc_${base}.json"
done

echo "=== Verifying Part 2 SQL as written (expected failures) ==="
run_expect_fail "part2_as_written_altinity_case1" "$ALTINITY_HOST" "$ALTINITY_PORT" "$ALTINITY_USER" "$ALTINITY_PASSWORD" "$ALTINITY_DATABASE" "sql/part2/01_case1_as_written.sql" "results/clickhouse/part2_as_written_altinity_case1.json"
run_expect_fail "part2_as_written_altinity_case2" "$ALTINITY_HOST" "$ALTINITY_PORT" "$ALTINITY_USER" "$ALTINITY_PASSWORD" "$ALTINITY_DATABASE" "sql/part2/02_case2_as_written.sql" "results/clickhouse/part2_as_written_altinity_case2.json"

echo "=== Verifying corrected Part 2 SQL ==="
run_sql_file "part2_fixed_altinity_case1" "$ALTINITY_HOST" "$ALTINITY_PORT" "$ALTINITY_USER" "$ALTINITY_PASSWORD" "$ALTINITY_DATABASE" "sql/part2/03_case1_fixed_altinity.sql" "results/clickhouse/part2_fixed_altinity_case1.json"
run_sql_file "part2_fixed_altinity_case2" "$ALTINITY_HOST" "$ALTINITY_PORT" "$ALTINITY_USER" "$ALTINITY_PASSWORD" "$ALTINITY_DATABASE" "sql/part2/04_case2_fixed_altinity.sql" "results/clickhouse/part2_fixed_altinity_case2.json"
run_sql_file "part2_fixed_clickhouse_inc_case1" "$CLICKHOUSE_INC_HOST" "$CLICKHOUSE_INC_PORT" "$CLICKHOUSE_INC_USER" "$CLICKHOUSE_INC_PASSWORD" "$CLICKHOUSE_INC_DATABASE" "sql/part2/05_case1_fixed_clickhouse_inc.sql" "results/clickhouse/part2_fixed_clickhouse_inc_case1.json"
run_sql_file "part2_fixed_clickhouse_inc_case2" "$CLICKHOUSE_INC_HOST" "$CLICKHOUSE_INC_PORT" "$CLICKHOUSE_INC_USER" "$CLICKHOUSE_INC_PASSWORD" "$CLICKHOUSE_INC_DATABASE" "sql/part2/06_case2_fixed_clickhouse_inc.sql" "results/clickhouse/part2_fixed_clickhouse_inc_case2.json"

echo "=== Summary ==="
echo "PASS: $pass_count"
echo "FAIL: $fail_count"

if [[ "$fail_count" -gt 0 ]]; then
  exit 1
fi
