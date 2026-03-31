- Connect to clickhouse server though MCP connection
- Do not use direct HTTP by any tools like curl.
- Use the `ontime` database to answer analytical questions
- Use `ontime-semantic-layer` skill for schema inspection, join guidance, and dimension semantics.
- write correct and efficient ClickHouse SQL 
- Before writing any SQL artifact, self-verify every SQL statement you intend to save.
- Run a cheap debug execution for each query first, usually with a small `LIMIT`, a narrow `WHERE` filter, or both applied inside the main data-reading subquery or CTE.
- Treat successful execution as mandatory. Fix any syntax, type, aggregate, window, join, or unknown-column errors in a loop until every saved query runs successfully.
- Do not write unchecked SQL.

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
Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?

List the recurrence count for each of the 10 routes explicitly before summarizing any tiers or categories, and make sure any category totals add up to 10.

### q2
What geographic pattern do the top itineraries show?

Base the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.

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

The maximum observed same-aircraft, same-flight-number itinerary length is 8 hops. The 10 most recent unique max-hop routes are all Southwest (WN) itineraries, with the newest on 2024-12-01: ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA flown by N957WN on flight 366. Across these 10 rows, route recurrence ranges from 1 day to 46 days, led by LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN (46), MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC (40), HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK (20), and BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK (12).

- Rows returned: 10
- Columns: aircraft_id, flight_number, carrier, flight_date, hop_count, route_recurrence_count, route

| aircraft_id | flight_number | carrier | flight_date | hop_count | route_recurrence_count | route |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |

> Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?

List the recurrence count for each of the 10 routes explicitly before summarizing any tiers or categories, and make sure any category totals add up to 10.

Recurrence counts for the 10 routes are: ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA = 1; CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN = 4; ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN = 2; MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX = 5; MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC = 40; LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN = 46; SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS = 7; HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK = 20; BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK = 12; BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX = 5. Overall, these top itineraries skew toward recurring patterns rather than pure one-offs: 4 are recurring (10+ days), 5 are occasional (2-9 days), and 1 is a one-off, which sums to all 10 routes.

- Rows returned: 10
- Columns: route, route_recurrence_count, recurrence_tier, tier_count

| route | route_recurrence_count | recurrence_tier | tier_count |
| --- | --- | --- | --- |
| ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | 1 | one-off | 1 |

> What geographic pattern do the top itineraries show?

Base the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.

The geographic pattern is entirely domestic U.S. and strongly tied to Southwest/West Coast plus Sun Belt and Mid-Atlantic nodes. Every airport in the 10 routes is a US airport, and the most repeated points are BWI, DAL, DEN, LAS, MSY, and OAK with 5 mentions each, followed by BNA and PHX with 4 each. California is especially prominent through OAK, BUR, LAX, SAN, PSP, SJC, and SMF, so these longest same-flight-number days look like long domestic zig-zag chains linking eastern and southeastern stops into western endpoints and western connection banks.

- Rows returned: 45
- Columns: airport_code, state_code, country_code, geo_bucket, airport_mentions

