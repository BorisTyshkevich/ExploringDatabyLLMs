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
- use the other section proof queries as supporting queries when they materially improve the narrative or supporting panels
- anchor the hero narrative to the top-ranked itinerary even when another itinerary is selected in the table
- show a lead-itinerary map that remains present even before airport-coordinate enrichment succeeds
- treat the first row returned by the primary query as the default selected itinerary on initial load
- derive hop count, stop sequence, and repeated-route comparisons from the result set
- include a narrative hero about the lead itinerary and the broader geographic pattern of the top itineraries
- label the map as airport-coordinate enrichment in the query ledger
- reuse the enrichment results for any itinerary selected from the primary result set without issuing a new per-click enrichment query
- include KPI cards for tail number, flight number, date, hop count, and route repetition context, with the date shown as its own visible KPI value
- keep the KPI strip synced to the currently selected itinerary
- include a legend plus both a route sequence/detail panel and an itinerary table below the map
- make itinerary table rows clickable so selecting a row redraws the map and refreshes the route sequence/detail panel for that itinerary
- make the selected-row map behavior explicit: when the selected itinerary differs from Rank 1, the map title, plotted route, markers, bounds, and route detail panel must visibly update to that selected itinerary rather than leaving the lead route drawn
- keep the map/detail/KPI selection state separate from the anchored hero state
- show a clear active-row state for the selected itinerary that is distinct from simple hover styling
- prefer the `Route` value from the primary query as the per-row itinerary representation for redraws
- if enrichment fails or the selected itinerary lacks enough coordinates, keep the map card visible with degraded-state messaging for that selected itinerary, report the degraded map in the ledger, and continue rendering the non-map analysis
- derive the ordered itinerary sequence for map redraws and the route detail panel by splitting `Route` on `-`

### Data Source

SQL query for primary data source:

```sql
WITH dedup_legs AS (
    SELECT
        FlightDate,
        ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,
        Flight_Number_Reporting_Airline AS flight_number,
        IATA_CODE_Reporting_Airline AS carrier,
        OriginCode,
        DestCode,
        coalesce(CRSDepTime, DepTime, toUInt16(9999)) AS dep_sort,
        min(coalesce(CRSArrTime, ArrTime, toUInt16(9999))) AS arr_sort
    FROM ontime.fact_ontime
    WHERE Cancelled = 0
    GROUP BY FlightDate, aircraft_id, flight_number, carrier, OriginCode, DestCode, dep_sort
), itineraries AS (
    SELECT
        FlightDate,
        aircraft_id,
        flight_number,
        carrier,
        count() AS hop_count,
        min(dep_sort) AS first_dep_sort,
        arraySort(groupArray((dep_sort, arr_sort, OriginCode, DestCode))) AS legs
    FROM dedup_legs
    GROUP BY FlightDate, aircraft_id, flight_number, carrier
    HAVING count() > 1
), route_recurrence AS (
    SELECT
        arrayStringConcat(arrayConcat(arrayMap(x -> x.3, legs), [arrayElement(legs, length(legs)).4]), '-') AS route,
        uniqExact(FlightDate) AS route_recurrence_count
    FROM itineraries
    GROUP BY route
), max_hops AS (
    SELECT max(hop_count) AS max_hop_count FROM itineraries
), top_unique AS (
    SELECT
        i.FlightDate,
        i.aircraft_id,
        i.flight_number,
        i.carrier,
        i.hop_count,
        i.first_dep_sort,
        arrayStringConcat(arrayConcat(arrayMap(x -> x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') AS route,
        rr.route_recurrence_count,
        row_number() OVER (
            PARTITION BY arrayStringConcat(arrayConcat(arrayMap(x -> x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-')
            ORDER BY i.FlightDate DESC, i.first_dep_sort DESC, i.carrier DESC, i.flight_number DESC, i.aircraft_id DESC
        ) AS route_rank
    FROM itineraries i
    CROSS JOIN max_hops m
    INNER JOIN route_recurrence rr
        ON arrayStringConcat(arrayConcat(arrayMap(x -> x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') = rr.route
    WHERE i.hop_count = m.max_hop_count
)
SELECT
    if(aircraft_id = '', 'unknown', aircraft_id) AS aircraft_id,
    flight_number,
    carrier,
    FlightDate AS flight_date,
    hop_count,
    route_recurrence_count,
    route
FROM top_unique
WHERE route_rank = 1
ORDER BY flight_date DESC, first_dep_sort DESC, carrier DESC, flight_number DESC, aircraft_id DESC
LIMIT 10
```

