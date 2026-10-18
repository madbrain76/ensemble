# Chapter 8 grader review — PASS for the pinned candidate

## Verdict and boundary

**PASS. No blocking candidate defect found; no production correction requested.** This verdict covers only the frozen thin `-ch 8` publication adapter and initial `max_tool_rounds` typed-presence check. It is not Chapter 8 student acceptance, browser/visual/TTS acceptance, or authorization to widen grading. No student session or implementation was launched. No repository source, grader, manuscript, guide, settings, keys, existing evidence, or unrelated work was edited. All runtime controls ended before sealing.

The assignment was `/tmp/cadence-ch08-grader-review-task.txt`. Current AGENTS, coding/reviewer skills, shared execution/readiness guidance, Chapter 8, inherited Chapter 7, and all three candidate files were read fully. The other previously fully read governing prerequisites were carried only after unchanged-pin verification; `pins-before.json` identifies all 16 unchanged input entries and the three declared candidate overrides. The accepted qualified entry reconciliation remains unchanged.

## Exact reviewed candidate

| Item | SHA-256 |
|---|---|
| `cmd/grade-edition3/main.go` | `5c03c0a92fc95b7ee7707ac54cbd560cd9a7da3da215acf5271cbf777972df63` |
| `cmd/grade-edition3/main_test.go` | `7a2ec5f820a831b3eea86727545cf51a43c69876e144a26c02864d863a5986b7` |
| `internal/gradeedition2/ch08.go` | `aa5c10df77ed05edf2cfb6ec82debee88fae4fe23849734ed3b7b8016b56c92d` |
| Frozen `/tmp/cadence-ch08-grader` (10,831,346 bytes) | `529187f45b272aeb025670651fa9157ea1f7e1262170ae5593a517afa07a0645` |

Pins matched before and after review. Existing candidate/diff/root receipts were not rewritten. The frozen publication binary—not a rebuilt substitute—performed all independent publication controls.

## Source review findings

- **Dispatch and identity forwarding:** `cmd/grade-edition3/main.go` adds 8 to existing help, accepted-chapter and validation conditions, then calls `gradeedition2.Chapter8(bin, core, gui, guiModel, consumer, workflow, streamExercise, streamModel, example, models)`. The order matches the existing exported signature. There is no new target inference, historical path default, merged core/GUI identity, fallback model, or copied evaluator. Existing Chapter 7 dispatch remains unchanged.
- **Early validation:** explicit module directories with regular `go.mod`, nonblank stream/GUI model declarations, fixture JSON parsing, target cardinality, and chapter-specific flag-presence restrictions precede `grade.Build`. Empty values still count as supplied flags when rejecting flags belonging to a later chapter. The new chapter follows the established cleanup, JSON/text report, failure-exit, and error-exit paths.
- **Narrow completeness fix:** `internal/gradeedition2/ch08.go:152–160` now requires initial `max_tool_rounds` to decode as a JSON number and satisfy the existing integer-representation idiom used for font size. Missing, null, string, bool, and fractional values do not satisfy it. This adds no authored product default, clamp endpoint, budget-setting API, or new score bucket. The existing four settings weights remain 25/20/20/15, with inherited Chapter 7 parity worth 20. Explicit zero/false, partial-update preservation, cross-client broadcasts and restart sequencing are unchanged.
- **Tests and maintainability:** the candidate extends the existing validation tests rather than adding a framework or replacing a real target with a canned successful grader. No new public API, goroutine, process-lifetime policy, persistence implementation, or dependency is introduced. The original narrow scope is maintained.

## Independent runtime outcomes

All commands were sequential, local, deadline-bounded, credential-scrubbed and offline for Go dependencies. `clean-env.json` fixes `GOMAXPROCS=2`, `GOFLAGS=-p=1 -mod=readonly`, `GOPROXY=off`, isolated HOME/TMP/work directories, and known existing caches. No browser/native UI, paid/live provider, live-8088 server, agent child, shared cleanup, or standalone Chapter 1 lifetime rerun was used.

### Complete publication routes

| Frozen binary route | Score | Passed | Exit | Seconds |
|---|---:|---|---:|---:|
| `-ch 8` historical positive | 100/100 | true | 0 | 17.295 |
| `-ch 7` forwarding regression | 100/100 | true | 0 | 16.274 |

Both reports are retained in full. Chapter 8 contains all four settings checks plus successful `ch7-parity`; Chapter 7 contains its existing GUI checks plus successful `ch6-parity`. These are complete existing routes, not settings-only replacements. Neither timed out or left a process group.

Both used the hash-verified retained historical root CLI (built from `solutions/ch08`, target `.`), core `solutions/ch08`, separate GUI `solutions/ch08/gui`, consumer `solutions/ch08/ch05`, workflow `solutions/ch08/ch06`, stream exercise `solutions/ch08/ch07`, and explicit GUI/stream model `claude-sonnet-4-6`. Explicit consumer-case and media-model fixture bytes were copied from hash-verified reconciliation artifacts. Exact argv, absolute paths, environment, executable/input hashes, scores and exits are in `historical-ch{8,7}.receipt.json`, their stdout reports, `pins-before.json`, and the two fixture JSON files. `builds.jsonl` records internal build targets and resulting binary hashes.

