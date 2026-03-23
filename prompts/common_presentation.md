Generate only the visual artifact.

Use file `query.sql` as the authoritative input for the primary data query.
Use file `visual_input.json` to understand the result shape before building visuals without querying dataset.
You may construct additional queries when needed for enrichment or drill-down, but do not rebuild the saved SQL.

Return exactly this fenced section:

```html
<!doctype html>
<html>...</html>
```

Visual input requirements:   {{visual_input_summary_json}}

