- Connect to clickhouse server though MCP connection
- Do not use direct HTTP by any tools like curl.
- Use the `ontime` database to answer analytical questions
- Use `ontime-semantic-layer` skill for schema inspection, join guidance, and dimension semantics.
- write correct and efficient ClickHouse SQL 
- Before writing any SQL artifact, self-verify every SQL statement you intend to save.
- Run a cheap debug execution for each query first, usually with a small `LIMIT`, a narrow `WHERE` filter, or both applied inside the main data-reading subquery or CTE.
- Treat successful execution as mandatory. Fix any syntax, type, aggregate, window, join, or unknown-column errors in a loop until every saved query runs successfully.
- Do not write unchecked SQL.

Create the presentation artifact using the proper `*-analyst-dashboard` skill.

### Rules

- Question title: `Highest daily hops for one aircraft on one flight number`
- Visual mode: `dynamic`
- Presentation target: `html`
- Visual type: `html_map`
- Derive KPIs, chart values, table rows, filters, and highlights from the actual analytical data. Do not invent or hardcode them.
- Respect the declared visual mode and visual type shown below.
- Follow question-specific visual guidance after the shared contract. Put reusable runtime behavior in shared page code, not in prose comments.

- use main proof query as the primary saved SQL already provided in the prompt
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
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline AS Carrier,
        FlightDate,
        count() AS hops,
        arrayStringConcat(
            arrayConcat(
                arraySort((x, t) -> t, groupArray(OriginCode), groupArray(assumeNotNull(CRSDepTime))),
                [argMax(DestCode, assumeNotNull(CRSDepTime))]
            ),
            '-'
        ) AS Route
    FROM ontime.fact_ontime
    WHERE Tail_Number != '' AND Cancelled = 0
    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate
),
best_by_route AS (
    SELECT
        Route,
        max(hops) AS max_hops,
        argMax(Tail_Number, FlightDate) AS aircraft_id,
        argMax(Flight_Number_Reporting_Airline, FlightDate) AS flight_number,
        argMax(Carrier, FlightDate) AS carrier,
        max(FlightDate) AS flight_date,
        argMax(hops, FlightDate) AS latest_hops
    FROM itineraries
    GROUP BY Route
)
SELECT
    aircraft_id,
    flight_number,
    carrier,
    flight_date,
    max_hops AS hop_count,
    latest_hops AS num_flights,
    Route
