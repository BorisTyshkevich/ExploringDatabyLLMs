- Connect to clickhouse server only though MCP connection
- Use the `ontime` database to answer analytical questions
- Use `ontime-semantic-layer` skill for schema inspection, join guidance, and dimension semantics.
- write correct and efficient ClickHouse SQL 
- Before writing any SQL artifact, self-verify every SQL statement you intend to save. Run a cheap debug execution for each query first, usually with a small `LIMIT`, a narrow `WHERE` filter, or both applied inside the main data-reading subquery or CTE.
- Treat successful execution as mandatory. Fix any syntax, type, aggregate, window, join, or unknown-column errors in a loop until every saved query runs successfully. Do not write not tested SQL to the artifact.

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

Across the most recent five years in the OnTime data, the longest same-aircraft, same-flight-number daily itineraries reached 8 hops. The 10 most recent unique routes at that maximum include ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA on 2024-12-01 (route recurrence 1 day), CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN on 2024-02-18 (4 days), and ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN on 2023-04-30 (2 days).

- Rows returned: 10
- Columns: aircraft_id, flight_number, carrier, flight_date, hop_count, route_recurrence_count, Route

| aircraft_id | flight_number | carrier | flight_date | hop_count | route_recurrence_count | Route |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |
```

Saved answer.raw.json:

```json
{
  "subquestions": [
    {
      "id": "main",
      "answer_markdown": "Across the most recent five years in the OnTime data, the longest same-aircraft, same-flight-number daily itineraries reached 8 hops. The 10 most recent unique routes at that maximum include ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA on 2024-12-01 (route recurrence 1 day), CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN on 2024-02-18 (4 days), and ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN on 2023-04-30 (2 days).",
      "sql": "WITH\n    (SELECT max(FlightDate) FROM ontime.fact_ontime) AS max_fd,\n    addYears(max_fd, -5) AS start_fd,\n    base AS (\n        SELECT\n            FlightDate,\n            Carrier,\n            FlightNum,\n            if(Tail_Number = '', '', Tail_Number) AS aircraft_id,\n            OriginCode,\n            DestCode,\n            coalesce(CRSDepTime, DepTime, 0) AS dep_sort,\n            coalesce(CRSArrTime, ArrTime, 0) AS arr_sort\n        FROM ontime.fact_ontime\n        WHERE FlightDate > start_fd\n          AND Cancelled = 0\n          AND FlightNum != ''\n          AND Carrier != ''\n    ),\n    dedup AS (\n        SELECT\n            FlightDate,\n            Carrier,\n            FlightNum,\n            aircraft_id,\n            dep_sort,\n            min(arr_sort) AS arr_sort,\n            min(OriginCode) AS OriginCode,\n            min(DestCode) AS DestCode\n        FROM base\n        GROUP BY FlightDate, Carrier, FlightNum, aircraft_id, dep_sort\n    ),\n    itineraries AS (\n        SELECT\n            FlightDate,\n            Carrier,\n            FlightNum,\n            aircraft_id,\n            min(dep_sort) AS first_dep_time,\n            arraySort(x -> (x.1, x.2, x.3, x.4), groupArray((dep_sort, arr_sort, OriginCode, DestCode))) AS legs,\n            length(legs) AS hop_count,\n            arrayStringConcat(arrayConcat(arrayMap(x -> x.3, legs), [legs[length(legs)].4]), '-') AS Route\n        FROM dedup\n        GROUP BY FlightDate, Carrier, FlightNum, aircraft_id\n        HAVING hop_count > 1\n    ),\n    route_days AS (\n        SELECT\n            Route,\n            countDistinct(FlightDate) AS route_recurrence_count\n        FROM itineraries\n        GROUP BY Route\n    ),\n    ranked AS (\n        SELECT\n            i.aircraft_id,\n            i.FlightNum,\n            i.Carrier,\n            i.FlightDate,\n            i.hop_count,\n            r.route_recurrence_count,\n            i.Route,\n            i.first_dep_time,\n            row_number() OVER (\n                PARTITION BY i.Route\n                ORDER BY i.FlightDate DESC, i.first_dep_time DESC, i.hop_count DESC, i.Carrier, i.FlightNum, i.aircraft_id\n            ) AS rn\n        FROM itineraries AS i\n        INNER JOIN route_days AS r ON i.Route = r.Route\n    )\nSELECT\n    aircraft_id,\n    FlightNum AS flight_number,\n    Carrier AS carrier,\n    FlightDate AS flight_date,\n    hop_count,\n    route_recurrence_count,\n    Route\nFROM ranked\nWHERE rn = 1\nORDER BY hop_count DESC, flight_date DESC, first_dep_time DESC, Route\nLIMIT 10\n"
    }
  ]
}
```

Saved analysis.json:

```json
{
  "sql": "WITH\n    (SELECT max(FlightDate) FROM ontime.fact_ontime) AS max_fd,\n    addYears(max_fd, -5) AS start_fd,\n    base AS (\n        SELECT\n            FlightDate,\n            Carrier,\n            FlightNum,\n            if(Tail_Number = '', '', Tail_Number) AS aircraft_id,\n            OriginCode,\n            DestCode,\n            coalesce(CRSDepTime, DepTime, 0) AS dep_sort,\n            coalesce(CRSArrTime, ArrTime, 0) AS arr_sort\n        FROM ontime.fact_ontime\n        WHERE FlightDate \u003e start_fd\n          AND Cancelled = 0\n          AND FlightNum != ''\n          AND Carrier != ''\n    ),\n    dedup AS (\n        SELECT\n            FlightDate,\n            Carrier,\n            FlightNum,\n            aircraft_id,\n            dep_sort,\n            min(arr_sort) AS arr_sort,\n            min(OriginCode) AS OriginCode,\n            min(DestCode) AS DestCode\n        FROM base\n        GROUP BY FlightDate, Carrier, FlightNum, aircraft_id, dep_sort\n    ),\n    itineraries AS (\n        SELECT\n            FlightDate,\n            Carrier,\n            FlightNum,\n            aircraft_id,\n            min(dep_sort) AS first_dep_time,\n            arraySort(x -\u003e (x.1, x.2, x.3, x.4), groupArray((dep_sort, arr_sort, OriginCode, DestCode))) AS legs,\n            length(legs) AS hop_count,\n            arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, legs), [legs[length(legs)].4]), '-') AS Route\n        FROM dedup\n        GROUP BY FlightDate, Carrier, FlightNum, aircraft_id\n        HAVING hop_count \u003e 1\n    ),\n    route_days AS (\n        SELECT\n            Route,\n            countDistinct(FlightDate) AS route_recurrence_count\n        FROM itineraries\n        GROUP BY Route\n    ),\n    ranked AS (\n        SELECT\n            i.aircraft_id,\n            i.FlightNum,\n            i.Carrier,\n            i.FlightDate,\n            i.hop_count,\n            r.route_recurrence_count,\n            i.Route,\n            i.first_dep_time,\n            row_number() OVER (\n                PARTITION BY i.Route\n                ORDER BY i.FlightDate DESC, i.first_dep_time DESC, i.hop_count DESC, i.Carrier, i.FlightNum, i.aircraft_id\n            ) AS rn\n        FROM itineraries AS i\n        INNER JOIN route_days AS r ON i.Route = r.Route\n    )\nSELECT\n    aircraft_id,\n    FlightNum AS flight_number,\n    Carrier AS carrier,\n    FlightDate AS flight_date,\n    hop_count,\n    route_recurrence_count,\n    Route\nFROM ranked\nWHERE rn = 1\nORDER BY hop_count DESC, flight_date DESC, first_dep_time DESC, Route\nLIMIT 10",
  "report_markdown": "",
  "subquestions": [
    {
      "id": "main",
      "subquestion": "Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.\nDefine uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.\nDo not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.\nCount hops from distinct same-day legs, not raw source rows; do not let duplicate or conflicting same-time rows inflate hop count or create artifact routes.\n\nReturn:\n\n- aircraft id\n- flight number\n- carrier\n- flight date\n- hop count\n- route recurrence count: total number of days across the analyzed window on which this exact Route string was flown by any aircraft\n- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`",
      "answer_markdown": "Across the most recent five years in the OnTime data, the longest same-aircraft, same-flight-number daily itineraries reached 8 hops. The 10 most recent unique routes at that maximum include ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA on 2024-12-01 (route recurrence 1 day), CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN on 2024-02-18 (4 days), and ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN on 2023-04-30 (2 days).",
      "sql": "WITH\n    (SELECT max(FlightDate) FROM ontime.fact_ontime) AS max_fd,\n    addYears(max_fd, -5) AS start_fd,\n    base AS (\n        SELECT\n            FlightDate,\n            Carrier,\n            FlightNum,\n            if(Tail_Number = '', '', Tail_Number) AS aircraft_id,\n            OriginCode,\n            DestCode,\n            coalesce(CRSDepTime, DepTime, 0) AS dep_sort,\n            coalesce(CRSArrTime, ArrTime, 0) AS arr_sort\n        FROM ontime.fact_ontime\n        WHERE FlightDate \u003e start_fd\n          AND Cancelled = 0\n          AND FlightNum != ''\n          AND Carrier != ''\n    ),\n    dedup AS (\n        SELECT\n            FlightDate,\n            Carrier,\n            FlightNum,\n            aircraft_id,\n            dep_sort,\n            min(arr_sort) AS arr_sort,\n            min(OriginCode) AS OriginCode,\n            min(DestCode) AS DestCode\n        FROM base\n        GROUP BY FlightDate, Carrier, FlightNum, aircraft_id, dep_sort\n    ),\n    itineraries AS (\n        SELECT\n            FlightDate,\n            Carrier,\n            FlightNum,\n            aircraft_id,\n            min(dep_sort) AS first_dep_time,\n            arraySort(x -\u003e (x.1, x.2, x.3, x.4), groupArray((dep_sort, arr_sort, OriginCode, DestCode))) AS legs,\n            length(legs) AS hop_count,\n            arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.3, legs), [legs[length(legs)].4]), '-') AS Route\n        FROM dedup\n        GROUP BY FlightDate, Carrier, FlightNum, aircraft_id\n        HAVING hop_count \u003e 1\n    ),\n    route_days AS (\n        SELECT\n            Route,\n            countDistinct(FlightDate) AS route_recurrence_count\n        FROM itineraries\n        GROUP BY Route\n    ),\n    ranked AS (\n        SELECT\n            i.aircraft_id,\n            i.FlightNum,\n            i.Carrier,\n            i.FlightDate,\n            i.hop_count,\n            r.route_recurrence_count,\n            i.Route,\n            i.first_dep_time,\n            row_number() OVER (\n                PARTITION BY i.Route\n                ORDER BY i.FlightDate DESC, i.first_dep_time DESC, i.hop_count DESC, i.Carrier, i.FlightNum, i.aircraft_id\n            ) AS rn\n        FROM itineraries AS i\n        INNER JOIN route_days AS r ON i.Route = r.Route\n    )\nSELECT\n    aircraft_id,\n    FlightNum AS flight_number,\n    Carrier AS carrier,\n    FlightDate AS flight_date,\n    hop_count,\n    route_recurrence_count,\n    Route\nFROM ranked\nWHERE rn = 1\nORDER BY hop_count DESC, flight_date DESC, first_dep_time DESC, Route\nLIMIT 10"
    }
  ]
}
```

Saved visual_input.json:

```json
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
```

Proof queries are saved as files in the run directory. Read the SQL files you need to verify grain, filters, metrics, and ranking logic:

- `queries/main.sql`

Executed query results are saved as files in the run directory. Read the result files you need for verification instead of assuming the report summary is complete:

- `results/main.json`