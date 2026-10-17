# Edition 3 Chapter 6 — independent whole-candidate code review

Date: 2026-10-17. Role: independent reviewer, not student or author. This is the whole-candidate code review, **not** another entry/adapter audit.

## Verdict: REQUIRED CHANGES — do not accept or tag this candidate

The frozen candidate passes the existing Chapter 6 grader, Chapter 1 lifetime grader, and its root/module gates. It nevertheless violates the published Gemini signed/opaque-part parity contract. One cohesive blocking correction is required: preserve Gemini opaque/signed boundaries instead of rejecting or merging them into neighboring parts.

This is a measured public-contract failure, not a numeric coverage threshold or a request for a broader architecture. The other bounded controls below passed. No student/application/grader/manuscript/guidance code was changed by this reviewer. Zero formal student correction rounds had been used at review entry; this report does not itself launch one.

## 1. Authority and preserved identities

- Assignment: `/tmp/cadence-ch06-code-review-task.txt`, SHA256 `0f7c1f4d0b802949e18cdb6b327e8233187de641f90e25116f80da804950b926`.
- Candidate receipt: `docs/edition-3-notes/ch06-candidate.json`, SHA256 `55f613021fc62b332ec935d6e59e01fd95d6a34bfca3272d5d4539de815432c3`.
- Student commit: `800da21d215b9ca48de41516be4773962fe38ccd`; tree `d54d3cce7e708e393d1b5576bbf3867c037847f4`; clean worktree; 98 exact pinned source files.
- Frozen core: `/tmp/e3-ch06-frozen-core`, SHA256 `b2fd9a630d5587d810c1069e44f63d9dd64daebe402f99893cb760b601043db5`.
- Exact grading binary: `/tmp/cadence-ch06-grader`, SHA256 `8818400c89bb1663cad72faeca99bc92d64fb9a26e9f023a54d5e319c0e2a798`.
- Adapter production source remains SHA256 `4227ad0a0d96f17b541b60f85e7ac19acd09edfbb79e50c35153aab58f7c5ee7`; tests remain `6f443849bef87fcf2b5e97566bff545088e9919db15212450560d76f11c9159a`.
- Historical Chapter 6 grader remains `edfd0ef2ad9e9841bf531ba552a8983ac182fa8b1f61612596343aa102d080b7`; public chapter remains `ed007e7e7059400939a6735534d7a953d4575e83230a6f68fc35583e8a17bd55`.
- All 474 candidate artifact entries and the complete source-package artifact list were verified twice. Frozen shape/module/package identities are retained in `input-verification.json` and the original verified source/evidence manifests; no private provider payload was copied into this report.
- Final integrity also verifies the exact unchanged prior entry-evidence 98-file set and adapter-review 86-file set, their checksum indexes, and both prior reports.

All new evidence is under `docs/edition-3-notes/ch06-code-review-evidence/` (`E` below), with owned scratch identified by `E/scratch-path.txt`. Main's pre-existing/concurrent publication and manuscript changes were preserved, not attributed to this candidate, and not reverted or committed. No pinned authority input drift occurred.

## 2. Full prerequisite and whole-source coverage

The following were read in full during this turn; omitted combined-output intervals were recovered:

| Prerequisite group | Coverage |
|---|---|
| Repository/reviewer instructions | `AGENTS.md`; `.agents/skills/ensemble-coding/SKILL.md`; `.agents/skills/edition-3-code-reviewer/SKILL.md` |
| Current Edition 3 procedure | `docs/edition-3-execution.md`, `edition-3-readiness.md`, `edition-3-carryover.md`, public `edition-3-student.md` |
| Required historical lessons/craft | `docs/agentic-codebook-workflow.md`; `docs/edition-2-carryover.md`; `book/voice.md`; `book/chapter-writing-procedure.md` |
| Public contracts | Complete Edition 3 Chapters 4, 5, and 6, including accepted wake-up/batch behavior |
| Current package/history | Complete assignment, candidate receipt, entry audit, adapter review, student review and `COMMANDS.md`; complete source/artifact manifests parsed and every listed artifact hash/size checked |

Current student-first procedure supersedes historical author-first instructions. No runtime skill registration was invented. No broad reading or implementation of later chapters was undertaken.

All **24 production Go files** were reviewed, including unchanged inherited history, jobs, tool dispatch/process ownership, public common types, public façade, CLI/display, external modules, and the non-running GUI placeholder. `production-source-review.json` records each path, line count, and hash. The diff and inherited test migrations were inspected; root/module tests were run fresh. Production Go delta is 1,147 additions/98 removals across 15 files; test Go delta is 666 additions/30 removals. These are scope facts, not acceptance metrics. The larger 67-file overall diff includes evidence/documentation rather than 67 production changes.

## 3. Blocking finding R1: Gemini loses opaque/signed part boundaries

