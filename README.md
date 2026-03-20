# ExploringDataByLLMs

This repository contains two main active areas:

- [`qforge`](docs/qforge.md)
  - the Go-based benchmark harness for LLM-generated analytics workflows
  - covers prompt assembly, SQL execution, report and visual generation, compare runs, and CLI usage
- [`datasets/ontime/download`](datasets/ontime/download/README.md)
  - the dataset rebuild and loading workflow for `ontime.ontime`
  - covers source archives, table layout, loader commands, and schema rules

Related repo areas:

- [`prompts`](prompts)
  - shared and question-specific prompt assets used by `qforge`
- [`Skills`](Skills)
  - currently includes only the [`ontime-analyst-dashboard`](Skills/ontime-analyst-dashboard/SKILL.md) skill for `visual.html` generation
- [`cmd/qforge`](cmd/qforge)
  - CLI entrypoint for the active harness
- [`internal`](internal)
  - harness implementation packages

Skills are separate from prompt templates:

- prompts provide question-specific and shared benchmark instructions
- skills provide focused implementation guidance for a narrower task (visual dashboard generation)

Start with [`qforge.md`](docs/qforge.md) if you want to run benchmarks or compare providers.
Start with [`datasets/ontime/download/README.md`](datasets/ontime/download/README.md) if you want to rebuild or inspect the OnTime dataset.

Published run artifacts live in the sibling runs site: [ExploringDatabyLLMs-runs](https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/).
Example compare report: [q001 on 2026-03-20](https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-20/q001_hops_per_day/compare_report.md).
