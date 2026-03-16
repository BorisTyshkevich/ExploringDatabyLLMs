CREATE TABLE IF NOT EXISTS default.airports_bts
(
    `AIRPORT_SEQ_ID` UInt32,
    `AIRPORT_ID` UInt32,
    `AIRPORT` String,
    `DISPLAY_AIRPORT_NAME` String,
    `DISPLAY_AIRPORT_CITY_NAME_FULL` String,
    `AIRPORT_WAC_SEQ_ID2` Nullable(UInt32),
    `AIRPORT_WAC` Nullable(UInt16),
    `AIRPORT_COUNTRY_NAME` String,
    `AIRPORT_COUNTRY_CODE_ISO` String,
    `AIRPORT_STATE_NAME` String,
    `AIRPORT_STATE_CODE` String,
    `AIRPORT_STATE_FIPS` String,
    `CITY_MARKET_SEQ_ID` Nullable(UInt32),
    `CITY_MARKET_ID` Nullable(UInt32),
    `DISPLAY_CITY_MARKET_NAME_FULL` String,
    `CITY_MARKET_WAC_SEQ_ID2` Nullable(UInt32),
    `CITY_MARKET_WAC` Nullable(UInt16),
    `LAT_DEGREES` Nullable(Int16),
    `LAT_HEMISPHERE` String,
    `LAT_MINUTES` Nullable(Int16),
    `LAT_SECONDS` Nullable(Int16),
    `LATITUDE` Nullable(Float64),
    `LON_DEGREES` Nullable(Int16),
    `LON_HEMISPHERE` String,
    `LON_MINUTES` Nullable(Int16),
    `LON_SECONDS` Nullable(Int16),
    `LONGITUDE` Nullable(Float64),
    `UTC_LOCAL_TIME_VARIATION` String,
    `AIRPORT_START_DATE` Nullable(Date),
    `AIRPORT_THRU_DATE` Nullable(Date),
    `AIRPORT_IS_CLOSED` UInt8,
    `AIRPORT_IS_LATEST` UInt8
)
ENGINE = MergeTree
ORDER BY (`AIRPORT_ID`, `AIRPORT_SEQ_ID`, ifNull(`AIRPORT_START_DATE`, toDate('1900-01-01')))
SETTINGS index_granularity = 8192
