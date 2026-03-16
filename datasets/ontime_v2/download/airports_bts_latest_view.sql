CREATE VIEW IF NOT EXISTS default.airports_bts_latest AS
SELECT *
FROM default.airports_bts
WHERE AIRPORT_IS_LATEST = 1
