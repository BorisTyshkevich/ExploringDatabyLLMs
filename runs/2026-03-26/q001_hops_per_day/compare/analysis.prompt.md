You are writing a rich comparison note for one qforge benchmark question and one experiment day.

Use only verified local artifacts and, if necessary, the configured MCP access for direct validation. Do not invent behavior, metrics, SQL differences, or output differences.

Question context:

- Question ID: `q001`
- Question slug: `q001_hops_per_day`
- Question title: `Highest daily hops for one aircraft on one flight number`
- Day: `2026-03-26`

Primary structured compare artifact:

- Local path: `2026-03-26/q001_hops_per_day/compare/compare.json`
- Published URL: `https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-26/q001_hops_per_day/compare/compare.json`

Question files:

- Question prompt: `prompts/q001_hops_per_day/report_prompt.md` (`https://github.com/boristyshkevich/ExploringDatabyLLMs/blob/main/prompts/q001_hops_per_day/report_prompt.md`)
- Visual prompt: `prompts/q001_hops_per_day/visual_prompt.md` (`https://github.com/boristyshkevich/ExploringDatabyLLMs/blob/main/prompts/q001_hops_per_day/visual_prompt.md`)

Published run artifact links to use in the final Markdown:

### claude / opus
- `run-002`
  - Published links: main.sql: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-26/q001_hops_per_day/claude/opus/run-002/queries/main.sql | q1.sql: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-26/q001_hops_per_day/claude/opus/run-002/queries/q1.sql | q2.sql: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-26/q001_hops_per_day/claude/opus/run-002/queries/q2.sql | report.md: https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-26%2Fq001_hops_per_day%2Fclaude%2Fopus%2Frun-002%2Freport.md | review.md: https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-26%2Fq001_hops_per_day%2Fclaude%2Fopus%2Frun-002%2Freview.md | result.json: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-26/q001_hops_per_day/claude/opus/run-002/result.json | visual.html: https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/2026-03-26/q001_hops_per_day/claude/opus/run-002/visual.html
  - Local verification: `2026-03-26/q001_hops_per_day/claude/opus/run-002/queries/main.sql`, `2026-03-26/q001_hops_per_day/claude/opus/run-002/queries/q1.sql`, `2026-03-26/q001_hops_per_day/claude/opus/run-002/queries/q2.sql`, `2026-03-26/q001_hops_per_day/claude/opus/run-002/report.md`, `2026-03-26/q001_hops_per_day/claude/opus/run-002/review.md`, `2026-03-26/q001_hops_per_day/claude/opus/run-002/result.json`, `2026-03-26/q001_hops_per_day/claude/opus/run-002/visual.html`

### claude / sonnet
- `run-001`
  - Published links: main.sql: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-26/q001_hops_per_day/claude/sonnet/run-001/queries/main.sql | q1.sql: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-26/q001_hops_per_day/claude/sonnet/run-001/queries/q1.sql | q2.sql: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-26/q001_hops_per_day/claude/sonnet/run-001/queries/q2.sql | report.md: https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-26%2Fq001_hops_per_day%2Fclaude%2Fsonnet%2Frun-001%2Freport.md | review.md: https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-26%2Fq001_hops_per_day%2Fclaude%2Fsonnet%2Frun-001%2Freview.md | result.json: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-26/q001_hops_per_day/claude/sonnet/run-001/result.json | visual.html: https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/2026-03-26/q001_hops_per_day/claude/sonnet/run-001/visual.html
  - Local verification: `2026-03-26/q001_hops_per_day/claude/sonnet/run-001/queries/main.sql`, `2026-03-26/q001_hops_per_day/claude/sonnet/run-001/queries/q1.sql`, `2026-03-26/q001_hops_per_day/claude/sonnet/run-001/queries/q2.sql`, `2026-03-26/q001_hops_per_day/claude/sonnet/run-001/report.md`, `2026-03-26/q001_hops_per_day/claude/sonnet/run-001/review.md`, `2026-03-26/q001_hops_per_day/claude/sonnet/run-001/result.json`, `2026-03-26/q001_hops_per_day/claude/sonnet/run-001/visual.html`

### codex / gpt-5.4
- `run-001`
  - Published links: main.sql: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-001/queries/main.sql | q1.sql: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-001/queries/q1.sql | q2.sql: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-001/queries/q2.sql | report.md: https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-26%2Fq001_hops_per_day%2Fcodex%2Fgpt-5.4%2Frun-001%2Freport.md | review.md: https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-26%2Fq001_hops_per_day%2Fcodex%2Fgpt-5.4%2Frun-001%2Freview.md | result.json: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-001/result.json | visual.html: https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-001/visual.html
  - Local verification: `2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-001/queries/main.sql`, `2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-001/queries/q1.sql`, `2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-001/queries/q2.sql`, `2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-001/report.md`, `2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-001/review.md`, `2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-001/result.json`, `2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-001/visual.html`

Run directories:

- 2026-03-26/q001_hops_per_day/claude/opus/run-002
- 2026-03-26/q001_hops_per_day/claude/sonnet/run-001
- 2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-001

