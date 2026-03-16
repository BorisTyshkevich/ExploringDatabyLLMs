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


if __name__ == "__main__":
    unittest.main()
