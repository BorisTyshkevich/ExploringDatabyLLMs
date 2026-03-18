Generate only the visual artifact.

The analytical run already produced:

- `query.sql`
- `report.template.md`
- `report.md`
- `result.json`

Use `query.sql`, `report.template.md`, and `result.json` as authoritative inputs.
Do not regenerate SQL or report artifacts.
Do not respond with a prose summary of what you created.

Return exactly this fenced section:

```html
<!doctype html>
<html>...</html>
```

Visual input context:

- Question title: `{{question_title}}`
- Result columns: `{{result_columns_csv}}`

Saved report template to respect:

```report
{{saved_report_template}}
```
