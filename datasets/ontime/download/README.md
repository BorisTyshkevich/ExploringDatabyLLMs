# OnTime V2

This module rebuilds the BTS Reporting Carrier On-Time dataset into a new ClickHouse fact table, `ontime.fact_ontime`, without touching the legacy `default.ontime` table.

## Source

Official monthly PREZIP archives from TranStats:

- `https://transtats.bts.gov/PREZIP/`
- file pattern:
  - `On_Time_Reporting_Carrier_On_Time_Performance_1987_present_{YYYY}_{M}.zip`

The loader validates that each archive contains one CSV plus optional `readme.html`.

Official BTS airport dimension source:

- live TranStats export for `Master Coordinate`
- table info: `https://www.transtats.bts.gov/Fields.asp?gnoyr_VQ=FLL`
- export form: `https://www.transtats.bts.gov/DL_SelectFields.aspx?gnoyr_VQ=FLL&QO_fu146_anzr=N8vn6v10+f722146+gnoyr5`
- archived documentation reference: `https://rosap.ntl.bts.gov/view/dot/58890`

## Tables

- `ontime.fact_ontime`
- `ontime.stage_ontime`
- `ontime.dim_airports_bts_full`
- `ontime.dim_airports`

Both use yearly partitions. The stage table is rebuilt for a target year and then published with:

```sql
ALTER TABLE ontime.fact_ontime
REPLACE PARTITION <year>
FROM ontime.stage_ontime
```

Before publish, the loader runs:

```sql
OPTIMIZE TABLE ontime.stage_ontime PARTITION <year> FINAL DEDUPLICATE
```

This removes exact duplicate rows within the staged year before `REPLACE PARTITION`. It does not collapse rows that share the same analytical leg key but differ in other columns.

## Loader Commands

Create tables:

```bash
python3 datasets/ontime/download/load_ontime_v2.py create-tables --connection demo
```

List available months from TranStats:

```bash
python3 datasets/ontime/download/load_ontime_v2.py list-available
```

Load one year:

```bash
python3 datasets/ontime/download/load_ontime_v2.py load-year --year 2025 --connection demo
```

Backfill a range:

```bash
python3 datasets/ontime/download/load_ontime_v2.py backfill --start-year 1987 --end-year 2025 --connection demo
```

Allow new unknown source columns temporarily:

```bash
python3 datasets/ontime/download/load_ontime_v2.py load-year --year 2025 --allow-new-columns
```

## Type Rules

- `FlightDate` is the only stored calendar source column.
- `Year`, `Quarter`, `Month`, `DayofMonth`, and `DayOfWeek` are aliases derived from `FlightDate`.
- IDs use `0` as the missing sentinel.
- Metrics and HHMM operational time fields stay `Nullable`.
- Strings use empty string when source values are blank.

## Deduplication

- Raw normalized rows are inserted into `ontime.stage_ontime`.
- After the full target year is loaded, the loader runs `OPTIMIZE TABLE ... PARTITION <year> FINAL DEDUPLICATE` on the stage table.
- This step removes only exact duplicate rows for the target year.
- Month metadata in `.cache/meta/*.json` records the staged year row count before deduplication, after deduplication, and the number of exact duplicate rows removed.

## Existing-Table Analysis

Generate a markdown report describing the current live `ontime` table and historical repo references:

```bash
python3 datasets/ontime/download/analyze_existing_ontime.py --connection demo
```

## BTS Airport Dimension

The official airport dimension is loaded separately from `default.airports`. It preserves BTS `Master Coordinate` history in `ontime.dim_airports_bts_full` and exposes `ontime.dim_airports` as the cleaned semantic airport view with a single latest row per airport code.

Create the airport table and latest view:

```bash
python3 datasets/ontime/download/load_airports_bts.py create-tables --connection demo
```

Download and inspect the current official export:

```bash
python3 datasets/ontime/download/load_airports_bts.py download --connection demo
```

Load the current official export into ClickHouse:

```bash
python3 datasets/ontime/download/load_airports_bts.py load --connection demo
```

Verify table counts, uniqueness, and OnTime join coverage:

```bash
python3 datasets/ontime/download/load_airports_bts.py verify --connection demo
```

Example join using the latest airport attributes:

```sql
SELECT
    replaceAll(toString(o.OriginCode), '\0', '') AS OriginCode,
    any(a.DisplayAirportName) AS AirportName,
    any(a.Latitude) AS Latitude,
    any(a.Longitude) AS Longitude,
    any(a.UtcLocalTimeVariation) AS UtcLocalTimeVariation
FROM ontime.fact_ontime AS o
LEFT JOIN ontime.dim_airports AS a
    ON o.OriginAirportID = a.AirportID
GROUP BY OriginCode
ORDER BY OriginCode
LIMIT 20
```

Example code-based lookup using the cleaned latest airport view:

```sql
SELECT
    AirportCode,
    DisplayAirportName,
    Latitude,
    Longitude
FROM ontime.dim_airports
WHERE AirportCode IN ('ISP', 'BWI', 'SEA')
ORDER BY AirportCode
```
