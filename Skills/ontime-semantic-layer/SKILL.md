---
name: ontime-semantic-layer
description: OnTime dataset schema inspection, join guidance, and airport-dimension semantics for SQL and visual prompt work.
---

Use this skill whenever work targets the OnTime dataset and needs table discovery, join guidance, or airport-field interpretation.

## First step

- Inspect the live schema with `SHOW TABLES FROM ontime`
- Run `DESCRIBE TABLE` for the tables you actually plan to use
- Use `SHOW CREATE TABLE` only when column defaults, aliases, engine details, or other DDL-specific behavior matters

## Canonical tables

- `ontime.fact_ontime`
- `ontime.dim_airports`  - actual data
- `ontime.dim_airports_bts_full` - BTS historic data

## Join guidance

- Prefer `ontime.fact_ontime.OriginAirportID = ontime.dim_airports.AirportID`
- Prefer `ontime.fact_ontime.DestAirportID = ontime.dim_airports.AirportID`
- Use `OriginCode` and `DestCode` for readable filters, grouping, and display
- Use `ontime.dim_airports.AirportCode` for airport-code lookups or fallback joins
- Do not use `OriginWorldAreaCode`, `DestWorldAreaCode`, or `WorldAreaCode` as join keys

## Airport dimension fields

Use `ontime.dim_airports` for current airport reference data, especially:

- `AirportCode`
- `DisplayAirportName`
- `Latitude`
- `Longitude`
- `UtcLocalTimeVariation`


