WITH itineraries AS (
  SELECT
    Tail_Number,
    IATA_CODE_Reporting_Airline AS Carrier,
    Flight_Number_Reporting_Airline AS FlightNum,
    FlightDate,
    countDistinct(OriginCode, DestCode, CRSDepTime) AS hop_count,
    arrayStringConcat(
      arrayConcat(
        arrayMap(t -> t.2, arraySort(t -> t.1, groupUniqArray((coalesce(CRSDepTime, 0), OriginCode)))),
        [argMax(DestCode, coalesce(CRSDepTime, 0))]
      ),
      '-'
    ) AS Route
  FROM ontime.fact_ontime
  GROUP BY Tail_Number, Carrier, FlightNum, FlightDate
),
max_hop AS (SELECT max(hop_count) AS max_hop_count FROM itineraries),
route_recurrence AS (
  SELECT Route, countDistinct(FlightDate) AS recurrence_count
  FROM itineraries GROUP BY Route
),
top_itineraries AS (
  SELECT i.Route, i.FlightDate, i.hop_count
  FROM itineraries i CROSS JOIN max_hop WHERE i.hop_count = max_hop.max_hop_count
),
top10_routes AS (
  SELECT t.Route, r.recurrence_count
  FROM top_itineraries t JOIN route_recurrence r ON t.Route = r.Route
  GROUP BY t.Route, r.recurrence_count
  ORDER BY max(t.FlightDate) DESC
  LIMIT 10
)
SELECT Route, recurrence_count
FROM top10_routes
ORDER BY recurrence_count DESC
