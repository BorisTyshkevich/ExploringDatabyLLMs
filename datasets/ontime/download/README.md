# OnTime V2

This module rebuilds the BTS Reporting Carrier On-Time dataset into a new ClickHouse table, `ontime.ontime`, without touching the legacy `default.ontime` table.

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

- `ontime.ontime`
- `ontime.ontime_stage`
- `ontime.airports_bts`
- `ontime.airports_latest`

Both use yearly partitions. The stage table is rebuilt for a target year and then published with:

```sql
ALTER TABLE ontime.ontime
REPLACE PARTITION <year>
FROM ontime.ontime_stage
```

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

## Existing-Table Analysis

Generate a markdown report describing the current live `ontime` table and historical repo references:

```bash
python3 datasets/ontime/download/analyze_existing_ontime.py --connection demo
```

## BTS Airport Dimension

The official airport dimension is loaded separately from `default.airports`. It preserves BTS `Master Coordinate` history in `ontime.airports_bts` and exposes `ontime.airports_latest` as the cleaned semantic airport view with a single latest row per airport code.

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
    replaceAll(toString(o.Origin), '\0', '') AS Origin,
    any(a.name) AS AirportName,
    any(a.latitude) AS Latitude,
    any(a.longitude) AS Longitude,
    any(a.utc_local_time_variation) AS UtcLocalTimeVariation
FROM ontime.ontime AS o
LEFT JOIN ontime.airports_latest AS a
    ON o.OriginAirportID = a.airport_id
GROUP BY Origin
ORDER BY Origin
LIMIT 20
```

Example code-based lookup using the cleaned latest airport view:

```sql
SELECT
    code,
    name,
    latitude,
    longitude
FROM ontime.airports_latest
WHERE code IN ('ISP', 'BWI', 'SEA')
ORDER BY code
```
