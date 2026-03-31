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
- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft
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
- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft
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
- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft
- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`

Within the most recent 5 years of available data, the highest observed same-aircraft, same-flight-number itineraries reached 8 hops after deduplicating same-day leg rows. The most recent top route was WN 366 on 2024-12-01 by aircraft N957WN, flying ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA. Among the top 10 unique 8-hop routes, recurrence ranged from 1 day to 46 days, with LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN recurring most often at 46 days and MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC next at 40 days.

- Rows returned: 10
- Columns: aircraft_id, flight_number, carrier, flight_date, hop_count, route_recurrence_count, Route

| aircraft_id | flight_number | carrier | flight_date | hop_count | route_recurrence_count | Route |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |

> Which airports act as the key connectors, origins, and termini within the top 10 unique maximum-hop itineraries?
Classify airport appearances by route position and return airport code, airport name, city/state, total appearances, origin appearances, intermediate-stop appearances, final-destination appearances, and share of itineraries containing that airport.

Across the top 10 unique maximum-hop itineraries, the strongest connector airports were DAL, DEN, LAS, MSY, OAK, and BWI, each appearing in 5 of the 10 itineraries. DAL was purely an intermediate stop in all 5 appearances, while BWI and MSY were the most common origins at 2 appearances each. The most common termini were OAK and LAX with 2 final-destination appearances each. BWI, DAL, DEN, LAS, MSY, and OAK each appeared in 50% of the itineraries.

- Rows returned: 45
- Columns: airport_code, airport_name, city_state, total_appearances, origin_appearances, intermediate_stop_appearances, final_destination_appearances, share_of_itineraries_containing_airport

| airport_code | airport_name | city_state | total_appearances | origin_appearances | intermediate_stop_appearances | final_destination_appearances | share_of_itineraries_containing_airport |
| --- | --- | --- | --- | --- | --- | --- | --- |
| BWI | Baltimore/Washington International Thurgood Marshall | Baltimore, MD, MD | 5 | 2 | 3 | 0 | 0.5 |

> How geographically extreme is each of the top 10 unique maximum-hop itineraries?
Return total flown distance, unique airports, unique city markets, unique states, unique local-time offsets, and whether the route is entirely domestic, then summarize which routes are the most geographically expansive.

All top 10 unique maximum-hop itineraries were entirely domestic. The most geographically expansive by total flown distance was BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK at 5,169 miles, spanning 9 airports, 8 city markets, 9 states, and 3 local-time offsets. Several routes were broader in time-zone spread, with CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN, ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN, HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK, MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC, and BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK each touching 4 distinct local-time offsets.

- Rows returned: 10
- Columns: Route, aircraft_id, flight_number, carrier, flight_date, hop_count, total_flown_distance, unique_airports, unique_city_markets, unique_states, unique_local_time_offsets, entirely_domestic

| Route | aircraft_id | flight_number | carrier | flight_date | hop_count | total_flown_distance | unique_airports | unique_city_markets | unique_states | unique_local_time_offsets | entirely_domestic |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX | N225WN | 3530 | WN | 2021-08-08T00:00:00Z | 8 | 5169 | 9 | 8 | 9 | 3 | 1 |
```

Saved answer.raw.json:

