- Connect to clickhouse server though MCP connection
- Do not use direct HTTP by any tools like curl.
- Use the `ontime` database to answer analytical questions
- Use `ontime-semantic-layer` skill for schema inspection, join guidance, and dimension semantics.
- write correct and efficient ClickHouse SQL 
- Before writing any SQL artifact, self-verify every SQL statement you intend to save.
- Run a cheap debug execution for each query first, usually with a small `LIMIT`, a narrow `WHERE` filter, or both applied inside the main data-reading subquery or CTE.
- Treat successful execution as mandatory. Fix any syntax, type, aggregate, window, join, or unknown-column errors in a loop until every saved query runs successfully.
- Do not write unchecked SQL.

## Output

Write one JSON object to `answer.raw.json` file with shape:

{
  "subquestions": [
    {
      "id": "Section id such as main, q1, q2, or q0.",
      "answer_markdown": "A concise prose answer to that subquestion.",
      "sql": "-- one SQL statement that proves the answer"
    }
  ]
}

## Rules:

- Provide one proof sql query for each parsed section question.
- Return one object in `subquestions` for every parsed section id.
- Preserve each required `id` exactly.
- Each `answer_markdown` must directly answer that section's question in concise business-readable prose.
- All prose claims must be directly traceable to a row in an executed query result.
- When a conclusion depends on shares, concentration, or other denominator-based comparisons, keep the denominator-bearing totals inspectable in the proof query result.
- Use the most recent 5 years by default unless the question-specific prompt explicitly asks for a narrower or wider time window.

## Do not

- Do not merge several section questions into one unioned or row-typed SQL result unless the question-specific guidance explicitly requires that.
- Do not collapse a section to a final aggregate when the presentation needs the underlying ranked, time-series, or drilldown rows.
- Do not emit result rows or any data output.
- Do not output fenced Markdown artifacts.

## Question title 
`Highest daily hops for one aircraft on one flight number`

## Question-specific guidance

Read every top-level `###` section in the question-specific guidance below.

### main
Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.
Define uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.
Do not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.
Count hops from distinct same-day legs, not raw source rows; do not let duplicate or conflicting same-time rows inflate hop count or create artifact routes.

Return:

- aircraft id
- flight number
- carrier
- flight date
- hop count
- route recurrence count: total number of days across the analyzed window on which this exact Route string was flown by any aircraft
- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`

### q1

Which airports act as the key connectors, origins, and termini within the top 10 unique maximum-hop itineraries?
Classify airport appearances by route position and return airport code, airport name, city/state, total appearances, origin appearances, intermediate-stop appearances, final-destination appearances, and share of itineraries containing that airport.

### q2

How geographically extreme is each of the top 10 unique maximum-hop itineraries?
Return total flown distance, unique airports, unique city markets, unique states, unique local-time offsets, and whether the route is entirely domestic, then summarize which routes are the most geographically expansive.

You are reviewing a qforge analysis run after SQL generation, harness execution, and report rendering.

Your job is to judge whether the analysis artifacts actually answer the question correctly and whether the SQL evidence supports the written claims.

Return the final review by writing `review.md` in the run directory.

`review.md` must use exactly this top-level structure:

```md
# Analysis Review
Verdict: PASS

## Summary
...

## Findings
...

## Suggested Prompt Fixes
...
```

Rules:

- Set `Verdict: PASS` only when the analysis is materially aligned with the question.
- Set `Verdict: WARN` when the analysis is materially aligned overall but has limited evidence-support, wording, or minor artifact issues that should not block downstream use.
- Set `Verdict: FAIL` when any substantive correctness, grain, metric, or evidence-support problem exists that makes the run unreliable.
- Base your judgment only on the provided question prompt and the verified run artifacts.
- When artifact file paths are listed, read the files directly from the run directory instead of relying only on prompt excerpts.
- Check whether all required `### main` / `### qN` sections were answered directly.
- Check whether SQL grain and returned metrics match the prompt.
- Check whether any explicitly requested returned fields are actually present in the executed result, and only treat missing identifiers as defects when the prompt clearly requires them to be non-empty.
- Check whether prose claims are supported by the executed query results.
- Check for duplicated entities, inconsistent counts, missing requested metrics, unsupported inference, and proof-query/result mismatches.
- Treat `report.md` as a monitoring artifact. Do not treat its one-row example table or abbreviated preview as a defect when the underlying query result preserves the required ranked or detailed rows.
- In `## Findings`, cite concrete artifact names such as `queries/main.sql`, `queries/q1.sql`, `results/q1.json`, or `report.md`.
- In `## Suggested Prompt Fixes`, propose prompt-level changes only when they would reduce the observed failure mode.
- When proposing a prompt fix, be concrete. Prefer 1-3 short replacement or insertion snippets that could be added to the question prompt, not vague advice.
- Tie each prompt fix to the specific failure you found. Explain what ambiguity in the current prompt allowed the bad result and how the new wording would close that gap.
- If the prompt is already sufficiently specific and the failure is not plausibly prompt-driven, say `None.` instead of inventing weak prompt edits.
- Do not suggest code changes to qforge in this review. Focus on the run and the prompt.
- Write Markdown only to `review.md`.

Question-specific guidance:

