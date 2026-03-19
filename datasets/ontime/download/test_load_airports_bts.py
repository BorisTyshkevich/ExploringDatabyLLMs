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
        self.assertEqual(got["code"], "JFK")
        self.assertEqual(got["airport_id"], "1001")
        self.assertEqual(got["city_market_id"], "31703")
        self.assertEqual(got["latitude"], "40.6413")
        self.assertEqual(got["start_date"], "2007-07-01")
        self.assertEqual(got["thru_date"], r"\N")

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
        self.assertEqual(load_airports_bts.serialize_target_value("city_market_id", r"\N"), "0")
        self.assertEqual(load_airports_bts.serialize_target_value("latitude", r"\N"), "0")
        self.assertEqual(load_airports_bts.serialize_target_value("thru_date", r"\N"), "1970-01-01")

    def test_apply_semantic_latest_picks_single_open_latest_row_per_code(self) -> None:
        rows = [
            {
                "code": "AUS",
                "airport_id": "16440",
                "airport_seq_id": "1",
                "name": "Robert Mueller Municipal",
                "city_name": "Austin, TX",
                "city_market_id": "1",
                "city_market_name": "Austin, TX",
                "wac": "1",
                "country_name": "United States",
                "country_code_iso": "US",
                "state_name": "Texas",
                "state_code": "TX",
                "state_fips": "48",
                "latitude": "30.3",
                "longitude": "-97.7",
                "utc_local_time_variation": "-0600",
                "start_date": "1970-01-01",
                "thru_date": "1999-05-31",
                "is_closed": "1",
                "is_latest": "1",
            },
            {
                "code": "AUS",
                "airport_id": "10423",
                "airport_seq_id": "2",
                "name": "Austin - Bergstrom International",
                "city_name": "Austin, TX",
                "city_market_id": "1",
                "city_market_name": "Austin, TX",
                "wac": "1",
                "country_name": "United States",
                "country_code_iso": "US",
                "state_name": "Texas",
                "state_code": "TX",
                "state_fips": "48",
                "latitude": "30.2",
                "longitude": "-97.6",
                "utc_local_time_variation": "-0600",
                "start_date": "1999-06-01",
                "thru_date": "1970-01-01",
                "is_closed": "0",
                "is_latest": "1",
            },
        ]

        got = load_airports_bts.apply_semantic_latest(rows)

        self.assertEqual([row["is_latest"] for row in got], ["0", "1"])

    def test_apply_semantic_latest_uses_newer_start_date_then_airport_id_as_tiebreaker(self) -> None:
        rows = [
            {
                "code": "ABC",
                "airport_id": "100",
                "airport_seq_id": "1",
                "name": "Older Open Airport",
                "city_name": "",
                "city_market_id": "0",
                "city_market_name": "",
                "wac": "0",
                "country_name": "",
                "country_code_iso": "",
                "state_name": "",
                "state_code": "",
                "state_fips": "",
                "latitude": "0",
                "longitude": "0",
                "utc_local_time_variation": "",
                "start_date": "2000-01-01",
                "thru_date": "1970-01-01",
                "is_closed": "0",
                "is_latest": "1",
            },
            {
                "code": "ABC",
                "airport_id": "200",
                "airport_seq_id": "2",
                "name": "Newer Open Airport",
                "city_name": "",
                "city_market_id": "0",
                "city_market_name": "",
                "wac": "0",
                "country_name": "",
                "country_code_iso": "",
                "state_name": "",
                "state_code": "",
                "state_fips": "",
                "latitude": "0",
                "longitude": "0",
                "utc_local_time_variation": "",
                "start_date": "2010-01-01",
                "thru_date": "1970-01-01",
                "is_closed": "0",
                "is_latest": "0",
            },
            {
                "code": "XYZ",
                "airport_id": "300",
                "airport_seq_id": "3",
                "name": "Lower Airport ID",
                "city_name": "",
                "city_market_id": "0",
                "city_market_name": "",
                "wac": "0",
                "country_name": "",
                "country_code_iso": "",
                "state_name": "",
                "state_code": "",
                "state_fips": "",
                "latitude": "0",
                "longitude": "0",
                "utc_local_time_variation": "",
                "start_date": "2010-01-01",
                "thru_date": "1970-01-01",
                "is_closed": "0",
                "is_latest": "1",
            },
            {
                "code": "XYZ",
                "airport_id": "400",
                "airport_seq_id": "4",
                "name": "Higher Airport ID",
                "city_name": "",
                "city_market_id": "0",
                "city_market_name": "",
                "wac": "0",
                "country_name": "",
                "country_code_iso": "",
                "state_name": "",
                "state_code": "",
                "state_fips": "",
                "latitude": "0",
                "longitude": "0",
                "utc_local_time_variation": "",
                "start_date": "2010-01-01",
                "thru_date": "1970-01-01",
                "is_closed": "0",
                "is_latest": "0",
            },
        ]

        got = load_airports_bts.apply_semantic_latest(rows)

        self.assertEqual([row["is_latest"] for row in got if row["code"] == "ABC"], ["0", "1"])
        self.assertEqual([row["is_latest"] for row in got if row["code"] == "XYZ"], ["0", "1"])


if __name__ == "__main__":
    unittest.main()
