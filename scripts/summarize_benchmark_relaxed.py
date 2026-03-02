#!/usr/bin/env python3
import json
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
EXEC_DIR = ROOT / "results" / "benchmark" / "exec"
OUT_JSON = ROOT / "results" / "benchmark" / "benchmark_relaxed_summary.json"
OUT_MD = ROOT / "results" / "benchmark" / "benchmark_relaxed_report.md"

MODELS = ["openai_gpt5", "openai_gpt53codex", "anthropic_sonnet", "anthropic_opus"]
CASES = [
    "case1_top_carrier_2019",
    "case2_avg_depdelay_delta_atl_2021_07",
    "case3_worst_airport_2021",
    "case4_worst_winter_carrier_airport",
    "case5_peak_aa_month",
    "case6_peak_route_season",
]


def clean_str(v):
    return v.replace("\x00", "").strip() if isinstance(v, str) else v


def relaxed_match(case_id, expected, got):
    if case_id == "case6_peak_route_season" and "route" not in got and "Origin" in got and "Dest" in got:
        got = {
            "route": f"{clean_str(got['Origin'])}-{clean_str(got['Dest'])}",
            "season": got.get("season"),
            "avg_dep_delay": got.get("avg_dep_delay"),
            "flights": got.get("flights"),
        }

    for key in ["carrier", "airport", "route", "season"]:
        if key in got:
            got[key] = clean_str(got[key])

    for key in ["carrier", "airport", "route", "season", "Year", "Month", "flights"]:
        if key in expected and got.get(key) != expected.get(key):
            return False

    for key in ["avg_dep_delay", "delayed_arrival_pct"]:
        if key in expected:
            if key not in got:
                return False
            if abs(float(got[key]) - float(expected[key])) > 0.1:
                return False

    return True


def main():
    expected = {}
    for case_id in CASES:
        expected_row = json.loads((EXEC_DIR / f"expected_{case_id}.json").read_text())[0]
        expected[case_id] = expected_row

    result = {"models": {}}

    for model in MODELS:
        strict = 0
        relaxed = 0
        details = []
        for case_id in CASES:
            got_doc = json.loads((EXEC_DIR / f"{model}_{case_id}.json").read_text())
            got_row = got_doc["data"][0]
            strict_ok = got_row == expected[case_id]
            relaxed_ok = relaxed_match(case_id, expected[case_id], dict(got_row))
            strict += 1 if strict_ok else 0
            relaxed += 1 if relaxed_ok else 0
            details.append({
                "case_id": case_id,
                "strict": strict_ok,
                "relaxed": relaxed_ok,
                "expected": expected[case_id],
                "got": got_row,
            })

        result["models"][model] = {
            "strict_correct": strict,
            "relaxed_correct": relaxed,
            "total": len(CASES),
            "strict_rate": round(strict / len(CASES), 4),
            "relaxed_rate": round(relaxed / len(CASES), 4),
            "details": details,
        }

    OUT_JSON.write_text(json.dumps(result, indent=2))

    lines = [
        "# Relaxed vs Strict Benchmark Summary",
        "",
        "| Model | Strict Correct | Relaxed Correct | Total | Strict Rate | Relaxed Rate |",
        "|---|---:|---:|---:|---:|---:|",
    ]

    for model in MODELS:
        m = result["models"][model]
        lines.append(
            f"| {model} | {m['strict_correct']} | {m['relaxed_correct']} | {m['total']} | {m['strict_rate']:.2%} | {m['relaxed_rate']:.2%} |"
        )

    OUT_MD.write_text("\n".join(lines) + "\n")


if __name__ == "__main__":
    main()
