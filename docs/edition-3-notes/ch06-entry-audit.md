# Edition 3 Chapter 6 — independent entry audit

## Decision and scope

**The existing Chapter 6 grading API is usable with separated artifact identities. The measured artifact-conflation problem requires a thin root CLI adapter, not a grader rewrite.** This is an entry qualification of that composition, not Chapter 6 student acceptance, and not certification that the historical reference satisfies every current public requirement.

No Chapter 6 student has been launched. The accepted Chapter 5 worktree was immutable throughout. No application, student, manuscript, grader, guidance, protected solution, historical source, tag or commit was changed. Only this report, `ch06-entry-evidence/`, and the owned scratch directory were written. No children, GUI, paid/live provider calls, credential reads, broad cleanup, or provider/mutation campaign were used. Tests used local fakes and a credential-excluding environment. Normal Go build-cache writes and grader-owned temporary fixture lifecycle were not redirected into source trees.

Main input: `b67386ab80cf2afe7985db5a32ab99ec2311f090` (tracked diff empty before/after). Protected `solutions/edition-3` tree: `4ff207e66df6d49cbdea462df2f6753d421779d8` unchanged. Pre-existing untracked repository files were neither adopted nor removed.

Student input: `f60ebf3aad2ece4251face9d95c106dd667abcd6`, tree `3de5d0d24dd9f593ec017fceccfbff1b209c5d57`, dereferenced tag `edition-3-ch05-r1`; clean and byte-manifest unchanged. The already accepted race-built CLI at `/tmp/cadence-ch05-wakeup-evidence/candidate-r2/frozen/ensemble` was reused after verifying SHA-256 `e51138885a359eaf801258a52bc21f28ca41691683256647f8c9911be410b233`. Its existing provenance is in `ch05-wakeup-student-evidence/final-source-verification.json`. No unchanged Chapter 1 actual-62-second control was separately rerun for readiness.

## Prerequisites and precedence

The following were read completely, with oversized combined outputs recovered by bounded interval reads. `ch06-entry-evidence/inputs.json` gives exact byte counts, line counts and SHA-256 identities; final verification confirms these files stayed unchanged.

| Prerequisite | Read coverage / application |
|---|---|
| `AGENTS.md` | Full current file; instruction hierarchy, read-only review, retained user changes, test/evidence discipline. |
| `.agents/skills/ensemble-coding/SKILL.md` | Full file, including cumulative architecture and offline/fake-test rules. Read directly as repository guidance; no invented runtime skill registration. |
| `.agents/skills/edition-3-code-reviewer/SKILL.md` | Full file; independent code-review role and required historical prerequisites. |
| `docs/edition-3-execution.md` | Full file; current **student-first** procedure controls. |
| `docs/edition-3-readiness.md` | Full file; no Chapter 6 student readiness inferred from a score alone. |
| `docs/edition-3-carryover.md` | Full file; cumulative architecture and boundary obligations. |
| `docs/edition-3-student.md` | Full public brief; no private answer source is supplied to a student. |
| `docs/agentic-codebook-workflow.md` | Full historical workflow, including appendix/execution records; retired author-first ordering is historical context, not current authority. |
| `docs/edition-2-carryover.md` | Full historical carryover, including later/deferred topics. |
| `book/voice.md` | Full voice guide. No author/manuscript work undertaken. |
| `book/chapter-writing-procedure.md` | Full procedure; current Edition 3 procedure supersedes old author-first staging. |
| `book/edition-3/chapter-04.md` | Full chapter: public facade/star topology, common ownership, request limits/lifetime, signed/source-bound replay, media and external consumer. |
| `book/edition-3/chapter-05.md` | Full chapter: actor ownership, causal feedback, responsive observers, wake-up versus termination, parked batch suffix, final callbacks and batch correlation. |
| `book/edition-3/chapter-06.md` | Full chapter and exercise, including SSE framing, callback/final semantics, CLI flushing and required source/integration proof. |

