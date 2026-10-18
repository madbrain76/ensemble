# Chapter 9 — ownership/test plan (approval required)

## Basis and boundary

Read the complete current public `AGENTS.md`, `.agents/skills/ensemble-coding/SKILL.md`, `.agents/skills/edition-3-coder/SKILL.md`, `docs/edition-3-execution.md`, `docs/edition-3-student.md`, `docs/edition-3-carryover.md`, and `book/edition-3/chapter-09.md`, recovering clipped sections. The current student brief supersedes historical workflow/grader commands. Assistance: these public instructions and `/tmp/cadence-ch09-student-plan-task.txt`; not independent invention.

Read-only inspection of `/Users/bill/projects/ensemble-edition-3-student` confirms clean status, HEAD **b8f22ba5b3c9d1fc52d1a114810b90b77928684d**, tree **fa896caa404b87b08de61eb38ce079eb38138fbe**. Annotated `edition-3-ch08-r1` object **bbae919d6c595f879eb27f602155c975e1cb8489** peels to that HEAD. There are 126 tracked source/support files, not 126 production files. No private evidence, historical implementations, credentials or grader internals were inspected. No builds/tests/grades or live exercises ran. No genuine technical blocker; implementation and later live execution await root authorization. `rg` is unavailable; standard grep suffices.

## Smallest ownership change

- **Composition:** extend `ensemble.go`'s existing Agent construction inputs with skills directory, selected primary and application-variable callbacks. Keep no-spec construction backward compatible. Translate `EN_SKILLS_DIR`, `EN_PRIMARY_SKILL`, `EN_CUSTOM_VAR` in `cmd/ch01/main.go` and the GUI construction path (`gui/main.go`), not inside the registry. No CLI flattening.
- **Agent:** add an Agent-owned registry in a small `internal/skills` package, connected through `internal/common` only as needed. It owns discovered definitions, active/pending sets, variable renderers and first-enabling metadata; no global mutable registry. Extend `internal/agent/agent.go`, `dispatch.go` and the existing actor coordination only where necessary. Keep one tool admission path, one durable event log, existing jobs/results and single-Agent history ownership.
- **Contract:** parse only the taught scalar/frontmatter and inline/block list forms; preserve body horizontal rules; validate required boundary/name/type, duplicate identities and path-bearing errors. Distinguish primary, loadable and dependency-only skills. Validate the entire requested dependency graph, tool names, types, availability and cycles before any activation. Activate dependencies before their dependent; failed loads change nothing. Creation loads the primary graph and renders immutable system text once.
- **Lifecycle:** install `load_skill`/`unload_skill` through existing tool machinery; enabled capabilities govern both declarations and actual dispatch, including control-tool shortcuts. Hidden guessed names cannot execute. Offered loadables provide progressive disclosure; do not mistake dependency reachability for public loadability. Successful loads return dependency/body instructions through ordinary tool results. Repeated loads are idempotent and preserve first provenance. Initial tool sequence is zero; dynamic provenance uses the actual enabling `load_skill` `tool_called` sequence, not a new counter. Unload only marks pending: instructions, declarations and admission remain intact.
- **Rendering/Engine/GUI:** render `$TOOLS`, `$SKILLS` and application variables at point of use, preserve unknown placeholders, install initial callbacks before primary rendering, and invoke callbacks outside registry/Agent locks. `RegisterVar` affects future loads, never frozen system text. Engine retains its existing renderer/HTTP path, consuming the Agent's enabled declarations without owning skills. GUI retains `hub.go` call/result replay and `app.js`/`tool-readable.js` action cards; only construction plumbing is expected. No new panel, compaction, memory or caching feature.

## Later verification targets

Run sequential `go test ./...`, `go test -race ./...`, `go vet ./...` in these actual modules: root, `gui`, `examples/ch04-consumer`, `examples/ch06` (workflow), `examples/ch07` (stream exercise). Build actual CLI `./cmd/ch01`, GUI `./cmd` within `gui`, and each example's `.`. Retain inherited budget/settings, streaming/media, jobs, pause/replay and consumer behavior; no invented coverage threshold.

Add focused parser/registry tests plus public Agent integration tests for both list forms/body rules; malformed/duplicate/type errors; dependency ordering, missing dependencies/cycles and rollback; two-Agent isolation; immutable primary text; dynamic lists/variables and reentrant callbacks; hidden dispatch refusal; loadable progression; repeat-load provenance and pending unload. Use local fake HTTP to inspect real rendered vendor requests and ordinary load results, actual process fixtures to prove executable capability changes, and GUI WebSocket tests through the existing prompt/replay wire path. Sensitivity controls in owned copies: bypass dispatch admission and separately mutate primary text on load; require the corresponding behavioral assertions to fail and pass again after restoration. Compilation/setup failures do not count.

After freezing one CLI, run the current public adapter from the book root (absolute paths below, not the chapter's historical command):

```sh
S=/Users/bill/projects/ensemble-edition-3-student
# BIN is the frozen CLI in owned ignored local storage.
go run ./cmd/grade-edition3 -ch 9 -json -core "$S" \
  -consumer "$S/examples/ch04-consumer" -workflow "$S/examples/ch06" \
  -stream-exercise "$S/examples/ch07" -stream-model claude-sonnet-4-6 \
  -gui "$S/gui" -gui-model claude-sonnet-4-6 "$BIN"
go run ./cmd/grade-edition3 -ch 1 -json "$BIN"
```

Both require `passed=true` AND exit 0 on that SAME binary. The chosen model is already declared. Allow the separate real 62-second header/body waits plus suite overhead; a 90-second observation is insufficient and silence is not failure.

## Manual evidence and handoff

After a separate pinned finite plan and root execution authorization, exercise the chapter's scratch primary/code-tools/hidden-skill example through CLI and the existing GUI: see load instructions, use the newly enabled file tool on an intended scratch file, reject a hidden skill, mark unload pending and demonstrate continued capability; check readable/replayed action cards and inherited controls. Use one authorized declared model if needed, not provider discovery or a three-vendor campaign. No browser/native/audio/MIC/provider access now.

Retain failures, commands, exit statuses and exact source/test provenance. Keep binaries, archives and private raw traces in owned ignored local storage, never portable main evidence; publish only suitable summaries. Later builds are sequential and stop below 1 GiB free without deleting shared caches. Freeze source commit/tree and binary hash; rerun affected gates after changes. Final source commit only if permitted; root alone accepts/tags. Formal correctness rounds: 0 of at most 3. No children or manuscript work. This turn writes only this plan and stops for approval.
