You are writing a rich comparison note for one qforge benchmark question and one experiment day.

Use only verified local artifacts and, if necessary, the configured MCP access for direct validation. Do not invent behavior, metrics, SQL differences, or output differences.

Question context:

- Question ID: `q002`
- Question slug: `q002_top_carrier_by_flights_leadership`
- Question title: `Yearly carrier leadership by completed flights`
- Day: `2026-03-17`

Primary structured compare artifact:

- Local path: `runs/2026-03-17/q002_top_carrier_by_flights_leadership/compare/compare.json`
- Published URL: `https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-17/q002_top_carrier_by_flights_leadership/compare/compare.json`

Question files:

- SQL question prompt: `prompts/q002_top_carrier_by_flights_leadership/prompt.md` (`https://github.com/boristyshkevich/ExploringDatabyLLMs/blob/main/prompts/q002_top_carrier_by_flights_leadership/prompt.md`)
- Report prompt: `prompts/q002_top_carrier_by_flights_leadership/report_prompt.md` (`https://github.com/boristyshkevich/ExploringDatabyLLMs/blob/main/prompts/q002_top_carrier_by_flights_leadership/report_prompt.md`)
- Visual prompt: `prompts/q002_top_carrier_by_flights_leadership/visual_prompt.md` (`https://github.com/boristyshkevich/ExploringDatabyLLMs/blob/main/prompts/q002_top_carrier_by_flights_leadership/visual_prompt.md`)
- Compare contract: `prompts/q002_top_carrier_by_flights_leadership/compare.yaml` (`https://github.com/boristyshkevich/ExploringDatabyLLMs/blob/main/prompts/q002_top_carrier_by_flights_leadership/compare.yaml`)

Published run artifact links to use in the final Markdown:

### claude / opus
- `run-001`
  - Published links: query.sql: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-17/q002_top_carrier_by_flights_leadership/claude/opus/run-001/query.sql | report.md: https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-17%2Fq002_top_carrier_by_flights_leadership%2Fclaude%2Fopus%2Frun-001%2Freport.md | result.json: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-17/q002_top_carrier_by_flights_leadership/claude/opus/run-001/result.json | visual.html: https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/2026-03-17/q002_top_carrier_by_flights_leadership/claude/opus/run-001/visual.html
  - Local verification: `runs/2026-03-17/q002_top_carrier_by_flights_leadership/claude/opus/run-001/query.sql`, `runs/2026-03-17/q002_top_carrier_by_flights_leadership/claude/opus/run-001/report.md`, `runs/2026-03-17/q002_top_carrier_by_flights_leadership/claude/opus/run-001/result.json`, `runs/2026-03-17/q002_top_carrier_by_flights_leadership/claude/opus/run-001/visual.html`

### codex / gpt-5.4
- `run-001`
  - Published links: query.sql: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-17/q002_top_carrier_by_flights_leadership/codex/gpt-5.4/run-001/query.sql | report.md: https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-17%2Fq002_top_carrier_by_flights_leadership%2Fcodex%2Fgpt-5.4%2Frun-001%2Freport.md | result.json: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-17/q002_top_carrier_by_flights_leadership/codex/gpt-5.4/run-001/result.json | visual.html: https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/2026-03-17/q002_top_carrier_by_flights_leadership/codex/gpt-5.4/run-001/visual.html
  - Local verification: `runs/2026-03-17/q002_top_carrier_by_flights_leadership/codex/gpt-5.4/run-001/query.sql`, `runs/2026-03-17/q002_top_carrier_by_flights_leadership/codex/gpt-5.4/run-001/report.md`, `runs/2026-03-17/q002_top_carrier_by_flights_leadership/codex/gpt-5.4/run-001/result.json`, `runs/2026-03-17/q002_top_carrier_by_flights_leadership/codex/gpt-5.4/run-001/visual.html`

