# Chapter 10 retention amendment — formal text review

**Verdict: PASS — no material findings. No author revision requested.**

This is the one formal review of the authorized retention amendment, not a second review or regrade of the historical chapter. Reviewed manuscript: `book/edition-3/chapter-10.md`, 34,877 bytes, SHA-256 `1c2ac0b080e05b53e34152f5669dfc1aff5157c0e06f76d8b4d48fe8079d9992`.

## Findings

- **The two-setting contract is clear and usable.** The opening requirements and §10.5 distinguish `diagnostic_max_records` (default 200, explicit zero unlimited, negative invalid) from `log_mode` (`full` or `dialogue_only`). Runtime policy is separate from archive representation. No last-N canonical-history policy has been introduced.
- **The recovery explanation remains sound.** Lines 417–440 distinguish available canonical history, derived dialogue with source-sequence provenance, the authoritative snapshot, and the complete post-`as_of` canonical recovery tail. Dialogue is not fabricated Event replay; switching modes cannot regenerate discarded originals. Context is expressly not compacted, so the dialogue mode is not presented as secret erasure. Exactly-once tail application and continued sequence identity remain explicit. The envelope and refusal sections retain legacy-full compatibility and atomic rejection of invalid input.
- **Diagnostic safety is taught at the right ownership boundary.** Lines 446–460 count whole records rather than lines, order recovery/save success before cleanup, serialize cleanup with appends, and explain descriptor-offset and old-inode hazards. The prose claims neither continuous disk containment nor general cross-process/power-loss safety. Read-only inspection of the accepted policy, writer, archive, projection and lifecycle paths supports these distinctions; this was not a new code audit or test run.
- **Chapter 14 remains pending.** Lines 463–467 require its future successful-`micro_handoff` save-then-trim integration and policy alignment without claiming implementation or tests here. The pending-entry obligation agrees.
- **Voice and evidence are preserved.** The opening is unchanged from the accepted baseline. The added subsection connects disk exhaustion to the chapter's snapshot/tail reasoning rather than becoming an operational diary. Bill's Gemini-fork experience is attributed testimony, with storage internals and snapshot behavior expressly unverified (lines 402–406). Historical grades, interrupted checks, provider sessions and GUI evidence retain their qualifications. They are not presented as fresh proof of the retention binary. The student's out-of-scope inherited lifetime run is not credited as authorized acceptance evidence.

## Evidence and limits

Read the entire current chapter, baseline, all 21 pinned prerequisite/evidence files, public accepted `CHAPTER10.md`, and the bounded public source listed in `read-ledger.json`. Capped passages were recovered; partial source reads are marked as such. Frozen source HEAD/tree/tag resolve to the root-accepted retention revision; its worktree was clean. The existing acceptance records remain the authority for code acceptance, not this review.

Lint is **carried, not rerun**: exact chapter identity, basename, actual production lint-source identity, prior receipt and output were verified. The prior exit was zero; soft warnings remain 12/9 negation forms and 2724/1200 person gap. They do not establish a material amendment defect or justify stripping useful explanation. See `lint-carry.json`.

No source or existing evidence was edited. All 21 pinned inputs and 34 protected paths were verified unchanged. No children or other assistance; no tests, builds, grades, provider/browser/network exercise, private probes, runtime-agent or paused-Chapter-11 inspection, future implementation inspection, cleanup, publication or runtime installation. Read-only command jobs finished. Available disk space exceeded the 1 GiB floor.

**Root action:** independently verify this exact seal and unchanged inputs, then perform its closure/integration decisions. No reviewer rereview is requested; any sole author revision and closure remain root-owned.