### main
Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.
Define uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.
Do not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.
Count hops from distinct same-day legs, not raw source rows; do not let duplicate or conflicting same-time rows inflate hop count or create artifact routes.

Return:

- aircraft id
- flight number
- carrier
- flight date
- hop count
- route recurrence count: total number of days across the analyzed window on which this exact Route string was flown by any aircraft
- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`

### q1

Which airports act as the key connectors, origins, and termini within the top 10 unique maximum-hop itineraries?
Classify airport appearances by route position and return airport code, airport name, city/state, total appearances, origin appearances, intermediate-stop appearances, final-destination appearances, and share of itineraries containing that airport.

### q2

How geographically extreme is each of the top 10 unique maximum-hop itineraries?
Return total flown distance, unique airports, unique city markets, unique states, unique local-time offsets, and whether the route is entirely domestic, then summarize which routes are the most geographically expansive.

Generated report.md:

```md
# Highest daily hops for one aircraft on one flight number

> Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.
Define uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.
Do not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.
Count hops from distinct same-day legs, not raw source rows; do not let duplicate or conflicting same-time rows inflate hop count or create artifact routes.

Return:

- aircraft id
- flight number
- carrier
- flight date
- hop count
- route recurrence count: total number of days across the analyzed window on which this exact Route string was flown by any aircraft
- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`

Across the most recent five years ending 2025-11-30, the longest same-aircraft, same-flight-number itineraries all reached 8 hops. The most recent top 10 unique routes are all Southwest (WN) examples, led by ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA on 2024-12-01; recurrence across the window ranges from 1 day to 46 days, with LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN the most recurrent of the top 10.

- Rows returned: 10
- Columns: aircraft_id, flight_number, carrier, FlightDate, hop_count, route_recurrence_count, Route

| aircraft_id | flight_number | carrier | FlightDate | hop_count | route_recurrence_count | Route |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |

> Which airports act as the key connectors, origins, and termini within the top 10 unique maximum-hop itineraries?
Classify airport appearances by route position and return airport code, airport name, city/state, total appearances, origin appearances, intermediate-stop appearances, final-destination appearances, and share of itineraries containing that airport.

Dallas Love Field, Denver, Las Vegas, Baltimore/Washington, New Orleans, and Oakland are the key nodes, each appearing in 5 of the 10 itineraries. DAL is the strongest pure connector with 5 intermediate-stop appearances; BWI and MSY are the most common origins at 2 each; OAK and LAX are the most common final destinations at 2 each.

- Rows returned: 45
- Columns: airport_code, airport_name, city_state, total_appearances, origin_appearances, intermediate_stop_appearances, final_destination_appearances, share_of_itineraries_containing_airport

| airport_code | airport_name | city_state | total_appearances | origin_appearances | intermediate_stop_appearances | final_destination_appearances | share_of_itineraries_containing_airport |
| --- | --- | --- | --- | --- | --- | --- | --- |
| DAL | Dallas Love Field | Dallas, TX | 5 | 0 | 5 | 0 | 0.5 |

> How geographically extreme is each of the top 10 unique maximum-hop itineraries?
Return total flown distance, unique airports, unique city markets, unique states, unique local-time offsets, and whether the route is entirely domestic, then summarize which routes are the most geographically expansive.

All 10 top routes are entirely domestic. The most geographically expansive by flown distance is BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX at 5,169 miles; the broadest state coverage is 9 states on BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX, MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX, and ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA; five routes span 4 local-time offsets, while every route touches 9 unique airports.

- Rows returned: 10
- Columns: Route, total_flown_distance, unique_airports, unique_city_markets, unique_states, unique_local_time_offsets, entirely_domestic

| Route | total_flown_distance | unique_airports | unique_city_markets | unique_states | unique_local_time_offsets | entirely_domestic |
| --- | --- | --- | --- | --- | --- | --- |
| BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX | 5169 | 9 | 8 | 9 | 3 | Yes |
```

Saved answer.raw.json:

