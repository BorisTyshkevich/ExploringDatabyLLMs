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
- Check whether all dashboard questions were answered directly.
- Check whether SQL grain and returned metrics match the prompt.
- Check whether prose claims are supported by the executed query results.
- Check for duplicated entities, inconsistent counts, missing requested metrics, unsupported inference, and proof-query/result mismatches.
- Treat `report.md` as a monitoring artifact. Do not treat its one-row example table or abbreviated preview as a defect when the underlying query result preserves the required ranked or detailed rows.
- In `## Findings`, cite concrete artifact names such as `queries/q1.sql`, `results/q1.json`, or `report.md`.
- In `## Suggested Prompt Fixes`, propose prompt-level changes only when they would reduce the observed failure mode.
- When proposing a prompt fix, be concrete. Prefer 1-3 short replacement or insertion snippets that could be added to the question prompt, not vague advice.
- Tie each prompt fix to the specific failure you found. Explain what ambiguity in the current prompt allowed the bad result and how the new wording would close that gap.
- If the prompt is already sufficiently specific and the failure is not plausibly prompt-driven, say `None.` instead of inventing weak prompt edits.
- Do not suggest code changes to qforge in this review. Focus on the run and the prompt.
- Write Markdown only to `review.md`.

Question-specific guidance:

Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.
Define uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.

Return:

- aircraft id
- flight number
- carrier
- flight date
- hop count
- number of flights
- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`

## Dashboard Questions

- Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?
- What geographic pattern do the top itineraries show?

Generated report.md:

```md
# Highest daily hops for one aircraft on one flight number

> Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?

The top itineraries are overwhelmingly recurring scheduled patterns. The highest-frequency 8-hop routes each appear dozens to over a hundred times across different dates, spanning multi-month windows. For example, the route `FLL-JAX-IND-MDW-MCI-DAL-ABQ-LAX-SJC` appeared 114 times between March and August 2008, and `CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN` (WN flight 3149) appears across multiple weeks in 2024. These are classic Southwest Airlines through-plane rotation schedules that repeat the same airport sequence day after day on fixed flight numbers.

- Rows returned: 10
- Columns: Route, max_hops, occurrences, first_seen, last_seen

| Route | max_hops | occurrences | first_seen | last_seen |
| --- | --- | --- | --- | --- |
| FLL-JAX-IND-MDW-MCI-DAL-ABQ-LAX-SJC | 8 | 114 | 2008-03-10T00:00:00Z | 2008-08-21T00:00:00Z |

> What geographic pattern do the top itineraries show?

All top itineraries are Southwest Airlines (WN) cross-country daisy-chains that traverse the full width of the continental United States. Routes typically originate on the East Coast or Gulf Coast (e.g., FLL, PVD, BDL, LGA, HOU, MSY), pass through Midwest hubs (MDW, STL, MCI, CMH), continue through Texas/Southwest hubs (DAL, ABQ, PHX, ELP), and terminate on the West Coast (SJC, OAK, ONT, SEA, RNO). Each 8-hop route covers 6–7 different states. The pattern reflects Southwest's hub-and-spoke-free model where a single aircraft flies a continuous east-to-west (or west-to-east) rotation through its point-to-point network in a single day.

- Rows returned: 10
- Columns: Route, occurrences, state_count, states_sample

