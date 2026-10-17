# Evidence index: one independent formal text review

Verdict and scope: [job-wakeup-amendment-review.md](../job-wakeup-amendment-review.md). This directory contains text-review receipts only, not implementation or runtime proof.

## Required read coverage

Read completely, with oversized/capped output recovered using bounded ranges:

- `AGENTS.md` (382 lines), including its current Edition 3 overlay and complete remainder.
- `.agents/skills/edition-3-text-reviewer/SKILL.md` (12 lines).
- `.agents/skills/edition-3-author/SKILL.md` (12 lines).
- `docs/edition-3-execution.md` (53 lines).
- `docs/agentic-codebook-workflow.md` (2,434 lines; final length check caught and recovered lines 2417–2434).
- `docs/edition-2-carryover.md` (1,110 lines).
- `book/voice.md` (632 lines).
- `book/chapter-writing-procedure.md` (622 lines).
- Entire `book/edition-3/chapter-04.md` (1,052 lines), unchanged across the two amendment freezes.
- Entire initial Chapter 5 and entire replacement `book/edition-3/chapter-05.md` (1,218 final lines). Replacement reread ranges were 1–310, 311–640, 641–960, and 961–EOF.
- Entire initial and refreshed `docs/edition-3-notes/job-wakeup-amendment-author.md`.
- Complete initial and refreshed author `amendment.diff`, plus complete `queued-tail-clarification.diff`.
- Complete original `ch04-student-review.md`, initial `ch04-code-review.md`, and correction `ch04-student-review-r1.md` under `docs/edition-3-notes/`.
- Complete original `ch05-student-review.md`, initial `ch05-code-review.md`, and correction `ch05-student-review-r1.md` under `docs/edition-3-notes/`.
- Author inventory, baseline/final lint receipts, whitespace receipt, and refreshed `text-validation.json`.

Historical reports were context for existing contracts and evidence limits, not implementation review or proof of new behavior. No private conversation logs, paid calls, code edits, grader runs, new live campaign, or children were used.

## Reviewer receipts

- [`frozen-pins.sha256`](frozen-pins.sha256): exact final verdict targets, explicitly confirmed by root in `main-17`.
- [`root-pin-confirmation.txt`](root-pin-confirmation.txt): confirmation provenance and scope.
- [`inventory-verification.txt`](inventory-verification.txt): independent `shasum -a 256 -c` output, all 28 entries OK, exit 0.
- [`verified-amendment.diff`](verified-amendment.diff): independently generated `diff -u` over the author's captured baseline/final chapters. All hunks byte-match author `amendment.diff` after removing timestamp/path headers.
- [`text-verification.json`](text-verification.json): final all-true checks for five pins, inventory, current/final snapshot identity, diff hunks, line counts, preserved openings, historical sections, and both prefaces equal to HEAD.
- [`guidance-pins.sha256`](guidance-pins.sha256): SHA-256 of required guidance read in this review.
- [`receipt-attempt-note.txt`](receipt-attempt-note.txt): both failed local receipt attempts and their resolution.
- [`failed-expected-pins.sha256`](failed-expected-pins.sha256) and [`failed-text-verification.json`](failed-text-verification.json): retained WRONG reviewer expectations and resulting failed checks. These are **not** the verdict pins. Root confirmed the actual intended freeze before they were corrected; no manuscript changes occurred.

## Reproducible bounded checks

From the repository root:

```sh
shasum -a 256 -c docs/edition-3-notes/job-wakeup-amendment-review-evidence/frozen-pins.sha256
shasum -a 256 -c docs/edition-3-notes/job-wakeup-amendment-evidence/authored-files.sha256
```

The independent receipt script compared current chapter bytes to the author `final/` copies; recomputed unified diffs against `baseline/`; counted actual added/deleted lines; compared chapter opening/historical sections; and used read-only `git --no-pager show HEAD:<preface>` for the two preserved prefaces. It did not infer authorship from the dirty worktree. Production, grader, guidance, reference and old acceptance paths were not written.

## Author evidence retained separately

All links below are author-provided and remain unmodified:

- [Author report](../job-wakeup-amendment-author.md).
- [Full author inventory](../job-wakeup-amendment-evidence/authored-files.sha256).
- [Full baseline-to-final diff](../job-wakeup-amendment-evidence/amendment.diff).
- [Queued-tail-only clarification](../job-wakeup-amendment-evidence/queued-tail-clarification.diff).
- [Baseline lint](../job-wakeup-amendment-evidence/baseline-prose-lint.txt) and [final lint](../job-wakeup-amendment-evidence/final-prose-lint.txt): inspected receipts, not new reviewer lint runs.
- [Text checks](../job-wakeup-amendment-evidence/text-validation.json) and [whitespace receipt](../job-wakeup-amendment-evidence/whitespace.txt).
- [Superseded initial freeze](../job-wakeup-amendment-evidence/superseded/initial-draft/): original pins/evidence preserved rather than rewritten as the final target.

## Limit of the verdict

PASS means the specified amendment is clear, bounded, consistent with the authorized contract, and honest about evidence. It does not mean runtime fixes exist, that old grades prove new behavior, that the parallel Chapter 5 R1–R3 implementation review is complete, or that any later student extension can omit its real public-interface regressions. No further author revision is requested by this review.
