Both `query.sql` and `visual.html` have been successfully generated and saved to the workspace.

The `visual.html` file implements the requested dynamic-mode dashboard and includes:
- JWE storage and tokenized execution against the Altinity MCP endpoint.
- Execution of the primary multi-hop itinerary analysis query.
- A secondary, explicit coordinate enrichment query against `ontime.dim_airports` that caches results per session.
- A robust query ledger that tracks status and full SQL for both the primary execution and the geocoding enrichment.
- A functional Leaflet `html_map` implementation highlighting the selected route sequence, with graceful degradation messaging if coordinates are missing.
- Interactive row selection triggering map and route-sequence panel redraws while keeping top-level KPIs anchored.
