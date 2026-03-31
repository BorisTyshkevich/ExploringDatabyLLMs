```sql
WITH annual_totals AS (
    SELECT
        Year,
        count() AS YearCompletedFlights
    FROM default.ontime_v2
    WHERE Cancelled = 0
    GROUP BY Year
),
annual_carrier_counts AS (
    SELECT
        Year,
        Reporting_Airline,
        count() AS CompletedFlights
    FROM default.ontime_v2
    WHERE Cancelled = 0
    GROUP BY
        Year,
        Reporting_Airline
),
ranked_carriers AS (
    SELECT
        acc.Year,
        acc.Reporting_Airline,
        row_number() OVER (
            PARTITION BY acc.Year
            ORDER BY acc.CompletedFlights DESC, acc.Reporting_Airline ASC
        ) AS RankInYear,
        acc.CompletedFlights,
        100.0 * acc.CompletedFlights / at.YearCompletedFlights AS SharePct
    FROM annual_carrier_counts AS acc
    INNER JOIN annual_totals AS at
        ON acc.Year = at.Year
),
leader_runner_up AS (
    SELECT
        Year,
        maxIf(Reporting_Airline, RankInYear = 1) AS LeaderReportingAirline,
        maxIf(Reporting_Airline, RankInYear = 2) AS RunnerUpReportingAirline,
        maxIf(SharePct, RankInYear = 1) AS LeaderSharePct,
        maxIf(SharePct, RankInYear = 2) AS RunnerUpSharePct
    FROM ranked_carriers
    WHERE RankInYear <= 2
    GROUP BY Year
),
leader_transitions AS (
    SELECT
        Year,
        LeaderReportingAirline,
        RunnerUpReportingAirline,
        LeaderSharePct - RunnerUpSharePct AS LeaderShareGapPctPts,
        prior_year_leader AS PriorYearLeaderReportingAirline,
        if(prior_year_leader IS NULL, 0, LeaderReportingAirline != prior_year_leader) AS LeaderChanged,
        if(
            prior_year_leader IS NULL,
            cast(NULL, 'Nullable(Float64)'),
            LeaderSharePct - prior_year_share_pct
        ) AS LeaderShareChangePctPts
    FROM (
        SELECT
            Year,
            LeaderReportingAirline,
            RunnerUpReportingAirline,
            LeaderSharePct,
            RunnerUpSharePct,
            lagInFrame(toNullable(LeaderReportingAirline)) OVER w AS prior_year_leader,
            lagInFrame(toNullable(LeaderSharePct)) OVER w AS prior_year_share_pct
        FROM leader_runner_up
        WINDOW w AS (ORDER BY Year ASC ROWS BETWEEN 1 PRECEDING AND CURRENT ROW)
    )
)
SELECT
    'carrier_year' AS RowType,
    rc.Year,
    rc.Reporting_Airline,
    rc.RankInYear,
    rc.CompletedFlights,
    rc.SharePct,
    if(rc.RankInYear = 1, lt.LeaderReportingAirline, cast(NULL, 'Nullable(String)')) AS LeaderReportingAirline,
    if(rc.RankInYear = 1, lt.RunnerUpReportingAirline, cast(NULL, 'Nullable(String)')) AS RunnerUpReportingAirline,
    if(rc.RankInYear = 1, lt.LeaderShareGapPctPts, cast(NULL, 'Nullable(Float64)')) AS LeaderShareGapPctPts,
    if(rc.RankInYear = 1, lt.PriorYearLeaderReportingAirline, cast(NULL, 'Nullable(String)')) AS PriorYearLeaderReportingAirline,
    if(rc.RankInYear = 1, toUInt8(lt.LeaderChanged), cast(NULL, 'Nullable(UInt8)')) AS LeaderChanged,
    if(rc.RankInYear = 1, lt.LeaderShareChangePctPts, cast(NULL, 'Nullable(Float64)')) AS LeaderShareChangePctPts
FROM ranked_carriers AS rc
INNER JOIN leader_transitions AS lt
    ON rc.Year = lt.Year
WHERE rc.RankInYear <= 5
ORDER BY
    rc.Year ASC,
    RowType,
    rc.RankInYear ASC,
    rc.Reporting_Airline ASC
```