```json
{
  "subquestions": [
    {
      "id": "main",
      "answer_markdown": "Across the most recent five years ending 2025-11-30, the longest same-aircraft, same-flight-number itineraries all reached 8 hops. The most recent top 10 unique routes are all Southwest (WN) examples, led by ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA on 2024-12-01; recurrence across the window ranges from 1 day to 46 days, with LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN the most recurrent of the top 10.",
      "sql": "WITH params AS (\n    SELECT addYears(toDate((SELECT max(FlightDate) FROM ontime.fact_ontime)), -5) AS start_date,\n           (SELECT max(FlightDate) FROM ontime.fact_ontime) AS end_date\n),\nlegs_raw AS (\n    SELECT\n        FlightDate,\n        ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,\n        Flight_Number_Reporting_Airline AS flight_number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginAirportID,\n        DestAirportID,\n        OriginCode,\n        DestCode,\n        coalesce(CRSDepTime, DepTime, 0) AS dep_hhmm,\n        coalesce(CRSArrTime, ArrTime, 0) AS arr_hhmm,\n        coalesce(DepTime, CRSDepTime, 0) AS dep_sort_hhmm,\n        Distance,\n        row_number() OVER (\n            PARTITION BY FlightDate, ifNull(nullIf(Tail_Number, ''), ''), Flight_Number_Reporting_Airline,\n                         IATA_CODE_Reporting_Airline, OriginAirportID, DestAirportID,\n                         coalesce(CRSDepTime, DepTime, 0), coalesce(CRSArrTime, ArrTime, 0)\n            ORDER BY Cancelled ASC, Diverted ASC, coalesce(ActualElapsedTime, CRSElapsedTime, 1000000) ASC\n        ) AS rn\n    FROM ontime.fact_ontime\n    WHERE FlightDate >= (SELECT start_date FROM params)\n      AND FlightDate <= (SELECT end_date FROM params)\n      AND Cancelled = 0\n),\nlegs AS (\n    SELECT *,\n           toDateTime(FlightDate) + toIntervalMinute(intDiv(dep_sort_hhmm,100)*60 + (dep_sort_hhmm % 100)) AS dep_ts\n    FROM legs_raw\n    WHERE rn = 1\n),\nitineraries AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        arraySort(groupArray((dep_ts, OriginCode, DestCode, OriginAirportID, DestAirportID, toUInt32(ifNull(Distance,0))))) AS legs_sorted\n    FROM legs\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier\n    HAVING length(legs_sorted) > 1\n       AND arrayAll(i -> legs_sorted[i].5 = legs_sorted[i + 1].4, range(1, length(legs_sorted)))\n),\nroutes AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        length(legs_sorted) AS hop_count,\n        legs_sorted[1].1 AS first_dep_ts,\n        arrayStringConcat(arrayConcat(arrayMap(x -> x.2, legs_sorted), [legs_sorted[length(legs_sorted)].3]), '-') AS Route\n    FROM itineraries\n),\nroute_days AS (\n    SELECT Route, uniqExact(FlightDate) AS route_recurrence_count\n    FROM routes\n    GROUP BY Route\n),\nunique_ranked AS (\n    SELECT r.*, rd.route_recurrence_count,\n           row_number() OVER (PARTITION BY Route ORDER BY first_dep_ts DESC, FlightDate DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rn\n    FROM routes r\n    INNER JOIN route_days rd USING (Route)\n)\nSELECT aircraft_id, flight_number, carrier, FlightDate, hop_count, route_recurrence_count, Route\nFROM unique_ranked\nWHERE route_rn = 1\nORDER BY hop_count DESC, first_dep_ts DESC\nLIMIT 10"
    },
    {
      "id": "q1",
      "answer_markdown": "Dallas Love Field, Denver, Las Vegas, Baltimore/Washington, New Orleans, and Oakland are the key nodes, each appearing in 5 of the 10 itineraries. DAL is the strongest pure connector with 5 intermediate-stop appearances; BWI and MSY are the most common origins at 2 each; OAK and LAX are the most common final destinations at 2 each.",
      "sql": "WITH params AS (\n    SELECT addYears(toDate((SELECT max(FlightDate) FROM ontime.fact_ontime)), -5) AS start_date,\n           (SELECT max(FlightDate) FROM ontime.fact_ontime) AS end_date\n),\nlegs_raw AS (\n    SELECT\n        FlightDate,\n        ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,\n        Flight_Number_Reporting_Airline AS flight_number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginAirportID,\n        DestAirportID,\n        OriginCode,\n        DestCode,\n        coalesce(CRSDepTime, DepTime, 0) AS dep_hhmm,\n        coalesce(CRSArrTime, ArrTime, 0) AS arr_hhmm,\n        coalesce(DepTime, CRSDepTime, 0) AS dep_sort_hhmm,\n        Distance,\n        row_number() OVER (\n            PARTITION BY FlightDate, ifNull(nullIf(Tail_Number, ''), ''), Flight_Number_Reporting_Airline,\n                         IATA_CODE_Reporting_Airline, OriginAirportID, DestAirportID,\n                         coalesce(CRSDepTime, DepTime, 0), coalesce(CRSArrTime, ArrTime, 0)\n            ORDER BY Cancelled ASC, Diverted ASC, coalesce(ActualElapsedTime, CRSElapsedTime, 1000000) ASC\n        ) AS rn\n    FROM ontime.fact_ontime\n    WHERE FlightDate >= (SELECT start_date FROM params)\n      AND FlightDate <= (SELECT end_date FROM params)\n      AND Cancelled = 0\n),\nlegs AS (\n    SELECT *,\n           toDateTime(FlightDate) + toIntervalMinute(intDiv(dep_sort_hhmm,100)*60 + (dep_sort_hhmm % 100)) AS dep_ts\n    FROM legs_raw\n    WHERE rn = 1\n),\nitineraries AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        arraySort(groupArray((dep_ts, OriginCode, DestCode, OriginAirportID, DestAirportID, toUInt32(ifNull(Distance,0))))) AS legs_sorted\n    FROM legs\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier\n    HAVING length(legs_sorted) > 1\n       AND arrayAll(i -> legs_sorted[i].5 = legs_sorted[i + 1].4, range(1, length(legs_sorted)))\n),\nroutes AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        length(legs_sorted) AS hop_count,\n        legs_sorted[1].1 AS first_dep_ts,\n        arrayMap(x -> x.2, legs_sorted) AS origin_codes,\n        arrayMap(x -> x.4, legs_sorted) AS origin_airport_ids,\n        arrayMap(x -> x.5, legs_sorted) AS dest_airport_ids,\n        arrayStringConcat(arrayConcat(arrayMap(x -> x.2, legs_sorted), [legs_sorted[length(legs_sorted)].3]), '-') AS Route\n    FROM itineraries\n),\nroute_days AS (\n    SELECT Route, uniqExact(FlightDate) AS route_recurrence_count\n    FROM routes\n    GROUP BY Route\n),\nunique_ranked AS (\n    SELECT r.*, rd.route_recurrence_count,\n           row_number() OVER (PARTITION BY Route ORDER BY first_dep_ts DESC, FlightDate DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rn\n    FROM routes r\n    INNER JOIN route_days rd USING (Route)\n),\ntop_routes AS (\n    SELECT *\n    FROM unique_ranked\n    WHERE route_rn = 1\n    ORDER BY hop_count DESC, first_dep_ts DESC\n    LIMIT 10\n),\nroute_airports AS (\n    SELECT\n        Route,\n        arrayJoin(arrayZip(\n            arrayConcat(origin_codes, [arrayElement(splitByChar('-', Route), length(splitByChar('-', Route)))]),\n            arrayConcat(origin_airport_ids, [dest_airport_ids[length(dest_airport_ids)]]),\n            arrayMap(i -> if(i = 1, 'origin', if(i = length(origin_codes) + 1, 'final_destination', 'intermediate_stop')), arrayEnumerate(arrayConcat(origin_codes, [arrayElement(splitByChar('-', Route), length(splitByChar('-', Route)))])))\n        )) AS ap\n    FROM top_routes\n)\nSELECT\n    ap.1 AS airport_code,\n    any(a.DisplayAirportName) AS airport_name,\n    any(a.CityName) AS city_state,\n    count() AS total_appearances,\n    countIf(ap.3 = 'origin') AS origin_appearances,\n    countIf(ap.3 = 'intermediate_stop') AS intermediate_stop_appearances,\n    countIf(ap.3 = 'final_destination') AS final_destination_appearances,\n    round(uniqExact(Route) / 10.0, 3) AS share_of_itineraries_containing_airport\nFROM route_airports\nLEFT JOIN ontime.dim_airports a\n    ON ap.2 = a.AirportID AND a.IsLatest = 1\nGROUP BY airport_code\nORDER BY total_appearances DESC, intermediate_stop_appearances DESC, airport_code"
    },
    {
      "id": "q2",
      "answer_markdown": "All 10 top routes are entirely domestic. The most geographically expansive by flown distance is BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX at 5,169 miles; the broadest state coverage is 9 states on BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX, MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX, and ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA; five routes span 4 local-time offsets, while every route touches 9 unique airports.",
      "sql": "WITH params AS (\n    SELECT addYears(toDate((SELECT max(FlightDate) FROM ontime.fact_ontime)), -5) AS start_date,\n           (SELECT max(FlightDate) FROM ontime.fact_ontime) AS end_date\n),\nlegs_raw AS (\n    SELECT\n        FlightDate,\n        ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,\n        Flight_Number_Reporting_Airline AS flight_number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginAirportID,\n        DestAirportID,\n        OriginCode,\n        DestCode,\n        coalesce(CRSDepTime, DepTime, 0) AS dep_hhmm,\n        coalesce(CRSArrTime, ArrTime, 0) AS arr_hhmm,\n        coalesce(DepTime, CRSDepTime, 0) AS dep_sort_hhmm,\n        Distance,\n        row_number() OVER (\n            PARTITION BY FlightDate, ifNull(nullIf(Tail_Number, ''), ''), Flight_Number_Reporting_Airline,\n                         IATA_CODE_Reporting_Airline, OriginAirportID, DestAirportID,\n                         coalesce(CRSDepTime, DepTime, 0), coalesce(CRSArrTime, ArrTime, 0)\n            ORDER BY Cancelled ASC, Diverted ASC, coalesce(ActualElapsedTime, CRSElapsedTime, 1000000) ASC\n        ) AS rn\n    FROM ontime.fact_ontime\n    WHERE FlightDate >= (SELECT start_date FROM params)\n      AND FlightDate <= (SELECT end_date FROM params)\n      AND Cancelled = 0\n),\nlegs AS (\n    SELECT *,\n           toDateTime(FlightDate) + toIntervalMinute(intDiv(dep_sort_hhmm,100)*60 + (dep_sort_hhmm % 100)) AS dep_ts\n    FROM legs_raw\n    WHERE rn = 1\n),\nitineraries AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        arraySort(groupArray((dep_ts, OriginCode, DestCode, OriginAirportID, DestAirportID, toUInt32(ifNull(Distance,0))))) AS legs_sorted\n    FROM legs\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier\n    HAVING length(legs_sorted) > 1\n       AND arrayAll(i -> legs_sorted[i].5 = legs_sorted[i + 1].4, range(1, length(legs_sorted)))\n),\nroutes AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        length(legs_sorted) AS hop_count,\n        legs_sorted[1].1 AS first_dep_ts,\n        arrayMap(x -> x.4, legs_sorted) AS origin_airport_ids,\n        arrayMap(x -> x.5, legs_sorted) AS dest_airport_ids,\n        arrayMap(x -> toUInt32(x.6), legs_sorted) AS leg_distances,\n        arrayStringConcat(arrayConcat(arrayMap(x -> x.2, legs_sorted), [legs_sorted[length(legs_sorted)].3]), '-') AS Route\n    FROM itineraries\n),\nroute_days AS (\n    SELECT Route, uniqExact(FlightDate) AS route_recurrence_count\n    FROM routes\n    GROUP BY Route\n),\nunique_ranked AS (\n    SELECT r.*, rd.route_recurrence_count,\n           row_number() OVER (PARTITION BY Route ORDER BY first_dep_ts DESC, FlightDate DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rn\n    FROM routes r\n    INNER JOIN route_days rd USING (Route)\n),\ntop_routes AS (\n    SELECT *,\n           arrayDistinct(arrayConcat(origin_airport_ids, [dest_airport_ids[length(dest_airport_ids)]])) AS route_airport_ids,\n           arraySum(leg_distances) AS total_flown_distance\n    FROM unique_ranked\n    WHERE route_rn = 1\n    ORDER BY hop_count DESC, first_dep_ts DESC\n    LIMIT 10\n),\nroute_airport_dim AS (\n    SELECT\n        tr.Route,\n        tr.total_flown_distance,\n        tr.route_airport_ids,\n        a.CityMarketID,\n        a.StateCode,\n        a.UtcLocalTimeVariation,\n        a.CountryCodeISO\n    FROM top_routes tr\n    ARRAY JOIN tr.route_airport_ids AS airport_id\n    LEFT JOIN ontime.dim_airports a\n        ON airport_id = a.AirportID AND a.IsLatest = 1\n)\nSELECT\n    Route,\n    any(total_flown_distance) AS total_flown_distance,\n    length(any(route_airport_ids)) AS unique_airports,\n    uniqExact(CityMarketID) AS unique_city_markets,\n    uniqExact(StateCode) AS unique_states,\n    uniqExact(UtcLocalTimeVariation) AS unique_local_time_offsets,\n    if(countIf(CountryCodeISO != 'US') = 0, 'Yes', 'No') AS entirely_domestic\nFROM route_airport_dim\nGROUP BY Route\nORDER BY total_flown_distance DESC, unique_states DESC\nLIMIT 10"
    }
  ]
}
```

