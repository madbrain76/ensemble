# Chapter 6 — independent correction-round-1 re-review

## Verdict

**ACCEPTED — the demonstrated Gemini boundary/provenance defect is closed on this exact R1 candidate. No further code correction is required by this re-review.** This is the independent code-review verdict, not repository acceptance, a tag, promotion of protected solution material, or authorization to start the next student. Root must inspect this report, the saved final tail and index, then join before making that decision.

This continues the same independent reviewer/context and the original whole review. It is not an entry/adapter re-audit or a fresh authoring pass. One formal correction round has been used. No student work was requested or replayed by this reviewer.

## Exact authority and inputs

| Input | Verified identity |
|---|---|
| Authoritative re-review task | `/tmp/cadence-ch06-r1-rereview-task.txt`, SHA256 `23a9cf0928045f83a4b2f750f00f47fdec5aa33488348a1589a5c324f9d4be5b` |
| Candidate receipt | `docs/edition-3-notes/ch06-candidate-r1.json`, SHA256 `d883d22f2708fae768932a65c3720dee417562a216ec39f628e4563cae9f5181` |
| Corrected student | clean `f595c787675e2a6343053ec42ea6c7e1e02a930e`, tree `6fb36168020a22b8993412812987338aa4b6e010`, **99 tracked files** |
| Original reviewed candidate | `800da21d215b9ca48de41516be4773962fe38ccd`, tree `d54d3cce7e708e393d1b5576bbf3867c037847f4`, 98 files |
| Exact frozen R1 core | `/tmp/e3-ch06-correction-r1/core`, SHA256 `06c26681699a1cc67bfb1eb28280bd4c323e22923395e7714aefeb47edd5a0d4` |
| Exact frozen R1 shape | `/tmp/e3-ch06-correction-r1/shape`, SHA256 `d52f16f4174e72559b0bce9cd72d74878e38d5003a1f44ec2166e7c2b0df8451` |

The two immutable source archives match their Git trees. All **310 R1 artifact entries** and all **474 original entries** were re-hashed successfully, including private-file observations by hash only. The source/evidence manifests and receipt, not an inferred path or a freshly rebuilt replacement executable, supplied the executable authority. Both frozen binaries' build information identifies the corrected commit with `vcs.modified=false`.

The approved adapter and historical grader pins, current Chapter 6 manuscript, and all prior entry/adapter/whole-review records remain unchanged. Protected prior evidence exact sets are 98, 86 and 184 files respectively; their complete checksum indexes verify. Final details are in `integrity.json`, `final-state.json`, `core-build-info.txt`, `shape-build-info.txt` and `prerequisite-hashes.json` in the evidence directory below.

## Prerequisites and review scope

The full current AGENTS, coding skill, Edition 3 reviewer skill, execution/readiness/carryover/public student guidance, historical workflow/carryover, voice guide, chapter procedure, Chapters 2–3 provenance/event contracts and complete Chapter 6 were read. Combined-output omissions were recovered, including a standalone full Chapter 2 read. The original whole review and qualified probes, the complete R1 student report and command ledger, both R1 manifests, exact two-file diff and complete corrected parser/new tests were examined. After the GUI interruption the full active instructions were refreshed and every prerequisite hash was checked unchanged; the same saved independent context was retained.

This re-review closes the original whole review's sole blocking finding. Unchanged whole-program obligations remain supported by that preserved source review and evidence, fresh root/module gates and selected integrations below. No claim is made that a numerical grade exhausts the public contract or that every possible provider frame is certified.

## Delta and architecture assessment

Only two files differ from the original candidate:

- `internal/engine/stream.go`: **20 added / 6 removed lines**, inside the existing `gemini` accumulator method; the rest of this file is byte-identical.
- `chapter6_gemini_boundary_test.go`: **305 added lines** of public-facade regression tests.

There is no parser API, global-ID scheme, compatibility layer, model row/default, controller, actor/history owner, scheduler, transport lifetime or application feature added. The shared plain decoder, renderer, source gating, CLI and external modules are unchanged. The repair uses the existing response-local accumulator and shared canonical decoder rather than creating a second response representation.

The ownership rule at `internal/engine/stream.go:374–415` is narrow and readable: explicit array positions delimit parts; only a first entry may continue compatible unsigned content; a signature closes its binding. A pure trailing signature may finish an unsigned text/thought/tool part, but may not overwrite a signed part or absorb a later explicit array entry. Standalone opaque material receives its own position and local ID and emits no opaque-data delta. Local variables/comments explain the boundary rather than introducing framework scaffolding. The size is proportionate to the fix; most added code is the regression proof. No speculative later-chapter work or test weakening was found.

## Closure of the original finding

