You are writing a rich comparison note for one qforge benchmark question and one experiment day.

Use only verified local artifacts and, if necessary, the configured MCP access for direct validation. Do not invent behavior, metrics, SQL differences, or output differences.

Question context:

- Question ID: `{{question_id}}`
- Question slug: `{{question_slug}}`
- Question title: `{{question_title}}`
- Day: `{{compare_day}}`

Primary structured compare artifact:

- Local path: `{{compare_json_path}}`
- Published URL: `{{compare_json_url}}`

Question files:

- Question prompt: `{{question_prompt_path}}` (`{{question_prompt_url}}`)
- Visual prompt: `{{visual_prompt_path}}` (`{{visual_prompt_url}}`)
- Compare contract: `{{compare_contract_path}}` (`{{compare_contract_url}}`)

Published run artifact links to use in the final Markdown:

{{published_run_artifacts_md}}

Run directories:

{{run_dirs_md}}

Query SQL files:

{{query_sql_paths_md}}

Report Markdown files:

{{report_md_paths_md}}

Visual HTML files:

{{visual_html_paths_md}}

Result JSON files:

{{result_json_paths_md}}

Deterministic compare summary:

{{compare_summary_md}}

Your job:

- write one evidence-based Markdown report suitable for `compare_report.md`
- use the real local artifacts above as the source of truth
- verify whether outputs actually differ before claiming they differ
- quantify differences when they exist
- mention performance differences only from verified query-log metrics
- describe SQL-shape differences only when supported by the actual `query.sql` files
- cite `report.md` and `visual.html` artifacts when discussing presentation outputs
- prefer links to local artifacts instead of long pasted SQL
- use only the published URLs provided above in the final Markdown for run artifacts; never emit absolute filesystem paths
- for `report.md`, use the provided `md.html?file=...` URL
- for `query.sql` and `result.json`, use the provided GitHub blob URL
- for `visual.html`, use the provided GitHub Pages file URL
- in sections 6, 9, and 10, group content by provider/model and then by run id
- keep the note concise but complete enough for a blog-style benchmark write-up

Required sections:

1. `# qNNN Experiment Note`
2. `## Question`
3. `## Why this question is useful`
4. `## Experiment setup`
5. `## Result summary`
6. `## Full SQL artifacts`
7. `## Real output differences`
8. `## SQL comparison`
9. `## Presentation artifacts`
10. `## Execution stats`
11. `## Takeaway`

Rules:

- If results are identical, say so explicitly and support that with verified evidence.
- If differences are localized to one field or row type, say that precisely.
- Do not use vague judgments like “better” or “worse” without concrete evidence.
- Do not mention files that you did not verify.
- In sections 6, 9, and 10, prefer short provider-grouped subsections with per-run bullets.
- Return only one fenced Markdown block.

Return exactly this fenced section:

```markdown
# qNNN Experiment Note
...
```
