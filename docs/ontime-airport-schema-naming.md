# OnTime Airport Schema Naming Rationale

This note explains the airport-focused naming cleanup applied to the OnTime dataset.

## What Changed

The dataset now uses star-schema table names:

- `ontime.fact_ontime`
- `ontime.stage_ontime`
- `ontime.dim_airports`
- `ontime.dim_airports_bts_full`

The fact table also renamed the airport code and geography fields:

- `Origin` -> `OriginCode`
- `Dest` -> `DestCode`
- `OriginWac` -> `OriginWorldAreaCode`
- `DestWac` -> `DestWorldAreaCode`

The airport dimension moved from snake_case to explicit CamelCase names such as:

- `AirportCode`
- `AirportID`
- `AirportSeqID`
- `DisplayAirportName`
- `WorldAreaCode`
- `UtcLocalTimeVariation`

## Why These Changes

The old names were compact but ambiguous.

- `Origin` and `Dest` looked like airport entities, not airport code fields.
- `Wac` was too terse for most readers and especially weak for model-generated SQL.
- `airports_latest` and `airports_bts` described implementation history, not their semantic role.
- snake_case airport dimension columns did not line up with the CamelCase OnTime fact schema.

The new names make three things explicit:

1. table role: fact vs dimension
2. field semantics: code vs ID vs classification
3. preferred join path: stable airport IDs first, airport codes second

## Join Guidance

Use airport IDs for enrichment joins:

```sql
SELECT
    f.OriginCode,
    a.DisplayAirportName,
    a.Latitude,
    a.Longitude
FROM ontime.fact_ontime AS f
LEFT JOIN ontime.dim_airports AS a
    ON f.OriginAirportID = a.AirportID
LIMIT 20
```

Use airport codes for filtering, grouping, or fallback lookups:

```sql
SELECT
    AirportCode,
    DisplayAirportName
FROM ontime.dim_airports
WHERE AirportCode IN ('SEA', 'BWI', 'ISP')
ORDER BY AirportCode
```

Do not use `OriginWorldAreaCode`, `DestWorldAreaCode`, or `WorldAreaCode` as airport join keys. They are geography classification fields, not unique airport identifiers.

## Why `AirportID` Still Matters

`OriginCode` and `DestCode` are easier to read, but they are not the most stable keys across long time ranges.

- airport codes can change
- airport codes can be reused
- airport names and attributes can evolve over time

`AirportID` remains the canonical join key because it is the stable DOT airport identifier already carried by the fact table.

## Why This Helps Humans and Models

The rename improves generated SQL and hand-written SQL in the same ways:

- clearer `SELECT` lists
- fewer mistaken joins on WAC fields
- fewer mistaken assumptions that `Origin` or `Dest` are rich airport objects
- easier prompt instructions because the names describe their own intent

In short, the schema is now more descriptive without changing the underlying analytical grain.

## Old vs New Quick Reference

| Old | New |
| --- | --- |
| `ontime.ontime` | `ontime.fact_ontime` |
| `ontime.ontime_stage` | `ontime.stage_ontime` |
| `ontime.airports_latest` | `ontime.dim_airports` |
| `ontime.airports_bts` | `ontime.dim_airports_bts_full` |
| `Origin` | `OriginCode` |
| `Dest` | `DestCode` |
| `OriginWac` | `OriginWorldAreaCode` |
| `DestWac` | `DestWorldAreaCode` |
| `code` | `AirportCode` |
| `airport_id` | `AirportID` |
| `name` | `DisplayAirportName` |
| `wac` | `WorldAreaCode` |

## Migration Intent

This was an intentional breaking rename.

The goal was not backward compatibility. The goal was a schema that is easier to inspect, easier to teach in prompts, and safer for future analytical and blog-facing examples.