**Location:** candidate `internal/engine/stream.go:359–419`, especially `385–400`.

A part without text/function-call content is either attached as a signature to `a.last`, or rejected when no previous part exists. Adjacent same-kind text/thinking blocks are merged regardless of an already attached signature, and a later signature overwrites the earlier one.

The published contract requires the streamed path to produce the same ordered canonical parts and opaque provenance as the plain response. It explicitly requires standalone Gemini thought/signature/opaque spans and source-bound signed ordinary text. A trailing signature continuation is supported, but that does not authorize treating every signature-only payload as continuation or extending a signed part across a new boundary.

The external public-API test `TestGeminiSignedAndOpaqueBoundaries` pairs equivalent plain and streamed responses using only reviewer-created fixture sentinels:

| Case | Observed streamed result versus successful plain parse |
|---|---|
| Standalone opaque signature first, then visible text | `unsupported Gemini stream part`; the plain parser accepts the opaque part |
| Signed ordinary text followed by unsigned ordinary text | Two parts become one; later text inherits the earlier signed part's provenance |
| Two separately signed ordinary text parts | Two parts become one; only the later signature remains |
| Signed thought, standalone opaque boundary, then visible text | Three parts become two; the intervening opaque signature replaces the earlier thought signature |

**Qualified reproduction:** `independent-public-final-r2.meta.json` / `.stdout` / `.stderr`; exact test source in `public_checks_test.go`. `go test -race -count=1 -v ./...` exits **1**. All four Gemini subcases fail. The other seven top-level independent tests pass; there is no race diagnostic. Harness `go vet ./...` exits **0**.

The tests compare parsed canonical data, not arbitrary JSON key ordering. The standalone-opaque-first case also fails directly, independently of adjacent-text interpretation. The issue can change source-bound replay payloads, although no real-provider rejection or downstream effect is claimed: this review made no paid calls.

**Required correction boundary:** repair this Gemini boundary/provenance handling and add regressions for these paired cases while retaining legitimate trailing-signature continuation, response-local ID correlation, and ordinary plain/SSE behavior. Do not add another parser API, global parser IDs, an actor/controller rewrite, or invented model capability rows. No implementation patch or answer code is supplied here. The existing manuscript is sufficiently explicit; no public clarification is necessary to establish this obligation.

## 4. Fresh required gates on the exact candidate

Every command, working directory, exit, duration, and full stdout/stderr is retained in its `*.meta.json`, `*.stdout`, and `*.stderr` triplet. `commands.json` is the combined index. Commands were sequential, with credential-free environment allowlisting, `GOMAXPROCS=2`, and `GOFLAGS=-p=1`.

| Gate | Exact outcome |
|---|---|
| Root `go test -race -count=1 ./...` on exact extracted source | Exit **0**, 144.703 s |
| Root `go vet ./...` | Exit **0** |
| `examples/ch04-consumer` race/vet | Both exit **0**; no module-local test files where Go reports that fact |
| `examples/ch06` race/vet | Both exit **0** |
| `examples/ch07` race/vet | Both exit **0**; no module-local tests claimed |
| Exact adapter, Chapter 6, pinned frozen core, distinct student core/consumer/workflow/stream-exercise roots, explicit `claude-sonnet-4-6`, declared case/media inputs | **100/100, passed=true, exit 0**, 13.150 s |
| Same exact core binary, Chapter 1 | **100/100, passed=true, exit 0**, 124.062 s |

Chapter 1 completed both actual **1m2s delayed headers and 1m2s delayed body** checks. No shorter timeout, substitute fixture, cancellation, or reused prior score is represented as a fresh result. Source review additionally checked caller-owned request context and body closure; the numeric lifetime controls alone do not prove arbitrary duration.

The Chapter 6 grader itself says complete opaque replay, gated flush, and SSE framing remain source/integration obligations. Its 100/100 therefore does **not** contradict R1 or certify all Gemini boundary shapes.

## 5. Additional bounded behavioral evidence

