# Edition 3 Chapter 5 — post-student author report

## Frozen submission

**Result:** two small teaching additions in `book/edition-3/chapter-05.md`; no deletions or rewrites. This is the normal post-student author pass, not another review of the accepted wake-up amendment. The candidate is ready for the parent's one formal text review; this author does not claim text-review acceptance.

- Exact baseline: `/tmp/cadence-ch05-poststudent-baseline/chapter-05.md`, 1,218 lines, SHA-256 `4f4c2595ef65a08511d5cc68553f853feaf854e4cb0fe7d4d9094d4ae7f5a344`.
- Frozen candidate: `book/edition-3/chapter-05.md`, 1,233 lines, SHA-256 `475d807b1fcbd9e3bbb475afc08af7baea6cae438e77d0c404c9030782e7aba9`.
- Exact diff: [`ch05-author-evidence/chapter.diff`](ch05-author-evidence/chapter.diff), 15 added lines, zero removed lines; two prose paragraphs, 134 additional lint-counted words.
- Exact UTF-8 chapter copies: [`baseline/chapter-05.md`](ch05-author-evidence/baseline/chapter-05.md) and [`final/chapter-05.md`](ch05-author-evidence/final/chapter-05.md).
- Explicit evidence hash inventory: [`ch05-author-evidence/SHA256SUMS`](ch05-author-evidence/SHA256SUMS). Paths in this inventory are relative to the evidence directory. It lists every evidence file except itself; the submission supplies the inventory's own hash and this report's hash to avoid circular self-hashes.

The accepted student code identity remains the reported `f60ebf3aad2ece4251face9d95c106dd667abcd6`, immutable tag `edition-3-ch05-r1`. No implementation source was inspected, changed, or tested by this author. Code acceptance is taken from the authoritative reports and acceptance receipts, not independently re-established here. The existing amendment at `4f1701a76f57ec57d362d4f1a30264f30b65ddf8` and its completed formal text review remain intact.

## Reading and scope

Read the complete current `AGENTS.md`, `.agents/skills/edition-3-author/SKILL.md`, `docs/edition-3-execution.md`, `docs/edition-3-readiness.md`, `docs/edition-3-student.md`, `docs/edition-3-carryover.md`, `docs/agentic-codebook-workflow.md`, `docs/edition-2-carryover.md`, `book/voice.md`, `book/chapter-writing-procedure.md`, and the full current Chapter 5 before editing. Large/capped reads were recovered in bounded ranges. The initial combined shared-procedure read and combined craft/chapter read were not treated as sufficient by themselves.

Read the full authoritative chapter reports: `ch05-student-review.md`, `ch05-student-review-r1.md`, `ch05-code-review.md`, `ch05-code-review-r1.md`, `ch05-wakeup-student.md`, `ch05-wakeup-code-review.md`, `ch05-accepted.json`, `ch05-wakeup-candidate.json`, and `ch05-execution.md`, all under `docs/edition-3-notes/`. Also read the complete `job-wakeup-amendment.md`, `job-wakeup-amendment-author.md`, and `job-wakeup-amendment-review.md` in that directory. Chapter 4 lines 270–440 were read as context for job observation, cursor, process, and `send_input` behavior; Chapter 4 was not edited or reviewed anew.

The exact input-document hashes are retained in [`inputs.sha256`](ch05-author-evidence/inputs.sha256). The coding skill was not needed: this assignment stayed in prose, did not inspect implementation source, and ran only the required existing prose linter. No child agents or reviewers were commissioned.

## Changed teaching and demonstrated cause

### 1. A text-only HTTP reply is still a pending-hint boundary

Candidate lines 594–599, in §5.6 immediately after the existing explanation of HTTP-boundary hint delivery, add the missing consequence: preserve the assistant text, apply the pending hint, and issue another request in the same turn rather than returning while that hint is still waiting for a later `Ask`.

Evidence:

- [`ch05-code-review.md`](ch05-code-review.md), R1, describes the controlled text-only HTTP counterexample: a hint admitted while the first request was in flight did not appear in a second request in that turn. It survived until a later turn.
- [`ch05-student-review-r1.md`](ch05-student-review-r1.md) records the correction: drain pending hints before text-only completion, preserve the assistant text, and continue the same actor turn.
- [`ch05-code-review-r1.md`](ch05-code-review-r1.md) accepts the corrected same-turn behavior and its focused regression evidence.