### gemini / gemini-3.1-pro-preview
- `run-001`
  - Published links: query.sql: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-17/q002_top_carrier_by_flights_leadership/gemini/gemini-3.1-pro-preview/run-001/query.sql | report.md: https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-17%2Fq002_top_carrier_by_flights_leadership%2Fgemini%2Fgemini-3.1-pro-preview%2Frun-001%2Freport.md | result.json: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-17/q002_top_carrier_by_flights_leadership/gemini/gemini-3.1-pro-preview/run-001/result.json | visual.html: https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/2026-03-17/q002_top_carrier_by_flights_leadership/gemini/gemini-3.1-pro-preview/run-001/visual.html
  - Local verification: `runs/2026-03-17/q002_top_carrier_by_flights_leadership/gemini/gemini-3.1-pro-preview/run-001/query.sql`, `runs/2026-03-17/q002_top_carrier_by_flights_leadership/gemini/gemini-3.1-pro-preview/run-001/report.md`, `runs/2026-03-17/q002_top_carrier_by_flights_leadership/gemini/gemini-3.1-pro-preview/run-001/result.json`, `runs/2026-03-17/q002_top_carrier_by_flights_leadership/gemini/gemini-3.1-pro-preview/run-001/visual.html`

Run directories:

- runs/2026-03-17/q002_top_carrier_by_flights_leadership/claude/opus/run-001
- runs/2026-03-17/q002_top_carrier_by_flights_leadership/codex/gpt-5.4/run-001
- runs/2026-03-17/q002_top_carrier_by_flights_leadership/gemini/gemini-3.1-pro-preview/run-001

Query SQL files:

- runs/2026-03-17/q002_top_carrier_by_flights_leadership/claude/opus/run-001/query.sql
- runs/2026-03-17/q002_top_carrier_by_flights_leadership/codex/gpt-5.4/run-001/query.sql
- runs/2026-03-17/q002_top_carrier_by_flights_leadership/gemini/gemini-3.1-pro-preview/run-001/query.sql

Report Markdown files:

- runs/2026-03-17/q002_top_carrier_by_flights_leadership/claude/opus/run-001/report.md
- runs/2026-03-17/q002_top_carrier_by_flights_leadership/codex/gpt-5.4/run-001/report.md
- runs/2026-03-17/q002_top_carrier_by_flights_leadership/gemini/gemini-3.1-pro-preview/run-001/report.md

Visual HTML files:

- runs/2026-03-17/q002_top_carrier_by_flights_leadership/claude/opus/run-001/visual.html
- runs/2026-03-17/q002_top_carrier_by_flights_leadership/codex/gpt-5.4/run-001/visual.html
- runs/2026-03-17/q002_top_carrier_by_flights_leadership/gemini/gemini-3.1-pro-preview/run-001/visual.html

Result JSON files:

- runs/2026-03-17/q002_top_carrier_by_flights_leadership/claude/opus/run-001/result.json
- runs/2026-03-17/q002_top_carrier_by_flights_leadership/codex/gpt-5.4/run-001/result.json
- runs/2026-03-17/q002_top_carrier_by_flights_leadership/gemini/gemini-3.1-pro-preview/run-001/result.json

Deterministic compare summary:

# qforge Compare Report

- Generated: `2026-03-18T09:18:41Z`
- Day: `2026-03-17`
- Question: `q002`

## q002: Yearly carrier leadership by completed flights

- Status: all runs succeeded.
- Row counts: mismatch (195, 234).
- Fastest successful run: codex/gpt-5.4/run-001 at 586 ms.
- Lowest read volume: claude/opus/run-001 at 921,230,348 rows.
- Lowest memory usage: claude/opus/run-001 at 346.2 MiB.
| runner | model | run | status | rows | duration | read rows | memory | warnings |
| --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: |
| claude | opus | run-001 | ok | 195 | 594 ms | 921,230,348 | 346.2 MiB | 0 |
| codex | gpt-5.4 | run-001 | ok | 195 | 586 ms | 921,230,348 | 349.7 MiB | 0 |
| gemini | gemini-3.1-pro-preview | run-001 | ok | 234 | 1.46 s | 1,612,153,109 | 584.9 MiB | 0 |

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