| Obligation | Fresh evidence and result |
|---|---|
| OpenAI sparse/interleaved tool indices, argument accumulation, zero/local IDs, exact final correlation | `TestOpenAISparseInterleavedToolsAndUsage`: pass; provider-index ordering retained; parser reuse does not leak prior state |
| Anthropic indexed opaque/text blocks and usage/cache parity | `TestAnthropicIndexedOpaqueParity`: pass; sparse IDs are legal, not forced into compact/global numbering |
| SSE comments/unknown fields, multiline `data:` newline join, CRLF, frame reset, EOF flush, malformed/incomplete rejection | `TestSSEFramingAndErrors` plus reviewed reader and root stream tests: pass; not a claim about every possible SSE/network shape |
| Wire-based dispatch and whole-part callbacks for JSON responses | Real HTTP mode tests return JSON even for streaming-capable requests; plain fallback yields one visible delta/final and correct text |
| Request flags and Gemini path/query selection | All three local provider adapters, enabled/disabled modes: pass; OpenAI usage option checked |
| Unlisted-model plain fallback | `TestRealRequestModesAndUnknownFallback/unlisted_model`: pass, no streaming request, one visible delta/final |
| Partial-response error and durable-history correctness | `TestPartialErrorNeverCommitsResponse`: pass; error does not persist the partial response; subsequent response succeeds |
| Interrupted stream and late old frames | `TestLateFramesAfterInterruptAreFenced`: pass; fresh request succeeds, no old final/late display or partial durable response |
| Explicit multi-call interruption under streaming | `TestStreamingMultiCallInterruptPreservesPairs`: pass; first retained job stays running, later tool does not execute, both original calls have exactly one result, no fabricated success, new request completes |
| Actual external Chapter 6 shape exercise | Same frozen shape binary: stream and plain both exit **0**, identical fixture answer; stream 2 text deltas versus plain 1; measured explicit-zero usage retained |

The shape fixture recorded first streamed delta at 4.223 ms versus 166.627 ms total. This is bounded local evidence, not a latency requirement or a claim about a live provider.

### Actual CLI gated arrival, input admission, and trace sanitization

`cli-gated-control.py` uses only a loopback, explicitly scripted Anthropic-shaped server and the exact frozen core. The server sends the first explanatory text, then withholds tool data until the controller has seen that text on CLI stdout and the hint acknowledgement. Results in `gated-cli.json`:

- CLI text at **12.714 ms**, hint acknowledgement at **12.738 ms**, server release at **12.746 ms**: actual arrival/flush proof, not post-hoc final ordering.
- All three requests request streaming; the hint is absent from request 1, present in request 2, absent from request 3.
- Tool file effect succeeds; final prose is printed once; the candidate exits **0**.
- Actual `api.log` frame tracing is present. A reflected fixture API-key sentinel is absent from trace/stderr; no header/query/settings/key content was harvested.

The script deliberately uses literal answers that do not echo the hint. It exits **1** on the retained stronger `hint_in_answer` assertion after all public-behavior checks pass. That nonzero result is preserved and is **not** counted as a passing model-answer exercise or a candidate contract failure.

### Accepted Chapter 5 behavior under streaming

A scratch **fixture-only** overlay emits SSE for the original Chapter 5 tests; production Go bytes are unchanged, and the complete `chapter5_wakeup_test.go` is byte-identical to the pinned candidate. No assertions were weakened. All 10 selected race tests pass (`streaming-ch05-wakeup-and-interruption-r2`, exit **0**): mixed batches with signed replay, setter/dispatch boundary, later process wait, wait timer/available output, stale timer fencing, completion races, failed retained reports, addressed process ownership, pending-context/reliable queue, and retained-job interruption/shutdown. `wakeup-selection.json`, the wrapper source, and final integrity prove the selection and unchanged source.

## 6. Source/architecture assessment beyond the finding

- The single public `Parse(response, callbacks)` seam remains the translation boundary. There is no `Parser.Stream` service or globally unique parser-ID requirement. `Response` identity scopes actor/display correlation; zero part IDs remain valid.
- HTTP reading remains worker work. The actor owns generation/response reduction, turns, tool dispatch, complete-response commit, pending context, and reliable request completion. Deltas do not become durable conversation events. Public observers remain best-effort and nonblocking by contract.
- The CLI's buffered-input correction is in the actual CLI path; the reviewer tested prompt/hint admission while a real response body was held open. Presentation consumes deltas/finals without introducing a second agent loop. `EN_DISABLE_STREAMING` and the exercise's `CH07_NO_STREAM` are the actual public controls.
- Complete final parts are committed before final observations. Original tool batches are retained, skipped/unstarted calls receive error results rather than fake success, and interrupted retained jobs remain root-owned. Fresh streaming controls cover these inherited Chapter 5 obligations.
- Renderer provenance gates distinguish vendor/surface/model, and unsafe opaque cross-source replay is dropped rather than improvised. Existing signed ordinary-text and opaque tests remain relevant; the mixed-batch overlay's initial wrong-model failure also demonstrates this gate. **Gemini reconstruction before that gate is where R1 breaks provenance units.**
- Measured usage, explicit zero, unknown-versus-zero rejection, and cache components were reviewed across decoder/accumulator/render paths. Root tests and paired local checks passed the exercised cases. No claim is made that an arbitrary unsupported stream shape produces usable usage.
- New state is localized to provider accumulation, response correlation, and presentation. No extra scheduler, transport, retry framework, GUI implementation, fallback codec guessing, or unrelated feature platform was introduced. No broad refactor or simplification is required before addressing R1.
- Inherited test migrations adapt `Render`/`Parse` to the new public seam; their behavioral assertions were inspected, and the full root race run passed. Passing migrated tests are not substituted for the new independent boundary checks.

