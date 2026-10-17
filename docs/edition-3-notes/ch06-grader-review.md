# Edition 3 Chapter 6 — independent adapter review

## Verdict and boundary

**PASS — the exact frozen adapter candidate is accepted for artifact plumbing. No required adapter fix found.** This is not Chapter 6 student acceptance, a new entry audit, a complete streaming-coverage certification, or authorization to launch the student. Root must inspect the saved final tail and join this review before accepting the verdict.

Review date: 2026-10-17. Independent reviewer retained from the completed entry audit. No adapter, grader, application, student, manuscript, public-guidance, or original-solution edits were made by this reviewer. Only this report, `ch06-grader-review-evidence/`, and owned `/tmp/ensemble-ch06-adapter-review-c6v5ft04` artifacts were written. Normal toolchain/grader temporary files are managed by those tools. No children, paid calls, credential reads, GUI, tags, commits, pushes, or cache cleanup.

## Exact inputs and prerequisite coverage

Main remains at `b67386ab80cf2afe7985db5a32ab99ec2311f090`; its tracked diff remained exactly the root-supplied two-file adapter/test candidate throughout this review. It is not described as a clean tracked tree. Protected `solutions/edition-3` FINAL tree remains `4ff207e66df6d49cbdea462df2f6753d421779d8`.

| Artifact | SHA-256 |
|---|---|
| `cmd/grade-edition3/main.go` | `4227ad0a0d96f17b541b60f85e7ac19acd09edfbb79e50c35153aab58f7c5ee7` |
| `cmd/grade-edition3/main_test.go` | `6f443849bef87fcf2b5e97566bff545088e9919db15212450560d76f11c9159a` |
| Exact candidate `/tmp/cadence-ch06-grader` | `8818400c89bb1663cad72faeca99bc92d64fb9a26e9f023a54d5e319c0e2a798` |
| Existing Chapter 6 grader | `edfd0ef2ad9e9841bf531ba552a8983ac182fa8b1f61612596343aa102d080b7` |
| Public Chapter 6 | `ed007e7e7059400939a6735534d7a953d4575e83230a6f68fc35583e8a17bd55` |
| Accepted f60 CLI | `e51138885a359eaf801258a52bc21f28ca41691683256647f8c9911be410b233` |

Candidate hashes were checked before and after measurement. Root's recorded build, pre-edit pins, initial guard failure, initial test file, and final test-only refinement are preserved in the evidence. The production source did not change after the candidate build. This review did not rebuild or substitute a different production executable.

Current mandatory instructions were refreshed: full `AGENTS.md`, `.agents/skills/ensemble-coding/SKILL.md`, `.agents/skills/edition-3-code-reviewer/SKILL.md`, and the complete public `book/edition-3/chapter-06.md`; oversized instruction output was recovered. The retained full entry-audit reads cover `docs/edition-3-execution.md`, readiness, carryover, public student brief, inherited Chapters 4/5, and the required full historical workflow, Edition 2 carryover, voice guide, and chapter-writing procedure. Their hashes were verified unchanged; no compaction/replacement of those prerequisite reads occurred between these bounded turns. Current student-first instructions supersede historical author-first phases. No runtime skill registration was invented.

The original entry report and its exact 98-file evidence set were verified byte-for-byte unchanged, including all 97 checksum entries. The accepted student remains clean at `f60ebf3aad2ece4251face9d95c106dd667abcd6`, tree `3de5d0d24dd9f593ec017fceccfbff1b209c5d57`, tag `edition-3-ch05-r1` resolving to that commit. Its entire saved source manifest still matches. No Chapter 6 student exists or was launched.

## Independent source/test review

The production change is small and direct: 27 added/12 removed lines, 111 to 126 total lines. Tests add 71/remove 1 line. There is no generic dispatch framework, new grading point, historical grader edit, or application implementation.

- Adds only Chapter 6 recognition, two Chapter 6 flags, existing-style validation, and the existing `Chapter6(bin, core, consumer, workflow, exercise, model, example, models)` call.
- Explicitly requires core, consumer, workflow, and streaming exercise directories with regular `go.mod` files, plus a nonblank model. It does not infer source identity from the positional executable target.
- Reads optional consumer-case/media-model JSON with the existing readers and forwards both resulting values. No replacement schema or interpretation is introduced.
- Uses `FlagSet.Visit`, so the new flags are rejected on Chapters 1/4/5 even when explicitly empty. Presence is not confused with a nonempty value.
- All chapter/argument/root/model/JSON validation is before `grade.Build`; help returns 0 during flag parsing before building or launching anything.
- Preserves the existing build, cleanup, report encoding, error handling, and pass/fail exit conventions. A textual comparison verifies the entire build/dispatch/output/exit suffix is exactly the old suffix after removing the single added Chapter 6 branch.
- Root's refined tests correctly isolate new-flag rejection on otherwise valid older-chapter argument sets. The initial test version's weaker isolation is not credited as runtime proof. Independent compiled-sentinel checks below provide that proof on the exact executable.

