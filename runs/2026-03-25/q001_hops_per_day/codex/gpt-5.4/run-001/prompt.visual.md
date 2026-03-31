- Connect to clickhouse server though MCP connection
- Do not use direct HTTP by any tools like curl.
- Use the `ontime` database to answer analytical questions
- Use `ontime-semantic-layer` skill for schema inspection, join guidance, and dimension semantics.
- write correct and efficient ClickHouse SQL 
- Before finalizing your answer, self-verify the query with a quick debug execution, usually with a small `LIMIT` or `WHERE` filter in a data reading subquery or CTE. Fix any errors in a loop until done.

Create the presentation artifact using the proper `*-analyst-dashboard` skill.

### Rules

- Question title: `Highest daily hops for one aircraft on one flight number`
- Visual mode: `dynamic`
- Presentation target: `html`
- Visual type: `html_map`
- Derive KPIs, chart values, table rows, filters, and highlights from the actual analytical data. Do not invent or hardcode them.
- Respect the declared visual mode and visual type shown below.
- Follow question-specific visual guidance after the shared contract. Put reusable runtime behavior in shared page code, not in prose comments.

- use the first dashboard-question proof query as the primary saved SQL already provided in the prompt
- use the other dashboard-question proof queries as supporting queries when they materially improve the narrative or supporting panels
- anchor the hero narrative and KPI strip to the top-ranked itinerary even when another itinerary is selected in the table
- show a lead-itinerary map that remains present even before airport-coordinate enrichment succeeds
- treat the first row returned by the primary query as the default selected itinerary on initial load
- derive hop count, stop sequence, and repeated-route comparisons from the result set
- include a narrative hero about the lead itinerary and the broader geographic pattern of the top itineraries
- label the map as airport-coordinate enrichment in the query ledger
- reuse the enrichment results for any itinerary selected from the primary result set without issuing a new per-click enrichment query
- include KPI cards for tail number, flight number, date, hop count, and route repetition context, with the date shown as its own visible KPI value
- keep the KPI strip anchored to the top-ranked result even when the selected itinerary changes
- include a legend plus both a route sequence/detail panel and an itinerary table below the map
- make itinerary table rows clickable so selecting a row redraws the map and refreshes the route sequence/detail panel for that itinerary
- make the selected-row map behavior explicit: when the selected itinerary differs from Rank 1, the map title, plotted route, markers, bounds, and route detail panel must visibly update to that selected itinerary rather than leaving the lead route drawn
- keep the map/detail selection state separate from the anchored hero and KPI state
- show a clear active-row state for the selected itinerary that is distinct from simple hover styling
- prefer the `Route` value from the primary query as the per-row itinerary representation for redraws
- if enrichment fails or the selected itinerary lacks enough coordinates, keep the map card visible with degraded-state messaging for that selected itinerary, report the degraded map in the ledger, and continue rendering the non-map analysis
- derive the ordered itinerary sequence for map redraws and the route detail panel by splitting `Route` on `-`

### Data Source

SQL query for primary data source:

```sql
WITH itineraries AS (
    SELECT
        FlightDate AS flight_date,
        TailNum AS aircraft_id,
        FlightNum AS flight_number,
        Carrier AS carrier,
        arraySort(
            x -> (x.1, x.2, x.3),
            groupArray((coalesce(DepTime, CRSDepTime, toUInt16(0)), OriginCode, DestCode))
        ) AS legs
    FROM ontime.fact_ontime
    WHERE Cancelled = 0
      AND TailNum != ''
      AND FlightNum != ''
    GROUP BY
        flight_date,
        aircraft_id,
        flight_number,
        carrier
    HAVING count() > 1
),
ranked AS (
    SELECT
        flight_date,
        aircraft_id,
        flight_number,
        carrier,
        length(legs) AS hop_count,
        arrayStringConcat(arrayMap(x -> x.2, legs), '-') || '-' || arrayElement(legs, -1).3 AS Route
    FROM itineraries
),
max_hops AS (
    SELECT max(hop_count) AS max_hop_count
    FROM ranked
),
unique_routes AS (
    SELECT
        flight_date,
        aircraft_id,
        flight_number,
        carrier,
        hop_count,
        Route,
        row_number() OVER (
            PARTITION BY Route
            ORDER BY flight_date DESC, aircraft_id DESC, flight_number DESC, carrier DESC
        ) AS route_rank
    FROM ranked
    WHERE hop_count = (SELECT max_hop_count FROM max_hops)
)
SELECT
    aircraft_id,
    flight_number,
    carrier,
    flight_date,
    hop_count,
    Route
FROM unique_routes
WHERE route_rank = 1
ORDER BY flight_date DESC, Route ASC
LIMIT 10
```

