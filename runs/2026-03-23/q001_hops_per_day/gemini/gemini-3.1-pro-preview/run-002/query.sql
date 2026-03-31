SELECT
    Tail_Number AS `Aircraft ID`,
    Flight_Number_Reporting_Airline AS `Flight Number`,
    Reporting_Airline AS `Carrier`,
    FlightDate AS `Date`,
    arrayStringConcat(
        arrayPushBack(
            arrayMap(x -> x.2, arraySort(groupArray((DepTime, OriginCode)))),
            argMax(DestCode, DepTime)
        ),
        '-'
    ) AS Route,
    count() AS hop_count
FROM ontime.fact_ontime
WHERE Cancelled = 0 
  AND Tail_Number != ''
GROUP BY
    `Aircraft ID`,
    `Flight Number`,
    `Carrier`,
    `Date`
ORDER BY
    hop_count DESC,
    `Date` DESC
LIMIT 10