#!/usr/bin/env python3
import json
import os
import re
import shlex
import subprocess
import sys
import time
from dataclasses import dataclass
from pathlib import Path
from typing import Any

ROOT = Path(__file__).resolve().parents[1]
ENV_FILE = ROOT / ".env"
ENV_EXAMPLE_FILE = ROOT / ".env.example"
RESULT_DIR = ROOT / "results" / "benchmark"
SQL_DIR = RESULT_DIR / "sql"
RAW_DIR = RESULT_DIR / "raw"
EXEC_DIR = RESULT_DIR / "exec"


@dataclass
class Model:
    key: str
    vendor: str
    runner: str
    model: str


MODELS = [
    Model("openai_gpt5", "openai", "codex", "gpt-5"),
    Model("openai_gpt53codex", "openai", "codex", "gpt-5.3-codex"),
    Model("anthropic_sonnet", "anthropic", "claude", "sonnet"),
    Model("anthropic_opus", "anthropic", "claude", "opus"),
]

CASES = [
    {
        "id": "case1_top_carrier_2019",
        "question": "Which carrier operated the most flights in 2019? Return carrier and flights.",
        "expected_sql": """
SELECT trimRight(toString(Carrier)) AS carrier, count() AS flights
FROM ontime
WHERE Year = 2019
GROUP BY carrier
ORDER BY flights DESC
LIMIT 1
""".strip(),
    },
    {
        "id": "case2_avg_depdelay_delta_atl_2021_07",
        "question": "What was the average departure delay for Delta flights out of Atlanta in July 2021? Return one rounded value avg_dep_delay.",
        "expected_sql": """
SELECT round(avg(DepDelay), 4) AS avg_dep_delay
FROM ontime
WHERE trimRight(toString(Carrier)) = 'DL'
  AND trimRight(toString(Origin)) = 'ATL'
  AND Year = 2021
  AND Month = 7
""".strip(),
    },
    {
        "id": "case3_worst_airport_2021",
        "question": "In 2021, which origin airport had the highest delayed arrival percentage? Use delayed arrivals as ArrDel15=1, require at least 10000 flights. Return airport, delayed_arrival_pct, flights.",
        "expected_sql": """
SELECT trimRight(toString(Origin)) AS airport,
       round(100.0 * avg(ArrDel15), 2) AS delayed_arrival_pct,
       count() AS flights
FROM ontime
WHERE Year = 2021
GROUP BY airport
HAVING flights >= 10000
ORDER BY delayed_arrival_pct DESC
LIMIT 1
""".strip(),
    },
    {
        "id": "case4_worst_winter_carrier_airport",
        "question": "Across winters (December, January, February) from 2019 to 2021, which carrier-airport pair had the highest average departure delay among groups with at least 5000 flights? Return carrier, airport, avg_dep_delay, flights.",
        "expected_sql": """
SELECT trimRight(toString(Carrier)) AS carrier,
       trimRight(toString(Origin)) AS airport,
       round(avg(DepDelay), 2) AS avg_dep_delay,
       count() AS flights
FROM ontime
WHERE Month IN (12, 1, 2)
  AND Year BETWEEN 2019 AND 2021
GROUP BY carrier, airport
HAVING flights >= 5000
ORDER BY avg_dep_delay DESC
LIMIT 1
""".strip(),
    },
    {
        "id": "case5_peak_aa_month",
        "question": "Between 2018 and 2021, which year-month had the highest average departure delay for AA? Return Year, Month, avg_dep_delay.",
        "expected_sql": """
SELECT Year,
       Month,
       round(avg(DepDelay), 2) AS avg_dep_delay
FROM ontime
WHERE trimRight(toString(Carrier)) = 'AA'
  AND Year BETWEEN 2018 AND 2021
GROUP BY Year, Month
ORDER BY avg_dep_delay DESC, Year DESC, Month DESC
LIMIT 1
""".strip(),
    },
    {
        "id": "case6_peak_route_season",
        "question": "For years 2019-2021, take top 50 routes by flight count (Origin-Dest), then find the route-season with highest average departure delay. Seasons: Winter(12,1,2), Spring(3,4,5), Summer(6,7,8), Fall(9,10,11). Return route, season, avg_dep_delay, flights.",
        "expected_sql": """
WITH routes AS (
    SELECT trimRight(toString(Origin)) AS origin,
           trimRight(toString(Dest)) AS dest,
           count() AS flights
    FROM ontime
    WHERE Year BETWEEN 2019 AND 2021
    GROUP BY origin, dest
    ORDER BY flights DESC
    LIMIT 50
)
SELECT concat(r.origin, '-', r.dest) AS route,
       multiIf(o.Month IN (12,1,2), 'Winter', o.Month IN (3,4,5), 'Spring', o.Month IN (6,7,8), 'Summer', 'Fall') AS season,
       round(avg(o.DepDelay), 2) AS avg_dep_delay,
       count() AS flights
FROM ontime o
INNER JOIN routes r
    ON trimRight(toString(o.Origin)) = r.origin
   AND trimRight(toString(o.Dest)) = r.dest
WHERE o.Year BETWEEN 2019 AND 2021
GROUP BY route, season
HAVING flights >= 1000
ORDER BY avg_dep_delay DESC
LIMIT 1
""".strip(),
    },
]


