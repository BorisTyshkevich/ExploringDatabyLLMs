WITH legs AS (
    SELECT
        o.Tail_Number,
        o.Flight_Number_Reporting_Airline,
        o.IATA_CODE_Reporting_Airline,
        o.FlightDate,
        assumeNotNull(o.DepTime) AS DepTime,
        o.Origin,
        o.Dest,
        coalesce(ao.name, o.Origin) AS OriginName,
        coalesce(ad.name, o.Dest) AS DestName
    FROM ontime.ontime AS o
    LEFT JOIN ontime.airports_latest AS ao ON o.Origin = ao.code
    LEFT JOIN ontime.airports_latest AS ad ON o.Dest = ad.code
    WHERE o.Tail_Number != ''
      AND o.DepTime IS NOT NULL
      AND o.Cancelled = 0
),
itineraries AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        count() AS Hops,
        arraySort(x -> x.1, groupArray(tuple(DepTime, Origin, Dest, OriginName, DestName))) AS legs_sorted
    FROM legs
    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate
    HAVING Hops >= 6
)
SELECT
    Tail_Number,
    Flight_Number_Reporting_Airline,
    IATA_CODE_Reporting_Airline,
    FlightDate,
    Hops,
    arrayStringConcat(
        arrayConcat(
            arrayMap(
                x -> concat(
                    x.2, ' (', x.4, ') ',
                    lpad(toString(x.1 DIV 100), 2, '0'), ':',
                    lpad(toString(x.1 MOD 100), 2, '0')
                ),
                legs_sorted
            ),
            [concat(legs_sorted[length(legs_sorted)].3, ' (', legs_sorted[length(legs_sorted)].5, ')')]
        ),
        ' -> '
    ) AS Route
FROM itineraries
ORDER BY Hops DESC, FlightDate DESC
LIMIT 10
