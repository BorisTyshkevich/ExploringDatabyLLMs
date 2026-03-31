# React Visual Artifacts Vs Plain HTML/JS

This note captures a practical conclusion for qforge-style visual artifact generation: React is often a better target for model generation and iterative editing, while plain HTML/CSS/JS is usually a better target for the final browser artifact when the page is mostly static.

## Short Answer

- For model generation, React is often easier to produce and revise.
- For browser execution, plain HTML/CSS/JS is usually lighter and faster for static artifacts.
- A good compromise is to generate in a component-oriented form, then prerender or export to mostly static output when interaction needs are limited.

## Why React Can Be Easier For Model Generation

Several recent UI code-generation papers point toward hierarchical, declarative, component-level generation as a strong fit for LLMs:

- [`UICoder`](https://arxiv.org/abs/2406.07739), June 2024
  - Focuses on user-interface code generation and shows that targeted UI-oriented training improves generation quality.
- [`Bridging Design and Development with Automated Declarative UI Code Generation`](https://arxiv.org/abs/2409.11667), September 2024
  - Supports the broader idea that declarative UI representations are a good generation target.
- [`Design2Code`](https://aclanthology.org/2025.naacl-long.199/), NAACL 2025
  - Benchmarks multimodal front-end generation and shows that layout fidelity and element recall remain hard, which increases the value of structured intermediate representations.
- [`FrontendBench`](https://arxiv.org/abs/2506.13832), June 2025
  - Provides a scalable benchmark for front-end generation and evaluation, reinforcing that front-end work benefits from clearer structural targets.
- [`ScreenCoder`](https://arxiv.org/abs/2507.22827), July 2025
  - Uses modular planning with hierarchical layout reasoning before producing HTML/CSS.

Practical inference:

- React gives the model natural edit units such as components, props, and local state.
- Iterative fixes are often easier when a page is already decomposed into components.
- Large monolithic HTML/JS outputs tend to be longer, flatter, and harder to patch reliably.

This is an inference from the papers above, not a direct benchmark result saying "React beats HTML/JS for LLM generation" in all settings.

## Why Plain HTML/JS Can Be Better For Browser Execution

For delivered browser performance, the tradeoff usually goes the other direction:

- React-style client apps often ship more JavaScript.
- They may pay reconciliation and hydration costs.
- Static HTML/CSS avoids most framework runtime overhead for non-interactive content.

Relevant source:

- [`Million.js: A Fast Compiler-Augmented Virtual DOM for the Web`](https://arxiv.org/abs/2202.08409), arXiv version January 1, 2023
  - Reports substantial improvements over popular virtual-DOM libraries and includes a real-world case where the migrated app loaded faster after moving away from React.

Practical inference:

- If the artifact is mostly a report, dashboard snapshot, or static narrative visual, plain HTML/CSS with minimal JavaScript will usually execute faster and load less code.
- If the page needs complex client-side interaction, React can still be the better engineering tradeoff despite runtime overhead.

## Recommendation For This Repo

For this repository's visual artifacts:

- Prefer a component-oriented generation mindset because it is easier for models to structure and revise.
- Prefer static final output when the artifact does not require substantial client-side interaction.
- Avoid paying the full React runtime cost for pages that are effectively presentational.

In other words:

- Authoring target: React-like component structure is often the easier generation target.
- Delivery target: static HTML/CSS/JS is often the better browser artifact.

## Early Repository Result

An initial paired benchmark in this repository did not support the hypothesis that React would be faster for model authoring.

Case:

- question: `q001_hops_per_day`
- runner/model: `claude/sonnet`
- date: March 24, 2026
- comparison:
  - HTML run: [`run-002`](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-24/q001_hops_per_day/claude/sonnet/run-002/manifest.json)
  - React run: [`run-003`](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-24/q001_hops_per_day/claude/sonnet/run-003/manifest.json)

Observed timings from `manifest.json`:

- HTML `presentation_provider_duration_ms`: `478164`
- React `presentation_provider_duration_ms`: `502791`
- Difference: React slower by `24627 ms`, about `5%`

Build step:

- HTML `presentation_build_duration_ms`: `0`
- React `presentation_build_duration_ms`: `9140`

End-to-end presentation total:

- HTML: `478164 ms`
- React: `511931 ms`
- Difference: React slower by `33767 ms`, about `7%`

Interpretation:

- For this dynamic map-heavy dashboard, React was not faster for model generation.
- The extra build step made the end-to-end presentation path slower still.
- This is one benchmark case only, not yet a general conclusion for all dashboard types.

Recommended next comparisons:

- static non-map dashboard
- dynamic non-map dashboard
- dynamic map dashboard

If React remains slower across all three categories, the original hypothesis becomes much weaker.

## Important Caveat

I did not find a clean, paper-backed apples-to-apples comparison that directly measures:

- LLM generation speed for React vs plain HTML/JS
- LLM edit reliability for React vs plain HTML/JS
- final browser runtime performance for both outputs on the same benchmarked tasks

So the conclusion here is evidence-based but still partly inferential:

- generation favors structure and components
- delivery favors less runtime and less JavaScript

This repository's first direct comparison adds an important practical caveat:

- better structure for the model does not automatically mean faster generation in practice
- for at least one complex dashboard, React generation was slower than direct HTML generation
