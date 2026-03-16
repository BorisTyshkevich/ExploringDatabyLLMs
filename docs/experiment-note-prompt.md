# Experiment Note Prompt

Use this prompt when you need a short, evidence-based report for one benchmark question and one experiment day.

## Goal

Write a Markdown artifact under `docs/` named like:

- `q002-experiment-note.md`
- `q004-experiment-note.md`

The note should be suitable for a blog article or internal write-up.

## Inputs

You must work from real local artifacts only:

- question prompt files under `questions/<question-slug>/`
- generated SQL files under `runs/<day>/<question-slug>/<runner>/<model>/run-XXX/query.sql`
- generated `result.json`
- compare artifacts such as `runs/q002-compare.md` and `runs/q002-compare.json`
- direct validation queries against the dataset when needed

Do not invent behavior or assume outputs are identical unless you verified that directly.

## Required sections

Write the note with these sections:

1. `# qNNN Experiment Note`
2. `## Question`
3. `## Why this question is useful`
4. `## Experiment setup`
5. `## Result summary`
6. `## Full SQL artifacts`
7. `## Real output differences`
8. `## SQL comparison`
9. `## Execution stats`
10. `## Takeaway`

Add extra sections only when they are justified by evidence, for example:

- `## Data interpretation`
- `## What a direct data check shows`
- `## Prompt correction made during the experiment`

## Hard rules

- Use only verified facts from local files or build and run a query over the dataset using mcp server.
- If you claim results differ, quantify the difference.
- If you claim results are identical, verify that directly.
- Do not describe model output quality in vague terms like “better” or “worse” without concrete evidence.
- Do not hallucinate execution metrics.
- Do not hallucinate SQL differences.
- Prefer links to full local SQL artifacts rather than pasting long SQL into the note.
- If a difference is localized to one field, say so explicitly.
- If a prompt ambiguity was fixed during the experiment, describe the ambiguity and the fix.

## Required comparison work

Before writing the note:

1. Locate all relevant `query.sql` files.
2. Locate all relevant `result.json` files.
3. Compare the result payloads directly.
4. Identify whether differences are:
   - row-count differences
   - ranking differences
   - field-value differences
   - formatting/nullability differences
5. If helpful, compute normalized row hashes for each result.
6. Read the compare report to extract:
   - status
   - row count
   - duration
   - read rows
   - memory usage
7. Run direct validation SQL when the prompt or outputs suggest ambiguity.

## Preferred reporting style

- Keep the note short enough to be blog-usable.
- Use tables for:
  - provider execution stats
  - pairwise difference summaries
  - spot-check metric comparisons
- Use bullet lists for:
  - strengths
  - concrete problems
  - prompt fixes
- Focus on what was actually learned from the experiment, not just artifact inventory.

## Suggested analysis checklist

- Is the question itself satisfiable on the dataset?
- Do all providers return the same number of rows?
- Are the top-ranked entities the same?
- Do any models violate a subtle business rule?
- Are numeric differences only due to rounding or due to different logic?
- Does one query scan materially more data than the others?
- Is there a prompt ambiguity that should be fixed?

## Example framing for real differences

Good:

- `Claude vs Codex differ on 157 of 195 rows, mainly because Claude repeats leader-transition fields on non-leader rows while Codex nulls them out.`
- `All three q004 outputs match exactly except for P90DepDelayMinutes.`
- `Codex and Gemini scanned roughly twice as many rows as Claude on q004.`

Bad:

- `Codex understood the task best.`
- `Gemini was weaker overall.`
- `The outputs were basically the same.`

## Output requirement

Save the final note into `docs/` and, if asked, copy it to a second repo without changing the content.