Data example/snippet:

{
  "question_title": "Highest daily hops for one aircraft on one flight number",
  "result_columns": null,
  "row_count": 3,
  "mode_hint": "This visual pass receives only verified subquestion answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "q1",
      "subquestion": "Which itinerary is the highest-hop example, and what does it look like?",
      "answer_markdown": "The highest-hop example is an 8-hop Southwest itinerary flown by aircraft N957WN as flight 366 on 2024-12-01, with route ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA.",
      "sql": "WITH itineraries AS ( SELECT FlightDate AS flight_date, TailNum AS aircraft_id, FlightNum AS flight_number, Carrier AS carrier, arraySort(x -\u003e (x.1, x.2, x.3), groupArray((coalesce(DepTime, CRSDepTime, toUInt16(0)), OriginCode, DestCode))) AS legs FROM ontime.fact_ontime WHERE Cancelled = 0 AND TailNum != '' AND FlightNum != '' GROUP BY flight_date, aircraft_id, flight_number, carrier HAVING count() \u003e 1 ), ranked AS ( SELECT flight_date, aircraft_id, flight_number, carrier, length(legs) AS hop_count, arrayStringConcat(arrayMap(x -\u003e x.2, legs), '-') || '-' || arrayElement(legs, -1).3 AS Route FROM itineraries ), max_hops AS ( SELECT max(hop_count) AS max_hop_count FROM ranked ), unique_routes AS ( SELECT flight_date, aircraft_id, flight_number, carrier, hop_count, Route, row_number() OVER (PARTITION BY Route ORDER BY flight_date DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rank FROM ranked WHERE hop_count = (SELECT max_hop_count FROM max_hops) ) SELECT aircraft_id, flight_number, carrier, flight_date, hop_count, Route FROM unique_routes WHERE route_rank = 1 ORDER BY flight_date DESC, Route ASC LIMIT 1",
      "row_count": 1,
      "result_columns": [
        "aircraft_id",
        "flight_number",
        "carrier",
        "flight_date",
        "hop_count",
        "Route"
      ],
      "first_row": {
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "aircraft_id": "N957WN",
        "carrier": "WN",
        "flight_date": "2024-12-01T00:00:00Z",
        "flight_number": "366",
        "hop_count": 8
      }
    },
    {
      "id": "q2",
      "subquestion": "Which of the top-ranked itineraries is the most recent?",
      "answer_markdown": "The most recent top-ranked itinerary is also the latest 8-hop unique route: N957WN operating Southwest flight 366 on 2024-12-01 over ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA.",
      "sql": "WITH itineraries AS ( SELECT FlightDate AS flight_date, TailNum AS aircraft_id, FlightNum AS flight_number, Carrier AS carrier, arraySort(x -\u003e (x.1, x.2, x.3), groupArray((coalesce(DepTime, CRSDepTime, toUInt16(0)), OriginCode, DestCode))) AS legs FROM ontime.fact_ontime WHERE Cancelled = 0 AND TailNum != '' AND FlightNum != '' GROUP BY flight_date, aircraft_id, flight_number, carrier HAVING count() \u003e 1 ), ranked AS ( SELECT flight_date, aircraft_id, flight_number, carrier, length(legs) AS hop_count, arrayStringConcat(arrayMap(x -\u003e x.2, legs), '-') || '-' || arrayElement(legs, -1).3 AS Route FROM itineraries ), max_hops AS ( SELECT max(hop_count) AS max_hop_count FROM ranked ), unique_routes AS ( SELECT flight_date, aircraft_id, flight_number, carrier, hop_count, Route, row_number() OVER (PARTITION BY Route ORDER BY flight_date DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rank FROM ranked WHERE hop_count = (SELECT max_hop_count FROM max_hops) ) SELECT aircraft_id, flight_number, carrier, flight_date, hop_count, Route FROM unique_routes WHERE route_rank = 1 ORDER BY flight_date DESC, Route ASC LIMIT 1",
      "row_count": 1,
      "result_columns": [
        "aircraft_id",
        "flight_number",
        "carrier",
        "flight_date",
        "hop_count",
        "Route"
      ],
      "first_row": {
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "aircraft_id": "N957WN",
        "carrier": "WN",
        "flight_date": "2024-12-01T00:00:00Z",
        "flight_number": "366",
        "hop_count": 8
      }
    },
    {
      "id": "q3",
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?",
      "answer_markdown": "They look mostly like recurring patterns rather than one-offs. Among the 10 most recent unique 8-hop routes, 9 recur on multiple dates, averaging 14.1 occurrences each, although the single most recent example appears only once.",
      "sql": "WITH itineraries AS ( SELECT FlightDate AS flight_date, TailNum AS aircraft_id, FlightNum AS flight_number, Carrier AS carrier, arraySort(x -\u003e (x.1, x.2, x.3), groupArray((coalesce(DepTime, CRSDepTime, toUInt16(0)), OriginCode, DestCode))) AS legs FROM ontime.fact_ontime WHERE Cancelled = 0 AND TailNum != '' AND FlightNum != '' GROUP BY flight_date, aircraft_id, flight_number, carrier HAVING count() \u003e 1 ), ranked AS ( SELECT flight_date, aircraft_id, flight_number, carrier, length(legs) AS hop_count, arrayStringConcat(arrayMap(x -\u003e x.2, legs), '-') || '-' || arrayElement(legs, -1).3 AS Route FROM itineraries ), max_hops AS ( SELECT max(hop_count) AS max_hop_count FROM ranked ), unique_routes AS ( SELECT flight_date, aircraft_id, flight_number, carrier, hop_count, Route, row_number() OVER (PARTITION BY Route ORDER BY flight_date DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rank FROM ranked WHERE hop_count = (SELECT max_hop_count FROM max_hops) ), top10 AS ( SELECT flight_date, aircraft_id, flight_number, carrier, hop_count, Route FROM unique_routes WHERE route_rank = 1 ORDER BY flight_date DESC, Route ASC LIMIT 10 ), route_occurrences AS ( SELECT t.Route, count() AS occurrences FROM top10 t INNER JOIN ranked r ON r.Route = t.Route AND r.hop_count = t.hop_count GROUP BY t.Route ) SELECT count() AS route_count, countIf(occurrences = 1) AS one_off_routes, countIf(occurrences \u003e 1) AS recurring_routes, round(avg(occurrences), 2) AS avg_occurrences, max(occurrences) AS max_occurrences FROM route_occurrences",
      "row_count": 1,
      "result_columns": [
        "route_count",
        "one_off_routes",
        "recurring_routes",
        "avg_occurrences",
        "max_occurrences"
      ],
      "first_row": {
        "avg_occurrences": 14.1,
        "max_occurrences": 46,
        "one_off_routes": 1,
        "recurring_routes": 9,
        "route_count": 10
      }
    }
  ]
}

