# Edition 3 Chapter 2 — student correction round 1

## Status and exact candidate

Implemented only the three required initial-review corrections, R1–R3. This is a **candidate for independent re-review**, not a declaration of acceptance. No acceptance tag was created and no next chapter was started.

- Student repository: `/Users/bill/projects/ensemble-edition-3-student`
- New clean commit: **`89cc874b3c32e65ee43f44a72ad9ca48c9dd3323`**
- Reviewed parent commit, preserved: `74eb0265bd8bdbc7857d579704300b317b2365da`
- Accepted Chapter 1 predecessor/tag peel, unchanged: `1aad3975842d502b9144b7b9bddfd30ca7309db7` (`edition-3-ch01-r1`)
- New CLI: `/Users/bill/projects/ensemble/docs/edition-3-notes/ch02-student-r1-evidence/student-cli`
- New CLI SHA-256: **`b9b29b73e0da79393dc8c499aaee1437e73cfcb260da57b1ee1fef7f4325624a`**
- Source archive SHA-256 (`git archive --format=tar HEAD`, student root cwd): **`7c9452a578e718f3b86dc45b2ee574ea4027801a0c3685839ae68e27197feb08`**
- `go version -m` records Go 1.25.4, darwin/arm64, the new revision, and `vcs.modified=false`.

Build command, student root cwd:

```sh
go build -o /Users/bill/projects/ensemble/docs/edition-3-notes/ch02-student-r1-evidence/student-cli ./cmd/ch01
```

Only three source-repository files changed: `internal/engine/parse.go`, `internal/engine/render.go`, and the new `chapter2_revision_test.go`. Public architecture, ownership, composition root, interfaces, import boundaries, GUI stub, command surface, and configuration remain unchanged. No tool loop, new Context field, pairing framework, native Responses support, or unrelated feature was added.

## Findings, fixes, and public-root regressions

### R1 — chronological Gemini name resolution

The renderer previously collected every tool-call name before rendering any dialogue. Two ordinary sequential ID-less calls therefore shared the empty-string key, and the later name replaced the first before its result was rendered.

The map is now populated while visiting each tool call in chronological rendering order, before either its opaque replay shortcut or normal wire assembly. Result resolution sees only calls already visited, not calls later in the dialogue. The existing result fallback is unchanged. This retains the simple renderer-local map and adds no persistent identity state.

`TestChapter2GeminiChronologicalNames` uses the public root, parses two Gemini responses with absent IDs and distinct function names, appends each response followed by its result, dumps and reloads the log through another public-root Agent, and renders the full dialogue. It asserts `lookup_weather` then `lookup_stock`, preserved absent IDs in facts, deterministic rendering, and no mutation of the log.

### R2 — measured cost for unusable OpenAI/Gemini replies

OpenAI and Gemini assigned normalized usage after their empty `choices`/`candidates` checks. A successfully decoded response reporting 11 input and 2 output tokens returned an error without carrying those measured counters to Engine attempt accounting.

Both parsers now assign their existing normalized counters immediately after wire decoding/model identity handling, before the empty-list check. No counters are estimated, and normalization formulas are unchanged. The existing Engine accounting consumes the measured usage even when parsing returns an unusable-answer error. Agent error handling still does not append a completed response.

`TestChapter2MeasuredEmptyReplies` exercises actual root-owned HTTP clients against local `httptest` servers for **all three vendors**. It first commits a successful 3-input/1-output response, then receives two empty replies measured at 11/2 each. It checks:

- each failed Ask errors and returns no answer;
- cumulative Engine usage becomes 14/3 then 25/5;
- exactly one HTTP request per Ask, without retries;
- committed text history remains unchanged after failure;
- Context usage remains the successful response's 3/1;
- there is still exactly one `response_ended`, plus one `error_occurred` per failure.

Thus completed-response reduction and measured attempt accounting remain explicitly distinct. No fabricated `ResponseEvent` was introduced.

### R3 — multiple OpenAI opaque parts

The parser already represented `reasoning_content` and `reasoning` as separate opaque parts. The renderer assigned both destination fields from each part, including absent nil fields, so the later part erased the earlier value.

The renderer now updates only a field present in that opaque part. The existing exact vendor/model/surface check remains in place.

`TestChapter2OpenAIMultipleOpaqueParts` parses one response containing both already-supported fields, appends it through public History, dumps/reloads it, reduces Context, and renders it to the exact issuer. It verifies both fields survive in either opaque-part order, along with visible text. Separate copied-Context tests change vendor, model, or surface and require both opaque values to remain excluded. Stored provenance/facts must remain unchanged.

## Sensitivity evidence