Query SQL files:

- 2026-03-26/q001_hops_per_day/claude/opus/run-002/queries/main.sql
- 2026-03-26/q001_hops_per_day/claude/opus/run-002/queries/q1.sql
- 2026-03-26/q001_hops_per_day/claude/opus/run-002/queries/q2.sql
- 2026-03-26/q001_hops_per_day/claude/sonnet/run-001/queries/main.sql
- 2026-03-26/q001_hops_per_day/claude/sonnet/run-001/queries/q1.sql
- 2026-03-26/q001_hops_per_day/claude/sonnet/run-001/queries/q2.sql
- 2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-001/queries/main.sql
- 2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-001/queries/q1.sql
- 2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-001/queries/q2.sql

Report Markdown files:

- 2026-03-26/q001_hops_per_day/claude/opus/run-002/report.md
- 2026-03-26/q001_hops_per_day/claude/sonnet/run-001/report.md
- 2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-001/report.md

Review Markdown files:

- 2026-03-26/q001_hops_per_day/claude/opus/run-002/review.md
- 2026-03-26/q001_hops_per_day/claude/sonnet/run-001/review.md
- 2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-001/review.md

Visual HTML files:

- 2026-03-26/q001_hops_per_day/claude/opus/run-002/visual.html
- 2026-03-26/q001_hops_per_day/claude/sonnet/run-001/visual.html
- 2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-001/visual.html

Result JSON files:

- 2026-03-26/q001_hops_per_day/claude/opus/run-002/result.json
- 2026-03-26/q001_hops_per_day/claude/sonnet/run-001/result.json
- 2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-001/result.json

Deterministic compare summary:

# qforge Compare Report

- Generated: `2026-03-26T13:58:02Z`
- Day: `2026-03-26`
- Question: `q001`

## q001: Highest daily hops for one aircraft on one flight number

- Status: 1 run(s) did not finish cleanly: claude/sonnet/run-001.
- Result rows (manifest): all runs returned 65 rows.
- Fastest successful run: claude/opus/run-002 at 122.44 s.
- Lowest read volume: claude/opus/run-002 at 3,075,342,183 rows.
- Lowest memory usage: claude/opus/run-002 at 49.1 GiB.
- Multi-query note: manifest result rows are the sum of all subquery result rows; synthetic result.json rows are summary rows only.
| runner | model | run | target | review | review md | status | result rows (manifest) | sql gen | visual gen | build | query time | read rows | bytes read | peak memory | warnings |
| --- | --- | --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| claude | opus | run-002 | html | PASS | [review.md](https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-26%2Fq001_hops_per_day%2Fclaude%2Fopus%2Frun-002%2Freview.md) | ok | 65 | 654.34 s | 241.31 s | n/a | 122.44 s | 3,075,342,183 | 38.0 GiB | 49.1 GiB | 0 |
| claude | sonnet | run-001 | html | WARN | [review.md](https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-26%2Fq001_hops_per_day%2Fclaude%2Fsonnet%2Frun-001%2Freview.md) | partial | 65 | 676.02 s | 679.17 s | n/a | 214.39 s | 3,720,509,274 | 43.9 GiB | 59.2 GiB | 0 |
| codex | gpt-5.4 | run-001 | html | PASS | [review.md](https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-26%2Fq001_hops_per_day%2Fcodex%2Fgpt-5.4%2Frun-001%2Freview.md) | ok | 65 | n/a | n/a | n/a | 306.33 s | 3,386,975,101 | 52.8 GiB | 57.0 GiB | 0 |

Your job:

- write one evidence-based Markdown report suitable for `compare_report.md`
- use the real local artifacts above as the source of truth
- verify whether outputs actually differ before claiming they differ
- quantify differences when they exist
- mention performance differences only from verified query-log metrics
- use `review.md` when discussing correctness, prompt-alignment issues, or known failures
- describe SQL-shape differences only when supported by the actual `query.sql` files
- cite `report.md` and `visual.html` artifacts when discussing presentation outputs
- prefer links to local artifacts instead of long pasted SQL
- use only the published URLs provided above in the final Markdown for run artifacts; never emit absolute filesystem paths
- for `report.md`, use the provided `md.html?file=...` URL
- for `query.sql` and `result.json`, use the provided GitHub blob URL
- for `visual.html`, use the provided GitHub Pages file URL
- in sections 6 and 9, group content by provider/model and then by run id
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
- In sections 6 and 9, prefer short provider-grouped subsections with per-run bullets.
- In section 10 (`## Execution stats`), use a Markdown table as the primary presentation, not prose bullets.
- The execution-stats table should have one row per run and include at least: provider/model, run id, status, query time, rows read, bytes read, peak memory, SQL generation time, visual generation time, and total run duration when available.
- After the execution-stats table, add at most one short paragraph calling out the most important performance spread or anomaly.
- Return only one fenced Markdown block.

Return exactly this fenced section:

```markdown
# qNNN Experiment Note
...
```