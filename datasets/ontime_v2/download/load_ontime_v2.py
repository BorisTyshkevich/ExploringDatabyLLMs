#!/usr/bin/env python3

from __future__ import annotations

import argparse
import csv
import io
import json
import re
import subprocess
import sys
import urllib.error
import urllib.request
import zipfile
from dataclasses import dataclass
from datetime import UTC, datetime
from decimal import Decimal, InvalidOperation
from pathlib import Path


BASE_URL = "https://transtats.bts.gov/PREZIP"
PREFIX = "On_Time_Reporting_Carrier_On_Time_Performance_1987_present"
SCRIPT_DIR = Path(__file__).resolve().parent
CACHE_DIR = SCRIPT_DIR / ".cache"
DOWNLOAD_DIR = CACHE_DIR / "downloads"
META_DIR = CACHE_DIR / "meta"
DEBUG = False


STORED_COLUMNS = [
    "FlightDate",
    "Reporting_Airline",
    "DOT_ID_Reporting_Airline",
    "IATA_CODE_Reporting_Airline",
    "Tail_Number",
    "Flight_Number_Reporting_Airline",
    "OriginAirportID",
    "OriginAirportSeqID",
    "OriginCityMarketID",
    "Origin",
    "OriginCityName",
    "OriginState",
    "OriginStateFips",
    "OriginStateName",
    "OriginWac",
    "DestAirportID",
    "DestAirportSeqID",
    "DestCityMarketID",
    "Dest",
    "DestCityName",
    "DestState",
    "DestStateFips",
    "DestStateName",
    "DestWac",
    "CRSDepTime",
    "DepTime",
    "DepDelay",
    "DepDelayMinutes",
    "DepDel15",
    "DepartureDelayGroups",
    "DepTimeBlk",
    "TaxiOut",
    "WheelsOff",
    "WheelsOn",
    "TaxiIn",
    "CRSArrTime",
    "ArrTime",
    "ArrDelay",
    "ArrDelayMinutes",
    "ArrDel15",
    "ArrivalDelayGroups",
    "ArrTimeBlk",
    "Cancelled",
    "CancellationCode",
    "Diverted",
    "CRSElapsedTime",
    "ActualElapsedTime",
    "AirTime",
    "Flights",
    "Distance",
    "DistanceGroup",
    "CarrierDelay",
    "WeatherDelay",
    "NASDelay",
    "SecurityDelay",
    "LateAircraftDelay",
    "FirstDepTime",
    "TotalAddGTime",
    "LongestAddGTime",
    "DivAirportLandings",
    "DivReachedDest",
    "DivActualElapsedTime",
    "DivArrDelay",
    "DivDistance",
    "Div1Airport",
    "Div1AirportID",
    "Div1AirportSeqID",
    "Div1WheelsOn",
    "Div1TotalGTime",
    "Div1LongestGTime",
    "Div1WheelsOff",
    "Div1TailNum",
    "Div2Airport",
    "Div2AirportID",
    "Div2AirportSeqID",
    "Div2WheelsOn",
    "Div2TotalGTime",
    "Div2LongestGTime",
    "Div2WheelsOff",
    "Div2TailNum",
    "Div3Airport",
    "Div3AirportID",
    "Div3AirportSeqID",
    "Div3WheelsOn",
    "Div3TotalGTime",
    "Div3LongestGTime",
    "Div3WheelsOff",
    "Div3TailNum",
    "Div4Airport",
    "Div4AirportID",
    "Div4AirportSeqID",
    "Div4WheelsOn",
    "Div4TotalGTime",
    "Div4LongestGTime",
    "Div4WheelsOff",
    "Div4TailNum",
    "Div5Airport",
    "Div5AirportID",
    "Div5AirportSeqID",
    "Div5WheelsOn",
    "Div5TotalGTime",
    "Div5LongestGTime",
    "Div5WheelsOff",
    "Div5TailNum",
]