Saved analysis.json:

```json
{
  "sql": "WITH params AS (\n    SELECT addYears(toDate((SELECT max(FlightDate) FROM ontime.fact_ontime)), -5) AS start_date,\n           (SELECT max(FlightDate) FROM ontime.fact_ontime) AS end_date\n),\nlegs_raw AS (\n    SELECT\n        FlightDate,\n        ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,\n        Flight_Number_Reporting_Airline AS flight_number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginAirportID,\n        DestAirportID,\n        OriginCode,\n        DestCode,\n        coalesce(CRSDepTime, DepTime, 0) AS dep_hhmm,\n        coalesce(CRSArrTime, ArrTime, 0) AS arr_hhmm,\n        coalesce(DepTime, CRSDepTime, 0) AS dep_sort_hhmm,\n        Distance,\n        row_number() OVER (\n            PARTITION BY FlightDate, ifNull(nullIf(Tail_Number, ''), ''), Flight_Number_Reporting_Airline,\n                         IATA_CODE_Reporting_Airline, OriginAirportID, DestAirportID,\n                         coalesce(CRSDepTime, DepTime, 0), coalesce(CRSArrTime, ArrTime, 0)\n            ORDER BY Cancelled ASC, Diverted ASC, coalesce(ActualElapsedTime, CRSElapsedTime, 1000000) ASC\n        ) AS rn\n    FROM ontime.fact_ontime\n    WHERE FlightDate \u003e= (SELECT start_date FROM params)\n      AND FlightDate \u003c= (SELECT end_date FROM params)\n      AND Cancelled = 0\n),\nlegs AS (\n    SELECT *,\n           toDateTime(FlightDate) + toIntervalMinute(intDiv(dep_sort_hhmm,100)*60 + (dep_sort_hhmm % 100)) AS dep_ts\n    FROM legs_raw\n    WHERE rn = 1\n),\nitineraries AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        arraySort(groupArray((dep_ts, OriginCode, DestCode, OriginAirportID, DestAirportID, toUInt32(ifNull(Distance,0))))) AS legs_sorted\n    FROM legs\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier\n    HAVING length(legs_sorted) \u003e 1\n       AND arrayAll(i -\u003e legs_sorted[i].5 = legs_sorted[i + 1].4, range(1, length(legs_sorted)))\n),\nroutes AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        length(legs_sorted) AS hop_count,\n        legs_sorted[1].1 AS first_dep_ts,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.2, legs_sorted), [legs_sorted[length(legs_sorted)].3]), '-') AS Route\n    FROM itineraries\n),\nroute_days AS (\n    SELECT Route, uniqExact(FlightDate) AS route_recurrence_count\n    FROM routes\n    GROUP BY Route\n),\nunique_ranked AS (\n    SELECT r.*, rd.route_recurrence_count,\n           row_number() OVER (PARTITION BY Route ORDER BY first_dep_ts DESC, FlightDate DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rn\n    FROM routes r\n    INNER JOIN route_days rd USING (Route)\n)\nSELECT aircraft_id, flight_number, carrier, FlightDate, hop_count, route_recurrence_count, Route\nFROM unique_ranked\nWHERE route_rn = 1\nORDER BY hop_count DESC, first_dep_ts DESC\nLIMIT 10",
  "report_markdown": "",
  "subquestions": [
    {
      "id": "main",
      "subquestion": "Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.\nDefine uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.\nDo not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.\nCount hops from distinct same-day legs, not raw source rows; do not let duplicate or conflicting same-time rows inflate hop count or create artifact routes.\n\nReturn:\n\n- aircraft id\n- flight number\n- carrier\n- flight date\n- hop count\n- route recurrence count: total number of days across the analyzed window on which this exact Route string was flown by any aircraft\n- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`",
      "answer_markdown": "Across the most recent five years ending 2025-11-30, the longest same-aircraft, same-flight-number itineraries all reached 8 hops. The most recent top 10 unique routes are all Southwest (WN) examples, led by ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA on 2024-12-01; recurrence across the window ranges from 1 day to 46 days, with LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN the most recurrent of the top 10.",
      "sql": "WITH params AS (\n    SELECT addYears(toDate((SELECT max(FlightDate) FROM ontime.fact_ontime)), -5) AS start_date,\n           (SELECT max(FlightDate) FROM ontime.fact_ontime) AS end_date\n),\nlegs_raw AS (\n    SELECT\n        FlightDate,\n        ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,\n        Flight_Number_Reporting_Airline AS flight_number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginAirportID,\n        DestAirportID,\n        OriginCode,\n        DestCode,\n        coalesce(CRSDepTime, DepTime, 0) AS dep_hhmm,\n        coalesce(CRSArrTime, ArrTime, 0) AS arr_hhmm,\n        coalesce(DepTime, CRSDepTime, 0) AS dep_sort_hhmm,\n        Distance,\n        row_number() OVER (\n            PARTITION BY FlightDate, ifNull(nullIf(Tail_Number, ''), ''), Flight_Number_Reporting_Airline,\n                         IATA_CODE_Reporting_Airline, OriginAirportID, DestAirportID,\n                         coalesce(CRSDepTime, DepTime, 0), coalesce(CRSArrTime, ArrTime, 0)\n            ORDER BY Cancelled ASC, Diverted ASC, coalesce(ActualElapsedTime, CRSElapsedTime, 1000000) ASC\n        ) AS rn\n    FROM ontime.fact_ontime\n    WHERE FlightDate \u003e= (SELECT start_date FROM params)\n      AND FlightDate \u003c= (SELECT end_date FROM params)\n      AND Cancelled = 0\n),\nlegs AS (\n    SELECT *,\n           toDateTime(FlightDate) + toIntervalMinute(intDiv(dep_sort_hhmm,100)*60 + (dep_sort_hhmm % 100)) AS dep_ts\n    FROM legs_raw\n    WHERE rn = 1\n),\nitineraries AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        arraySort(groupArray((dep_ts, OriginCode, DestCode, OriginAirportID, DestAirportID, toUInt32(ifNull(Distance,0))))) AS legs_sorted\n    FROM legs\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier\n    HAVING length(legs_sorted) \u003e 1\n       AND arrayAll(i -\u003e legs_sorted[i].5 = legs_sorted[i + 1].4, range(1, length(legs_sorted)))\n),\nroutes AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        length(legs_sorted) AS hop_count,\n        legs_sorted[1].1 AS first_dep_ts,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.2, legs_sorted), [legs_sorted[length(legs_sorted)].3]), '-') AS Route\n    FROM itineraries\n),\nroute_days AS (\n    SELECT Route, uniqExact(FlightDate) AS route_recurrence_count\n    FROM routes\n    GROUP BY Route\n),\nunique_ranked AS (\n    SELECT r.*, rd.route_recurrence_count,\n           row_number() OVER (PARTITION BY Route ORDER BY first_dep_ts DESC, FlightDate DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rn\n    FROM routes r\n    INNER JOIN route_days rd USING (Route)\n)\nSELECT aircraft_id, flight_number, carrier, FlightDate, hop_count, route_recurrence_count, Route\nFROM unique_ranked\nWHERE route_rn = 1\nORDER BY hop_count DESC, first_dep_ts DESC\nLIMIT 10"
    },
    {
      "id": "q1",
      "subquestion": "Which airports act as the key connectors, origins, and termini within the top 10 unique maximum-hop itineraries?\nClassify airport appearances by route position and return airport code, airport name, city/state, total appearances, origin appearances, intermediate-stop appearances, final-destination appearances, and share of itineraries containing that airport.",
      "answer_markdown": "Dallas Love Field, Denver, Las Vegas, Baltimore/Washington, New Orleans, and Oakland are the key nodes, each appearing in 5 of the 10 itineraries. DAL is the strongest pure connector with 5 intermediate-stop appearances; BWI and MSY are the most common origins at 2 each; OAK and LAX are the most common final destinations at 2 each.",
      "sql": "WITH params AS (\n    SELECT addYears(toDate((SELECT max(FlightDate) FROM ontime.fact_ontime)), -5) AS start_date,\n           (SELECT max(FlightDate) FROM ontime.fact_ontime) AS end_date\n),\nlegs_raw AS (\n    SELECT\n        FlightDate,\n        ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,\n        Flight_Number_Reporting_Airline AS flight_number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginAirportID,\n        DestAirportID,\n        OriginCode,\n        DestCode,\n        coalesce(CRSDepTime, DepTime, 0) AS dep_hhmm,\n        coalesce(CRSArrTime, ArrTime, 0) AS arr_hhmm,\n        coalesce(DepTime, CRSDepTime, 0) AS dep_sort_hhmm,\n        Distance,\n        row_number() OVER (\n            PARTITION BY FlightDate, ifNull(nullIf(Tail_Number, ''), ''), Flight_Number_Reporting_Airline,\n                         IATA_CODE_Reporting_Airline, OriginAirportID, DestAirportID,\n                         coalesce(CRSDepTime, DepTime, 0), coalesce(CRSArrTime, ArrTime, 0)\n            ORDER BY Cancelled ASC, Diverted ASC, coalesce(ActualElapsedTime, CRSElapsedTime, 1000000) ASC\n        ) AS rn\n    FROM ontime.fact_ontime\n    WHERE FlightDate \u003e= (SELECT start_date FROM params)\n      AND FlightDate \u003c= (SELECT end_date FROM params)\n      AND Cancelled = 0\n),\nlegs AS (\n    SELECT *,\n           toDateTime(FlightDate) + toIntervalMinute(intDiv(dep_sort_hhmm,100)*60 + (dep_sort_hhmm % 100)) AS dep_ts\n    FROM legs_raw\n    WHERE rn = 1\n),\nitineraries AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        arraySort(groupArray((dep_ts, OriginCode, DestCode, OriginAirportID, DestAirportID, toUInt32(ifNull(Distance,0))))) AS legs_sorted\n    FROM legs\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier\n    HAVING length(legs_sorted) \u003e 1\n       AND arrayAll(i -\u003e legs_sorted[i].5 = legs_sorted[i + 1].4, range(1, length(legs_sorted)))\n),\nroutes AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        length(legs_sorted) AS hop_count,\n        legs_sorted[1].1 AS first_dep_ts,\n        arrayMap(x -\u003e x.2, legs_sorted) AS origin_codes,\n        arrayMap(x -\u003e x.4, legs_sorted) AS origin_airport_ids,\n        arrayMap(x -\u003e x.5, legs_sorted) AS dest_airport_ids,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.2, legs_sorted), [legs_sorted[length(legs_sorted)].3]), '-') AS Route\n    FROM itineraries\n),\nroute_days AS (\n    SELECT Route, uniqExact(FlightDate) AS route_recurrence_count\n    FROM routes\n    GROUP BY Route\n),\nunique_ranked AS (\n    SELECT r.*, rd.route_recurrence_count,\n           row_number() OVER (PARTITION BY Route ORDER BY first_dep_ts DESC, FlightDate DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rn\n    FROM routes r\n    INNER JOIN route_days rd USING (Route)\n),\ntop_routes AS (\n    SELECT *\n    FROM unique_ranked\n    WHERE route_rn = 1\n    ORDER BY hop_count DESC, first_dep_ts DESC\n    LIMIT 10\n),\nroute_airports AS (\n    SELECT\n        Route,\n        arrayJoin(arrayZip(\n            arrayConcat(origin_codes, [arrayElement(splitByChar('-', Route), length(splitByChar('-', Route)))]),\n            arrayConcat(origin_airport_ids, [dest_airport_ids[length(dest_airport_ids)]]),\n            arrayMap(i -\u003e if(i = 1, 'origin', if(i = length(origin_codes) + 1, 'final_destination', 'intermediate_stop')), arrayEnumerate(arrayConcat(origin_codes, [arrayElement(splitByChar('-', Route), length(splitByChar('-', Route)))])))\n        )) AS ap\n    FROM top_routes\n)\nSELECT\n    ap.1 AS airport_code,\n    any(a.DisplayAirportName) AS airport_name,\n    any(a.CityName) AS city_state,\n    count() AS total_appearances,\n    countIf(ap.3 = 'origin') AS origin_appearances,\n    countIf(ap.3 = 'intermediate_stop') AS intermediate_stop_appearances,\n    countIf(ap.3 = 'final_destination') AS final_destination_appearances,\n    round(uniqExact(Route) / 10.0, 3) AS share_of_itineraries_containing_airport\nFROM route_airports\nLEFT JOIN ontime.dim_airports a\n    ON ap.2 = a.AirportID AND a.IsLatest = 1\nGROUP BY airport_code\nORDER BY total_appearances DESC, intermediate_stop_appearances DESC, airport_code"
    },
    {
      "id": "q2",
      "subquestion": "How geographically extreme is each of the top 10 unique maximum-hop itineraries?\nReturn total flown distance, unique airports, unique city markets, unique states, unique local-time offsets, and whether the route is entirely domestic, then summarize which routes are the most geographically expansive.",
      "answer_markdown": "All 10 top routes are entirely domestic. The most geographically expansive by flown distance is BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX at 5,169 miles; the broadest state coverage is 9 states on BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX, MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX, and ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA; five routes span 4 local-time offsets, while every route touches 9 unique airports.",
      "sql": "WITH params AS (\n    SELECT addYears(toDate((SELECT max(FlightDate) FROM ontime.fact_ontime)), -5) AS start_date,\n           (SELECT max(FlightDate) FROM ontime.fact_ontime) AS end_date\n),\nlegs_raw AS (\n    SELECT\n        FlightDate,\n        ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,\n        Flight_Number_Reporting_Airline AS flight_number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginAirportID,\n        DestAirportID,\n        OriginCode,\n        DestCode,\n        coalesce(CRSDepTime, DepTime, 0) AS dep_hhmm,\n        coalesce(CRSArrTime, ArrTime, 0) AS arr_hhmm,\n        coalesce(DepTime, CRSDepTime, 0) AS dep_sort_hhmm,\n        Distance,\n        row_number() OVER (\n            PARTITION BY FlightDate, ifNull(nullIf(Tail_Number, ''), ''), Flight_Number_Reporting_Airline,\n                         IATA_CODE_Reporting_Airline, OriginAirportID, DestAirportID,\n                         coalesce(CRSDepTime, DepTime, 0), coalesce(CRSArrTime, ArrTime, 0)\n            ORDER BY Cancelled ASC, Diverted ASC, coalesce(ActualElapsedTime, CRSElapsedTime, 1000000) ASC\n        ) AS rn\n    FROM ontime.fact_ontime\n    WHERE FlightDate \u003e= (SELECT start_date FROM params)\n      AND FlightDate \u003c= (SELECT end_date FROM params)\n      AND Cancelled = 0\n),\nlegs AS (\n    SELECT *,\n           toDateTime(FlightDate) + toIntervalMinute(intDiv(dep_sort_hhmm,100)*60 + (dep_sort_hhmm % 100)) AS dep_ts\n    FROM legs_raw\n    WHERE rn = 1\n),\nitineraries AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        arraySort(groupArray((dep_ts, OriginCode, DestCode, OriginAirportID, DestAirportID, toUInt32(ifNull(Distance,0))))) AS legs_sorted\n    FROM legs\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier\n    HAVING length(legs_sorted) \u003e 1\n       AND arrayAll(i -\u003e legs_sorted[i].5 = legs_sorted[i + 1].4, range(1, length(legs_sorted)))\n),\nroutes AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        length(legs_sorted) AS hop_count,\n        legs_sorted[1].1 AS first_dep_ts,\n        arrayMap(x -\u003e x.4, legs_sorted) AS origin_airport_ids,\n        arrayMap(x -\u003e x.5, legs_sorted) AS dest_airport_ids,\n        arrayMap(x -\u003e toUInt32(x.6), legs_sorted) AS leg_distances,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.2, legs_sorted), [legs_sorted[length(legs_sorted)].3]), '-') AS Route\n    FROM itineraries\n),\nroute_days AS (\n    SELECT Route, uniqExact(FlightDate) AS route_recurrence_count\n    FROM routes\n    GROUP BY Route\n),\nunique_ranked AS (\n    SELECT r.*, rd.route_recurrence_count,\n           row_number() OVER (PARTITION BY Route ORDER BY first_dep_ts DESC, FlightDate DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rn\n    FROM routes r\n    INNER JOIN route_days rd USING (Route)\n),\ntop_routes AS (\n    SELECT *,\n           arrayDistinct(arrayConcat(origin_airport_ids, [dest_airport_ids[length(dest_airport_ids)]])) AS route_airport_ids,\n           arraySum(leg_distances) AS total_flown_distance\n    FROM unique_ranked\n    WHERE route_rn = 1\n    ORDER BY hop_count DESC, first_dep_ts DESC\n    LIMIT 10\n),\nroute_airport_dim AS (\n    SELECT\n        tr.Route,\n        tr.total_flown_distance,\n        tr.route_airport_ids,\n        a.CityMarketID,\n        a.StateCode,\n        a.UtcLocalTimeVariation,\n        a.CountryCodeISO\n    FROM top_routes tr\n    ARRAY JOIN tr.route_airport_ids AS airport_id\n    LEFT JOIN ontime.dim_airports a\n        ON airport_id = a.AirportID AND a.IsLatest = 1\n)\nSELECT\n    Route,\n    any(total_flown_distance) AS total_flown_distance,\n    length(any(route_airport_ids)) AS unique_airports,\n    uniqExact(CityMarketID) AS unique_city_markets,\n    uniqExact(StateCode) AS unique_states,\n    uniqExact(UtcLocalTimeVariation) AS unique_local_time_offsets,\n    if(countIf(CountryCodeISO != 'US') = 0, 'Yes', 'No') AS entirely_domestic\nFROM route_airport_dim\nGROUP BY Route\nORDER BY total_flown_distance DESC, unique_states DESC\nLIMIT 10"
    }
  ]
}
```

Saved visual_input.json:

```json
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
```

Proof queries are saved as files in the run directory. Read the SQL files you need to verify grain, filters, metrics, and ranking logic:

- `queries/main.sql`
- `queries/q1.sql`
- `queries/q2.sql`

Executed query results are saved as files in the run directory. Read the result files you need for verification instead of assuming the report summary is complete:

- `results/main.json`
- `results/q1.json`
- `results/q2.json`