```json
{
  "subquestions": [
    {
      "id": "main",
      "answer_markdown": "Within the most recent 5 years of available data, the highest observed same-aircraft, same-flight-number itineraries reached 8 hops after deduplicating same-day leg rows. The most recent top route was WN 366 on 2024-12-01 by aircraft N957WN, flying ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA. Among the top 10 unique 8-hop routes, recurrence ranged from 1 day to 46 days, with LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN recurring most often at 46 days and MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC next at 40 days.",
      "sql": "WITH max_dt AS (\n    SELECT max(FlightDate) AS mx\n    FROM ontime.fact_ontime\n),\nbase AS (\n    SELECT\n        f.FlightDate,\n        f.IATA_CODE_Reporting_Airline AS carrier,\n        f.Flight_Number_Reporting_Airline AS flight_number,\n        ifNull(nullIf(f.Tail_Number, ''), '') AS aircraft_id,\n        f.OriginCode AS origin_code,\n        f.DestCode AS dest_code,\n        f.OriginAirportID AS origin_airport_id,\n        f.DestAirportID AS dest_airport_id,\n        f.OriginCityMarketID AS origin_city_market_id,\n        f.DestCityMarketID AS dest_city_market_id,\n        coalesce(f.CRSDepTime, f.DepTime, 9999) AS dep_sort_time,\n        coalesce(f.Distance, 0) AS distance,\n        row_number() OVER (\n            PARTITION BY\n                f.FlightDate,\n                f.IATA_CODE_Reporting_Airline,\n                f.Flight_Number_Reporting_Airline,\n                ifNull(nullIf(f.Tail_Number, ''), ''),\n                f.OriginCode,\n                f.DestCode,\n                coalesce(f.CRSDepTime, f.DepTime, 9999)\n            ORDER BY\n                f.DepTime DESC NULLS LAST,\n                f.ArrTime DESC NULLS LAST,\n                coalesce(f.Distance, 0) DESC\n        ) AS rn\n    FROM ontime.fact_ontime AS f\n    CROSS JOIN max_dt\n    WHERE f.Cancelled = 0\n      AND f.FlightDate >= addYears(mx, -5)\n),\nlegs AS (\n    SELECT *\n    FROM base\n    WHERE rn = 1\n),\nitineraries AS (\n    SELECT\n        FlightDate,\n        carrier,\n        flight_number,\n        aircraft_id,\n        arraySort(groupArray((\n            dep_sort_time,\n            origin_code,\n            dest_code,\n            origin_airport_id,\n            dest_airport_id,\n            origin_city_market_id,\n            dest_city_market_id,\n            distance\n        ))) AS legs_arr,\n        min(dep_sort_time) AS first_dep_time\n    FROM legs\n    GROUP BY FlightDate, carrier, flight_number, aircraft_id\n),\nroutes AS (\n    SELECT\n        FlightDate,\n        carrier,\n        flight_number,\n        aircraft_id,\n        first_dep_time,\n        length(legs_arr) AS hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                [tupleElement(arrayElement(legs_arr, 1), 2)],\n                arrayMap(x -> tupleElement(x, 3), legs_arr)\n            ),\n            '-'\n        ) AS Route\n    FROM itineraries\n),\nroute_recurrence AS (\n    SELECT\n        Route,\n        uniqExact(FlightDate) AS route_recurrence_count\n    FROM routes\n    GROUP BY Route\n),\nranked_unique AS (\n    SELECT\n        r.FlightDate,\n        r.carrier,\n        r.flight_number,\n        r.aircraft_id,\n        r.first_dep_time,\n        r.hop_count,\n        rr.route_recurrence_count,\n        r.Route,\n        row_number() OVER (\n            PARTITION BY r.Route\n            ORDER BY\n                r.hop_count DESC,\n                r.FlightDate DESC,\n                r.first_dep_time DESC,\n                r.carrier,\n                r.flight_number,\n                r.aircraft_id\n        ) AS route_pick_rank\n    FROM routes AS r\n    INNER JOIN route_recurrence AS rr\n        ON r.Route = rr.Route\n)\nSELECT\n    if(aircraft_id = '', 'unknown', aircraft_id) AS aircraft_id,\n    flight_number,\n    carrier,\n    FlightDate AS flight_date,\n    hop_count,\n    route_recurrence_count,\n    Route\nFROM ranked_unique\nWHERE route_pick_rank = 1\nORDER BY hop_count DESC, flight_date DESC, first_dep_time DESC, Route\nLIMIT 10"
    },
    {
      "id": "q1",
      "answer_markdown": "Across the top 10 unique maximum-hop itineraries, the strongest connector airports were DAL, DEN, LAS, MSY, OAK, and BWI, each appearing in 5 of the 10 itineraries. DAL was purely an intermediate stop in all 5 appearances, while BWI and MSY were the most common origins at 2 appearances each. The most common termini were OAK and LAX with 2 final-destination appearances each. BWI, DAL, DEN, LAS, MSY, and OAK each appeared in 50% of the itineraries.",
      "sql": "WITH max_dt AS (\n    SELECT max(FlightDate) AS mx\n    FROM ontime.fact_ontime\n),\nbase AS (\n    SELECT\n        f.FlightDate,\n        f.IATA_CODE_Reporting_Airline AS carrier,\n        f.Flight_Number_Reporting_Airline AS flight_number,\n        ifNull(nullIf(f.Tail_Number, ''), '') AS aircraft_id,\n        f.OriginCode AS origin_code,\n        f.DestCode AS dest_code,\n        f.OriginAirportID AS origin_airport_id,\n        f.DestAirportID AS dest_airport_id,\n        coalesce(f.CRSDepTime, f.DepTime, 9999) AS dep_sort_time,\n        coalesce(f.Distance, 0) AS distance,\n        row_number() OVER (\n            PARTITION BY\n                f.FlightDate,\n                f.IATA_CODE_Reporting_Airline,\n                f.Flight_Number_Reporting_Airline,\n                ifNull(nullIf(f.Tail_Number, ''), ''),\n                f.OriginCode,\n                f.DestCode,\n                coalesce(f.CRSDepTime, f.DepTime, 9999)\n            ORDER BY\n                f.DepTime DESC NULLS LAST,\n                f.ArrTime DESC NULLS LAST,\n                coalesce(f.Distance, 0) DESC\n        ) AS rn\n    FROM ontime.fact_ontime AS f\n    CROSS JOIN max_dt\n    WHERE f.Cancelled = 0\n      AND f.FlightDate >= addYears(mx, -5)\n),\nlegs AS (\n    SELECT *\n    FROM base\n    WHERE rn = 1\n),\nitineraries AS (\n    SELECT\n        FlightDate,\n        carrier,\n        flight_number,\n        aircraft_id,\n        arraySort(groupArray((\n            dep_sort_time,\n            origin_code,\n            dest_code,\n            origin_airport_id,\n            dest_airport_id,\n            distance\n        ))) AS legs_arr,\n        min(dep_sort_time) AS first_dep_time\n    FROM legs\n    GROUP BY FlightDate, carrier, flight_number, aircraft_id\n),\nroutes AS (\n    SELECT\n        FlightDate,\n        carrier,\n        flight_number,\n        aircraft_id,\n        first_dep_time,\n        length(legs_arr) AS hop_count,\n        legs_arr,\n        arrayStringConcat(\n            arrayConcat(\n                [tupleElement(arrayElement(legs_arr, 1), 2)],\n                arrayMap(x -> tupleElement(x, 3), legs_arr)\n            ),\n            '-'\n        ) AS Route\n    FROM itineraries\n),\nroute_recurrence AS (\n    SELECT\n        Route,\n        uniqExact(FlightDate) AS route_recurrence_count\n    FROM routes\n    GROUP BY Route\n),\nranked_unique AS (\n    SELECT\n        r.*,\n        rr.route_recurrence_count,\n        row_number() OVER (\n            PARTITION BY r.Route\n            ORDER BY\n                r.hop_count DESC,\n                r.FlightDate DESC,\n                r.first_dep_time DESC,\n                r.carrier,\n                r.flight_number,\n                r.aircraft_id\n        ) AS route_pick_rank\n    FROM routes AS r\n    INNER JOIN route_recurrence AS rr\n        ON r.Route = rr.Route\n),\ntop_routes AS (\n    SELECT *\n    FROM ranked_unique\n    WHERE route_pick_rank = 1\n    ORDER BY hop_count DESC, FlightDate DESC, first_dep_time DESC, Route\n    LIMIT 10\n),\nroute_airports AS (\n    SELECT\n        Route,\n        arrayJoin(\n            arrayZip(\n                arrayConcat(\n                    [tupleElement(arrayElement(legs_arr, 1), 2)],\n                    arrayMap(x -> tupleElement(x, 3), legs_arr)\n                ),\n                arrayEnumerate(\n                    arrayConcat(\n                        [tupleElement(arrayElement(legs_arr, 1), 2)],\n                        arrayMap(x -> tupleElement(x, 3), legs_arr)\n                    )\n                ),\n                arrayMap(\n                    x -> length(legs_arr) + 1,\n                    arrayConcat(\n                        [tupleElement(arrayElement(legs_arr, 1), 2)],\n                        arrayMap(x -> tupleElement(x, 3), legs_arr)\n                    )\n                )\n            )\n        ) AS airport_pos\n    FROM top_routes\n)\nSELECT\n    trimBoth(toString(airport_pos.1)) AS airport_code,\n    any(a.DisplayAirportName) AS airport_name,\n    any(concat(a.CityName, ', ', a.StateCode)) AS city_state,\n    count() AS total_appearances,\n    countIf(airport_pos.2 = 1) AS origin_appearances,\n    countIf(airport_pos.2 > 1 AND airport_pos.2 < airport_pos.3) AS intermediate_stop_appearances,\n    countIf(airport_pos.2 = airport_pos.3) AS final_destination_appearances,\n    round(uniqExact(Route) / 10.0, 3) AS share_of_itineraries_containing_airport\nFROM route_airports AS ra\nLEFT JOIN ontime.dim_airports AS a\n    ON a.AirportCode = ra.airport_pos.1\n   AND a.IsLatest = 1\nGROUP BY airport_code\nORDER BY total_appearances DESC, share_of_itineraries_containing_airport DESC, airport_code"
    },
    {
      "id": "q2",
      "answer_markdown": "All top 10 unique maximum-hop itineraries were entirely domestic. The most geographically expansive by total flown distance was BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK at 5,169 miles, spanning 9 airports, 8 city markets, 9 states, and 3 local-time offsets. Several routes were broader in time-zone spread, with CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN, ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN, HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK, MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC, and BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK each touching 4 distinct local-time offsets.",
      "sql": "WITH max_dt AS (\n    SELECT max(FlightDate) AS mx\n    FROM ontime.fact_ontime\n),\nbase AS (\n    SELECT\n        f.FlightDate,\n        f.IATA_CODE_Reporting_Airline AS carrier,\n        f.Flight_Number_Reporting_Airline AS flight_number,\n        ifNull(nullIf(f.Tail_Number, ''), '') AS aircraft_id,\n        f.OriginCode AS origin_code,\n        f.DestCode AS dest_code,\n        f.OriginAirportID AS origin_airport_id,\n        f.DestAirportID AS dest_airport_id,\n        f.OriginCityMarketID AS origin_city_market_id,\n        f.DestCityMarketID AS dest_city_market_id,\n        coalesce(f.CRSDepTime, f.DepTime, 9999) AS dep_sort_time,\n        coalesce(f.Distance, 0) AS distance,\n        row_number() OVER (\n            PARTITION BY\n                f.FlightDate,\n                f.IATA_CODE_Reporting_Airline,\n                f.Flight_Number_Reporting_Airline,\n                ifNull(nullIf(f.Tail_Number, ''), ''),\n                f.OriginCode,\n                f.DestCode,\n                coalesce(f.CRSDepTime, f.DepTime, 9999)\n            ORDER BY\n                f.DepTime DESC NULLS LAST,\n                f.ArrTime DESC NULLS LAST,\n                coalesce(f.Distance, 0) DESC\n        ) AS rn\n    FROM ontime.fact_ontime AS f\n    CROSS JOIN max_dt\n    WHERE f.Cancelled = 0\n      AND f.FlightDate >= addYears(mx, -5)\n),\nlegs AS (\n    SELECT *\n    FROM base\n    WHERE rn = 1\n),\nitineraries AS (\n    SELECT\n        FlightDate,\n        carrier,\n        flight_number,\n        aircraft_id,\n        arraySort(groupArray((\n            dep_sort_time,\n            origin_code,\n            dest_code,\n            origin_airport_id,\n            dest_airport_id,\n            origin_city_market_id,\n            dest_city_market_id,\n            distance\n        ))) AS legs_arr,\n        min(dep_sort_time) AS first_dep_time\n    FROM legs\n    GROUP BY FlightDate, carrier, flight_number, aircraft_id\n),\nroutes AS (\n    SELECT\n        FlightDate,\n        carrier,\n        flight_number,\n        aircraft_id,\n        first_dep_time,\n        length(legs_arr) AS hop_count,\n        legs_arr,\n        arrayStringConcat(\n            arrayConcat(\n                [tupleElement(arrayElement(legs_arr, 1), 2)],\n                arrayMap(x -> tupleElement(x, 3), legs_arr)\n            ),\n            '-'\n        ) AS Route\n    FROM itineraries\n),\nroute_recurrence AS (\n    SELECT\n        Route,\n        uniqExact(FlightDate) AS route_recurrence_count\n    FROM routes\n    GROUP BY Route\n),\nranked_unique AS (\n    SELECT\n        r.*,\n        rr.route_recurrence_count,\n        row_number() OVER (\n            PARTITION BY r.Route\n            ORDER BY\n                r.hop_count DESC,\n                r.FlightDate DESC,\n                r.first_dep_time DESC,\n                r.carrier,\n                r.flight_number,\n                r.aircraft_id\n        ) AS route_pick_rank\n    FROM routes AS r\n    INNER JOIN route_recurrence AS rr\n        ON r.Route = rr.Route\n),\ntop_routes AS (\n    SELECT *\n    FROM ranked_unique\n    WHERE route_pick_rank = 1\n    ORDER BY hop_count DESC, FlightDate DESC, first_dep_time DESC, Route\n    LIMIT 10\n),\nexpanded AS (\n    SELECT\n        Route,\n        hop_count,\n        FlightDate,\n        carrier,\n        flight_number,\n        if(aircraft_id = '', 'unknown', aircraft_id) AS aircraft_id,\n        arrayConcat(\n            [tupleElement(arrayElement(legs_arr, 1), 2)],\n            arrayMap(x -> tupleElement(x, 3), legs_arr)\n        ) AS airport_codes_seq,\n        arrayConcat(\n            [tupleElement(arrayElement(legs_arr, 1), 4)],\n            arrayMap(x -> tupleElement(x, 5), legs_arr)\n        ) AS airport_ids_seq,\n        arrayConcat(\n            [tupleElement(arrayElement(legs_arr, 1), 6)],\n            arrayMap(x -> tupleElement(x, 7), legs_arr)\n        ) AS city_market_ids_seq,\n        arrayMap(x -> tupleElement(x, 8), legs_arr) AS leg_distances\n    FROM top_routes\n),\nroute_airports AS (\n    SELECT\n        e.Route,\n        e.hop_count,\n        e.FlightDate,\n        e.carrier,\n        e.flight_number,\n        e.aircraft_id,\n        e.airport_codes_seq,\n        e.airport_ids_seq,\n        e.city_market_ids_seq,\n        e.leg_distances,\n        arrayJoin(arrayZip(e.airport_codes_seq, e.airport_ids_seq)) AS ap\n    FROM expanded AS e\n)\nSELECT\n    ra.Route,\n    any(ra.aircraft_id) AS aircraft_id,\n    any(ra.flight_number) AS flight_number,\n    any(ra.carrier) AS carrier,\n    any(ra.FlightDate) AS flight_date,\n    any(ra.hop_count) AS hop_count,\n    arraySum(any(ra.leg_distances)) AS total_flown_distance,\n    length(arrayDistinct(any(ra.airport_codes_seq))) AS unique_airports,\n    length(arrayDistinct(any(ra.city_market_ids_seq))) AS unique_city_markets,\n    uniqExact(a.StateCode) AS unique_states,\n    uniqExact(a.UtcLocalTimeVariation) AS unique_local_time_offsets,\n    countIf(a.CountryCodeISO != 'US') = 0 AS entirely_domestic\nFROM route_airports AS ra\nLEFT JOIN ontime.dim_airports AS a\n    ON a.AirportID = ra.ap.2\n   AND a.IsLatest = 1\nGROUP BY ra.Route\nORDER BY total_flown_distance DESC, unique_airports DESC, unique_city_markets DESC, Route"
    }
  ]
}
```

