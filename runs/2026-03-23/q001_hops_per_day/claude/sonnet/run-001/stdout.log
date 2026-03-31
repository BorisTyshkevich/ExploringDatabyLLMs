Both files written:

- **`query.sql`** — single statement using a CTE to group legs by `(Tail_Number, FlightNum, Carrier, FlightDate)`, sorts legs chronologically by `CRSDepTime`, then builds the `Route` string via `arrayConcat([first_origin], all_dests)` joined with ` → `. Orders by `hop_count DESC, FlightDate DESC LIMIT 10`.

- **`report.template.md`** — template with `{{data_overview_md}}` and `{{result_table_md}}` placeholders plus a fixed analytical takeaway about high hop counts indicating aggressive aircraft utilization on sequential point-to-point routes.
