# Proposal: Better Guidance For Additional Questions In `visual_prompt.md`

## Summary

Keep additional question headings in `visual_prompt.md` as prose-only author guidance.

This proposal does not change current qforge behavior. It recommends clearer documentation and stronger authoring conventions so question prompts can suggest extra dashboard panels or visual-only SQL without implying that qforge parses or verifies those sections as first-class query contracts.

## Current Model

Today, qforge treats `visual_prompt.md` differently from `report_prompt.md`.

- `report_prompt.md` may define structured analysis sections that feed the saved analysis contract
- `visual_prompt.md` is passed through as presentation guidance
- extra headings inside `visual_prompt.md` are visible to the model, but qforge does not normalize them into saved subquestions

That means a visual prompt may ask for extra panels such as:

- connector analysis
- geographic comparison
- operational stress summaries

The authored dashboard may decide to create visual-only lookup or enrichment SQL for those panels, but those queries are:

- second-pass browser queries
- not saved as reviewed proof queries
- not written to `queries/*.sql`
- not surfaced as `Question.Subquestions`

## Why Keep It Prose-Only

The prose-only model is preferable for now because it matches actual runtime behavior and keeps the harness boundary clear.

Benefits:

- avoids pretending that question-local panel semantics are part of the qforge contract
- lets prompt authors describe intended dashboard analysis without overfitting tests to one implementation
- keeps shared runtime behavior in shared prompt assets instead of scattering it across question prompts
- avoids adding parser, schema, or metadata complexity for a workflow that is still exploratory

## Failure Mode This Proposal Addresses

The recent q001 cleanup exposed a boundary mismatch:

- the visual prompt was rewritten to ask broader additional questions
- tests were still asserting exact lookup-query mechanics for one specific panel
- those tests effectively treated prose author guidance as if qforge owned and verified it

That creates brittle tests and discourages useful prompt evolution.

The better boundary is:

- test shared presentation-contract guarantees
- let question-local additional questions remain prose unless qforge explicitly grows a new structured feature

## Recommended Authoring Convention

Document and encourage a clearer prose pattern for additional visual-only questions.

Recommended style:

- use explicit headings such as `### additional question: key connectors` or another clearly presentation-oriented heading style
- describe the intended panel or comparison outcome
- describe any important user-facing degraded behavior only when it is essential to the question
- avoid over-specifying exact query lifecycle mechanics unless the dashboard truly depends on them
- rely on shared prompt assets for common rules such as query ledger behavior, browser-side verification, date reruns, and general lookup-query handling

Good examples:

- “Ask an additional question about which airports act as key connectors across the top itineraries.”
- “Add a panel comparing the geographic spread of the top itineraries.”
- “If this panel depends on extra browser-side SQL, keep the main analysis visible if the panel fails.”

Less useful examples:

- repeating the full shared lookup-query contract in every question prompt
- making tests depend on a specific ledger label or rerun trigger for one question-local panel
- implying that the extra question becomes a reviewed proof query

## Explicit Non-Goals

This proposal does not recommend:

- parsing additional visual headings into structured metadata
- adding a new `visual_subquestions` schema
- saving visual-only panel queries as analysis artifacts
- extending `Question.Subquestions` to include `visual_prompt.md`

If qforge ever needs machine-readable visual panel contracts, that should be a separate future proposal with a dedicated design.

## Recommendation

Keep `visual_prompt.md` additional questions prose-only.

Improve docs and test boundaries so prompt authors can ask for richer panels without implying unsupported harness semantics.