FROM best_by_route
ORDER BY max_hops DESC, flight_date DESC
LIMIT 10
```

Data example/snippet:

{
  "question_title": "Highest daily hops for one aircraft on one flight number",
  "result_columns": null,
  "row_count": 2,
  "mode_hint": "This visual pass receives only verified subquestion answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "q1",
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?",
      "answer_markdown": "The top itineraries are overwhelmingly recurring scheduled patterns. The highest-frequency 8-hop routes each appear dozens to over a hundred times across different dates, spanning multi-month windows. For example, the route `FLL-JAX-IND-MDW-MCI-DAL-ABQ-LAX-SJC` appeared 114 times between March and August 2008, and `CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN` (WN flight 3149) appears across multiple weeks in 2024. These are classic Southwest Airlines through-plane rotation schedules that repeat the same airport sequence day after day on fixed flight numbers.",
      "sql": "WITH itineraries AS (SELECT Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline AS Carrier, FlightDate, count() AS hops, arrayStringConcat(arrayConcat(arraySort((x, t) -\u003e t, groupArray(OriginCode), groupArray(assumeNotNull(CRSDepTime))), [argMax(DestCode, assumeNotNull(CRSDepTime))]), '-') AS Route FROM ontime.fact_ontime WHERE Tail_Number != '' AND Cancelled = 0 GROUP BY Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate), route_stats AS (SELECT Route, max(hops) AS max_hops, count() AS occurrences, min(FlightDate) AS first_seen, max(FlightDate) AS last_seen FROM itineraries GROUP BY Route) SELECT Route, max_hops, occurrences, first_seen, last_seen FROM route_stats ORDER BY max_hops DESC, occurrences DESC LIMIT 10",
      "row_count": 10,
      "result_columns": [
        "Route",
        "max_hops",
        "occurrences",
        "first_seen",
        "last_seen"
      ],
      "first_row": {
        "Route": "FLL-JAX-IND-MDW-MCI-DAL-ABQ-LAX-SJC",
        "first_seen": "2008-03-10T00:00:00Z",
        "last_seen": "2008-08-21T00:00:00Z",
        "max_hops": 8,
        "occurrences": 114
      }
    },
    {
      "id": "q2",
      "subquestion": "What geographic pattern do the top itineraries show?",
      "answer_markdown": "All top itineraries are Southwest Airlines (WN) cross-country daisy-chains that traverse the full width of the continental United States. Routes typically originate on the East Coast or Gulf Coast (e.g., FLL, PVD, BDL, LGA, HOU, MSY), pass through Midwest hubs (MDW, STL, MCI, CMH), continue through Texas/Southwest hubs (DAL, ABQ, PHX, ELP), and terminate on the West Coast (SJC, OAK, ONT, SEA, RNO). Each 8-hop route covers 6–7 different states. The pattern reflects Southwest's hub-and-spoke-free model where a single aircraft flies a continuous east-to-west (or west-to-east) rotation through its point-to-point network in a single day.",
      "sql": "WITH itineraries AS (SELECT FlightDate, count() AS hops, arrayStringConcat(arrayConcat(arraySort((x, t) -\u003e t, groupArray(OriginCode), groupArray(assumeNotNull(CRSDepTime))), [argMax(DestCode, assumeNotNull(CRSDepTime))]), '-') AS Route, arrayDistinct(groupArray(OriginState)) AS states_covered FROM ontime.fact_ontime WHERE Tail_Number != '' AND Cancelled = 0 GROUP BY FlightDate, Tail_Number, Flight_Number_Reporting_Airline HAVING hops = 8), route_stats AS (SELECT Route, count() AS occurrences, any(states_covered) AS states_sample FROM itineraries GROUP BY Route) SELECT Route, occurrences, length(states_sample) AS state_count, states_sample FROM route_stats ORDER BY occurrences DESC LIMIT 10",
      "row_count": 10,
      "result_columns": [
        "Route",
        "occurrences",
        "state_count",
        "states_sample"
      ],
      "first_row": {
        "Route": "FLL-JAX-IND-MDW-MCI-DAL-ABQ-LAX-SJC",
        "occurrences": 99,
        "state_count": 7,
        "states_sample": [
          "NM",
          "TX",
          "FL",
          "IN",
          "CA",
          "MO",
          "IL"
        ]
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
  "row_count": 2,
  "mode_hint": "This visual pass receives only verified subquestion answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "q1",
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?",
      "answer_markdown": "The top itineraries are overwhelmingly recurring scheduled patterns. The highest-frequency 8-hop routes each appear dozens to over a hundred times across different dates, spanning multi-month windows. For example, the route `FLL-JAX-IND-MDW-MCI-DAL-ABQ-LAX-SJC` appeared 114 times between March and August 2008, and `CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN` (WN flight 3149) appears across multiple weeks in 2024. These are classic Southwest Airlines through-plane rotation schedules that repeat the same airport sequence day after day on fixed flight numbers.",
      "sql": "WITH itineraries AS (SELECT Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline AS Carrier, FlightDate, count() AS hops, arrayStringConcat(arrayConcat(arraySort((x, t) -\u003e t, groupArray(OriginCode), groupArray(assumeNotNull(CRSDepTime))), [argMax(DestCode, assumeNotNull(CRSDepTime))]), '-') AS Route FROM ontime.fact_ontime WHERE Tail_Number != '' AND Cancelled = 0 GROUP BY Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate), route_stats AS (SELECT Route, max(hops) AS max_hops, count() AS occurrences, min(FlightDate) AS first_seen, max(FlightDate) AS last_seen FROM itineraries GROUP BY Route) SELECT Route, max_hops, occurrences, first_seen, last_seen FROM route_stats ORDER BY max_hops DESC, occurrences DESC LIMIT 10",
      "row_count": 10,
      "result_columns": [
        "Route",
        "max_hops",
        "occurrences",
        "first_seen",
        "last_seen"
      ],
      "first_row": {
        "Route": "FLL-JAX-IND-MDW-MCI-DAL-ABQ-LAX-SJC",
        "first_seen": "2008-03-10T00:00:00Z",
        "last_seen": "2008-08-21T00:00:00Z",
        "max_hops": 8,
        "occurrences": 114
      }
    },
    {
      "id": "q2",
      "subquestion": "What geographic pattern do the top itineraries show?",
      "answer_markdown": "All top itineraries are Southwest Airlines (WN) cross-country daisy-chains that traverse the full width of the continental United States. Routes typically originate on the East Coast or Gulf Coast (e.g., FLL, PVD, BDL, LGA, HOU, MSY), pass through Midwest hubs (MDW, STL, MCI, CMH), continue through Texas/Southwest hubs (DAL, ABQ, PHX, ELP), and terminate on the West Coast (SJC, OAK, ONT, SEA, RNO). Each 8-hop route covers 6–7 different states. The pattern reflects Southwest's hub-and-spoke-free model where a single aircraft flies a continuous east-to-west (or west-to-east) rotation through its point-to-point network in a single day.",
      "sql": "WITH itineraries AS (SELECT FlightDate, count() AS hops, arrayStringConcat(arrayConcat(arraySort((x, t) -\u003e t, groupArray(OriginCode), groupArray(assumeNotNull(CRSDepTime))), [argMax(DestCode, assumeNotNull(CRSDepTime))]), '-') AS Route, arrayDistinct(groupArray(OriginState)) AS states_covered FROM ontime.fact_ontime WHERE Tail_Number != '' AND Cancelled = 0 GROUP BY FlightDate, Tail_Number, Flight_Number_Reporting_Airline HAVING hops = 8), route_stats AS (SELECT Route, count() AS occurrences, any(states_covered) AS states_sample FROM itineraries GROUP BY Route) SELECT Route, occurrences, length(states_sample) AS state_count, states_sample FROM route_stats ORDER BY occurrences DESC LIMIT 10",
      "row_count": 10,
      "result_columns": [
        "Route",
        "occurrences",
        "state_count",
        "states_sample"
      ],
      "first_row": {
        "Route": "FLL-JAX-IND-MDW-MCI-DAL-ABQ-LAX-SJC",
        "occurrences": 99,
        "state_count": 7,
        "states_sample": [
          "NM",
          "TX",
          "FL",
          "IN",
          "CA",
          "MO",
          "IL"
        ]
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