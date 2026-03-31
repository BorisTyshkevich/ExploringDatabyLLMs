SELECT FlightDate                                                AS Date, Tail_Number AS AircraftID,
       Flight_Number_Reporting_Airline                           AS FlightNumber, Reporting_Airline AS Carrier,
       count()                                                   AS Hops,
       arrayStringConcat(arrayMap(x -> x.2 || ' (' || leftPad(toString(x.1), 4, '0') || ')',
                                  arraySort(groupArray(tuple(DepTime, Origin, Dest)))), ' - ') || ' - ' ||
       arraySort(groupArray(tuple(DepTime, Origin, Dest)))[-1].3 AS Route
FROM ontime.ontime
WHERE Cancelled = 0 AND Tail_Number != '' AND Tail_Number != 'UNKNOW' AND DepTime IS NOT NULL
GROUP BY Date, AircraftID, FlightNumber, Carrier
ORDER BY Hops DESC, Date DESC
LIMIT 10