Saved analysis.json:

```json
{
  "sql": "WITH max_dt AS (\n    SELECT max(FlightDate) AS mx\n    FROM ontime.fact_ontime\n),\nbase AS (\n    SELECT\n        f.FlightDate,\n        f.IATA_CODE_Reporting_Airline AS carrier,\n        f.Flight_Number_Reporting_Airline AS flight_number,\n        ifNull(nullIf(f.Tail_Number, ''), '') AS aircraft_id,\n        f.OriginCode AS origin_code,\n        f.DestCode AS dest_code,\n        f.OriginAirportID AS origin_airport_id,\n        f.DestAirportID AS dest_airport_id,\n        f.OriginCityMarketID AS origin_city_market_id,\n        f.DestCityMarketID AS dest_city_market_id,\n        coalesce(f.CRSDepTime, f.DepTime, 9999) AS dep_sort_time,\n        coalesce(f.Distance, 0) AS distance,\n        row_number() OVER (\n            PARTITION BY\n                f.FlightDate,\n                f.IATA_CODE_Reporting_Airline,\n                f.Flight_Number_Reporting_Airline,\n                ifNull(nullIf(f.Tail_Number, ''), ''),\n                f.OriginCode,\n                f.DestCode,\n                coalesce(f.CRSDepTime, f.DepTime, 9999)\n            ORDER BY\n                f.DepTime DESC NULLS LAST,\n                f.ArrTime DESC NULLS LAST,\n                coalesce(f.Distance, 0) DESC\n        ) AS rn\n    FROM ontime.fact_ontime AS f\n    CROSS JOIN max_dt\n    WHERE f.Cancelled = 0\n      AND f.FlightDate \u003e= addYears(mx, -5)\n),\nlegs AS (\n    SELECT *\n    FROM base\n    WHERE rn = 1\n),\nitineraries AS (\n    SELECT\n        FlightDate,\n        carrier,\n        flight_number,\n        aircraft_id,\n        arraySort(groupArray((\n            dep_sort_time,\n            origin_code,\n            dest_code,\n            origin_airport_id,\n            dest_airport_id,\n            origin_city_market_id,\n            dest_city_market_id,\n            distance\n        ))) AS legs_arr,\n        min(dep_sort_time) AS first_dep_time\n    FROM legs\n    GROUP BY FlightDate, carrier, flight_number, aircraft_id\n),\nroutes AS (\n    SELECT\n        FlightDate,\n        carrier,\n        flight_number,\n        aircraft_id,\n        first_dep_time,\n        length(legs_arr) AS hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                [tupleElement(arrayElement(legs_arr, 1), 2)],\n                arrayMap(x -\u003e tupleElement(x, 3), legs_arr)\n            ),\n            '-'\n        ) AS Route\n    FROM itineraries\n),\nroute_recurrence AS (\n    SELECT\n        Route,\n        uniqExact(FlightDate) AS route_recurrence_count\n    FROM routes\n    GROUP BY Route\n),\nranked_unique AS (\n    SELECT\n        r.FlightDate,\n        r.carrier,\n        r.flight_number,\n        r.aircraft_id,\n        r.first_dep_time,\n        r.hop_count,\n        rr.route_recurrence_count,\n        r.Route,\n        row_number() OVER (\n            PARTITION BY r.Route\n            ORDER BY\n                r.hop_count DESC,\n                r.FlightDate DESC,\n                r.first_dep_time DESC,\n                r.carrier,\n                r.flight_number,\n                r.aircraft_id\n        ) AS route_pick_rank\n    FROM routes AS r\n    INNER JOIN route_recurrence AS rr\n        ON r.Route = rr.Route\n)\nSELECT\n    if(aircraft_id = '', 'unknown', aircraft_id) AS aircraft_id,\n    flight_number,\n    carrier,\n    FlightDate AS flight_date,\n    hop_count,\n    route_recurrence_count,\n    Route\nFROM ranked_unique\nWHERE route_pick_rank = 1\nORDER BY hop_count DESC, flight_date DESC, first_dep_time DESC, Route\nLIMIT 10",
  "report_markdown": "",
  "subquestions": [
    {
      "id": "main",
      "subquestion": "Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.\nDefine uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.\nDo not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.\nCount hops from distinct same-day legs, not raw source rows; do not let duplicate or conflicting same-time rows inflate hop count or create artifact routes.\n\nReturn:\n\n- aircraft id\n- flight number\n- carrier\n- flight date\n- hop count\n- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft\n- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`",
      "answer_markdown": "Within the most recent 5 years of available data, the highest observed same-aircraft, same-flight-number itineraries reached 8 hops after deduplicating same-day leg rows. The most recent top route was WN 366 on 2024-12-01 by aircraft N957WN, flying ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA. Among the top 10 unique 8-hop routes, recurrence ranged from 1 day to 46 days, with LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN recurring most often at 46 days and MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC next at 40 days.",
      "sql": "WITH max_dt AS (\n    SELECT max(FlightDate) AS mx\n    FROM ontime.fact_ontime\n),\nbase AS (\n    SELECT\n        f.FlightDate,\n        f.IATA_CODE_Reporting_Airline AS carrier,\n        f.Flight_Number_Reporting_Airline AS flight_number,\n        ifNull(nullIf(f.Tail_Number, ''), '') AS aircraft_id,\n        f.OriginCode AS origin_code,\n        f.DestCode AS dest_code,\n        f.OriginAirportID AS origin_airport_id,\n        f.DestAirportID AS dest_airport_id,\n        f.OriginCityMarketID AS origin_city_market_id,\n        f.DestCityMarketID AS dest_city_market_id,\n        coalesce(f.CRSDepTime, f.DepTime, 9999) AS dep_sort_time,\n        coalesce(f.Distance, 0) AS distance,\n        row_number() OVER (\n            PARTITION BY\n                f.FlightDate,\n                f.IATA_CODE_Reporting_Airline,\n                f.Flight_Number_Reporting_Airline,\n                ifNull(nullIf(f.Tail_Number, ''), ''),\n                f.OriginCode,\n                f.DestCode,\n                coalesce(f.CRSDepTime, f.DepTime, 9999)\n            ORDER BY\n                f.DepTime DESC NULLS LAST,\n                f.ArrTime DESC NULLS LAST,\n                coalesce(f.Distance, 0) DESC\n        ) AS rn\n    FROM ontime.fact_ontime AS f\n    CROSS JOIN max_dt\n    WHERE f.Cancelled = 0\n      AND f.FlightDate \u003e= addYears(mx, -5)\n),\nlegs AS (\n    SELECT *\n    FROM base\n    WHERE rn = 1\n),\nitineraries AS (\n    SELECT\n        FlightDate,\n        carrier,\n        flight_number,\n        aircraft_id,\n        arraySort(groupArray((\n            dep_sort_time,\n            origin_code,\n            dest_code,\n            origin_airport_id,\n            dest_airport_id,\n            origin_city_market_id,\n            dest_city_market_id,\n            distance\n        ))) AS legs_arr,\n        min(dep_sort_time) AS first_dep_time\n    FROM legs\n    GROUP BY FlightDate, carrier, flight_number, aircraft_id\n),\nroutes AS (\n    SELECT\n        FlightDate,\n        carrier,\n        flight_number,\n        aircraft_id,\n        first_dep_time,\n        length(legs_arr) AS hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                [tupleElement(arrayElement(legs_arr, 1), 2)],\n                arrayMap(x -\u003e tupleElement(x, 3), legs_arr)\n            ),\n            '-'\n        ) AS Route\n    FROM itineraries\n),\nroute_recurrence AS (\n    SELECT\n        Route,\n        uniqExact(FlightDate) AS route_recurrence_count\n    FROM routes\n    GROUP BY Route\n),\nranked_unique AS (\n    SELECT\n        r.FlightDate,\n        r.carrier,\n        r.flight_number,\n        r.aircraft_id,\n        r.first_dep_time,\n        r.hop_count,\n        rr.route_recurrence_count,\n        r.Route,\n        row_number() OVER (\n            PARTITION BY r.Route\n            ORDER BY\n                r.hop_count DESC,\n                r.FlightDate DESC,\n                r.first_dep_time DESC,\n                r.carrier,\n                r.flight_number,\n                r.aircraft_id\n        ) AS route_pick_rank\n    FROM routes AS r\n    INNER JOIN route_recurrence AS rr\n        ON r.Route = rr.Route\n)\nSELECT\n    if(aircraft_id = '', 'unknown', aircraft_id) AS aircraft_id,\n    flight_number,\n    carrier,\n    FlightDate AS flight_date,\n    hop_count,\n    route_recurrence_count,\n    Route\nFROM ranked_unique\nWHERE route_pick_rank = 1\nORDER BY hop_count DESC, flight_date DESC, first_dep_time DESC, Route\nLIMIT 10"
    },
    {
      "id": "q1",
      "subquestion": "Which airports act as the key connectors, origins, and termini within the top 10 unique maximum-hop itineraries?\nClassify airport appearances by route position and return airport code, airport name, city/state, total appearances, origin appearances, intermediate-stop appearances, final-destination appearances, and share of itineraries containing that airport.",
      "answer_markdown": "Across the top 10 unique maximum-hop itineraries, the strongest connector airports were DAL, DEN, LAS, MSY, OAK, and BWI, each appearing in 5 of the 10 itineraries. DAL was purely an intermediate stop in all 5 appearances, while BWI and MSY were the most common origins at 2 appearances each. The most common termini were OAK and LAX with 2 final-destination appearances each. BWI, DAL, DEN, LAS, MSY, and OAK each appeared in 50% of the itineraries.",
      "sql": "WITH max_dt AS (\n    SELECT max(FlightDate) AS mx\n    FROM ontime.fact_ontime\n),\nbase AS (\n    SELECT\n        f.FlightDate,\n        f.IATA_CODE_Reporting_Airline AS carrier,\n        f.Flight_Number_Reporting_Airline AS flight_number,\n        ifNull(nullIf(f.Tail_Number, ''), '') AS aircraft_id,\n        f.OriginCode AS origin_code,\n        f.DestCode AS dest_code,\n        f.OriginAirportID AS origin_airport_id,\n        f.DestAirportID AS dest_airport_id,\n        coalesce(f.CRSDepTime, f.DepTime, 9999) AS dep_sort_time,\n        coalesce(f.Distance, 0) AS distance,\n        row_number() OVER (\n            PARTITION BY\n                f.FlightDate,\n                f.IATA_CODE_Reporting_Airline,\n                f.Flight_Number_Reporting_Airline,\n                ifNull(nullIf(f.Tail_Number, ''), ''),\n                f.OriginCode,\n                f.DestCode,\n                coalesce(f.CRSDepTime, f.DepTime, 9999)\n            ORDER BY\n                f.DepTime DESC NULLS LAST,\n                f.ArrTime DESC NULLS LAST,\n                coalesce(f.Distance, 0) DESC\n        ) AS rn\n    FROM ontime.fact_ontime AS f\n    CROSS JOIN max_dt\n    WHERE f.Cancelled = 0\n      AND f.FlightDate \u003e= addYears(mx, -5)\n),\nlegs AS (\n    SELECT *\n    FROM base\n    WHERE rn = 1\n),\nitineraries AS (\n    SELECT\n        FlightDate,\n        carrier,\n        flight_number,\n        aircraft_id,\n        arraySort(groupArray((\n            dep_sort_time,\n            origin_code,\n            dest_code,\n            origin_airport_id,\n            dest_airport_id,\n            distance\n        ))) AS legs_arr,\n        min(dep_sort_time) AS first_dep_time\n    FROM legs\n    GROUP BY FlightDate, carrier, flight_number, aircraft_id\n),\nroutes AS (\n    SELECT\n        FlightDate,\n        carrier,\n        flight_number,\n        aircraft_id,\n        first_dep_time,\n        length(legs_arr) AS hop_count,\n        legs_arr,\n        arrayStringConcat(\n            arrayConcat(\n                [tupleElement(arrayElement(legs_arr, 1), 2)],\n                arrayMap(x -\u003e tupleElement(x, 3), legs_arr)\n            ),\n            '-'\n        ) AS Route\n    FROM itineraries\n),\nroute_recurrence AS (\n    SELECT\n        Route,\n        uniqExact(FlightDate) AS route_recurrence_count\n    FROM routes\n    GROUP BY Route\n),\nranked_unique AS (\n    SELECT\n        r.*,\n        rr.route_recurrence_count,\n        row_number() OVER (\n            PARTITION BY r.Route\n            ORDER BY\n                r.hop_count DESC,\n                r.FlightDate DESC,\n                r.first_dep_time DESC,\n                r.carrier,\n                r.flight_number,\n                r.aircraft_id\n        ) AS route_pick_rank\n    FROM routes AS r\n    INNER JOIN route_recurrence AS rr\n        ON r.Route = rr.Route\n),\ntop_routes AS (\n    SELECT *\n    FROM ranked_unique\n    WHERE route_pick_rank = 1\n    ORDER BY hop_count DESC, FlightDate DESC, first_dep_time DESC, Route\n    LIMIT 10\n),\nroute_airports AS (\n    SELECT\n        Route,\n        arrayJoin(\n            arrayZip(\n                arrayConcat(\n                    [tupleElement(arrayElement(legs_arr, 1), 2)],\n                    arrayMap(x -\u003e tupleElement(x, 3), legs_arr)\n                ),\n                arrayEnumerate(\n                    arrayConcat(\n                        [tupleElement(arrayElement(legs_arr, 1), 2)],\n                        arrayMap(x -\u003e tupleElement(x, 3), legs_arr)\n                    )\n                ),\n                arrayMap(\n                    x -\u003e length(legs_arr) + 1,\n                    arrayConcat(\n                        [tupleElement(arrayElement(legs_arr, 1), 2)],\n                        arrayMap(x -\u003e tupleElement(x, 3), legs_arr)\n                    )\n                )\n            )\n        ) AS airport_pos\n    FROM top_routes\n)\nSELECT\n    trimBoth(toString(airport_pos.1)) AS airport_code,\n    any(a.DisplayAirportName) AS airport_name,\n    any(concat(a.CityName, ', ', a.StateCode)) AS city_state,\n    count() AS total_appearances,\n    countIf(airport_pos.2 = 1) AS origin_appearances,\n    countIf(airport_pos.2 \u003e 1 AND airport_pos.2 \u003c airport_pos.3) AS intermediate_stop_appearances,\n    countIf(airport_pos.2 = airport_pos.3) AS final_destination_appearances,\n    round(uniqExact(Route) / 10.0, 3) AS share_of_itineraries_containing_airport\nFROM route_airports AS ra\nLEFT JOIN ontime.dim_airports AS a\n    ON a.AirportCode = ra.airport_pos.1\n   AND a.IsLatest = 1\nGROUP BY airport_code\nORDER BY total_appearances DESC, share_of_itineraries_containing_airport DESC, airport_code"
    },
    {
      "id": "q2",
      "subquestion": "How geographically extreme is each of the top 10 unique maximum-hop itineraries?\nReturn total flown distance, unique airports, unique city markets, unique states, unique local-time offsets, and whether the route is entirely domestic, then summarize which routes are the most geographically expansive.",
      "answer_markdown": "All top 10 unique maximum-hop itineraries were entirely domestic. The most geographically expansive by total flown distance was BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK at 5,169 miles, spanning 9 airports, 8 city markets, 9 states, and 3 local-time offsets. Several routes were broader in time-zone spread, with CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN, ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN, HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK, MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC, and BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK each touching 4 distinct local-time offsets.",
      "sql": "WITH max_dt AS (\n    SELECT max(FlightDate) AS mx\n    FROM ontime.fact_ontime\n),\nbase AS (\n    SELECT\n        f.FlightDate,\n        f.IATA_CODE_Reporting_Airline AS carrier,\n        f.Flight_Number_Reporting_Airline AS flight_number,\n        ifNull(nullIf(f.Tail_Number, ''), '') AS aircraft_id,\n        f.OriginCode AS origin_code,\n        f.DestCode AS dest_code,\n        f.OriginAirportID AS origin_airport_id,\n        f.DestAirportID AS dest_airport_id,\n        f.OriginCityMarketID AS origin_city_market_id,\n        f.DestCityMarketID AS dest_city_market_id,\n        coalesce(f.CRSDepTime, f.DepTime, 9999) AS dep_sort_time,\n        coalesce(f.Distance, 0) AS distance,\n        row_number() OVER (\n            PARTITION BY\n                f.FlightDate,\n                f.IATA_CODE_Reporting_Airline,\n                f.Flight_Number_Reporting_Airline,\n                ifNull(nullIf(f.Tail_Number, ''), ''),\n                f.OriginCode,\n                f.DestCode,\n                coalesce(f.CRSDepTime, f.DepTime, 9999)\n            ORDER BY\n                f.DepTime DESC NULLS LAST,\n                f.ArrTime DESC NULLS LAST,\n                coalesce(f.Distance, 0) DESC\n        ) AS rn\n    FROM ontime.fact_ontime AS f\n    CROSS JOIN max_dt\n    WHERE f.Cancelled = 0\n      AND f.FlightDate \u003e= addYears(mx, -5)\n),\nlegs AS (\n    SELECT *\n    FROM base\n    WHERE rn = 1\n),\nitineraries AS (\n    SELECT\n        FlightDate,\n        carrier,\n        flight_number,\n        aircraft_id,\n        arraySort(groupArray((\n            dep_sort_time,\n            origin_code,\n            dest_code,\n            origin_airport_id,\n            dest_airport_id,\n            origin_city_market_id,\n            dest_city_market_id,\n            distance\n        ))) AS legs_arr,\n        min(dep_sort_time) AS first_dep_time\n    FROM legs\n    GROUP BY FlightDate, carrier, flight_number, aircraft_id\n),\nroutes AS (\n    SELECT\n        FlightDate,\n        carrier,\n        flight_number,\n        aircraft_id,\n        first_dep_time,\n        length(legs_arr) AS hop_count,\n        legs_arr,\n        arrayStringConcat(\n            arrayConcat(\n                [tupleElement(arrayElement(legs_arr, 1), 2)],\n                arrayMap(x -\u003e tupleElement(x, 3), legs_arr)\n            ),\n            '-'\n        ) AS Route\n    FROM itineraries\n),\nroute_recurrence AS (\n    SELECT\n        Route,\n        uniqExact(FlightDate) AS route_recurrence_count\n    FROM routes\n    GROUP BY Route\n),\nranked_unique AS (\n    SELECT\n        r.*,\n        rr.route_recurrence_count,\n        row_number() OVER (\n            PARTITION BY r.Route\n            ORDER BY\n                r.hop_count DESC,\n                r.FlightDate DESC,\n                r.first_dep_time DESC,\n                r.carrier,\n                r.flight_number,\n                r.aircraft_id\n        ) AS route_pick_rank\n    FROM routes AS r\n    INNER JOIN route_recurrence AS rr\n        ON r.Route = rr.Route\n),\ntop_routes AS (\n    SELECT *\n    FROM ranked_unique\n    WHERE route_pick_rank = 1\n    ORDER BY hop_count DESC, FlightDate DESC, first_dep_time DESC, Route\n    LIMIT 10\n),\nexpanded AS (\n    SELECT\n        Route,\n        hop_count,\n        FlightDate,\n        carrier,\n        flight_number,\n        if(aircraft_id = '', 'unknown', aircraft_id) AS aircraft_id,\n        arrayConcat(\n            [tupleElement(arrayElement(legs_arr, 1), 2)],\n            arrayMap(x -\u003e tupleElement(x, 3), legs_arr)\n        ) AS airport_codes_seq,\n        arrayConcat(\n            [tupleElement(arrayElement(legs_arr, 1), 4)],\n            arrayMap(x -\u003e tupleElement(x, 5), legs_arr)\n        ) AS airport_ids_seq,\n        arrayConcat(\n            [tupleElement(arrayElement(legs_arr, 1), 6)],\n            arrayMap(x -\u003e tupleElement(x, 7), legs_arr)\n        ) AS city_market_ids_seq,\n        arrayMap(x -\u003e tupleElement(x, 8), legs_arr) AS leg_distances\n    FROM top_routes\n),\nroute_airports AS (\n    SELECT\n        e.Route,\n        e.hop_count,\n        e.FlightDate,\n        e.carrier,\n        e.flight_number,\n        e.aircraft_id,\n        e.airport_codes_seq,\n        e.airport_ids_seq,\n        e.city_market_ids_seq,\n        e.leg_distances,\n        arrayJoin(arrayZip(e.airport_codes_seq, e.airport_ids_seq)) AS ap\n    FROM expanded AS e\n)\nSELECT\n    ra.Route,\n    any(ra.aircraft_id) AS aircraft_id,\n    any(ra.flight_number) AS flight_number,\n    any(ra.carrier) AS carrier,\n    any(ra.FlightDate) AS flight_date,\n    any(ra.hop_count) AS hop_count,\n    arraySum(any(ra.leg_distances)) AS total_flown_distance,\n    length(arrayDistinct(any(ra.airport_codes_seq))) AS unique_airports,\n    length(arrayDistinct(any(ra.city_market_ids_seq))) AS unique_city_markets,\n    uniqExact(a.StateCode) AS unique_states,\n    uniqExact(a.UtcLocalTimeVariation) AS unique_local_time_offsets,\n    countIf(a.CountryCodeISO != 'US') = 0 AS entirely_domestic\nFROM route_airports AS ra\nLEFT JOIN ontime.dim_airports AS a\n    ON a.AirportID = ra.ap.2\n   AND a.IsLatest = 1\nGROUP BY ra.Route\nORDER BY total_flown_distance DESC, unique_airports DESC, unique_city_markets DESC, Route"
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
      "subquestion": "Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.\nDefine uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.\nDo not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.\nCount hops from distinct same-day legs, not raw source rows; do not let duplicate or conflicting same-time rows inflate hop count or create artifact routes.\n\nReturn:\n\n- aircraft id\n- flight number\n- carrier\n- flight date\n- hop count\n- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft\n- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`",
      "answer_markdown": "Within the most recent 5 years of available data, the highest observed same-aircraft, same-flight-number itineraries reached 8 hops after deduplicating same-day leg rows. The most recent top route was WN 366 on 2024-12-01 by aircraft N957WN, flying ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA. Among the top 10 unique 8-hop routes, recurrence ranged from 1 day to 46 days, with LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN recurring most often at 46 days and MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC next at 40 days.",
      "sql": "WITH max_dt AS (\n    SELECT max(FlightDate) AS mx\n    FROM ontime.fact_ontime\n),\nbase AS (\n    SELECT\n        f.FlightDate,\n        f.IATA_CODE_Reporting_Airline AS carrier,\n        f.Flight_Number_Reporting_Airline AS flight_number,\n        ifNull(nullIf(f.Tail_Number, ''), '') AS aircraft_id,\n        f.OriginCode AS origin_code,\n        f.DestCode AS dest_code,\n        f.OriginAirportID AS origin_airport_id,\n        f.DestAirportID AS dest_airport_id,\n        f.OriginCityMarketID AS origin_city_market_id,\n        f.DestCityMarketID AS dest_city_market_id,\n        coalesce(f.CRSDepTime, f.DepTime, 9999) AS dep_sort_time,\n        coalesce(f.Distance, 0) AS distance,\n        row_number() OVER (\n            PARTITION BY\n                f.FlightDate,\n                f.IATA_CODE_Reporting_Airline,\n                f.Flight_Number_Reporting_Airline,\n                ifNull(nullIf(f.Tail_Number, ''), ''),\n                f.OriginCode,\n                f.DestCode,\n                coalesce(f.CRSDepTime, f.DepTime, 9999)\n            ORDER BY\n                f.DepTime DESC NULLS LAST,\n                f.ArrTime DESC NULLS LAST,\n                coalesce(f.Distance, 0) DESC\n        ) AS rn\n    FROM ontime.fact_ontime AS f\n    CROSS JOIN max_dt\n    WHERE f.Cancelled = 0\n      AND f.FlightDate \u003e= addYears(mx, -5)\n),\nlegs AS (\n    SELECT *\n    FROM base\n    WHERE rn = 1\n),\nitineraries AS (\n    SELECT\n        FlightDate,\n        carrier,\n        flight_number,\n        aircraft_id,\n        arraySort(groupArray((\n            dep_sort_time,\n            origin_code,\n            dest_code,\n            origin_airport_id,\n            dest_airport_id,\n            origin_city_market_id,\n            dest_city_market_id,\n            distance\n        ))) AS legs_arr,\n        min(dep_sort_time) AS first_dep_time\n    FROM legs\n    GROUP BY FlightDate, carrier, flight_number, aircraft_id\n),\nroutes AS (\n    SELECT\n        FlightDate,\n        carrier,\n        flight_number,\n        aircraft_id,\n        first_dep_time,\n        length(legs_arr) AS hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                [tupleElement(arrayElement(legs_arr, 1), 2)],\n                arrayMap(x -\u003e tupleElement(x, 3), legs_arr)\n            ),\n            '-'\n        ) AS Route\n    FROM itineraries\n),\nroute_recurrence AS (\n    SELECT\n        Route,\n        uniqExact(FlightDate) AS route_recurrence_count\n    FROM routes\n    GROUP BY Route\n),\nranked_unique AS (\n    SELECT\n        r.FlightDate,\n        r.carrier,\n        r.flight_number,\n        r.aircraft_id,\n        r.first_dep_time,\n        r.hop_count,\n        rr.route_recurrence_count,\n        r.Route,\n        row_number() OVER (\n            PARTITION BY r.Route\n            ORDER BY\n                r.hop_count DESC,\n                r.FlightDate DESC,\n                r.first_dep_time DESC,\n                r.carrier,\n                r.flight_number,\n                r.aircraft_id\n        ) AS route_pick_rank\n    FROM routes AS r\n    INNER JOIN route_recurrence AS rr\n        ON r.Route = rr.Route\n)\nSELECT\n    if(aircraft_id = '', 'unknown', aircraft_id) AS aircraft_id,\n    flight_number,\n    carrier,\n    FlightDate AS flight_date,\n    hop_count,\n    route_recurrence_count,\n    Route\nFROM ranked_unique\nWHERE route_pick_rank = 1\nORDER BY hop_count DESC, flight_date DESC, first_dep_time DESC, Route\nLIMIT 10",
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
    },
    {
      "id": "q1",
      "subquestion": "Which airports act as the key connectors, origins, and termini within the top 10 unique maximum-hop itineraries?\nClassify airport appearances by route position and return airport code, airport name, city/state, total appearances, origin appearances, intermediate-stop appearances, final-destination appearances, and share of itineraries containing that airport.",
      "answer_markdown": "Across the top 10 unique maximum-hop itineraries, the strongest connector airports were DAL, DEN, LAS, MSY, OAK, and BWI, each appearing in 5 of the 10 itineraries. DAL was purely an intermediate stop in all 5 appearances, while BWI and MSY were the most common origins at 2 appearances each. The most common termini were OAK and LAX with 2 final-destination appearances each. BWI, DAL, DEN, LAS, MSY, and OAK each appeared in 50% of the itineraries.",
      "sql": "WITH max_dt AS (\n    SELECT max(FlightDate) AS mx\n    FROM ontime.fact_ontime\n),\nbase AS (\n    SELECT\n        f.FlightDate,\n        f.IATA_CODE_Reporting_Airline AS carrier,\n        f.Flight_Number_Reporting_Airline AS flight_number,\n        ifNull(nullIf(f.Tail_Number, ''), '') AS aircraft_id,\n        f.OriginCode AS origin_code,\n        f.DestCode AS dest_code,\n        f.OriginAirportID AS origin_airport_id,\n        f.DestAirportID AS dest_airport_id,\n        coalesce(f.CRSDepTime, f.DepTime, 9999) AS dep_sort_time,\n        coalesce(f.Distance, 0) AS distance,\n        row_number() OVER (\n            PARTITION BY\n                f.FlightDate,\n                f.IATA_CODE_Reporting_Airline,\n                f.Flight_Number_Reporting_Airline,\n                ifNull(nullIf(f.Tail_Number, ''), ''),\n                f.OriginCode,\n                f.DestCode,\n                coalesce(f.CRSDepTime, f.DepTime, 9999)\n            ORDER BY\n                f.DepTime DESC NULLS LAST,\n                f.ArrTime DESC NULLS LAST,\n                coalesce(f.Distance, 0) DESC\n        ) AS rn\n    FROM ontime.fact_ontime AS f\n    CROSS JOIN max_dt\n    WHERE f.Cancelled = 0\n      AND f.FlightDate \u003e= addYears(mx, -5)\n),\nlegs AS (\n    SELECT *\n    FROM base\n    WHERE rn = 1\n),\nitineraries AS (\n    SELECT\n        FlightDate,\n        carrier,\n        flight_number,\n        aircraft_id,\n        arraySort(groupArray((\n            dep_sort_time,\n            origin_code,\n            dest_code,\n            origin_airport_id,\n            dest_airport_id,\n            distance\n        ))) AS legs_arr,\n        min(dep_sort_time) AS first_dep_time\n    FROM legs\n    GROUP BY FlightDate, carrier, flight_number, aircraft_id\n),\nroutes AS (\n    SELECT\n        FlightDate,\n        carrier,\n        flight_number,\n        aircraft_id,\n        first_dep_time,\n        length(legs_arr) AS hop_count,\n        legs_arr,\n        arrayStringConcat(\n            arrayConcat(\n                [tupleElement(arrayElement(legs_arr, 1), 2)],\n                arrayMap(x -\u003e tupleElement(x, 3), legs_arr)\n            ),\n            '-'\n        ) AS Route\n    FROM itineraries\n),\nroute_recurrence AS (\n    SELECT\n        Route,\n        uniqExact(FlightDate) AS route_recurrence_count\n    FROM routes\n    GROUP BY Route\n),\nranked_unique AS (\n    SELECT\n        r.*,\n        rr.route_recurrence_count,\n        row_number() OVER (\n            PARTITION BY r.Route\n            ORDER BY\n                r.hop_count DESC,\n                r.FlightDate DESC,\n                r.first_dep_time DESC,\n                r.carrier,\n                r.flight_number,\n                r.aircraft_id\n        ) AS route_pick_rank\n    FROM routes AS r\n    INNER JOIN route_recurrence AS rr\n        ON r.Route = rr.Route\n),\ntop_routes AS (\n    SELECT *\n    FROM ranked_unique\n    WHERE route_pick_rank = 1\n    ORDER BY hop_count DESC, FlightDate DESC, first_dep_time DESC, Route\n    LIMIT 10\n),\nroute_airports AS (\n    SELECT\n        Route,\n        arrayJoin(\n            arrayZip(\n                arrayConcat(\n                    [tupleElement(arrayElement(legs_arr, 1), 2)],\n                    arrayMap(x -\u003e tupleElement(x, 3), legs_arr)\n                ),\n                arrayEnumerate(\n                    arrayConcat(\n                        [tupleElement(arrayElement(legs_arr, 1), 2)],\n                        arrayMap(x -\u003e tupleElement(x, 3), legs_arr)\n                    )\n                ),\n                arrayMap(\n                    x -\u003e length(legs_arr) + 1,\n                    arrayConcat(\n                        [tupleElement(arrayElement(legs_arr, 1), 2)],\n                        arrayMap(x -\u003e tupleElement(x, 3), legs_arr)\n                    )\n                )\n            )\n        ) AS airport_pos\n    FROM top_routes\n)\nSELECT\n    trimBoth(toString(airport_pos.1)) AS airport_code,\n    any(a.DisplayAirportName) AS airport_name,\n    any(concat(a.CityName, ', ', a.StateCode)) AS city_state,\n    count() AS total_appearances,\n    countIf(airport_pos.2 = 1) AS origin_appearances,\n    countIf(airport_pos.2 \u003e 1 AND airport_pos.2 \u003c airport_pos.3) AS intermediate_stop_appearances,\n    countIf(airport_pos.2 = airport_pos.3) AS final_destination_appearances,\n    round(uniqExact(Route) / 10.0, 3) AS share_of_itineraries_containing_airport\nFROM route_airports AS ra\nLEFT JOIN ontime.dim_airports AS a\n    ON a.AirportCode = ra.airport_pos.1\n   AND a.IsLatest = 1\nGROUP BY airport_code\nORDER BY total_appearances DESC, share_of_itineraries_containing_airport DESC, airport_code",
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
        "airport_code": "BWI",
        "airport_name": "Baltimore/Washington International Thurgood Marshall",
        "city_state": "Baltimore, MD, MD",
        "final_destination_appearances": 0,
        "intermediate_stop_appearances": 3,
        "origin_appearances": 2,
        "share_of_itineraries_containing_airport": 0.5,
        "total_appearances": 5
      }
    },
    {
      "id": "q2",
      "subquestion": "How geographically extreme is each of the top 10 unique maximum-hop itineraries?\nReturn total flown distance, unique airports, unique city markets, unique states, unique local-time offsets, and whether the route is entirely domestic, then summarize which routes are the most geographically expansive.",
      "answer_markdown": "All top 10 unique maximum-hop itineraries were entirely domestic. The most geographically expansive by total flown distance was BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK at 5,169 miles, spanning 9 airports, 8 city markets, 9 states, and 3 local-time offsets. Several routes were broader in time-zone spread, with CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN, ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN, HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK, MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC, and BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK each touching 4 distinct local-time offsets.",
      "sql": "WITH max_dt AS (\n    SELECT max(FlightDate) AS mx\n    FROM ontime.fact_ontime\n),\nbase AS (\n    SELECT\n        f.FlightDate,\n        f.IATA_CODE_Reporting_Airline AS carrier,\n        f.Flight_Number_Reporting_Airline AS flight_number,\n        ifNull(nullIf(f.Tail_Number, ''), '') AS aircraft_id,\n        f.OriginCode AS origin_code,\n        f.DestCode AS dest_code,\n        f.OriginAirportID AS origin_airport_id,\n        f.DestAirportID AS dest_airport_id,\n        f.OriginCityMarketID AS origin_city_market_id,\n        f.DestCityMarketID AS dest_city_market_id,\n        coalesce(f.CRSDepTime, f.DepTime, 9999) AS dep_sort_time,\n        coalesce(f.Distance, 0) AS distance,\n        row_number() OVER (\n            PARTITION BY\n                f.FlightDate,\n                f.IATA_CODE_Reporting_Airline,\n                f.Flight_Number_Reporting_Airline,\n                ifNull(nullIf(f.Tail_Number, ''), ''),\n                f.OriginCode,\n                f.DestCode,\n                coalesce(f.CRSDepTime, f.DepTime, 9999)\n            ORDER BY\n                f.DepTime DESC NULLS LAST,\n                f.ArrTime DESC NULLS LAST,\n                coalesce(f.Distance, 0) DESC\n        ) AS rn\n    FROM ontime.fact_ontime AS f\n    CROSS JOIN max_dt\n    WHERE f.Cancelled = 0\n      AND f.FlightDate \u003e= addYears(mx, -5)\n),\nlegs AS (\n    SELECT *\n    FROM base\n    WHERE rn = 1\n),\nitineraries AS (\n    SELECT\n        FlightDate,\n        carrier,\n        flight_number,\n        aircraft_id,\n        arraySort(groupArray((\n            dep_sort_time,\n            origin_code,\n            dest_code,\n            origin_airport_id,\n            dest_airport_id,\n            origin_city_market_id,\n            dest_city_market_id,\n            distance\n        ))) AS legs_arr,\n        min(dep_sort_time) AS first_dep_time\n    FROM legs\n    GROUP BY FlightDate, carrier, flight_number, aircraft_id\n),\nroutes AS (\n    SELECT\n        FlightDate,\n        carrier,\n        flight_number,\n        aircraft_id,\n        first_dep_time,\n        length(legs_arr) AS hop_count,\n        legs_arr,\n        arrayStringConcat(\n            arrayConcat(\n                [tupleElement(arrayElement(legs_arr, 1), 2)],\n                arrayMap(x -\u003e tupleElement(x, 3), legs_arr)\n            ),\n            '-'\n        ) AS Route\n    FROM itineraries\n),\nroute_recurrence AS (\n    SELECT\n        Route,\n        uniqExact(FlightDate) AS route_recurrence_count\n    FROM routes\n    GROUP BY Route\n),\nranked_unique AS (\n    SELECT\n        r.*,\n        rr.route_recurrence_count,\n        row_number() OVER (\n            PARTITION BY r.Route\n            ORDER BY\n                r.hop_count DESC,\n                r.FlightDate DESC,\n                r.first_dep_time DESC,\n                r.carrier,\n                r.flight_number,\n                r.aircraft_id\n        ) AS route_pick_rank\n    FROM routes AS r\n    INNER JOIN route_recurrence AS rr\n        ON r.Route = rr.Route\n),\ntop_routes AS (\n    SELECT *\n    FROM ranked_unique\n    WHERE route_pick_rank = 1\n    ORDER BY hop_count DESC, FlightDate DESC, first_dep_time DESC, Route\n    LIMIT 10\n),\nexpanded AS (\n    SELECT\n        Route,\n        hop_count,\n        FlightDate,\n        carrier,\n        flight_number,\n        if(aircraft_id = '', 'unknown', aircraft_id) AS aircraft_id,\n        arrayConcat(\n            [tupleElement(arrayElement(legs_arr, 1), 2)],\n            arrayMap(x -\u003e tupleElement(x, 3), legs_arr)\n        ) AS airport_codes_seq,\n        arrayConcat(\n            [tupleElement(arrayElement(legs_arr, 1), 4)],\n            arrayMap(x -\u003e tupleElement(x, 5), legs_arr)\n        ) AS airport_ids_seq,\n        arrayConcat(\n            [tupleElement(arrayElement(legs_arr, 1), 6)],\n            arrayMap(x -\u003e tupleElement(x, 7), legs_arr)\n        ) AS city_market_ids_seq,\n        arrayMap(x -\u003e tupleElement(x, 8), legs_arr) AS leg_distances\n    FROM top_routes\n),\nroute_airports AS (\n    SELECT\n        e.Route,\n        e.hop_count,\n        e.FlightDate,\n        e.carrier,\n        e.flight_number,\n        e.aircraft_id,\n        e.airport_codes_seq,\n        e.airport_ids_seq,\n        e.city_market_ids_seq,\n        e.leg_distances,\n        arrayJoin(arrayZip(e.airport_codes_seq, e.airport_ids_seq)) AS ap\n    FROM expanded AS e\n)\nSELECT\n    ra.Route,\n    any(ra.aircraft_id) AS aircraft_id,\n    any(ra.flight_number) AS flight_number,\n    any(ra.carrier) AS carrier,\n    any(ra.FlightDate) AS flight_date,\n    any(ra.hop_count) AS hop_count,\n    arraySum(any(ra.leg_distances)) AS total_flown_distance,\n    length(arrayDistinct(any(ra.airport_codes_seq))) AS unique_airports,\n    length(arrayDistinct(any(ra.city_market_ids_seq))) AS unique_city_markets,\n    uniqExact(a.StateCode) AS unique_states,\n    uniqExact(a.UtcLocalTimeVariation) AS unique_local_time_offsets,\n    countIf(a.CountryCodeISO != 'US') = 0 AS entirely_domestic\nFROM route_airports AS ra\nLEFT JOIN ontime.dim_airports AS a\n    ON a.AirportID = ra.ap.2\n   AND a.IsLatest = 1\nGROUP BY ra.Route\nORDER BY total_flown_distance DESC, unique_airports DESC, unique_city_markets DESC, Route",
      "date_field_hint": "flight_date",
      "row_count": 10,
      "result_columns": [
        "Route",
        "aircraft_id",
        "flight_number",
        "carrier",
        "flight_date",
        "hop_count",
        "total_flown_distance",
        "unique_airports",
        "unique_city_markets",
        "unique_states",
        "unique_local_time_offsets",
        "entirely_domestic"
      ],
      "first_row": {
        "Route": "BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX",
        "aircraft_id": "N225WN",
        "carrier": "WN",
        "entirely_domestic": 1,
        "flight_date": "2021-08-08T00:00:00Z",
        "flight_number": "3530",
        "hop_count": 8,
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