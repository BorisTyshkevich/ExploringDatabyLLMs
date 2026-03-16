#!/usr/bin/env python3

from __future__ import annotations

import argparse
import csv
import http.cookiejar
import io
import json
import re
import subprocess
import sys
import urllib.error
import urllib.parse
import urllib.request
import zipfile
from dataclasses import dataclass
from datetime import UTC, datetime
from decimal import Decimal, InvalidOperation
from pathlib import Path


BASE_DOWNLOAD_URL = "https://www.transtats.bts.gov/DL_SelectFields.aspx?gnoyr_VQ=FLL&QO_fu146_anzr=N8vn6v10+f722146+gnoyr5"
BASE_INFO_URL = "https://www.transtats.bts.gov/Fields.asp?gnoyr_VQ=FLL"
SCRIPT_DIR = Path(__file__).resolve().parent
CACHE_DIR = SCRIPT_DIR / ".cache" / "airports_bts"
DOWNLOAD_DIR = CACHE_DIR / "downloads"
META_DIR = CACHE_DIR / "meta"
DEBUG = False


EXPECTED_COLUMNS = [
    "AIRPORT_SEQ_ID",
    "AIRPORT_ID",
    "AIRPORT",
    "DISPLAY_AIRPORT_NAME",
    "DISPLAY_AIRPORT_CITY_NAME_FULL",
    "AIRPORT_WAC_SEQ_ID2",
    "AIRPORT_WAC",
    "AIRPORT_COUNTRY_NAME",
    "AIRPORT_COUNTRY_CODE_ISO",
    "AIRPORT_STATE_NAME",
    "AIRPORT_STATE_CODE",
    "AIRPORT_STATE_FIPS",
    "CITY_MARKET_SEQ_ID",
    "CITY_MARKET_ID",
    "DISPLAY_CITY_MARKET_NAME_FULL",
    "CITY_MARKET_WAC_SEQ_ID2",
    "CITY_MARKET_WAC",
    "LAT_DEGREES",
    "LAT_HEMISPHERE",
    "LAT_MINUTES",
    "LAT_SECONDS",
    "LATITUDE",
    "LON_DEGREES",
    "LON_HEMISPHERE",
    "LON_MINUTES",
    "LON_SECONDS",
    "LONGITUDE",
    "UTC_LOCAL_TIME_VARIATION",
    "AIRPORT_START_DATE",
    "AIRPORT_THRU_DATE",
    "AIRPORT_IS_CLOSED",
    "AIRPORT_IS_LATEST",
]

STRING_COLUMNS = {
    "AIRPORT",
    "DISPLAY_AIRPORT_NAME",
    "DISPLAY_AIRPORT_CITY_NAME_FULL",
    "AIRPORT_COUNTRY_NAME",
    "AIRPORT_COUNTRY_CODE_ISO",
    "AIRPORT_STATE_NAME",
    "AIRPORT_STATE_CODE",
    "AIRPORT_STATE_FIPS",
    "DISPLAY_CITY_MARKET_NAME_FULL",
    "LAT_HEMISPHERE",
    "LON_HEMISPHERE",
    "UTC_LOCAL_TIME_VARIATION",
}

UINT_COLUMNS = {
    "AIRPORT_SEQ_ID",
    "AIRPORT_ID",
    "AIRPORT_IS_CLOSED",
    "AIRPORT_IS_LATEST",
}

NULLABLE_UINT_COLUMNS = {
    "AIRPORT_WAC_SEQ_ID2",
    "AIRPORT_WAC",
    "CITY_MARKET_SEQ_ID",
    "CITY_MARKET_ID",
    "CITY_MARKET_WAC_SEQ_ID2",
    "CITY_MARKET_WAC",
}

NULLABLE_INT_COLUMNS = {
    "LAT_DEGREES",
    "LAT_MINUTES",
    "LAT_SECONDS",
    "LON_DEGREES",
    "LON_MINUTES",
    "LON_SECONDS",
}

NULLABLE_FLOAT_COLUMNS = {
    "LATITUDE",
    "LONGITUDE",
}

DATE_COLUMNS = {
    "AIRPORT_START_DATE",
    "AIRPORT_THRU_DATE",
}


@dataclass(frozen=True)
class DownloadArtifact:
    zip_path: Path
    csv_name: str
    downloaded_at: str
    source_url: str = BASE_DOWNLOAD_URL