| Original failure | Corrected evidence |
|---|---|
| Signature-only opaque part first was rejected | Accepted as its own canonical opaque part, with source and final ID; no private-data delta. |
| Signed ordinary text absorbed later unsigned text | Distinct ordered text parts; the original signature remains bound only to its signed text. |
| Two signed texts collapsed and lost/overwrote a binding | Both parts and both exact synthetic signature bindings survive. |
| Signed thought lost its signature to a following opaque marker | Thought and opaque marker remain distinct, followed by the correct ordinary text. |

The four original qualified independent probes were copied **byte-for-byte**, not weakened. They pass against R1. An additional independent public-only module supplies hand-written canonical expectations, not merely agreement between two parsers: four boundary cases in one-frame and split-frame forms, legitimate text/thought/tool signature continuation, explicit adjacent unsigned text and explicit text followed by opaque material — **13 defined layouts**. This is a finite fixture set, not a numeric coverage threshold.

For those layouts the checks cover ordered parts/types/native data and signature bindings; exact From/vendor/surface/resolved-model facts; measured input/cache/output usage; one response event; exactly-once finals; local zero IDs and parser reuse; delta/final association and content; no signature bytes as deltas; same-source native replay; cross-model/vendor/surface private-data filtering; and retained ordinary text. The fresh root race run also executes the student's real-HTTP/public-Ask privacy regressions. Ordinary answers do not acquire opaque/private thinking data.

### Independent causal/sensitivity control

The same frozen independent test source was compiled against a scratch R1 tree, then against that tree with **only** `stream.go` restored to the pinned original implementation, then with the corrected file restored. No student source or frozen application binary was changed.

| Control | Compile exit | Behavioral exit | Meaning |
|---|---:|---:|---|
| Corrected R1, full independent set | 0 | 0 | Original probes, expanded canonical/replay matrix and retained integration controls pass. |
| Original implementation, same boundary tests | 0 | 1 | Original four failures and expanded boundary failures reproduced; legitimate continuity controls remain meaningful. |
| Corrected bytes restored, full independent set | 0 | 0 | All 99 control-source files match R1 again; assertions are unchanged. |

The original production file SHA256 is `f070beb523fc0546110b213227d5b5ec163ccf694aecd91d247cbfda78a3309f`; corrected file SHA256 is `37f860cbfb84296153f5d408c583897a87a29452e6420758f0e09937ed29fed0`. The compiled original control differs from the corrected executable. The restored test executable is **byte-identical** to the corrected positive: `7b8b78953b1f71e1053ee456e3a2db35feed975b3c3cb96000659e8dd0dca7d2`. Exact source/module/test/binary hashes, original-control diff and command exits are retained in `control-results.json`, `original-restoration-control.diff` and the `independent-*` receipts. Independent-module vet exited 0. No broad mutation/provider campaign was run.

## Fresh independent gates on pinned R1

All processes ran sequentially with bounded Go concurrency, read-only module mode, local toolchain, `GOPROXY=off`, `GOSUMDB=off`, an allowlisted environment without inherited credentials, and owned TMPDIR. No timeout wrapper shortened a lifetime gate.

| Command receipt | Exit | Seconds |
|---|---:|---:|
| `root-race` | 0 | 143.441 |
| `root-vet` | 0 | 0.608 |
| `ch04-consumer-race` | 0 | 0.301 |
| `ch04-consumer-vet` | 0 | 0.111 |
| `ch06-race` | 0 | 1.434 |
| `ch06-vet` | 0 | 0.118 |
| `ch07-race` | 0 | 0.130 |
| `ch07-vet` | 0 | 0.110 |
| `grade-ch06` | 0 | 13.483 |
| `grade-ch01` | 0 | 124.059 |

The external module race/test commands report no test files where applicable: those are compile/link gates, not invented unit coverage; vet is separate. Their actual module behavior is exercised by the cumulative grade and exact shape integration.

Chapter 6: **passed=true, 100/100, exit 0**, with all six existing checks passing. Chapter 1: **passed=true, 100/100, exit 0**, on the **same frozen core**, including both genuine **62-second delayed-header and delayed-body phases**, completed naturally. The score reports retain their historical Edition 2 Chapter 6 title; the actual command was the approved `grade-edition3 -ch 6` adapter with separate explicit core/consumer/workflow/exercise identities, declared `claude-sonnet-4-6`, exact R1 media map and frozen executable. No scores or grader checks were changed.

`commands.json` and every `*.meta.json` retain exact cwd, argv, environment, timestamps, duration, exits, disk observations and output hashes; adjacent stdout/stderr files are complete. Full grade outputs remain in `grade-ch06.stdout` and `grade-ch01.stdout`.

## Pertinent integration and inherited-contract coverage

