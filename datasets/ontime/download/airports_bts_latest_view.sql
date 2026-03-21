CREATE OR REPLACE VIEW ontime.dim_airports AS
SELECT *
FROM ontime.dim_airports_bts_full
WHERE IsLatest = 1