Data example/snippet:

{
  "question_title": "Highest daily hops for one aircraft on one flight number",
  "result_columns": null,
  "row_count": 3,
  "mode_hint": "This visual pass receives only verified section answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "main",
      "subquestion": "Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.\nDefine uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.\nDo not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.\nCount hops from distinct same-day legs, not raw source rows; do not let duplicate or conflicting same-time rows inflate hop count or create artifact routes.\n\nReturn:\n\n- aircraft id\n- flight number\n- carrier\n- flight date\n- hop count\n- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft\n- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`",
      "answer_markdown": "The maximum observed same-aircraft, same-flight-number itinerary length is 8 hops. The 10 most recent unique max-hop routes are all Southwest (WN) itineraries, with the newest on 2024-12-01: ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA flown by N957WN on flight 366. Across these 10 rows, route recurrence ranges from 1 day to 46 days, led by LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN (46), MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC (40), HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK (20), and BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK (12).",
      "sql": "WITH dedup_legs AS (\n    SELECT\n        FlightDate,\n        ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,\n        Flight_Number_Reporting_Airline AS flight_number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginCode,\n        DestCode,\n        coalesce(CRSDepTime, DepTime, toUInt16(9999)) AS dep_sort,\n        min(coalesce(CRSArrTime, ArrTime, toUInt16(9999))) AS arr_sort\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier, OriginCode, DestCode, dep_sort\n), itineraries AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        count() AS hop_count,\n        min(dep_sort) AS first_dep_sort,\n        arraySort(groupArray((dep_sort, arr_sort, OriginCode, DestCode))) AS legs\n    FROM dedup_legs\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier\n    HAVING count() \u003e 1\n), route_recurrence AS (\n    SELECT\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, legs), [arrayElement(legs, length(legs)).4]), '-') AS route,\n        uniqExact(FlightDate) AS route_recurrence_count\n    FROM itineraries\n    GROUP BY route\n), max_hops AS (\n    SELECT max(hop_count) AS max_hop_count FROM itineraries\n), top_unique AS (\n    SELECT\n        i.FlightDate,\n        i.aircraft_id,\n        i.flight_number,\n        i.carrier,\n        i.hop_count,\n        i.first_dep_sort,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') AS route,\n        rr.route_recurrence_count,\n        row_number() OVER (\n            PARTITION BY arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-')\n            ORDER BY i.FlightDate DESC, i.first_dep_sort DESC, i.carrier DESC, i.flight_number DESC, i.aircraft_id DESC\n        ) AS route_rank\n    FROM itineraries i\n    CROSS JOIN max_hops m\n    INNER JOIN route_recurrence rr\n        ON arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') = rr.route\n    WHERE i.hop_count = m.max_hop_count\n)\nSELECT\n    if(aircraft_id = '', 'unknown', aircraft_id) AS aircraft_id,\n    flight_number,\n    carrier,\n    FlightDate AS flight_date,\n    hop_count,\n    route_recurrence_count,\n    route\nFROM top_unique\nWHERE route_rank = 1\nORDER BY flight_date DESC, first_dep_sort DESC, carrier DESC, flight_number DESC, aircraft_id DESC\nLIMIT 10",
      "row_count": 10,
      "result_columns": [
        "aircraft_id",
        "flight_number",
        "carrier",
        "flight_date",
        "hop_count",
        "route_recurrence_count",
        "route"
      ],
      "first_row": {
        "aircraft_id": "N957WN",
        "carrier": "WN",
        "flight_date": "2024-12-01T00:00:00Z",
        "flight_number": "366",
        "hop_count": 8,
        "route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "route_recurrence_count": 1
      }
    },
    {
      "id": "q1",
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?\n\nList the recurrence count for each of the 10 routes explicitly before summarizing any tiers or categories, and make sure any category totals add up to 10.",
      "answer_markdown": "Recurrence counts for the 10 routes are: ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA = 1; CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN = 4; ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN = 2; MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX = 5; MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC = 40; LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN = 46; SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS = 7; HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK = 20; BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK = 12; BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX = 5. Overall, these top itineraries skew toward recurring patterns rather than pure one-offs: 4 are recurring (10+ days), 5 are occasional (2-9 days), and 1 is a one-off, which sums to all 10 routes.",
      "sql": "WITH main_top10 AS (\n    WITH dedup_legs AS (\n        SELECT\n            FlightDate,\n            ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,\n            Flight_Number_Reporting_Airline AS flight_number,\n            IATA_CODE_Reporting_Airline AS carrier,\n            OriginCode,\n            DestCode,\n            coalesce(CRSDepTime, DepTime, toUInt16(9999)) AS dep_sort,\n            min(coalesce(CRSArrTime, ArrTime, toUInt16(9999))) AS arr_sort\n        FROM ontime.fact_ontime\n        WHERE Cancelled = 0\n        GROUP BY FlightDate, aircraft_id, flight_number, carrier, OriginCode, DestCode, dep_sort\n    ), itineraries AS (\n        SELECT\n            FlightDate,\n            aircraft_id,\n            flight_number,\n            carrier,\n            count() AS hop_count,\n            min(dep_sort) AS first_dep_sort,\n            arraySort(groupArray((dep_sort, arr_sort, OriginCode, DestCode))) AS legs\n        FROM dedup_legs\n        GROUP BY FlightDate, aircraft_id, flight_number, carrier\n        HAVING count() \u003e 1\n    ), route_recurrence AS (\n        SELECT\n            arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, legs), [arrayElement(legs, length(legs)).4]), '-') AS route,\n            uniqExact(FlightDate) AS route_recurrence_count\n        FROM itineraries\n        GROUP BY route\n    ), max_hops AS (\n        SELECT max(hop_count) AS max_hop_count FROM itineraries\n    ), top_unique AS (\n        SELECT\n            i.FlightDate,\n            i.first_dep_sort,\n            arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') AS route,\n            rr.route_recurrence_count,\n            row_number() OVER (\n                PARTITION BY arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-')\n                ORDER BY i.FlightDate DESC, i.first_dep_sort DESC, i.carrier DESC, i.flight_number DESC, i.aircraft_id DESC\n            ) AS route_rank\n        FROM itineraries i\n        CROSS JOIN max_hops m\n        INNER JOIN route_recurrence rr\n            ON arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') = rr.route\n        WHERE i.hop_count = m.max_hop_count\n    )\n    SELECT route, route_recurrence_count, FlightDate, first_dep_sort\n    FROM top_unique\n    WHERE route_rank = 1\n    ORDER BY FlightDate DESC, first_dep_sort DESC\n    LIMIT 10\n)\nSELECT\n    route,\n    route_recurrence_count,\n    multiIf(route_recurrence_count = 1, 'one-off', route_recurrence_count \u003c 10, 'occasional', 'recurring') AS recurrence_tier,\n    count() OVER (PARTITION BY multiIf(route_recurrence_count = 1, 'one-off', route_recurrence_count \u003c 10, 'occasional', 'recurring')) AS tier_count\nFROM main_top10\nORDER BY FlightDate DESC, first_dep_sort DESC",
      "row_count": 10,
      "result_columns": [
        "route",
        "route_recurrence_count",
        "recurrence_tier",
        "tier_count"
      ],
      "first_row": {
        "recurrence_tier": "one-off",
        "route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "route_recurrence_count": 1,
        "tier_count": 1
      }
    },
    {
      "id": "q2",
      "subquestion": "What geographic pattern do the top itineraries show?\n\nBase the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.",
      "answer_markdown": "The geographic pattern is entirely domestic U.S. and strongly tied to Southwest/West Coast plus Sun Belt and Mid-Atlantic nodes. Every airport in the 10 routes is a US airport, and the most repeated points are BWI, DAL, DEN, LAS, MSY, and OAK with 5 mentions each, followed by BNA and PHX with 4 each. California is especially prominent through OAK, BUR, LAX, SAN, PSP, SJC, and SMF, so these longest same-flight-number days look like long domestic zig-zag chains linking eastern and southeastern stops into western endpoints and western connection banks.",
      "sql": "WITH dedup_legs AS (\n    SELECT\n        FlightDate,\n        ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,\n        Flight_Number_Reporting_Airline AS flight_number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginCode,\n        DestCode,\n        coalesce(CRSDepTime, DepTime, toUInt16(9999)) AS dep_sort,\n        min(coalesce(CRSArrTime, ArrTime, toUInt16(9999))) AS arr_sort\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier, OriginCode, DestCode, dep_sort\n), itineraries AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        count() AS hop_count,\n        min(dep_sort) AS first_dep_sort,\n        arraySort(groupArray((dep_sort, arr_sort, OriginCode, DestCode))) AS legs\n    FROM dedup_legs\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier\n    HAVING count() \u003e 1\n), max_hops AS (\n    SELECT max(hop_count) AS max_hop_count FROM itineraries\n), top_unique AS (\n    SELECT\n        i.FlightDate,\n        i.first_dep_sort,\n        i.legs,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') AS route,\n        row_number() OVER (\n            PARTITION BY arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-')\n            ORDER BY i.FlightDate DESC, i.first_dep_sort DESC, i.carrier DESC, i.flight_number DESC, i.aircraft_id DESC\n        ) AS route_rank\n    FROM itineraries i\n    CROSS JOIN max_hops m\n    WHERE i.hop_count = m.max_hop_count\n), top10 AS (\n    SELECT *\n    FROM top_unique\n    WHERE route_rank = 1\n    ORDER BY FlightDate DESC, first_dep_sort DESC\n    LIMIT 10\n), airport_mentions AS (\n    SELECT arrayJoin(arrayConcat(arrayMap(x -\u003e x.3, legs), [arrayElement(legs, length(legs)).4])) AS airport_code\n    FROM top10\n)\nSELECT\n    a.AirportCode AS airport_code,\n    a.StateCode AS state_code,\n    a.CountryCodeISO AS country_code,\n    multiIf(a.Longitude \u003c= -115, 'far_west', a.Longitude \u003c= -95, 'interior_west_central', 'east') AS geo_bucket,\n    count() AS airport_mentions\nFROM airport_mentions am\nLEFT JOIN ontime.dim_airports a\n    ON am.airport_code = toString(a.AirportCode)\n   AND a.IsLatest = 1\nGROUP BY airport_code, state_code, country_code, geo_bucket\nORDER BY airport_mentions DESC, airport_code\nLIMIT 100",
      "row_count": 45,
      "result_columns": [
        "airport_code",
        "state_code",
        "country_code",
        "geo_bucket",
        "airport_mentions"
      ],
      "first_row": {
        "airport_code": "BWI",
        "airport_mentions": 5,
        "country_code": "US",
        "geo_bucket": "east",
        "state_code": "MD"
      }
    }
  ]
}