The pre-existing prose already explained timely hint admission and draining before HTTP calls. It did not explicitly connect a text-only response to the decision whether that turn can finish. The added paragraph names the actual failed boundary and its repair without adding a new API, setting, exercise requirement, or checklist item. Its last sentence explains why admission alone did not help the model.

### 2. Output readiness does not finish every observation

Candidate lines 655–661, also in §5.6, distinguish a notification that prompts a condition recheck from the condition that returns a tool observation. The added paragraph retains the demonstrated PTY failure: treating every new byte as a return condition caused `send_input` to report echo from a still-running process. It separates the already-taught immediate read of unseen bytes on resumed `wait_for_job` from `send_input`'s pattern/completion/callback-delay behavior unless an actual hint or question wakes it.

Evidence:

- [`ch05-wakeup-student.md`](ch05-wakeup-student.md), “Attempts, assistance and limits of credit,” preserves the inherited Chapter 2 PTY failure at intermediate commit `ccf82d8`: `running=true` and cursor 11 after echo instead of the eventual completed result. It records the repair that kept `send_input`'s established observation policy while allowing resumed `wait_for_job` to read available output immediately.
- [`ch05-wakeup-code-review.md`](ch05-wakeup-code-review.md), “Output readiness, cursor, timer and process compatibility” and “Valid original-production sensitivity,” independently retains that actual bad intermediate and distinguishes the two observation policies. The reviewer reports original-production sensitivity and accepted final behavior; this author did not rerun either.
- The accepted amendment already teaches buffered output, cursor preservation, durable notifications, and waking observations without stopping work. Its exact sentences and Bill's story remain unchanged.

This addition qualifies a plausible overgeneralization at the new boundary. It does not change the amendment's wake-up story or contract and does not add a new requirement for Chapter 4.

## Deliberate no-change decisions

- **R2, queued caller cancellation and durable replies:** §5.4 already gives the buffered one-result private reply channel and explains that an abandoned caller must not strand the actor. The public-method discussion and §5.8 already connect caller cancellation to queued admission and final settlement. The exercise already covers queued callers abandoning their waits. Another paragraph would repeat explicit teaching rather than clarify an omission.
- **R3, stop all managed processes before joining noncancellable handlers:** §5.7 already states the ordering directly and explains why a blocked Go handler must not leave a killable shell alive. It also retains the question-pause/cancel interaction. No new wording is needed.
- No new framework, headings, exercise requirements, future-chapter preparation, or style cleanup was introduced. Existing explanations, useful failures, examples, pseudocode, tables, stories, and evidence qualifications were not removed to satisfy lint.

## Preservation checks

[`preservation-checks.json`](ch05-author-evidence/preservation-checks.json) records passing checks at that earlier check point that:

- removing exactly the two new paragraphs restores the exact baseline bytes;
- all original chapter text therefore survives unchanged, including the already-reviewed amendment and Bill's preface;
- all headings, fenced code blocks, and table lines are identical;
- repository candidate, retained candidate, and lint candidate match exactly;
- the retained baseline and exact supplied baseline lint input match;
- protected input-document hashes, Chapter 4, `book/preface.md`, and `book/edition-3/preface.md` are unchanged;
- before/after status adds only the Chapter 5 tracked modification, and the staged diff remains empty;
- saved evidence at that check is UTF-8.

There were pre-existing tracked changes in `cmd/grade-edition3/main.go`, `cmd/grade-edition3/main_test.go`, `docs/edition-3-execution.md`, and `docs/edition-3-student.md`, plus extensive existing untracked work. They were left alone. Status snapshots are retained, not presented as a clean-checkout claim. The evidence directory already existed at the “before” snapshot because this author had just created it. The report itself was written after the “after” snapshot. The protected input hashes are not a whole-repository integrity inventory, and no independent hash comparison of every unrelated untracked artifact is claimed. A later check found the execution record had gained a concurrent parent checkpoint note, confirmed by the parent in message `main-21` (retained in `parent-confirmation.txt`); its before/after hashes and the final qualified comparison are recorded in `final-initial-failed.json` and `final-checks.txt`. The earlier preservation result is retained as a point-in-time observation, not misrepresented as final equality for that file.

