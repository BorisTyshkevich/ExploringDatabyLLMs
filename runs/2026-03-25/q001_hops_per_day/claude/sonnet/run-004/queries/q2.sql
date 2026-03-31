SELECT
    multiIf(d.Longitude < -100, 'West', d.Longitude < -85, 'Central', 'East') AS region,
    count() AS airport_count,
    groupArray(a.code) AS airports
FROM (
    SELECT DISTINCT arrayJoin([
        'ISP','BWI','MYR','BNA','VPS','DAL','LAS','OAK','SEA',
        'CLE','PNS','HOU','MCI','PHX','BUR','DEN',
        'ELP','LIT','ATL','RIC','MDW','SAN',
        'MSY','CMH','RDU','DTW','LAX',
        'TPA','ORD','SLC','SJC',
        'LGA','STL','ICT','COS','JAN',
        'SMF','PSP','RNO',
        'ABQ','MAF',
        'FLL','MCO','MEM','IAD'
    ]) AS code
) a
LEFT JOIN ontime.dim_airports d ON a.code = d.AirportCode
GROUP BY region
ORDER BY min(d.Longitude)