LEGACY_NAME_MAP = {
    "UniqueCarrier": "Reporting_Airline",
    "AirlineID": "DOT_ID_Reporting_Airline",
    "Carrier": "IATA_CODE_Reporting_Airline",
    "TailNum": "Tail_Number",
    "FlightNum": "Flight_Number_Reporting_Airline",
}

IGNORED_SOURCE_COLUMNS = {"", "Year", "Quarter", "Month", "DayofMonth", "DayOfWeek"}

STRING_COLUMNS = {
    "Reporting_Airline",
    "IATA_CODE_Reporting_Airline",
    "Tail_Number",
    "Flight_Number_Reporting_Airline",
    "Origin",
    "OriginCityName",
    "OriginState",
    "OriginStateFips",
    "OriginStateName",
    "Dest",
    "DestCityName",
    "DestState",
    "DestStateFips",
    "DestStateName",
    "DepTimeBlk",
    "ArrTimeBlk",
    "CancellationCode",
    "Div1Airport",
    "Div1TailNum",
    "Div2Airport",
    "Div2TailNum",
    "Div3Airport",
    "Div3TailNum",
    "Div4Airport",
    "Div4TailNum",
    "Div5Airport",
    "Div5TailNum",
}

NON_NULL_INT_COLUMNS = {
    "DOT_ID_Reporting_Airline",
    "OriginAirportID",
    "OriginAirportSeqID",
    "OriginCityMarketID",
    "OriginWac",
    "DestAirportID",
    "DestAirportSeqID",
    "DestCityMarketID",
    "DestWac",
    "DepDel15",
    "ArrDel15",
    "Cancelled",
    "Diverted",
    "Flights",
    "DistanceGroup",
    "Div1AirportID",
    "Div1AirportSeqID",
    "Div2AirportID",
    "Div2AirportSeqID",
    "Div3AirportID",
    "Div3AirportSeqID",
    "Div4AirportID",
    "Div4AirportSeqID",
    "Div5AirportID",
    "Div5AirportSeqID",
    "DivReachedDest",
}

NULLABLE_UINT_COLUMNS = {
    "CRSDepTime",
    "DepTime",
    "DepDelayMinutes",
    "TaxiOut",
    "WheelsOff",
    "WheelsOn",
    "TaxiIn",
    "CRSArrTime",
    "ArrTime",
    "ArrDelayMinutes",
    "Distance",
    "CarrierDelay",
    "WeatherDelay",
    "NASDelay",
    "SecurityDelay",
    "LateAircraftDelay",
    "FirstDepTime",
    "DivAirportLandings",
    "DivActualElapsedTime",
    "DivDistance",
    "Div1WheelsOn",
    "Div1WheelsOff",
    "Div2WheelsOn",
    "Div2WheelsOff",
    "Div3WheelsOn",
    "Div3WheelsOff",
    "Div4WheelsOn",
    "Div4WheelsOff",
    "Div5WheelsOn",
    "Div5WheelsOff",
}

NULLABLE_INT_COLUMNS = {
    "DepDelay",
    "DepartureDelayGroups",
    "ArrDelay",
    "ArrivalDelayGroups",
    "CRSElapsedTime",
    "ActualElapsedTime",
    "AirTime",
    "TotalAddGTime",
    "LongestAddGTime",
    "DivActualElapsedTime",
    "DivArrDelay",
    "Div1TotalGTime",
    "Div1LongestGTime",
    "Div2TotalGTime",
    "Div2LongestGTime",
    "Div3TotalGTime",
    "Div3LongestGTime",
    "Div4TotalGTime",
    "Div4LongestGTime",
    "Div5TotalGTime",
    "Div5LongestGTime",
}

HHMM_COLUMNS = {
    "CRSDepTime",
    "DepTime",
    "WheelsOff",
    "WheelsOn",
    "CRSArrTime",
    "ArrTime",
    "FirstDepTime",
    "Div1WheelsOn",
    "Div1WheelsOff",
    "Div2WheelsOn",
    "Div2WheelsOff",
    "Div3WheelsOn",
    "Div3WheelsOff",
    "Div4WheelsOn",
    "Div4WheelsOff",
    "Div5WheelsOn",
    "Div5WheelsOff",
}


