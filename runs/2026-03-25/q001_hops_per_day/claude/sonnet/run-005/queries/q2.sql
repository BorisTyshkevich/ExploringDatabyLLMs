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
top_itineraries AS (
  SELECT i.Route, i.FlightDate, i.hop_count
  FROM itineraries i CROSS JOIN max_hop WHERE i.hop_count = max_hop.max_hop_count
),
top10_routes AS (
  SELECT Route
  FROM top_itineraries
  GROUP BY Route
  ORDER BY max(FlightDate) DESC
  LIMIT 10
),
all_airports AS (
  SELECT DISTINCT arrayJoin(splitByChar('-', Route)) AS ap
  FROM top10_routes
)
SELECT
  a.ap AS airport_code,
  d.DisplayAirportName AS airport_name,
  d.Latitude AS lat,
  d.Longitude AS lon
FROM all_airports a
LEFT JOIN ontime.dim_airports d ON a.ap = d.AirportCode
ORDER BY airport_code
