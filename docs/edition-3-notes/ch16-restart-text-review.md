# Chapter 16 restart — formal text review

## Verdict: PASS

Exactly one formal post-student text review. No required revision and no author response requested. Final state: **HOLD**. This is not authorization for another review, live OAuth work, publication, installation, or advancement to another chapter.

## Scope and basis

Read the full current reviewer skill and author overlay, `AGENTS.md`, `docs/agentic-codebook-workflow.md`, `docs/edition-2-carryover.md`, `book/voice.md`, `book/chapter-writing-procedure.md`, and `docs/edition-3-execution.md`; recovered clipped ranges. The current lean guide and this assignment govern over dated campaign, mandatory-three-round, and author-first instructions.

Read the complete current `book/edition-3/chapter-16.md` and the Chapter 16 restart student, code-review, and author notes. Compared the manuscript with `/private/tmp/ch16-author-baseline/chapter-16.md`. Inspected the relevant root entry changes against HEAD (`cmd/grade-edition3`, `internal/gradeedition2/ch16*`, including the added contract tests), the GUI Hub change, and its deterministic R3 test. This is a text review of the accepted code outcome, not another code-review or validation round.

## Findings

- **The author delta is appropriately narrow.** It identifies the earlier student evidence, explains the summary catch-up repair, and records the restart results and live-work deferral. The baseline comparison shows no rewrite of the opening, TL;DR, established rationale, or historical failure account. The preface was not edited by this review.

- **The summary wording is acceptable in context** (`chapter-16.md:640–654`). “Text supplied only by the final summary” occurs within the GUI observation/live-delta discussion and describes text the catch-up path obtains from the final summary card, beyond the prefix already observed there. It need not mean that the provider emitted no earlier summary deltas. The source distinction matters: `solutions/edition-3/main/gui/hub.go` previously caught up unseen final suffixes for ordinary text only; R3 includes `reasoning_summary` under the same identity, prefix comparison, cursor, and ordering rules. Skipped or coalesced partial checkpoints can therefore require this recovery even when the provider did emit deltas. `gui/chapter16_r3_test.go` controls immutable checkpoints for both entirely unseen and prefix-already-seen summaries, checking missing-byte recovery, ordering, identity, repeat-checkpoint idempotence, reconnect finals, and unchanged canonical History. That test establishes these GUI cases, not provider-only-final delivery as the cause of the failed grade. The chapter neither asserts that transport cause nor promises that every final establishes incremental delivery; it explicitly says the opposite at line 654. A wording preference alone does not justify an author correction.

- **The root entry supports the documented local exercise.** Chapter 16 routes an explicit GUI target to the existing seven-check grader and rejects unrelated parity/cache flags. The driver supplies an isolated synthetic issuer/provider environment, uses the actual plain-line prompt contract, and observes the GUI's error frames. The summary check in `internal/gradeedition2/ch16.go:73–89` measures separate complete `part_delta` content and the causal gate; a surviving final card alone does not satisfy it. The author has not recast this local grade as live-account or native-UI proof.

- **Results and limitations remain accurately separated** (`chapter-16.md:874–933`). The Chapter 16 third attempt is the accepted 100/100 result; 15/100 and 85/100 remain failed attempts. The first Chapter 1 regression grade is recorded separately. The carried handoff reports 100/true/exit 0 for each accepted grade and real 62+62-second waits for Chapter 1. The historical paid subscription run still records its completed task, zero public-summary deltas, launcher kill, and lack of graceful-shutdown/saved-conversation proof. The new restart section quotes Bill's live-OAuth deferral and expressly disclaims fresh live account/provider, native GUI, browser, and listening observations.

## Evidence limits and preserved work

The accepted code-review handoff records three production rounds, subsequent test-only timing corrections, and the 216-file source export at `solutions/edition-3/ch16`. Local root/GUI race and vet results are carried; inherited module and JavaScript coverage retains the limitations stated in that handoff. No additional source acceptance or coverage claim is made here.

Reused the coordinator's same-basename baseline/final lint receipts rather than rerunning lint: both exit 0, no hard failures, preexisting soft warnings retained (negation 13→14, allowance 11→12, person gap 3937→3940). Those warning changes do not establish a blocking defect in this bounded delta.

This review ran no tests, grades, lint, network/provider requests, native/browser work, or live authentication. It read/copied no credentials, changed no settings, production code, manuscript, preface, or historical receipts, created no children or audit bundle, and made no commit. Disk checks remained above the 1 GiB floor. The sole written file is this review note.

**Required revisions: none. PASS; HOLD.**
