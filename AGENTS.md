# AGENTS.md

This file applies to the entire repository rooted here.

## Working style
- Keep changes small and focused on the user's request.
- Prefer updating existing files over introducing new structure unless needed.
- Preserve the current style and naming patterns of nearby code.
- Do not touch large datasets, generated outputs, or notebooks unless the user asks.

## Files and edits
- Use `rg`/`rg --files` for fast code search.
- Read large files in chunks before editing.
- Reuse existing scripts, templates, and utilities when available.
- Update documentation when behavior, commands, or setup steps change.

## Validation
- Run the narrowest relevant check first, then broaden only if useful.
- Avoid fixing unrelated failures discovered during validation.

## Git hygiene
- Do not commit, create branches, or rewrite history unless the user asks.
- Keep generated files and local artifacts out of diffs; update `.gitignore` only when a new artifact source is introduced by your change.
