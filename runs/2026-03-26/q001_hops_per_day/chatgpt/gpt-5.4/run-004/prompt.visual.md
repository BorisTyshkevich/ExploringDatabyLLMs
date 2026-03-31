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

### main

- use main proof query as the primary saved SQL already provided in the prompt
- use the other reviewed section proof queries as supporting queries when they materially improve the narrative or supporting panels
- anchor the hero narrative to the top-ranked itinerary even when another itinerary is selected in the table
- show a lead-itinerary map that remains present even before airport-coordinate lookup succeeds
- treat the first row returned by the primary query as the default selected itinerary on initial load
- derive hop count, stop sequence, and repeated-route comparisons from the result set
- include a narrative hero about the lead itinerary and the broader geographic pattern of the top itineraries
- label the map as airport-coordinate lookup in the query ledger
- reuse the lookup results for any itinerary selected from the primary result set without issuing a new per-click lookup query
- include KPI cards for tail number, flight number, date, hop count, and route repetition context, with the date shown as its own visible KPI value
- keep the KPI strip synced to the currently selected itinerary
- include a legend plus both a route sequence/detail panel and an itinerary table below the map
- make itinerary table rows clickable so selecting a row redraws the map and refreshes the route sequence/detail panel for that itinerary
- make the selected-row map behavior explicit: when the selected itinerary differs from Rank 1, the map title, plotted route, markers, bounds, and route detail panel must visibly update to that selected itinerary rather than leaving the lead route drawn
- keep the map/detail/KPI selection state separate from the anchored hero state
- show a clear active-row state for the selected itinerary that is distinct from simple hover styling
- prefer the `Route` value from the primary query as the per-row itinerary representation for redraws
- if airport-coordinate lookup fails or the selected itinerary lacks enough coordinates, keep the map card visible with degraded-state messaging for that selected itinerary, report the degraded map in the ledger, and continue rendering the non-map analysis
- derive the ordered itinerary sequence for map redraws and the route detail panel by splitting `Route` on `-`

### operational-stress

Which airports act as the key connectors, origins, and termini within the top 10 unique maximum-hop itineraries?
Classify airport appearances by route position and return airport code, airport name, city/state, total appearances, origin appearances, intermediate-stop appearances, final-destination appearances, and share of itineraries containing that airport.

- author a sql query 
- include a pane below the selection-driven itinerary detail 
- have the operational-stress lookup query use the currently selected itinerary context from the primary result set and fetch per-airport and per-leg average departure delay, average arrival delay, 15-plus-minute delay rate, diversion incidence, and stop-position stress signals
- label that pane query as an operational-stress lookup in the query ledger
- make itinerary table row selection also refresh the operational-stress lookup pane for that same itinerary
- if the operational-stress lookup fails, keep the operational-stress pane visible with degraded-state messaging for the selected itinerary, report that degraded pane in the ledger, and continue rendering the rest of the dashboard

### Data Source

SQL query for primary data source:

