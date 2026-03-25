# Feature Request: `sql_plus_transform` Analysis Mode

## Summary

Add an optional analysis mode that keeps one authoritative SQL extraction query, then allows a local transform step to derive secondary outputs for `report.md` and `visual.html` without requiring one independent proof SQL query per dashboard question.

Working name:

- `analysis_mode: sql_plus_transform`

This is a proposal only. It does not change current qforge behavior.

## Motivation

The current `multi_query_json` mode is strong for proof-oriented workflows because each dashboard question has:

- one direct prose answer
- one proof SQL statement

That gives good provenance, but it also has a limitation:

- each subquestion SQL is executed independently by qforge

As implemented today, qforge does not take the result rows from the first proof query and feed them into later subquestions. It executes every proof query separately, saves each result, and renders `report.md` from the saved prose answers plus lightweight query summaries.

This makes some prompt guidance awkward, especially when several dashboard questions are really lightweight derivations of the same ranked candidate set.

Examples:

- “Which itinerary is the highest-hop example?”
- “Which of the top-ranked itineraries is the most recent?”
- “Do the top itineraries look recurring or one-off?”

In some questions, one primary SQL query could return a compact result set, and later answers could be computed locally from that result instead of requiring separate database queries.

## Proposed Mode

`sql_plus_transform` would separate analysis into two layers:

1. authoritative extraction in SQL
2. deterministic local derivation in a transform step

Proposed required artifacts:

- `query.sql`
- `transform.js`
- `answer.raw.json`

Optional alternative:

- `query.sql`
- `transform.js`
- `report.template.md`

The first variant is closer to current `multi_query_json`. The second variant is closer to current template-based modes.

## High-Level Flow

1. The provider writes `query.sql`.
2. qforge executes `query.sql` and writes canonical `result.json`.
3. The provider writes `transform.js`.
4. qforge runs `transform.js` locally against `result.json`.
5. The transform emits a deterministic derived artifact such as `derived.json`.
6. qforge uses `derived.json` to render `report.md`.
7. The visual phase may use `result.json` plus `derived.json`.

## Why This Might Help

### Better fit for reuse-heavy questions

Some questions naturally have:

- one main extraction query
- several lightweight derived answers
- dashboard state that should be computed from the same small result set

In these cases, local transform logic may be simpler and more efficient than forcing every derived question back into separate SQL.

### Cleaner visual handoff

The visual phase often needs:

- ranked rows
- derived labels
- grouped summaries
- comparison metrics

If those are computed once in a deterministic transform step, the visual can consume stable derived structures instead of reconstructing them again.

### Less pressure on prompt wording

Today, some prompt rules imply reuse across proof queries, but the harness executes those queries independently. A transform phase would make true reuse explicit and supported.

## Main Tradeoffs

### Pros

- Allows true reuse of one main query result without rescanning source tables for lightweight derived questions.
- Better fit for dashboard-oriented questions where several answers come from the same compact result.
- Makes some report and visual shaping logic easier to express than in SQL alone.
- Can reduce repeated query cost when the main result set is already sufficient.

### Cons

- Weakens the simplicity of the current “one subquestion, one proof SQL” provenance model.
- Splits analytical logic across SQL and JavaScript, which increases review complexity.
- Adds runtime and sandbox complexity because qforge would need to execute local code safely.
- Creates new failure modes in transform execution, artifact validation, and schema compatibility.
- Makes review harder because correctness would depend on both extracted data and transform logic.
- Encourages models to generate executable code instead of only SQL and prose, which may reduce reliability.

## Provenance Risk

The strongest argument against this mode is analytical auditability.

Current `multi_query_json` gives a direct mapping:

- question
- proof SQL
- executed result
- prose answer

With a transform step, the mapping becomes:

- extraction SQL
- raw result
- transform code
- derived result
- prose answer

That is more flexible, but also less transparent.

If this feature is pursued, provenance should remain explicit.

Suggested mitigation:

- treat `query.sql` as the only source-of-truth extraction artifact
- require `transform.js` to be deterministic and side-effect free
- save `derived.json`
- include both extraction and transform metadata in `manifest.json`
- ensure review prompts reference both files

## Suggested Contract

One possible contract:

### Provider outputs

- `query.sql`
- `transform.js`
- `answer.raw.json`

### `answer.raw.json` shape

```json
{
  "subquestions": [
    {
      "subquestion": "Question text copied from the contract",
      "answer_markdown": "Direct prose answer",
      "derived_key": "name of the derived object that proves this answer"
    }
  ]
}
```

### Transform input

qforge provides `result.json` to `transform.js`.

### Transform output

`transform.js` writes `derived.json`, for example:

```json
{
  "tables": {
    "top_itineraries": [],
    "route_repetition": []
  },
  "metrics": {
    "max_hops": 8,
    "most_recent_top_date": "2021-09-26"
  },
  "proofs": {
    "q1": {},
    "q2": {},
    "q3": {}
  }
}
```

This would let each prose answer point to a specific derived proof object without requiring every proof to be an independent SQL query.

## Alternative: Keep SQL-Only, But Add Better Semantics

Before adding a new mode, qforge could instead tighten the wording around `multi_query_json`:

- each proof query is executed independently
- later proof queries cannot consume earlier executed JSON results
- shared candidate-set logic may still be repeated inside SQL where useful

This is much cheaper than adding a new mode and may solve most prompt confusion.

## Recommended Scope If Implemented Later

Keep the first version narrow:

- Node.js only
- one input file: `result.json`
- one output file: `derived.json`
- no network access
- no filesystem writes except the declared output
- no package installation
- deterministic runtime only

This would limit complexity and keep the feature reviewable.

## Open Questions

- Should `transform.js` be allowed only in a new analysis mode, or optionally within `template_files`?
- Should the report renderer consume `derived.json` directly, or should the provider still own prose answers in `answer.raw.json`?
- How should review prompts verify transform correctness?
- Should visuals read `derived.json` directly, or only through a normalized visual-input artifact?
- Is the provenance cost acceptable for benchmark-style questions?

## Recommendation

Do not replace `multi_query_json` with this mode.

If pursued, add `sql_plus_transform` as an additional mode for questions where:

- one primary query is authoritative
- later answers are lightweight deterministic derivations
- visual shaping benefits from local post-processing

For the current codebase, the lowest-cost next step is not implementation. It is documenting the distinction between:

- independent proof-query execution in `multi_query_json`
- hypothetical future local derivation from a primary query result

