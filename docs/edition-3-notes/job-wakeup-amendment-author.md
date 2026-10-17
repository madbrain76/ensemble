# Foreground job-observation wakeup amendment: author report

Status: clarified text draft frozen for root's one independent formal TEXT review. This is not acceptance, implementation completion, or the normal post-student author phase. No reviewer was spawned. The root retains the later student extension and independent code review, and at most one author revision after this amendment's formal text review.

## Authorization and scope

Bill explicitly authorized this bounded contract amendment, including necessary Chapter 4 clarification. That authorization overrides the ordinary post-student-only timing only for this task. The further root instruction supplied Bill's observed join/question failure and requested a short grounded account of it. A subsequent explicit user clarification, received before formal review completed, authorized the narrow undispatched-batch-tail rule and regression. This update is not a response to reviewer findings or a normal post-student author pass; the root still owns one formal review.

Edited manuscripts:

- `book/edition-3/chapter-04.md`
- `book/edition-3/chapter-05.md`

New author material is this report and bounded text-validation receipts under `docs/edition-3-notes/job-wakeup-amendment-evidence/`. No source, grader, guidance, frozen answer, old accepted receipt, or student repository was edited. No commit, stage, push, installation, provider call, child agent, runtime experiment, or grader run was performed. The existing root prose linter was run, not changed.

The worktree already contains unrelated modifications and untracked records. They were not cleaned up, adopted, or credited to this amendment. The exact baseline snapshots, not a claim that the whole worktree is clean, define the text delta.

Bill's preface remains verbatim. Both `book/preface.md` and `book/edition-3/preface.md` are byte-equal to their HEAD blobs, SHA-256 `915cf8c531169c5aacc4cb81bd743eeba8a181d20bc779d936ac5f7ab9998d96`. Chapter 4's opening before §4.1 and Chapter 5's opening before the TL;DR also remain byte-equal to their captured baselines.

## Required context and historical boundaries

Read the complete current `AGENTS.md`, `.agents/skills/edition-3-author/SKILL.md`, `docs/edition-3-execution.md`, `docs/agentic-codebook-workflow.md`, `docs/edition-2-carryover.md`, `book/voice.md`, `book/chapter-writing-procedure.md`, and both assigned chapters; recovered capped spans using smaller reads. Read `docs/edition-3-notes/ch04-completion.md` and `docs/edition-3-notes/ch05-execution.md` for the applicable acceptance and correction boundaries. No private conversation logs were used.

Chapter 4's existing accepted job contract remains its historical contract. Chapter 5's separately frozen three-defect correction and its recorded tests/review remain separate evidence. Neither those results, the existing `not-deaf` marker, nor the historical demonstrations are relabeled as proving this new wakeup obligation. The amendment does not invalidate the old acceptance or claim any runtime has been repaired. Chapter 5 from §5.11 onward, including its historical verification account, is unchanged.

## Existing requirements clarified, not replaced

Chapter 4 already supplies the uniform job mechanism, `tool_limits`, `wait_for_job`, `send_input`, `kill_job`, bounded reports, retained output and the unseen-output cursor. The existing default table, literal-next-call rule, precedence, second-setter consumption/reset paragraph and established tool names remain intact.

The bounded additions make these consequences explicit:

- Quick `read_file`, long `run_command` and application-registered ordinary tools all enter the job path. Later `send_message`, `join_agents` and MCP adapters inherit that seam when introduced; this does not add their frameworks to the chapter.
- Controls can identify an existing target but allocate no new job handle. A shared result shape does not give files or joins stdin; `send_input` actually delivers input only to supporting jobs and reports unsupported input honestly.
- A callback delay is maximum observation silence, not a timeout on underlying work. Initial observations and subsequent inspections use the shared limits. Finished results and available unseen output can be returned without starting another handler.
- A setter followed by its target within one model response is meaningful: dispatch consumes pending settings in response order, including controls, argument errors and tool failures. Call-specific overrides win only for that call. This orders dispatch, not job completion, and introduces no parallel framework.

No defaults, unlimited-wait semantics, new API names, or early Chapter 4 actor-hint requirement were added.

## New Chapter 5 obligation