```sql
WITH params AS (
    SELECT addYears(toDate((SELECT max(FlightDate) FROM ontime.fact_ontime)), -5) AS start_date,
           (SELECT max(FlightDate) FROM ontime.fact_ontime) AS end_date
),
legs_raw AS (
    SELECT
        FlightDate,
        ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,
        Flight_Number_Reporting_Airline AS flight_number,
        IATA_CODE_Reporting_Airline AS carrier,
        OriginAirportID,
        DestAirportID,
        OriginCode,
        DestCode,
        coalesce(CRSDepTime, DepTime, 0) AS dep_hhmm,
        coalesce(CRSArrTime, ArrTime, 0) AS arr_hhmm,
        coalesce(DepTime, CRSDepTime, 0) AS dep_sort_hhmm,
        Distance,
        row_number() OVER (
            PARTITION BY FlightDate, ifNull(nullIf(Tail_Number, ''), ''), Flight_Number_Reporting_Airline,
                         IATA_CODE_Reporting_Airline, OriginAirportID, DestAirportID,
                         coalesce(CRSDepTime, DepTime, 0), coalesce(CRSArrTime, ArrTime, 0)
            ORDER BY Cancelled ASC, Diverted ASC, coalesce(ActualElapsedTime, CRSElapsedTime, 1000000) ASC
        ) AS rn
    FROM ontime.fact_ontime
    WHERE FlightDate >= (SELECT start_date FROM params)
      AND FlightDate <= (SELECT end_date FROM params)
      AND Cancelled = 0
),
legs AS (
    SELECT *,
           toDateTime(FlightDate) + toIntervalMinute(intDiv(dep_sort_hhmm,100)*60 + (dep_sort_hhmm % 100)) AS dep_ts
    FROM legs_raw
    WHERE rn = 1
),
itineraries AS (
    SELECT
        FlightDate,
        aircraft_id,
        flight_number,
        carrier,
        arraySort(groupArray((dep_ts, OriginCode, DestCode, OriginAirportID, DestAirportID, toUInt32(ifNull(Distance,0))))) AS legs_sorted
    FROM legs
    GROUP BY FlightDate, aircraft_id, flight_number, carrier
    HAVING length(legs_sorted) > 1
       AND arrayAll(i -> legs_sorted[i].5 = legs_sorted[i + 1].4, range(1, length(legs_sorted)))
),
routes AS (
    SELECT
        FlightDate,
        aircraft_id,
        flight_number,
        carrier,
        length(legs_sorted) AS hop_count,
        legs_sorted[1].1 AS first_dep_ts,
        arrayStringConcat(arrayConcat(arrayMap(x -> x.2, legs_sorted), [legs_sorted[length(legs_sorted)].3]), '-') AS Route
    FROM itineraries
),
route_days AS (
    SELECT Route, uniqExact(FlightDate) AS route_recurrence_count
    FROM routes
    GROUP BY Route
),
unique_ranked AS (
    SELECT r.*, rd.route_recurrence_count,
           row_number() OVER (PARTITION BY Route ORDER BY first_dep_ts DESC, FlightDate DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rn
    FROM routes r
    INNER JOIN route_days rd USING (Route)
)
SELECT aircraft_id, flight_number, carrier, FlightDate, hop_count, route_recurrence_count, Route
FROM unique_ranked
WHERE route_rn = 1
ORDER BY hop_count DESC, first_dep_ts DESC
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
      "subquestion": "Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.\nDefine uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.\nDo not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.\nCount hops from distinct same-day legs, not raw source rows; do not let duplicate or conflicting same-time rows inflate hop count or create artifact routes.\n\nReturn:\n\n- aircraft id\n- flight number\n- carrier\n- flight date\n- hop count\n- route recurrence count: total number of days across the analyzed window on which this exact Route string was flown by any aircraft\n- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`",
      "answer_markdown": "Across the most recent five years ending 2025-11-30, the longest same-aircraft, same-flight-number itineraries all reached 8 hops. The most recent top 10 unique routes are all Southwest (WN) examples, led by ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA on 2024-12-01; recurrence across the window ranges from 1 day to 46 days, with LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN the most recurrent of the top 10.",
      "sql": "WITH params AS (\n    SELECT addYears(toDate((SELECT max(FlightDate) FROM ontime.fact_ontime)), -5) AS start_date,\n           (SELECT max(FlightDate) FROM ontime.fact_ontime) AS end_date\n),\nlegs_raw AS (\n    SELECT\n        FlightDate,\n        ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,\n        Flight_Number_Reporting_Airline AS flight_number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginAirportID,\n        DestAirportID,\n        OriginCode,\n        DestCode,\n        coalesce(CRSDepTime, DepTime, 0) AS dep_hhmm,\n        coalesce(CRSArrTime, ArrTime, 0) AS arr_hhmm,\n        coalesce(DepTime, CRSDepTime, 0) AS dep_sort_hhmm,\n        Distance,\n        row_number() OVER (\n            PARTITION BY FlightDate, ifNull(nullIf(Tail_Number, ''), ''), Flight_Number_Reporting_Airline,\n                         IATA_CODE_Reporting_Airline, OriginAirportID, DestAirportID,\n                         coalesce(CRSDepTime, DepTime, 0), coalesce(CRSArrTime, ArrTime, 0)\n            ORDER BY Cancelled ASC, Diverted ASC, coalesce(ActualElapsedTime, CRSElapsedTime, 1000000) ASC\n        ) AS rn\n    FROM ontime.fact_ontime\n    WHERE FlightDate \u003e= (SELECT start_date FROM params)\n      AND FlightDate \u003c= (SELECT end_date FROM params)\n      AND Cancelled = 0\n),\nlegs AS (\n    SELECT *,\n           toDateTime(FlightDate) + toIntervalMinute(intDiv(dep_sort_hhmm,100)*60 + (dep_sort_hhmm % 100)) AS dep_ts\n    FROM legs_raw\n    WHERE rn = 1\n),\nitineraries AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        arraySort(groupArray((dep_ts, OriginCode, DestCode, OriginAirportID, DestAirportID, toUInt32(ifNull(Distance,0))))) AS legs_sorted\n    FROM legs\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier\n    HAVING length(legs_sorted) \u003e 1\n       AND arrayAll(i -\u003e legs_sorted[i].5 = legs_sorted[i + 1].4, range(1, length(legs_sorted)))\n),\nroutes AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        length(legs_sorted) AS hop_count,\n        legs_sorted[1].1 AS first_dep_ts,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.2, legs_sorted), [legs_sorted[length(legs_sorted)].3]), '-') AS Route\n    FROM itineraries\n),\nroute_days AS (\n    SELECT Route, uniqExact(FlightDate) AS route_recurrence_count\n    FROM routes\n    GROUP BY Route\n),\nunique_ranked AS (\n    SELECT r.*, rd.route_recurrence_count,\n           row_number() OVER (PARTITION BY Route ORDER BY first_dep_ts DESC, FlightDate DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rn\n    FROM routes r\n    INNER JOIN route_days rd USING (Route)\n)\nSELECT aircraft_id, flight_number, carrier, FlightDate, hop_count, route_recurrence_count, Route\nFROM unique_ranked\nWHERE route_rn = 1\nORDER BY hop_count DESC, first_dep_ts DESC\nLIMIT 10",
      "is_primary": true,
      "date_field_hint": "FlightDate",
      "row_count": 10,
      "result_columns": [
        "aircraft_id",
        "flight_number",
        "carrier",
        "FlightDate",
        "hop_count",
        "route_recurrence_count",
        "Route"
      ],
      "first_row": {
        "FlightDate": "2024-12-01T00:00:00Z",
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "aircraft_id": "N957WN",
        "carrier": "WN",
        "flight_number": "366",
        "hop_count": 8,
        "route_recurrence_count": 1
      }
    },
    {
      "id": "q1",
      "subquestion": "Which airports act as the key connectors, origins, and termini within the top 10 unique maximum-hop itineraries?\nClassify airport appearances by route position and return airport code, airport name, city/state, total appearances, origin appearances, intermediate-stop appearances, final-destination appearances, and share of itineraries containing that airport.",
      "answer_markdown": "Dallas Love Field, Denver, Las Vegas, Baltimore/Washington, New Orleans, and Oakland are the key nodes, each appearing in 5 of the 10 itineraries. DAL is the strongest pure connector with 5 intermediate-stop appearances; BWI and MSY are the most common origins at 2 each; OAK and LAX are the most common final destinations at 2 each.",
      "sql": "WITH params AS (\n    SELECT addYears(toDate((SELECT max(FlightDate) FROM ontime.fact_ontime)), -5) AS start_date,\n           (SELECT max(FlightDate) FROM ontime.fact_ontime) AS end_date\n),\nlegs_raw AS (\n    SELECT\n        FlightDate,\n        ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,\n        Flight_Number_Reporting_Airline AS flight_number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginAirportID,\n        DestAirportID,\n        OriginCode,\n        DestCode,\n        coalesce(CRSDepTime, DepTime, 0) AS dep_hhmm,\n        coalesce(CRSArrTime, ArrTime, 0) AS arr_hhmm,\n        coalesce(DepTime, CRSDepTime, 0) AS dep_sort_hhmm,\n        Distance,\n        row_number() OVER (\n            PARTITION BY FlightDate, ifNull(nullIf(Tail_Number, ''), ''), Flight_Number_Reporting_Airline,\n                         IATA_CODE_Reporting_Airline, OriginAirportID, DestAirportID,\n                         coalesce(CRSDepTime, DepTime, 0), coalesce(CRSArrTime, ArrTime, 0)\n            ORDER BY Cancelled ASC, Diverted ASC, coalesce(ActualElapsedTime, CRSElapsedTime, 1000000) ASC\n        ) AS rn\n    FROM ontime.fact_ontime\n    WHERE FlightDate \u003e= (SELECT start_date FROM params)\n      AND FlightDate \u003c= (SELECT end_date FROM params)\n      AND Cancelled = 0\n),\nlegs AS (\n    SELECT *,\n           toDateTime(FlightDate) + toIntervalMinute(intDiv(dep_sort_hhmm,100)*60 + (dep_sort_hhmm % 100)) AS dep_ts\n    FROM legs_raw\n    WHERE rn = 1\n),\nitineraries AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        arraySort(groupArray((dep_ts, OriginCode, DestCode, OriginAirportID, DestAirportID, toUInt32(ifNull(Distance,0))))) AS legs_sorted\n    FROM legs\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier\n    HAVING length(legs_sorted) \u003e 1\n       AND arrayAll(i -\u003e legs_sorted[i].5 = legs_sorted[i + 1].4, range(1, length(legs_sorted)))\n),\nroutes AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        length(legs_sorted) AS hop_count,\n        legs_sorted[1].1 AS first_dep_ts,\n        arrayMap(x -\u003e x.2, legs_sorted) AS origin_codes,\n        arrayMap(x -\u003e x.4, legs_sorted) AS origin_airport_ids,\n        arrayMap(x -\u003e x.5, legs_sorted) AS dest_airport_ids,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.2, legs_sorted), [legs_sorted[length(legs_sorted)].3]), '-') AS Route\n    FROM itineraries\n),\nroute_days AS (\n    SELECT Route, uniqExact(FlightDate) AS route_recurrence_count\n    FROM routes\n    GROUP BY Route\n),\nunique_ranked AS (\n    SELECT r.*, rd.route_recurrence_count,\n           row_number() OVER (PARTITION BY Route ORDER BY first_dep_ts DESC, FlightDate DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rn\n    FROM routes r\n    INNER JOIN route_days rd USING (Route)\n),\ntop_routes AS (\n    SELECT *\n    FROM unique_ranked\n    WHERE route_rn = 1\n    ORDER BY hop_count DESC, first_dep_ts DESC\n    LIMIT 10\n),\nroute_airports AS (\n    SELECT\n        Route,\n        arrayJoin(arrayZip(\n            arrayConcat(origin_codes, [arrayElement(splitByChar('-', Route), length(splitByChar('-', Route)))]),\n            arrayConcat(origin_airport_ids, [dest_airport_ids[length(dest_airport_ids)]]),\n            arrayMap(i -\u003e if(i = 1, 'origin', if(i = length(origin_codes) + 1, 'final_destination', 'intermediate_stop')), arrayEnumerate(arrayConcat(origin_codes, [arrayElement(splitByChar('-', Route), length(splitByChar('-', Route)))])))\n        )) AS ap\n    FROM top_routes\n)\nSELECT\n    ap.1 AS airport_code,\n    any(a.DisplayAirportName) AS airport_name,\n    any(a.CityName) AS city_state,\n    count() AS total_appearances,\n    countIf(ap.3 = 'origin') AS origin_appearances,\n    countIf(ap.3 = 'intermediate_stop') AS intermediate_stop_appearances,\n    countIf(ap.3 = 'final_destination') AS final_destination_appearances,\n    round(uniqExact(Route) / 10.0, 3) AS share_of_itineraries_containing_airport\nFROM route_airports\nLEFT JOIN ontime.dim_airports a\n    ON ap.2 = a.AirportID AND a.IsLatest = 1\nGROUP BY airport_code\nORDER BY total_appearances DESC, intermediate_stop_appearances DESC, airport_code",
      "date_field_hint": "airport_name",
      "row_count": 45,
      "result_columns": [
        "airport_code",
        "airport_name",
        "city_state",
        "total_appearances",
        "origin_appearances",
        "intermediate_stop_appearances",
        "final_destination_appearances",
        "share_of_itineraries_containing_airport"
      ],
      "first_row": {
        "airport_code": "DAL",
        "airport_name": "Dallas Love Field",
        "city_state": "Dallas, TX",
        "final_destination_appearances": 0,
        "intermediate_stop_appearances": 5,
        "origin_appearances": 0,
        "share_of_itineraries_containing_airport": 0.5,
        "total_appearances": 5
      }
    },
    {
      "id": "q2",
      "subquestion": "How geographically extreme is each of the top 10 unique maximum-hop itineraries?\nReturn total flown distance, unique airports, unique city markets, unique states, unique local-time offsets, and whether the route is entirely domestic, then summarize which routes are the most geographically expansive.",
      "answer_markdown": "All 10 top routes are entirely domestic. The most geographically expansive by flown distance is BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX at 5,169 miles; the broadest state coverage is 9 states on BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX, MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX, and ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA; five routes span 4 local-time offsets, while every route touches 9 unique airports.",
      "sql": "WITH params AS (\n    SELECT addYears(toDate((SELECT max(FlightDate) FROM ontime.fact_ontime)), -5) AS start_date,\n           (SELECT max(FlightDate) FROM ontime.fact_ontime) AS end_date\n),\nlegs_raw AS (\n    SELECT\n        FlightDate,\n        ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,\n        Flight_Number_Reporting_Airline AS flight_number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginAirportID,\n        DestAirportID,\n        OriginCode,\n        DestCode,\n        coalesce(CRSDepTime, DepTime, 0) AS dep_hhmm,\n        coalesce(CRSArrTime, ArrTime, 0) AS arr_hhmm,\n        coalesce(DepTime, CRSDepTime, 0) AS dep_sort_hhmm,\n        Distance,\n        row_number() OVER (\n            PARTITION BY FlightDate, ifNull(nullIf(Tail_Number, ''), ''), Flight_Number_Reporting_Airline,\n                         IATA_CODE_Reporting_Airline, OriginAirportID, DestAirportID,\n                         coalesce(CRSDepTime, DepTime, 0), coalesce(CRSArrTime, ArrTime, 0)\n            ORDER BY Cancelled ASC, Diverted ASC, coalesce(ActualElapsedTime, CRSElapsedTime, 1000000) ASC\n        ) AS rn\n    FROM ontime.fact_ontime\n    WHERE FlightDate \u003e= (SELECT start_date FROM params)\n      AND FlightDate \u003c= (SELECT end_date FROM params)\n      AND Cancelled = 0\n),\nlegs AS (\n    SELECT *,\n           toDateTime(FlightDate) + toIntervalMinute(intDiv(dep_sort_hhmm,100)*60 + (dep_sort_hhmm % 100)) AS dep_ts\n    FROM legs_raw\n    WHERE rn = 1\n),\nitineraries AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        arraySort(groupArray((dep_ts, OriginCode, DestCode, OriginAirportID, DestAirportID, toUInt32(ifNull(Distance,0))))) AS legs_sorted\n    FROM legs\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier\n    HAVING length(legs_sorted) \u003e 1\n       AND arrayAll(i -\u003e legs_sorted[i].5 = legs_sorted[i + 1].4, range(1, length(legs_sorted)))\n),\nroutes AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        length(legs_sorted) AS hop_count,\n        legs_sorted[1].1 AS first_dep_ts,\n        arrayMap(x -\u003e x.4, legs_sorted) AS origin_airport_ids,\n        arrayMap(x -\u003e x.5, legs_sorted) AS dest_airport_ids,\n        arrayMap(x -\u003e toUInt32(x.6), legs_sorted) AS leg_distances,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.2, legs_sorted), [legs_sorted[length(legs_sorted)].3]), '-') AS Route\n    FROM itineraries\n),\nroute_days AS (\n    SELECT Route, uniqExact(FlightDate) AS route_recurrence_count\n    FROM routes\n    GROUP BY Route\n),\nunique_ranked AS (\n    SELECT r.*, rd.route_recurrence_count,\n           row_number() OVER (PARTITION BY Route ORDER BY first_dep_ts DESC, FlightDate DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rn\n    FROM routes r\n    INNER JOIN route_days rd USING (Route)\n),\ntop_routes AS (\n    SELECT *,\n           arrayDistinct(arrayConcat(origin_airport_ids, [dest_airport_ids[length(dest_airport_ids)]])) AS route_airport_ids,\n           arraySum(leg_distances) AS total_flown_distance\n    FROM unique_ranked\n    WHERE route_rn = 1\n    ORDER BY hop_count DESC, first_dep_ts DESC\n    LIMIT 10\n),\nroute_airport_dim AS (\n    SELECT\n        tr.Route,\n        tr.total_flown_distance,\n        tr.route_airport_ids,\n        a.CityMarketID,\n        a.StateCode,\n        a.UtcLocalTimeVariation,\n        a.CountryCodeISO\n    FROM top_routes tr\n    ARRAY JOIN tr.route_airport_ids AS airport_id\n    LEFT JOIN ontime.dim_airports a\n        ON airport_id = a.AirportID AND a.IsLatest = 1\n)\nSELECT\n    Route,\n    any(total_flown_distance) AS total_flown_distance,\n    length(any(route_airport_ids)) AS unique_airports,\n    uniqExact(CityMarketID) AS unique_city_markets,\n    uniqExact(StateCode) AS unique_states,\n    uniqExact(UtcLocalTimeVariation) AS unique_local_time_offsets,\n    if(countIf(CountryCodeISO != 'US') = 0, 'Yes', 'No') AS entirely_domestic\nFROM route_airport_dim\nGROUP BY Route\nORDER BY total_flown_distance DESC, unique_states DESC\nLIMIT 10",
      "date_field_hint": "Route",
      "row_count": 10,
      "result_columns": [
        "Route",
        "total_flown_distance",
        "unique_airports",
        "unique_city_markets",
        "unique_states",
        "unique_local_time_offsets",
        "entirely_domestic"
      ],
      "first_row": {
        "Route": "BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX",
        "entirely_domestic": "Yes",
        "total_flown_distance": 5169,
        "unique_airports": 9,
        "unique_city_markets": 8,
        "unique_local_time_offsets": 3,
        "unique_states": 9
      }
    }
  ]
}

