This file is repository documentation for the OnTime semantic model. Prompt generation does not inline this file.
For runtime prompt guidance, schema inspection, and join behavior, use the `ontime-semantic-layer` skill.

Use `ontime.fact_ontime` as the primary fact table for flight operations.

Use `ontime.dim_airports` as the semantic airport dimension for current airport reference data, including:

- `AirportCode`
- `DisplayAirportName`
- `Latitude`
- `Longitude`
- `UtcLocalTimeVariation`

`ontime.dim_airports` can be used directly to get coordinates and other columns for enrichment in application code (such as JavaScript or python)
or used for SQL JOINs.

Preferred sql joins:

- `ontime.fact_ontime.OriginAirportID = ontime.dim_airports.AirportID`
- `ontime.fact_ontime.DestAirportID = ontime.dim_airports.AirportID`

For `ontime.fact_ontime`, use `OriginAirportID` and `DestAirportID` when enriching airport names, coordinates, or other airport attributes inside SQL.
Do not join `ontime.fact_ontime` legs directly to `ontime.dim_airports` by airport code when the airport id columns are already available.

Fallback sql joins:

- use `ontime.fact_ontime.OriginCode = ontime.dim_airports.AirportCode`
- use `ontime.fact_ontime.DestCode = ontime.dim_airports.AirportCode`
