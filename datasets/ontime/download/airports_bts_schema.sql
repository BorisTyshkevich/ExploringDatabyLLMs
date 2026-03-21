CREATE TABLE IF NOT EXISTS ontime.dim_airports_bts_full
(
    `AirportCode` FixedString(3) COMMENT 'Three-letter airport code from the BTS Master Coordinate export.',
    `AirportID` UInt32 COMMENT 'Stable DOT airport identifier used by ontime.fact_ontime OriginAirportID and DestAirportID.',
    `AirportSeqID` UInt32 COMMENT 'Time-specific DOT airport sequence identifier for this version of the airport record.',
    `DisplayAirportName` String COMMENT 'Airport display name from the BTS Master Coordinate export.',
    `CityName` String COMMENT 'Full city name associated with the airport record in the BTS export.',
    `CityMarketID` UInt32 COMMENT 'DOT city market identifier used to group airports serving the same market. Zero means missing in the source export.',
    `CityMarketName` String COMMENT 'Full city market display name from the BTS export.',
    `WorldAreaCode` UInt16 COMMENT 'World area code for the airport. Zero means missing in the source export.',
    `CountryName` String COMMENT 'Country name for the airport record.',
    `CountryCodeISO` String COMMENT 'ISO country code for the airport record.',
    `StateName` String COMMENT 'State or province name for the airport record.',
    `StateCode` String COMMENT 'State or province code for the airport record.',
    `StateFips` String COMMENT 'State FIPS code for the airport record when present.',
    `Latitude` Float64 COMMENT 'Airport latitude in decimal degrees. Zero means missing in the source export.',
    `Longitude` Float64 COMMENT 'Airport longitude in decimal degrees. Zero means missing in the source export.',
    `UtcLocalTimeVariation` String COMMENT 'UTC offset string from BTS, for example -0500.',
    `StartDate` Date COMMENT 'Start date when this airport record version became effective. 1970-01-01 means missing in the source export.',
    `ThruDate` Date COMMENT 'End date when this airport record version stopped being effective. 1970-01-01 means missing in the source export.',
    `IsClosed` UInt8 COMMENT 'Closure flag from BTS where 1 means the airport record is closed.',
    `IsLatest` UInt8 COMMENT 'Loader-computed latest-record flag where 1 marks the single selected current row for an airport code.'
)
ENGINE = MergeTree
ORDER BY (`AirportCode`, `AirportID`, `IsLatest`)
COMMENT 'Simplified airport dimension modeled from the BTS Master Coordinate history export, with a cleaned latest flag per code.'