| airport_code | state_code | country_code | geo_bucket | airport_mentions |
| --- | --- | --- | --- | --- |
| BWI | MD | US | east | 5 |
```

Saved answer.raw.json:

```json
{
  "subquestions": [
    {
      "id": "main",
      "answer_markdown": "The maximum observed same-aircraft, same-flight-number itinerary length is 8 hops. The 10 most recent unique max-hop routes are all Southwest (WN) itineraries, with the newest on 2024-12-01: ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA flown by N957WN on flight 366. Across these 10 rows, route recurrence ranges from 1 day to 46 days, led by LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN (46), MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC (40), HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK (20), and BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK (12).",
      "sql": "WITH dedup_legs AS (\n    SELECT\n        FlightDate,\n        ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,\n        Flight_Number_Reporting_Airline AS flight_number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginCode,\n        DestCode,\n        coalesce(CRSDepTime, DepTime, toUInt16(9999)) AS dep_sort,\n        min(coalesce(CRSArrTime, ArrTime, toUInt16(9999))) AS arr_sort\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier, OriginCode, DestCode, dep_sort\n), itineraries AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        count() AS hop_count,\n        min(dep_sort) AS first_dep_sort,\n        arraySort(groupArray((dep_sort, arr_sort, OriginCode, DestCode))) AS legs\n    FROM dedup_legs\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier\n    HAVING count() > 1\n), route_recurrence AS (\n    SELECT\n        arrayStringConcat(arrayConcat(arrayMap(x -> x.3, legs), [arrayElement(legs, length(legs)).4]), '-') AS route,\n        uniqExact(FlightDate) AS route_recurrence_count\n    FROM itineraries\n    GROUP BY route\n), max_hops AS (\n    SELECT max(hop_count) AS max_hop_count FROM itineraries\n), top_unique AS (\n    SELECT\n        i.FlightDate,\n        i.aircraft_id,\n        i.flight_number,\n        i.carrier,\n        i.hop_count,\n        i.first_dep_sort,\n        arrayStringConcat(arrayConcat(arrayMap(x -> x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') AS route,\n        rr.route_recurrence_count,\n        row_number() OVER (\n            PARTITION BY arrayStringConcat(arrayConcat(arrayMap(x -> x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-')\n            ORDER BY i.FlightDate DESC, i.first_dep_sort DESC, i.carrier DESC, i.flight_number DESC, i.aircraft_id DESC\n        ) AS route_rank\n    FROM itineraries i\n    CROSS JOIN max_hops m\n    INNER JOIN route_recurrence rr\n        ON arrayStringConcat(arrayConcat(arrayMap(x -> x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') = rr.route\n    WHERE i.hop_count = m.max_hop_count\n)\nSELECT\n    if(aircraft_id = '', 'unknown', aircraft_id) AS aircraft_id,\n    flight_number,\n    carrier,\n    FlightDate AS flight_date,\n    hop_count,\n    route_recurrence_count,\n    route\nFROM top_unique\nWHERE route_rank = 1\nORDER BY flight_date DESC, first_dep_sort DESC, carrier DESC, flight_number DESC, aircraft_id DESC\nLIMIT 10"
    },
    {
      "id": "q1",
      "answer_markdown": "Recurrence counts for the 10 routes are: ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA = 1; CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN = 4; ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN = 2; MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX = 5; MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC = 40; LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN = 46; SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS = 7; HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK = 20; BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK = 12; BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX = 5. Overall, these top itineraries skew toward recurring patterns rather than pure one-offs: 4 are recurring (10+ days), 5 are occasional (2-9 days), and 1 is a one-off, which sums to all 10 routes.",
      "sql": "WITH main_top10 AS (\n    WITH dedup_legs AS (\n        SELECT\n            FlightDate,\n            ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,\n            Flight_Number_Reporting_Airline AS flight_number,\n            IATA_CODE_Reporting_Airline AS carrier,\n            OriginCode,\n            DestCode,\n            coalesce(CRSDepTime, DepTime, toUInt16(9999)) AS dep_sort,\n            min(coalesce(CRSArrTime, ArrTime, toUInt16(9999))) AS arr_sort\n        FROM ontime.fact_ontime\n        WHERE Cancelled = 0\n        GROUP BY FlightDate, aircraft_id, flight_number, carrier, OriginCode, DestCode, dep_sort\n    ), itineraries AS (\n        SELECT\n            FlightDate,\n            aircraft_id,\n            flight_number,\n            carrier,\n            count() AS hop_count,\n            min(dep_sort) AS first_dep_sort,\n            arraySort(groupArray((dep_sort, arr_sort, OriginCode, DestCode))) AS legs\n        FROM dedup_legs\n        GROUP BY FlightDate, aircraft_id, flight_number, carrier\n        HAVING count() > 1\n    ), route_recurrence AS (\n        SELECT\n            arrayStringConcat(arrayConcat(arrayMap(x -> x.3, legs), [arrayElement(legs, length(legs)).4]), '-') AS route,\n            uniqExact(FlightDate) AS route_recurrence_count\n        FROM itineraries\n        GROUP BY route\n    ), max_hops AS (\n        SELECT max(hop_count) AS max_hop_count FROM itineraries\n    ), top_unique AS (\n        SELECT\n            i.FlightDate,\n            i.first_dep_sort,\n            arrayStringConcat(arrayConcat(arrayMap(x -> x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') AS route,\n            rr.route_recurrence_count,\n            row_number() OVER (\n                PARTITION BY arrayStringConcat(arrayConcat(arrayMap(x -> x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-')\n                ORDER BY i.FlightDate DESC, i.first_dep_sort DESC, i.carrier DESC, i.flight_number DESC, i.aircraft_id DESC\n            ) AS route_rank\n        FROM itineraries i\n        CROSS JOIN max_hops m\n        INNER JOIN route_recurrence rr\n            ON arrayStringConcat(arrayConcat(arrayMap(x -> x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') = rr.route\n        WHERE i.hop_count = m.max_hop_count\n    )\n    SELECT route, route_recurrence_count, FlightDate, first_dep_sort\n    FROM top_unique\n    WHERE route_rank = 1\n    ORDER BY FlightDate DESC, first_dep_sort DESC\n    LIMIT 10\n)\nSELECT\n    route,\n    route_recurrence_count,\n    multiIf(route_recurrence_count = 1, 'one-off', route_recurrence_count < 10, 'occasional', 'recurring') AS recurrence_tier,\n    count() OVER (PARTITION BY multiIf(route_recurrence_count = 1, 'one-off', route_recurrence_count < 10, 'occasional', 'recurring')) AS tier_count\nFROM main_top10\nORDER BY FlightDate DESC, first_dep_sort DESC"
    },
    {
      "id": "q2",
      "answer_markdown": "The geographic pattern is entirely domestic U.S. and strongly tied to Southwest/West Coast plus Sun Belt and Mid-Atlantic nodes. Every airport in the 10 routes is a US airport, and the most repeated points are BWI, DAL, DEN, LAS, MSY, and OAK with 5 mentions each, followed by BNA and PHX with 4 each. California is especially prominent through OAK, BUR, LAX, SAN, PSP, SJC, and SMF, so these longest same-flight-number days look like long domestic zig-zag chains linking eastern and southeastern stops into western endpoints and western connection banks.",
      "sql": "WITH dedup_legs AS (\n    SELECT\n        FlightDate,\n        ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,\n        Flight_Number_Reporting_Airline AS flight_number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginCode,\n        DestCode,\n        coalesce(CRSDepTime, DepTime, toUInt16(9999)) AS dep_sort,\n        min(coalesce(CRSArrTime, ArrTime, toUInt16(9999))) AS arr_sort\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier, OriginCode, DestCode, dep_sort\n), itineraries AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        count() AS hop_count,\n        min(dep_sort) AS first_dep_sort,\n        arraySort(groupArray((dep_sort, arr_sort, OriginCode, DestCode))) AS legs\n    FROM dedup_legs\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier\n    HAVING count() > 1\n), max_hops AS (\n    SELECT max(hop_count) AS max_hop_count FROM itineraries\n), top_unique AS (\n    SELECT\n        i.FlightDate,\n        i.first_dep_sort,\n        i.legs,\n        arrayStringConcat(arrayConcat(arrayMap(x -> x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') AS route,\n        row_number() OVER (\n            PARTITION BY arrayStringConcat(arrayConcat(arrayMap(x -> x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-')\n            ORDER BY i.FlightDate DESC, i.first_dep_sort DESC, i.carrier DESC, i.flight_number DESC, i.aircraft_id DESC\n        ) AS route_rank\n    FROM itineraries i\n    CROSS JOIN max_hops m\n    WHERE i.hop_count = m.max_hop_count\n), top10 AS (\n    SELECT *\n    FROM top_unique\n    WHERE route_rank = 1\n    ORDER BY FlightDate DESC, first_dep_sort DESC\n    LIMIT 10\n), airport_mentions AS (\n    SELECT arrayJoin(arrayConcat(arrayMap(x -> x.3, legs), [arrayElement(legs, length(legs)).4])) AS airport_code\n    FROM top10\n)\nSELECT\n    a.AirportCode AS airport_code,\n    a.StateCode AS state_code,\n    a.CountryCodeISO AS country_code,\n    multiIf(a.Longitude <= -115, 'far_west', a.Longitude <= -95, 'interior_west_central', 'east') AS geo_bucket,\n    count() AS airport_mentions\nFROM airport_mentions am\nLEFT JOIN ontime.dim_airports a\n    ON am.airport_code = toString(a.AirportCode)\n   AND a.IsLatest = 1\nGROUP BY airport_code, state_code, country_code, geo_bucket\nORDER BY airport_mentions DESC, airport_code\nLIMIT 100"
    }
  ]
}
```

Saved analysis.json:

```json
{
  "sql": "WITH dedup_legs AS (\n    SELECT\n        FlightDate,\n        ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,\n        Flight_Number_Reporting_Airline AS flight_number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginCode,\n        DestCode,\n        coalesce(CRSDepTime, DepTime, toUInt16(9999)) AS dep_sort,\n        min(coalesce(CRSArrTime, ArrTime, toUInt16(9999))) AS arr_sort\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier, OriginCode, DestCode, dep_sort\n), itineraries AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        count() AS hop_count,\n        min(dep_sort) AS first_dep_sort,\n        arraySort(groupArray((dep_sort, arr_sort, OriginCode, DestCode))) AS legs\n    FROM dedup_legs\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier\n    HAVING count() \u003e 1\n), route_recurrence AS (\n    SELECT\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, legs), [arrayElement(legs, length(legs)).4]), '-') AS route,\n        uniqExact(FlightDate) AS route_recurrence_count\n    FROM itineraries\n    GROUP BY route\n), max_hops AS (\n    SELECT max(hop_count) AS max_hop_count FROM itineraries\n), top_unique AS (\n    SELECT\n        i.FlightDate,\n        i.aircraft_id,\n        i.flight_number,\n        i.carrier,\n        i.hop_count,\n        i.first_dep_sort,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') AS route,\n        rr.route_recurrence_count,\n        row_number() OVER (\n            PARTITION BY arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-')\n            ORDER BY i.FlightDate DESC, i.first_dep_sort DESC, i.carrier DESC, i.flight_number DESC, i.aircraft_id DESC\n        ) AS route_rank\n    FROM itineraries i\n    CROSS JOIN max_hops m\n    INNER JOIN route_recurrence rr\n        ON arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') = rr.route\n    WHERE i.hop_count = m.max_hop_count\n)\nSELECT\n    if(aircraft_id = '', 'unknown', aircraft_id) AS aircraft_id,\n    flight_number,\n    carrier,\n    FlightDate AS flight_date,\n    hop_count,\n    route_recurrence_count,\n    route\nFROM top_unique\nWHERE route_rank = 1\nORDER BY flight_date DESC, first_dep_sort DESC, carrier DESC, flight_number DESC, aircraft_id DESC\nLIMIT 10",
  "report_markdown": "",
  "subquestions": [
    {
      "id": "main",
      "subquestion": "Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.\nDefine uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.\nDo not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.\nCount hops from distinct same-day legs, not raw source rows; do not let duplicate or conflicting same-time rows inflate hop count or create artifact routes.\n\nReturn:\n\n- aircraft id\n- flight number\n- carrier\n- flight date\n- hop count\n- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft\n- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`",
      "answer_markdown": "The maximum observed same-aircraft, same-flight-number itinerary length is 8 hops. The 10 most recent unique max-hop routes are all Southwest (WN) itineraries, with the newest on 2024-12-01: ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA flown by N957WN on flight 366. Across these 10 rows, route recurrence ranges from 1 day to 46 days, led by LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN (46), MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC (40), HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK (20), and BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK (12).",
      "sql": "WITH dedup_legs AS (\n    SELECT\n        FlightDate,\n        ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,\n        Flight_Number_Reporting_Airline AS flight_number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginCode,\n        DestCode,\n        coalesce(CRSDepTime, DepTime, toUInt16(9999)) AS dep_sort,\n        min(coalesce(CRSArrTime, ArrTime, toUInt16(9999))) AS arr_sort\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier, OriginCode, DestCode, dep_sort\n), itineraries AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        count() AS hop_count,\n        min(dep_sort) AS first_dep_sort,\n        arraySort(groupArray((dep_sort, arr_sort, OriginCode, DestCode))) AS legs\n    FROM dedup_legs\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier\n    HAVING count() \u003e 1\n), route_recurrence AS (\n    SELECT\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, legs), [arrayElement(legs, length(legs)).4]), '-') AS route,\n        uniqExact(FlightDate) AS route_recurrence_count\n    FROM itineraries\n    GROUP BY route\n), max_hops AS (\n    SELECT max(hop_count) AS max_hop_count FROM itineraries\n), top_unique AS (\n    SELECT\n        i.FlightDate,\n        i.aircraft_id,\n        i.flight_number,\n        i.carrier,\n        i.hop_count,\n        i.first_dep_sort,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') AS route,\n        rr.route_recurrence_count,\n        row_number() OVER (\n            PARTITION BY arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-')\n            ORDER BY i.FlightDate DESC, i.first_dep_sort DESC, i.carrier DESC, i.flight_number DESC, i.aircraft_id DESC\n        ) AS route_rank\n    FROM itineraries i\n    CROSS JOIN max_hops m\n    INNER JOIN route_recurrence rr\n        ON arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') = rr.route\n    WHERE i.hop_count = m.max_hop_count\n)\nSELECT\n    if(aircraft_id = '', 'unknown', aircraft_id) AS aircraft_id,\n    flight_number,\n    carrier,\n    FlightDate AS flight_date,\n    hop_count,\n    route_recurrence_count,\n    route\nFROM top_unique\nWHERE route_rank = 1\nORDER BY flight_date DESC, first_dep_sort DESC, carrier DESC, flight_number DESC, aircraft_id DESC\nLIMIT 10"
    },
    {
      "id": "q1",
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?\n\nList the recurrence count for each of the 10 routes explicitly before summarizing any tiers or categories, and make sure any category totals add up to 10.",
      "answer_markdown": "Recurrence counts for the 10 routes are: ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA = 1; CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN = 4; ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN = 2; MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX = 5; MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC = 40; LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN = 46; SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS = 7; HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK = 20; BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK = 12; BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX = 5. Overall, these top itineraries skew toward recurring patterns rather than pure one-offs: 4 are recurring (10+ days), 5 are occasional (2-9 days), and 1 is a one-off, which sums to all 10 routes.",
      "sql": "WITH main_top10 AS (\n    WITH dedup_legs AS (\n        SELECT\n            FlightDate,\n            ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,\n            Flight_Number_Reporting_Airline AS flight_number,\n            IATA_CODE_Reporting_Airline AS carrier,\n            OriginCode,\n            DestCode,\n            coalesce(CRSDepTime, DepTime, toUInt16(9999)) AS dep_sort,\n            min(coalesce(CRSArrTime, ArrTime, toUInt16(9999))) AS arr_sort\n        FROM ontime.fact_ontime\n        WHERE Cancelled = 0\n        GROUP BY FlightDate, aircraft_id, flight_number, carrier, OriginCode, DestCode, dep_sort\n    ), itineraries AS (\n        SELECT\n            FlightDate,\n            aircraft_id,\n            flight_number,\n            carrier,\n            count() AS hop_count,\n            min(dep_sort) AS first_dep_sort,\n            arraySort(groupArray((dep_sort, arr_sort, OriginCode, DestCode))) AS legs\n        FROM dedup_legs\n        GROUP BY FlightDate, aircraft_id, flight_number, carrier\n        HAVING count() \u003e 1\n    ), route_recurrence AS (\n        SELECT\n            arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, legs), [arrayElement(legs, length(legs)).4]), '-') AS route,\n            uniqExact(FlightDate) AS route_recurrence_count\n        FROM itineraries\n        GROUP BY route\n    ), max_hops AS (\n        SELECT max(hop_count) AS max_hop_count FROM itineraries\n    ), top_unique AS (\n        SELECT\n            i.FlightDate,\n            i.first_dep_sort,\n            arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') AS route,\n            rr.route_recurrence_count,\n            row_number() OVER (\n                PARTITION BY arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-')\n                ORDER BY i.FlightDate DESC, i.first_dep_sort DESC, i.carrier DESC, i.flight_number DESC, i.aircraft_id DESC\n            ) AS route_rank\n        FROM itineraries i\n        CROSS JOIN max_hops m\n        INNER JOIN route_recurrence rr\n            ON arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') = rr.route\n        WHERE i.hop_count = m.max_hop_count\n    )\n    SELECT route, route_recurrence_count, FlightDate, first_dep_sort\n    FROM top_unique\n    WHERE route_rank = 1\n    ORDER BY FlightDate DESC, first_dep_sort DESC\n    LIMIT 10\n)\nSELECT\n    route,\n    route_recurrence_count,\n    multiIf(route_recurrence_count = 1, 'one-off', route_recurrence_count \u003c 10, 'occasional', 'recurring') AS recurrence_tier,\n    count() OVER (PARTITION BY multiIf(route_recurrence_count = 1, 'one-off', route_recurrence_count \u003c 10, 'occasional', 'recurring')) AS tier_count\nFROM main_top10\nORDER BY FlightDate DESC, first_dep_sort DESC"
    },
    {
      "id": "q2",
      "subquestion": "What geographic pattern do the top itineraries show?\n\nBase the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.",
      "answer_markdown": "The geographic pattern is entirely domestic U.S. and strongly tied to Southwest/West Coast plus Sun Belt and Mid-Atlantic nodes. Every airport in the 10 routes is a US airport, and the most repeated points are BWI, DAL, DEN, LAS, MSY, and OAK with 5 mentions each, followed by BNA and PHX with 4 each. California is especially prominent through OAK, BUR, LAX, SAN, PSP, SJC, and SMF, so these longest same-flight-number days look like long domestic zig-zag chains linking eastern and southeastern stops into western endpoints and western connection banks.",
      "sql": "WITH dedup_legs AS (\n    SELECT\n        FlightDate,\n        ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,\n        Flight_Number_Reporting_Airline AS flight_number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginCode,\n        DestCode,\n        coalesce(CRSDepTime, DepTime, toUInt16(9999)) AS dep_sort,\n        min(coalesce(CRSArrTime, ArrTime, toUInt16(9999))) AS arr_sort\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier, OriginCode, DestCode, dep_sort\n), itineraries AS (\n    SELECT\n        FlightDate,\n        aircraft_id,\n        flight_number,\n        carrier,\n        count() AS hop_count,\n        min(dep_sort) AS first_dep_sort,\n        arraySort(groupArray((dep_sort, arr_sort, OriginCode, DestCode))) AS legs\n    FROM dedup_legs\n    GROUP BY FlightDate, aircraft_id, flight_number, carrier\n    HAVING count() \u003e 1\n), max_hops AS (\n    SELECT max(hop_count) AS max_hop_count FROM itineraries\n), top_unique AS (\n    SELECT\n        i.FlightDate,\n        i.first_dep_sort,\n        i.legs,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') AS route,\n        row_number() OVER (\n            PARTITION BY arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-')\n            ORDER BY i.FlightDate DESC, i.first_dep_sort DESC, i.carrier DESC, i.flight_number DESC, i.aircraft_id DESC\n        ) AS route_rank\n    FROM itineraries i\n    CROSS JOIN max_hops m\n    WHERE i.hop_count = m.max_hop_count\n), top10 AS (\n    SELECT *\n    FROM top_unique\n    WHERE route_rank = 1\n    ORDER BY FlightDate DESC, first_dep_sort DESC\n    LIMIT 10\n), airport_mentions AS (\n    SELECT arrayJoin(arrayConcat(arrayMap(x -\u003e x.3, legs), [arrayElement(legs, length(legs)).4])) AS airport_code\n    FROM top10\n)\nSELECT\n    a.AirportCode AS airport_code,\n    a.StateCode AS state_code,\n    a.CountryCodeISO AS country_code,\n    multiIf(a.Longitude \u003c= -115, 'far_west', a.Longitude \u003c= -95, 'interior_west_central', 'east') AS geo_bucket,\n    count() AS airport_mentions\nFROM airport_mentions am\nLEFT JOIN ontime.dim_airports a\n    ON am.airport_code = toString(a.AirportCode)\n   AND a.IsLatest = 1\nGROUP BY airport_code, state_code, country_code, geo_bucket\nORDER BY airport_mentions DESC, airport_code\nLIMIT 100"
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
```

Proof queries are saved as files in the run directory. Read the SQL files you need to verify grain, filters, metrics, and ranking logic:

- `queries/main.sql`
- `queries/q1.sql`
- `queries/q2.sql`

Executed query results are saved as files in the run directory. Read the result files you need for verification instead of assuming the report summary is complete:

- `results/main.json`
- `results/q1.json`
- `results/q2.json`