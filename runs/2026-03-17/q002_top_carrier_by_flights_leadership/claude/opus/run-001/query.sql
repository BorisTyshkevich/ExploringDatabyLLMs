WITH 
-- Step 1: Annual totals for completed flights
annual_totals AS (
    SELECT 
        Year,
        count() AS TotalCompletedFlights
    FROM default.ontime_v2
    WHERE Cancelled = 0
    GROUP BY Year
),

-- Step 2: Carrier-year aggregates for completed flights
carrier_year_stats AS (
    SELECT 
        Year,
        Reporting_Airline,
        count() AS CompletedFlights
    FROM default.ontime_v2
    WHERE Cancelled = 0
    GROUP BY Year, Reporting_Airline
),

-- Step 3: Ranked carrier-year with share calculation
ranked_carriers AS (
    SELECT 
        cy.Year,
        cy.Reporting_Airline,
        cy.CompletedFlights,
        round(100.0 * cy.CompletedFlights / at.TotalCompletedFlights, 4) AS SharePct,
        row_number() OVER (PARTITION BY cy.Year ORDER BY cy.CompletedFlights DESC, cy.Reporting_Airline ASC) AS RankInYear
    FROM carrier_year_stats cy
    JOIN annual_totals at ON cy.Year = at.Year
),

-- Step 4: Extract leaders (rank 1) and runners-up (rank 2) per year
year_leaders AS (
    SELECT 
        Year,
        argMax(Reporting_Airline, RankInYear = 1) AS LeaderReportingAirline,
        argMax(SharePct, RankInYear = 1) AS LeaderSharePct,
        argMax(Reporting_Airline, RankInYear = 2) AS RunnerUpReportingAirline,
        argMax(SharePct, RankInYear = 2) AS RunnerUpSharePct
    FROM ranked_carriers
    WHERE RankInYear IN (1, 2)
    GROUP BY Year
),

-- Step 5: Compute leader transitions with YoY share change
leader_transitions AS (
    SELECT 
        Year,
        LeaderReportingAirline,
        RunnerUpReportingAirline,
        round(LeaderSharePct - RunnerUpSharePct, 4) AS LeaderShareGapPctPts,
        LeaderSharePct,
        lagInFrame(LeaderReportingAirline) OVER (ORDER BY Year) AS PriorYearLeaderReportingAirline,
        lagInFrame(LeaderSharePct) OVER (ORDER BY Year) AS PriorYearLeaderSharePct,
        row_number() OVER (ORDER BY Year) AS YearSeq
    FROM year_leaders
),

-- Step 6: Final leader transition metrics
leader_metrics AS (
    SELECT 
        Year,
        LeaderReportingAirline,
        RunnerUpReportingAirline,
        LeaderShareGapPctPts,
        LeaderSharePct,
        PriorYearLeaderReportingAirline,
        CASE WHEN YearSeq = 1 THEN NULL ELSE round(LeaderSharePct - PriorYearLeaderSharePct, 4) END AS LeaderShareChangePctPts,
        CASE 
            WHEN YearSeq = 1 THEN 0
            WHEN LeaderReportingAirline != PriorYearLeaderReportingAirline THEN 1
            ELSE 0
        END AS LeaderChanged
    FROM leader_transitions
),

-- Step 7: Top 5 carriers per year for bump chart
top5_carriers AS (
    SELECT *
    FROM ranked_carriers
    WHERE RankInYear <= 5
)

-- Final output combining carrier_year rows
SELECT 
    'carrier_year' AS RowType,
    tc.Year,
    tc.Reporting_Airline,
    tc.RankInYear,
    tc.CompletedFlights,
    tc.SharePct,
    lm.LeaderReportingAirline,
    lm.RunnerUpReportingAirline,
    lm.LeaderShareGapPctPts,
    lm.PriorYearLeaderReportingAirline,
    lm.LeaderChanged,
    lm.LeaderShareChangePctPts
FROM top5_carriers tc
JOIN leader_metrics lm ON tc.Year = lm.Year
ORDER BY 
    tc.Year ASC,
    RowType ASC,
    tc.RankInYear ASC,
    tc.Reporting_Airline ASC
