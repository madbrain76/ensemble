# Edition 3 Chapter 2 — correction-round-one code review

**Disposition: ACCEPTED for the frozen candidate below. R1–R3 are closed; no further required code corrections.** This is the independent code-review decision, not a tag, commit, or authorization to bypass the parent's final evidence/closure check. One of three permitted student correction rounds has been used. No author/text phase was started by this review.

## Frozen identities and scope

- Student candidate: `89cc874b3c32e65ee43f44a72ad9ca48c9dd3323`; initial review candidate: `74eb0265bd8bdbc7857d579704300b317b2365da`; accepted Chapter 1 predecessor: `1aad3975842d502b9144b7b9bddfd30ca7309db7`.
- Authoritative repository: `5351583208988f8370937cd0494727645efa8b46`; unchanged Chapter 2 SHA-256: `1a67ee9e3b6ae5f0f776a30f58908c6d53b53be2d8136ef44661728a677dc2d2`.
- Candidate source archive SHA-256: `7c9452a578e718f3b86dc45b2ee574ea4027801a0c3685839ae68e27197feb08`.
- Submitted and independently rebuilt CLI SHA-256: `b9b29b73e0da79393dc8c499aaee1437e73cfcb260da57b1ee1fef7f4325624a`. Explicit build target: `./cmd/ch01` in the student root module. Build metadata names the candidate revision and `vcs.modified=false`.
- Fully read round-one report: `ch02-student-review-r1.md`, SHA-256 `e2220672d634f5423b24690d531c902ab768d2d8d8f7ae61fb84cb82b33836d4`; fully read index: `ch02-student-r1-evidence/README.md`; parent pin: `ch02-candidate-r1.json`.
- Prior corrected review: `ch02-code-review.md`, SHA-256 `b57f390b22ebb905acefafd9b309c2f76a6447d0fe54411f63dc54b34312675b`. Initial report and all original evidence remain unchanged.

The current reviewer skill was reread. All governing files previously read in full, including the current chapter, were checked against the prior verified hashes; none changed. This is a focused re-review of R1–R3 and regressions their fixes could introduce, not a reopened architectural or feature review. Only `internal/engine/parse.go`, `internal/engine/render.go`, and new `chapter2_revision_test.go` differ from the initial candidate. Source/test hashes:

| File | SHA-256 |
|---|---|
| `internal/engine/parse.go` | `cdc363ac873c68bb0e318a9f87e370161bbf2a48153c5da3b72521469a2d47f0` |
| `internal/engine/render.go` | `c09a8df325d6432f88e685fdf6f87e4294281c90799850235a58fc2638fd0d96` |
| `chapter2_revision_test.go` | `a1bf6eca87441fc806a09987db7e4e94eb98a6e2015541dc811a6819e39b43fb` |

Full source hashes and binary/archive identity are retained in `ch02-code-review-r1-evidence/identity.json`; final checks are in `final-integrity.json`. The original 126-entry and new 135-entry submitted manifests independently verified, including sizes and SHA-256 values. The prior review's 63 checksum entries also verified. Counts are derived from actual files, not inferred from the student's prose.

## Findings closed

### R1 — chronological Gemini result names: closed

`internal/engine/render.go:392–400` updates the call-name lookup as each call is visited, **before** the exact-provenance opaque replay shortcut; `:422–430` resolves a result from that already-visited state. `:442–454` starts an empty lookup and processes dialogue in order rather than pre-populating it from future calls. This fixes the demonstrated ordinary sequential ID-less pairing without adding a Context field, changing issued IDs, or inventing a generalized correlation mechanism.

Public basis remains `book/edition-3/chapter-02.md:1003–1008`: recover the result's function name from its matching call earlier in the dialogue. New public-root `TestChapter2GeminiChronologicalNames` (`chapter2_revision_test.go:18–75`) parses native ID-less calls, dumps/loads real facts, checks both result names, repeated deterministic rendering, and log immutability. The unchanged reviewer public regression now passes. Reverting only R1 in scratch reproduces `[lookup_stock lookup_stock]`; restoring the fix passes. No new required invariant violation was found in the changed path.

