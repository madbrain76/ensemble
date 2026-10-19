# Chapter 10 public retention amendment — author draft

**DRAFT_PUBLIC_AMENDMENT_PENDING_INCREMENTAL_STUDENT_VALIDATION**

Chapter SHA-256: `e52a8ab1e33737c82a3d972ae158f263f12b86c2f07a6664dfe2a2872aab9675`  
Immutable before SHA-256: `c156708aae6fde042cb34340da239279702041b245e096efb03a5d19a7178aca`

## Scope and direction

Changed only `book/edition-3/chapter-10.md` plus new evidence in this `author/`
directory. The first prose-only turn stopped before edits. This resumed draft
follows the expanded direction and Bill's later clarification: **compacted log
means dialogue-only**, not arbitrary last-N canonical events. It is new student
work, not a claim that accepted Chapter 10 already implements these policies.
Old acceptance/history and `edition-3-ch10-r1` were not rewritten.

The draft adds `diagnostic_max_records` (200 default, zero unlimited, negative
invalid), startup and successful-save boundaries, whole-record writer-owned
cleanup, and failure preservation. A startup with no valid recovery save must
establish one successfully before destructive cleanup. Diagnostic policy never
trims canonical saves, journals, memories or worker state.

Saved history defaults to `log_mode=full`, `log_retention=0`. Positive canonical
retention keeps the required post-anchor tail even when it exceeds N. Explicit
`dialogue_only` mode uses optional `log_mode`/`dialogue_log` save fields and an
honestly derived `source_seq`/`role`/`text` archive, not rewritten Events. Mixed
assistant records contribute visible text only; tool calls/results and opaque
or private replay parts are excluded. Snapshot Context stays authoritative;
complete canonical recovery events stay separately in `log`, and a later
journal must remain unfiltered. Old four-field full-mode saves remain valid.

This intentionally chooses a text dialogue projection, not a new replay or
multimedia archive. It neither compacts Context nor caps dialogue archive size.
Chapter 14's existing `log_retention` is retained for canonical history, not
repurposed into the dialogue filter. Returning to full mode cannot recreate
removed originals. Incomplete full-history rebuilds must be rejected; available
checkpoint diagnosis and tail replay do not reconstruct discarded earlier steps.

Bill's Gemini fork disk-exhaustion report is attributed in one sentence; storage
internals and snapshot presence remain unverified. His reported Git push of
memory logs is not verified backup evidence and authorizes no deletion.

## Required Chapter 14 entry obligation — not implemented

Root must carry this obligation explicitly into the Chapter 14 entry contract
before its student work, aligning TL;DR persistence/settings rules and §§14.12–14.13:

- After `micro_handoff` completes a successful recovery save, call the **same**
  configured diagnostic-retention operation. Keep snapshot-before-journal-trim
  ordering; never trim the required post-`as_of` recovery tail or let a filtered
  dialogue archive stand in for the journal.
- Align saved-history mode/projection with Chapter 14's canonical `log_retention`;
  preserve snapshot state, sequence allocation and apply-once replay in both modes.
- Test successful handoff with more than 200 multiline diagnostic records,
  a positive override and explicit zero; assert whole-record order and safe
  continued appends. Inject save failure and prove diagnostics/recovery history
  are unchanged. Verify restart, required-tail preservation and no duplicate
  application, plus exclusion of mixed tool/opaque parts from dialogue projection.

No Chapter 14 text/runtime changed and no handoff test was run. Chapter 10's
incremental validation must first cover its settings persistence/validation,
startup and normal-save success/failure, live-writer safety, both archive modes,
complete-log equivalence, bounded-tail restart and unavailable-replay refusal.
Root inspects this draft, then commissions incremental student work and one
formal text review **after** that experience. No extra formal review occurred.

## Checks and limits

Full prerequisite reads are recorded in `read-ledger.json`; only voice and
execution guidance were carried from the prior full read, with exact unchanged
input hashes. Full Chapters 10 and 14 were read; Chapter 8 ownership excerpts
are labelled partial. No implementation, grader or private-audit content was
inspected; protected-file checks were hashing only.

Baseline and final same-basename `go run ./cmd/lintprose -v` both exited **0**,
offline, bounded to 120 seconds each. Baseline person-gap warning: 2469/1200.
Draft warnings: person gap 2758/1200; negation forms 13/9. These are soft prose
warnings, not hidden passes or reasons for unrelated rewrites. Exact commands,
input hashes and output are in the lint receipts. `git diff --check` passed.

All 34 original protected tracked pins and the immutable before-copy match.
The opener, exercise/grading commands and table, all fenced blocks except the
necessary SaveFile schema extension, and §§10.6–10.7 including history are
byte-preserved. Chapter 14 has no worktree diff. No runtime tests, implementation
gates, provider/browser/network activity, children, commit, or student work;
student remains paused. Runtime behavior and new acceptance remain unverified.

Evidence: `changes.diff`, `read-ledger.json`, `lint-baseline.txt`,
`lint-draft.txt`, `lint-receipt.json`, `current-chapter.sha256`,
`preservation.json`, and self-excluded `inventory.json`.