Reviewed grading sources include both root CLIs, `internal/gradeedition2/ch06.go`, `ch06_stream.go`, `ch06_test.go`, inherited Chapters 3–5 composition/topology/workflow/media code, and the original `internal/grade/ch07*` evaluator/fixture implementation. Matching historical code was consulted privately for provenance and the two selected scratch defects, not as a student answer.

## Public contract, not withdrawn shorthand

The published seam is **one `Parse(resp *http.Response, cb StreamCallbacks) error`** for SSE and ordinary JSON. Plain mode supplies one delta per part, then the same semantic finalization. There is no public requirement for a separate `Parser.Stream` entry point.

Parser part IDs are response-scoped; zero is valid. Delta/final IDs must correlate within that response; actor log IDs may be mapped afterward. `ch06_stream.go` already keys comparison aliases by `(response, id)`, accepts explicit zero, rejects missing/negative/noninteger IDs, and excludes tool-result finals using existing event identity. A global parser-ID mandate would be incorrect.

The public CLI switch is `EN_DISABLE_STREAMING`; the separate streaming exercise uses `CH07_NO_STREAM`. `ENSEMBLE_DISABLE` is not the contract. Known models use their actual stream capability; unknown models remain usable in plain mode without a made-up capability entry.

## Smallest valid grading composition

Current `cmd/grade-edition3` supports only Chapters 1, 4 and 5. It already separates positional build target from explicit `-core`, `-consumer`, `-workflow`, `-consumer-case`, and `-media-models` inputs. That shape supports the accepted root public-library plus nested `cmd/ch01` layout.

The required call already exists unchanged:

`gradeedition2.Chapter6(bin, core, consumer, workflow, exercise, model, example, media)`

It invokes corrected Chapter 5, which invokes corrected Chapter 4 and inherited earlier checks, then adds five streaming checks and a 20-point `ch5-parity` gate. Do **not** run a second separate Chapter 5/4 grader just to reconstruct the same cumulative result. A separate Chapter 1 readiness rerun is not needed for this entry task.

The evidence-only caller is `ch06-entry-evidence/api/main.go`. It reads the existing case/media formats, calls the existing API, emits its report, and exits 1 if any check fails. It contains no replacement grading logic. The exact pre-gofmt build source is retained as `api-source-initial.txt`; `verify-harness-format.*` confirms it is gofmt-identical to the saved formatted source. All measured runs used the same harness binary hash below.

For the accepted layout, keep these distinct:

- CLI build target: `<student>/cmd/ch01`, or an independently verified frozen executable.
- Core source identity: `<student>` (its `go.mod`), **not** the binary or `cmd/ch01` directory.
- Chapter 4 external consumer: `<student>/examples/ch04-consumer`, with its declared case JSON.
- Chapter 5 workflow: `<student>/examples/ch06`, with the accepted media mapping.
- Chapter 6 streaming exercise: the student's separately declared public-consumer module, once one exists; do not invent a directory or reuse the actor exercise as if it were submitted.
- Streaming model: a declared known model; retain the separate unknown-model control.

### Proposed adapter only — not implemented here

Extend the existing root command's Chapter 4/5 path to Chapter 6, accept explicit streaming-exercise and known-stream-model inputs (e.g. historical `-stream-exercise` and `-stream-model` names), validate the supplied module roots, and dispatch to the existing API. Retain positional executable/build target, `-core`, existing case/media inputs, JSON/text reports, cleanup and 0/pass, 1/failed-check, 2/invocation-error semantics. Do not infer core from target, impose a main package at module root, hard-code private historical subdirectories, duplicate the inherited graders, or alter scores/fixtures to accommodate a CLI bug. Adapter source and focused CLI routing/validation tests remain a separate root-owned change/review task.

This recommendation concerns **artifact plumbing only**. The source/integration obligations below remain necessary after an adapter yields 100/100.

## Measured results

