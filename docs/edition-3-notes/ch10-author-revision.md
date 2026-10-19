# Chapter 10: sole permitted author revision

**Completed for root's direct verification. No second text review or text
acceptance is claimed.** This revision addresses only formal-review R1/R2.

## Corrections

**R1 — failed-to-pass history.** Replaced the unsupported minimal-running-config
and declared-model causal account with the documented default-load fixture
mismatch: the generated `CURRENT_RUNTIME_PROMPT` had a terminal newline while
the unchanged exact expectation did not. The chapter now says that correcting
the fixture input resolved the mismatch, that no student production defect was
established, that no prompt trimming was requested, and that production source
stayed unchanged. The final Chapter 10 **100/100**, same-frozen-CLI Chapter 1
success, delayed-header/body phases, and runtime-configuration ownership remain.
The original **80/100** failure remains separately stated.

**R2 — unsupported first person.** Replaced both reviewed first-person
assertions with third-person technical prose. Replay verification still explains
persistence's early placement, and old/new card identity still has the same
ownership requirement without prescribing a new public identity API.

There are exactly three passage replacements. Reversing them reproduces the
immutable reviewed chapter byte-for-byte. Commands, chapter policy, surrounding
story and all other prose are unchanged. See `changes.json` and
`reviewed-to-final.diff` in the new revision evidence directory.

## Reads and checks

The revision task, complete formal text review, complete student final
reconciliation, and current chapter were read before writing. Seventeen other
mandatory full reads were carried from the same uncompacted context only after
verifying their hashes unchanged. The read ledger records all identities.

The final `go run ./cmd/lintprose -v book/edition-3/chapter-10.md` exits **0**, with
**no hard failures** and **six negation forms**. Removing unsupported first
person raises the person-gap warning from **1,297** to **2,469** words against
the 1,200 threshold. That warning is retained and disclosed as expressly
permitted; no invented attribution or unrelated prose was added to reduce it.
The exact argv, cwd, minimal environment, timing and output are preserved in the
new lint receipts. The chapter-only `git diff --check` exits **0**.

No implementation source was opened or changed in this revision. No application,
grader, browser or provider execution, experiments, children or extra planning
stage occurred. All commands finished. Disk remained above the 1 GiB floor.

## Immutable history and seal

Before and after editing, the original **24-member author inventory** was
verified with its `book/edition-3/chapter-10.md` member resolved through
`docs/edition-3-notes/ch10-text-review-input-chapter.md`. The original author
seal is unchanged and **does not validate the revised current chapter at its
original path**. Both original seals, all **seven reviewer inventory members**,
the immutable input chapter and its manifest remain unchanged (35 preservation
records, with overlapping inventory paths recorded separately).

The old author reports retain their historical, now-superseded account; they
were not rewritten. This report and the new chapter carry the R1 correction.
The new `ch10-author-revision-evidence/seal.json` inventories the revised chapter,
both new reports and all new revision evidence, excluding itself. Root owns the
direct R1/R2 condition check and any subsequent text-acceptance decision.
