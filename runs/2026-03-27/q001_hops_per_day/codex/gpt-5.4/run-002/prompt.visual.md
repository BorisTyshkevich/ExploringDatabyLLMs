- Connect to clickhouse server only though MCP connection
- Use the `ontime` database to answer analytical questions
- Use `ontime-semantic-layer` skill for schema inspection, join guidance, and dimension semantics.
- write correct and efficient ClickHouse SQL 
- Before writing any SQL artifact, self-verify every SQL statement you intend to save. Run a cheap debug execution for each query first, usually with a small `LIMIT`, a narrow `WHERE` filter, or both applied inside the main data-reading subquery or CTE.
- Treat successful execution as mandatory. Fix any syntax, type, aggregate, window, join, or unknown-column errors in a loop until every saved query runs successfully. Do not write not tested SQL to the artifact.

Create the presentation artifact using the proper `*-analyst-dashboard` skill.

### Rules

- Question title: `Highest daily hops for one aircraft on one flight number`
- Visual mode: `dynamic`
- Presentation target: `html`
- Visual type: `html_map`
- Derive KPIs, chart values, table rows, filters, and highlights from the actual analytical data. Do not invent or hardcode them.
- Respect the declared visual mode and visual type shown below.
- Follow question-specific visual guidance after the shared contract. Put reusable runtime behavior in shared page code, not in prose comments.

### main query

- use main query as the primary source of information for visualizing
- use the other/supporting queries when they materially improve the narrative or supporting panels
- anchor the hero narrative to the top-ranked itinerary even when another itinerary is selected in the table
- show a lead-itinerary map only after airport-coordinate lookup succeeds
- treat the first row returned by the main query as the default selected itinerary on initial load
- derive hop count, stop sequence, and repeated-route comparisons from the result set
- include a narrative hero about the lead itinerary and the broader geographic pattern of the top itineraries
- label the map as airport-coordinate lookup in the query ledger
- reuse the lookup results for any itinerary selected from the main query result set without issuing a new per-click lookup query
- include KPI cards for tail number, flight number, date, hop count, and route repetition context, with the date shown as its own visible KPI value
- keep the KPI strip synced to the currently selected itinerary
- include a legend plus both a route sequence/detail panel and an itinerary table below the map
- make itinerary table rows clickable so selecting a row redraws the map and refreshes the route sequence/detail panel for that itinerary
- make the selected-row map behavior explicit: when the selected itinerary differs from Rank 1, the map title, plotted route, markers, bounds, and route detail panel must visibly update to that selected itinerary rather than leaving the lead route drawn
- keep the map/detail/KPI selection state separate from the anchored hero state
- place selection-driven itinerary detail right after KPI cards, but before map panel.
- show a clear active-row state for the selected itinerary that is distinct from simple hover styling
- prefer the `Route` value from the main query as the per-row itinerary representation for redraws
- if airport-coordinate lookup fails or the selected itinerary lacks enough coordinates, keep the map card visible with degraded-state messaging for that selected itinerary, report the degraded map in the ledger, and continue rendering the non-map analysis
- derive the ordered itinerary sequence for map redraws and the route detail panel by splitting `Route` on `-`

### operational-stress

Which airports or legs are the main operational stress points within the top 10 unique maximum-hop itineraries?
Return per-airport and per-leg average departure delay, average arrival delay, rate of 15-plus-minute delays, and diversion incidence, and identify the stop positions most associated with disruption.

### key connectors

Which airports act as the key connectors, origins, and termini within the top 10 unique maximum-hop itineraries?
Classify airport appearances by route position and return airport code, airport name, city/state, total appearances, origin appearances, intermediate-stop appearances, final-destination appearances, and share of itineraries containing that airport.

### geographically extreme

How geographically extreme is each of the top 10 unique maximum-hop itineraries?
Return total flown distance, unique airports, unique city markets, unique states, unique local-time offsets, and whether the route is entirely domestic, then summarize which routes are the most geographically expansive.

### Data Source

SQL query for primary data source:

