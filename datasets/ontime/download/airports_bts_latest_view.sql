CREATE OR REPLACE VIEW ontime.airports_latest AS
SELECT *
FROM ontime.airports_bts
WHERE is_latest = 1