- **Exact R1 Gemini shape:** both streaming and plain modes exited 0 and returned the same ordinary answer and **three finalized parts**, with two visible text deltas and known-zero usage. Streaming first delta was **4.217 ms** versus **161.018 ms** observed request duration. This is local fixture timing, not provider latency. The inspected fixture's assertions were unchanged; only its output root was moved into owned scratch. The executable was the exact frozen R1 shape, not a reviewer rewrite.
- **Exact R1 CLI arrival/admission:** the retained gated fixture observed the explanation at **12.273 ms**, hint acknowledgment at **12.294 ms**, and only then released tool data at **12.304 ms**. It verified next-request-only hint delivery, correct local tool effect and API trace activity without the synthetic reflected key. The core exited 0. Only the executable pin and owned output path changed from the qualified prior fixture; assertions were not relaxed.
- **Stricter final-answer check remains failed:** that same scripted CLI fixture deliberately does not echo the hint in its final answer after another request. `hint_in_answer=false`, checker exit **1**, is retained and **not credited as a pass**. It does not contradict the public next-request-only pending-context contract and is not a repaired live-model compliance claim. All other public fixture checks passed.
- **Fresh retained public controls:** sparse/interleaved OpenAI tool correlation; indexed Anthropic opaque parity; SSE framing/malformed rejection; real request-mode flags for all three adapters; unknown-model plain fallback; partial-error history exclusion; late-frame generation fencing; and streaming multi-call interruption/accounting all pass in the corrected and restored full independent runs.
- **Carry-forward, explicitly not freshly replayed:** the earlier whole review's ten-test streaming Chapter 5 wake-up overlay remains valid evidence for unchanged actor/history/Anthropic paths. Fresh root race tests cover Chapter 5 behavior; fresh streaming interruption and actual CLI hint tests supplement it. No claim is made that the entire earlier overlay was rerun in R1.

Together with byte-preserved source owners and the fresh lifetime gates, these checks retain the one-Parse seam, response-scoped IDs (zero valid), source-bound opaque replay/signed ordinary text, honest usage, actual Content-Type parsing, request-owned HTTP lifetime, reliable completion, clean history ownership and Chapter 5 wake-up/interruption contracts. Public flags remain `EN_DISABLE_STREAMING` and `CH07_NO_STREAM`. No withdrawn shorthand or invented callback API is required.

## Assistance, failures and live-evidence boundary

The student received the bounded four-case finding/contract, not reviewer solution code. Its retained setup failures are real: nonexistent `Context.Turns`, missing positive MaxTokens, and a refused wrong source anchor followed by misleadingly named old-code failure receipts. The actual applied repair and successful final gates are separately identified. No additional student correction was requested here.

This reviewer likewise preserves uncredited setup/bookkeeping mistakes: inferred nonexistent informational paths; a copied read-only recorder; incorrectly transcribed hash literals; differing manifest size fields; a wrong informational metadata key; the first CLI pin guard failure before target launch; and the first finalizer's assumed internal receiver/receipt schema. Authority checks were corrected against pinned receipts, not against desired test results. The initial combined integration command remains exit 1; `integration-results.json` explicitly reconciles its completed shape result with the later qualified CLI receipt rather than claiming that initial sequence passed.

**No paid call, real key/settings read, install, GUI action, child agent, cache cleanup, commit, tag or push was performed by this reviewer.** All **13 real live exchanges carried in the R1 manifest belong to the original attempt**. They are preserved historical observations, not fresh R1 live passes and not proof of changed Gemini bytes against every real provider. The original Anthropic final-answer failure remains visible. No actual provider rejection or new provider behavior is asserted from synthetic fixtures. Private payloads were not imported into the report or used as new student context.

## Interruption recovery and final preservation

Bill's GUI restart interrupted this reviewer after the behavioral commands and successful `final-integrity-r2` had completed, before this report/checksum index existed. On resume the same context/task/ID was retained, all 28 pre-interruption receipt output hashes were verified, no relevant process was running, and no behavioral gate was repeated. The restart backup was not read or imported. The later hash-only `recovery-final-integrity` also exited 0.

The source remains clean at the exact corrected commit/tree. Both candidate artifact inventories, frozen binaries, grader pins, original manifests/reports and all prior review checksum sets still verify. Scratch `api.log`/`cr/io` test byproducts are retained and separately hashed, not counted as candidate source. Unrelated concurrent publication/manuscript work was preserved and excluded from candidate attribution.

Minimum recorded free disk during the measured receipts was **951,214,080 bytes**; the later restart-era increase has unknown cause and is not attributed to cleanup. All work remained sequential and disk-conscious.

## Final handoff

Evidence: **`docs/edition-3-notes/ch06-code-review-r1-evidence/INDEX.md`** and **`SHA256SUMS`**. `SETUP-AND-FAILURES.md`, `commands.json`, `control-results.json`, `integration-results.json`, `recovery-state.json`, `integrity.json` and `final-state.json` identify the exact proofs and qualifications. The final checksum index covers the complete new evidence set; report and index hashes accompany the submission.

**Final independent verdict: ACCEPTED for this exact R1 Chapter 6 code candidate.** No further implementation change or public-contract clarification is required by this review. Acceptance/promotion/tagging and progression remain root decisions after full saved-tail inspection and join. No jobs, tests, graders, provider requests, children or questions are pending; this independent context remains reserved for later code-review work.