@dataclass(frozen=True)
class MonthRef:
    year: int
    month: int

    @property
    def stem(self) -> str:
        return f"{self.year}_{self.month:02d}"

    @property
    def zip_name(self) -> str:
        return f"{PREFIX}_{self.year}_{self.month}.zip"

    @property
    def url(self) -> str:
        return f"{BASE_URL}/{self.zip_name}"


def log(message: str) -> None:
    if not DEBUG:
        return
    now = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    print(f"[{now}] {message}", file=sys.stderr, flush=True)


def ensure_dirs() -> None:
    DOWNLOAD_DIR.mkdir(parents=True, exist_ok=True)
    META_DIR.mkdir(parents=True, exist_ok=True)


def run_clickhouse(connection: str, query: str | None = None, query_file: Path | None = None) -> subprocess.CompletedProcess[str]:
    args = ["clickhouse-client", "--connection", connection]
    if query_file is not None:
        args.extend(["--queries-file", str(query_file)])
    elif query is not None:
        args.extend(["--query", query])
    else:
        raise ValueError("query or query_file is required")
    target = str(query_file) if query_file is not None else (query[:160] + ("..." if query and len(query) > 160 else ""))
    log(f"clickhouse start: {target}")
    proc = subprocess.run(args, text=True, capture_output=True, check=False)
    log(f"clickhouse done rc={proc.returncode}: {target}")
    return proc


def list_available_months() -> list[MonthRef]:
    log(f"fetching PREZIP index: {BASE_URL}/")
    with urllib.request.urlopen(f"{BASE_URL}/", timeout=60) as response:
        html = response.read().decode("utf-8", errors="replace")
    matches = set(re.findall(rf"{PREFIX}_(\d{{4}})_(\d{{1,2}})\.zip", html))
    months = [MonthRef(int(year), int(month)) for year, month in matches]
    return sorted(months, key=lambda item: (item.year, item.month))


def selected_months(args: argparse.Namespace) -> list[MonthRef]:
    available = list_available_months()
    if args.command == "list-available":
        return available

    if args.command == "load-year":
        months = [item for item in available if item.year == args.year]
    else:
        start_year = args.start_year
        end_year = args.end_year
        months = [item for item in available if start_year <= item.year <= end_year]

    if getattr(args, "month", None) is not None:
        months = [item for item in months if item.month == args.month]
    return months


def download_month(month_ref: MonthRef, force: bool) -> Path:
    ensure_dirs()
    path = DOWNLOAD_DIR / month_ref.zip_name
    if path.exists() and not force:
        if zipfile.is_zipfile(path):
            log(f"download skip exists: {path.name}")
            return path
        log(f"download purge invalid cache: {path.name}")
        path.unlink()
    log(f"download start: {month_ref.url}")
    tmp_path = path.with_suffix(path.suffix + ".part")
    if tmp_path.exists():
        tmp_path.unlink()
    with urllib.request.urlopen(month_ref.url, timeout=300) as response, tmp_path.open("wb") as handle:
        while True:
            chunk = response.read(1024 * 1024)
            if not chunk:
                break
            handle.write(chunk)
    if not zipfile.is_zipfile(tmp_path):
        tmp_path.unlink(missing_ok=True)
        raise RuntimeError(f"downloaded file is not a valid zip: {month_ref.url}")
    tmp_path.replace(path)
    log(f"download done: {path.name} size={path.stat().st_size}")
    return path


def month_meta_path(month_ref: MonthRef) -> Path:
    return META_DIR / f"{month_ref.stem}.json"


def canonical_name(name: str) -> str:
    return LEGACY_NAME_MAP.get(name, name)


def detect_csv_encoding(zf: zipfile.ZipFile, csv_name: str) -> str:
    sample = zf.read(csv_name)[:200000]
    for encoding in ("utf-8", "cp1252", "latin-1"):
        try:
            sample.decode(encoding)
            return encoding
        except UnicodeDecodeError:
            continue
    raise RuntimeError(f"could not decode {csv_name} with supported encodings")


