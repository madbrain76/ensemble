# Chapter 6 post-student text-review evidence

## Verdict and authority

**REQUIRED REVISION — one manual CLI build-recipe correction only.** The report is `docs/edition-3-notes/ch06-text-review.md`. This is the one formal post-student text review. Root alone accepts/closes and may request the one minimal author revision. No second review, code rerun, live proof, or publication authorization is implied.

## Exact inventory rule

`SHA256SUMS` contains repository-root-relative paths. Its exact set is **every regular file recursively inside `docs/edition-3-notes/ch06-text-review-evidence/` except `SHA256SUMS` itself, plus the one external report `docs/edition-3-notes/ch06-text-review.md`**. No other external file is indexed here. The manifest intentionally does not hash itself; its own hash accompanies submission. All indexed payloads are safe UTF-8 text (empty output files allowed), with no NUL bytes, private provider payloads or keys. The two chapter snapshots preserve exact original bytes, including line endings.

The manifest freezes the initial reviewed draft, not a later author revision. No existing author/student/reviewer evidence was rewritten. Re-running `validate_review.py` in place would overwrite its generated result files, so use a separate scratch evidence root if reproducing it after freeze; the saved first-attempt output and exact command receipts are the submitted proof.

## Evidence map

- `baseline/chapter-06.md`, `final/chapter-06.md`: same-basename baseline and frozen reviewed draft.
- `baseline-prose-lint.*`, `final-prose-lint.*`, `lint-commands.txt`: complete diagnostics, exit statuses, exact sequential commands with `GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off`.
- `chapter.diff`, `chapter.git.diff`, `chapter-diff-check.*`, `chapter-git-diff.exit`: narrow diff and whitespace receipts; portable/Git diffs equal the author's frozen copies.
- `pins.json`, `initial-state.json`, `author-manifest-verification.json`, `read-inputs.json`: accepted identities, initial preservation hashes, all 31 author manifest entries/exact set, and fully read current safe inputs.
- `validate_review.py`, `text-validation-attempt1.*`: successful independent text-only checks, with no behavioral execution.
- `local-links.json`, `preservation.json`, `selected-protected-check.json`, `freeze-verification.json`: seven local links, ten concurrent publication hashes and Bill's preface, 23 selected protected inputs (permitted root execution-log append distinguished), final frozen-input reconciliation.
- `chapter-review-excerpts.txt`, `root-clarification.json`: exact finding locations and the answered consequential-ambiguity question. Root confirmed the accepted core-library/nested-CLI package distinction.
- `report-qualification-excerpts.txt`: safe excerpts from the full original/R1 reports, retaining the old-commit/live-failure and GUI-recovery qualifications. The full reports, not the excerpts alone, were read.
- `SETUP-AND-LIMITATIONS.md`, `setup-inventory-probe.txt`: reviewer guessed-filename failures and capped-output recovery, kept separate from successful chapter checks and historical author failures.
- `git-status-before.txt`, `git-status-final.txt`, `git-status-at-freeze.txt`: worktree observations, not attribution of unrelated changes. Root-owned concurrent records and publication work were not edited by this reviewer.
- `completed-jobs.json`: terminal prior command receipts and answered clarification. No children or background provider/test work were launched. Final inventory/readback commands finish before submission.

## Payload paths

The external report and this `INDEX.md` are also indexed. Other payloads:

- `SETUP-AND-LIMITATIONS.md`
- `author-manifest-verification.json`
- `baseline/chapter-06.md`
- `baseline-prose-lint.exit`
- `baseline-prose-lint.txt`
- `chapter-diff-check.exit`
- `chapter-diff-check.txt`
- `chapter-git-diff.exit`
- `chapter-review-excerpts.txt`
- `chapter.diff`
- `chapter.git.diff`
- `completed-jobs.json`
- `final/chapter-06.md`
- `final-prose-lint.exit`
- `final-prose-lint.txt`
- `freeze-verification.json`
- `git-status-at-freeze.txt`
- `git-status-before.txt`
- `git-status-final.txt`
- `initial-state.json`
- `lint-commands.txt`
- `local-links.json`
- `pins.json`
- `preservation.json`
- `read-inputs.json`
- `report-qualification-excerpts.txt`
- `root-clarification.json`
- `selected-protected-check.json`
- `setup-inventory-probe.txt`
- `text-validation-attempt1.exit`
- `text-validation-attempt1.txt`
- `validate_review.py`