The explicit-model requirement is appropriate. The pinned historical source declares `claude-sonnet-4-6` with `StreamAll` in `internal/engine/media.go:LookupModel`. The adapter supplies no invented capability row or historical default. The public contract remains **one `Parse(resp, callbacks)` seam, response-scoped parser IDs including zero, `EN_DISABLE_STREAMING`, and `CH07_NO_STREAM`**; the adapter changes none of it. No public clarification or guidance edit is required by this delta.

## Exact-candidate measurements

All graders ran sequentially. `run.py` records full argv, cwd, controlled environment, UTC start/end, exit, duration, free disk, and unabridged stdout/stderr with hashes. It has no wall-clock cancellation deadline. Only grader-provided local fakes were used. Minimum recorded free space was **2,262,949,888 bytes**; no cleanup was performed.

The historical positive is the unchanged pinned scratch source qualified by the entry audit, not an Edition 3 student implementation. Its binary, core module, consumer, workflow, and streaming module are distinct arguments:

```sh
/tmp/cadence-ch06-grader -json -ch 6 \
  -core /tmp/ensemble-ch06-entry-7_7zbtqs/solutions/ch06 \
  -consumer /tmp/ensemble-ch06-entry-7_7zbtqs/solutions/ch06/ch05 \
  -workflow /tmp/ensemble-ch06-entry-7_7zbtqs/solutions/ch06/ch06 \
  -stream-exercise /tmp/ensemble-ch06-entry-7_7zbtqs/solutions/ch06/ch07 \
  -stream-model claude-sonnet-4-6 \
  /tmp/ensemble-ch06-entry-7_7zbtqs/positive
```

Working directory: `/Users/bill/projects/ensemble`. Historical directory names are private fixture identities, not prescribed student paths. Full historical commit/tree/source manifests remain pinned in entry evidence; final verification reconfirmed every source byte. Exact binary hashes are also in `source-binary-pins.json`.

| Recorded run | Score | Passed | Exit | Seconds |
|---|---:|:---:|---:|---:|
| `positive` | 100/100 | true | 0 | 13.430 |
| `wrong-stream-model` | 35/100 | false | 1 | 12.340 |
| `wrong-stream-exercise` | 90/100 | false | 1 | 12.084 |
| `custom-case-media` | 80/100 | false | 1 | 12.484 |
| `legacy-ch4` | 160/160 | true | 0 | 8.600 |
| `legacy-ch5` | 100/100 | true | 0 | 12.071 |
| `legacy-ch1` | 100/100 | true | 0 | 127.147 |

The complete measurement sequence exited 0: each intentional negative returned its expected exit 1, and each positive returned 0. Scores are reported on their actual scales, not normalized or reinterpreted.

### Meaningful forwarding sensitivity

1. **Stream model:** a nonblank, syntactically valid but unlisted `adapter-review-unlisted-stream-model` reached the unchanged grader rather than being rejected by the adapter. It lost `stream-deltas`, `deltas-match-final`, `tool-params-streamed`, and `delivery-not-content`. Actual diagnostics included `saw 2 text deltas, want at least 8`, both modes producing five deltas, and request streaming flags disagreeing with the selected delivery. This is expected plain fallback for an unlisted model, used deliberately as a wrong-model exercise control—not a public-contract defect in unknown-model fallback.
2. **Stream exercise:** a valid `go.mod` directory with no source or core dependency passed adapter shape validation and then lost `delivery-not-content`. Actual diagnostic: `stream exercise dependency: go list in /tmp/ensemble-ch06-adapter-review-c6v5ft04/module: exit status 1: no Go files ...`. This proves the supplied exercise path reaches structural/exercise qualification. No student/example implementation was created to force this result.
3. **Consumer case and media map:** valid JSON supplied a deliberately wrong expected consumer result/diagnostic and an unlisted Anthropic media model. The single pertinent joint control lost `ch5-parity` with **both independent nested observations**: actual `HELLO FRAMEWORK` versus expected `ADAPTER CASE FORWARDING EXPECTATION`, absent custom diagnostic text, and `model "adapter-review-unlisted-media-model" does not support image media`. Thus neither optional input was silently replaced with defaults. This is a control failure, not a candidate defect.

These runtime negatives complement, rather than replace, direct inspection of all eight forwarded API arguments.

### Pre-build/launch gates on the exact binary

