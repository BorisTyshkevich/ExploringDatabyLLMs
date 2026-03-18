CREATE TABLE IF NOT EXISTS ontime.airports
(
    `code` FixedString(3) COMMENT 'Three-letter airport code from the BTS Master Coordinate export.',
    `airport_id` UInt32 COMMENT 'Stable DOT airport identifier used by ontime.ontime OriginAirportID and DestAirportID.',
    `airport_seq_id` UInt32 COMMENT 'Time-specific DOT airport sequence identifier for this version of the airport record.',
    `name` String COMMENT 'Airport display name from the BTS Master Coordinate export.',
    `city_name` String COMMENT 'Full city name associated with the airport record in the BTS export.',
    `city_market_id` UInt32 COMMENT 'DOT city market identifier used to group airports serving the same market. Zero means missing in the source export.',
    `city_market_name` String COMMENT 'Full city market display name from the BTS export.',
    `wac` UInt16 COMMENT 'World area code for the airport. Zero means missing in the source export.',
    `country_name` String COMMENT 'Country name for the airport record.',
    `country_code_iso` String COMMENT 'ISO country code for the airport record.',
    `state_name` String COMMENT 'State or province name for the airport record.',
    `state_code` String COMMENT 'State or province code for the airport record.',
    `state_fips` String COMMENT 'State FIPS code for the airport record when present.',
    `latitude` Float64 COMMENT 'Airport latitude in decimal degrees. Zero means missing in the source export.',
    `longitude` Float64 COMMENT 'Airport longitude in decimal degrees. Zero means missing in the source export.',
    `utc_local_time_variation` String COMMENT 'UTC offset string from BTS, for example -0500.',
    `start_date` Date COMMENT 'Start date when this airport record version became effective. 1970-01-01 means missing in the source export.',
    `thru_date` Date COMMENT 'End date when this airport record version stopped being effective. 1970-01-01 means missing in the source export.',
    `is_closed` UInt8 COMMENT 'Closure flag from BTS where 1 means the airport record is closed.',
    `is_latest` UInt8 COMMENT 'Latest-record flag from BTS where 1 marks the current record for the airport code.'
)
ENGINE = MergeTree
ORDER BY (`code`, `airport_id`, `is_latest`)
COMMENT 'Simplified airport dimension modeled from the BTS Master Coordinate history export.'
