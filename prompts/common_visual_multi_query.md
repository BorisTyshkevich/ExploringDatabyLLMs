Create browser-ready HTML `visual.html`.

Return exactly one fenced `html` block containing the full document. If your environment can also write files directly, save the same content to `visual.html`. Do not open the artifact view frame.

### Rules

- Question title: `{{question_title}}`
- Visual mode: `{{visual_mode}}`
- Visual type: `{{visual_type}}`
- `visual.html` is the canonical human-facing artifact. Use the verified subquestion answers and proof-query previews below to build both the prose and the charts/tables.
- Treat the verified analysis package below as the only analytical input for this visual pass.
- Do not query ClickHouse, use MCP tools, or run new SQL during this visual pass.
- Do not assume access to `query.sql`, `result.json`, browser-side live queries, or full result tables unless they are explicitly included in that package.
- Build a self-contained page with inline CSS and JavaScript.
- Do not load remote scripts or stylesheets.
- Embed the verified analysis package, or the exact subset you use from it, directly in the HTML in a `<script type="application/json">` data block.
- Parse that embedded JSON in browser JavaScript and derive all displayed values from it.
- Do not fabricate additional rows, synthetic time series, approximate heatmap cells, or extrapolated values that are not present in the embedded data.
- Derive narrative claims, KPIs, chart values, and highlights from the verified subquestion answers and query previews. Do not invent or hardcode them.
- Follow question-specific visual guidance after the shared contract. Put reusable runtime behavior in shared page code, not in prose comments.

{{visual_prompt_md}}

### Verified Analysis Package

Use this JSON package as the authoritative input for the visual:

{{visual_input_summary_json}}