A compiled Go sentinel writes an owned marker and exits 99 when launched. Its direct positive self-check did so; the marker was moved aside and retained. **All 61 subsequent candidate gate cases produced the expected exits and left the launch marker absent.** Full argv/output/exit/duration receipts are in `runtime-gates.json`.

Covered: missing and explicitly empty required values; missing roots; roots without `go.mod`; `go.mod` directories rather than regular files; a file used as a root; whitespace-only model; malformed/missing/directory/schema-invalid optional JSON; each new flag on each of Chapters 1/4/5 with both empty and nonempty values; help with exit 0; missing/multiple targets; flags placed after the positional target; unknown flags; unsupported/malformed chapter selectors. Invalid cases exited 2. Bad optional JSON combined with a nonexistent target returned the JSON parse error, independently confirming validation takes precedence over target building. Source order independently establishes that all validation precedes `grade.Build`.

### Existing-route preservation and package checks

Fresh sequential `go test -count=1 ./cmd/grade-edition3`, `go test -race -count=1 ./cmd/grade-edition3`, `go vet ./cmd/grade-edition3`, and `git --no-pager diff --check` all exited 0. Sentinel compilation exited 0.

The exact candidate's actual Chapter 4 and 5 routes passed on the matching historical positive. Chapter 1 used the **same accepted f60 frozen CLI**, freshly reverified by SHA-256. Its actual header-delay and body-delay checks each waited **1m2s** and passed; the full run finished naturally in 127.147 seconds. No shortened delay, changed lifetime, fake fresh result, or test cancellation was used. Existing argument/dispatch semantics were also inspected against tracked HEAD; prior entry measurements are cited only for fixture qualification and unchanged source identity, not presented as new runtime results.

## Failed setup and test-proof history

Preserved without credit as passing evidence:

- Root's initial guarded edit failed before writing because its expected help/implementation anchor was wrong. No grading or production answer code occurred under that failure; `root-initial-guard-failure.txt` is retained.
- Root's first passing test set did not fully isolate older-chapter new-flag rejection. The initial test source, refinement diff, original receipts, and final race/vet/whitespace rerun are retained. Production remained unchanged after its first build.
- The review's first gate harness expected an unnecessarily specific error prefix. The candidate correctly returned exit 2, `grader: unexpected end of JSON input`, without launching the sentinel. The harness assertion stopped the sequence before any grading. Initial script and all receipts are retained; the corrected completed 61-case matrix is the credited run.
- An unrun measurement draft omitted `-json`; this was caught and corrected before any grading. Both versions are retained.
- Read-only setup mistakes—two interactive Git pagers, default temporary-directory selection before any scratch contents, unavailable `rg`, and one incorrect observation handle—are recorded in `setup-notes.md`. Only owned read-only pager jobs were terminated. No grader/test was interrupted, duplicated, repaired, or hidden to obtain a score.

## Limits and carry-forward

This adapter does not improve or certify the existing grader's remaining coverage. The entry report's required public/source/integration work remains binding: actual stream-arrival explanation before tool execution; gated early human CLI flush; SSE framing; complete exact final correlation; unknown-model plain fallback; source-bound opaque replay and signed ordinary text; request/body lifetime and cancellation; truthful aggregate exercise behavior; and accepted Chapter 5 wake-up/parked-batch interruption under streaming. Passing this historical fixture is not proof of a future student's implementation or of all those properties.

No provider/mutation campaign, new student implementation, or application design review was performed. The wrong-exercise probe covers a source-less wrong module rather than every possible module-identity mismatch. The combined optional-input control demonstrates both forwarded values through distinct diagnostics, not every valid configuration. Internal grader-built transient exercise binaries are not falsely claimed as individually captured; their pinned source/build commands are retained, while the candidate, historical core, f60 CLI, and compiled sentinel have exact binary hashes.

The source-only dependency/module check intentionally remains in the historical grader; the adapter requires only explicit regular `go.mod` roots. No arbitrary coverage threshold, extra grading point, private folder-name requirement, or new public obligation is introduced.

## Frozen handoff

Evidence index: `ch06-grader-review-evidence/INDEX.md`. `commands.json`, per-command `.meta.json/.stdout/.stderr`, `runtime-gates.json`, `results.json`, `source-binary-pins.json`, `source-comparison.json`, and `final-verification.json` contain the receipts. `REPORT-SHA256` binds this report; `SHA256SUMS` binds the exact evidence file set. Full final saved tail and manifest hashes are printed at submission.

**Final submitted verdict: PASS for this exact thin adapter only; no required fix.** All planned unit/build/behavioral/observation jobs are complete; no children or pending jobs remain. Student, historical material, protected FINAL tree, public guidance/manuscript, original graders, and entry report/evidence remain unchanged. Independent context is reserved for later adapter/code review, not student authoring.
