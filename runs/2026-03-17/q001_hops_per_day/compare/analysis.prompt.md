You are writing a rich comparison note for one qforge benchmark question and one experiment day.

Use only verified local artifacts and, if necessary, the configured MCP access for direct validation. Do not invent behavior, metrics, SQL differences, or output differences.

Question context:

- Question ID: `q001`
- Question slug: `q001_hops_per_day`
- Question title: `Highest daily hops for one aircraft on one flight number`
- Day: `2026-03-17`

Primary structured compare artifact:

- Local path: `runs/2026-03-17/q001_hops_per_day/compare/compare.json`
- Published URL: `https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-17/q001_hops_per_day/compare/compare.json`

Question files:

- SQL question prompt: `prompts/q001_hops_per_day/prompt.md` (`https://github.com/boristyshkevich/ExploringDatabyLLMs/blob/main/prompts/q001_hops_per_day/prompt.md`)
- Report prompt: `prompts/q001_hops_per_day/report_prompt.md` (`https://github.com/boristyshkevich/ExploringDatabyLLMs/blob/main/prompts/q001_hops_per_day/report_prompt.md`)
- Visual prompt: `prompts/q001_hops_per_day/visual_prompt.md` (`https://github.com/boristyshkevich/ExploringDatabyLLMs/blob/main/prompts/q001_hops_per_day/visual_prompt.md`)
- Compare contract: `prompts/q001_hops_per_day/compare.yaml` (`https://github.com/boristyshkevich/ExploringDatabyLLMs/blob/main/prompts/q001_hops_per_day/compare.yaml`)

Published run artifact links to use in the final Markdown:

### claude / opus
- `run-003`
  - Published links: query.sql: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-17/q001_hops_per_day/claude/opus/run-003/query.sql | report.md: https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-17%2Fq001_hops_per_day%2Fclaude%2Fopus%2Frun-003%2Freport.md | result.json: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-17/q001_hops_per_day/claude/opus/run-003/result.json | visual.html: https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/2026-03-17/q001_hops_per_day/claude/opus/run-003/visual.html
  - Local verification: `runs/2026-03-17/q001_hops_per_day/claude/opus/run-003/query.sql`, `runs/2026-03-17/q001_hops_per_day/claude/opus/run-003/report.md`, `runs/2026-03-17/q001_hops_per_day/claude/opus/run-003/result.json`, `runs/2026-03-17/q001_hops_per_day/claude/opus/run-003/visual.html`

### codex / gpt-5.4
- `run-003`
  - Published links: query.sql: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-17/q001_hops_per_day/codex/gpt-5.4/run-003/query.sql | report.md: https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-17%2Fq001_hops_per_day%2Fcodex%2Fgpt-5.4%2Frun-003%2Freport.md | result.json: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-17/q001_hops_per_day/codex/gpt-5.4/run-003/result.json | visual.html: https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/2026-03-17/q001_hops_per_day/codex/gpt-5.4/run-003/visual.html
  - Local verification: `runs/2026-03-17/q001_hops_per_day/codex/gpt-5.4/run-003/query.sql`, `runs/2026-03-17/q001_hops_per_day/codex/gpt-5.4/run-003/report.md`, `runs/2026-03-17/q001_hops_per_day/codex/gpt-5.4/run-003/result.json`, `runs/2026-03-17/q001_hops_per_day/codex/gpt-5.4/run-003/visual.html`

### gemini / gemini-3.1-pro-preview
- `run-003`
  - Published links: query.sql: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-17/q001_hops_per_day/gemini/gemini-3.1-pro-preview/run-003/query.sql | report.md: https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-17%2Fq001_hops_per_day%2Fgemini%2Fgemini-3.1-pro-preview%2Frun-003%2Freport.md | result.json: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-17/q001_hops_per_day/gemini/gemini-3.1-pro-preview/run-003/result.json | visual.html: https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/2026-03-17/q001_hops_per_day/gemini/gemini-3.1-pro-preview/run-003/visual.html
  - Local verification: `runs/2026-03-17/q001_hops_per_day/gemini/gemini-3.1-pro-preview/run-003/query.sql`, `runs/2026-03-17/q001_hops_per_day/gemini/gemini-3.1-pro-preview/run-003/report.md`, `runs/2026-03-17/q001_hops_per_day/gemini/gemini-3.1-pro-preview/run-003/result.json`, `runs/2026-03-17/q001_hops_per_day/gemini/gemini-3.1-pro-preview/run-003/visual.html`

Run directories:

- runs/2026-03-17/q001_hops_per_day/claude/opus/run-003
- runs/2026-03-17/q001_hops_per_day/codex/gpt-5.4/run-003
- runs/2026-03-17/q001_hops_per_day/gemini/gemini-3.1-pro-preview/run-003

Query SQL files:

- runs/2026-03-17/q001_hops_per_day/claude/opus/run-003/query.sql
- runs/2026-03-17/q001_hops_per_day/codex/gpt-5.4/run-003/query.sql
- runs/2026-03-17/q001_hops_per_day/gemini/gemini-3.1-pro-preview/run-003/query.sql

Report Markdown files:

- runs/2026-03-17/q001_hops_per_day/claude/opus/run-003/report.md
- runs/2026-03-17/q001_hops_per_day/codex/gpt-5.4/run-003/report.md
- runs/2026-03-17/q001_hops_per_day/gemini/gemini-3.1-pro-preview/run-003/report.md

Visual HTML files:

- runs/2026-03-17/q001_hops_per_day/claude/opus/run-003/visual.html
- runs/2026-03-17/q001_hops_per_day/codex/gpt-5.4/run-003/visual.html
- runs/2026-03-17/q001_hops_per_day/gemini/gemini-3.1-pro-preview/run-003/visual.html

Result JSON files:

- runs/2026-03-17/q001_hops_per_day/claude/opus/run-003/result.json
- runs/2026-03-17/q001_hops_per_day/codex/gpt-5.4/run-003/result.json
- runs/2026-03-17/q001_hops_per_day/gemini/gemini-3.1-pro-preview/run-003/result.json

Deterministic compare summary:

# qforge Compare Report

- Generated: `2026-03-17T20:20:04Z`
- Day: `2026-03-17`
- Question: `q001`

## q001: Highest daily hops for one aircraft on one flight number

- Status: 1 run(s) did not finish cleanly: gemini/gemini-3.1-pro-preview/run-003.
- Row counts: all runs returned 10 rows.
- Fastest successful run: codex/gpt-5.4/run-003 at 7.88 s.
- Lowest read volume: claude/opus/run-003 at 193,061,941 rows.
- Lowest memory usage: claude/opus/run-003 at 35.1 GiB.
| runner | model | run | status | rows | duration | read rows | memory | warnings |
| --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: |
| claude | opus | run-003 | ok | 10 | 18.63 s | 193,061,941 | 35.1 GiB | 0 |
| codex | gpt-5.4 | run-003 | ok | 10 | 7.88 s | 193,061,941 | 45.3 GiB | 0 |
| gemini | gemini-3.1-pro-preview | run-003 | partial | 10 | 15.53 s | 230,307,587 | 48.9 GiB | 0 |

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