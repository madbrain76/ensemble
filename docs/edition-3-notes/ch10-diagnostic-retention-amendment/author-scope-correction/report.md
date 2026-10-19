# Chapter 10 — pre-student scope correction

**DRAFT_PUBLIC_AMENDMENT_PENDING_INCREMENTAL_STUDENT_VALIDATION**

Root withdrew its provisional import of Chapter 14's positive-N canonical
`log_retention` policy. This is a small root-owned scope correction before
incremental student work, **not** a formal text-review revision; none has occurred.

Chapter 10 now requires only two new policy settings: `diagnostic_max_records`
and `log_mode`. Full mode preserves **all available canonical history**. Removed
the third-setting clauses, positive-N example and associated Chapter 10 policy
obligation. Diagnostic 200/default, explicit zero/unlimited, negative validation,
mode validation, derived dialogue archive, old-format compatibility, required
recovery tail and save-failure ordering remain. Existing truncated-input,
snapshot/tail, apply-once and sequence obligations remain; no new positive-N
canonical-retention implementation or tests are requested for Chapter 10.

Before editing, preserved the initial draft at `../author-initial-draft.md`:
`e52a8ab1e33737c82a3d972ae158f263f12b86c2f07a6664dfe2a2872aab9675`. The old `../author/` artifacts and inventory are byte-unchanged.
Their chapter SHA and lint pins now resolve to this immutable alias, **not** the
corrected current chapter. The old report's third-setting and positive-N Chapter
10 requirements are superseded by this correction, not silently rewritten.
The original `../chapter-before.md` accepted baseline remains unchanged.

## Future Chapter 14 obligation — unimplemented

Chapter 14 retains its existing `log_retention` policy and must align it with
full/dialogue-only modes. At entry, root must carry forward the successful
`micro_handoff` save hook to the same diagnostic-retention operation, preserving
snapshot-before-journal-trim ordering and the complete recovery tail. Test that
hook's configured diagnostic count (including zero), multiline whole records,
continued writer appends and save-failure non-destruction; verify restart,
sequence/apply-once behavior and exclusion of tool/opaque parts from dialogue.
No Chapter 14 edits, hook implementation or tests occurred here.

## Evidence and limits

Current chapter SHA-256: `1c2ac0b080e05b53e34152f5669dfc1aff5157c0e06f76d8b4d48fe8079d9992`.
`changes.diff` is relative to the immutable initial draft, not the accepted
baseline. Prose decreased by 683 bytes; no code fences changed.
All prior read hashes matched before editing; `verification.json` carries their
original coverage, including the explicitly partial Chapter 8 read. No new
prerequisite read loop. All 34 original protected pins still match.

Bounded offline same-basename prose lint exited **0**. Soft warnings: negation
forms **12/9**, person gap **2724/1200** (initial draft: 13/9 and 2758/1200).
Exact invocation/hash/output: `lint-receipt.json`, `lint.txt`; initial-draft lint
remains in the old sealed evidence. `git diff --check` passed. `inventory.json`
seals the new evidence and excludes itself.

Only Chapter 10, the new immutable alias and this new evidence directory were
written. No runtime tests, implementation gates, student/runtime/grader/guides
edits, providers, browser, children or commit. Student remains paused. Root
inspection and incremental student validation precede the one formal text review.
