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
max_hop AS (
  SELECT max(hop_count) AS max_hop_count FROM itineraries
),
route_recurrence AS (
  SELECT Route, countDistinct(FlightDate) AS recurrence_count
  FROM itineraries GROUP BY Route
),
top_itineraries AS (
  SELECT i.Tail_Number, i.Carrier, i.FlightNum, i.FlightDate, i.hop_count, i.Route
  FROM itineraries i CROSS JOIN max_hop WHERE i.hop_count = max_hop.max_hop_count
),
top_with_recurrence AS (
  SELECT t.Tail_Number, t.Carrier, t.FlightNum, t.FlightDate, t.hop_count, r.recurrence_count, t.Route
  FROM top_itineraries t JOIN route_recurrence r ON t.Route = r.Route
)
SELECT
  argMax(Tail_Number, FlightDate) AS aircraft_id,
  argMax(FlightNum, FlightDate) AS flight_number,
  argMax(Carrier, FlightDate) AS carrier,
  max(FlightDate) AS flight_date,
  max(hop_count) AS hop_count,
  any(recurrence_count) AS recurrence_count,
  Route
FROM top_with_recurrence
GROUP BY Route
ORDER BY flight_date DESC
LIMIT 10