Historical input is `HEAD:solutions/ch06`, tree `5351f4519f9a4b8f97ff8c6d23547c9746db5d03`, exported by pinned `git archive` into `/tmp/ensemble-ch06-entry-7_7zbtqs/solutions/ch06`. Final verification compares every source file to the pinned Git bytes. Consumer/workflow/stream modules are this snapshot's `ch05`, `ch06`, `ch07`, respectively; the known model used is `claude-sonnet-4-6`.

Each evidence stem below has full `.stdout`, `.stderr`, and `.meta.json`: exact argv, cwd, UTC start, exit, elapsed seconds, disk before/after, and output hashes. `commands.txt` consolidates them. All graders ran sequentially with `GOMAXPROCS=2`, `GOFLAGS=-p=1`; the runner imposes no timeout. Short tool waits were observations, never cancellation.

| Evidence stem | Score | Passed | Exit | Seconds | Interpretation |
|---|---:|---|---:|---:|---|
| `positive-api` | 100/100 | true | 0 | 13.198 | Explicit binary/core identities; pinned historical positive. |
| `positive-old-cli-source-valid` | 100/100 | true | 0 | 12.597 | Historical CLI with module-root positional source target. |
| `positive-old-cli-binary-valid` | 70/100 | false | 1 | 8.498 | Same frozen historical executable used by API, but old CLI conflates target and source identity. |
| `ch05-nonsubmission-baseline` | 20/100 | false | 1 | 163.086 | Accepted Chapter 5 predecessor; **not a Chapter 6 submission**. `ch5-parity` passes. |
| `mutant-disable` | 90/100 | false | 1 | 13.001 | Real disable/request-mode defect; `delivery-not-content` fails. |
| `mutant-correlation` | 50/100 | false | 1 | 12.696 | Real final-ID defect; text/thinking/tool correlation checks fail. |
| `restored-disable` | 100/100 | true | 0 | 12.757 | Source bytes restored, rebuilt, all checks pass. |
| `restored-correlation` | 100/100 | true | 0 | 12.532 | Source bytes restored, rebuilt, all checks pass. |

All initial, mutant, and restored builds exited 0. No pinned-reference fixture failure was concealed or repaired.

### Artifact-conflation diagnosis

The old command passes its positional target both to `grade.Build` and as `core`. With the source-root target this happens to work. With the **identical API-positive frozen binary** it loses `delivery-not-content` (10) and `ch5-parity` (20). Safe actual diagnostics include:

- `stream exercise dependency: core and consumer need distinct module declarations`
- `ch4-parity: star-topology: core and consumer need distinct module declarations`
- `not-deaf: workflow dependency: core and consumer need distinct module declarations`
- `hub-clean: go list in /tmp/ensemble-ch06-entry-7_7zbtqs/positive: fork/exec /usr/local/go/bin/go: not a directory:`

The streaming content checks still pass. This controlled source-root/identical-binary/API comparison demonstrates a launcher identity error, not a streaming implementation defect. The accepted-library baseline separately confirms that explicit `core=<student>` works for its real nested CLI layout: its inherited Chapter 5 parity remains passing.

Two initial old-CLI attempts incorrectly used `-chapter`/`-core` flags. Both exited 2 (`flag provided but not defined: -chapter`), produced no grade, and are retained as `positive-old-cli-source.*` and `positive-old-cli-binary.*`. They are **invocation failures**, not positive or negative grading controls. The corrected commands use `-ch 6` and a positional target. A later portable-format verification replaced a shell process-substitution attempt unsupported by `/bin/sh`; see `invocation-notes.md`. Neither issue changed implementation or grader behavior.

### Explicit predecessor nonsubmission

The baseline used the verified frozen Chapter 5 binary and actual accepted core, consumer and workflow sources, with copied accepted case/media JSON. Its `-exercise` was `/tmp/ensemble-ch06-entry-7_7zbtqs/absent-ch06-stream-exercise`; that path did not exist before or after the run. It was not created, borrowed from a solution, or presented as a failed student module.

