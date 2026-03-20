# ExploringDataByLLMs

This repository contains two main active areas:

- [`qforge`](/Users/bvt/work/ExploringDatabyLLMs/docs/qforge.md)
  - the Go-based benchmark harness for LLM-generated analytics workflows
  - covers prompt assembly, SQL execution, report and visual generation, compare runs, and CLI usage
- [`datasets/ontime/download`](/Users/bvt/work/ExploringDatabyLLMs/datasets/ontime/download/README.md)
  - the dataset rebuild and loading workflow for `ontime.ontime`
  - covers source archives, table layout, loader commands, and schema rules

Related repo areas:

- [`prompts`](/Users/bvt/work/ExploringDatabyLLMs/prompts)
  - shared and question-specific prompt assets used by `qforge`
- [`Skills`](/Users/bvt/work/ExploringDatabyLLMs/Skills)
  - currently includes only the [`ontime-analyst-dashboard`](/Users/bvt/work/ExploringDatabyLLMs/Skills/ontime-analyst-dashboard/SKILL.md) skill for `visual.html` generation
- [`cmd/qforge`](/Users/bvt/work/ExploringDatabyLLMs/cmd/qforge)
  - CLI entrypoint for the active harness
- [`internal`](/Users/bvt/work/ExploringDatabyLLMs/internal)
  - harness implementation packages

Skills are separate from prompt templates:

- prompts provide question-specific and shared benchmark instructions
- skills provide focused implementation guidance for a narrower task (visual dashboard generation)

Start with [`qforge.md`](/Users/bvt/work/ExploringDatabyLLMs/docs/qforge.md) if you want to run benchmarks or compare providers.
Start with [`datasets/ontime/download/README.md`](/Users/bvt/work/ExploringDatabyLLMs/datasets/ontime/download/README.md) if you want to rebuild or inspect the OnTime dataset.
