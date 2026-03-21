import unittest

import load_airports_bts


class LoadAirportsBTSTest(unittest.TestCase):
    def test_extract_hidden_fields(self) -> None:
        html = """
        <input type="hidden" name="__VIEWSTATE" value="state" />
        <input type="hidden" name="__VIEWSTATEGENERATOR" value="gen" />
        <input type="hidden" name="__EVENTVALIDATION" value="event" />
        """
        self.assertEqual(
            load_airports_bts.extract_hidden_fields(html),
            {
                "__VIEWSTATE": "state",
                "__VIEWSTATEGENERATOR": "gen",
                "__EVENTVALIDATION": "event",
            },
        )

    def test_extract_hidden_fields_rejects_missing_values(self) -> None:
        with self.assertRaises(RuntimeError):
            load_airports_bts.extract_hidden_fields('<input type="hidden" name="__VIEWSTATE" value="state" />')

    def test_extract_checkbox_names_matches_expected_export_columns(self) -> None:
        html = "".join(
            f'<input id="{name}" type="checkbox" name="{name}" checked="checked" />'
            for name in load_airports_bts.EXPECTED_COLUMNS
        )
        self.assertEqual(load_airports_bts.extract_checkbox_names(html), load_airports_bts.EXPECTED_COLUMNS)

    def test_normalize_value_formats_dates(self) -> None:
        self.assertEqual(load_airports_bts.normalize_value("AIRPORT_START_DATE", "7/1/2007 12:00:00 AM"), "2007-07-01")

    def test_transform_row_projects_modeled_schema(self) -> None:
        got = load_airports_bts.transform_row(
            {
                "AIRPORT_SEQ_ID": "10001",
                "AIRPORT_ID": "1001",
                "AIRPORT": "JFK",
                "DISPLAY_AIRPORT_NAME": "John F Kennedy International",
                "DISPLAY_AIRPORT_CITY_NAME_FULL": "New York, NY",
                "AIRPORT_WAC": "1",
                "AIRPORT_COUNTRY_NAME": "United States",
                "AIRPORT_COUNTRY_CODE_ISO": "US",
                "AIRPORT_STATE_NAME": "New York",
                "AIRPORT_STATE_CODE": "NY",
                "AIRPORT_STATE_FIPS": "36",
                "CITY_MARKET_ID": "31703",
                "DISPLAY_CITY_MARKET_NAME_FULL": "New York, NY",
                "LATITUDE": "40.6413",
                "LONGITUDE": "-73.7781",
                "UTC_LOCAL_TIME_VARIATION": "-0500",
                "AIRPORT_START_DATE": "7/1/2007 12:00:00 AM",
                "AIRPORT_THRU_DATE": "",
                "AIRPORT_IS_CLOSED": "0",
                "AIRPORT_IS_LATEST": "1",
            }
        )
        self.assertEqual(got["AirportCode"], "JFK")
        self.assertEqual(got["AirportID"], "1001")
        self.assertEqual(got["CityMarketID"], "31703")
        self.assertEqual(got["Latitude"], "40.6413")
        self.assertEqual(got["StartDate"], "2007-07-01")
        self.assertEqual(got["ThruDate"], r"\N")

    def test_transform_row_rejects_non_iata_length_code(self) -> None:
        with self.assertRaises(RuntimeError):
            load_airports_bts.transform_row(
                {
                    "AIRPORT_SEQ_ID": "1",
                    "AIRPORT_ID": "1",
                    "AIRPORT": "ABCD",
                    "DISPLAY_AIRPORT_NAME": "Bad",
                    "DISPLAY_AIRPORT_CITY_NAME_FULL": "Bad City",
                    "AIRPORT_WAC": "",
                    "AIRPORT_COUNTRY_NAME": "",
                    "AIRPORT_COUNTRY_CODE_ISO": "",
                    "AIRPORT_STATE_NAME": "",
                    "AIRPORT_STATE_CODE": "",
                    "AIRPORT_STATE_FIPS": "",
                    "CITY_MARKET_ID": "",
                    "DISPLAY_CITY_MARKET_NAME_FULL": "",
                    "LATITUDE": "",
                    "LONGITUDE": "",
                    "UTC_LOCAL_TIME_VARIATION": "",
                    "AIRPORT_START_DATE": "",
                    "AIRPORT_THRU_DATE": "",
                    "AIRPORT_IS_CLOSED": "0",
                    "AIRPORT_IS_LATEST": "1",
                }
            )

    def test_serialize_target_value_uses_zero_defaults_for_missing_modeled_values(self) -> None:
        self.assertEqual(load_airports_bts.serialize_target_value("CityMarketID", r"\N"), "0")
        self.assertEqual(load_airports_bts.serialize_target_value("Latitude", r"\N"), "0")
        self.assertEqual(load_airports_bts.serialize_target_value("ThruDate", r"\N"), "1970-01-01")

    def test_apply_semantic_latest_picks_single_open_latest_row_per_code(self) -> None:
        rows = [
            {
                "AirportCode": "AUS",
                "AirportID": "16440",
                "AirportSeqID": "1",
                "DisplayAirportName": "Robert Mueller Municipal",
                "CityName": "Austin, TX",
                "CityMarketID": "1",
                "CityMarketName": "Austin, TX",
                "WorldAreaCode": "1",
                "CountryName": "United States",
                "CountryCodeISO": "US",
                "StateName": "Texas",
                "StateCode": "TX",
                "StateFips": "48",
                "Latitude": "30.3",
                "Longitude": "-97.7",
                "UtcLocalTimeVariation": "-0600",
                "StartDate": "1970-01-01",
                "ThruDate": "1999-05-31",
                "IsClosed": "1",
                "IsLatest": "1",
            },
            {
                "AirportCode": "AUS",
                "AirportID": "10423",
                "AirportSeqID": "2",
                "DisplayAirportName": "Austin - Bergstrom International",
                "CityName": "Austin, TX",
                "CityMarketID": "1",
                "CityMarketName": "Austin, TX",
                "WorldAreaCode": "1",
                "CountryName": "United States",
                "CountryCodeISO": "US",
                "StateName": "Texas",
                "StateCode": "TX",
                "StateFips": "48",
                "Latitude": "30.2",
                "Longitude": "-97.6",
                "UtcLocalTimeVariation": "-0600",
                "StartDate": "1999-06-01",
                "ThruDate": "1970-01-01",
                "IsClosed": "0",
                "IsLatest": "1",
            },
        ]

        got = load_airports_bts.apply_semantic_latest(rows)

        self.assertEqual([row["IsLatest"] for row in got], ["0", "1"])

    def test_apply_semantic_latest_uses_newer_start_date_then_airport_id_as_tiebreaker(self) -> None:
        rows = [
            {
                "AirportCode": "ABC",
                "AirportID": "100",
                "AirportSeqID": "1",
                "DisplayAirportName": "Older Open Airport",
                "CityName": "",
                "CityMarketID": "0",
                "CityMarketName": "",
                "WorldAreaCode": "0",
                "CountryName": "",
                "CountryCodeISO": "",
                "StateName": "",
                "StateCode": "",
                "StateFips": "",
                "Latitude": "0",
                "Longitude": "0",
                "UtcLocalTimeVariation": "",
                "StartDate": "2000-01-01",
                "ThruDate": "1970-01-01",
                "IsClosed": "0",
                "IsLatest": "1",
            },
            {
                "AirportCode": "ABC",
                "AirportID": "200",
                "AirportSeqID": "2",
                "DisplayAirportName": "Newer Open Airport",
                "CityName": "",
                "CityMarketID": "0",
                "CityMarketName": "",
                "WorldAreaCode": "0",
                "CountryName": "",
                "CountryCodeISO": "",
                "StateName": "",
                "StateCode": "",
                "StateFips": "",
                "Latitude": "0",
                "Longitude": "0",
                "UtcLocalTimeVariation": "",
                "StartDate": "2010-01-01",
                "ThruDate": "1970-01-01",
                "IsClosed": "0",
                "IsLatest": "0",
            },
            {
                "AirportCode": "XYZ",
                "AirportID": "300",
                "AirportSeqID": "3",
                "DisplayAirportName": "Lower Airport ID",
                "CityName": "",
                "CityMarketID": "0",
                "CityMarketName": "",
                "WorldAreaCode": "0",
                "CountryName": "",
                "CountryCodeISO": "",
                "StateName": "",
                "StateCode": "",
                "StateFips": "",
                "Latitude": "0",
                "Longitude": "0",
                "UtcLocalTimeVariation": "",
                "StartDate": "2010-01-01",
                "ThruDate": "1970-01-01",
                "IsClosed": "0",
                "IsLatest": "1",
            },
            {
                "AirportCode": "XYZ",
                "AirportID": "400",
                "AirportSeqID": "4",
                "DisplayAirportName": "Higher Airport ID",
                "CityName": "",
                "CityMarketID": "0",
                "CityMarketName": "",
                "WorldAreaCode": "0",
                "CountryName": "",
                "CountryCodeISO": "",
                "StateName": "",
                "StateCode": "",
                "StateFips": "",
                "Latitude": "0",
                "Longitude": "0",
                "UtcLocalTimeVariation": "",
                "StartDate": "2010-01-01",
                "ThruDate": "1970-01-01",
                "IsClosed": "0",
                "IsLatest": "0",
            },
        ]

        got = load_airports_bts.apply_semantic_latest(rows)

        self.assertEqual([row["IsLatest"] for row in got if row["AirportCode"] == "ABC"], ["0", "1"])
        self.assertEqual([row["IsLatest"] for row in got if row["AirportCode"] == "XYZ"], ["0", "1"])


if __name__ == "__main__":
    unittest.main()