def log(message: str) -> None:
    if not DEBUG:
        return
    now = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    print(f"[{now}] {message}", file=sys.stderr, flush=True)


def ensure_dirs() -> None:
    DOWNLOAD_DIR.mkdir(parents=True, exist_ok=True)
    META_DIR.mkdir(parents=True, exist_ok=True)


def build_opener() -> urllib.request.OpenerDirector:
    cookie_jar = http.cookiejar.CookieJar()
    return urllib.request.build_opener(urllib.request.HTTPCookieProcessor(cookie_jar))


def fetch_text(opener: urllib.request.OpenerDirector, url: str) -> str:
    request = urllib.request.Request(url, headers={"User-Agent": "Mozilla/5.0"})
    with opener.open(request, timeout=60) as response:
        return response.read().decode("utf-8", errors="replace")


def extract_hidden_fields(html: str) -> dict[str, str]:
    fields = {
        name: value
        for name, value in re.findall(
            r'name="(__VIEWSTATE|__VIEWSTATEGENERATOR|__EVENTVALIDATION)"[^>]*value="([^"]*)"',
            html,
        )
    }
    missing = sorted({"__VIEWSTATE", "__VIEWSTATEGENERATOR", "__EVENTVALIDATION"} - fields.keys())
    if missing:
        raise RuntimeError(f"missing hidden form fields: {missing}")
    return fields


def extract_checkbox_names(html: str) -> list[str]:
    names = [name for _, name in re.findall(r'<input id="([A-Z0-9_]+)" type="checkbox" name="([A-Z0-9_]+)"', html)]
    if sorted(names) != sorted(EXPECTED_COLUMNS):
        missing = sorted(set(EXPECTED_COLUMNS) - set(names))
        extra = sorted(set(names) - set(EXPECTED_COLUMNS))
        raise RuntimeError(f"unexpected export checkbox set: missing={missing} extra={extra}")
    return names


def export_zip(force: bool) -> DownloadArtifact:
    ensure_dirs()
    opener = build_opener()
    html = fetch_text(opener, BASE_DOWNLOAD_URL)
    hidden_fields = extract_hidden_fields(html)
    checkbox_names = extract_checkbox_names(html)

    timestamp = datetime.now(UTC).strftime("%Y%m%d_%H%M%S")
    zip_path = DOWNLOAD_DIR / f"T_MASTER_CORD_{timestamp}.zip"
    meta_path = META_DIR / "latest_download.json"
    if meta_path.exists() and not force:
        metadata = json.loads(meta_path.read_text())
        cached_zip = Path(metadata["zip_path"])
        if cached_zip.exists() and zipfile.is_zipfile(cached_zip):
            return DownloadArtifact(
                zip_path=cached_zip,
                csv_name=metadata["csv_name"],
                downloaded_at=metadata["downloaded_at"],
                source_url=metadata.get("source_url", BASE_DOWNLOAD_URL),
            )

    form = dict(hidden_fields)
    for name in checkbox_names:
        form[name] = "on"
    form.update(
        {
            "btnDownload": "Download",
            "cboGeography": "All",
            "cboYear": "All",
            "cboPeriod": "All",
        }
    )
    encoded = urllib.parse.urlencode(form).encode()
    request = urllib.request.Request(
        BASE_DOWNLOAD_URL,
        data=encoded,
        method="POST",
        headers={
            "User-Agent": "Mozilla/5.0",
            "Referer": BASE_DOWNLOAD_URL,
            "Content-Type": "application/x-www-form-urlencoded",
        },
    )
    log("starting BTS export download")
    with opener.open(request, timeout=300) as response:
        content_type = response.headers.get("Content-Type", "")
        if "zip" not in content_type.lower():
            body = response.read(500).decode("utf-8", errors="replace")
            raise RuntimeError(f"expected zip export, got {content_type!r}: {body[:200]}")
        blob = response.read()

    zip_path.write_bytes(blob)
    if not zipfile.is_zipfile(zip_path):
        zip_path.unlink(missing_ok=True)
        raise RuntimeError("downloaded BTS export is not a valid zip file")
    with zipfile.ZipFile(zip_path) as archive:
        csv_names = [name for name in archive.namelist() if name.lower().endswith(".csv")]
        if len(csv_names) != 1:
            raise RuntimeError(f"expected exactly one CSV in export zip, found {csv_names}")
        csv_name = csv_names[0]

    artifact = DownloadArtifact(zip_path=zip_path, csv_name=csv_name, downloaded_at=datetime.now(UTC).isoformat())
    meta_path.write_text(
        json.dumps(
            {
                "zip_path": str(artifact.zip_path),
                "csv_name": artifact.csv_name,
                "downloaded_at": artifact.downloaded_at,
                "source_url": artifact.source_url,
                "table_info_url": BASE_INFO_URL,
            },
            indent=2,
            sort_keys=True,
        )
        + "\n"
    )
    return artifact


