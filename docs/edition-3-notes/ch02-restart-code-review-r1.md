# Chapter 2 restart R1 — final independent review

## Recommendation

**Recommend acceptance of the held Chapter 2 candidate. No remaining public-contract blocker is established.** The seven correction findings closed within the three correction rounds. Source/design review, targeted regressions, corrected Chapter 2 grade, standalone Chapter 1/lifetime grade, and the authorized live text exercise support acceptance of the same frozen candidate. Root retains the final acceptance decision.

This is the same reviewer context and note, not a new review campaign. Final follow-up read the current student note and actual saved results; it ran nothing. The old rejected attempt and its report remain separate.

## Source and correction judgment

R1 restored the public event/Context vocabulary and star ownership: Agent orchestration, History causal facts/replay, Engine vendor projections/transport accounting, common shared contracts, and thin clients. No material-bloat or fresh-restart objection remains.

The seven closed findings concern: separate live versus historical usage; request-tail ephemera and consumption; Gemini IDs and stable render-only fallback; honest unknown producer identity; unknown usage and measured charges through parser failures; `ErrorOccurred` preserving dialogue; and ordinary signed Gemini text remaining portable without incompatible signatures.

The last correction is concrete in `solutions/edition-3/main/internal/engine/parse.go:113–120,132–166`: invalid cache operands differ from omitted optional zeros, and `if !readOK || !writeOK` / `if !readOK` prevent fabricated derived input charges while preserving independent measured buckets. `chapter02_round3_test.go:24–100` covers 18 dependency cases including live/historical accounting and replay; its supplied log records all passing. Earlier Round 2's corrected-fixture log was a single-subtest rerun, not a full-suite rerun.

## Grader and inherited-test qualification

The original frozen-binary Chapter 2 result remains preserved at `/tmp/ch02-final-validation/grade-ch02.json` and `.stderr`: **85/100, false, exit 1**. Its sole failed check combined two missing explicit replay-target assumptions with rejection of valid unrelated Anthropic text-string content.

I independently accepted the narrow grader repair: `internal/grade/ch02_harness.go:370–389` declares the offline exact target while preserving the live fake override; `internal/gradeedition2/ch02.go:127–168` accepts string content without removing the result IDs, payload, ordering, or text assertions. The focused environment and interleaving positive/negative tests passed. Paired frozen outputs showed unknown targets omitting signatures and exact targets retaining them, plus the correct interleaved result-bearing message.

I ran the corrected Chapter 2 grade **once**, against the unchanged `/tmp/ch02-final-validation/frozen-cli`, under a clean environment without inherited provider configuration. Saved `grade-ch02-corrected-reviewer.{json,stderr,exit}` reports **100/100, passed=true, exit 0**, all checks passing, empty stderr. No student rebuild accompanied this repair.

The inherited test-only updates in `ensemble_test.go` and `cmd/ch01/main_test.go` correctly adapt to causal failed attempts, independent partial measurements, credential-free offline rendering, and four-bucket protocol output. Completed-answer, request-count, no-HTTP-before-validation, usage, diagnostic/secret, and owner-isolation guarantees remain. Original race failures and updated passing output were inspected; they were not hidden or rerun by the reviewer.

## Final results inspected

All paths here are under `/tmp/ch02-final-validation`:

- `grade-ch01.json` and `.exit`: **100/100, passed=true, exit 0**, including completed answers and usage after **62-second response-header and response-body delays**. The same frozen executable was used.
- `root-race-updated.log`: passing updated race run, exit 0. Root/GUI vet and build success are supplied validation results, not reviewer executions. `gui-test.log` explicitly reports **no test files**.
- `live/live.log`: success in 20.05 seconds, exactly **3 discovery GETs and 9 POSTs**, three successful chat turns per provider, no retries/fallbacks. The proxy records independent outgoing-content checks at **1/3/5 messages**, native output cap 1024, actual-answer recall, and four-bucket native/log accounting agreement.

Saved causal logs confirm returned identities match the selected `claude-sonnet-5-5`, `gpt-4.1-2025-04-14`, and `gemini-3.8-flash`. Recalled names are Quillfire-7, FrostPine, and VoxelDrift. Summed input/cache-write/cache-read/output totals are respectively **325/0/0/31**, **222/0/0/9**, and **175/0/0/655**.

The saved keyless dump equals the actual Anthropic log byte-for-byte. All three saved offline requests preserve that conversation; the live log records two identical successful renders per vendor. These are offline cross-vendor projections, not additional provider conversations.

## Limits, not invented blockers

No browser exercise or authored GUI tests; no live tools, media, signature replay, or nonzero-cache exercise. Those are not claimed as covered by the text run. Chat retains the inherited `claude>` label for all vendors and input/output-only status; four-bucket **protocol and causal accounting** were checked independently, not inferred from that display. These cosmetic/display limits do not establish a failure of the chapter's stated contract.

The final results check made no source edits, builds, test/grade reruns, paid requests, manifests, or new report package. Only this existing review note was updated; the failed grader result and all qualification remain visible.

## Root acceptance

Accepted code and required exercises after direct comparison with the original solution and the same independent reviewer’s final recommendation. The star architecture is preserved; additional provenance, usage and replay handling serves current requirements rather than speculative machinery. Three correction rounds were used (one fresh re-execution and two bounded corrections). The source snapshot is `solutions/edition-3/ch02`; Chapter 1 remains unchanged. Prose review and chapter commit follow separately. The live harness demonstrated sequential request quotas, not a generally synchronized concurrent quota guard.