def inspect_archive(path: Path, strict: bool) -> dict[str, object]:
    with zipfile.ZipFile(path) as zf:
        csv_names = [name for name in zf.namelist() if name.lower().endswith(".csv")]
        if len(csv_names) != 1:
            raise RuntimeError(f"{path.name}: expected exactly one CSV, found {csv_names}")
        csv_name = csv_names[0]
        readme_names = [name for name in zf.namelist() if name.lower().endswith("readme.html")]
        encoding = detect_csv_encoding(zf, csv_name)
        with zf.open(csv_name) as raw:
            reader = csv.reader(io.TextIOWrapper(raw, encoding=encoding, newline=""))
            header = next(reader)

    trimmed_header = [canonical_name(name) for name in header if name not in IGNORED_SOURCE_COLUMNS]
    duplicates = sorted({name for name in trimmed_header if trimmed_header.count(name) > 1})
    if duplicates:
        raise RuntimeError(f"{path.name}: duplicate canonical columns {duplicates}")

    known = set(STORED_COLUMNS)
    extras = sorted(name for name in trimmed_header if name not in known)
    missing = sorted(name for name in STORED_COLUMNS if name not in trimmed_header)
    if extras and strict:
        raise RuntimeError(f"{path.name}: unexpected source columns {extras}")

    info = {
        "csv_name": csv_name,
        "encoding": encoding,
        "readme_names": readme_names,
        "header": header,
        "canonical_header": trimmed_header,
        "extras": extras,
        "missing": missing,
        "size_bytes": path.stat().st_size,
        "inspected_at": datetime.now(UTC).isoformat(),
    }
    log(
        "inspect done: "
        f"{path.name} csv={csv_name} encoding={encoding} cols={len(trimmed_header)} "
        f"missing={len(missing)} extras={len(extras)}"
    )
    return info


def to_int(value: str) -> int:
    try:
        return int(Decimal(value))
    except (InvalidOperation, ValueError) as exc:
        raise RuntimeError(f"cannot parse numeric value {value!r}") from exc


def to_uint(value: str, column: str) -> int:
    parsed = to_int(value)
    if parsed < 0:
        raise RuntimeError(f"negative value {value!r} is invalid for unsigned column {column}")
    return parsed


def normalize_value(column: str, raw: str) -> str:
    value = raw.strip()
    if column == "FlightDate":
        if not value:
            raise RuntimeError("FlightDate is required")
        return value
    if column in STRING_COLUMNS:
        return value
    if column in NON_NULL_INT_COLUMNS:
        return "0" if not value else str(to_uint(value, column))
    if column in NULLABLE_UINT_COLUMNS:
        if column in HHMM_COLUMNS and value and not value.isdigit():
            return r"\N"
        return r"\N" if not value else str(to_uint(value, column))
    if column in NULLABLE_INT_COLUMNS:
        return r"\N" if not value else str(to_int(value))
    raise RuntimeError(f"unclassified column {column}")