def inspect_export(artifact: DownloadArtifact) -> dict[str, object]:
    with zipfile.ZipFile(artifact.zip_path) as archive, archive.open(artifact.csv_name) as raw:
        text = io.TextIOWrapper(raw, encoding="utf-8", newline="")
        reader = csv.reader(text)
        header = next(reader)
        if header != EXPECTED_COLUMNS:
            raise RuntimeError(f"unexpected BTS header: {header}")
        row_count = sum(1 for _ in reader)

    info = {
        "zip_path": str(artifact.zip_path),
        "csv_name": artifact.csv_name,
        "source_url": artifact.source_url,
        "table_info_url": BASE_INFO_URL,
        "header": EXPECTED_COLUMNS,
        "row_count": row_count,
        "downloaded_at": artifact.downloaded_at,
        "inspected_at": datetime.now(UTC).isoformat(),
    }
    (META_DIR / "latest_inspection.json").write_text(json.dumps(info, indent=2, sort_keys=True) + "\n")
    return info


def run_clickhouse(connection: str, query: str | None = None, query_file: Path | None = None) -> subprocess.CompletedProcess[str]:
    args = ["clickhouse-client", "--connection", connection]
    if query_file is not None:
        args.extend(["--queries-file", str(query_file)])
    elif query is not None:
        args.extend(["--query", query])
    else:
        raise ValueError("query or query_file is required")
    return subprocess.run(args, text=True, capture_output=True, check=False)


def create_tables(connection: str) -> None:
    for sql_file in (SCRIPT_DIR / "airports_bts_schema.sql", SCRIPT_DIR / "airports_bts_latest_view.sql"):
        result = run_clickhouse(connection, query_file=sql_file)
        if result.returncode != 0:
            raise RuntimeError(f"failed to apply {sql_file.name}: {result.stderr.strip()}")


def to_uint(value: str, column: str) -> str:
    try:
        parsed = int(Decimal(value))
    except (InvalidOperation, ValueError) as exc:
        raise RuntimeError(f"cannot parse unsigned integer {value!r} for {column}") from exc
    if parsed < 0:
        raise RuntimeError(f"negative value {value!r} is invalid for {column}")
    return str(parsed)


def to_int(value: str, column: str) -> str:
    try:
        return str(int(Decimal(value)))
    except (InvalidOperation, ValueError) as exc:
        raise RuntimeError(f"cannot parse integer {value!r} for {column}") from exc


def to_float(value: str, column: str) -> str:
    try:
        return str(float(value))
    except ValueError as exc:
        raise RuntimeError(f"cannot parse float {value!r} for {column}") from exc


def to_date(value: str, column: str) -> str:
    try:
        return datetime.strptime(value, "%m/%d/%Y %I:%M:%S %p").strftime("%Y-%m-%d")
    except ValueError as exc:
        raise RuntimeError(f"cannot parse date {value!r} for {column}") from exc


def normalize_value(column: str, raw: str) -> str:
    value = raw.strip()
    if column in STRING_COLUMNS:
        return value
    if column in UINT_COLUMNS:
        if not value:
            raise RuntimeError(f"{column} is required")
        return to_uint(value, column)
    if column in NULLABLE_UINT_COLUMNS:
        return r"\N" if not value else to_uint(value, column)
    if column in NULLABLE_INT_COLUMNS:
        return r"\N" if not value else to_int(value, column)
    if column in NULLABLE_FLOAT_COLUMNS:
        return r"\N" if not value else to_float(value, column)
    if column in DATE_COLUMNS:
        return r"\N" if not value else to_date(value, column)
    raise RuntimeError(f"unclassified column {column}")