### Multi-query additions

- The saved SQL shown below is the primary section query for this page.
- The verified analysis package includes named supporting section queries that may be used for enrichment, drill-down, or secondary visuals when the question-specific prompt calls for them.
- Treat the supporting queries in the verified analysis package as first-class runtime queries, not just narrative context.
- Prefill editable SQL panels for the primary query and each supporting query from the verified package.
- Use section answers as narrative framing, but derive displayed KPIs, charts, tables, and interactions from live browser execution of the primary saved SQL and any supporting queries you actually run.
- When the question-specific visual prompt explicitly asks for a lookup query, author that query only in the visual pass; do not treat it as part of the reviewed analysis package.
- Do not assume auxiliary lookup or enrichment schema details from memory. Use only columns you have checked against the live endpoint or semantic-layer guidance.
- If you run supporting queries, record them in the same visible query ledger as the primary query.
- If you add a lookup query in the visual pass, record it in that same unified query ledger with its own label and status.
- When the date selector changes, rerun the primary query and every supporting query that supports the active date range.
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
      "subquestion": "Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.\nDefine uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.\nDo not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.\nCount hops from distinct same-day legs, not raw source rows; do not let duplicate or conflicting same-time rows inflate hop count or create artifact routes.\n\nReturn:\n\n- aircraft id\n- flight number\n- carrier\n- flight date\n- hop count\n- route recurrence count: total number of days across the analyzed window on which this exact Route string was flown by any aircraft\n- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`",
      "answer_markdown": "Across the most recent five years ending 2025-11-30, the longest same-aircraft, same-flight-number itineraries all reached 8 hops. The most recent top 10 unique routes are all Southwest (WN) examples, led by ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA on 2024-12-01; recurrence across the window ranges from 1 day to 46 days, with LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN the most recurrent of the top 10.",
      "sql": "WITH params AS (\n    SELECT addYears(toDate((SELECT max(FlightDate) FROM ontime.fact_ontime)), -5) AS start_date,\n           (SELECT max(FlightDate) FROM ontime.fact_ontime) AS end_date\n),\nlegs_raw AS (\n    SELECT\n        FlightDate,\n        ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,\n        Flight_Number_Reporting_Airline AS flight_number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginAirportID,\n        DestAirportID,\n        OriginCode,\n        DestCode,\n        coalesce(CRSDepTime, DepTime, 0) AS dep_hhmm,\n        coalesce(CRSArrTime, ArrTime, 0) AS arr_hhmm,\n        coalesce(DepTime, CRSDepTime, 0) AS dep_sort_hhmm,\n        Distance,\n        row_number() OVER (\n            PARTITION BY FlightDate, ifNull(nullIf(Tail_Number, ''), ''), Flight_Number_Reporting_Airline,\n                         IATA_CODE_Reporting_Airline, OriginAirportID, DestAirportID,\n                         coalesce(CRSDepTime, DepTime, 0), coalesce(CRSArrTime, ArrTime, 0)\n            ORDER BY Cancelled ASC, Diverted ASC, coalesce(ActualElapsedTime, CRSElapsedTime, 1000000) ASC\n        ) AS rn\n    FROM ontime.fact_ontime\n    WHERE FlightDate \u003e= (SELECT start_date FROM params)\n      AND FlightDate \u003c= (SELECT end_date FROM params)\n      AND Cancelled = 0\n),\nlegs AS (\n    SELECT *,\n           toDateTime(FlightDate) + toIntervalMinute(intDiv(dep_sort_hhmm,100)*60 + (dep_sort_hhmm % 100)) AS dep_ts\n    FROM legs_raw\n    WHERE rn = 1\n),\nitineraries AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        arraySort(groupArray((dep_ts, OriginCode, DestCode, OriginAirportID, DestAirportID, toUInt32(ifNull(Distance,0))))) AS legs_sorted\n    FROM legs\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier\n    HAVING length(legs_sorted) \u003e 1\n       AND arrayAll(i -\u003e legs_sorted[i].5 = legs_sorted[i + 1].4, range(1, length(legs_sorted)))\n),\nroutes AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        length(legs_sorted) AS hop_count,\n        legs_sorted[1].1 AS first_dep_ts,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.2, legs_sorted), [legs_sorted[length(legs_sorted)].3]), '-') AS Route\n    FROM itineraries\n),\nroute_days AS (\n    SELECT Route, uniqExact(FlightDate) AS route_recurrence_count\n    FROM routes\n    GROUP BY Route\n),\nunique_ranked AS (\n    SELECT r.*, rd.route_recurrence_count,\n           row_number() OVER (PARTITION BY Route ORDER BY first_dep_ts DESC, FlightDate DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rn\n    FROM routes r\n    INNER JOIN route_days rd USING (Route)\n)\nSELECT aircraft_id, flight_number, carrier, FlightDate, hop_count, route_recurrence_count, Route\nFROM unique_ranked\nWHERE route_rn = 1\nORDER BY hop_count DESC, first_dep_ts DESC\nLIMIT 10",
      "is_primary": true,
      "date_field_hint": "FlightDate",
      "row_count": 10,
      "result_columns": [
        "aircraft_id",
        "flight_number",
        "carrier",
        "FlightDate",
        "hop_count",
        "route_recurrence_count",
        "Route"
      ],
      "first_row": {
        "FlightDate": "2024-12-01T00:00:00Z",
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "aircraft_id": "N957WN",
        "carrier": "WN",
        "flight_number": "366",
        "hop_count": 8,
        "route_recurrence_count": 1
      }
    },
    {
      "id": "q1",
      "subquestion": "Which airports act as the key connectors, origins, and termini within the top 10 unique maximum-hop itineraries?\nClassify airport appearances by route position and return airport code, airport name, city/state, total appearances, origin appearances, intermediate-stop appearances, final-destination appearances, and share of itineraries containing that airport.",
      "answer_markdown": "Dallas Love Field, Denver, Las Vegas, Baltimore/Washington, New Orleans, and Oakland are the key nodes, each appearing in 5 of the 10 itineraries. DAL is the strongest pure connector with 5 intermediate-stop appearances; BWI and MSY are the most common origins at 2 each; OAK and LAX are the most common final destinations at 2 each.",
      "sql": "WITH params AS (\n    SELECT addYears(toDate((SELECT max(FlightDate) FROM ontime.fact_ontime)), -5) AS start_date,\n           (SELECT max(FlightDate) FROM ontime.fact_ontime) AS end_date\n),\nlegs_raw AS (\n    SELECT\n        FlightDate,\n        ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,\n        Flight_Number_Reporting_Airline AS flight_number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginAirportID,\n        DestAirportID,\n        OriginCode,\n        DestCode,\n        coalesce(CRSDepTime, DepTime, 0) AS dep_hhmm,\n        coalesce(CRSArrTime, ArrTime, 0) AS arr_hhmm,\n        coalesce(DepTime, CRSDepTime, 0) AS dep_sort_hhmm,\n        Distance,\n        row_number() OVER (\n            PARTITION BY FlightDate, ifNull(nullIf(Tail_Number, ''), ''), Flight_Number_Reporting_Airline,\n                         IATA_CODE_Reporting_Airline, OriginAirportID, DestAirportID,\n                         coalesce(CRSDepTime, DepTime, 0), coalesce(CRSArrTime, ArrTime, 0)\n            ORDER BY Cancelled ASC, Diverted ASC, coalesce(ActualElapsedTime, CRSElapsedTime, 1000000) ASC\n        ) AS rn\n    FROM ontime.fact_ontime\n    WHERE FlightDate \u003e= (SELECT start_date FROM params)\n      AND FlightDate \u003c= (SELECT end_date FROM params)\n      AND Cancelled = 0\n),\nlegs AS (\n    SELECT *,\n           toDateTime(FlightDate) + toIntervalMinute(intDiv(dep_sort_hhmm,100)*60 + (dep_sort_hhmm % 100)) AS dep_ts\n    FROM legs_raw\n    WHERE rn = 1\n),\nitineraries AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        arraySort(groupArray((dep_ts, OriginCode, DestCode, OriginAirportID, DestAirportID, toUInt32(ifNull(Distance,0))))) AS legs_sorted\n    FROM legs\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier\n    HAVING length(legs_sorted) \u003e 1\n       AND arrayAll(i -\u003e legs_sorted[i].5 = legs_sorted[i + 1].4, range(1, length(legs_sorted)))\n),\nroutes AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        length(legs_sorted) AS hop_count,\n        legs_sorted[1].1 AS first_dep_ts,\n        arrayMap(x -\u003e x.2, legs_sorted) AS origin_codes,\n        arrayMap(x -\u003e x.4, legs_sorted) AS origin_airport_ids,\n        arrayMap(x -\u003e x.5, legs_sorted) AS dest_airport_ids,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.2, legs_sorted), [legs_sorted[length(legs_sorted)].3]), '-') AS Route\n    FROM itineraries\n),\nroute_days AS (\n    SELECT Route, uniqExact(FlightDate) AS route_recurrence_count\n    FROM routes\n    GROUP BY Route\n),\nunique_ranked AS (\n    SELECT r.*, rd.route_recurrence_count,\n           row_number() OVER (PARTITION BY Route ORDER BY first_dep_ts DESC, FlightDate DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rn\n    FROM routes r\n    INNER JOIN route_days rd USING (Route)\n),\ntop_routes AS (\n    SELECT *\n    FROM unique_ranked\n    WHERE route_rn = 1\n    ORDER BY hop_count DESC, first_dep_ts DESC\n    LIMIT 10\n),\nroute_airports AS (\n    SELECT\n        Route,\n        arrayJoin(arrayZip(\n            arrayConcat(origin_codes, [arrayElement(splitByChar('-', Route), length(splitByChar('-', Route)))]),\n            arrayConcat(origin_airport_ids, [dest_airport_ids[length(dest_airport_ids)]]),\n            arrayMap(i -\u003e if(i = 1, 'origin', if(i = length(origin_codes) + 1, 'final_destination', 'intermediate_stop')), arrayEnumerate(arrayConcat(origin_codes, [arrayElement(splitByChar('-', Route), length(splitByChar('-', Route)))])))\n        )) AS ap\n    FROM top_routes\n)\nSELECT\n    ap.1 AS airport_code,\n    any(a.DisplayAirportName) AS airport_name,\n    any(a.CityName) AS city_state,\n    count() AS total_appearances,\n    countIf(ap.3 = 'origin') AS origin_appearances,\n    countIf(ap.3 = 'intermediate_stop') AS intermediate_stop_appearances,\n    countIf(ap.3 = 'final_destination') AS final_destination_appearances,\n    round(uniqExact(Route) / 10.0, 3) AS share_of_itineraries_containing_airport\nFROM route_airports\nLEFT JOIN ontime.dim_airports a\n    ON ap.2 = a.AirportID AND a.IsLatest = 1\nGROUP BY airport_code\nORDER BY total_appearances DESC, intermediate_stop_appearances DESC, airport_code",
      "date_field_hint": "airport_name",
      "row_count": 45,
      "result_columns": [
        "airport_code",
        "airport_name",
        "city_state",
        "total_appearances",
        "origin_appearances",
        "intermediate_stop_appearances",
        "final_destination_appearances",
        "share_of_itineraries_containing_airport"
      ],
      "first_row": {
        "airport_code": "DAL",
        "airport_name": "Dallas Love Field",
        "city_state": "Dallas, TX",
        "final_destination_appearances": 0,
        "intermediate_stop_appearances": 5,
        "origin_appearances": 0,
        "share_of_itineraries_containing_airport": 0.5,
        "total_appearances": 5
      }
    },
    {
      "id": "q2",
      "subquestion": "How geographically extreme is each of the top 10 unique maximum-hop itineraries?\nReturn total flown distance, unique airports, unique city markets, unique states, unique local-time offsets, and whether the route is entirely domestic, then summarize which routes are the most geographically expansive.",
      "answer_markdown": "All 10 top routes are entirely domestic. The most geographically expansive by flown distance is BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX at 5,169 miles; the broadest state coverage is 9 states on BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX, MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX, and ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA; five routes span 4 local-time offsets, while every route touches 9 unique airports.",
      "sql": "WITH params AS (\n    SELECT addYears(toDate((SELECT max(FlightDate) FROM ontime.fact_ontime)), -5) AS start_date,\n           (SELECT max(FlightDate) FROM ontime.fact_ontime) AS end_date\n),\nlegs_raw AS (\n    SELECT\n        FlightDate,\n        ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,\n        Flight_Number_Reporting_Airline AS flight_number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginAirportID,\n        DestAirportID,\n        OriginCode,\n        DestCode,\n        coalesce(CRSDepTime, DepTime, 0) AS dep_hhmm,\n        coalesce(CRSArrTime, ArrTime, 0) AS arr_hhmm,\n        coalesce(DepTime, CRSDepTime, 0) AS dep_sort_hhmm,\n        Distance,\n        row_number() OVER (\n            PARTITION BY FlightDate, ifNull(nullIf(Tail_Number, ''), ''), Flight_Number_Reporting_Airline,\n                         IATA_CODE_Reporting_Airline, OriginAirportID, DestAirportID,\n                         coalesce(CRSDepTime, DepTime, 0), coalesce(CRSArrTime, ArrTime, 0)\n            ORDER BY Cancelled ASC, Diverted ASC, coalesce(ActualElapsedTime, CRSElapsedTime, 1000000) ASC\n        ) AS rn\n    FROM ontime.fact_ontime\n    WHERE FlightDate \u003e= (SELECT start_date FROM params)\n      AND FlightDate \u003c= (SELECT end_date FROM params)\n      AND Cancelled = 0\n),\nlegs AS (\n    SELECT *,\n           toDateTime(FlightDate) + toIntervalMinute(intDiv(dep_sort_hhmm,100)*60 + (dep_sort_hhmm % 100)) AS dep_ts\n    FROM legs_raw\n    WHERE rn = 1\n),\nitineraries AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        arraySort(groupArray((dep_ts, OriginCode, DestCode, OriginAirportID, DestAirportID, toUInt32(ifNull(Distance,0))))) AS legs_sorted\n    FROM legs\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier\n    HAVING length(legs_sorted) \u003e 1\n       AND arrayAll(i -\u003e legs_sorted[i].5 = legs_sorted[i + 1].4, range(1, length(legs_sorted)))\n),\nroutes AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        length(legs_sorted) AS hop_count,\n        legs_sorted[1].1 AS first_dep_ts,\n        arrayMap(x -\u003e x.4, legs_sorted) AS origin_airport_ids,\n        arrayMap(x -\u003e x.5, legs_sorted) AS dest_airport_ids,\n        arrayMap(x -\u003e toUInt32(x.6), legs_sorted) AS leg_distances,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.2, legs_sorted), [legs_sorted[length(legs_sorted)].3]), '-') AS Route\n    FROM itineraries\n),\nroute_days AS (\n    SELECT Route, uniqExact(FlightDate) AS route_recurrence_count\n    FROM routes\n    GROUP BY Route\n),\nunique_ranked AS (\n    SELECT r.*, rd.route_recurrence_count,\n           row_number() OVER (PARTITION BY Route ORDER BY first_dep_ts DESC, FlightDate DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rn\n    FROM routes r\n    INNER JOIN route_days rd USING (Route)\n),\ntop_routes AS (\n    SELECT *,\n           arrayDistinct(arrayConcat(origin_airport_ids, [dest_airport_ids[length(dest_airport_ids)]])) AS route_airport_ids,\n           arraySum(leg_distances) AS total_flown_distance\n    FROM unique_ranked\n    WHERE route_rn = 1\n    ORDER BY hop_count DESC, first_dep_ts DESC\n    LIMIT 10\n),\nroute_airport_dim AS (\n    SELECT\n        tr.Route,\n        tr.total_flown_distance,\n        tr.route_airport_ids,\n        a.CityMarketID,\n        a.StateCode,\n        a.UtcLocalTimeVariation,\n        a.CountryCodeISO\n    FROM top_routes tr\n    ARRAY JOIN tr.route_airport_ids AS airport_id\n    LEFT JOIN ontime.dim_airports a\n        ON airport_id = a.AirportID AND a.IsLatest = 1\n)\nSELECT\n    Route,\n    any(total_flown_distance) AS total_flown_distance,\n    length(any(route_airport_ids)) AS unique_airports,\n    uniqExact(CityMarketID) AS unique_city_markets,\n    uniqExact(StateCode) AS unique_states,\n    uniqExact(UtcLocalTimeVariation) AS unique_local_time_offsets,\n    if(countIf(CountryCodeISO != 'US') = 0, 'Yes', 'No') AS entirely_domestic\nFROM route_airport_dim\nGROUP BY Route\nORDER BY total_flown_distance DESC, unique_states DESC\nLIMIT 10",
      "date_field_hint": "Route",
      "row_count": 10,
      "result_columns": [
        "Route",
        "total_flown_distance",
        "unique_airports",
        "unique_city_markets",
        "unique_states",
        "unique_local_time_offsets",
        "entirely_domestic"
      ],
      "first_row": {
        "Route": "BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX",
        "entirely_domestic": "Yes",
        "total_flown_distance": 5169,
        "unique_airports": 9,
        "unique_city_markets": 8,
        "unique_local_time_offsets": 3,
        "unique_states": 9
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