### R2 — measured empty-reply usage: closed

`internal/engine/parse.go:129–133,208–211` now normalizes the already-decoded measured usage before rejecting empty choices/candidates. The existing formulas, resolved provenance, and answer-error behavior remain intact. Unchanged `internal/engine/engine.go:106–112` accumulates all returned measurement once; `internal/agent/agent.go:99–106` captures an error instead of committing a successful response on failure.

Public basis remains the inherited explicit measured-failure rule, `book/edition-3/chapter-01.md:349–353`. `TestChapter2MeasuredEmptyReplies` (`chapter2_revision_test.go:77–141`) uses real public roots, agents and local HTTP servers across all three vendors: success followed by two measured failures; exact cumulative accounting; one request per Ask; unchanged committed text history; only successful-response usage in the reducer; correct completed/error event counts. The reviewer's prior OpenAI/Gemini regressions pass. A scratch R2 revert fails measured totals for both affected vendors; restore passes. Independent local CLI checks additionally require nonzero exit, `input=11 output=2`, and exactly message/request/error events for all three vendors. No estimated usage, fake completed response, retry, or raw failed-body archive was introduced.

### R3 — multiple supported OpenAI opaque fields: closed

`internal/engine/render.go:295–309` retains the existing exact-provenance gate and assigns each supported field only when that field is present. A later opaque part no longer clears the other field with nil. There is no new wire API or native Responses support.

Public basis remains `book/edition-3/chapter-02.md:939–947`, carrying opaque material back to its issuing model. `TestChapter2OpenAIMultipleOpaqueParts` (`chapter2_revision_test.go:143–220`) exercises Parse → durable facts → Dump/Load → Context → Render, both part orders, both fields, vendor/model/surface exclusion independently, and unchanged stored facts. The unchanged external reviewer regression passes. A scratch R3 revert reproduces the lost field; restore passes. Local CLI exact-model and other-model render checks pass.

## Independent execution and sensitivity

All reviewer processes use an allowlisted environment with no provider credentials or credential-file reads; dependency fetching is disabled. Exact argv, cwd, watchdogs, exits, durations and output are retained in the new evidence directory. Actual grader targets are always the new pinned executable, not a guessed project directory.

| Check | Actual result |
|---|---|
| Explicit `go build -o $SCRATCH/student-cli ./cmd/ch01` | exit 0; byte-identical to submitted binary |
| `go test -race -count=1 -timeout=180s -coverpkg=./... -coverprofile=… ./...` at student root | exit 0, **126.672 s**; includes the real 62-second header/body source probes |
| Three exact new regression tests, root package, `-race -count=1` | exit 0, **1.272 s** |
| `go vet ./...` | exit 0 |
| Previous external public-import probe, unchanged test bytes, real public root | all tests pass, exit 0, **1.367 s** |
| Explicitly built `grade-edition2 -ch 2 -json $NEW_CLI` | **100/100, `passed:true`, exit 0, 0.417 s** |
| Explicitly built `grade-edition3 -ch 1 -json $NEW_CLI` | **100/100, `passed:true`, exit 0, 124.056 s**; actual 62-second header and body phases |
| Three isolated fix-revert controls | **R1, R2, R3 each behavioral exit 1 when reverted, exit 0 after restoration**; no compilation failures counted as sensitivity |
| Affected capture replay / targeted keyless CLI / local error-cost checks | script exit 0; detailed counts below |

Both graders used **the same binary**, with its SHA-256 checked before/after each run. The independent source build is also byte-identical. The 210-second outer source-suite and 180-second grader watchdogs did not expire; no shortened substitute probes. Test output confirms the focused selectors actually ran the named tests.