The combined new regressions were run before production fixes. They failed behaviorally on all three reported defects; the Anthropic branch of R2 already passed. See `regressions-before-fixes.txt` and its exit file.

`check-sensitivity.py` copies only student source into an isolated temporary scratch directory beneath the new evidence directory. It independently reverts each fix, invokes its exact public-root regression, then restores the original bytes and reruns. It rejects compilation failures or selectors that run no tests. It hashes the candidate source before/after to prove production files were untouched. Scratch is removed after completion.

| Fix reverted in scratch | Regression | Reverted | Restored |
|---|---|---:|---:|
| R1 | `TestChapter2GeminiChronologicalNames` | behavioral failure, exit 1 | pass, exit 0 |
| R2 | `TestChapter2MeasuredEmptyReplies` | behavioral failures on OpenAI/Gemini, exit 1 | all vendors pass, exit 0 |
| R3 | `TestChapter2OpenAIMultipleOpaqueParts` | behavioral failure, exit 1 | pass, exit 0 |

Records: `R1-sensitivity.txt`, `R2-sensitivity.txt`, `R3-sensitivity.txt`, `sensitivity-summary.json`, `sensitivity-run.txt`/`.exit`. Final focused regressions were also rerun on the committed source with `-v -count=1`, all passing (`final-focused-regressions.txt`).

## Final validation on the new candidate

All grader commands below targeted the **same new CLI**, whose path/hash appear above. Both were run from `/Users/bill/projects/ensemble`, the authoritative root module, not `agent/`. Source remained frozen during validation.

| Cwd | Command | Result |
|---|---|---|
| authoritative root | `go run ./cmd/grade-edition2 -ch 2 -json /Users/bill/projects/ensemble/docs/edition-3-notes/ch02-student-r1-evidence/student-cli` | exit 0, `passed:true`, 100/100 |
| authoritative root | `go run ./cmd/grade-edition3 -ch 1 -json /Users/bill/projects/ensemble/docs/edition-3-notes/ch02-student-r1-evidence/student-cli` | exit 0, `passed:true`, 100/100 |
| student root | `go test -race -count=1 ./...` | exit 0; root suite 126.106 seconds |
| student root | `go test -count=1 -coverpkg=./... -coverprofile=/Users/bill/projects/ensemble/docs/edition-3-notes/ch02-student-r1-evidence/final-cover.out ./...` | exit 0; root suite 124.712 seconds |
| student root | `go tool cover -func=.../final-cover.out` | exit 0, merged 89.4% |
| student root | `go vet ./...` | exit 0 |
| student root | `git diff --check` | exit 0 |
| student root | `go list -f '{{.ImportPath}}: {{join .Imports " "}}' ./...` | exit 0; no sibling implementation dependencies; CLI imports public root |
| new evidence `import-smoke/` module | `go run .` | exit 0; external import/public Ensemble ownership checks pass |
| student `gui/` module | `go test ./...`; `go build ./...` | both exit 0; no GUI runtime claims |

The Edition 3 grader explicitly reports successful **1m2s response-header delay** and **1m2s response-body delay** checks. Neither was replaced by cached results or interrupted with a short deadline.

Source coverage denominator: **749 covered production statements / 838 selected statements = 89.379%, displayed as 89.4%**. The root module includes CLI source. The separate GUI module, external grader processes, standalone evidence invocations, and old live runs are not counted. **89 statements remain uncovered.** No numeric minimum is asserted; current rules specify no numeric coverage threshold. `coverage-denominator.json`, the raw profile, and function report are retained.

## Keyless and fake-protocol checks with the new CLI

`check-offline.py` uses a minimal environment without real provider credentials and a local HTTP sentinel. It never reads the authorized key store or any other credential source. Final run exit is 0.

- Each of the three original saved live logs dumps byte-for-byte unchanged in a fresh new-CLI process.
- Each original log renders to all three vendors twice under the new CLI (nine byte-identical pairs). Each pair also matches a fresh local render performed by the preserved old binary under the same target configuration. These are current **offline** comparisons, not new live calls.
- Synthetic R1 CLI fixture: two sequential ID-less call/result pairs produce the correct distinct response names.
- Synthetic R3 CLI fixture: exact-provenance rendering retains both opaque fields; another model receives neither.
- The original synthetic redaction log still omits its fake private result when rendered to each vendor.
- All these offline operations produce **zero HTTP requests** to the sentinel.
- Three additional actual CLI invocations use only a local fake HTTP server and the non-secret literal `fake-local-only`. Each vendor returns an empty response with measured 11-input/2-output usage. Each CLI exits nonzero as expected, reports measured cost in diagnostics, and saves exactly `message_received`, `request_sent`, `error_occurred`—no completed response. There are exactly three local fake requests and no retries.
- Original evidence files and the original review were hashed before/after; all remain byte-identical. The student worktree is clean after checks.