def insert_month(
    connection: str,
    table: str,
    zip_path: Path,
    inspected: dict[str, object],
    max_bad_rows: int,
) -> tuple[int, int]:
    csv_name = str(inspected["csv_name"])
    encoding = str(inspected["encoding"])
    header = list(inspected["header"])
    insert_query = (
        f"INSERT INTO {table} ({', '.join(STORED_COLUMNS)}) FORMAT TabSeparated"
    )
    log(f"insert start: {zip_path.name} -> {table}")
    proc = subprocess.Popen(
        ["clickhouse-client", "--connection", connection, "--query", insert_query],
        stdin=subprocess.PIPE,
        text=True,
    )
    assert proc.stdin is not None

    rows = 0
    skipped_rows = 0
    with zipfile.ZipFile(zip_path) as zf, zf.open(csv_name) as raw:
        reader = csv.reader(io.TextIOWrapper(raw, encoding=encoding, newline=""))
        next(reader)
        canonical_header = [canonical_name(name) for name in header]
        active_pairs = [(idx, name) for idx, name in enumerate(canonical_header) if name not in IGNORED_SOURCE_COLUMNS]
        for line_no, row in enumerate(reader, start=2):
            if len(row) < len(header):
                row.extend([""] * (len(header) - len(row)))
            if len(row) > len(header):
                row = row[: len(header)]
            source = {name: row[idx] for idx, name in active_pairs}
            try:
                normalized = [normalize_value(column, source.get(column, "")) for column in STORED_COLUMNS]
            except RuntimeError as exc:
                skipped_rows += 1
                if skipped_rows > max_bad_rows:
                    raise RuntimeError(
                        f"{zip_path.name}: exceeded malformed-row limit ({max_bad_rows}) at line {line_no}: {exc}"
                    ) from exc
                log(f"skip malformed row: file={zip_path.name} line={line_no} skipped={skipped_rows} error={exc}")
                continue
            proc.stdin.write("\t".join(normalized))
            proc.stdin.write("\n")
            rows += 1
            if DEBUG and rows % 500000 == 0:
                log(f"insert progress: {zip_path.name} rows={rows}")

    proc.stdin.close()
    returncode = proc.wait()
    if returncode != 0:
        raise RuntimeError(f"clickhouse insert failed for {zip_path.name} with exit code {returncode}")
    log(f"insert done: {zip_path.name} rows={rows} skipped={skipped_rows}")
    return rows, skipped_rows


def save_meta(month_ref: MonthRef, payload: dict[str, object]) -> None:
    month_meta_path(month_ref).write_text(json.dumps(payload, indent=2, sort_keys=True) + "\n")
    log(f"meta saved: {month_ref.stem}")


def create_tables(connection: str) -> None:
    for sql_file in (SCRIPT_DIR / "schema.sql", SCRIPT_DIR / "stage_schema.sql"):
        log(f"apply ddl: {sql_file.name}")
        result = run_clickhouse(connection, query_file=sql_file)
        if result.returncode != 0:
            raise RuntimeError(f"failed to apply {sql_file.name}: {result.stderr.strip()}")


