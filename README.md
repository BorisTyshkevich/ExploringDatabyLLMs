# ExploringDataByLLMs

This repository contains two main active areas:

- [`qforge`](/Users/bvt/work/ExploringDatabyLLMs/qforge.md)
  - the Go-based benchmark harness for LLM-generated analytics workflows
  - covers prompt assembly, SQL execution, report and visual generation, compare runs, and CLI usage
- [`datasets/ontime_v2/download`](/Users/bvt/work/ExploringDatabyLLMs/datasets/ontime_v2/download/README.md)
  - the dataset rebuild and loading workflow for `ontime.ontime`
  - covers source archives, table layout, loader commands, and schema rules

Related repo areas:

- [`prompts`](/Users/bvt/work/ExploringDatabyLLMs/prompts)
  - shared and question-specific prompt assets used by `qforge`
- [`cmd/qforge`](/Users/bvt/work/ExploringDatabyLLMs/cmd/qforge)
  - CLI entrypoint for the active harness
- [`internal`](/Users/bvt/work/ExploringDatabyLLMs/internal)
  - harness implementation packages

Start with [`qforge.md`](/Users/bvt/work/ExploringDatabyLLMs/qforge.md) if you want to run benchmarks or compare providers.
Start with [`datasets/ontime_v2/download/README.md`](/Users/bvt/work/ExploringDatabyLLMs/datasets/ontime_v2/download/README.md) if you want to rebuild or inspect the OnTime v2 dataset.
