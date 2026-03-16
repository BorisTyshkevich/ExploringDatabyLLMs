# OnTime V2

This module rebuilds the BTS Reporting Carrier On-Time dataset into a new ClickHouse table, `default.ontime_v2`, without touching the legacy `default.ontime` table.

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

- `default.ontime_v2`
- `default.ontime_v2_stage`
- `default.airports_bts`
- `default.airports_bts_latest`

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

## BTS Airport Dimension

The official airport dimension is loaded separately from `default.airports`. It preserves the BTS `Master Coordinate` history in `default.airports_bts` and exposes `default.airports_bts_latest` as a latest-only helper view filtered to `AIRPORT_IS_LATEST = 1`.

Create the BTS table and latest view:

```bash
python3 datasets/ontime_v2/download/load_airports_bts.py create-tables --connection demo
```

Download and inspect the current official export:

```bash
python3 datasets/ontime_v2/download/load_airports_bts.py download --connection demo
```

Load the current official export into ClickHouse:

```bash
python3 datasets/ontime_v2/download/load_airports_bts.py load --connection demo
```

Verify table counts and OnTime join coverage:

```bash
python3 datasets/ontime_v2/download/load_airports_bts.py verify --connection demo
```

Example join using the latest BTS airport attributes:

```sql
SELECT
    replaceAll(toString(o.Origin), '\0', '') AS Origin,
    any(a.DISPLAY_AIRPORT_NAME) AS AirportName,
    any(a.LATITUDE) AS Latitude,
    any(a.LONGITUDE) AS Longitude,
    any(a.UTC_LOCAL_TIME_VARIATION) AS UtcLocalTimeVariation
FROM default.ontime_v2 AS o
LEFT JOIN default.airports_bts_latest AS a
    ON o.OriginAirportID = a.AIRPORT_ID
GROUP BY Origin
ORDER BY Origin
LIMIT 20
```