def load_year(
    connection: str,
    year: int,
    months: list[MonthRef],
    force_download: bool,
    strict: bool,
    max_bad_rows_per_month: int,
) -> None:
    target_months = [item for item in months if item.year == year]
    if not target_months:
        raise RuntimeError(f"no source months available for {year}")

    log(f"load year start: {year} months={len(target_months)}")
    truncate = run_clickhouse(connection, query="TRUNCATE TABLE default.ontime_v2_stage")
    if truncate.returncode != 0:
        raise RuntimeError(truncate.stderr.strip())

    year_rows = 0
    for month_ref in target_months:
        log(f"month start: {month_ref.stem}")
        archive = download_month(month_ref, force_download)
        inspected = inspect_archive(archive, strict)
        inserted_rows, skipped_rows = insert_month(
            connection,
            "default.ontime_v2_stage",
            archive,
            inspected,
            max_bad_rows_per_month,
        )
        save_meta(
            month_ref,
            {
                **inspected,
                "inserted_rows": inserted_rows,
                "skipped_rows": skipped_rows,
                "year": month_ref.year,
                "month": month_ref.month,
                "source_url": month_ref.url,
            },
        )
        year_rows += inserted_rows
        log(f"month done: {month_ref.stem} rows={inserted_rows} skipped={skipped_rows}")

    count_query = f"SELECT count() FROM default.ontime_v2_stage WHERE Year = {year}"
    counted = run_clickhouse(connection, query=count_query)
    if counted.returncode != 0:
        raise RuntimeError(counted.stderr.strip())
    stage_rows = int(counted.stdout.strip() or "0")
    if stage_rows != year_rows:
        raise RuntimeError(f"stage row count mismatch for {year}: inserted={year_rows}, counted={stage_rows}")
    log(f"stage validated: year={year} rows={stage_rows}")

    replace = run_clickhouse(
        connection,
        query=f"ALTER TABLE default.ontime_v2 REPLACE PARTITION {year} FROM default.ontime_v2_stage",
    )
    if replace.returncode != 0:
        raise RuntimeError(replace.stderr.strip())
    log(f"publish done: year={year}")

    cleanup = run_clickhouse(connection, query="TRUNCATE TABLE default.ontime_v2_stage")
    if cleanup.returncode != 0:
        raise RuntimeError(cleanup.stderr.strip())
    log(f"stage cleanup done: year={year}")


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        description=(
            "Load BTS On-Time Reporting Carrier data into ontime_v2. "
            "Typical year controls: load-year --year YYYY, "
            "or backfill --start-year YYYY --end-year YYYY."
        )
    )
    parser.add_argument("--connection", default="demo", help="clickhouse-client connection name")
    parser.add_argument("--debug", action="store_true", help="print timestamped progress logs to stderr")
    parser.add_argument("--start-year", type=int, help="top-level shorthand for backfill start year")
    parser.add_argument("--end-year", type=int, help="top-level shorthand for backfill end year")
    parser.add_argument(
        "--max-bad-rows-per-month",
        type=int,
        default=100,
        help="skip up to this many malformed source rows per month before failing",
    )
    subparsers = parser.add_subparsers(dest="command", required=True)

    create_parser = subparsers.add_parser("create-tables")
    create_parser.add_argument("--connection", dest="subcommand_connection")

    list_parser = subparsers.add_parser("list-available")
    list_parser.add_argument("--connection", dest="subcommand_connection")

    year_parser = subparsers.add_parser("load-year")
    year_parser.add_argument("--connection", dest="subcommand_connection")
    year_parser.add_argument("--year", type=int, required=True)
    year_parser.add_argument("--month", type=int)
    year_parser.add_argument("--force-download", action="store_true")
    year_parser.add_argument("--allow-new-columns", action="store_true")

    range_parser = subparsers.add_parser("backfill")
    range_parser.add_argument("--connection", dest="subcommand_connection")
    range_parser.add_argument("--start-year", type=int, required=True)
    range_parser.add_argument("--end-year", type=int, required=True)
    range_parser.add_argument("--force-download", action="store_true")
    range_parser.add_argument("--allow-new-columns", action="store_true")

    return parser


def main() -> int:
    global DEBUG
    argv = sys.argv[1:]
    known_commands = {"create-tables", "list-available", "load-year", "backfill"}
    if not any(arg in known_commands for arg in argv):
        if "--start-year" in argv or "--end-year" in argv:
            insert_at = len(argv)
            for flag in ("--start-year", "--end-year"):
                if flag in argv:
                    insert_at = min(insert_at, argv.index(flag))
            argv = [*argv[:insert_at], "backfill", *argv[insert_at:]]

    parser = build_parser()
    args = parser.parse_args(argv)
    if getattr(args, "subcommand_connection", None):
        args.connection = args.subcommand_connection
    DEBUG = args.debug

    try:
        if args.command == "create-tables":
            create_tables(args.connection)
            return 0

        months = selected_months(args)
        if args.command == "list-available":
            for month_ref in months:
                print(month_ref.stem)
            return 0

        strict = not args.allow_new_columns
        if args.command == "load-year":
            create_tables(args.connection)
            load_year(
                args.connection,
                args.year,
                months,
                args.force_download,
                strict,
                args.max_bad_rows_per_month,
            )
            return 0

        if args.command == "backfill":
            create_tables(args.connection)
            for year in sorted({item.year for item in months}):
                load_year(
                    args.connection,
                    year,
                    months,
                    args.force_download,
                    strict,
                    args.max_bad_rows_per_month,
                )
            return 0
    except (RuntimeError, urllib.error.URLError) as exc:
        print(str(exc), file=sys.stderr)
        return 1

    parser.error(f"unsupported command {args.command}")
    return 2


if __name__ == "__main__":
    raise SystemExit(main())
