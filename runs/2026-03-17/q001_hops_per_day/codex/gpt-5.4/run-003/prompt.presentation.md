You are running inside qforge.

Generate valid ClickHouse SQL.
Use the configured MCP server for all data access.
Do not construct raw OpenAPI URLs manually.
Stay within the configured dataset scope.
Do not reference tables outside the allowed dataset constraints.

Dataset constraints:

- Use `default.ontime_v2` as the primary fact table.
- Do not reference `default.ontime`.

You are generating presentation artifacts inside qforge.

The harness already executed the saved SQL and produced `result.json`.
Do not invent KPIs or data values.
The report is a Markdown template that qforge will render from `result.json`.

Return exactly these fenced sections:

```report
Use placeholders only where data is needed.
Allowed placeholders: {{row_count}}, {{generated_at}}, {{columns_csv}}, {{question_title}}, {{data_overview_md}}, {{result_table_md}}
```

```html
<!doctype html>
<html>...</html>
```

Report rules:

- The report must be Markdown.
- The report must be a template, not a data-filled summary.
- Prefer `{{data_overview_md}}` and `{{result_table_md}}` for JSON-derived sections.
- Keep the report concise and analytical.

Question title: `Highest daily hops for one aircraft on one flight number`

Result columns: `Aircraft ID, Flight Number, Carrier, Date, Route`

Question-specific report guidance:

Explain:

- the maximum hop count observed and whether it appears to be a repeated operating pattern or a one-off itinerary,
- the single most recent itinerary among the maximum-hop rows, including carrier, flight number, date, and full route,
- and any notable route repetition or clustering visible across the top 10 longest itineraries.

Create `visual.html` using the `ontime-analyst-dashboard` skill.

The returned `visual.html` must be final browser-ready HTML. qforge will not patch or rewrite it after generation.

General visual rules:

- Derive KPIs, chart values, table rows, filters, and highlights from the actual analytical data. Do not invent or hardcode them.
- Respect the declared visual mode and visual type shown below.
- Follow question-specific visual guidance after the shared contract. Put reusable runtime behavior in shared page code, not in prose comments.
- Keep report logic out of `visual.html`; qforge already renders `report.md` separately.

Saved SQL to preserve in the final page contract:

```sql
WITH
    groupArray((
        addMinutes(addHours(toDateTime(FlightDate), intDiv(DepTime, 100)), DepTime % 100),
        concat(formatDateTime(addMinutes(addHours(toDateTime(FlightDate), intDiv(DepTime, 100)), DepTime % 100), '%H:%i'), ' ', Origin, '->', Dest)
    )) AS legs
SELECT
    Tail_Number AS `Aircraft ID`,
    Flight_Number_Reporting_Airline AS `Flight Number`,
    IATA_CODE_Reporting_Airline AS `Carrier`,
    FlightDate AS `Date`,
    arrayStringConcat(arrayMap(x -> x.2, arraySort(x -> x.1, legs)), ' | ') AS `Route`
FROM default.ontime_v2
WHERE Cancelled = 0
  AND Diverted = 0
  AND DepTime IS NOT NULL
  AND Tail_Number != ''
  AND Flight_Number_Reporting_Airline != ''
GROUP BY
    Tail_Number,
    Flight_Number_Reporting_Airline,
    IATA_CODE_Reporting_Airline,
    FlightDate
ORDER BY
    count() DESC,
    FlightDate DESC,
    max(addMinutes(addHours(toDateTime(FlightDate), intDiv(DepTime, 100)), DepTime % 100)) DESC
LIMIT 10
```

Visual context:

- Question title: `Highest daily hops for one aircraft on one flight number`
- Visual mode: `dynamic`
- Visual type: `html_map`

Question-specific visual guidance:

The page must:

- show a lead-itinerary map that remains present even before airport-coordinate enrichment succeeds
- treat the first row returned by the primary query as the default selected itinerary on initial load
- derive hop count, stop sequence, and repeated-route comparisons from the result set
- run an explicit airport-coordinate enrichment query against `default.airports_bts` using airport codes parsed from the route strings
- label the map as airport-coordinate enrichment in the query ledger
- reuse the enrichment results for any itinerary selected from the primary result set without issuing a new per-click enrichment query
- include KPI cards for tail number, flight number, date, hop count, and route repetition context, with the date shown as its own visible KPI value
- keep the KPI strip anchored to the top-ranked result even when the selected itinerary changes
- include a legend plus both a route sequence/detail panel and an itinerary table below the map
- make itinerary table rows clickable so selecting a row redraws the map and refreshes the route sequence/detail panel for that itinerary
- show a clear active-row state for the selected itinerary that is distinct from simple hover styling
- if enrichment fails or the selected itinerary lacks enough coordinates, keep the map card visible with degraded-state messaging for that selected itinerary, report the degraded map in the ledger, and continue rendering the non-map analysis

Dynamic-mode additions:

- Build the page in dynamic mode using the `ontime-analyst-dashboard` skill contract.
- Execute the embedded saved SQL in the browser as the primary query.
- Keep the embedded saved SQL authoritative for the artifact.
- Surface every browser query in a visible query ledger.
- Do not embed the primary analytical dataset as `result.json` payloads or CSV snapshots.
- Keep additional browser queries limited to explicit enrichment or drill-down that materially improves the visualization.

The page must:

- show a lead-itinerary map that remains present even before airport-coordinate enrichment succeeds
- treat the first row returned by the primary query as the default selected itinerary on initial load
- derive hop count, stop sequence, and repeated-route comparisons from the result set
- run an explicit airport-coordinate enrichment query against `default.airports_bts` using airport codes parsed from the route strings
- label the map as airport-coordinate enrichment in the query ledger
- reuse the enrichment results for any itinerary selected from the primary result set without issuing a new per-click enrichment query
- include KPI cards for tail number, flight number, date, hop count, and route repetition context, with the date shown as its own visible KPI value
- keep the KPI strip anchored to the top-ranked result even when the selected itinerary changes
- include a legend plus both a route sequence/detail panel and an itinerary table below the map
- make itinerary table rows clickable so selecting a row redraws the map and refreshes the route sequence/detail panel for that itinerary
- show a clear active-row state for the selected itinerary that is distinct from simple hover styling
- if enrichment fails or the selected itinerary lacks enough coordinates, keep the map card visible with degraded-state messaging for that selected itinerary, report the degraded map in the ledger, and continue rendering the non-map analysis