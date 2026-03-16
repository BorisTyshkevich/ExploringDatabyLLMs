# OnTime V2

This module rebuilds the BTS Reporting Carrier On-Time dataset into a new ClickHouse table, `default.ontime_v2`, without touching the legacy `default.ontime` table.

## Source

Official monthly PREZIP archives from TranStats:

- `https://transtats.bts.gov/PREZIP/`
- file pattern:
  - `On_Time_Reporting_Carrier_On_Time_Performance_1987_present_{YYYY}_{M}.zip`

The loader validates that each archive contains one CSV plus optional `readme.html`.

## Tables

- `default.ontime_v2`
- `default.ontime_v2_stage`

Both use yearly partitions. The stage table is rebuilt for a target year and then published with:

```sql
ALTER TABLE default.ontime_v2
REPLACE PARTITION <year>
FROM default.ontime_v2_stage
```

## Loader Commands

Create tables:

```bash
python3 datasets/ontime_v2/load_ontime_v2.py create-tables --connection demo
```

List available months from TranStats:

```bash
python3 datasets/ontime_v2/load_ontime_v2.py list-available
```

Load one year:

```bash
python3 datasets/ontime_v2/load_ontime_v2.py load-year --year 2025 --connection demo
```

Backfill a range:

```bash
python3 datasets/ontime_v2/load_ontime_v2.py backfill --start-year 1987 --end-year 2025 --connection demo
```

Allow new unknown source columns temporarily:

```bash
python3 datasets/ontime_v2/load_ontime_v2.py load-year --year 2025 --allow-new-columns
```

## Type Rules

- `FlightDate` is the only stored calendar source column.
- `Year`, `Quarter`, `Month`, `DayofMonth`, and `DayOfWeek` are aliases derived from `FlightDate`.
- IDs use `0` as the missing sentinel.
- Metrics and HHMM operational time fields stay `Nullable`.
- Strings use empty string when source values are blank.

## Existing-Table Analysis

Generate a markdown report describing the current live `ontime` table and historical repo references:

```bash
python3 datasets/ontime_v2/analyze_existing_ontime.py --connection demo
```