def load_env() -> dict[str, str]:
    env_path = ENV_FILE if ENV_FILE.exists() else ENV_EXAMPLE_FILE
    env: dict[str, str] = {}
    for line in env_path.read_text().splitlines():
        line = line.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        k, v = line.split("=", 1)
        env[k.strip()] = v.strip()
    return env


def run(cmd: list[str], timeout: int = 300) -> subprocess.CompletedProcess[str]:
    return subprocess.run(cmd, capture_output=True, text=True, timeout=timeout)


def sanitize_sql(text: str) -> str:
    s = text.strip()
    s = re.sub(r"^```(?:sql)?", "", s, flags=re.IGNORECASE).strip()
    s = re.sub(r"```$", "", s, flags=re.IGNORECASE).strip()
    if "```" in s:
        blocks = re.findall(r"```(?:sql)?\s*(.*?)```", s, flags=re.IGNORECASE | re.DOTALL)
        if blocks:
            s = blocks[0].strip()
    s = s.strip().rstrip(";")
    return s


def build_prompt(question: str) -> str:
    return (
        "Generate one ClickHouse SQL query only (no markdown, no comments, no explanation).\n"
        "Database: default\n"
        "Table: ontime\n"
        "Important schema notes:\n"
        "- Column names are case-sensitive (Year, Month, Carrier, Origin, Dest, DepDelay, ArrDel15).\n"
        "- Carrier and Origin/Dest are fixed-width strings; safe comparisons use trimRight(toString(...)).\n"
        "- Return exactly the fields requested by the question.\n"
        "- Always include ORDER BY when using LIMIT.\n"
        "Question: "
        + question
    )


def generate_sql(model: Model, prompt: str) -> tuple[str, str, int]:
    if model.runner == "codex":
        out_file = RAW_DIR / f"{model.key}_last_message.txt"
        cmd = [
            "codex",
            "exec",
            "--skip-git-repo-check",
            "-m",
            model.model,
            "-c",
            'web_search="disabled"',
            "-c",
            'model_reasoning_effort="medium"',
            "-o",
            str(out_file),
            prompt,
        ]
        cp = run(cmd, timeout=600)
        raw = out_file.read_text() if out_file.exists() else ""
        return raw, cp.stderr + cp.stdout, cp.returncode

    if model.runner == "claude":
        cmd = [
            "claude",
            "-p",
            "--output-format",
            "json",
            "--model",
            model.model,
            prompt,
        ]
        cp = run(cmd, timeout=600)
        raw = ""
        if cp.returncode == 0:
            try:
                parsed = json.loads(cp.stdout)
                raw = parsed.get("result", "")
            except json.JSONDecodeError:
                raw = cp.stdout
        return raw, cp.stderr + cp.stdout, cp.returncode

    raise ValueError(f"Unknown runner: {model.runner}")


def exec_clickhouse_sql(env: dict[str, str], sql: str) -> tuple[bool, Any, str, float]:
    cmd = [
        "clickhouse-client",
        "--secure",
        "--host",
        env["ALTINITY_HOST"],
        "--port",
        env.get("ALTINITY_PORT", "9440"),
        "--user",
        env["ALTINITY_USER"],
        "--password",
        env.get("ALTINITY_PASSWORD", ""),
        "--database",
        env.get("ALTINITY_DATABASE", "default"),
        "--query",
        sql,
        "--format",
        "JSON",
    ]
    t0 = time.perf_counter()
    cp = run(cmd, timeout=300)
    elapsed = time.perf_counter() - t0
    if cp.returncode != 0:
        return False, None, cp.stderr.strip() or cp.stdout.strip(), elapsed
    try:
        parsed = json.loads(cp.stdout)
        return True, parsed.get("data", []), "", elapsed
    except json.JSONDecodeError:
        return False, None, "Invalid JSON result", elapsed


def json_norm(v: Any) -> str:
    return json.dumps(v, sort_keys=True, ensure_ascii=True)