## Required prose lint and whitespace

Both lint runs used the same `chapter-05.md` basename and ran sequentially from `/Users/bill/projects/ensemble` with `GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off`, preventing toolchain/dependency fetching. Exact commands, full stdout, full stderr, and individual exit codes are saved separately.

```text
go run ./cmd/lintprose -v /tmp/cadence-ch05-poststudent-baseline/chapter-05.md
go run ./cmd/lintprose -v /tmp/cadence-ch05-poststudent-author-final/chapter-05.md
```

| Measurement | Baseline | Candidate |
|---|---:|---:|
| Exit code | 0 | 0 |
| Words | 7,854 | 7,988 |
| Paragraphs used for short-closer threshold | 143 | 145 |
| Negation forms | 38 | 38 |
| Short closers | 35 | 35 |
| Pointer closers | 1 | 1 |
| Em dashes | 0 | 0 |
| LLM crutches | 0 | 0 |
| Heading shape run | 1 | 1 |
| Person gap after line 0 | 4,083 words | 4,217 words |

Both retain the same soft-warning categories: words exceed the 4,000–7,500 target range, negation forms exceed 15 (one per 500 words), and the person gap exceeds 1,200 words. The additions increase the word count and the existing person gap by 134; no new warning category appears. The saved stdout reports these warnings and successful checks, not a separate printed PASS verdict. Exit 0 is not a verdict on teaching quality. No teaching was deleted to remove these warnings. Both complete stdout streams were read; both stderr files are empty. The stdout byte counts are 2,955 (baseline) and 2,959 (candidate).

`git --no-pager diff --check -- book/edition-3/chapter-05.md` exited 0 with empty stdout/stderr. Its exact outcomes are retained as `whitespace.*.txt`. Additional final UTF-8, report/chapter whitespace, scope, and inventory verification results are retained in `final-checks.txt` and the final submission.

## Failed checks, limitations, and work not run

- The first final cross-check caught inaccurate candidate SHA and lint figures in this report's draft. The draft had incorrectly stated 8,106→8,240 words, 156→158 paragraphs, 50 negations, a 32% negation ratio, a 1,082-word person gap, and thresholds not present in the saved lint. Those claims were replaced with the actual saved stdout figures above; unsupported semantic-break/code-line/heading-count rows were removed. The chapter and retained candidate had always matched. The same check detected the concurrent `ch05-execution.md` append. The failed exit 1 is preserved in `final-initial-failed.json`, with corrections explained in `validation-notes.txt`. This was a report/evidence correction before frozen submission, not another chapter revision or review round.
- An initial preservation assertion incorrectly expected the entire checkout's tracked diff to contain only Chapter 5. It failed because of the four pre-existing tracked edits listed above. The original failure is retained in [`preservation-initial-failed.json`](ch05-author-evidence/preservation-initial-failed.json). Python raised `AssertionError` and exited 1; subsequent read-only commands in that compound shell made its final shell exit 0. The corrected before/after scope check passed. No files were reverted or cleaned up to make it pass.
- An early `ls -ld` probe for the not-yet-created report and evidence directory returned 1 with stderr suppressed. This was an absence check, not a failed validation.
- Capped reading outputs were recovered; no capped lint output or unobserved running job is being represented as complete. See [`validation-notes.txt`](ch05-author-evidence/validation-notes.txt).
- No paid/provider calls, installation, deployment, publication, push, staging, commit, tag, child agents, or formal text review occurred in this pass.
- No implementation tests, race tests, graders, student builds, source builds other than the required prose lint, or full-book lint were run. Existing code/grader findings are attributed to their reports, not claimed as fresh author measurements.
- The accepted source, graders, Chapter 4 and other chapters, guide, reference, original manuscripts, and prior evidence/reports were not edited by this author.

The chapter is frozen at the candidate hash above. The parent owns exactly one formal text-review assignment and any request for the one permitted author revision. This author has not started a review loop, requested optional feedback, or advanced to another chapter.