### Multi-query additions

- The saved SQL shown below is the primary section query for this page.
- The verified analysis package includes named supporting section queries that may be used for enrichment, drill-down, or secondary visuals when the question-specific prompt calls for them.
- Use section answers as narrative framing, but derive displayed KPIs, charts, tables, and interactions from live browser execution of the primary saved SQL and any supporting queries you actually run.
- Do not assume auxiliary lookup or enrichment schema details from memory. Use only columns you have checked against the live endpoint or semantic-layer guidance.
- If you run supporting queries, record them in the same visible query ledger as the primary query.
- The dashboard does not need to mirror `report.md`; it should combine narrative and interactive analysis.

### Verified Analysis Package

Use this JSON package as the supporting context for the visual:

{
  "question_title": "Highest daily hops for one aircraft on one flight number",
  "result_columns": null,
  "row_count": 3,
  "mode_hint": "This visual pass receives only verified section answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "main",
      "subquestion": "Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.\nDefine uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.\nDo not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.\nCount hops from distinct same-day legs, not raw source rows; do not let duplicate or conflicting same-time rows inflate hop count or create artifact routes.\n\nReturn:\n\n- aircraft id\n- flight number\n- carrier\n- flight date\n- hop count\n- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft\n- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`",
      "answer_markdown": "The maximum observed same-aircraft, same-flight-number itinerary length is 8 hops. The 10 most recent unique max-hop routes are all Southwest (WN) itineraries, with the newest on 2024-12-01: ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA flown by N957WN on flight 366. Across these 10 rows, route recurrence ranges from 1 day to 46 days, led by LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN (46), MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC (40), HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK (20), and BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK (12).",
      "sql": "WITH dedup_legs AS (\n    SELECT\n        FlightDate,\n        ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,\n        Flight_Number_Reporting_Airline AS flight_number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginCode,\n        DestCode,\n        coalesce(CRSDepTime, DepTime, toUInt16(9999)) AS dep_sort,\n        min(coalesce(CRSArrTime, ArrTime, toUInt16(9999))) AS arr_sort\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier, OriginCode, DestCode, dep_sort\n), itineraries AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        count() AS hop_count,\n        min(dep_sort) AS first_dep_sort,\n        arraySort(groupArray((dep_sort, arr_sort, OriginCode, DestCode))) AS legs\n    FROM dedup_legs\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier\n    HAVING count() \u003e 1\n), route_recurrence AS (\n    SELECT\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, legs), [arrayElement(legs, length(legs)).4]), '-') AS route,\n        uniqExact(FlightDate) AS route_recurrence_count\n    FROM itineraries\n    GROUP BY route\n), max_hops AS (\n    SELECT max(hop_count) AS max_hop_count FROM itineraries\n), top_unique AS (\n    SELECT\n        i.FlightDate,\n        i.aircraft_id,\n        i.flight_number,\n        i.carrier,\n        i.hop_count,\n        i.first_dep_sort,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') AS route,\n        rr.route_recurrence_count,\n        row_number() OVER (\n            PARTITION BY arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-')\n            ORDER BY i.FlightDate DESC, i.first_dep_sort DESC, i.carrier DESC, i.flight_number DESC, i.aircraft_id DESC\n        ) AS route_rank\n    FROM itineraries i\n    CROSS JOIN max_hops m\n    INNER JOIN route_recurrence rr\n        ON arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') = rr.route\n    WHERE i.hop_count = m.max_hop_count\n)\nSELECT\n    if(aircraft_id = '', 'unknown', aircraft_id) AS aircraft_id,\n    flight_number,\n    carrier,\n    FlightDate AS flight_date,\n    hop_count,\n    route_recurrence_count,\n    route\nFROM top_unique\nWHERE route_rank = 1\nORDER BY flight_date DESC, first_dep_sort DESC, carrier DESC, flight_number DESC, aircraft_id DESC\nLIMIT 10",
      "row_count": 10,
      "result_columns": [
        "aircraft_id",
        "flight_number",
        "carrier",
        "flight_date",
        "hop_count",
        "route_recurrence_count",
        "route"
      ],
      "first_row": {
        "aircraft_id": "N957WN",
        "carrier": "WN",
        "flight_date": "2024-12-01T00:00:00Z",
        "flight_number": "366",
        "hop_count": 8,
        "route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "route_recurrence_count": 1
      }
    },
    {
      "id": "q1",
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?\n\nList the recurrence count for each of the 10 routes explicitly before summarizing any tiers or categories, and make sure any category totals add up to 10.",
      "answer_markdown": "Recurrence counts for the 10 routes are: ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA = 1; CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN = 4; ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN = 2; MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX = 5; MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC = 40; LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN = 46; SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS = 7; HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK = 20; BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK = 12; BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX = 5. Overall, these top itineraries skew toward recurring patterns rather than pure one-offs: 4 are recurring (10+ days), 5 are occasional (2-9 days), and 1 is a one-off, which sums to all 10 routes.",
      "sql": "WITH main_top10 AS (\n    WITH dedup_legs AS (\n        SELECT\n            FlightDate,\n            ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,\n            Flight_Number_Reporting_Airline AS flight_number,\n            IATA_CODE_Reporting_Airline AS carrier,\n            OriginCode,\n            DestCode,\n            coalesce(CRSDepTime, DepTime, toUInt16(9999)) AS dep_sort,\n            min(coalesce(CRSArrTime, ArrTime, toUInt16(9999))) AS arr_sort\n        FROM ontime.fact_ontime\n        WHERE Cancelled = 0\n        GROUP BY FlightDate, aircraft_id, flight_number, carrier, OriginCode, DestCode, dep_sort\n    ), itineraries AS (\n        SELECT\n            FlightDate,\n            aircraft_id,\n            flight_number,\n            carrier,\n            count() AS hop_count,\n            min(dep_sort) AS first_dep_sort,\n            arraySort(groupArray((dep_sort, arr_sort, OriginCode, DestCode))) AS legs\n        FROM dedup_legs\n        GROUP BY FlightDate, aircraft_id, flight_number, carrier\n        HAVING count() \u003e 1\n    ), route_recurrence AS (\n        SELECT\n            arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, legs), [arrayElement(legs, length(legs)).4]), '-') AS route,\n            uniqExact(FlightDate) AS route_recurrence_count\n        FROM itineraries\n        GROUP BY route\n    ), max_hops AS (\n        SELECT max(hop_count) AS max_hop_count FROM itineraries\n    ), top_unique AS (\n        SELECT\n            i.FlightDate,\n            i.first_dep_sort,\n            arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') AS route,\n            rr.route_recurrence_count,\n            row_number() OVER (\n                PARTITION BY arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-')\n                ORDER BY i.FlightDate DESC, i.first_dep_sort DESC, i.carrier DESC, i.flight_number DESC, i.aircraft_id DESC\n            ) AS route_rank\n        FROM itineraries i\n        CROSS JOIN max_hops m\n        INNER JOIN route_recurrence rr\n            ON arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') = rr.route\n        WHERE i.hop_count = m.max_hop_count\n    )\n    SELECT route, route_recurrence_count, FlightDate, first_dep_sort\n    FROM top_unique\n    WHERE route_rank = 1\n    ORDER BY FlightDate DESC, first_dep_sort DESC\n    LIMIT 10\n)\nSELECT\n    route,\n    route_recurrence_count,\n    multiIf(route_recurrence_count = 1, 'one-off', route_recurrence_count \u003c 10, 'occasional', 'recurring') AS recurrence_tier,\n    count() OVER (PARTITION BY multiIf(route_recurrence_count = 1, 'one-off', route_recurrence_count \u003c 10, 'occasional', 'recurring')) AS tier_count\nFROM main_top10\nORDER BY FlightDate DESC, first_dep_sort DESC",
      "row_count": 10,
      "result_columns": [
        "route",
        "route_recurrence_count",
        "recurrence_tier",
        "tier_count"
      ],
      "first_row": {
        "recurrence_tier": "one-off",
        "route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "route_recurrence_count": 1,
        "tier_count": 1
      }
    },
    {
      "id": "q2",
      "subquestion": "What geographic pattern do the top itineraries show?\n\nBase the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.",
      "answer_markdown": "The geographic pattern is entirely domestic U.S. and strongly tied to Southwest/West Coast plus Sun Belt and Mid-Atlantic nodes. Every airport in the 10 routes is a US airport, and the most repeated points are BWI, DAL, DEN, LAS, MSY, and OAK with 5 mentions each, followed by BNA and PHX with 4 each. California is especially prominent through OAK, BUR, LAX, SAN, PSP, SJC, and SMF, so these longest same-flight-number days look like long domestic zig-zag chains linking eastern and southeastern stops into western endpoints and western connection banks.",
      "sql": "WITH dedup_legs AS (\n    SELECT\n        FlightDate,\n        ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,\n        Flight_Number_Reporting_Airline AS flight_number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginCode,\n        DestCode,\n        coalesce(CRSDepTime, DepTime, toUInt16(9999)) AS dep_sort,\n        min(coalesce(CRSArrTime, ArrTime, toUInt16(9999))) AS arr_sort\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier, OriginCode, DestCode, dep_sort\n), itineraries AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        count() AS hop_count,\n        min(dep_sort) AS first_dep_sort,\n        arraySort(groupArray((dep_sort, arr_sort, OriginCode, DestCode))) AS legs\n    FROM dedup_legs\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier\n    HAVING count() \u003e 1\n), max_hops AS (\n    SELECT max(hop_count) AS max_hop_count FROM itineraries\n), top_unique AS (\n    SELECT\n        i.FlightDate,\n        i.first_dep_sort,\n        i.legs,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') AS route,\n        row_number() OVER (\n            PARTITION BY arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-')\n            ORDER BY i.FlightDate DESC, i.first_dep_sort DESC, i.carrier DESC, i.flight_number DESC, i.aircraft_id DESC\n        ) AS route_rank\n    FROM itineraries i\n    CROSS JOIN max_hops m\n    WHERE i.hop_count = m.max_hop_count\n), top10 AS (\n    SELECT *\n    FROM top_unique\n    WHERE route_rank = 1\n    ORDER BY FlightDate DESC, first_dep_sort DESC\n    LIMIT 10\n), airport_mentions AS (\n    SELECT arrayJoin(arrayConcat(arrayMap(x -\u003e x.3, legs), [arrayElement(legs, length(legs)).4])) AS airport_code\n    FROM top10\n)\nSELECT\n    a.AirportCode AS airport_code,\n    a.StateCode AS state_code,\n    a.CountryCodeISO AS country_code,\n    multiIf(a.Longitude \u003c= -115, 'far_west', a.Longitude \u003c= -95, 'interior_west_central', 'east') AS geo_bucket,\n    count() AS airport_mentions\nFROM airport_mentions am\nLEFT JOIN ontime.dim_airports a\n    ON am.airport_code = toString(a.AirportCode)\n   AND a.IsLatest = 1\nGROUP BY airport_code, state_code, country_code, geo_bucket\nORDER BY airport_mentions DESC, airport_code\nLIMIT 100",
      "row_count": 45,
      "result_columns": [
        "airport_code",
        "state_code",
        "country_code",
        "geo_bucket",
        "airport_mentions"
      ],
      "first_row": {
        "airport_code": "BWI",
        "airport_mentions": 5,
        "country_code": "US",
        "geo_bucket": "east",
        "state_code": "MD"
      }
    }
  ]
}

### Dynamic-mode additions

- Use this endpoint template for every browser query: `https://mcp.demo.altinity.cloud/{JWE}/openapi/execute_query?query=...`
- Keep JWE in `localStorage['OnTimeAnalystDashboard::auth::jwe']`.
- Do not embed the primary analytical dataset as `result.json` payloads or CSV snapshots.
- Before writing `visual.html`, self-verify every browser-side SQL statement you intend to ship, including primary, supporting, enrichment, drill-down, and lookup queries.
- For each query, run a cheap live check against the real endpoint and schema first, usually with a small `LIMIT`, a narrow `WHERE` filter, or both when that preserves the query shape.
- Treat successful execution as mandatory. Fix any syntax, type, aggregate, join, or unknown-column errors in a loop until every shipped browser query runs successfully.

Create browser-ready HTML `visual.html`.

Write the file or provide a download link. Do not include the HTML source in the response. Do not open the artifact view frame.