See `offline-checks-summary.json` for argv, cwd, process statuses, original artifact hashes and fake request paths. All stdout/stderr records and normalized synthetic logs are retained. No raw credential/header/provider payload archives were created.

## Original live evidence attribution — not rerun

**There were zero new paid calls in this correction round. The new CLI is not claimed to have been live-validated.** Original six live calls remain attributed to commit `74eb0265bd8bdbc7857d579704300b317b2365da` and its original binary. Their UTC capture times, as stored in the original logs, are:

| Provider/surface | Original event time range, 2026-10-16 UTC | Responses | Input / output |
|---|---|---:|---:|
| Anthropic Messages, `claude-haiku-4-5-20251001` | 21:01:06.414619–21:01:10.268633 | 2 | 279 / 197 |
| OpenAI ChatCompletions, returned `gpt-4.1-nano-2025-04-14` | 21:01:10.294747–21:01:14.154838 | 2 | 244 / 147 |
| Gemini generateContent, `gemini-2.5-flash-lite` | 21:01:14.177626–21:01:15.862874 | 2 | 192 / 120 |

All original cache-write/cache-read values were measured zero. `original-live-attribution.json` records these facts. None of R1–R3 invalidates those successful text-only captures: they contained no function pairs or reasoning-field handoff and were not empty failed replies. New offline comparisons show their rendered requests remain unchanged. This does **not** extrapolate the old live runs to untested new live behavior.

The original launcher usage-regex failure and the independently audited offline correction remain retained in the old evidence directory, without edits. This round does not relabel that old failed launcher exit.

## Assistance, input separation, failures, and feedback

Assistance received this round was the supplied initial-review message identifying R1–R3, relevant public chapter obligations, acceptance of the architecture in principle, the absence of a numeric coverage minimum, and permission not to repeat paid live calls for these local corrections. Those findings were reproduced through new tests before fixing. No further question or outside answer was needed.

Only permitted public chapter excerpts, current student source, and student-owned original evidence were inspected. No private review/evidence files, historical answers, grader internals, or future chapters were read. Graders were executed without inspecting their implementations. No authoritative repository content was staged or committed. Only the new review file and the new round-one evidence directory were written there.

Retained development/evidence failures:

1. New regressions initially failed on all three findings, as intended (`regressions-before-fixes.*`).
2. Two guarded edit scripts stopped on incorrect exact-text anchors before writing source; the unchanged failing regression rerun is retained in `edit-anchor-failure-regressions.*`. Subsequent exact edits passed. These were editing errors, not successful fixes.
3. One requested read used a nonexistent `internal/common/enums.go` filename; a directory listing identified `common.go`. No source changed.
4. The first local evidence harness incorrectly passed `--protocol` instead of using the documented default JSONL mode. The CLI correctly rejected it before any fake request. The original script, failed output/exit, and stderr are retained as `first-attempt-*`; source was not changed for the harness error.
5. A retry stopped on its clean-worktree guard because `go build ./...` in the GUI module produced the untracked executable `gui/ensemble-gui`. That exact generated artifact was moved to the new evidence directory as `gui-stub`; no reset, clean, or broad deletion was used. The stopped run is retained as `second-attempt-*`. The subsequent local run passed.

Concise teaching feedback: the existing requirements were sufficient. The gaps were in student test selection, not in a need to change the chapter or grader. In addition to one round trip, sequential ID-less calls expose accidental future lookups; cross-vendor empty-response tests distinguish measured attempt cost from completed conversation facts; and two independently preserved opaque fields expose destructive merge behavior that one-field tests miss. The scratch revert checks demonstrate that the new tests discriminate these faults. Explicitly checking the real CLI syntax before composing an evidence harness would have avoided this round's harness-only failure.

## Remaining limits

No new live provider behavior, nonzero live caching, live function execution, streaming, media loading, or GUI runtime connection was tested or added. R1 validation covers the required ordinary sequential ID-less pairing, not a new generalized correlation scheme. R2 distinguishes fully decoded measured empty replies from successful committed responses; no raw failed wire responses are archived. R3 covers the two already-supported ChatCompletions reasoning fields and retains exact provenance exclusion. Source coverage is not exhaustive, as quantified above.

Final source worktree and source/binary hashes were rechecked after all execution. Original candidate/report/evidence are preserved. Ready for independent exact-revision re-review of `89cc874b3c32e65ee43f44a72ad9ca48c9dd3323`.
