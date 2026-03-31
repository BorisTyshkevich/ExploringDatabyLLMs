WITH itineraries AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline AS Carrier,
        FlightDate,
        count() AS hops,
        arrayStringConcat(
            arrayConcat(
                arraySort((x, t) -> t, groupArray(OriginCode), groupArray(assumeNotNull(CRSDepTime))),
                [argMax(DestCode, assumeNotNull(CRSDepTime))]
            ),
            '-'
        ) AS Route
    FROM ontime.fact_ontime
    WHERE Tail_Number != '' AND Cancelled = 0
    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate
),
best_by_route AS (
    SELECT
        Route,
        max(hops) AS max_hops,
        argMax(Tail_Number, FlightDate) AS aircraft_id,
        argMax(Flight_Number_Reporting_Airline, FlightDate) AS flight_number,
        argMax(Carrier, FlightDate) AS carrier,
        max(FlightDate) AS flight_date,
        argMax(hops, FlightDate) AS latest_hops
    FROM itineraries
    GROUP BY Route
)
SELECT
    aircraft_id,
    flight_number,
    carrier,
    flight_date,
    max_hops AS hop_count,
    latest_hops AS num_flights,
    Route
FROM best_by_route
ORDER BY max_hops DESC, flight_date DESC
LIMIT 10