## 7. Small fault-sensitivity set

Exactly two deliberate implementation faults were tested in isolated scratch copies. They are not candidate defects and were never applied to student/main source.

| Fault | Compiled changed-byte result | Targeted behavior | Restored control |
|---|---|---|---|
| Ignore `Config.DisableStreaming` | Compile exit **0**, distinct source hash | Real HTTP enabled/disabled check exits **1**, three disabled-provider cases fail | Exact original file and complete scratch tree restored; rebuild **0**, behavior **0** |
| Offset streamed final part IDs | Compile exit **0**, distinct source hash | Exact delta/final correlation check exits **1** | Exact original file and complete scratch tree restored; rebuild **0**, behavior **0** |

The shared baseline compiled and passed first. `mutation-results.json` retains all original/mutant/restored source hashes and all five test-binary hashes; `.diff` files and full command receipts retain the changed bytes and diagnostics. No mutation campaign, provider sweep, grader changes, or replacement scoring was performed.

## 8. Prior assistance/failure history and reviewer setup failures

The assignment/student ledger were not sanitized into a success-only story. The student initially scored 70 and repaired the actual CLI buffered-input admission path. The retained Anthropic live attempt failed a final-answer hint-compliance assertion; later reported live passes are existing student evidence, not reviewer-run paid observations. Prior root adapter setup/test-isolation failures remain historical evidence; this turn neither re-ran nor rewrote that accepted adapter audit. No formal student correction round was consumed by this reviewer.

The original live checker tests `hint_in_answer` over output; the portable failure summary says the hint reached the next request but was not repeated in a later answer. Source plus the gated three-request control establishes the public **next-request-only** boundary. Thus the failed stronger assertion is preserved as an experiment/checker limit, not repaired by extending hint lifetime or inventing a live answer. The separate frozen-shape local pass is not substituted for that failure.

Reviewer failures are fully retained and uncredited (`SETUP-AND-FAILURES.md`):

- Initial receipt verifier assumed the wrong schema; corrected before extraction/tests.
- Review module initially needed `go mod tidy`; no behavior ran then.
- Early independent assertions wrongly expected OpenAI arrival order, compact Anthropic IDs, or omitted required fixture usage/system facts. Initial sources/outputs remain; qualified tests use the published ordering/ID contract and complete fixtures.
- A guessed non-existent engine convenience method and later a mistaken Boolean interpretation of `Jobs.Get` caused compile/vet failures. Correct public API calls were used in the qualified runs; no behavior was credited under those failures.
- First streaming wake-up overlay hardcoded the wrong reported model. Correct source-bound replay rejection caused its signed-block failure. The wrapper was corrected to preserve the fixture's original model, **without changing production or test assertions**; the full selected suite then passed. Preliminary speculation about JSON ordering was withdrawn after direct inspection.
- First gated-CLI trace assertion expected a label the logger does not use. The actual sanitized frame record is checked in the qualified rerun; the final-answer assertion still intentionally fails.
- Initial final-integrity check incorrectly expected a runtime-populated scratch directory to contain only source files. All 98 source hashes already matched. The 295 local `api.log`/`cr/io/*` byproducts are retained and indexed separately; no cleanup hid them.
- Informational filename/glob/read mistakes are listed as setup-only. They did not launch grading with guessed inputs or alter protected files.

## 9. Limits, carry-forward, and freeze

- Required next step is one bounded Gemini correction and regression evidence, followed by the ordinary independent review of that exact corrected candidate. This report does not accept/tag the source or authorize student launch by itself.
- No new paid/live call, settings/key read, GUI, child agent, cache deletion, main/student commit, tag, push, or application/manuscript/grader/guidance edit occurred. Existing private provider artifacts were hash-verified without publishing their contents.
- Root race/module gates and two exact-binary graders are fresh. Existing reported live cases are retained provenance only. Local scripted fixtures are labeled as such. No arbitrary coverage percentage or extra model behavior was made a public requirement.
- This is bounded evidence, not proof of every provider revision, network failure shape, or scheduling interleaving. The contract-specific Gemini failures are sufficient to block acceptance despite the positive controls.
- `E/INDEX.md` links the receipts, source/provenance ledgers, failed attempts, and mutation controls. `E/SHA256SUMS` covers the complete frozen evidence set; final report/index hashes are supplied with submission.
- Final integrity verifies clean student/source/package pins, unchanged prior audit evidence, unchanged current prerequisites, unchanged overlay production/assertions, and fully restored mutant trees. Unrelated main changes are preserved and unattributed.
- All grader/test/controller jobs are finished; no children were created and **no jobs are pending**. Owned scratch and saved review outputs are frozen. Root must inspect the final saved tail and join before accepting this review.
