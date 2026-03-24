# q004 Run Review Notes

## Review Summary

This run is directionally good, but it is not fully valid against the current q004 question.

- The headline answer still appears correct: `MDW` is the worst qualifying airport by departure OTP.
- The saved ranking query is not truly airport-level and can split one airport into multiple rows.
- One requested business metric is missing from the ranked output.
- The report makes a stronger claim about "major hubs" than the proof query actually demonstrates.

## Findings

### 1. Ranking query is not actually one row per airport

The saved ranking query groups by:

- `f.OriginCode`
- `d.DisplayAirportName`
- `f.OriginCityName`

That can split the same airport into multiple ranked rows when descriptive fields change over time.

Observed in this run:

- `q1.json` has `161` rows
- `q2.json` says there are `160` qualifying airports
- `AUS` appears twice in `q1.json`
  - one historical row before Bergstrom
  - one modern Bergstrom-era row

So the ranked proof set is malformed even though the top airport remains `MDW`.

### 2. The prompt asked for average departure delay, but the query returns only average delay when late

The ranked queries return:

- `avg_dep_delay_when_late`

but not:

- average departure delay across all completed departures

Those are not interchangeable metrics.

### 3. The "high-delay" measure is weaker than it should be

The query uses:

- `quantileExact(0.90)(DepDelayMinutes)`

without limiting to delayed departures.

That means on-time or zero-delay flights can compress the tail measure. If the question is meant to capture the worse end of delay severity, a delayed-flight tail metric is more defensible.

### 4. The "mostly major hubs" conclusion is partly inference, not proof

The third proof query returns:

- airport identity
- traffic
- OTP and delay metrics

It does **not** return an explicit hub/non-hub classification. So the report should either:

- label that statement as an inference from the airport set and traffic levels, or
- use a proof query that directly supports the classification claim

## Independent Verification

A corrected airport-level aggregation still ranks `MDW` worst at `76.52%` OTP, so the headline answer survives.

The problem is not the top-line winner. The problem is that the saved proof query and report are not fully aligned with the question's required grain and metrics.

## Proposed Prompt Fix For q004

Suggested changes for:

- `/Users/bvt/work/ExploringDatabyLLMs/prompts/q004_worst_origin_airport_otp_thresholded/report_prompt.md`

### Replace the metric section with

```md
For each airport, quantify:

- completed departures
- departure on-time performance
- average departure delay across all completed departures
- a high-delay measure that reflects the worse end of the delay distribution for delayed departures
- the first and last dates represented in the data
```

### Add this grain constraint after the metric section

```md
The ranking must be one row per origin airport. Do not split the same airport into multiple ranked rows because of airport-name changes, city labels, or dimension-table history. If descriptive fields vary over time, still aggregate at the airport level and choose one representative display name for the output row.
```

### Tighten the proof-query section to

```md
Provide one proof query for each required business question. Across those proof queries, include enough evidence to support both:

- a ranked view of the weakest qualifying origin airports
- a comparison view that makes the gap between the very worst airports and the middle of the ranked set easy to judge

The proof query behind the worst-airport question should preserve the ranked airport-level rows needed for the dashboard, not just a single top airport or summary statistic.

Use the same qualifying-airport definition consistently across the ranking and spread questions.
```

### Tighten the third dashboard question guidance to

```md
- Are the weakest airports mostly major hubs, or is the bottom group more mixed?

In the report, answer that question using only evidence returned by the proof query. If you make an interpretation such as "mostly major hubs," make it explicit that it is an inference from the returned airport set and traffic levels, not a directly queried classification.
```

## Practical Effect Of The Fix

This should prevent:

- duplicate-airport rows caused by descriptive-field grouping
- missing `avg_dep_delay` in the ranked output
- a diluted high-delay metric
- overstated prose that goes beyond the proof query