The previous external probe was copied unchanged from the initial reviewer evidence into a scratch consumer module importing only public `ensemble`; it was not weakened to accommodate the fixes. Scratch sensitivity reverts are the copied round-one `check-sensitivity.py`, inspected before execution. Actual records are `R1-sensitivity.txt`, `R2-sensitivity.txt`, `R3-sensitivity.txt` and `sensitivity-summary.json`; source hashes match the pin after restoration. These are the **three new fix-revert controls**, not the initial five unrelated controls.

**Coverage:** independently unioning duplicate instrumentation blocks produces **749/838 statements = 89.379% (89.4% displayed), 89 uncovered**, matching the report. The denominator is root-module production code including CLI, excluding test source, the separate GUI module, external grader/CLI processes and historical live calls. Raw profile and per-file counts are preserved. Current rules impose no numeric source-coverage threshold; no percentage is used as a substitute for the targeted checks or source review.

The production delta is small and ownership-preserving: no new interface, constructor, route, transport, global state, or framework. The required separate GUI import stub and inherited architecture remain unchanged from the initial review. Existing source tests and both graders continue to pass the inherited Ref/stub, fail-closed, ID, provenance/signature, Anthropic ordering, replay/state, isolation and cancellation checks. Prior source-review conclusions stand except for the three defects now corrected; this is not a claim of exhaustive coverage or a security proof. The new comments explain chronology and attempt accounting rather than hiding behavior behind abstraction.

## Retained live evidence and affected offline checks

**No new paid/live provider calls were made by the student in this round or by this reviewer.** The six retained live calls belong to initial commit `74eb026…`, as the new report explicitly states. The original launcher exit 1 caused by the final verifier's wrong regex remains preserved; it is not relabeled as a successful launcher run. The three original live CLI exits were 0, and the earlier offline verifier/evidence interpretation remains separately qualified.

Retention is appropriate for this delta: model routing, HTTP endpoints, auth and transport are unchanged; fixes affect chronological replay, already-decoded error accounting and retention of already-supported opaque fields. Independent execution of an inspected copy of `check-offline.py` on the **new binary** verifies:

- **3** dumps of original captured logs, byte-equal to their source files;
- **18** new-binary renders (two per origin/target pair) and **9** old-binary comparison renders, equal for all nine captured-log pairs;
- **6** affected keyless renders: one R1 sequential-name fixture, R3 exact/other-model fixtures, and three redaction renders;
- **zero** HTTP requests throughout those offline operations;
- **3** local fake CLI requests, one per vendor, each expected **exit 1** with measured cost and error-only capture; exact paths recorded;
- **39 total CLI executions: 36 exit 0, three expected local-error exit 1**. Original evidence files unchanged.

These counts are asserted against `offline-checks-summary.json` in `offline-counts.json`. They prove affected offline behavior and local protocol/accounting behavior, **not new-binary live-provider validation**. The original ordinary-text calls also do not prove live signed tool calls, nonzero live caching, media or streaming. No unnecessary paid repetition is required for these local fixes.

## Preserved failures, limits, and next step

Expected negative outcomes are retained, not hidden: three reverted-fix test failures and three local CLI error exits. No source test, grader or public regression failed on the unchanged new candidate. One read-only inspection typo requested nonexistent root `agent.go`; the correct `internal/agent/agent.go` section was recovered, with the typo recorded in `inspection-note.txt`. No failed setup was counted as a candidate bug or a passed test.

Not repeated this round: initial five mutation controls, unchanged GUI module build/test, private first-edition comparison, unrelated paid live calls, or a wholesale initial audit. The earlier validated evidence remains intact. No GUI runtime, credential settings, native Responses, new feature, manuscript change, grader rewrite, install, restart, push, commit or tag was attempted.

**Next step:** parent verifies the saved final review/evidence and candidate pin, then may perform the normal Chapter 2 acceptance/tag/phase handoff. No additional student correction round is requested. This decision is limited to the exact candidate and contract above; future source changes require a fresh pin and appropriately scoped checks.