The result has `ch5-parity` 20/20 and all five new streaming checks failing. Actual diagnostics include zero streaming deltas, request-mode disagreement, unknown-model plain completion lacking required delta observations, and absent exercise dependency. This shows the new gate distinguishes the predecessor from a Chapter 6 implementation. It says nothing adverse about the accepted Chapter 5 submission. The full 163-second run finished normally using the verified race build; it was not stopped at a short wait boundary.

### Selected sensitivity and restoration

Only two implementation mutations were selected. `mutations.py`, `mutations.json`, the two `.diff` files, complete mutant source manifests, build reports and `restoration.json` retain exact reproduction/identity evidence.

1. **Disable ignored:** the historical request-selection branch in `internal/engine/media.go` was changed so `Config.DisableStreaming` no longer selects plain mode. One file changed, SHA-256 `14fbe62e6acaae317f77e90e49d28ac8ed0a227e3d7370eee66721601c5362ca` → `ed47efc8d6abfca087c634294b96db798ced0f7f44418b2dee591fe1195f2671`. Compile succeeded. Actual failure: `streaming and plain runs both produced 104 deltas; the flag changed nothing`, plain request flag disagreement, and external exercise counts remaining `36/54/14` in disabled mode. Other Chapter 6 checks and inherited parity still pass.
2. **Final correlation broken:** the final callback in `internal/engine/stream.go` received an ID offset from its delta ID. One file changed, SHA-256 `435c213d77d8c1450d594955a0b392faadb03e81de10aac16cfe7fc7c41685c1` → `7a3add6690a5510bed7af923c66196bc42fed74d98f90c6a65ff603c3a4b15c4`. Compile succeeded. Actual failures include `part 5 finalized with text but no delta ever carried part_id 5`, and unmatched thinking/tool finals and delta groups. Counts alone and delivery parity still pass, demonstrating why correlation is a distinct assertion.

Both changed files were then copied back from the pinned positive. Full restored source manifests equal the pinned historical manifest, not merely the touched lines. Both scratch copies were rebuilt into distinct retained binaries and graded back to 100/100. Different scratch source paths cause different binary hashes; equality is asserted for source bytes, not assumed for binaries built at different paths. This is the entire mutation scope, not a numeric coverage campaign.

## Public requirements → existing checks → remaining proof

The six measured check IDs and weights are `stream-deltas` 20, `deltas-match-final` 20, `thinking-streamed` 15, `tool-params-streamed` 15, `delivery-not-content` 10, `ch5-parity` 20. Weights/minimum fixture counts are existing behavior, not an invented coverage threshold.