def main() -> int:
    env = load_env()
    RESULT_DIR.mkdir(parents=True, exist_ok=True)
    SQL_DIR.mkdir(parents=True, exist_ok=True)
    RAW_DIR.mkdir(parents=True, exist_ok=True)
    EXEC_DIR.mkdir(parents=True, exist_ok=True)

    expected_results: dict[str, Any] = {}
    print("Computing expected results from reference SQL...")
    for case in CASES:
        ok, data, err, elapsed = exec_clickhouse_sql(env, case["expected_sql"])
        if not ok:
            print(f"FAILED expected SQL for {case['id']}: {err}")
            return 1
        expected_results[case["id"]] = data
        (EXEC_DIR / f"expected_{case['id']}.json").write_text(json.dumps(data, indent=2))
        print(f"  {case['id']}: expected rows={len(data)} elapsed={elapsed:.3f}s")

    rows: list[dict[str, Any]] = []
    for model in MODELS:
        print(f"\\nModel: {model.key} ({model.model})")
        model_sql_dir = SQL_DIR / model.key
        model_sql_dir.mkdir(parents=True, exist_ok=True)

        for case in CASES:
            prompt = build_prompt(case["question"])
            raw, generation_log, gen_rc = generate_sql(model, prompt)
            sql = sanitize_sql(raw)

            raw_path = RAW_DIR / f"{model.key}_{case['id']}.txt"
            raw_path.write_text(raw)
            log_path = RAW_DIR / f"{model.key}_{case['id']}.log"
            log_path.write_text(generation_log)

            (model_sql_dir / f"{case['id']}.sql").write_text(sql + "\n")

            if gen_rc != 0 or not sql:
                row = {
                    "model_key": model.key,
                    "vendor": model.vendor,
                    "model": model.model,
                    "case_id": case["id"],
                    "generation_ok": False,
                    "execution_ok": False,
                    "correct": False,
                    "error": f"generation_failed rc={gen_rc}",
                    "elapsed_sec": None,
                }
                rows.append(row)
                print(f"  {case['id']}: generation FAILED")
                continue

            ok, data, err, elapsed = exec_clickhouse_sql(env, sql)
            result_path = EXEC_DIR / f"{model.key}_{case['id']}.json"
            result_path.write_text(json.dumps({"ok": ok, "data": data, "error": err, "elapsed_sec": elapsed}, indent=2))

            correct = ok and json_norm(data) == json_norm(expected_results[case["id"]])
            row = {
                "model_key": model.key,
                "vendor": model.vendor,
                "model": model.model,
                "case_id": case["id"],
                "generation_ok": True,
                "execution_ok": ok,
                "correct": correct,
                "error": err if not ok else "",
                "elapsed_sec": round(elapsed, 4),
            }
            rows.append(row)
            status = "OK" if correct else ("RUNS_WRONG" if ok else "EXEC_FAIL")
            print(f"  {case['id']}: {status} ({elapsed:.3f}s)")

    summary: dict[str, Any] = {
        "generated_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
        "dataset": "default.ontime on demo.demo.altinity.cloud",
        "models": [model.__dict__ for model in MODELS],
        "cases": [{"id": c["id"], "question": c["question"]} for c in CASES],
        "rows": rows,
    }

    # per-model aggregates
    aggregates = {}
    for model in MODELS:
        mrows = [r for r in rows if r["model_key"] == model.key]
        total = len(mrows)
        exec_ok = sum(1 for r in mrows if r["execution_ok"])
        correct = sum(1 for r in mrows if r["correct"])
        avg_elapsed = round(sum(r["elapsed_sec"] or 0 for r in mrows if r["elapsed_sec"] is not None) / max(exec_ok, 1), 4)
        aggregates[model.key] = {
            "total_cases": total,
            "execution_ok": exec_ok,
            "correct": correct,
            "execution_rate": round(exec_ok / total, 4) if total else 0,
            "correct_rate": round(correct / total, 4) if total else 0,
            "avg_elapsed_sec": avg_elapsed,
        }

    summary["aggregates"] = aggregates
    (RESULT_DIR / "benchmark_summary.json").write_text(json.dumps(summary, indent=2))

    # markdown report
    lines = []
    lines.append("# LLM SQL Benchmark Results")
    lines.append("")
    lines.append(f"Generated at: {summary['generated_at']}")
    lines.append("")
    lines.append("## Aggregate Scores")
    lines.append("")
    lines.append("| Model | Vendor | Execution OK | Correct | Total | Correct Rate | Avg Runtime (s) |")
    lines.append("|---|---|---:|---:|---:|---:|---:|")
    for model in MODELS:
        agg = aggregates[model.key]
        lines.append(
            f"| {model.key} | {model.vendor} | {agg['execution_ok']} | {agg['correct']} | {agg['total_cases']} | {agg['correct_rate']:.2%} | {agg['avg_elapsed_sec']} |"
        )

    lines.append("")
    lines.append("## Per-Case Outcomes")
    lines.append("")
    lines.append("| Model | Case | Status | Error |")
    lines.append("|---|---|---|---|")
    for row in rows:
        if row["correct"]:
            status = "correct"
        elif row["execution_ok"]:
            status = "runs_wrong"
        elif row["generation_ok"]:
            status = "exec_fail"
        else:
            status = "generation_fail"
        error = row["error"].replace("|", "\\|")[:160]
        lines.append(f"| {row['model_key']} | {row['case_id']} | {status} | {error} |")

    (RESULT_DIR / "benchmark_report.md").write_text("\n".join(lines) + "\n")

    print("\\nBenchmark complete.")
    print(f"Summary: {RESULT_DIR / 'benchmark_summary.json'}")
    print(f"Report:  {RESULT_DIR / 'benchmark_report.md'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