### Real historical GUI omission sensitivity

This was not the earlier wire-only proxy. A copy of the actual `solutions/ch08/gui` was compiled against the unchanged historical core. Only scratch module resolution and the initial replay payload were changed. In `gui.go`, while no settings file exists, the copied GUI omitted `max_tool_rounds` from `current_settings`; update/broadcast payloads and persistence code were untouched. After settings were saved, restart snapshots were complete. Byte restoration removed the mutation.

| Actual real-GUI stage | Compile | Score / Passed | Behavioral exit | Result |
|---|---|---|---:|---|
| Omission mutant + retained pre-candidate checker | success | 80/80, true | 0 | Old blind spot reproduced |
| Same omission mutant + candidate checker | success | 60/80, false | 1 | Only `settings-on-connect` failed |
| Byte-restored GUI + candidate checker | success | 80/80, true | 0 | Real positive restored |

The two mutant GUI build hashes are identical (`ac0181f44456df232d11014661b92e3aa2edd06a1508821668eb716841dea254`), excluding a target-build mismatch as the explanation. Original/restored `gui.go` is `02487da66c5e0cb1d53cb958261cfff70222c89d032744544fb40e4bf1746eb2`; mutant source is `d1c564d21463876554cc8485052d3dae8df4b35e3159a535e5505519cea10b12`. Source snapshots, complete reports, compile receipts, caller sources and old/new checker provenance are retained. No repository implementation was mutated.

### Compiled no-target/build sentinels

A compiled marker was first executed successfully in both modes: target marker exit 98 and `go` build-interceptor marker exit 97, each leaving its expected receipt. Then **all 32 frozen-adapter cases passed**, with **no target or build marker**:

- Help for Chapters 8, 7 and 1: exit 0.
- Chapter 8 missing core/consumer/workflow/stream/GUI identity; empty and whitespace-only stream/GUI model; nonregular `go.mod`; malformed consumer/media JSON; missing/extra target: exit 2 before build/target execution.
- Unsupported Chapters 2 and 9: exit 2.
- Representative Chapter 7 missing-GUI/model validation and older Chapters 1/4/5/6 rejection of later-chapter flags even with empty values: exit 2.

The fixed case list, every argv/exit/stdout/stderr, marker self-test proof and results are in `sentinel-plan.json`, `sentinel-results.json`, and named receipts. This is not a new legacy-wide behavioral matrix.

## Carried evidence and limits

Root's pinned selected `go test -race -count=1`, `go vet`, and publication build all exited 0. Their exact commands, environments, times and output hashes are cross-checked in `carried-root-controls.json`. They are **carried root evidence, not independently rerun results**. Root's existing preparation-failure disclosure is referenced, not hidden.

Standalone legacy Chapter 1/4/5/6 full publication gates were **not rerun** and no independent full-gate claim is made for them. Source continuity, carried selected validation tests, representative compiled invalid-input controls, and the complete Chapter 7 inherited parity run support this narrow regression review. No first-lifetime/cancellation rerun was substituted for those qualifications.

Settings snapshots/restart and inherited wire checks are not proof of effective next-turn budget execution, every settings/load/write-failure path, DOM behavior, visual preferences, browser reconnects or audible speech quality. Those remain existing student implementation/code/browser acceptance obligations. No new API/setter prerequisite, fresh Bill-approval gate, malformed-input/write-failure campaign or grading framework is recommended.

## Preserved failures and evidence integrity

Three reviewer bookkeeping errors are retained with their original scripts/receipts in `preparation-failures.json`: (1) initial setup used a nonexistent manifest path key and stopped before controls; (2) the mutation wrapper misspelled the established check ID, stopped after the valid negative, and restored source in `finally`; only the remaining restored stage was then run; (3) sealing first paired root logs by a zero-based filename guess, then matched their one-based files to recorded hashes. No failed candidate behavior was retried into a positive. Missing guessed read paths are separately qualified in `read-lookup-notes.json`.

`verification.json` verifies 5,098 protected baseline files, all 19 input pins with only the declared three candidate overrides, 120 student files, 25 protected tracked files, unchanged commits, and a clean student worktree. All 41 owned command receipts are terminal without timeouts or remaining process groups. Final status/resource checks are in `final-seal.json`.

Portable outputs are only this new report and `ch08-grader-review-evidence/`. `inventory.json` lists every portable member's hash/size, excluding itself. Portable members must be UTF-8 text without NULs; no executable/archive/pycache payload is exported. `local-artifacts.json` separately inventories retained owned scratch artifacts (36 entries, 16,486,982 bytes including symlink target accounting), including compiled binaries and the local-only NUL-delimited status baseline. Prior retained artifacts remain untouched and have separate provenance. No shared or unrelated artifact was deleted.