| Public requirement | Existing measured/reviewed coverage | Necessary remaining source/public/integration evidence |
|---|---|---|
| One `Parse(resp, callbacks)`; unified canonical finalization in SSE/plain modes | CLI stream/plain content observations; source exposes the intended route in the historical input. | Review actual student's interfaces and shared semantic decoder. Black-box equality does not rule out six drifting provider implementations, duplicate event construction, or observer-side parsing. |
| Real text streaming, explanation before tools | Many observed text deltas, exact reconstruction, finalized tool executed and result included in the next request. Inherited checks cover some final event ordering. | Gated stream test must observe explanatory text **while the response is still open and before tool dispatch**, not only sorted/final JSONL records. No current completed-output fixture proves arrival order. |
| Human-facing early CLI flush | No arrival-time assertion: `driveStream` collects output and evaluates after process completion. | Gate the server after an early text chunk; assert bytes available to the human CLI before releasing completion/tool blocks. Verify no duplicate final printing, thinking not mixed into answer, and no extra ordinary-delta wire protocol requirement imposed on the aggregate exercise. |
| SSE framing | Existing fake supplies happy-path frames through real HTTP/parser code. | Deterministic parser tests: arbitrary network chunk splits, multi-`data:` events/newline joining, LF/CRLF/CR handling as published, comments/empty frames, EOF without delimiter, malformed/error events and long frames. Review framing/provider-decoding separation. The fake is not a framing conformance suite. |
| Response-scoped IDs, including zero; exact final correlation | `ch06_stream.go` maps `(response,id)`, validates explicit IDs, rejects duplicate provider finals; `correlateStream` checks both directions and exact text, expected thinking content, tool name and decoded arguments. Selected mismatch mutant fails despite plausible counts. | Student parser/actor tests for zero, ID reuse across responses, mixed kinds, no missing/duplicate finalization, actual canonical event identity/order. Opaque final payload equality is not proved by an observable kind/identity alias. |
| Thinking separate from reply; nontext/opaque parts preserved | Thinking deltas are present, content checked against fake thought, no reply leak window, thinking/final ID correlation. | Signed ordinary text, opaque/redacted thinking, encrypted/opaque items, signature-only suffixes, and final callback/event preservation need actual source-bound replay tests. Presence of a thinking final does not prove payload integrity. |
| Tool name/argument deltas never executable until final | Correlated name/decoded arguments, complete scripted `run_command`, tool result returned in second provider request. | Gate partial argument fragments and assert no early execution. Verify canonical complete-call dispatch, order, error/failure facts and multi-call correlation under streaming. |
| `EN_DISABLE_STREAMING` and `CH07_NO_STREAM`; plain one-delta semantics | Every core request flag checked against selected mode; plain versus streaming counts/content; external on/off shape. Disable mutant fails real request/count assertions. | Source must wire both public switches through actual configuration/request selection. Public/external app should remain an observer, not a shadow transport/parser. Verify per-part/plain callback semantics beyond fixture totals. |
| Unknown-model plain fallback | Separate `unknown-stream-model-edition2` run must return exact answer, make no streaming request and emit two text deltas for the fixture's two text parts. Historical positive passes. | Ensure source preserves unknown-model fallback rather than hard-coding this test model or adding capability metadata merely to satisfy a fixture. Extend only if actual student changes demand it. |
| Stream/plain same final events, tool args and measured usage | Observable finals are deeply compared; text/tool name/args correlation; each run must have two requests and measured totals input 310/output 88/cache-read 0/cache-write 0, not merely equal invented counts. | Full event/provenance/opaque payload equivalence, signed ordinary text and source identity immutability need source/integration assertions; visible final maps do not expose all canonical event fields. |
| Source-bound opaque replay | Corrected inherited Chapter 2/4/5 parity carries earlier log/replay/media checks; streaming continuation exercises a real request. | Test unmodified same-source/same-model replay, wrong vendor/model and edited source-visible text exclusion, signature preservation on ordinary text, plus plain/stream equivalence. Do not infer full replay safety from the reference score. |
| Separate public-library aggregate exercise | `publicDependency` checks distinct module/public root dependency; shape app is built and run, final answer and `stream_stats` required; on/off text/thinking/tool counts compared; `first_delta_ms` must be numeric/nonnegative. | Source-review counts/timing collected from real callbacks, accurate final counts, real first-delta timestamp, no invented numbers, no private/internal dependency or framework implementation copied into the app. Timing presence is not timing accuracy. |
| Inherited request lifetime and cancellation | Chapter 5 → Chapter 4 parity preserved in all positive/mutant controls and accepted predecessor; inherited topology/actor/media checks still run. | Review/run focused Chapter 4 tests around `Complete`-owned context, caller-interrupt independence, `Cancel()` reuse, terminal `Close()`, policy/token/media propagation, and no deferred cancel at headers that cuts off a stream. Callback body lifetime now matters directly. |
| Newly accepted Chapter 5 wake-up and batch interruption under streaming | Baseline inherited parity passes, but those broad checks do not specifically prove the recent blocked-watch/parked-batch behavior during SSE. Accepted `chapter5_wakeup_test.go` and Chapter 5 public contracts were inspected. | Re-run/adapt the accepted focused feedback-arrival, unrelated completion, watched progress, partial-line pattern, kill/remove/shutdown wake-up, and ready-earlier-final/later-blocked-watch batch cases with real streaming where the new code touches them. Assert no polling/new tool call required, no replay/reordering of executed prefix or parked suffix, one live watch/cursor owner, progress versus terminal semantics, complete unique ordered correlation per batch, causal feedback exactly once, responsive sibling/observer/feedback, and race-safe cancellation/cleanup. Do not silently treat `ch5-parity` as this proof. |
| Active agents own mutable state; terminal callbacks/lifecycle | Inherited topology, workflow and observer checks; source-read actor contracts. | Review streaming callback ownership, no cross-actor mutation, no body/callback goroutine outliving request cleanup, and final callbacks on success/error/cancel without duplicate delivery. |