def load_export(connection: str, artifact: DownloadArtifact, inspection: dict[str, object]) -> dict[str, int]:
    truncate = run_clickhouse(connection, query="TRUNCATE TABLE default.airports_bts")
    if truncate.returncode != 0:
        raise RuntimeError(truncate.stderr.strip())

    insert_query = f"INSERT INTO default.airports_bts ({', '.join(EXPECTED_COLUMNS)}) FORMAT TabSeparated"
    proc = subprocess.Popen(
        ["clickhouse-client", "--connection", connection, "--query", insert_query],
        stdin=subprocess.PIPE,
        text=True,
    )
    assert proc.stdin is not None

    inserted_rows = 0
    with zipfile.ZipFile(artifact.zip_path) as archive, archive.open(artifact.csv_name) as raw:
        text = io.TextIOWrapper(raw, encoding="utf-8", newline="")
        reader = csv.DictReader(text)
        if reader.fieldnames != EXPECTED_COLUMNS:
            raise RuntimeError(f"unexpected BTS header during load: {reader.fieldnames}")
        for line_no, row in enumerate(reader, start=2):
            try:
                normalized = [normalize_value(column, row.get(column, "")) for column in EXPECTED_COLUMNS]
            except RuntimeError as exc:
                raise RuntimeError(f"{artifact.csv_name}: line {line_no}: {exc}") from exc
            proc.stdin.write("\t".join(normalized))
            proc.stdin.write("\n")
            inserted_rows += 1

    proc.stdin.close()
    return_code = proc.wait()
    if return_code != 0:
        raise RuntimeError(f"clickhouse insert failed with exit code {return_code}")

    counted = run_clickhouse(connection, query="SELECT count() FROM default.airports_bts")
    if counted.returncode != 0:
        raise RuntimeError(counted.stderr.strip())
    counted_rows = int(counted.stdout.strip() or "0")
    expected_rows = int(inspection["row_count"])
    if counted_rows != inserted_rows or counted_rows != expected_rows:
        raise RuntimeError(
            f"row count mismatch: inserted={inserted_rows}, counted={counted_rows}, expected={expected_rows}"
        )
    return {
        "inserted_rows": inserted_rows,
        "counted_rows": counted_rows,
    }


def query_tsv(connection: str, sql: str) -> list[list[str]]:
    result = subprocess.run(
        ["clickhouse-client", "--connection", connection, "--format", "TabSeparatedRaw", "--query", sql],
        text=True,
        capture_output=True,
        check=False,
    )
    if result.returncode != 0:
        raise RuntimeError(result.stderr.strip())
    lines = [line for line in result.stdout.splitlines() if line]
    return [line.split("\t") for line in lines]


