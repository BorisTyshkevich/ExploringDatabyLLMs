WITH legs AS (
    SELECT
        o.Tail_Number,
        o.Flight_Number_Reporting_Airline,
        o.IATA_CODE_Reporting_Airline,
        o.FlightDate,
        toDateTime(
            concat(
                toString(o.FlightDate),
                ' ',
                leftPad(toString(intDiv(o.DepTime, 100)), 2, '0'),
                ':',
                leftPad(toString(o.DepTime % 100), 2, '0'),
                ':00'
            )
        ) AS dep_ts,
        concat(
            leftPad(toString(intDiv(o.DepTime, 100)), 2, '0'),
            ':',
            leftPad(toString(o.DepTime % 100), 2, '0')
        ) AS dep_hhmm,
        o.Origin,
        o.Dest
    FROM ontime.ontime AS o
    WHERE o.Cancelled = 0
      AND o.Diverted = 0
      AND o.DepTime IS NOT NULL
      AND o.Tail_Number != ''
      AND o.Flight_Number_Reporting_Airline != ''
),
itineraries AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        count() AS Hops,
        arraySort(
            x -> (x.1, x.2, x.3),
            groupArray((dep_ts, toString(Origin), toString(Dest), dep_hhmm))
        ) AS ordered_legs,
        arrayStringConcat(
            arrayConcat([ordered_legs[1].2], arrayMap(x -> x.3, ordered_legs)),
            ' -> '
        ) AS Route,
        arrayStringConcat(
            arrayMap(x -> concat(x.4, ' ', x.2), ordered_legs),
            ' | '
        ) AS `Departure Times From Origin`
    FROM legs
    GROUP BY
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate
),
max_hops AS (
    SELECT max(Hops) AS value
    FROM itineraries
),
scored AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        Hops,
        Route,
        `Departure Times From Origin`,
        (SELECT value FROM max_hops) AS `Max Hops Overall`,
        (SELECT count() FROM itineraries WHERE Hops = (SELECT value FROM max_hops)) AS `Max Hop Itinerary Count`,
        count() OVER (PARTITION BY Hops, Route) AS `Same-Hops Route Count`
    FROM itineraries
),
top10 AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        Hops,
        Route,
        `Departure Times From Origin`,
        `Max Hops Overall`,
        `Max Hop Itinerary Count`,
        `Same-Hops Route Count`
    FROM scored
    ORDER BY Hops DESC, FlightDate DESC, Tail_Number DESC
    LIMIT 10
)
SELECT
    Tail_Number AS `Aircraft ID`,
    Flight_Number_Reporting_Airline AS `Flight Number`,
    IATA_CODE_Reporting_Airline AS Carrier,
    FlightDate AS Date,
    Hops,
    Route,
    `Departure Times From Origin`,
    `Max Hops Overall`,
    `Max Hop Itinerary Count`,
    `Same-Hops Route Count`,
    count() OVER (PARTITION BY Route) AS `Route Frequency In Top 10`
FROM top10
ORDER BY Hops DESC, Date DESC, `Aircraft ID` DESC