No extra provider sweep, new long control, full framing suite, gated human-CLI test, source-replay suite, or streaming wake-up/batch campaign was run in this finite entry task. They are explicitly remaining acceptance obligations, not measured passes. The historical reference is a matching grader positive, not an oracle for those obligations.

## Public clarification assessment

**No manuscript/public contract correction is demonstrated as necessary by this entry audit.** The single parser seam, response-local ID scope/zero, two disable switches, separate aggregate consumer, arrival-time explanation, gated flush and inherited actor semantics are already stated. Withdrawn handoff shorthand must stay withdrawn. The thin adapter should accept explicit module paths and known-model metadata without adding a private root-main or historical-folder requirement to the public exercise. If later student evidence reveals a real public ambiguity, record that specific gap before changing grading or asking for revisions; no hypothetical clarification is imposed now.

## Evidence index, hashes and completion

Start at `ch06-entry-evidence/INDEX.md`; the machine-readable result index is `results.json`. `SHA256SUMS` covers every frozen evidence file except itself; `REPORT-SHA256` binds this report. Exact inputs and source manifests, initial/final Git state, build information, full commands/exits, stdout/stderr, two mutation diffs, restoration verification and binary hashes are retained. The full original/mutant/restored core binaries remain in the owned scratch path, not in the student tree. The API internally builds short-lived consumer/workflow/exercise executables and cleans them up; their source identities and commands are retained, but this audit did not intercept/hash those internal temporary binaries. Do not label their hashes as independently captured.

Key SHA-256 identities:

- API harness binary: `a8ef301b3b0347b404003943289828aab1cfd15fdbe1ffbf0563b6fe46c02af2`.
- Historical old CLI: `ba3d865edca71e1f5c346c35c628b519e2421de18ad150740d0588e0defafbef`.
- Historical positive core: `c1537b9f56e7a3aac9ff6d84f0d42f1b37c68ff3618b08dbb296be1b07b8535f`.
- Disable mutant core: `ed50a4d253b6bddb3a595140da0566e5b458efceafebca8229d418e81273d345`.
- Correlation mutant core: `63cfc1bc70c0d8c019a3608d5d22d2f363bd0484cc293690cf0d662095226e59`.
- Restored disable core: `f4bac30ec22c18c68c04c5968d06ede8e294d56c67f7e7c25f0cbf4677c94e86`.
- Restored correlation core: `9a73d6886d3bdd6f60af06dfc491041660ac77076602c66f50a9a8bed7b63cc1`.

Recorded grading/build disk availability stayed above 2.4 billion bytes; final verification recorded 2,413,535,232 bytes available (about 2.25 GiB). No disk pressure intervention or cache/tmp cleanup was needed. All grader jobs, build jobs, restoration jobs and the parent-message acknowledgement completed; no children or pending audit jobs remain. Evidence is frozen for root inspection. No author work or adapter implementation is authorized by this report; root must inspect the saved report/artifacts/final tail and join the independent context before accepting the entry diagnosis or making a separate adapter change.