The TL;DR, inbound-event discussion, existing actor mechanism and exercise now require prompt foreground-observation wakeup for a human hint or addressed incoming question, including an initial job observation and a later `wait_for_job` with a chosen 3,600-second callback delay.

The observable contract is:

1. Admit and preserve the actual addressed event, including sender/recipient/request-turn correlation metadata. It is not a lossy observer notification and does not silently become a new blocking `Ask`/`Submit` turn.
2. End the current observation and cancel its timer, while the underlying work and stable job handle continue.
3. Resolve the outstanding tool call once under its original call id with bounded unseen output and the outcome actually observed. A running case may say “observation interrupted; job still running.” Never abandon the call, invent success/failure, resend the work, or force the actor turn to end.
4. Carry the event and correlated observation return into the next model request. Preserve the existing reliable private replies and save-before-reply rule.
5. Resume through `wait_for_job` on the same id, immediately returning available output or a retained terminal result. Preserve output retention and cursor behavior.
6. Completion/message races retain both facts, produce one return per observation call, and retire old timers so they cannot answer a later wait.

The concrete `go test ./...`/long-wait example explains the human steering case. An application-registered held join-like handler supplies a real seam for offline validation without adding child-agent/MCP frameworks. Future named adapters inherit it. Interrupt remains a turn-ending operation; a hint during HTTP still cannot change bytes already sent.

The added two-paragraph story attributes the observed parent stuck in `join_agents`, unanswered student questions, timeouts and `message_fleet` attempt to Bill's report. It expressly distinguishes “perhaps in desperation” from established motive and does not pretend to possess a timing or causal wire trace. The explanation identifies the dependency stall: students wait for an answer while the parent waits for completion. Extending a timeout does not solve that cycle. The parent must receive the question, answer it, and be able to resume the same join id. `message_fleet` and child spawning are not new exercise dependencies.

## Exercise and evidence boundary

The new exercise paragraphs require brief offline regressions through real Agent/job interfaces with fake boundaries: held command and registered handler; initial observation and later `wait_for_job`; long callback plus hint and question; bounded prompt delivery to the next fake-provider request; continuing same-id work; one handler invocation; available output and finished-job reads; cap/retention/cursor checks; message/completion orderings and a close race; reliable private replies; short controlled-timer checks for ghost returns.

The text does not claim these regressions have run, that current graders detect this behavior, or that another live campaign is needed. An hour-long configured delay is not a requirement to let an hour elapse. These are requirements for the later implementation/review phase, not implementation evidence from this author task.

## Exact manuscript hashes

| Path | Before SHA-256 | Frozen draft SHA-256 |
|---|---|---|
| `book/edition-3/chapter-04.md` | `ae7f6fbe1048400f857e8fb3eeefbe513d0a9bab492e1eff1d13a10ff266f37f` | `030c9a7ceb6c18f56596e840bb9df404fa418ea5f22c77bc7712b7ec1e6d8c55` |
| `book/edition-3/chapter-05.md` | `5d3db35e29da8a7fa8b343c01c88b3040339d75c347e845f434f124f3b17b0f2` | `4f4c2595ef65a08511d5cc68553f853feaf854e4cb0fe7d4d9094d4ae7f5a344` |

The exact before/after snapshots, `before.sha256`, `after.sha256` and `amendment.diff` are in the evidence directory. The manuscript diff is 28 additions/5 deletions for Chapter 4 and 132 additions/9 deletions for Chapter 5. Existing rationale and stories were retained rather than rewritten wholesale.

## Text validation results

Commands used the root's existing tool with the same chapter basenames for baseline and final:

```sh
go run ./cmd/lintprose docs/edition-3-notes/job-wakeup-amendment-evidence/baseline/chapter-04.md docs/edition-3-notes/job-wakeup-amendment-evidence/baseline/chapter-05.md
go run ./cmd/lintprose docs/edition-3-notes/job-wakeup-amendment-evidence/final/chapter-04.md docs/edition-3-notes/job-wakeup-amendment-evidence/final/chapter-05.md
git diff --check -- book/edition-3/chapter-04.md book/edition-3/chapter-05.md
```

