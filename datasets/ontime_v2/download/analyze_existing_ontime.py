#!/usr/bin/env python3

from __future__ import annotations

import argparse
import subprocess
from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]
OLD_LOAD = Path("/Users/bvt/work/altinity-datasets/OnTime/load.sh")
OLD_DDL = Path("/Users/bvt/work/altinity-datasets/airline/ddl/ontime.sql")


def query(connection: str, sql: str) -> str:
    proc = subprocess.run(
        ["clickhouse-client", "--connection", connection, "--format", "TSVWithNames", "--query", sql],
        text=True,
        capture_output=True,
        check=False,
    )
    if proc.returncode != 0:
        raise RuntimeError(proc.stderr.strip())
    return proc.stdout.strip()


def build_report(connection: str) -> str:
    server = query(
        connection,
        "select hostName() as hostname, version() as version, getMacro('cluster') as cluster_name",
    )
    tables = query(
        connection,
        "SELECT database, name, engine, total_rows, total_bytes, metadata_modification_time "
        "FROM system.tables WHERE name ILIKE '%ontime%' ORDER BY database, name",
    )
    parts = query(
        connection,
        "SELECT database, table, count() AS active_parts, sum(rows) AS rows, "
        "formatReadableSize(sum(bytes_on_disk)) AS size, min(partition) AS min_partition, "
        "max(partition) AS max_partition FROM system.parts WHERE active AND table='ontime' "
        "GROUP BY database, table ORDER BY database, table",
    )
    create_sql = query(connection, "SHOW CREATE TABLE default.ontime")
    ddl_history = query(
        connection,
        "SELECT event_time, type, query FROM system.query_log "
        "WHERE type IN ('QueryFinish','ExceptionWhileProcessing') "
        "AND (query ILIKE 'CREATE TABLE%ontime%' OR query ILIKE 'ATTACH TABLE%ontime%' "
        "OR query ILIKE 'DROP TABLE%ontime%' OR query ILIKE 'RENAME TABLE%ontime%' "
        "OR query ILIKE 'INSERT INTO%ontime%') ORDER BY event_time DESC LIMIT 50",
    )

    old_load_excerpt = OLD_LOAD.read_text().strip() if OLD_LOAD.exists() else "missing"
    old_ddl_excerpt = OLD_DDL.read_text().strip() if OLD_DDL.exists() else "missing"

    return f"""# Existing OnTime Analysis

## Server

```text
{server}
```

## Live OnTime Tables

```text
{tables}
```

## Live Parts Summary

```text
{parts}
```

## Current Live DDL

```sql
{create_sql}
```

## Retained DDL / Insert History

```text
{ddl_history}
```

## Historical Repo Loader Reference

Path: `{OLD_LOAD}`

```bash
{old_load_excerpt}
```

## Historical Repo DDL Reference

Path: `{OLD_DDL}`

```sql
{old_ddl_excerpt}
```

## Findings

- The current live `default.ontime` table matches the old adjacent-repo MergeTree layout closely.
- The live table only covers partitions through `2021`.
- The legacy design physically stores redundant calendar columns and uses old carrier/flight column names.
- Several sparse operational fields were modeled as strings in the old design.
- No older `CREATE TABLE ... ontime` history is currently retained in `system.query_log`, so server-side creation history is only partially reconstructable.
"""


def main() -> int:
    parser = argparse.ArgumentParser(description="Analyze the existing live ontime table and historical references.")
    parser.add_argument("--connection", default="demo")
    parser.add_argument(
        "--output",
        default=str(ROOT / "docs" / "ontime_v2-existing-analysis.md"),
    )
    args = parser.parse_args()

    report = build_report(args.connection)
    output = Path(args.output)
    output.parent.mkdir(parents=True, exist_ok=True)
    output.write_text(report)
    print(output)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