```sql
WITH
    (SELECT max(FlightDate) FROM ontime.fact_ontime) AS max_fd,
    addYears(max_fd, -5) AS start_fd,
    base AS (
        SELECT
            FlightDate,
            Carrier,
            FlightNum,
            if(Tail_Number = '', '', Tail_Number) AS aircraft_id,
            OriginCode,
            DestCode,
            coalesce(CRSDepTime, DepTime, 0) AS dep_sort,
            coalesce(CRSArrTime, ArrTime, 0) AS arr_sort
        FROM ontime.fact_ontime
        WHERE FlightDate > start_fd
          AND Cancelled = 0
          AND FlightNum != ''
          AND Carrier != ''
    ),
    dedup AS (
        SELECT
            FlightDate,
            Carrier,
            FlightNum,
            aircraft_id,
            dep_sort,
            min(arr_sort) AS arr_sort,
            min(OriginCode) AS OriginCode,
            min(DestCode) AS DestCode
        FROM base
        GROUP BY FlightDate, Carrier, FlightNum, aircraft_id, dep_sort
    ),
    itineraries AS (
        SELECT
            FlightDate,
            Carrier,
            FlightNum,
            aircraft_id,
            min(dep_sort) AS first_dep_time,
            arraySort(x -> (x.1, x.2, x.3, x.4), groupArray((dep_sort, arr_sort, OriginCode, DestCode))) AS legs,
            length(legs) AS hop_count,
            arrayStringConcat(arrayConcat(arrayMap(x -> x.3, legs), [legs[length(legs)].4]), '-') AS Route
        FROM dedup
        GROUP BY FlightDate, Carrier, FlightNum, aircraft_id
        HAVING hop_count > 1
    ),
    route_days AS (
        SELECT
            Route,
            countDistinct(FlightDate) AS route_recurrence_count
        FROM itineraries
        GROUP BY Route
    ),
    ranked AS (
        SELECT
            i.aircraft_id,
            i.FlightNum,
            i.Carrier,
            i.FlightDate,
            i.hop_count,
            r.route_recurrence_count,
            i.Route,
            i.first_dep_time,
            row_number() OVER (
                PARTITION BY i.Route
                ORDER BY i.FlightDate DESC, i.first_dep_time DESC, i.hop_count DESC, i.Carrier, i.FlightNum, i.aircraft_id
            ) AS rn
        FROM itineraries AS i
        INNER JOIN route_days AS r ON i.Route = r.Route
    )
SELECT
    aircraft_id,
    FlightNum AS flight_number,
    Carrier AS carrier,
    FlightDate AS flight_date,
    hop_count,
    route_recurrence_count,
    Route
FROM ranked
WHERE rn = 1
ORDER BY hop_count DESC, flight_date DESC, first_dep_time DESC, Route
LIMIT 10
```

Data example/snippet:

{
  "question_title": "Highest daily hops for one aircraft on one flight number",
  "result_columns": null,
  "row_count": 1,
  "mode_hint": "This visual pass receives only verified section answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "main",
      "subquestion": "Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.\nDefine uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.\nDo not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.\nCount hops from distinct same-day legs, not raw source rows; do not let duplicate or conflicting same-time rows inflate hop count or create artifact routes.\n\nReturn:\n\n- aircraft id\n- flight number\n- carrier\n- flight date\n- hop count\n- route recurrence count: total number of days across the analyzed window on which this exact Route string was flown by any aircraft\n- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`",
      "answer_markdown": "Across the most recent five years in the OnTime data, the longest same-aircraft, same-flight-number daily itineraries reached 8 hops. The 10 most recent unique routes at that maximum include ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA on 2024-12-01 (route recurrence 1 day), CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN on 2024-02-18 (4 days), and ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN on 2023-04-30 (2 days).",
      "sql": "WITH\n    (SELECT max(FlightDate) FROM ontime.fact_ontime) AS max_fd,\n    addYears(max_fd, -5) AS start_fd,\n    base AS (\n        SELECT\n            FlightDate,\n            Carrier,\n            FlightNum,\n            if(Tail_Number = '', '', Tail_Number) AS aircraft_id,\n            OriginCode,\n            DestCode,\n            coalesce(CRSDepTime, DepTime, 0) AS dep_sort,\n            coalesce(CRSArrTime, ArrTime, 0) AS arr_sort\n        FROM ontime.fact_ontime\n        WHERE FlightDate \u003e start_fd\n          AND Cancelled = 0\n          AND FlightNum != ''\n          AND Carrier != ''\n    ),\n    dedup AS (\n        SELECT\n            FlightDate,\n            Carrier,\n            FlightNum,\n            aircraft_id,\n            dep_sort,\n            min(arr_sort) AS arr_sort,\n            min(OriginCode) AS OriginCode,\n            min(DestCode) AS DestCode\n        FROM base\n        GROUP BY FlightDate, Carrier, FlightNum, aircraft_id, dep_sort\n    ),\n    itineraries AS (\n        SELECT\n            FlightDate,\n            Carrier,\n            FlightNum,\n            aircraft_id,\n            min(dep_sort) AS first_dep_time,\n            arraySort(x -\u003e (x.1, x.2, x.3, x.4), groupArray((dep_sort, arr_sort, OriginCode, DestCode))) AS legs,\n            length(legs) AS hop_count,\n            arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, legs), [legs[length(legs)].4]), '-') AS Route\n        FROM dedup\n        GROUP BY FlightDate, Carrier, FlightNum, aircraft_id\n        HAVING hop_count \u003e 1\n    ),\n    route_days AS (\n        SELECT\n            Route,\n            countDistinct(FlightDate) AS route_recurrence_count\n        FROM itineraries\n        GROUP BY Route\n    ),\n    ranked AS (\n        SELECT\n            i.aircraft_id,\n            i.FlightNum,\n            i.Carrier,\n            i.FlightDate,\n            i.hop_count,\n            r.route_recurrence_count,\n            i.Route,\n            i.first_dep_time,\n            row_number() OVER (\n                PARTITION BY i.Route\n                ORDER BY i.FlightDate DESC, i.first_dep_time DESC, i.hop_count DESC, i.Carrier, i.FlightNum, i.aircraft_id\n            ) AS rn\n        FROM itineraries AS i\n        INNER JOIN route_days AS r ON i.Route = r.Route\n    )\nSELECT\n    aircraft_id,\n    FlightNum AS flight_number,\n    Carrier AS carrier,\n    FlightDate AS flight_date,\n    hop_count,\n    route_recurrence_count,\n    Route\nFROM ranked\nWHERE rn = 1\nORDER BY hop_count DESC, flight_date DESC, first_dep_time DESC, Route\nLIMIT 10",
      "is_primary": true,
      "date_field_hint": "flight_date",
      "row_count": 10,
      "result_columns": [
        "aircraft_id",
        "flight_number",
        "carrier",
        "flight_date",
        "hop_count",
        "route_recurrence_count",
        "Route"
      ],
      "first_row": {
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "aircraft_id": "N957WN",
        "carrier": "WN",
        "flight_date": "2024-12-01T00:00:00Z",
        "flight_number": "366",
        "hop_count": 8,
        "route_recurrence_count": 1
      }
    }
  ]
}