### Multi-query additions

- The saved SQL shown below is the primary dashboard query for this page.
- The verified analysis package includes named supporting queries that may be used for enrichment, drill-down, or secondary visuals when the question-specific prompt calls for them.
- Use subquestion answers as narrative framing, but derive displayed KPIs, charts, tables, and interactions from live browser execution of the primary saved SQL and any supporting queries you actually run.
- If you run supporting queries, record them in the same visible query ledger as the primary query.
- The dashboard does not need to mirror `report.md`; it should combine narrative and interactive analysis.

### Verified Analysis Package

Use this JSON package as the supporting context for the visual:

{
  "question_title": "Highest daily hops for one aircraft on one flight number",
  "result_columns": null,
  "row_count": 3,
  "mode_hint": "This visual pass receives only verified subquestion answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "q1",
      "subquestion": "Which itinerary is the highest-hop example, and what does it look like?",
      "answer_markdown": "The highest-hop example is an 8-hop Southwest itinerary flown by aircraft N957WN as flight 366 on 2024-12-01, with route ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA.",
      "sql": "WITH itineraries AS ( SELECT FlightDate AS flight_date, TailNum AS aircraft_id, FlightNum AS flight_number, Carrier AS carrier, arraySort(x -\u003e (x.1, x.2, x.3), groupArray((coalesce(DepTime, CRSDepTime, toUInt16(0)), OriginCode, DestCode))) AS legs FROM ontime.fact_ontime WHERE Cancelled = 0 AND TailNum != '' AND FlightNum != '' GROUP BY flight_date, aircraft_id, flight_number, carrier HAVING count() \u003e 1 ), ranked AS ( SELECT flight_date, aircraft_id, flight_number, carrier, length(legs) AS hop_count, arrayStringConcat(arrayMap(x -\u003e x.2, legs), '-') || '-' || arrayElement(legs, -1).3 AS Route FROM itineraries ), max_hops AS ( SELECT max(hop_count) AS max_hop_count FROM ranked ), unique_routes AS ( SELECT flight_date, aircraft_id, flight_number, carrier, hop_count, Route, row_number() OVER (PARTITION BY Route ORDER BY flight_date DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rank FROM ranked WHERE hop_count = (SELECT max_hop_count FROM max_hops) ) SELECT aircraft_id, flight_number, carrier, flight_date, hop_count, Route FROM unique_routes WHERE route_rank = 1 ORDER BY flight_date DESC, Route ASC LIMIT 1",
      "row_count": 1,
      "result_columns": [
        "aircraft_id",
        "flight_number",
        "carrier",
        "flight_date",
        "hop_count",
        "Route"
      ],
      "first_row": {
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "aircraft_id": "N957WN",
        "carrier": "WN",
        "flight_date": "2024-12-01T00:00:00Z",
        "flight_number": "366",
        "hop_count": 8
      }
    },
    {
      "id": "q2",
      "subquestion": "Which of the top-ranked itineraries is the most recent?",
      "answer_markdown": "The most recent top-ranked itinerary is also the latest 8-hop unique route: N957WN operating Southwest flight 366 on 2024-12-01 over ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA.",
      "sql": "WITH itineraries AS ( SELECT FlightDate AS flight_date, TailNum AS aircraft_id, FlightNum AS flight_number, Carrier AS carrier, arraySort(x -\u003e (x.1, x.2, x.3), groupArray((coalesce(DepTime, CRSDepTime, toUInt16(0)), OriginCode, DestCode))) AS legs FROM ontime.fact_ontime WHERE Cancelled = 0 AND TailNum != '' AND FlightNum != '' GROUP BY flight_date, aircraft_id, flight_number, carrier HAVING count() \u003e 1 ), ranked AS ( SELECT flight_date, aircraft_id, flight_number, carrier, length(legs) AS hop_count, arrayStringConcat(arrayMap(x -\u003e x.2, legs), '-') || '-' || arrayElement(legs, -1).3 AS Route FROM itineraries ), max_hops AS ( SELECT max(hop_count) AS max_hop_count FROM ranked ), unique_routes AS ( SELECT flight_date, aircraft_id, flight_number, carrier, hop_count, Route, row_number() OVER (PARTITION BY Route ORDER BY flight_date DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rank FROM ranked WHERE hop_count = (SELECT max_hop_count FROM max_hops) ) SELECT aircraft_id, flight_number, carrier, flight_date, hop_count, Route FROM unique_routes WHERE route_rank = 1 ORDER BY flight_date DESC, Route ASC LIMIT 1",
      "row_count": 1,
      "result_columns": [
        "aircraft_id",
        "flight_number",
        "carrier",
        "flight_date",
        "hop_count",
        "Route"
      ],
      "first_row": {
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "aircraft_id": "N957WN",
        "carrier": "WN",
        "flight_date": "2024-12-01T00:00:00Z",
        "flight_number": "366",
        "hop_count": 8
      }
    },
    {
      "id": "q3",
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?",
      "answer_markdown": "They look mostly like recurring patterns rather than one-offs. Among the 10 most recent unique 8-hop routes, 9 recur on multiple dates, averaging 14.1 occurrences each, although the single most recent example appears only once.",
      "sql": "WITH itineraries AS ( SELECT FlightDate AS flight_date, TailNum AS aircraft_id, FlightNum AS flight_number, Carrier AS carrier, arraySort(x -\u003e (x.1, x.2, x.3), groupArray((coalesce(DepTime, CRSDepTime, toUInt16(0)), OriginCode, DestCode))) AS legs FROM ontime.fact_ontime WHERE Cancelled = 0 AND TailNum != '' AND FlightNum != '' GROUP BY flight_date, aircraft_id, flight_number, carrier HAVING count() \u003e 1 ), ranked AS ( SELECT flight_date, aircraft_id, flight_number, carrier, length(legs) AS hop_count, arrayStringConcat(arrayMap(x -\u003e x.2, legs), '-') || '-' || arrayElement(legs, -1).3 AS Route FROM itineraries ), max_hops AS ( SELECT max(hop_count) AS max_hop_count FROM ranked ), unique_routes AS ( SELECT flight_date, aircraft_id, flight_number, carrier, hop_count, Route, row_number() OVER (PARTITION BY Route ORDER BY flight_date DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rank FROM ranked WHERE hop_count = (SELECT max_hop_count FROM max_hops) ), top10 AS ( SELECT flight_date, aircraft_id, flight_number, carrier, hop_count, Route FROM unique_routes WHERE route_rank = 1 ORDER BY flight_date DESC, Route ASC LIMIT 10 ), route_occurrences AS ( SELECT t.Route, count() AS occurrences FROM top10 t INNER JOIN ranked r ON r.Route = t.Route AND r.hop_count = t.hop_count GROUP BY t.Route ) SELECT count() AS route_count, countIf(occurrences = 1) AS one_off_routes, countIf(occurrences \u003e 1) AS recurring_routes, round(avg(occurrences), 2) AS avg_occurrences, max(occurrences) AS max_occurrences FROM route_occurrences",
      "row_count": 1,
      "result_columns": [
        "route_count",
        "one_off_routes",
        "recurring_routes",
        "avg_occurrences",
        "max_occurrences"
      ],
      "first_row": {
        "avg_occurrences": 14.1,
        "max_occurrences": 46,
        "one_off_routes": 1,
        "recurring_routes": 9,
        "route_count": 10
      }
    }
  ]
}

### Dynamic-mode additions

- Use this endpoint template for every browser query: `https://mcp.demo.altinity.cloud/{JWE}/openapi/execute_query?query=...`
- Keep JWE in `localStorage['OnTimeAnalystDashboard::auth::jwe']`.
- Do not embed the primary analytical dataset as `result.json` payloads or CSV snapshots.

Create browser-ready HTML `visual.html`.

Write the file or provide a download link. Do not include the HTML source in the response. Do not open the artifact view frame.