| Route | occurrences | state_count | states_sample |
| --- | --- | --- | --- |
| FLL-JAX-IND-MDW-MCI-DAL-ABQ-LAX-SJC | 99 | 7 | [NM, TX, FL, IN, CA, MO, IL] |
```

Saved main.sql:

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

Saved answer.raw.json:

```json
{
  "subquestions": [
    {
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?",
      "answer_markdown": "The top itineraries are overwhelmingly recurring scheduled patterns. The highest-frequency 8-hop routes each appear dozens to over a hundred times across different dates, spanning multi-month windows. For example, the route `FLL-JAX-IND-MDW-MCI-DAL-ABQ-LAX-SJC` appeared 114 times between March and August 2008, and `CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN` (WN flight 3149) appears across multiple weeks in 2024. These are classic Southwest Airlines through-plane rotation schedules that repeat the same airport sequence day after day on fixed flight numbers.",
      "sql": "WITH itineraries AS (SELECT Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline AS Carrier, FlightDate, count() AS hops, arrayStringConcat(arrayConcat(arraySort((x, t) -> t, groupArray(OriginCode), groupArray(assumeNotNull(CRSDepTime))), [argMax(DestCode, assumeNotNull(CRSDepTime))]), '-') AS Route FROM ontime.fact_ontime WHERE Tail_Number != '' AND Cancelled = 0 GROUP BY Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate), route_stats AS (SELECT Route, max(hops) AS max_hops, count() AS occurrences, min(FlightDate) AS first_seen, max(FlightDate) AS last_seen FROM itineraries GROUP BY Route) SELECT Route, max_hops, occurrences, first_seen, last_seen FROM route_stats ORDER BY max_hops DESC, occurrences DESC LIMIT 10"
    },
    {
      "subquestion": "What geographic pattern do the top itineraries show?",
      "answer_markdown": "All top itineraries are Southwest Airlines (WN) cross-country daisy-chains that traverse the full width of the continental United States. Routes typically originate on the East Coast or Gulf Coast (e.g., FLL, PVD, BDL, LGA, HOU, MSY), pass through Midwest hubs (MDW, STL, MCI, CMH), continue through Texas/Southwest hubs (DAL, ABQ, PHX, ELP), and terminate on the West Coast (SJC, OAK, ONT, SEA, RNO). Each 8-hop route covers 6–7 different states. The pattern reflects Southwest's hub-and-spoke-free model where a single aircraft flies a continuous east-to-west (or west-to-east) rotation through its point-to-point network in a single day.",
      "sql": "WITH itineraries AS (SELECT FlightDate, count() AS hops, arrayStringConcat(arrayConcat(arraySort((x, t) -> t, groupArray(OriginCode), groupArray(assumeNotNull(CRSDepTime))), [argMax(DestCode, assumeNotNull(CRSDepTime))]), '-') AS Route, arrayDistinct(groupArray(OriginState)) AS states_covered FROM ontime.fact_ontime WHERE Tail_Number != '' AND Cancelled = 0 GROUP BY FlightDate, Tail_Number, Flight_Number_Reporting_Airline HAVING hops = 8), route_stats AS (SELECT Route, count() AS occurrences, any(states_covered) AS states_sample FROM itineraries GROUP BY Route) SELECT Route, occurrences, length(states_sample) AS state_count, states_sample FROM route_stats ORDER BY occurrences DESC LIMIT 10"
    }
  ]
}
```

Saved analysis.json:

```json
{
  "sql": "WITH itineraries AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline,\n        IATA_CODE_Reporting_Airline AS Carrier,\n        FlightDate,\n        count() AS hops,\n        arrayStringConcat(\n            arrayConcat(\n                arraySort((x, t) -\u003e t, groupArray(OriginCode), groupArray(assumeNotNull(CRSDepTime))),\n                [argMax(DestCode, assumeNotNull(CRSDepTime))]\n            ),\n            '-'\n        ) AS Route\n    FROM ontime.fact_ontime\n    WHERE Tail_Number != '' AND Cancelled = 0\n    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate\n),\nbest_by_route AS (\n    SELECT\n        Route,\n        max(hops) AS max_hops,\n        argMax(Tail_Number, FlightDate) AS aircraft_id,\n        argMax(Flight_Number_Reporting_Airline, FlightDate) AS flight_number,\n        argMax(Carrier, FlightDate) AS carrier,\n        max(FlightDate) AS flight_date,\n        argMax(hops, FlightDate) AS latest_hops\n    FROM itineraries\n    GROUP BY Route\n)\nSELECT\n    aircraft_id,\n    flight_number,\n    carrier,\n    flight_date,\n    max_hops AS hop_count,\n    latest_hops AS num_flights,\n    Route\nFROM best_by_route\nORDER BY max_hops DESC, flight_date DESC\nLIMIT 10",
  "report_markdown": "",
  "subquestions": [
    {
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?",
      "answer_markdown": "The top itineraries are overwhelmingly recurring scheduled patterns. The highest-frequency 8-hop routes each appear dozens to over a hundred times across different dates, spanning multi-month windows. For example, the route `FLL-JAX-IND-MDW-MCI-DAL-ABQ-LAX-SJC` appeared 114 times between March and August 2008, and `CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN` (WN flight 3149) appears across multiple weeks in 2024. These are classic Southwest Airlines through-plane rotation schedules that repeat the same airport sequence day after day on fixed flight numbers.",
      "sql": "WITH itineraries AS (SELECT Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline AS Carrier, FlightDate, count() AS hops, arrayStringConcat(arrayConcat(arraySort((x, t) -\u003e t, groupArray(OriginCode), groupArray(assumeNotNull(CRSDepTime))), [argMax(DestCode, assumeNotNull(CRSDepTime))]), '-') AS Route FROM ontime.fact_ontime WHERE Tail_Number != '' AND Cancelled = 0 GROUP BY Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate), route_stats AS (SELECT Route, max(hops) AS max_hops, count() AS occurrences, min(FlightDate) AS first_seen, max(FlightDate) AS last_seen FROM itineraries GROUP BY Route) SELECT Route, max_hops, occurrences, first_seen, last_seen FROM route_stats ORDER BY max_hops DESC, occurrences DESC LIMIT 10"
    },
    {
      "subquestion": "What geographic pattern do the top itineraries show?",
      "answer_markdown": "All top itineraries are Southwest Airlines (WN) cross-country daisy-chains that traverse the full width of the continental United States. Routes typically originate on the East Coast or Gulf Coast (e.g., FLL, PVD, BDL, LGA, HOU, MSY), pass through Midwest hubs (MDW, STL, MCI, CMH), continue through Texas/Southwest hubs (DAL, ABQ, PHX, ELP), and terminate on the West Coast (SJC, OAK, ONT, SEA, RNO). Each 8-hop route covers 6–7 different states. The pattern reflects Southwest's hub-and-spoke-free model where a single aircraft flies a continuous east-to-west (or west-to-east) rotation through its point-to-point network in a single day.",
      "sql": "WITH itineraries AS (SELECT FlightDate, count() AS hops, arrayStringConcat(arrayConcat(arraySort((x, t) -\u003e t, groupArray(OriginCode), groupArray(assumeNotNull(CRSDepTime))), [argMax(DestCode, assumeNotNull(CRSDepTime))]), '-') AS Route, arrayDistinct(groupArray(OriginState)) AS states_covered FROM ontime.fact_ontime WHERE Tail_Number != '' AND Cancelled = 0 GROUP BY FlightDate, Tail_Number, Flight_Number_Reporting_Airline HAVING hops = 8), route_stats AS (SELECT Route, count() AS occurrences, any(states_covered) AS states_sample FROM itineraries GROUP BY Route) SELECT Route, occurrences, length(states_sample) AS state_count, states_sample FROM route_stats ORDER BY occurrences DESC LIMIT 10"
    }
  ]
}
```

Saved visual_input.json:

```json
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
```

Proof queries are saved as files in the run directory. Read the SQL files you need to verify grain, filters, metrics, and ranking logic:

- `main.sql`
- `queries/q1.sql`
- `queries/q2.sql`

Executed query results are saved as files in the run directory. Read the result files you need for verification instead of assuming the report summary is complete:

- `results/main.json`
- `results/q1.json`
- `results/q2.json`