Both lint runs exited 0, with **zero hard failures**, not zero warnings. Existing linter source was inspected to distinguish hard failures from soft warnings. Final whitespace check exited 0 with no diagnostics.

| Soft metric | Chapter 4 baseline → final | Chapter 5 baseline → final |
|---|---|---|
| Words (target 4,000–7,500) | 8,369 → 8,603; warning remains | 6,613 → 7,854; new soft warning |
| Negation forms | 27 → 29; warning remains | 29 → 38; warning remains |
| Longest person gap | 5,381 → 5,482; warning remains | 4,431 → 4,083; warning remains |

Other lint metrics remain within their respective limits. Chapter 5's new soft word warning is disclosed, not treated as a hard failure or silently waived. The added story and mechanism carry the requested contract; no unrelated cuts were made merely to target a count. Formal review can judge load-bearing prose within the root's reserved review/revision bound.

`text-validation.json` records successful bounded text assertions for unchanged openings/default table/second-setter paragraph/historical ending, required new contract anchors, preface identity and exact final snapshot identity. These checks establish text facts only, not semantic completeness or runtime correctness.

A guarded intermediate Chapter 5 text replacement rejected a mismatched anchor before writing that chapter; it was retried against the actual current passages. An interactive validation pager was terminated and later diffs disabled paging. Neither event affected underlying application work. No failed runtime check was converted into a passing claim. The file-writing tool's report was not visible to the repository shell; this report was therefore written through the shell and checked there before submission.

## Direct user clarification: undispatched batch tail

This reopened pass changed only Chapter 5's wakeup mechanism (one paragraph) and its existing offline exercise (one regression bullet). Chapter 4 and all other manuscript text remain unchanged from the submitted draft. It adds 26 lines/268 words to that draft; no API or implementation framework was added.

The completed/dispatched/undispatched distinction is explicit. Completed calls retain real results. Started jobs retain their observed state and handle; only observation wakes. Remaining undispatched calls receive exact call-id-correlated `interrupted; not executed` errors without handlers, effects or new job ids, and are not automatically replayed. Original assistant calls/signatures remain intact, every call is paired before the next HTTP request, and a hint/question is not held behind another queued long wait. The actor's serial dispatch boundary settles races: started work cannot be falsely labeled unexecuted. Pending one-shot limits are consumed/reset by the literal next call's error; a setter that never executes cannot arm settings. Fresh model calls may run after wakeup while the original job continues; waiting for completion is not required to answer a student.

The added regression covers that mixed batch, side-effect suppression, pairing/signature preservation, no replay or new job ids for the tail, fresh calls during continuing work, dispatch/setter races and limits reset. These tests are specified, not run in this author task.

Initial submission evidence is preserved byte-for-byte under `job-wakeup-amendment-evidence/superseded/initial-draft/`, including its full report, final snapshots, original manifest and receipts. The archived manifest's SHA-256 remains `60953e0b264e15e3565c958bea0c8a5c962ded53c9d97ac155b9f63e6c4e4d41`; all 14 entries were checked against their archived equivalents. Its Chapter 5 snapshot remains `0aa8a2060025b59c538a8cb6f0a389428ae2b6673382fb7433f1496f16f0c913` and its report remains `7f5bfc71600d255a428519a169e03001590dcf21ecebf1b076c1f4a655c04420`.

`queued-tail-clarification.diff` records precisely the initial-submission-to-current change. Current final snapshots, full baseline-to-current diff, after hashes, lint receipt, text assertions and manifest have been refreshed, not passed off as the initial evidence. The original baseline remains unchanged. Final same-basename prose lint again exits 0 with zero hard failures, retaining the disclosed soft warnings; the updated metrics are above. Whitespace checking also exits 0. No paid calls, code implementation, children, commits or extra review rounds were introduced.

## Frozen handoff

The chapter bytes above are frozen for the root's independent formal TEXT review. All authored paths and their hashes, including this report and the bounded receipts, are listed in `job-wakeup-amendment-evidence/authored-files.sha256` (the manifest excludes itself). No acceptance claim is made. Any material unresolved issue belongs to the root; no extra reviewer or author loop was started.