### Multi-query additions

- The dashboard should combine narrative and interactive analysis.
- The verified analysis package includes the main query and named supporting section queries that may be used for building visual panels for particular questions.
- Prefill editable SQL panels for all queries in the verified package.
- Use section answers as narrative framing, but derive displayed KPIs, charts, tables, and interactions from live browser execution of the primary saved SQL and any supporting queries you actually run.

### additional questions

- When the question-specific visual prompt explicitly asks an additional question, author a lookup query and use it for building a visual panel
- label that panel query in the query ledger with the header name
- additional lookup query use the currently selected context from the primary/main query result set and fetches additional data to enrich
- lookup query should use a filter to read only rows from the main set that need to be enriched with additional information. 
- example WHERE clause: `where (c1,c2,c3) in ( ('vc1','vc2', 'vc3'), ('wc1','wc2', 'wc3') )`
- Don't try to calculate the main row set again by a complicated SQL subquery/CTE.
- When the date selector changes, rerun the main query and every enrichment query to support the active date range.
- if the enrichment lookup fails, keep the panel visible with degraded-state messaging, report that degraded panel in the ledger, and continue rendering the rest of the dashboard
- If you run supporting or lookup queries, record them in the same visible query ledger as the main query.

### Verified Analysis Package

Use this JSON package as the supporting context for the visual:

{
  "question_title": "Highest daily hops for one aircraft on one flight number",
  "result_columns": null,
  "row_count": 1,
  "mode_hint": "This visual pass receives only verified section answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "main",
      "subquestion": "Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.\nDefine uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.\nDo not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.\nCount hops from distinct same-day legs, not raw source rows; do not let duplicate or conflicting same-time rows inflate hop count or create artifact routes.\n\nReturn:\n\n- aircraft id\n- flight number\n- carrier\n- flight date\n- hop count\n- route recurrence count: total number of days across the analyzed window on which this exact Route string was flown by any aircraft\n- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`",
      "answer_markdown": "Across the most recent five years in the OnTime data, the longest same-aircraft, same-flight-number daily itineraries reached 8 hops. The 10 most recent unique routes at that maximum include ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA on 2024-12-01 (route recurrence 1 day), CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN on 2024-02-18 (4 days), and ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN on 2023-04-30 (2 days).",
      "sql": "WITH\n    (SELECT max(FlightDate) FROM ontime.fact_ontime) AS max_fd,\n    addYears(max_fd, -5) AS start_fd,\n    base AS (\n        SELECT\n            FlightDate,\n            Carrier,\n            FlightNum,\n            if(Tail_Number = '', '', Tail_Number) AS aircraft_id,\n            OriginCode,\n            DestCode,\n            coalesce(CRSDepTime, DepTime, 0) AS dep_sort,\n            coalesce(CRSArrTime, ArrTime, 0) AS arr_sort\n        FROM ontime.fact_ontime\n        WHERE FlightDate \u003e start_fd\n          AND Cancelled = 0\n          AND FlightNum != ''\n          AND Carrier != ''\n    ),\n    dedup AS (\n        SELECT\n            FlightDate,\n            Carrier,\n            FlightNum,\n            aircraft_id,\n            dep_sort,\n            min(arr_sort) AS arr_sort,\n            min(OriginCode) AS OriginCode,\n            min(DestCode) AS DestCode\n        FROM base\n        GROUP BY FlightDate, Carrier, FlightNum, aircraft_id, dep_sort\n    ),\n    itineraries AS (\n        SELECT\n            FlightDate,\n            Carrier,\n            FlightNum,\n            aircraft_id,\n            min(dep_sort) AS first_dep_time,\n            arraySort(x -\u003e (x.1, x.2, x.3, x.4), groupArray((dep_sort, arr_sort, OriginCode, DestCode))) AS legs,\n            length(legs) AS hop_count,\n            arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, legs), [legs[length(legs)].4]), '-') AS Route\n        FROM dedup\n        GROUP BY FlightDate, Carrier, FlightNum, aircraft_id\n        HAVING hop_count \u003e 1\n    ),\n    route_days AS (\n        SELECT\n            Route,\n            countDistinct(FlightDate) AS route_recurrence_count\n        FROM itineraries\n        GROUP BY Route\n    ),\n    ranked AS (\n        SELECT\n            i.aircraft_id,\n            i.FlightNum,\n            i.Carrier,\n            i.FlightDate,\n            i.hop_count,\n            r.route_recurrence_count,\n            i.Route,\n            i.first_dep_time,\n            row_number() OVER (\n                PARTITION BY i.Route\n                ORDER BY i.FlightDate DESC, i.first_dep_time DESC, i.hop_count DESC, i.Carrier, i.FlightNum, i.aircraft_id\n            ) AS rn\n        FROM itineraries AS i\n        INNER JOIN route_days AS r ON i.Route = r.Route\n    )\nSELECT\n    aircraft_id,\n    FlightNum AS flight_number,\n    Carrier AS carrier,\n    FlightDate AS flight_date,\n    hop_count,\n    route_recurrence_count,\n    Route\nFROM ranked\nWHERE rn = 1\nORDER BY hop_count DESC, flight_date DESC, first_dep_time DESC, Route\nLIMIT 10",
      "is_primary": true,
      "date_field_hint": "flight_date",
      "row_count": 10,
      "result_columns": [
        "aircraft_id",
        "flight_number",
        "carrier",
        "flight_date",
        "hop_count",
        "route_recurrence_count",
        "Route"
      ],
      "first_row": {
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "aircraft_id": "N957WN",
        "carrier": "WN",
        "flight_date": "2024-12-01T00:00:00Z",
        "flight_number": "366",
        "hop_count": 8,
        "route_recurrence_count": 1
      }
    }
  ]
}

### Dynamic-mode additions

- Use this endpoint template for every browser query: `https://mcp.demo.altinity.cloud/{JWE}/openapi/execute_query?query=...`
- Keep JWE in `localStorage['OnTimeAnalystDashboard::auth::jwe']`.
- Do not embed the primary analytical dataset as `result.json` payloads or CSV snapshots.
- Provide a visible start/end date selector for the analytical range.
- Drive the date selector through SQL reruns, not client-side filtering alone.
- Default the visual to the most recent 5 years relative to the latest available analytical date when a usable date field exists.
- If no safe date field can be detected for a query, keep the selector visible but disable it with a clear note for that query or view.
- Expose editable SQL controls for the primary query and every supporting query the page uses.
- Provide individual run buttons for editable queries and a `Run all` path when the page uses multiple queries.
- Keep one unified query ledger that records each execution with query id, role, effective date range, status, rows, and expandable SQL text.
- Prefer an explicit SQL wrapping or parameter-insertion strategy for date predicates instead of brittle string replacement.
- If a supporting query cannot be safely date-parameterized, keep it editable and manually runnable, and surface that limitation in the UI.
- Saved supporting queries come from reviewed analysis artifacts; lookup queries are visual-only second-pass queries authored directly into the page for a concrete panel, lookup, or drill-down need.
- A lookup query may depend on the currently selected primary-result context or current dashboard state when that dependency is explicit in the UI.
- If a lookup query fails, degrade only the dependent panel, keep the primary analysis visible, and record the lookup failure in the ledger and visible status UI.
- Before writing `visual.html`, self-verify every browser-side SQL statement you intend to ship, including primary, supporting, enrichment, drill-down, and lookup queries.
- For each query, run a cheap live check against the real endpoint and schema first, usually with a small `LIMIT`, a narrow `WHERE` filter, or both when that preserves the query shape.
- Treat successful execution as mandatory. Fix any syntax, type, aggregate, join, or unknown-column errors in a loop until every shipped browser query runs successfully.

Create browser-ready HTML `visual.html`.

Write the file or provide a download link. Do not include the HTML source in the response. Do not open the artifact view frame.