def verify(connection: str) -> dict[str, object]:
    row_counts = query_tsv(
        connection,
        """
        SELECT
            (SELECT count() FROM default.airports_bts),
            (SELECT count() FROM default.airports_bts_latest),
            (SELECT count() FROM default.airports_bts_latest WHERE AIRPORT_IS_LATEST != 1)
        """.strip(),
    )[0]

    origin_coverage = query_tsv(
        connection,
        """
        SELECT
            countDistinctIf(o.OriginAirportID, o.OriginAirportID != 0),
            countDistinctIf(o.OriginAirportID, o.OriginAirportID != 0 AND a.AIRPORT_ID IS NOT NULL)
        FROM default.ontime_v2 AS o
        LEFT JOIN default.airports_bts_latest AS a
            ON o.OriginAirportID = a.AIRPORT_ID
        """.strip(),
    )[0]
    dest_coverage = query_tsv(
        connection,
        """
        SELECT
            countDistinctIf(o.DestAirportID, o.DestAirportID != 0),
            countDistinctIf(o.DestAirportID, o.DestAirportID != 0 AND a.AIRPORT_ID IS NOT NULL)
        FROM default.ontime_v2 AS o
        LEFT JOIN default.airports_bts_latest AS a
            ON o.DestAirportID = a.AIRPORT_ID
        """.strip(),
    )[0]

    unmatched_origins = query_tsv(
        connection,
        """
        SELECT
            toString(o.OriginAirportID),
            any(replaceAll(toString(o.Origin), '\\0', '')),
            toString(count())
        FROM default.ontime_v2 AS o
        LEFT JOIN default.airports_bts_latest AS a
            ON o.OriginAirportID = a.AIRPORT_ID
        WHERE o.OriginAirportID != 0
          AND a.AIRPORT_ID IS NULL
        GROUP BY o.OriginAirportID
        ORDER BY count() DESC, o.OriginAirportID
        LIMIT 20
        """.strip(),
    )

    unmatched_dests = query_tsv(
        connection,
        """
        SELECT
            toString(o.DestAirportID),
            any(replaceAll(toString(o.Dest), '\\0', '')),
            toString(count())
        FROM default.ontime_v2 AS o
        LEFT JOIN default.airports_bts_latest AS a
            ON o.DestAirportID = a.AIRPORT_ID
        WHERE o.DestAirportID != 0
          AND a.AIRPORT_ID IS NULL
        GROUP BY o.DestAirportID
        ORDER BY count() DESC, o.DestAirportID
        LIMIT 20
        """.strip(),
    )

    matched_sample = query_tsv(
        connection,
        """
        SELECT
            replaceAll(toString(o.Origin), '\\0', '') AS Origin,
            toString(any(o.OriginAirportID)),
            toString(any(a.AIRPORT_ID)),
            any(a.DISPLAY_AIRPORT_NAME),
            any(a.UTC_LOCAL_TIME_VARIATION),
            toString(any(a.LATITUDE)),
            toString(any(a.LONGITUDE))
        FROM default.ontime_v2 AS o
        INNER JOIN default.airports_bts_latest AS a
            ON o.OriginAirportID = a.AIRPORT_ID
        GROUP BY Origin
        ORDER BY Origin
        LIMIT 10
        """.strip(),
    )

    current_airports_id_misjoin = query_tsv(
        connection,
        """
        SELECT count()
        FROM default.ontime_v2 AS o
        INNER JOIN default.airports AS a
            ON toString(o.OriginAirportID) = a.AirportID
        """.strip(),
    )[0][0]

    bts_meta = {
        "airports_bts_rows": int(row_counts[0]),
        "airports_bts_latest_rows": int(row_counts[1]),
        "latest_view_non_latest_rows": int(row_counts[2]),
        "origin_distinct_airport_ids": int(origin_coverage[0]),
        "origin_matched_airport_ids": int(origin_coverage[1]),
        "dest_distinct_airport_ids": int(dest_coverage[0]),
        "dest_matched_airport_ids": int(dest_coverage[1]),
        "legacy_airports_id_join_rows": int(current_airports_id_misjoin),
        "unmatched_origins": unmatched_origins,
        "unmatched_dests": unmatched_dests,
        "matched_sample": matched_sample,
        "verified_at": datetime.now(UTC).isoformat(),
    }
    (META_DIR / "latest_verification.json").write_text(json.dumps(bts_meta, indent=2, sort_keys=True) + "\n")
    return bts_meta


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(description="Load the official BTS Master Coordinate airport dimension into ClickHouse.")
    parser.add_argument("--connection", default="demo", help="clickhouse-client connection name")
    parser.add_argument("--debug", action="store_true", help="print progress logs to stderr")
    subparsers = parser.add_subparsers(dest="command", required=True)

    create_parser = subparsers.add_parser("create-tables")
    create_parser.add_argument("--connection", dest="subcommand_connection")

    download_parser = subparsers.add_parser("download")
    download_parser.add_argument("--connection", dest="subcommand_connection")
    download_parser.add_argument("--force-download", action="store_true")

    load_parser = subparsers.add_parser("load")
    load_parser.add_argument("--connection", dest="subcommand_connection")
    load_parser.add_argument("--force-download", action="store_true")

    verify_parser = subparsers.add_parser("verify")
    verify_parser.add_argument("--connection", dest="subcommand_connection")
    return parser


def main() -> int:
    global DEBUG
    parser = build_parser()
    args = parser.parse_args()
    if getattr(args, "subcommand_connection", None):
        args.connection = args.subcommand_connection
    DEBUG = args.debug

    try:
        if args.command == "create-tables":
            create_tables(args.connection)
            return 0

        if args.command == "download":
            artifact = export_zip(force=args.force_download)
            inspection = inspect_export(artifact)
            print(json.dumps(inspection, indent=2, sort_keys=True))
            return 0

        if args.command == "load":
            create_tables(args.connection)
            artifact = export_zip(force=args.force_download)
            inspection = inspect_export(artifact)
            result = load_export(args.connection, artifact, inspection)
            print(json.dumps({**inspection, **result}, indent=2, sort_keys=True))
            return 0

        if args.command == "verify":
            print(json.dumps(verify(args.connection), indent=2, sort_keys=True))
            return 0
    except (RuntimeError, urllib.error.URLError) as exc:
        print(str(exc), file=sys.stderr)
        return 1

    parser.error(f"unsupported command {args.command}")
    return 2


if __name__ == "__main__":
    raise SystemExit(main())
